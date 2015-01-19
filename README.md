# chinchilla
A simple Golang service that takes pictures from dropbox and puts them on facebook.

## NOTE: this repo has submodules
To check it out properly:

```bash
$ git clone --recursive https://github.com/mrgamer/chinchilla.git
```

## Example config.toml (REQUIRED)
Chinchilla requires yourself to register a dropbox app (more details later), and a facebook app.  
The service is tought to be backed by a mongoDB server, so also a mongoDB instance is **required**.

```toml
[dropbox]
key = "YOUR APP KEY HERE"
secret = "YOUR APP SECRET HERE"

[mongo]
addrs = ["127.0.0.1"]
database = "database-name"
```
