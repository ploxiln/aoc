#include <assert.h>
#include <ctype.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main(int argc, char *argv[]) {
	// grid will have blank rows and columns on the perimiter to allow over-edge reads
	char **grid = NULL;
	ssize_t ncol = 0;
	ssize_t nrow = 0;

	char *line = NULL;
	size_t lsz = 0;
	for (;;) {
		ssize_t r = getline(&line, &lsz, stdin);
		if (r == -1) {
			break;
		}
		// trim right
		while (r > 0 && isspace(line[r-1])) {
			r--;
			line[r] = '\0';
		}
		if (ncol == 0) {
			ncol = r;
			// blank first row
			grid = realloc(grid, 1 * sizeof(char*));
			grid[0] = calloc(ncol + 2, sizeof(char));
		} else {
			if (r != ncol) {
				printf("ERROR inconsistent width %ld != %ld\n", r, ncol);
				return 1;
			}
		}
		nrow++;
		int y = nrow;

		// include blank first row, and blank first+last column
		grid = realloc(grid, (nrow + 1) * sizeof(char*));
		grid[y] = calloc(ncol + 2, sizeof(char));
		memcpy(grid[y] + 1, line, ncol);
	}
	if (!feof(stdin)) {
		perror("ERROR reading stdin");
		return 1;
	}
	free(line);

	// blank last row
	grid = realloc(grid, (nrow + 2) * sizeof(char*));
	grid[nrow + 1] = calloc(ncol + 2, sizeof(char));

	printf("grid: %ld x %ld\n", ncol, nrow);

	// straightforward calculation, for each spot on the grid
	int accessible = 0;
	for (int y = 1; y <= nrow; y++) {
		for (int x = 1; x <= ncol; x++) {
			if (grid[y][x] != '@') {
				continue;
			}
			int adj = 0;
			adj += (int)(grid[y  ][x+1] == '@');
			adj += (int)(grid[y+1][x+1] == '@');
			adj += (int)(grid[y+1][x  ] == '@');
			adj += (int)(grid[y+1][x-1] == '@');
			adj += (int)(grid[y  ][x-1] == '@');
			adj += (int)(grid[y-1][x-1] == '@');
			adj += (int)(grid[y-1][x  ] == '@');
			adj += (int)(grid[y-1][x+1] == '@');
			if (adj < 4) {
				accessible++;
			}
		}
	}
	printf("accessible rolls: %d\n", accessible);

	for (ssize_t i = 0; i < nrow+2; i++) {
		free(grid[i]);
	}
	free(grid);
	return 0;
}
