netblock helps you turn off distracting websites, by redirecting them to localhost.

Example: while doing something important, you suddenly feel an urge to shop on Amazon. Don't worry, because you
have blocked amazon.com, you'll see an ugly error message in the browser instead. 

It blocks websites by listing the blocked domains in /etc/hosts file.

Install
-------

	go install github.com/tanel/netblock

Block web page
--------------

	netblock add amazon.com

Unblock web page
----------------

	netblock remove amazon.com

List web pages
--------------

	netblock list
