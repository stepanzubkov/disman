all: main

PREFIX=/usr
SRCS = main.go pam.go passwd.go env.go daemon.go xorg_server.go console.go config.go desktop_entry.go

main: $(SRCS)
	go build -o disman $(SRCS)

.PHONY: clean install

clean:
	rm -f main

install:
	install -v disman $(PREFIX)/bin
	cp -v ../res/pam /etc/pam.d/disman


# Install SysVinit configuration for PCLinuxOS
install-pclos:
	cp -v ../res/disman-sysvinit-pclos.conf /usr/share/X11/dm.d/1disman.conf
	sed -i 's/DISPLAYMANAGER=.*/DISPLAYMANAGER=disman/' /etc/sysconfig/desktop
	
# Install Runit configuration
install-runit:
	mkdir -vp /etc/sv/disman
	cp -v ../res/runit-run /etc/sv/disman/run 
