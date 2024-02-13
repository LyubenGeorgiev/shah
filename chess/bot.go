package chess

// parse UCI "go" command
func (e *Engine) Search() Move {
	clear_hash_table()
	e.TimeController = NewTimeController()

	// init start time
	e.Timeset = 1
	e.StartTime = GetTimeMs()
	e.StopTime = e.StartTime + 2000

	// search position

	e.nodes = 0

	// reset "time is up" flag
	e.Stopped = 0

	// reset follow PV flags
	e.follow_pv = 0
	e.score_pv = 0

	// clear helper data structures for search
	e.Searcher = Searcher{}

	var safeBestMove Move

	// iterative deepening
	for current_depth := 1; current_depth <= 64; current_depth++ {
		// if time is up
		if e.Stopped == 1 {
			// stop calculating and return best move so far
			break
		}

		// enable follow PV flag
		e.follow_pv = 1

		// find best move within a given position
		e.negamax(-INFINITY, INFINITY, current_depth)

		// Might have bugs if we dont do this!
		if e.Stopped == 1 {
			return safeBestMove
		} else {
			safeBestMove = e.pv_table[0][0]
		}
	}

	return safeBestMove
}
