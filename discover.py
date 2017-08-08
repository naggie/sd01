from socket import socket, AF_INET, SOCK_DGRAM

PORT = 51222
MAGIC = "dstv_magic"

# create UDP socket
s = socket(AF_INET, SOCK_DGRAM)
s.bind(('', PORT))

while 1:
    data, addr = s.recvfrom(len(MAGIC))
    if data == MAGIC:
        print "got service announcement from",addr[0]
