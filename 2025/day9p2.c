#include <assert.h>
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


typedef struct {
	int x, y;
} Point;

// length of span *including* both ends (tiles)
static
int span_incl(int A, int B) {
	return (A > B) ? (A + 1 - B)
	               : (B + 1 - A);
}

// rectangular area defined by opposite corner tiles A and B
static
int64_t area_defined(Point A, Point B) {
	return (int64_t) span_incl(A.x, B.x)
	     * (int64_t) span_incl(A.y, B.y);
}


// is span [U,V] *outside* span [A,B]
static
int span_outside(int A, int B, int U, int V) {
	return (U <= A && U <= B && V <= A && V <= B)
	    || (U >= A && U >= B && V >= A && V >= B);
}

// does segment [s1, s2] cut *inside* rectangle [r1, r2]
static
int segment_inside_rect(Point r1, Point r2, Point s1, Point s2) {
	assert(s1.x == s2.x || s1.y == s2.y);
	return !span_outside(r1.x, r2.x, s1.x, s2.x)
	    && !span_outside(r1.y, r2.y, s1.y, s2.y);
}


int main(int argc, char *argv[]) {
	Point *reds = NULL;
	size_t reds_cap = 0,
	       reds_len = 0;

	// read in all points
	while (1) {
		char line[48];
		if (fgets(line, sizeof(line), stdin) == NULL) {
			break;
		}
		Point p = {0};
		if (sscanf(line, "%d,%d", &p.x, &p.y) != 2) {
			printf("ERROR parsing line %zu: '%s'\n", reds_len, line);
			return 1;
		}
		append_vec(reds, reds_cap, reds_len, p);
	}
	assert(feof(stdin));

	printf("Rough Sketch:\n");
	char grid[100][100];
	int marker = 0;
	memset(grid, ' ', sizeof(grid));
	for (size_t i = 0; i < reds_len; i++) {
		int x = reds[i].x / 1000;
		int y = reds[i].y / 1000;
		assert(x < 100);
		assert(y < 100);
		if (grid[y][x] == ' ') {
			grid[y][x] = '0' + (char)marker;
			marker = (marker + 1) % 10;
			if (i == 0) {
				grid[y][x] = '*';
			}
		}
	}
	for (int y = 0; y < 100; y++) {
		for (int x = 0; x < 100; x++) {
			putchar(grid[y][x]);
		}
		putchar('\n');
	}

	// Requirement: rectangle must consist of red and green tiles (must be inside outline).
	// Heuristic: Reject rectangle if any (other) edges extend inside it.
	// It is possible for L shape of region to go completely outside this rect,
	// but for our particular input shape that doesn't happen for the biggest rects.
	//
	// This could be made robust by filling a bitmap, and checking center tile of each rect.
	int64_t area_max = 0;

	for (size_t i = 0; i < reds_len; i++) {
		for (size_t j = i+1; j < reds_len; j++) {
			int64_t area = area_defined(reds[i], reds[j]);
			if (area > area_max) {
				int fail = 0;
				for (size_t u = 0; u < reds_len; u++) {
					size_t v = (u + 1) % reds_len;
					if (segment_inside_rect(reds[i], reds[j], reds[u], reds[v])) {
						// printf("segment (%d,%d)-(%d,%d) inside rect (%d,%d)x(%d,%d)\n",
						//        reds[u].x, reds[u].y, reds[v].x, reds[v].y,
						//        reds[i].x, reds[i].y, reds[j].x, reds[j].y);
						fail = 1;
						break;
					}
				}
				if (!fail) {
					printf("bigger area found (%d,%d)x(%d,%d) = %ld\n",
					       reds[i].x, reds[i].y, reds[j].x, reds[j].y, area);
					area_max = area;
				}
			}
		}
	}

	printf("Greatest area: %ld\n", area_max);

	free(reds);
	return 0;
}
