# Examples

```
foo.service --------> bar.service (After=foo.service)

             activate
plugh.socket --------> foobar@.service

              activate
foobar.socket --------> foobar@.service


       foobar.target
             |  (Wants=,Before=)
             v
         baz.service  ------------> quuz.service
            /  \                 (Requires=,After=)
           /    \
          /      \
         /        \
        /          \
       v            v
qux.service   quux.service
 (PartOf=)   (Wants=,After=)

```

* Before=,Wants= https://unix.stackexchange.com/questions/507702/before-and-want-for-the-same-systemd-service

# Install example

    make install

## Test Socket activation

    systemctl start plugh.socket       # start plugh socket unit
    journalctl -u plugh.service -f     # monitor plugh service unit
    curl -i http://127.0.0.1:50004     # send traffic to plugh.socket and activate plugh.service

## notify and watchdog

设置`Type=notify`启用notify功能设置`WatchdogSec=20`启用看门狗服务。`Restart=`设为`on-failure`

## Test restart

    分别测试正常退出，异常退出，crush，超时
    curl http://127.0.0.1:50004?exit=normal
    curl http://127.0.0.1:50004?exit=abnormal
    curl http://127.0.0.1:50004?exit=abort
    curl http://127.0.0.1:50004?exit=timeout
