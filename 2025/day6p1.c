#include <assert.h>
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

#define VALROWS 4

#define append_vec(vec, vcap, vlen, el) do {             \
	if (vlen == vcap) {                              \
		vcap += 32;                              \
		vec = realloc(vec, vcap * sizeof(*vec)); \
		assert(vec != NULL);                     \
	}                                                \
	vec[vlen++] = el;                                \
} while (0)

int main(int argc, char *argv[]) {
	// N rows of numbers, then a row of operators
	int *vals[VALROWS] = {0};
	char *ops = NULL;
	size_t cols = 0;


	char *line = NULL;
	size_t lsz = 0;
	size_t row = 0;
	for (;;) {
		ssize_t r = getline(&line, &lsz, stdin);
		if (r == -1) {
			break;
		}
		size_t n = 0;
		size_t sz = 0;

		// parse space-separated columns for each row
		char *tok = NULL, *saveptr = NULL;
		for (tok = strtok_r(line,  " \n", &saveptr); tok != NULL;
		     tok = strtok_r(NULL,  " \n", &saveptr))
		{
			if (row < VALROWS) { // number row
				char *endptr = NULL;
				long val = strtol(tok, &endptr, 10);
				if (endptr == tok || *endptr != '\0') {
					printf("ERROR parsing row=%zu col=%zu : '%s'\n", row, n, tok);
					return 1;
				}
				append_vec(vals[row], sz, n, val);
			}
			else if (row == VALROWS) { // operator row
				if (strcmp(tok, "*") != 0 && strcmp(tok, "+") != 0) {
					printf("ERROR invalid operator col=%zu : '%s'\n", n, tok);
					return 1;
				}
				append_vec(ops, sz, n, tok[0]);
			}
			else {
				printf("ERROR too many rows (max = %d)\n", VALROWS + 1);
				return 1;
			}
		}
		printf("parsed row %zu (cols=%zu)\n", row, n);
		if (cols == 0) {
			cols = n;
		}
		assert(n == cols);
		row++;
	}
	assert(feof(stdin));
	assert(row == VALROWS+1);
	free(line);

	uint64_t total = 0;
	for (size_t i=0; i<cols; i++) {
		uint64_t subtot = (ops[i] == '*') ? 1 : 0;

		for (int j=0; j<VALROWS; j++) {
			uint64_t v = vals[j][i];
			subtot = (ops[i] == '*') ? (subtot * v)
			                         : (subtot + v);
		}
		total += subtot;
	}
	printf("total sum of answers: %lu\n", total);

	for (int j=0; j<VALROWS; j++) {
		free(vals[j]);
	}
	free(ops);
	return 0;
}
