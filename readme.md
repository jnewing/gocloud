
# gocloud

This aims to obe a quick way to setup and then keep your dyanmic ip updated on
Cloudflare through the use of their API.

I wrote this as kind of a intro to Go and becuase I needed a way to keep my Cloudflare records
up to day with my dynamic ip. Homelab stuff.

Usage:
--------------------------------------------------
`gocloud [--debug]`

Flags:

--debug
Writes a LOT of spammy messages to stdout, this does help with finding issues

Config
--------------------------------------------------
This does require a config.yml file to be placed in the same directory as the main exeacutable. I've included
an example file, config.example.yml you can make a copy of this and remove the .example part from the name.const

Be sure to read the comments in there it's fairly self explnitory. The gist of it is you need your Cloudflare email, and
Global API Key as well some zones you want to make sure are up to date.

How
--------------------------------------------------
How does it work? Well it loops through your list of zones (defined in config.yml) and calls the Cloudflare API to then fetch a
list of all "A" records for this DNS Zone.

We then iterate over these records looking for a matching name (again in your config.yml) if a name is matcched and IP set in
this type "A" record is different than your external IP (or IP set in config.yml) we set it to match. If the name is not matched OR
the IP already matches we don't do anything.

