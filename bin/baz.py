#!/bin/env python

import os
import socket
import sys
import time
import contextlib

def _abstractify(socket_name):
    if socket_name.startswith('@'):
        # abstract namespace socket
        socket_name = '\0%s' % socket_name[1:]
    return socket_name

def _sd_notify(unset_env, msg):
    notify_socket = os.getenv('NOTIFY_SOCKET')
    print('socket={}'.format(notify_socket))
    if notify_socket:
        sock = socket.socket(socket.AF_UNIX, socket.SOCK_DGRAM)
        with contextlib.closing(sock):
            try:
                sock.connect(_abstractify(notify_socket))
                sock.sendall(msg)
                if unset_env:
                    del os.environ['NOTIFY_SOCKET']
                print("Systemd notified")
            except EnvironmentError:
                print("Systemd notification failed")

def notify(remove=False):
    _sd_notify(remove, b'READY=1')


if __name__ == '__main__':
    timeout = float(sys.argv[1])
    time.sleep(timeout)
    notify()
