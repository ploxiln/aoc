// USAGE: day2p2 < input.txt

#include <assert.h>
#include <stdio.h>
#include <stdlib.h>

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

int main(int argv, char *argc[]) {
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

	// for part 1, replace position 1 with value 12, and replace position 2 with value 2
	assert(vec.used >= 3);
	vec.arr[1] = 12;
	vec.arr[2] = 2;

	// run the program
	size_t ip = 0;
	while (1) {
		assert(ip < vec.used);
		int inst = vec.arr[ip];
		if (inst == 99) {
			printf("PROGRAM DONE\n");
			break;
		}
		if (inst == 1 || inst == 2) {
			assert(vec.used > ip+3);
			int a1 = vec.arr[ip+1];
			int a2 = vec.arr[ip+2];
			int a3 = vec.arr[ip+3];
			int proglen = (int)vec.used;
			assert(proglen > a1);
			assert(proglen > a2);
			assert(proglen > a3);
			int res = (inst == 1) ? vec.arr[a1] + vec.arr[a2]
			         /* inst 2 */ : vec.arr[a1] * vec.arr[a2];

			// debug/progress
			printf("%3ld: %d %3d %3d %3d (%6d %c %6d = %6d)\n",
				ip, inst, a1, a2, a3, vec.arr[a1], inst == 1 ? '+' : '*', vec.arr[a2], res);

			vec.arr[a3] = res;
			ip += 4;
		} else {
			printf("ERROR invalid inst=%d @ ip=%zd\n", vec.arr[ip], ip);
			return 1;
		}
	}
	printf("mem addr 0: %d\n", vec.arr[0]);
}
