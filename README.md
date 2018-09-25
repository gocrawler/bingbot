# BingBot
[![Report Card](https://goreportcard.com/badge/github.com/AnikHasibul/bingbot)](https://goreportcard.com/report/github.com/AnikHasibul/BingBot#license)
[![Build Status](https://travis-ci.org/AnikHasibul/bingbot.svg?branch=master)](https://travis-ci.org/AnikHasibul/BingBot)

>A simple multi threaded bing search client for site collection.

** SAVE YOUR DORKS AT `dorks.txt`**

```txt
site:in php?id=
site:bd reg.php
site:gov.pk page.php?
```

Save your all dorks like this in `dorks.txt` file! and the command `go run bingbot.go`

It will update a count value in your terminal about *how many site has been collected till now* 

# For getting the collected sites:


> Each result set from each dork will be saved as the dork name such as `inurl%20php.txt`.

And this bot has a live http server for getting the sites from browser!

> For getting all sites as a clickable link visit **`http://localhost:1338/`**

> For exiting the bot you have to kill it or press `CTRL+C` or visit **`http://localhost:1338/exit`**

> For realoading or rescaning `dorks.txt` & `deny.txt` files without exiting server visit **`http://localhost:1338/reload`**

> For getting only domain name visit **`http://localhost:1338/domain`**


### Feel free to report any issue or hit the star!

> For any suggestion AnikHasibul@outlook.com

### TODO

*  Json Output
*  Colored Output
*  Smart web interface
