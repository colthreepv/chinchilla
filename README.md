# chinchilla
A simple Go service that takes pictures from dropbox and puts them on facebook.  
<a href="https://www.pinterest.com/pin/304837468499947408/"><img src="https://s-media-cache-ak0.pinimg.com/736x/3b/3e/ae/3b3eae66e5a656c7eceaf3ef56414a6d.jpg" height="500px"></a>

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
