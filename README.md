sd01 is a minimal service discovery protocol with strict implementation.

[![Build Status](https://travis-ci.org/naggie/sd01.svg?branch=master)](https://travis-ci.org/naggie/sd01)

Development status: Beta. Suitable for production but test coverage needs to
be improved.

sd01 is an alternative to [mDNS](https://en.wikipedia.org/wiki/Multicast_DNS)
which is implemented by Bonjour/Ahahi.


### Features
1. Advertisement of a IP + port for a given service by name
2. Timeout mechanism to remove services that have disappeared


### Non-features
1. No service metadata (instance name or capabilites, for example). This should
   be implemented out of band by the protocol the service uses itself.


#  Implementations

Current:

* Python 2/3 (reference implementation)
* Go


Planned:

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
particular named service. Note that device descriptions, names, versioning and
capabilities are deliberately left out; this should be implemented as an API on
the service itself.

sd01 works nicely with an RPC mechanism such as gRPC.

## Definitions

  * Host: a device or server running a service
  * Service: Something listening on a port running on a host.
  * Service name: Product name, max 55 characters. Do not encode version or
    capabilites.
  * Service port: The port the service is listening on

## Message

A host emits a sd01 message every 10 seconds. If an announcer has not
seen the sd01 message for 600 seconds, the host is considered non-existent.

```
sd01:[service_name]:[service_port]
```

Where, without brackets:

  1. The total message length is no more than 64 bytes (53 chars for service_name)
  2. The entire message is composed of ASCII characters only
  3. The message is prefixed with `sd01`
  4. The service port is an integer string


For example with a service name "DS Light controller" running on port 80:

```
sd01:DS light controller:80
```


Both the announcer and discoverer are aware of the `service_name` in advance.

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
  5. Valid sd01 messages from a different service_name are ignored
  5. Debug logging



# Use cases

## Multi Sensor
This sensor is stateless and should just work.

1. ESP8266 based sensor module connects to wifi, using DHCP
2. ESP8266 announces via sd01 `DS Multisensor` port `80`
3. Automation controller discovers a known multi sensor is listening on `192.168.59.32:80`, and knows the sensor speaks a specific HTTP JSON protocol
4. Automation controller queries the sensor, and recevies something like:
    ```
    {
       "id": "2ef8787",
       "temperature": 49,
       "humidity": 23
    }
    ```
5. The automation controller consumes the value, knowing the location of the sensor by looking up the ID in its database. It knows this sensor is located in the `Kitchen`

Note that the sensor is effectively stateless. It has its own preconfigured ID, which happens to be derrived from its MAC address.

