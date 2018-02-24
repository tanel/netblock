netblock helps you turn off distracting websites, by redirecting them to localhost (or any other host you choose).

Example: while doing something important, you suddenly feel an urge to shop on Amazon. Don't worry, because you
have blocked amazon.com, you'll see an ugly error message in the browser instead. 

It blocks websites by listing the blocked domains in /etc/hosts file.

Install
-------

	go install github.com/tanel/netblock

Block web page
--------------

	netblock amazon.com

Unblock web page
----------------

	netblock allow amazon.com

List web pages
--------------

	netblock
