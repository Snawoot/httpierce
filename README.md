# httpierce

Experimental shadowsocks plugin purposed to bypass captive portals of some mobile ISPs.

Based on simple-tls plugin code by @IrineSistiana.

## Installation


#### Pre-built binaries

Pre-built binaries available on [releases](https://github.com/Snawoot/httpierce/releases/latest) page.

#### From source

Alternatively, you may install httpierce from source:

```
go get github.com/Snawoot/httpierce
```

## Usage

### As a shadowsocks plugin

Just specify path to binary as plugin argument. On server pass `server` as plugin options


### As a standalone executable

```
$ httpierce -h
Usage of bin/httpierce:
  -V	VPN mode. Used by shadowsocks on Android
  -bind string
    	listen address
  -dst string
    	target address
  -server
    	server-side mode
  -timeout duration
    	connect timeout (default 10s)
```
