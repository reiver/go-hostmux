/*
Package hostmux provides a host oriented ("middleware") HTTP handler,
which can hand-off to other HTTP handler for each host.

Basic Usage

	hmux := hostmux.New()
	
	hmux.Host(handler1, "one.com")
	hmux.Host(handler2, "two.net", "www.two.net")
	hmux.Host(handler3, "three.org", "www.three.org", "web.three.org")

*/
package hostmux
