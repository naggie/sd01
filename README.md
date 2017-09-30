sd01 is a minimal service discovery protocol with strict implementation in the
following languages so far:

* Python 2/3


Planned languages:

* Embedded C
* Go
* ...


sd01 works by sending a message specific to the service name and service port
over a UDP broadcast, port 17823.




# Why?

There are many service discovery systems out there, most of them are much more
complicated than normally required. May systems implement a custom UDP
broadcast; sd01 is just a standard way to do so, cross-platform.


# Protocol

sd01 is deliberately minimal. It just exposes the IP address and the port of a
particular named service. Note that device descriptions, names and capabilities
are deliberately left out; this should be implemented as an API on the service
itself.

sd01 works nicely with an RPC mechanism such as gRPC.

## Definitions

  * Host: a device or server running a service
  * Service: Something listening on a port running on a host.
  * Service class: An ascii identifier corresponding to the service/project
    name. A version number could be appended. Max 23 characters.
  * Service port: The port the service is listening on

## Message

```
sd01[service_class][service_port]
```

Where, without brackets:

  1. The total message length is no more than 32 bytes (23 chars for service_class)
  2. The entire message is composed of ASCII characters only
  3. The message is prefixed with `sd01`
  4. The message the service port, 5 digit, zero padded


For example with a service class named "lightcontrollerv2" running on port 80:

```
sd01lightcontrollerv200080
```


Both the announcer and discoverer are aware of the `service_class` in advance. 
