PREFIX=/usr

BINARY=bin/gosocket bin/foo bin/bar
UNITS:=lib/systemd/system/foo.service
UNITS+=lib/systemd/system/bar.service
UNITS+=lib/systemd/system/gosocket.socket
UNITS+=lib/systemd/system/gosocket.service

install: $(BINARY)
	install -m 0755 $(BINARY) $(PREFIX)/bin/
	install $(UNITS) $(PREFIX)/lib/systemd/system/
	init q

uninstall:
	@for i in ${BINARY}; do \
            x=/usr/bin/$$(basename $$i); \
	    echo rm -f $$x; \
            /bin/rm -f $$x; \
        done

	@for i in ${UNITS}; do \
            x=/usr/lib/systemd/system/$$(basename $$i); \
	    echo rm -f $$x; \
            /bin/rm -f $$x; \
        done
	init q

bin/gosocket:
	go build -mod vendor -o $@ examples/socket/socket.go

clean:
	$(RM) bin/gosocket
