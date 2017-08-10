# Copyright (c) 2017 Callan bryant
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

Divulges IPv4 addresses of hosts on the same subnet with the same "magic" string
in Announce mode.

By design, sd01 does not support service descriptions. It is intended that the
device will be interrogated by the discoverer post-discovery via another
mechanism, and operate on a canonical port.*


Usage:
    # on devices that you wish to discover
    Announcer('my_project_magic_string').start()

    # on machine that must discover hosts
    d = Discoverer('my_project_magic_string')
    d.start()
    hosts = d.get_hosts(wait=True)

    # ...at any time, preferably after wait
    hosts = d.get_hosts()

`get_hosts` will only return hosts that are actively Announcing.

sd01 works using a UDP broadcast of a magic string on an automatically chosen
port over 10000. A port can be specified when you have an os-level firewall enabled.



* If multiple services are required to run on one host, and therefore use
different ports the following mechanism is suggested:

    1. Services start (in any order) and attempt bind to a base port
    2. If the bind fails, increment port number by one and try again up to a limit of 10 (or so)
    3. The discoverer should try to bind to the base port and all ports following up to the limit of 10

This way, services don't have to be configured with individual ports explicitly.

"""
# TODO IPv6 (multicast based) support
# example https://svn.python.org/projects/python/trunk/Demo/sockets/mcast.py

from socket import socket, AF_INET, SOCK_DGRAM, SOL_SOCKET, SO_BROADCAST
from threading import Thread, Lock
from binascii import crc32
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



def forever_IOError(fn):
    @wraps(fn)
    def _fn(*args,**kwargs):
        while True:
            try:
                return fn(*args,**kwargs)
            except IOError as e:
                # particularly e.errno == ENETUNREACH
                log.exception('Caught IOError, re-attempting.')
                sleep(5)

    return _fn



class Base(Thread):
    daemon = True

    def __init__(self, magic, interval=5, port=None):
        super(Base, self).__init__()
        self.magic = str(magic).encode('ascii')
        self.interval = int(interval)

        if self.interval < 1:
            raise ValueError('Interval must be more than 1')

        # User may have selected a port due to firewall implications.
        # otherwise, pick a deterministic high port number
        if not port:
            # & 0xffffffff to match python3 behavior
            self.port = 10**4 + crc32(self.magic) & 0xFFFFFFFF % 10**4


class Announcer(Base):
    @forever_IOError
    def run(self):
        # create UDP socket
        s = socket(AF_INET, SOCK_DGRAM)
        s.bind(('', 0))
        s.setsockopt(SOL_SOCKET, SO_BROADCAST, 1)

        while True:
            log.debug('Announcing on port %s with magic %s',
                      self.port, self.magic)
            s.sendto(self.magic, ('<broadcast>', self.port))
            sleep(self.interval)


class Discoverer(Base):
    def __init__(self, *args, **kwargs):
        super(Discoverer, self).__init__(*args, **kwargs)

        # map of host -> timestamp of last announcement
        self.hosts = dict()

        self.lock = Lock()
        self.running = False

    @forever_IOError
    def run(self):
        # create UDP socket
        s = socket(AF_INET, SOCK_DGRAM)
        s.bind(('', self.port))

        self.running = True

        while True:
            data, addr = s.recvfrom(len(self.magic))
            if data == self.magic:
                host = addr[0]
                log.debug('Discovered %s with magic %s', host, self.magic)
                with self.lock:
                    self.hosts[host] = time()


    def get_hosts(self, wait=False):
        if not self.running:
            raise RuntimeError(
                'You must call start() first to start listening for announcements')

        if wait:
            sleep(self.interval + 1)

        min_ts = time() - self.interval * 2

        with self.lock:
            return [h for h, ts in self.hosts.items() if ts > min_ts]
