### gput - CLI client for Rackspace Cloud Files

Actually this is not a true client, its usability is optimized for quickly uploading files into containers. 
It can take most parameters out of `~/.gput.ini` file for quickness and convenience.

#### Usage & Examples

`gput` uses `put` command by default unless something else (like `list` or `delete`) is secified. If a 
target container is CDN-enabled, it will print public URLs for HTTP and HTTPS.

```
Usage:
	gput <options> <command> <file>

Commands:
	put   :	Upload file into a container. This command executes by default.
	list  :	Llist containers or files within a container
	delete:	Delete file in a container
	gen   :	Generate a template config file

Options:
	-config : path to a config file
	-container : Cloud Files container name
	-key : Rackspace API key
	-region : Cloud Files region to use, like DFW, ORD, etc
	-ttl : Time to live in seconds. 0 (default) means forever
	-user : Rackspace API username

Examples:
	gput example.txt   (assumes the region and container are in a config file)
	gput -region DFW put example.txt new-name.txt
	gput -ttl 600 put example.txt
	gput -region DFW list
	gput -region DFW -container public-container list
	gput -region DFW list public-container
	gput -region DFW -container public-container delete example.txt
```

#### Config File

For convenience, `gput` assumes `put` as the default command. Most other options can be stored in a `~/.gput.ini` file:
```ini
[Auth]
key=329800735698586629295641978511
username=konzsevoy

[Cloud Files]
container=ev-public
region=dfw
```

#### Misc

This is my first Golang app. Let me know if there're glaring issues in my Go code. Thank you.
