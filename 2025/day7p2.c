#include <assert.h>
#include <ctype.h>
#include <stdio.h>
#include <stdint.h>
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
	// Part 2:    Need a *count* for each beam position,
	//            different histories that come together *add*.
	// Attempt 2: Numbers get HUGE, need 64bit counters, put
	//            them in separate shadow grid aka "beam[][]".

	int64_t **beam = calloc(rows, sizeof(int64_t*));
	for (size_t y = 0; y < rows; y++) {
		beam[y] = calloc(cols, sizeof(int64_t));
	}

	assert(rows > 2);
	for (size_t x = 0; x < cols; x++) {
		if (grid[0][x] == 'S') {
			assert(grid[1][x] == '.');
			beam[1][x] = 1;
			break;
		}
	}
	for (size_t y = 1; y < rows-1; y++) {
		for (size_t x = 0; x < cols; x++) {
			if (beam[y][x] > 0) {
				if (grid[y+1][x] == '^') { // splitter
					assert(x > 0);
					assert(x < cols-1);
					beam[y+1][x-1] += beam[y][x];
					beam[y+1][x+1] += beam[y][x];

					printf("SPLIT @ (%zu,%zu) l=%ld r=%ld\n",
					        x, y, beam[y+1][x-1], beam[y+1][x+1]);
				}
				else { // continue beam histories straight down
					beam[y+1][x] += beam[y][x];
				}
			}
		}
	}
	int64_t total_histories = 0;
	for (size_t x = 0; x < cols; x++) {
		total_histories += beam[rows-1][x];
	}
	printf("total-histories=%ld\n", total_histories);

	for (size_t y = 0; y < rows; y++) {
		free(grid[y]);
		free(beam[y]);
	}
	free(grid);
	free(beam);
	return 0;
}
