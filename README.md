* Yet Another Golang DNS  Experimentation Project *

This is a simple dns server which takes a configuration file, and implements the low-level packet building of DNS responses.

Currently:
 - Serves IPv4 A Records from a JSON configuration file.
 - TTLs are hard coded to 300 seconds.
 - Not much logging.
 - Some test coverage

Some interesting things to note about DNS:
- Based on recomendations [here](https://stackoverflow.com/a/4083071) there will basically always be only one query per flag.
- The c0 0c pointer is explained
	[here](https://ask.wireshark.org/questions/50806/help-understanding-dns-packet-data).
	Since the likelyhood of getting multiple queries in a single packet is
	essentially 0, and the header of the request is so specifically defined, it
	is basically possible to hard code this value as a pointer to the name part
	of the query section of the packet.
