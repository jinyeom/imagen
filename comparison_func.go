/*


comparison_func.go implementation of comparison functions for mGA.

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

// ComparisonFunc defines a type of function that compares two fitness scores
// (float64) and returns true if the first argument fitness score is more fit
// than the second, false if otherwise.
type ComparisonFunc func(float64, float64) bool

// DirectComparison returns a comparison function that returns true if the
// first argument fitness score is higher than the second.
func DirectComparison() ComparisonFunc {
	return func(score0, score1 float64) bool {
		return score0 > score1
	}
}

// InverseComparison returns a comparison function that returns true if the
// first argument fitness score is lower than the second.
func InverseComparison() ComparisonFunc {
	return func(score0, score1 float64) bool {
		return score0 < score1
	}
}
