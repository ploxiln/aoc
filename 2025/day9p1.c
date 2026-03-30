#include <assert.h>
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>

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

	// find areas defined by all combinations of red tiles, keep largest
	int64_t area_max = 0;

	for (size_t i = 0; i < reds_len; i++) {
		for (size_t j = i+1; j < reds_len; j++) {
			int64_t area = area_defined(reds[i], reds[j]);
			if (area > area_max) {
				printf("bigger area found (%d,%d)x(%d,%d) = %ld\n",
				       reds[i].x, reds[i].y, reds[j].x, reds[j].y, area);
				area_max = area;
			}
		}
	}
	printf("greatest area: %ld\n", area_max);

	free(reds);
	return 0;
}
