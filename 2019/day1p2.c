// reads input from stdin: ./day1p2 <input.txt

#include <stdio.h>

int main(int argc, char *argv[]) {
	long total = 0;

	// input parsing not as strict as possible, but AoC input is well-formed
	while (1) {
		char buf[64];
		char *line = fgets(buf, sizeof(buf), stdin);
		if (line == NULL) {
			break;
		}
		if (line[0] == '\0' || line[0] == '\n') {
			continue;
		}

		long mass = 0;
		if (sscanf(line, "%ld", &mass) != 1) {
			printf("ERROR parsing line: %s\n", line);
			return 1;
		}
		printf("> mass=%6ld  fuel:", mass);

		long fuel = mass / 3 - 2;
		while (fuel > 0 ) {
			printf(" %ld", fuel);
			total += fuel;
			fuel = fuel / 3 - 2;
		}
		printf("\n");
	}

	printf("total fuel: %ld\n", total);
	return 0;
}
