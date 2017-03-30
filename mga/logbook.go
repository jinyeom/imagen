package mga

import (
	"fmt"
	"os"
	"time"
)

type LogBook struct {
	Best *Genome  // best performing genome
	Log  []string // log of each tournament
}

func NewLogBook(numTournaments int) *LogBook {
	return &LogBook{
		Best: &Genome{},
		Log:  make([]string, 0, numTournaments),
	}
}

func (l *LogBook) Record(g1, g2 int, score1, score2, bestScore float64) {
	log := fmt.Sprintf("Tournament [%d (%.3f) and %d (%.3f)]; Best score: %.3f",
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
	f, err := os.Create(fmt.Sprintf("EAN_log_%d.txt", time.Now().UnixNano()))
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
