// USAGE: day8p1 [connections] < input.txt
#include <assert.h>
#include <ctype.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>


typedef struct {
	int x, y, z;
	int circuit;
} Point;

typedef struct {
	int idx_a, idx_b;
	int64_t dist_sq;
} Conn;

// squared distance (skip sqrt, expensive)
int64_t dist_sq(Point *a, Point *b) {
	int64_t dx = a->x - b->x;
	int64_t dy = a->y - b->y;
	int64_t dz = a->z - b->z;
	return (dx * dx) + (dy * dy) + (dz * dz);
}

#define append_vec(vec, vcap, vlen, el) do {             \
	if (vlen == vcap) {                              \
		vcap += 32;                              \
		vec = realloc(vec, vcap * sizeof(*vec)); \
		assert(vec != NULL);                     \
	}                                                \
	vec[vlen++] = el;                                \
} while (0)

# define append(vec, el) append_vec(vec, vec ## _cap, vec ## _len, el)

void insert_dist(Conn *arr, size_t cap, size_t *len, int idx_a, int idx_b, int64_t dist_sq) {
	size_t i = *len;
	while (i > 0 && arr[i-1].dist_sq > dist_sq) {
		i--;
	}
	if (i < *len) {
		// shift
		size_t count = *len - i - (*len < cap ? 0 : 1);
		memmove(arr+i+1, arr+i, count * sizeof(Conn));
	}
	if (i < *len || *len < cap) {
		// insert
		arr[i] = (Conn){idx_a, idx_b, dist_sq};
	}
	if (*len < cap) {
		(*len)++;
	}
}

// compare func for qsort(), integers, reverse (descending)
int cmpint_rev(const void *p1, const void *p2) {
	return *(int*)p2 - *(int*)p1;
}

int main(int argc, char *argv[]) {
	Point *points = NULL;
	size_t points_cap = 0,
	       points_len = 0;

	if (argc != 2) {
		printf("USAGE: day8p1 [connections] < input.txt\n");
		return 2;
	}
	// find N shortest connections
	int connections = atoi(argv[1]);

	while (1) {
		char line[48];
		if (fgets(line, sizeof(line), stdin) == NULL) {
			break;
		}
		Point p = {0};
		if (sscanf(line, "%d,%d,%d", &p.x, &p.y, &p.z) != 3) {
			printf("ERROR parsing line '%s'\n", line);
			return 1;
		}
		append(points, p);
	}
	assert(feof(stdin));
	printf("read %zu points ...\n", points_len);

	// ranked shortest connections
	Conn *shortest = calloc(connections, sizeof(Conn));
	size_t shortest_len = 0;

	// check them all, insertion-sort into ranked array
	for (size_t i = 0; i < points_len; i++) {
		for (size_t j = i+1; j < points_len; j++) {
			int64_t d = dist_sq(points + i, points + j);
			insert_dist(shortest, connections, &shortest_len, i, j, d);
		}
	}
	printf("ranked %zu shortest connections ...\n", shortest_len);

	// numbered circuits, store count of points (junction boxes) connected
	int *circuits = calloc(connections, sizeof(int));
	int clen = 1; // index 0 is unused

	// connect points (boxes) into circuits
	for (size_t i = 0; i < shortest_len; i++) {
		int ia = shortest[i].idx_a;
		int ib = shortest[i].idx_b;
		int ca = points[ia].circuit;
		int cb = points[ib].circuit;

		if (i < 10) { // debug
			Point a = points[ia], b = points[ib];
			printf("%zd: (%d,%d,%d) [%d] <-> (%d,%d,%d) [%d] : %5ld\n",
			       i, a.x, a.y, a.z, a.circuit, b.x, b.y, b.z, b.circuit, shortest[i].dist_sq);
		}
		if (ca == 0 && cb == 0) {
			// two un-linked points/boxes
			assert(clen < connections);
			points[ia].circuit = clen;
			points[ib].circuit = clen;
			circuits[clen] = 2;
			clen++;
		}
		else if (ca == 0 || cb == 0) {
			// an unlinked point/box joins a circuit
			int c = ca | cb;
			points[ia].circuit = c;
			points[ib].circuit = c;
			circuits[c] += 1;
		}
		else if (ca == cb) {
			// noop
		}
		else {
			// join two circuits (merge smaller into larger)
			printf("JOIN %d [%d] + %d [%d]\n", ca, circuits[ca], cb, circuits[cb]);
			int c = (circuits[ca] > circuits[cb]) ? ca : cb;
			int c2 = (c == ca) ? cb : ca;
			for (size_t j = 0; j < points_len; j++) {
				if (points[j].circuit == c2) {
					points[j].circuit = c;
				}
				// TODO short-circuit after circuits[c2] found?
			}
			circuits[c] += circuits[c2];
			circuits[c2] = 0;
		}
	}
	for (int i = 0; i < clen; i++) {
		if (circuits[i] > 0) {
			printf("circuit %d : %d boxes\n", i, circuits[i]);
		}
	}

	// find 3 biggest circuits
	assert(clen >= 3);
	qsort(circuits, clen, sizeof(circuits[0]), cmpint_rev);
	printf("biggest circuits: %d , %d , %d | product = %d\n",
	       circuits[0], circuits[1], circuits[2], circuits[0] * circuits[1] * circuits[2]);

	free(points);
	free(shortest);
	free(circuits);
	return 0;
}
