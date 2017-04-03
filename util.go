/*


util.go extra functions needed for mGA.

@licstart   The following is the entire license notice for
the Go code in this page.

Copyright (C) 2017 jin yeom

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.

As additional permission under GNU GPL version 3 section 7, you
may distribute non-source (e.g., minimized or compacted) forms of
that code without the copy of the GNU GPL normally required by
section 4, provided you include this license notice and a URL
through which recipients can access the Corresponding Source.

@licend    The above is the entire license notice
for the Go code in this page.


*/

package main

import (
	"math/rand"
)

// randWeight returns a random connection weight.
func randWeight() float64 {
	return rand.NormFloat64()
}

// randAFuncName returns a random activation function name.
func randAFuncName() string {
	options := []string{
		"identity",
		"sigmoid",
		"tanh",
		"relu",
		//"elu",
		//"abs",
		"sine",
		"gaussian",
	}
	return options[rand.Intn(len(options))]
}

func randGenome(population []*Genome) *Genome {
	return population[rand.Intn(len(population))]
}
