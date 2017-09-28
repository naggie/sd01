# Copyright (c) 2017 Callan Bryant
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.
"""
A bare minimal service discovery system.

Divulges IPv4 addresses and ports of services running on the same subnet with
the same "service id" ascii string in Announce mode.

By design, sd01 does not support service descriptions. It is intended that the
device will be interrogated by the discoverer post-discovery via another
mechanism such as a subsequent API call.*


Definitions:

  * Service: Something listening on a port running on a host.
  * Service class: An ascii identifier corresponding to the service/project
    name. A version number could be appended. Max 23 characters.
  * Service port: The port the service is listening on


Usage:
    # on devices that you wish to discover
    Announcer('my_project_name',service_port=80).start()

    # on machine that must discover services
    d = Discoverer('my_project_name')
    d.start()
    services = d.get_services(wait=True)

    # ...at any time, preferably after wait
    services = d.get_services()

`get_services` will only return services that are actively Announcing.

sd01 works using a UDP broadcast of a magic string on port 17823 every 5
seconds.

"""
# TODO IPv6 (multicast based) support
# example https://svn.python.org/projects/python/trunk/Demo/sockets/mcast.py

from socket import socket, AF_INET, SOCK_DGRAM, SOL_SOCKET, SO_BROADCAST
from threading import Thread, Lock
from logging import getLogger
from time import sleep
from functools import wraps

try:
    # Python 3.x
    from time import monotonic as time
except ImportError:
    # Python 2.x, not monotonic
    from time import time

log = getLogger(__name__)

# deterministic packet size regardless of port. Service ID max 25 chars --
# packet length is 32 bytes max to keep broadcast traffic low.
PACKET_FORMAT = 'sd01{service_class:.23}{service_port:0>5}'

# May be a problem with thousands of devices. A good compromise IMO -- 5
# seconds is an acceptable wait IMO. Results in 6.4 bytes per second per
# service.
INTERVAL = 5

PORT = 17823


def forever_IOError(fn):
    @wraps(fn)
    def _fn(*args, **kwargs):
        while True:
            try:
                return fn(*args, **kwargs)
            except IOError as e:
                # particularly e.errno == ENETUNREACH
                log.exception('Caught IOError, re-attempting.')
                sleep(5)

    return _fn


class Announcer(Thread):
    daemon = True

    def __init__(self, service_class, service_port):
        super(Announcer, self).__init__()
        service_class.encode('ascii')  # validate
        self.service_class = service_class
        self.service_port = int(service_port)

        if self.service_port < 0 or self.service_port > 65535:
            raise ValueError('Port number out of legal range')

    @forever_IOError
    def run(self):
        # create UDP socket
        s = socket(AF_INET, SOCK_DGRAM)
        s.bind(('', 0))
        s.setsockopt(SOL_SOCKET, SO_BROADCAST, 1)

        magic = PACKET_FORMAT.format(
            service_class=self.service_class,
            service_port=self.service_port,
        ).encode('ascii')

        while True:
            log.debug('Announcing on port %s with magic %s',
                      PORT, magic)
            s.sendto(magic, ('<broadcast>', PORT))
            sleep(INTERVAL)


class Discoverer(Thread):
    daemon = True

    def __init__(self, service_class):
        super(Discoverer, self).__init__()
        service_class.encode('ascii')  # validate
        self.service_class = service_class

        # map of (host,port) -> timestamp of last announcement
        self.services = dict()

        self.lock = Lock()
        self.running = False

    @forever_IOError
    def run(self):
        # create UDP socket
        s = socket(AF_INET, SOCK_DGRAM)
        s.bind(('', PORT))

        self.running = True

        magic = PACKET_FORMAT.format(
            service_class=self.service_class,
            service_port=0,
        ).encode('ascii')

        while True:
            data, addr = s.recvfrom(len(magic) + 5)
            if data.startswith(magic[:-5]):
                host = addr[0]

                try:
                    data = data.decode('ascii')
                except ValueError:
                    log.warn('Received invalid sd01 packet: non-ascii characters')
                    continue

                # no whitespace or decimals
                if not data[-5:].isdigit():
                    log.warn('Received invalid sd01 packet: invalid port number. Must be 5 digit, zero padded.')
                    continue

                port = int(data[-5:])

                if port < 0 or port > 65535:
                    log.warn(
                        'Received invalid sd01 packet: port number out of legal range')
                    continue

                log.debug('Discovered %s on port %s', host, port)

                with self.lock:
                    self.services[(host, port)] = time()

    def get_services(self, wait=False):
        '''Returns a list of tuples (host,port) for active services'''
        if not self.running:
            raise RuntimeError(
                'You must call start() first to start listening for announcements')

        if wait:
            sleep(INTERVAL + 1)

        min_ts = time() - INTERVAL * 2

        with self.lock:
            return [h for h, ts in self.services.items() if ts > min_ts]
