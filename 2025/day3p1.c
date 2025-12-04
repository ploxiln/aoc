#include <assert.h>
#include <ctype.h>
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	int sum = 0;
	int banks = 0;

	char *line = NULL;
	size_t lsz = 0;
	for (;;) {
		// process one line of input (one battery bank) at a time
		ssize_t r = getline(&line, &lsz, stdin);
		if (r == -1) {
			break;
		}
		// trim right
		while (r > 0 && isspace(line[r-1])) {
			r--;
			line[r] = '\0';
		}
		printf("bank(%d)=%s (%zd)\n", banks, line, r);
		banks++;

		for (ssize_t i = 0; i < r; i++) {
			if (!isdigit(line[i])) {
				printf("ERROR invalid char '%c' offset=%ld\n", line[i], i);
				return 2;
			}
		}
		// finding 2 digits which (in-order) make the max value
		char md1 = 0; // max digit 1
		char md2 = 0; // max digit 2
		size_t mi1 = 0; // index of md1
		size_t mi2 = 0; // index of md2

		// first digit: highest, but not in last position
		// (comparing digit char value is equivalent to comparing decimal value)
		for (ssize_t i = 0; i < r-1; i++)  {
			if (line[i] > md1) {
				md1 = line[i];
				mi1 = i;
			}
		}
		// second digit: highest after first
		for (ssize_t i = mi1+1; i < r; i++)  {
			if (line[i] > md2) {
				md2 = line[i];
				mi2 = i;
			}
		}
		// finally convert chars to decimal value
		int joltage = (md1 - '0') * 10 + (md2 - '0');
		printf("bank=%d max-joltage=%d (b-%zd, b-%zd)\n", banks, joltage, mi1, mi2);
		sum += joltage;
	}
	free(line);

	if (!feof(stdin)) {
		perror("ERROR reading stdin");
		return 1;
	}
	printf("joltage-sum=%d\n", sum);
	return 0;
}
