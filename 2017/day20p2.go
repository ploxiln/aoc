package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
)

type Vec3 struct {
	X, Y, Z int64
}

func (v *Vec3) Add(a Vec3) {
	v.X += a.X
	v.Y += a.Y
	v.Z += a.Z
}

type Particle struct {
	Pos, Vel, Acc Vec3

	// cache manhattanDist() of each Vec3 for sorting
	Dist, VelM, AccM int64

	// original input index
	Idx int
}

func (p *Particle) Step() {
	p.Vel.Add(p.Acc)
	p.Pos.Add(p.Vel)
	p.UpdateMags()
}
func (p *Particle) UpdateMags() {
	p.Dist = manhattanDist(p.Pos)
	p.VelM = manhattanDist(p.Vel)
	p.AccM = manhattanDist(p.Acc)
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

	idx := 0
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Bytes()
		m := particleRe.FindSubmatch(line)
		if m == nil {
			fmt.Printf("ERROR failed to parse: %q\n", line)
		}
		p := &Particle{
			Idx: idx,
			Pos: Vec3{mInt(m[1]), mInt(m[2]), mInt(m[3])},
			Vel: Vec3{mInt(m[4]), mInt(m[5]), mInt(m[6])},
			Acc: Vec3{mInt(m[7]), mInt(m[8]), mInt(m[9])},
		}
		p.UpdateMags()
		particles = append(particles, p)
		idx++
	}
	if err := scanner.Err(); err != nil {
		abort("ERROR reading stdin: %s", err)
	}
	fmt.Printf("Parsed %d particles\n", len(particles))

	fmt.Printf("First 5 particles:\n")
	for i, p := range particles[:5] {
		fmt.Printf("%3d: p=<%5d,%5d,%5d> v=<%4d,%4d,%4d> a=<%3d,%3d,%3d>\n", i,
			p.Pos.X, p.Pos.Y, p.Pos.Z,
			p.Vel.X, p.Vel.Y, p.Vel.Z,
			p.Acc.X, p.Acc.Y, p.Acc.Z,
		)
	}

	// Run particle movement simulation. Keep particles sorted by distance
	// from origin, to make detecting collisions easier. Run simulation
	// until particles are also in order of total velocity and acceleration.
	for n := 0; n < 2000; n++ {
		sort.Slice(particles, func(i, j int) bool {
			return particles[i].Dist < particles[j].Dist
		})

		ooVel := 0
		ooAcc := 0
		for i := 0; i+1 < len(particles); i++ {
			if particles[i].VelM > particles[i+1].VelM {
				ooVel += 1
			}
			if particles[i].AccM > particles[i+1].AccM {
				ooAcc += 1
			}
		}
		if ooVel + ooAcc > 0 {
			if n % 100 == 0 {
				fmt.Printf(
					"Step %4d: Particles not yet ordered by velocity (%d) or acceleration (%d)\n",
					n, ooVel, ooAcc)
			}
		} else {
			fmt.Printf("Step %4d: All particles now ordered by distance/velocity/acceleration\n", n)
			return
		}

		collisions := map[int]bool{}
		collide := func (i int) {
			if !collisions[i] {
				collisions[i] = true
				fmt.Printf("Collision: %3d @ <%5d,%5d,%5d>\n",
					particles[i].Idx,
					particles[i].Pos.X,
					particles[i].Pos.Y,
					particles[i].Pos.Z,
				)
			}
		}
		for i := 0; i < len(particles); i++ {
			for j := i+1; j < len(particles); j++ {
				if particles[i].Dist != particles[j].Dist {
					break
				}
				if particles[i].Pos == particles[j].Pos {
					collide(i)
					collide(j)
				}
			}
		}
		for i := range collisions {
			n := len(particles) - 1
			for collisions[n] {
				particles = particles[:n]
				n--
			}
			if i > n {
				continue
			}
			particles[i] = particles[n]
			particles = particles[:n]
		}
		if len(collisions) > 0 {
			fmt.Printf("Step %4d: %d collisions, %d particles remain\n", n, len(collisions), len(particles))
		}

		// move/update step
		for _, p := range particles {
			p.Step()
		}
	}
}
