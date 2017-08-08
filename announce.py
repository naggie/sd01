from time import sleep
from socket import socket, AF_INET, SOCK_DGRAM, SOL_SOCKET, SO_BROADCAST

PORT = 51222
MAGIC = "dstv_magic"

# create UDP socket
s = socket(AF_INET, SOCK_DGRAM)
s.bind(('', 0))
s.setsockopt(SOL_SOCKET, SO_BROADCAST, 1)

while 1:
    s.sendto(MAGIC, ('<broadcast>', PORT))
    print "sent service announcement"
    sleep(5)
