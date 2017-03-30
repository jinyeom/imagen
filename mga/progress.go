package mga

type ProgressBar struct {
	NumTournament int      // number of tournament
	Progress      chan int // channel for current progress
}
