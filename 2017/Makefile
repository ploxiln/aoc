# no built-in rules or variables
MAKEFLAGS = -r -R

ALL_PROGS = $(patsubst %.go, %, $(wildcard day*.go))

all: ${ALL_PROGS}

day%: day%.go
	@echo "GO build $@"
	@go build -ldflags="-w -s" $<

clean:
	@rm -vf ${ALL_PROGS}
