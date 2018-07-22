sd01 is a minimal service discovery protocol with strict implementation.

Developement status: Beta. Suitable for production but test coverage needs to
be improved. #10 may result inn a breaking change to the protocol.


### Features
1. Advertisement of a service identified by a service type
2. Advertisement of port of service
3. Heartbeat/timeout mechanism to detect currently active services


### Non-features
1. No service metadata (instance name or capabilites, for example). This should
   be implemented out of band by the protocol the service uses itself.


#  Implementations
* Python 2/3 (reference implementation)
* Go


Planned languages:

* Embedded C
* ...

sd01 reveals the host (IP) and port for any similar services running.

sd01 works by sending a message specific to the service name and service port
over a UDP broadcast, port 17823. As a UDP broadcast is used, only hosts on the
same subnet can discover. The service host (IP) is taken from the UDP packet
source attribute.

An IPv6 and/or IPv4 multicast extension is being considered for version 2.


# Why?

There are many service discovery systems out there, most of them are much more
complicated than normally required. Many systems implement a custom UDP
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

A host emits a sd01 message every 5 seconds. If an announcer has not seen the
sd01 message for 10 seconds, the host is considered offline.

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

## Reference implementation

`sd01.py` is the reference implementation. Other implementations are expected
to follow the same design which consists of:

  1. Unit tests for valid and invalid messages. The implementation must be
     strict to avoid replying on undefined behaviour.
  2. Threaded announcer and discoverer
  3. `Discoverer.get_services`, a thread-safe method to return online
     (host,port) pairs.
  4. Invalid sd01 messages are handled as errors (although these errors should
     not crash the program, they may be logged)
  5. Valid sd01 messages from a different service_class are ignored
  5. Debug logging
