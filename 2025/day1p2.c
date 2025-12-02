// reads input from stdin: ./day1p1 <day1.in.txt

#include <stdio.h>

int main(int argc, char *argv[]) {
	int pos = 50;  // initial position
	int zeros = 0; // number of stops on *or passes by* zero

	while (1) {
		// input lines very short
		char buf[32];
		char *line = fgets(buf, sizeof(buf), stdin);
		if (line == NULL) {
			break;
		}
		if (line[0] == '\0' || line[0] == '\n') {
			continue;
		}

		char dir = 0;
		int dist = 0;
		if (sscanf(line, "%1c%d", &dir, &dist) != 2) {
			printf("ERROR parsing line: %s\n", line);
			return 1;
		}

		switch (dir) {
		case 'L':
			if (pos == 0) {
				// loop just below will assume pos *passed* zero
				// (but no, it started on zero) and add this back
				zeros--;
			}
			pos -= dist;
			while (pos < 0) {
				pos += 100;
				zeros++;
			}
			if (pos == 0) {
				// stop on zero counts (now)
				zeros++;
			}
			break;

		case 'R':
			pos += dist;
			while (pos > 99) {
				pos -= 100;
				zeros++;
			}
			// if stop on zero, already counted it
			break;

		default:
			printf("ERROR invalid line: %s\n", line);
			return 1;
		}
	}
	printf("final position=%d zeros=%d\n", pos, zeros);
	return 0;
}
