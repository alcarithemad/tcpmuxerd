TCPMUX:
===
A Practical Defense Against TCP Port Number Exhaustion
---

When TCP was originally being devised, the Internet was small, and its uses (and therefore protocols) were few. It seemed like 65,536 ports would be enough to last forever. However, as time went on, and we assigned more and more of the TCP port number space, people began to realize that we needed a way to limit the number of ports involved in running a single daemon. While several stopgap measures have been proposed and implementated, only one true solution has been found: RFC 1078.
When TCP was new, it was not uncommon for protocols to use huge numbers of ports. Early adopters of SSL assigned a second port for SSL connections; examples of this include HTTP (80, 443), POP3 (110, 995), IMAP (143, 993), and SMTP (25, 465). An even more outrageous example, FTP, uses a control port (21), plus a random ephemeral port per active file transfer. On a typical modern server, capable of serving thousands of files simultaneously, that ties up thousands of ports for a single service.
Some protocol designers have replaced the need for bi-port protocols by designing protocol messages that restart the connection with an SSL negotation on the same port. One such example is SMTP's STARTTLS command. Other protocols have been supplanted by newer, more advanced protocols, capable of serving multiple purposes, like SSH, which can be used both as a remote shell (supplanting telnet) and for file transfers using rsync, scp, or sftp (supplanting FTP and its thousands of ports). However, people still continue to design and use inefficient protocols which require more than one port.
RFC 1078, on the other hand, specifies a real, viable, solution. Under RFC 1078, also known as TCPMUX, we only need one port number: 1, and through it we may connect to any service, even if we do not know its port number. By replacing port numbers with service names, we expand the number of services that can be run on a single server by an exponential amount. With TCPMUX, there is no danger of TCP port space exhaustion, and port-based firewalls are completely unnecessary!



BUILDING
===
Prerequisites:
---
Go 1

run `go build`

