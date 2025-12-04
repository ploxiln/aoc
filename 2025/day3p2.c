#include <assert.h>
#include <ctype.h>
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>

// digits to pick for each battery bank's "joltage"
#define NDIGITS 12

int main(int argc, char *argv[]) {
	int64_t sum = 0;
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
		char digits[NDIGITS+1] = {0};  // digits of result (and string null term)
		ssize_t look = 0;              // index to start looking for next digit

		// always best to use highest digit possible first,
		// but need to leave enough later digits for total NDIGITS
		// (comparing digit char value is equivalent to comparing decimal value)
		for (int i = 0; i < NDIGITS; i++) {
			ssize_t stop = r - (NDIGITS - (i+1));
			for (ssize_t j = look; j < stop; j++) {
				if (line[j] > digits[i]) {
					digits[i] = line[j];
					look = j + 1;
				}
			}
		}
		// finally convert chars to decimal value
		int64_t joltage = atoll(digits);
		printf("max-joltage=%ld\n", joltage);
		sum += joltage;
	}
	free(line);

	if (!feof(stdin)) {
		perror("ERROR reading stdin");
		return 1;
	}
	printf("joltage-sum=%ld\n", sum);
	return 0;
}
