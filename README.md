netblock helps you turn off distracting websites, by redirecting them to localhost.

Example: while doing something important, you suddenly feel an urge to visit test.com - don't worry, because you
have blocked test.com, you'll see an ugly error message in the browser instead. 

It blocks websites by listing the blocked domains in /etc/hosts file.

Install
-------

	go install github.com/tanel/netblock

Block web page
--------------

	netblock add test.com

Unblock web page
----------------

	netblock remove test.com

List web pages
--------------

	netblock list
