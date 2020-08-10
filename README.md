# install example

    make install

# Socket example

    # start gosocket socket unit
    systemctl start gosocket.socket

    # monitor gosocket service unit
    journalctl -u gosocket.service -f

    # send traffic to gosocket.socket
    curl http://127.0.0.1:50004

# notify and watchdog

设置`Type=notify`启用notify功能设置`WatchdogSec=20`启用看门狗服务。`Restart=`设为`on-failure`
