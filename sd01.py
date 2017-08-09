from socket import socket, AF_INET, SOCK_DGRAM, SOL_SOCKET, SO_BROADCAST
from threading import Thread,Lock
from binascii import crc32
from logging import getLogger
from time import sleep

log = getLogger(__name__)

try:
    # Python 3.x
    from time import monotonic as time
except ImportError:
    # Python 2.x, not monotonic
    from time import time

class Base(Thread):
    def __init__(self,magic,interval=5):
        super(Base,self).__init__()
        self.magic = str(magic).encode('ascii')
        self.interval = int(interval)

        if self.interval < 1:
            raise ValueError('Interval must be more than 1')

        # may remove this behavior due to firewall implications
        # picka deterministic high port number
        self.port = 10**4 + crc32(self.magic) % 10**4


class Announcer(Base):
    daemon = True
    def run(self):
        # create UDP socket
        s = socket(AF_INET, SOCK_DGRAM)
        s.bind(('', 0))
        s.setsockopt(SOL_SOCKET, SO_BROADCAST, 1)

        while True:
            log.debug('Announcing on port %s with magic %s',self.port,self.magic)
            s.sendto(self.magic, ('<broadcast>', self.port))
            sleep(self.interval)


class Discoverer(Base):
    daemon = True
    def __init__(self,*args,**kwargs):
        super(Discoverer,self).__init__(*args,**kwargs)

        # map of host -> timestamp of last announcement
        self.hosts = dict()

        self.lock = Lock()

    def run(self):
        # create UDP socket
        s = socket(AF_INET, SOCK_DGRAM)
        s.bind(('', self.port))

        while True:
            data, addr = s.recvfrom(len(self.magic))
            if data == self.magic:
                host = addr[0]
                log.debug('Discovered %s with magic %s',host,self.magic)
                self.hosts[host] = time()

    def get_hosts(self,wait=False):
        if wait:
            sleep(self.interval+1)

        min_ts = time() - self.interval * 2

        with self.lock:
            return [h for h,ts in self.hosts.items() if ts > min_ts]
