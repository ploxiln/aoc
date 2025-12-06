// Find intersection of two paths on 2d grid ... ideas:
// 1. allocate a big-enough grid, mark every square along one path ...
// 2. sparse grid in a hash-table ...
// 3. for every segment of path B, check for intersection with every segment of path A
//    ^^ yeah that seems better

#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// for simpler intersection test, order x/y "low" to "high"
struct segment {
	char dir;
	int dist;
	int xl, xh, yl, yh;
};

// parse comma-separated relative segments, convert to absolute x,y segments
ssize_t parse_grid_path(struct segment **dst, char *path) {
	size_t segments = 1;
	for (char *p = path; *p != '\0'; p++) {
		if (*p == ',') {
			segments++;
		}
	}
	struct segment *arr = calloc(segments, sizeof(struct segment));
	assert(arr != NULL);
	size_t used = 0;

	// track current position
	int x = 0, y = 0;

	char *tok = NULL, *saveptr = NULL;
	for (tok = strtok_r(path,  ",", &saveptr); tok != NULL;
	     tok = strtok_r(NULL,  ",", &saveptr))
	{
		char dir[2] = {0};
		int dist = 0;
		if (sscanf(tok, "%1[UDLR]%d", dir, &dist) != 2) {
			printf("ERROR parsing segment %zd '%s'\n", used, tok);
			free(arr);
			return -1;
		}
		struct segment seg = {
			.dir = dir[0],
			.dist = dist,
			.xl = x,
			.xh = x,
			.yl = y,
			.yh = y,
		};
		switch (dir[0]) {
		case 'D':
			y -= dist;
			seg.yl = y;
			break;
		case 'U':
			y += dist;
			seg.yh = y;
			break;
		case 'L':
			x -= dist;
			seg.xl = x;
			break;
		case 'R':
			x += dist;
			seg.xh = x;
			break;
		default:
			printf("ERROR impossible dir %c\n", dir[0]);
			free(arr);
			return -1;
		}
		assert(used < segments);
		memcpy(arr+used, &seg, sizeof(seg));
		used++;
	}
	*dst = arr;
	return (ssize_t)used;
}

int segments_intersect(struct segment *s0, struct segment *s1, int *x_out, int *y_out) {
	// find vertical and horizontal segments
	struct segment *sh = NULL, *sv = NULL;

	if (s0->dir == 'U' || s0->dir == 'D') {
		if (s1->dir == 'U' || s1->dir == 'D') {
			return 0; // both vertical
		}
		sv = s0;
		sh = s1;
	}
	else { // s0->dir is 'L' or 'R'
		if (s1->dir == 'L' || s1->dir == 'R') {
			return 0; // both horizontal
		}
		sh = s0;
		sv = s1;
	}
	if (sh->yl == 0 && sv->xl == 0) {
		return 0; // if intersect, at origin, skip
	}
	// finally, intersection test
	if (sv->yl <= sh->yl && sh->yh <= sv->yh &&
	    sh->xl <= sv->xl && sv->xh <= sh->xh   ) {
		// HIT
		*y_out = sh->yl;
		*x_out = sv->xl;
		return 1;
	}
	return 0;
}

int main(int argc, char *argv[]) {
	struct segment *wire[2] = {0};
	size_t          segs[2] = {0};
	int wi = 0; // wire 0 or 1

	// read in all wire segments, and convert relative movements to absolute coords
	char *buf = NULL;
	size_t bsz = 0;
	while (getline(&buf, &bsz, stdin) != -1) {
		if (wi >= 2) {
			printf("ERROR too many lines of input (max %d)\n", 2);
			return 1;
		}
		ssize_t n = parse_grid_path(wire + wi, buf);
		if (n < 0) {
			return 1;
		}
		printf("PARSED w%d segments=%zd\n", wi, n);
		segs[wi] = n;
		wi++;
	}
	if (!feof(stdin)) {
		perror("ERROR reading stdin");
		return 1;
	}
	free(buf);

	// for each segment of wire 0, check for intersection with every segment of wire 1
	int steps_min = 999999;

	int w0_steps = 0;
	for (size_t i=0; i < segs[0]; i++) {
		struct segment *s0 = wire[0] + i;
		w0_steps += s0->dist;

		int w1_steps = 0;
		for (size_t j=0; j < segs[1]; j++) {
			struct segment *s1 = wire[1] + j;
			w1_steps += s1->dist;

			int x = 0, y = 0;
			if(segments_intersect(s0, s1, &x, &y)) {
				// add total steps of each wire so far, but subtract
				// steps *past* the intersection on current segments
				int steps = w0_steps + w1_steps;
				switch (s0->dir) {
				case 'U': steps -= s0->yh -     y ; break;
				case 'D': steps -=     y  - s0->yl; break;
				case 'R': steps -= s0->xh -     x ; break;
				case 'L': steps -=     x  - s0->xl; break;
				}
				switch (s1->dir) {
				case 'U': steps -= s1->yh -     y ; break;
				case 'D': steps -=     y  - s1->yl; break;
				case 'R': steps -= s1->xh -     x ; break;
				case 'L': steps -=     x  - s1->xl; break;
				}
				if (steps < steps_min) {
					steps_min = steps;
				}
				printf("INTERSECT (%d, %d) steps=%d\n", x, y, steps);
			}
		}
	}
	printf("FEWEST STEPS to intersect = %d\n", steps_min);
	free(wire[0]);
	free(wire[1]);
	return 0;
}
