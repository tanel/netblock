# netblock

[![Build Status](https://travis-ci.org/tanel/netblock.svg?branch=master)](https://travis-ci.org/tanel/netblock) [![Go Report Card](https://goreportcard.com/badge/github.com/tanel/netblock)](https://goreportcard.com/report/github.com/tanel/netblock)

netblock helps you turn off distracting websites, by redirecting them to 0.0.0.0.

Example: while doing something important, you suddenly feel an urge to visit test.com - don't worry, because you
have blocked test.com, you'll see an ugly error message in the browser instead. 

It blocks websites by listing the blocked domains in /etc/hosts file.

Install
-------

	go install github.com/tanel/netblock

Block web page
--------------

	sudo netblock add test.com

Unblock web page
----------------

	sudo netblock remove test.com

List web pages
--------------

	netblock list
