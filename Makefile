all: main

SRCS = main.go pam.go

main: $(SRCS)
	go build $(SRCS)

.PHONY: clean

clean:
	rm -f main
