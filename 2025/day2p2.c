#include <assert.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

// is this string made up of N repeating parts?
int str_repeat(const char *s, int len, int parts) {
	if (len % parts != 0) {
		return 0; // no, not divisible
	}
	int pl = len / parts; // part length
	for (int i = 0; i < pl; i++) {
		for (int j = 1; j < parts; j++) {
			if (s[i] != s[j*pl+i]) {
				return 0; // no, mismatch
			}
		}
	}
	return 1; // parts all matched
}

// is this id "invalid" (repeating decimal digit pattern)
int id_invalid(uint64_t id) {
	char intstr[32];
	int len = snprintf(intstr, sizeof(intstr), "%lu", id);
	assert(len < 32);

	// check every equal division of the string
	for (int d = 2; d <= len; d++) {
		if (str_repeat(intstr, len, d)) {
			return 1; // invalid
		}
	}
	return 0; // valid
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
				printf("invalid ID: %lu\n", id);
				sum += id;
			}
		}
	}
	free(buf);

	printf("sum of invalid IDs : %lu\n", sum);
	return 0;
}
