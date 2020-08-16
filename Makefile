PREFIX=/usr

srcfoobar := $(wildcard examples/general/*.go)
srcsocket := $(wildcard examples/socket/*.go)

BINARY=bin/foobar bin/socket bin/foo bin/bar
UNITS:=lib/systemd/system/foo.service
UNITS+=lib/systemd/system/bar.service
UNITS+=lib/systemd/system/baz.service
UNITS+=lib/systemd/system/qux.service
UNITS+=lib/systemd/system/foobar@.service
UNITS+=lib/systemd/system/foobar.target
UNITS+=lib/systemd/system/plugh.socket
UNITS+=lib/systemd/system/plugh.service

install: $(BINARY)
	install -m 0755 $(BINARY) $(PREFIX)/bin/
	install -m 0644 $(UNITS)  $(PREFIX)/lib/systemd/system/
	init q

uninstall:
	@for i in ${BINARY}; do \
            x=/usr/bin/$$(basename $$i); \
	    echo rm -f $$x; \
            /bin/rm -f $$x; \
        done

	@for i in ${UNITS}; do \
            systemctl disable $$(basename $$i); \
            x=/usr/lib/systemd/system/$$(basename $$i); \
	    echo rm -f $$x; \
            /bin/rm -f $$x; \
        done
	init q

bin/foobar: $(srcfoobar)
	go build -mod vendor -o $@ $^

bin/socket:$(srcsocket)
	go build -mod vendor -o $@ $^

clean:
	$(RM) bin/socket bin/foobar
