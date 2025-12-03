#include <assert.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

// is this id "invalid" (repeating decimal digit pattern)
int id_invalid(uint64_t id) {
	char intstr[32];
	int len = snprintf(intstr, sizeof(intstr), "%lu", id);
	assert(len < 32);

	if (len % 2 != 0) {
		// cannot split in 2 halves
		return 0;
	}

	// check first and second half match
	for (int j = 0; j < len/2; j++) {
		if (intstr[j] != intstr[len/2 + j]) {
			return 0; // no, mismatch
		}
	}

	// yes, "invalid" pattern
	return 1;
}

int main(int argc, char *argv[]) {
	uint64_t sum = 0;
	// lazy: "%lu" (below) only correct for 64-bit unix (should use PRIu64 (ugly))

	char *buf = NULL;
	size_t bsz = 0;
	while (getdelim(&buf, &bsz, ',', stdin) != -1) {
		uint64_t low = 0, high = 0;
		if (sscanf(buf, "%lu-%lu,", &low, &high) != 2) {
			printf("ERROR invalid range: '%s'\n", buf);
			return 1;
		}
		printf("checking ID range: [%lu, %lu]\n", low, high);
		for (uint64_t id = low; id <= high; id++) {
			if (id_invalid(id)) {
				printf("invalid ID: %ld\n", id);
				sum += id;
			}
		}
	}
	free(buf);

	printf("sum of invalid IDs : %lu\n", sum);
	return 0;
}
