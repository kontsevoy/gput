### gput - CLI client for Rackspace Cloud Files

Actually this is not a full client and does not support all features of [Cloud Files](http://www.rackspace.com/cloud/files), a similar service to Amazon S3. Its usability is optimized for quickly uploading files into containers. It can take most parameters out of `~/.gput.ini` file for quickness.

#### Installation

Run: 
```
$ go install github.com/kontsevoy/gput
```

For convenience, here are latest binaries for [64-bit Linux](http://i.kontsevoy/gput/gput-linux-x86_64) or [Mac OSX](http://i.kontsevoy/gput/gput-darwin) if you cannot compile from source.

#### Why?

I needed a real world project to learn Go. I also needed a quicker way to share files, or sometimes just the std output of something, with other people. So I wrote `gput`. Here's how I quickly upload get a public URL for a file, which will only exist for 1 hour:

```
$ gput -ttl 600 build.log 
201 Created
http://4ae40c0c72b6359660b2-9a0aa15bf6dbebe33d1d027dbc9d76cd.r73.cf1.rackcdn.com/build.log
https://b2d1a1f1f5e5744aef4b-9a0aa15bf6dbebe33d1d027dbc9d76cd.ssl.cf1.rackcdn.com/build.log
```
Those URLs can be quickly copy&pasted into a chat window or sent via email. They will stop working after 1 hour (600 seconds) automatically.

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
	gput -region DFW put example.txt folder/name.txt
	gput -ttl 600 put example.txt
	gput -region DFW list
	gput -region DFW -container public-container list
	gput -region DFW list public-container
	gput -region DFW -container public-container delete example.txt
	gput delete folder/example.txt
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
