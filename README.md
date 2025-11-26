## My Advent of Code solutions

See https://adventofcode.com/

The purpose of publishing these solutions is just to share tips with friends.
This is not the most robust code I've written; the goal is rather to be
clear, concise, and stand-alone.

 * Error handling is typically just "abort" / "panic".
 * Input might not be strictly validated.
 * Still, possibility of memory-corruption (etc) is not OK.
 * Logging and option handling are very simplistic.
 * For part-2, I "fork" the part-1 solution into a separate file.
 * Common helper functions are also just copied into each solution as needed.
 * The problem statement is not repeated, read it on the adventofcode website
   (you already did it yourself, right...)
 * These solutions are not published promptly when the problem is solved,
   but rather a day or more later, with just a *bit* of cleanup (e.g. comments).

Each solution is named like: `2017/day4p1.go` (or `.c`)

To build: `cd 2017 && make day4p1`

To run: `cd 2017 && ./day4p1 < input.txt`

Some solutions take a command arg or two, if very short, instead of input on `stdin`.
For `stdin` input, you can also just paste it:

```sh
$ ./day4p1
... waits for input, paste here (and then <enter> if last line not blank)
<ctrl-d> (on blank line sends EOF, aka end-of-file)
```
