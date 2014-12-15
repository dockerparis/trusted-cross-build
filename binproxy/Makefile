BINS = /bin/sh /bin/bash /bin/ls /bin/ln /bin/chmod /bin/echo /bin/cp /bin/cat /usr/bin/apt-get /bin/env

all:
	for bin in $(BINS); do \
	    ln -sf /bin/wrapper `basename $$bin`; \
	done

clean:
	for bin in $(BINS); do \
	    rm -f `basename $$bin`; '
	done
