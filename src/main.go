/*
gocloud

This aims to obe a quick way to setup and then keep your dyanmic ip updated on
Cloudflare through the use of their API.

I wrote this as kind of a intro to Go and becuase I needed a way to keep my Cloudflare records
up to day with my dynamic ip. Homelab stuff.

Usage:
--------------------------------------------------
gocloud [--debug]

Flags:

--debug
Writes a LOT of spammy messages to stdout, this does help with finding issues

Config:
--------------------------------------------------
This does require a config.yml file to be placed in the same directory as the main exeacutable. I've included
an example file, config.example.yml you can make a copy of this and remove the .example part from the name.const

Be sure to read the comments in there it's fairly self explnitory. The gist of it is you need your Cloudflare email, and
Global API Key as well some zones you want to make sure are up to date.

How:
--------------------------------------------------
How does it work? Well it loops through your list of zones (defined in config.yml) and calls the Cloudflare API to then fetch a
list of all "A" records for this DNS Zone.

We then iterate over these records looking for a matching name (again in your config.yml) if a name is matcched and IP set in
this type "A" record is different than your external IP (or IP set in config.yml) we set it to match. If the name is not matched OR
the IP already matches we don't do anything.
*/

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"gocloud/src/cloudflare"
	"gocloud/src/config"
	"gocloud/src/helpers"
	"io"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Holds our global config values.
var conf config.Configs

// Global debug flag.
var debug *bool

// SetupLogger
// General logging setup, file, max size, age etc...
func SetupLogger() {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filepath.ToSlash("./logs/gocloud.log"),
		MaxSize:    5, // MB
		MaxBackups: 3,
		MaxAge:     30, // Days
		Compress:   true,
	}

	multiWriter := io.MultiWriter(lumberjackLogger)

	logFormatter := new(log.TextFormatter)
	logFormatter.TimestampFormat = time.RFC1123Z // or RFC3339
	logFormatter.FullTimestamp = true

	log.SetFormatter(logFormatter)
	log.SetLevel(log.DebugLevel)
	log.SetOutput(multiWriter)
}

// SetupConf
// Reads our config.yml file and places it into global conf varaible.
func SetupConf() {
	viper.SetConfigFile("config.yml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file. ", err)
	}

	verr := viper.Unmarshal(&conf)
	if verr != nil {
		log.Fatal("Unable to decode into struct. ", verr)
	}
}

// fetchDNSZone
// Fetches DNS Records from a zone using Cloudflare API
// https://developers.cloudflare.com/api/operations/dns-records-for-a-zone-list-dns-records
func fetchDNSZone(zone config.Zone) cloudflare.CFResp {
	resp := cloudflare.APICall("GET",
		conf.Cloudflare.API+"/zones/"+zone.ID+"/dns_records?type=A",
		conf.Cloudflare.Email,
		conf.Cloudflare.Key, nil)

	var CloudFlareResp cloudflare.CFResp
	json.Unmarshal(resp, &CloudFlareResp)

	return CloudFlareResp
}

// init
func init() {
	SetupLogger()
	SetupConf()

	debug = flag.Bool("debug", false, "Print debug output to stdout.")
}

// main
func main() {
	flag.Parse()

	if *debug {
		log.SetLevel(log.InfoLevel)
		log.SetOutput(os.Stdout)
	}

	extIP, err := helpers.FetchExtIP()

	if err != nil {
		log.Fatal("Unable to fetch external IP.")
	}

	for _, zone := range conf.Zones {
		log.Printf("[+] Zone ID: %s", zone.ID)

		if len(zone.IP) > 0 {
			*extIP = zone.IP
			log.Infof("[+] IP %s set, not using external IP.", zone.IP)
		}

		zoneDNS := fetchDNSZone(zone)

		for _, entry := range zoneDNS.DNSEntry {
			log.Infof("[+] Entry: %s", entry.ZoneName)

			for _, update := range zone.Update {
				if entry.Name == update {
					log.Infof("[+] Checking IP for Entry: %s", update)

					if entry.Content != *extIP {
						log.Infof("[+] Updating IP: %s -> %s", entry.Content, *extIP)

						payload, _ := json.Marshal(map[string]interface{}{
							"content": *extIP,
							"name":    entry.Name,
							"type":    entry.Type,
						})

						cloudflare.APICall("PATCH",
							conf.Cloudflare.API+"/zones/"+entry.ZoneID+"/dns_records/"+entry.ID,
							conf.Cloudflare.Email,
							conf.Cloudflare.Key,
							bytes.NewBuffer(payload),
						)
					} else {
						log.Infof("[+] IPs match.")
					}
				} else {
					log.Infof("[-] No matches for %s", update)
				}
			}
		}
	}
}
