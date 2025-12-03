// USAGE: day2p2 < input.txt

// In part-1, noun=12 verb=2. Tracing debug-logged execution (of my input program)
// showed that it boiled down to something like:
// ~= ((noun * a + b) * c + d) * e ... + verb
// New desired result is higher, so first guess bigger values for "noun",
// until result is within 100 of target, then fine-adjust with "verb".

#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define TARGET 19690720

typedef struct {
	int *arr;
	size_t used;
	size_t size;
} IntVector;

void appendIntVector(IntVector *vec, int value) {
	if (vec->used == vec->size) {
		vec->size = vec->size ? vec->size * 2 : 8;
		vec->arr = realloc(vec->arr, vec->size * sizeof(int));
		assert(vec->arr != NULL);
	}
	vec->arr[vec->used++] = value;
}

int runProg(IntVector *vec, int noun, int verb) {
	// copy to new mem array
	size_t used = vec->used;
	int *arr = malloc(used * sizeof(int));
	memcpy(arr, vec->arr, used * sizeof(int));

	assert(used >= 3);
	arr[1] = noun;
	arr[2] = verb;

	// run the program
	size_t ip = 0;
	while (1) {
		assert(ip < used);
		int inst = arr[ip];
		if (inst == 99) {
			break;
		}
		if (inst == 1 || inst == 2) {
			assert(used > ip+3);
			int a1 = arr[ip+1]; assert((int)used > a1);
			int a2 = arr[ip+2]; assert((int)used > a2);
			int a3 = arr[ip+3]; assert((int)used > a3);

			int res = (inst == 1) ? arr[a1] + arr[a2]
			         /* inst 2 */ : arr[a1] * arr[a2];
			arr[a3] = res;
			ip += 4;
		} else {
			printf("ERROR invalid inst=%d @ ip=%zd\n", arr[ip], ip);
			exit(1);
		}
	}
	int out = arr[0];
	free(arr);
	return out;
}

int main(int argc, char *argv[]) {
	IntVector vec = {0};

	char *buf = NULL;
	size_t bsz = 0;
	while (getdelim(&buf, &bsz, ',', stdin) != -1) {
		int val = 0;
		if (sscanf(buf, "%d", &val) != 1) {
			printf("ERROR invalid token (%zd): '%s'\n", vec.used, buf);
			return 1;
		}
		appendIntVector(&vec, val);
	}
	free(buf);
	printf("read %zd opcodes\n", vec.used);

	// "noun" and "verb" in range [0, 100)
	// "noun" has huge effect, "verb" just added at the end
	int noun = 0;
	int verb = 0;
	for (noun = 0; noun < 100; noun++) {
		int res = runProg(&vec, noun, verb);
		printf("noun=%2d  verb=%d  out=%d\n", noun, verb, res);
		if (TARGET - res < 100) {
			verb = TARGET - res;
			break;
		}
	}

	int out = runProg(&vec, noun, verb);
	if (out == TARGET) {
		printf("HIT TARGET  noun=%2d  verb=%d  out=%d\n", noun, verb, out);
		return 0;
	} else {
		printf("FAILED to hit target\n");
		return 1;
	}
}
