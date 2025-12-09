#include <assert.h>
#include <ctype.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define append_vec(vec, vcap, vlen, el) do {             \
	if (vlen == vcap) {                              \
		vcap += 32;                              \
		vec = realloc(vec, vcap * sizeof(*vec)); \
		assert(vec != NULL);                     \
	}                                                \
	vec[vlen++] = el;                                \
} while (0)

ssize_t trim_right(char *str, ssize_t len) {
	while (len > 0 && isspace(str[len-1])) {
		len--;
		str[len] = '\0';
	}
	return len;
}

int main(int argc, char *argv[]) {
	char **grid = NULL;
	size_t grsz = 0;
	size_t rows = 0;
	size_t cols = 0;

	// read input to populate "grid" of chars
	char *line = NULL;
	size_t lsz = 0;
	while (1) {
		ssize_t n = getline(&line, &lsz, stdin);
		if (n < 0) {
			break;
		}
		n = trim_right(line, n);

		if (cols == 0) {
			cols = n;
		}
		assert((size_t)n == cols);

		append_vec(grid, grsz, rows, strdup(line));
	}
	assert(feof(stdin));
	free(line);

	// Run the sim.
	// Idea:      Keep a list of beams to follow "down".
	// Lazy idea: Scan across all cols every row.
	int total_splits = 0;
	assert(rows > 2);

	for (size_t x = 0; x < cols; x++) {
		if (grid[0][x] == 'S') {
			assert(grid[1][x] == '.');
			grid[1][x] = '|';
			break;
		}
	}
	for (size_t y = 1; y < rows-1; y++) {
		int splits = 0;
		int continues = 0;
		for (size_t x = 0; x < cols; x++) {
			if (grid[y][x] == '|') {
				if (grid[y+1][x] == '^') { // splitter
					splits++;
					assert(x > 0);
					assert(x < cols-1);
					grid[y+1][x-1] = '|';
					grid[y+1][x+1] = '|';
				}
				else { // not splitter, should be '.' or '|'
					continues++;
					grid[y+1][x] = '|';
				}
			}
		}
		printf("row=%zu continues=%d splits=%d\n", y, continues, splits);
		total_splits += splits;
	}
	printf("total-splits=%d\n", total_splits);

	for (size_t y = 0; y < rows; y++) {
		free(grid[y]);
	}
	free(grid);
	return 0;
}
