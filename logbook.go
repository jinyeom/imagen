/*


logbook.go implementation of a logging system of mGA.

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
	"fmt"
	"os"
	"time"
)

// LogBook keeps track of each tournament and its result in mGA.
type LogBook struct {
	Best *Genome  // best performing genome
	Log  []string // log of each tournament
}

// NewLogBook creates a new LogBook, provided the number of tournaments.
func NewLogBook(numTournaments int) *LogBook {
	return &LogBook{
		Best: &Genome{},
		Log:  make([]string, 0, numTournaments),
	}
}

func (l *LogBook) Record(g1, g2 int, score1, score2, bestScore float64) {
	log := fmt.Sprintf("Tournament [%d (%f) and %d (%f)]; Best score: %f",
		g1, score1, g2, score2, bestScore)
	l.Log = append(l.Log, log)
}

func (l *LogBook) Summarize() {
	for _, log := range l.Log {
		fmt.Println(log)
	}
	fmt.Println("Best Genome:")
	fmt.Println(l.Best.ToString())
}

func (l *LogBook) Export() error {
	f, err := os.Create(fmt.Sprintf("imagen_%d.txt", time.Now().UnixNano()))
	if err != nil {
		return err
	}
	defer f.Close()

	for _, log := range l.Log {
		_, err := f.WriteString(log + "\n")
		if err != nil {
			return err
		}
	}

	f.WriteString(l.Best.ToString())

	return nil
}
