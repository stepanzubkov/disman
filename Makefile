all: main

SRCS = main.go pam.go passwd.go env.go daemon.go xorg_server.go console.go

main: $(SRCS)
	go build $(SRCS)

.PHONY: clean

clean:
	rm -f main
