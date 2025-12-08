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

int main(int argc, char *argv[]) {
	struct range *ranges = NULL;
	size_t ranges_cap = 0,
	       ranges_len = 0;

	int ranges_done = 0;
	int total_avail = 0;
	int fresh_avail = 0;

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
			// parse id, then check in the ranges
			uint64_t id = -1;
			if (sscanf(buf, "%lu", &id) != 1) {
				printf("ERROR failed to parse line '%s'\n", buf);
				return 1;
			}
			total_avail++;
			// could sort ranges and binary-search? but this is fast enough for now
			for (size_t i = 0; i < ranges_len; i++) {
				if (ranges[i].lo <= id && id <= ranges[i].hi) {
					fresh_avail++;
					break;
				}
			}
		}
	}
	if (!feof(stdin)) {
		perror("ERROR reading stdin");
		return 1;
	}
	printf("Total fresh: %d / %d\n", fresh_avail, total_avail);
	free(ranges);
	return 0;
}
