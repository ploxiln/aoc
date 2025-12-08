#include <assert.h>
#include <stdint.h>
#include <stdlib.h>
#include <stdio.h>

#define append_vec(vec, vcap, vlen, el) do {             \
	if (vlen == vcap) {                              \
		vcap += 32;                              \
		vec = realloc(vec, vcap * sizeof(*vec)); \
		assert(vec != NULL);                     \
	}                                                \
	vec[vlen++] = el;                                \
} while (0)

# define append(vec, el) append_vec(vec, vec ## _cap, vec ## _len, el)

struct range {
	uint64_t lo, hi;
};

int ranges_overlap(struct range *a, struct range *b) {
	// MEA CULPA - my original range-intersection expression was buggy,
	// and my small test inputs did not expose any bug, so after some hours
	// I finally turned to Claude sonnet, which immediately identified it.
	// These two lines are the only use of LLMs or any kind of assistance, FWIW.
	// My buggy original (which was inline below):
	//        (ranges[i].lo <= ranges[j].lo && ranges[j].lo <= ranges[i].hi)
	//     || (ranges[i].lo <= ranges[j].hi && ranges[j].hi <= ranges[i].hi)
	// Fixed by Claude sonnet:
	return b->lo <= a->hi
	    && a->lo <= b->hi;
}

uint64_t min_u64(uint64_t a, uint64_t b) {
	return a < b ? a : b;
}
uint64_t max_u64(uint64_t a, uint64_t b) {
	return a > b ? a : b;
}

int main(int argc, char *argv[]) {
	struct range *ranges = NULL;
	size_t ranges_cap = 0,
	       ranges_len = 0;

	int ranges_done = 0;
	while (1) {
		char buf[64];
		char *line = fgets(buf, sizeof(buf), stdin);
		if (line == NULL) {
			break;
		}
		if (!ranges_done && line[0] == '\n') {
			printf("Parsed %zu ranges\n", ranges_len);
			ranges_done = 1;
			continue;
		}
		if (!ranges_done) {
			// parse range
			uint64_t low = 0, high = 0;
			if (sscanf(buf, "%lu-%lu", &low, &high) != 2) {
				printf("ERROR failed to parse line '%s'\n", buf);
				return 1;
			}
			assert(low <= high);

			struct range r = {
				.lo = low,
				.hi = high,
			};
			append(ranges, r);
		} else {
			// parse second section of input like day5p1, but ignore
			uint64_t id = -1;
			if (sscanf(buf, "%lu", &id) != 1) {
				printf("ERROR failed to parse line '%s'\n", buf);
				return 1;
			}
		}
	}
	if (!feof(stdin)) {
		perror("ERROR reading stdin");
		return 1;
	}

	// Combine ranges that overlap.
	// After combining, "tombstone"/zero the later range,
	// instead of shuffling to compact the vector, for easier debug.
	int overlaps = 0;
	for (size_t i = 0; i < ranges_len; i++) {
		if (ranges[i].hi == 0) {
			continue; // tombstone - skip
		}
		int found = 0;

		// only check for overlap with later ranges, because
		// earlier ranges were already checked against this range
		for (size_t j = i+1; j < ranges_len; j++) {
			if (ranges[j].hi == 0) {
				continue; // tombstone - skip
			}
			if (ranges_overlap(ranges+i, ranges+j)) {
				printf("combining (%3zu)+(%3zu) : [%15lu,%15lu] + [%15lu,%15lu]\n",
				       i, j, ranges[i].lo, ranges[i].hi, ranges[j].lo, ranges[j].hi);
				found = 1;
				overlaps++;

				ranges[i].lo = min_u64(ranges[i].lo, ranges[j].lo);
				ranges[i].hi = max_u64(ranges[i].hi, ranges[j].hi);

				// tombstone for eliminated range
				ranges[j] = (struct range){0};
			}
		}
		if (found) {
			i--; // range updated, check for overlaps again
		}
	}
	printf("found %d range overlaps\n", overlaps);

	// if no ranges overlap, total fresh IDs is just sum of all range sizes
	uint64_t total_fresh = 0;
	for (size_t i = 0; i < ranges_len; i++) {
		if (ranges[i].hi == 0) {
			continue; // tombstone - skip
		}
		total_fresh += (ranges[i].hi - ranges[i].lo + 1);
	}
	printf("Total fresh: %lu\n", total_fresh);
	free(ranges);

	return 0;
}
