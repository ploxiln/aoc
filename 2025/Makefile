# no built-in rules or variables
MAKEFLAGS = -r -R

CC = cc
WARNINGS = -Wall -Wextra -Wshadow -Wwrite-strings -Wno-unused-parameter
# _XOPEN_SOURCE=700 needed for some libc functions like getline() or strtok_r() ...
ALL_CFLAGS = -std=c11 -D_XOPEN_SOURCE=700 -O1 -g ${WARNINGS} ${CFLAGS}
# include std math lib in case useful, and basic "hardening" options, why not
ALL_LDFLAGS = -lm -z now -z relro ${LDFLAGS}

ALL_PROGS = $(patsubst %.c, %, $(wildcard day*.c))

all: ${ALL_PROGS}

%: %.c
	@echo "CC $@"
	@${CC} ${ALL_CFLAGS} $< ${ALL_LDFLAGS} -o $@

clean:
	@rm -vf ${ALL_PROGS}

.PHONY: all clean
