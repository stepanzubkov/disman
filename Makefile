all: main

PREFIX=/usr
SRCS = main.go pam.go passwd.go env.go daemon.go xorg_server.go console.go

main: $(SRCS)
	go build -o disman $(SRCS)

.PHONY: clean install

clean:
	rm -f main

install:
	install -v disman $(PREFIX)/bin
