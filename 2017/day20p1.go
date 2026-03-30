package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Vec3 struct {
	X, Y, Z int64
}

type Particle struct {
	Pos, Vel, Acc Vec3
}

func abort(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func mInt(str []byte) int64 {
	v, err := strconv.ParseInt(string(str), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse int %q: %s", str, err))
	}
	return v
}

func abs64(i int64) int64 {
	// bit trick, no branch, maybe pointless :)
	return i - ((i * 2) & (i >> 63))
}

func manhattanDist(v Vec3) int64 {
	return abs64(v.X) + abs64(v.Y) + abs64(v.Z)
}

// relative init velocity component (magnitude, negative only if opposite acceleration)
func relInitVelComp(acc, vel int64) int64 {
	if acc != 0 && (acc < 0) != (vel < 0) {
		return 0 - abs64(vel)
	} else {
		return abs64(vel)
	}
}

// sum of magnitude of velocity components, but negative if opposite acceleration
// (for comparing particles with same "total" acceleration magnitude)
func relInitVelocity(p *Particle) int64 {
	v := int64(0)
	v += relInitVelComp(p.Acc.X, p.Vel.X)
	v += relInitVelComp(p.Acc.Y, p.Vel.Y)
	v += relInitVelComp(p.Acc.Z, p.Vel.Z)
	return v
}

func main() {
	particleRe, err := regexp.Compile(" *" +
		"p=<(-?[0-9]+),(-?[0-9]+),(-?[0-9]+)>, *" +
		"v=<(-?[0-9]+),(-?[0-9]+),(-?[0-9]+)>, *" +
		"a=<(-?[0-9]+),(-?[0-9]+),(-?[0-9]+)>\\s*",
	)
	if err != nil {
		abort("ERROR compiling regexp: %s", err)
	}

	particles := []*Particle{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Bytes()
		m := particleRe.FindSubmatch(line)
		if m == nil {
			fmt.Printf("ERROR failed to parse: %q\n", line)
		}
		particles = append(particles, &Particle{
			Pos: Vec3{mInt(m[1]), mInt(m[2]), mInt(m[3])},
			Vel: Vec3{mInt(m[4]), mInt(m[5]), mInt(m[6])},
			Acc: Vec3{mInt(m[7]), mInt(m[8]), mInt(m[9])},
		})
	}
	if err := scanner.Err(); err != nil {
		abort("ERROR reading stdin: %s", err)
	}
	fmt.Printf("Parsed %d particles\n", len(particles))

	fmt.Printf("First 5 particles:\n")
	for i, p := range particles[:5] {
		fmt.Printf("%3d: p=<%5d,%5d,%5d> v=<%4d,%4d,%4d> a=<%3d,%3d,%3d> mag(a)=%d\n", i,
			p.Pos.X, p.Pos.Y, p.Pos.Z,
			p.Vel.X, p.Vel.Y, p.Vel.Z,
			p.Acc.X, p.Acc.Y, p.Acc.Z,
			manhattanDist(p.Acc),
		)
	}

	// Theory: In the "long term", acceleration dominates.
	// Tie-break with inital velocity: particles with same total acceleration
	// will always have same difference in total velocity. Velocity components
	// magnitudes matter, but are counted negative (only) if opposite acceleration.
	// For "manhattan distance" acceleration, just add acceleration components.

	minAcc := manhattanDist(particles[0].Acc)
	for _, p := range particles {
		acc := manhattanDist(p.Acc)
		if acc < minAcc {
			minAcc = acc
		}
	}

	fmt.Printf("Min accelerating particles:\n")
	minId := -1
	minInitVel := int64(-1)
	for i, p := range particles {
		if manhattanDist(p.Acc) == minAcc {
			relInitVel := relInitVelocity(p)
			if minId == -1 || relInitVel < minInitVel {
				minId = i
				minInitVel = relInitVel
			}
			fmt.Printf("%3d: p=<%5d,%5d,%5d> v=<%4d,%4d,%4d> a=<%3d,%3d,%3d> relInitV=%d\n", i,
				p.Pos.X, p.Pos.Y, p.Pos.Z,
				p.Vel.X, p.Vel.Y, p.Vel.Z,
				p.Acc.X, p.Acc.Y, p.Acc.Z,
				relInitVel,
			)
		}
	}
	if minId != -1 {
		fmt.Printf("Long-term closest to origin: particle %d\n", minId)
	}
}
