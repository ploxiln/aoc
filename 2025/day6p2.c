#include <assert.h>
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

#define ROWS 5
#define MAXVALS 6

uint64_t eval_problem(char op, uint64_t vals[], int n) {
	printf("PROBLEM: vals=%d op=%c :", n, op);
	assert(op == '+' || op == '*');
	assert(n > 0);

	uint64_t subtot = (op == '*') ? 1 : 0;

	for (int j = 0; j < n; j++) {
		printf(" %lu", vals[j]);
		subtot = (op == '*') ? (subtot * vals[j])
		                     : (subtot + vals[j]);
	}
	printf(" = %lu\n", subtot);
	return subtot;
}

int main(int argc, char *argv[]) {
	// N rows of value digits (and spaces), then a row of operators (and spaces)
	// (don't bother parsing on ingest, to be parsed vertically while processing)
	char *grid[ROWS] = {0};
	size_t cols = 0;

	char *line = NULL;
	size_t lsz = 0;
	size_t row = 0;
	for (;;) {
		ssize_t r = getline(&line, &lsz, stdin);
		if (r < 0) {
			break;
		}
		// trim cr/nl
		while (r > 0 && (line[r-1] == '\n' || line[r-1] == '\r')) {
			line[--r] = '\0';
		}
		printf("got row %zu (cols=%zd)\n", row, r);

		if (row == ROWS) {
			printf("ERROR too many rows (max = %d)\n", ROWS);
			return 1;
		}
		if (cols == 0) {
			cols = r;
		}
		assert((size_t)r == cols);

		grid[row] = strdup(line);
		row++;
	}
	assert(feof(stdin));
	assert(row == ROWS);
	free(line);

	uint64_t total = 0;
	uint64_t vals[MAXVALS];
	int vused = 0;
	char op = '?';

	for (size_t i=0; i < cols; i++) {
		// calculate value for this column, if digits present
		uint64_t val = 0;
		for (int j=0; j < ROWS-1; j++) {
			char c = grid[j][i];
			if (c == ' ') {
				continue;
			}
			assert(c >= '0' && c <= '9');
			val *= 10;
			val += c - '0';
		}
		if (val > 0) {
			assert(vused < MAXVALS);
			vals[vused++] = val;
		}
		// operator char in this column, if present
		char opc = grid[ROWS-1][i];
		if (opc != ' ') {
			assert(op == '?');
			assert(opc == '+' || opc == '*');
			op = opc;
		}
		// blank column
		if (val == 0 && opc == ' ') {
			if (vused == 0 && op == '?') {
				continue; // first blank
			}
			total += eval_problem(op, vals, vused);
			vused = 0;
			op = '?';
		}
	}
	if (vused > 0) { // problem in last column
		total += eval_problem(op, vals, vused);
	}
	printf("total sum of answers: %lu\n", total);

	for (int j=0; j<ROWS; j++) {
		free(grid[j]);
	}
	return 0;
}
