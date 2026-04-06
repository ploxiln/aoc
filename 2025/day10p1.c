#include <assert.h>
#include <ctype.h>
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

#define countof(x) (sizeof(x) / sizeof((x)[0]))

// The "lights" are 1 bit per light, least-significant bit first.
// Each of "buttons" is similarly 1 bit per possible light index toggled.
typedef struct {
	uint16_t bits;
	uint16_t lights;
	uint16_t blen;
	uint16_t buttons[16];
} Machine;

static
uint16_t parse_lights(const char *str) {
	size_t n = strlen(str);
	assert(str[0  ] == '[');
	assert(str[n-1] == ']');

	uint16_t val = 0;
	for (size_t i = 1; i < n-1; i++) {
		if (str[i] == '#') {
			val |= 1 << (i-1);
		}
	}
	return val;
}

static
uint16_t parse_button(char *str) {
	size_t n = strlen(str);
	assert(str[0  ] == '(');
	assert(str[n-1] == ')');

	uint16_t val = 0;

	char *saveptr = NULL;
	char *tok = strtok_r(str, "(,)", &saveptr);
	while (tok != NULL) {
		// only up to 10 bits, 0-9, in actual input
		assert(isdigit(tok[0]));
		assert(tok[1] == '\0');
		val |= 1 << (tok[0] - '0');

		tok = strtok_r(NULL, "(,)", &saveptr);
	}
	return val;
}

static
void parse_machine(const char *line, Machine *m) {
	char *saveptr = NULL;
	char *lcpy = strdup(line);

	char *tok = strtok_r(lcpy, " ", &saveptr);
	assert(tok != NULL);
	m->lights = parse_lights(tok);
	m->bits = strlen(tok) - 2;

	while (1) {
		tok = strtok_r(NULL, " ", &saveptr);
		if (tok == NULL) {
			break;
		}
		if (tok[0] == '{') {
			// joltage requirements not yet needed
			break;
		}
		assert(m->blen < 16);
		m->buttons[m->blen] = parse_button(tok);
		m->blen++;
	}
	printf("parsed machine: lights=x%03hx buttons=x%03hx", m->lights, m->buttons[0]);
	for (uint16_t i = 1; i < m->blen; i++) {
		printf(",x%03hx", m->buttons[i]);
	}
	printf("\n");
	free(lcpy);
}

static
uint16_t press_buttons(Machine *m, uint16_t buttons) {
	uint16_t lights = 0;
	for (uint16_t i = 0; i < m->blen; i++) {
		if (buttons & (1 << i)) {
			lights ^= m->buttons[i];
		}
	}
	return lights;
}

static
uint16_t bitcount16(uint16_t val) {
	// gratuitous branchless optimization
	uint16_t c = val;
	c = (c & 0x5555) + ((c >> 1) & 0x5555);
	c = (c & 0x3333) + ((c >> 2) & 0x3333);
	c = (c & 0x0F0F) + ((c >> 4) & 0x0F0F);
	c = (c & 0x00FF) + ((c >> 8) & 0x00FF);
	return c;
}

static
uint16_t fewest_buttons_to_match(Machine *m) {
	uint16_t bcount = 255;
	uint16_t buttons = 0;

	uint16_t bcombos = 1 << m->blen;
	for (uint16_t x = 1; x < bcombos; x++) {
		if (m->lights == press_buttons(m, x)) {
			uint16_t bits = bitcount16(x);
			if (bits < bcount) {
				buttons = x;
				bcount = bits;
			}
		}
	}
	printf("machine lights=x%03hx press=x%03hx (%d)\n", m->lights, buttons, bcount);
	assert(bcount != 255);
	return bcount;
}

int main(int argc, char *argv[]) {
	int total_presses = 0;

	char *line = NULL;
	size_t lsz = 0;
	while (1) {
		ssize_t n = getline(&line, &lsz, stdin);
		if (n < 0) {
			break;
		}
		// trim right
		while (n > 0 && isspace(line[n-1])) {
			n--;
			line[n] = '\0';
		}

		Machine m = {0};
		parse_machine(line, &m);
		total_presses += fewest_buttons_to_match(&m);
	}
	assert(feof(stdin));
	free(line);
	printf("total-presses=%d\n", total_presses);
	return 0;
}
