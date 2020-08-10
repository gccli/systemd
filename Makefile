PREFIX=/usr

gosocket:

install:
	install -m 0755 bin/foo    $(PREFIX)/bin/
	install -m 0755 bin/bar    $(PREFIX)/bin/
	install -m 0755 bin/foo.socket  $(PREFIX)/bin/
