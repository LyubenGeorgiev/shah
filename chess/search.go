package chess

import "fmt"

const (
	INFINITY   = 50000
	MATE_VALUE = 49000
	MATE_SCORE = 48000
	MAX_PLY    = 64
)

// MVV LVA [attacker][victim]
var mvv_lva = [12][12]int{
	{105, 205, 305, 405, 505, 605, 105, 205, 305, 405, 505, 605},
	{104, 204, 304, 404, 504, 604, 104, 204, 304, 404, 504, 604},
	{103, 203, 303, 403, 503, 603, 103, 203, 303, 403, 503, 603},
	{102, 202, 302, 402, 502, 602, 102, 202, 302, 402, 502, 602},
	{101, 201, 301, 401, 501, 601, 101, 201, 301, 401, 501, 601},
	{100, 200, 300, 400, 500, 600, 100, 200, 300, 400, 500, 600},
	{105, 205, 305, 405, 505, 605, 105, 205, 305, 405, 505, 605},
	{104, 204, 304, 404, 504, 604, 104, 204, 304, 404, 504, 604},
	{103, 203, 303, 403, 503, 603, 103, 203, 303, 403, 503, 603},
	{102, 202, 302, 402, 502, 602, 102, 202, 302, 402, 502, 602},
	{101, 201, 301, 401, 501, 601, 101, 201, 301, 401, 501, 601},
	{100, 200, 300, 400, 500, 600, 100, 200, 300, 400, 500, 600},
}

type Searcher struct {
	// killer moves [id][ply]
	killer_moves [2][MAX_PLY]Move
	// history moves [piece][square]
	history_moves [12][64]int
	// PV length [ply]
	pv_length [MAX_PLY]int
	// PV table [ply][ply]
	pv_table [MAX_PLY][MAX_PLY]Move
	// follow PV & score PV move
	follow_pv Move
	score_pv  Move
}

// quiescence search
func (e *Engine) quiescence(alpha, beta int) int {
	// every 2047 nodes
	if (e.nodes & 2047) == 0 {
		e.timeout()
	}

	// increment nodes count
	e.nodes++

	// we are too deep, hence there's an overflow of arrays relying on max ply constant
	if e.ply > MAX_PLY-1 {
		// evaluate position
		return e.evaluate()
	}

	// evaluate position
	evaluation := e.evaluate()

	// fail-hard beta cutoff
	if evaluation >= beta {
		// node (position) fails high
		return beta
	}

	// found a better move
	if evaluation > alpha {
		// PV node (position)
		alpha = evaluation
	}

	// create move list instance
	var moves Moves

	// generate moves
	e.generateCaptureMoves(&moves)

	// sort moves
	e.sort_moves(&moves, 0)

	// loop over moves within a movelist
	for count := 0; count < moves.count; count++ {
		// preserve board state
		boardCopy := e.Board

		// increment ply
		e.ply++

		// make sure to make only legal moves
		if !e.MakeMove(moves.moves[count], true) {
			// decrement ply
			e.ply--

			// skip to next move
			continue
		}

		// score current move
		score := -e.quiescence(-beta, -alpha)

		// decrement ply
		e.ply--

		// take move back
		e.Board = boardCopy

		// reutrn 0 if time is up
		if e.Stopped == 1 {
			return 0
		}

		// found a better move
		if score > alpha {
			// PV node (position)
			alpha = score

			// fail-hard beta cutoff
			if score >= beta {
				// node (position) fails high
				return beta
			}
		}
	}

	// node (position) fails low
	return alpha
}

const (
	// full depth moves counter
	full_depth_moves = 4
	// depth limit to consider reduction
	reduction_limit = 3
)

// negamax alpha beta search
func (e *Engine) negamax(alpha, beta, depth int) int {
	// init PV length
	e.pv_length[e.ply] = e.ply

	// variable to store current move's score (from the static evaluation perspective)
	var score int

	// best move (to store in TT)
	var bestMove Move

	// define hash flag
	hash_flag := hash_flag_alpha

	// if position repetition occurs
	if e.ply > 0 && e.is_repetition() || e.Fifty >= 100 {
		// return draw score
		return 0
	}

	// a hack by Pedro Castro to figure out whether the current node is PV node or not
	pv_node := beta-alpha > 1

	// read hash entry if we're not in a root ply and hash entry is available
	// and current node is not a PV node
	if e.ply > 0 {
		score = e.read_hash_entry(alpha, beta, depth, &bestMove)
		if score != no_hash_entry && !pv_node {
			// if the move has already been searched (hence has a value)
			// we just return the score for this move without searching it
			return score
		}
	}

	// every 2047 nodes
	if (e.nodes & 2047) == 0 {
		// "listen" to the GUI/user input
		e.timeout()
	}

	// recursion escapre condition
	if depth == 0 {
		// run quiescence search
		return e.quiescence(alpha, beta)
	}

	// we are too deep, hence there's an overflow of arrays relying on max ply constant
	if e.ply > MAX_PLY-1 {
		// evaluate position
		return e.evaluate()
	}

	e.nodes++

	// is king in check
	var in_check bool
	if e.Board.Side == white {
		in_check = e.isSquareAttacked(square(e.Bitboards[K].GetLs1bIndex()), e.Board.Side.opposite())
	} else {
		in_check = e.isSquareAttacked(square(e.Bitboards[k].GetLs1bIndex()), e.Board.Side.opposite())
	}

	// increase search depth if the king has been exposed into a check
	if in_check {
		depth++
	}

	// legal moves counter
	legal_moves := 0

	// null move pruning
	if depth >= 3 && !in_check && e.ply != 0 {
		// preserve board state
		boardCopy := e.Board

		// increment ply
		e.ply++

		// increment repetition index & store hash key
		e.RepetitionTable = append(e.RepetitionTable, e.HashKey)

		// hash enpassant if available
		if e.Enpassant != no_sq {
			e.HashKey ^= enpassant_keys[e.Enpassant]
		}

		// reset enpassant capture square
		e.Enpassant = no_sq

		// switch the side, literally giving opponent an extra move to make
		e.Side = e.Side.opposite()

		// hash the side
		e.HashKey ^= side_key

		/* search moves with reduced depth to find beta cutoffs
		   depth - 1 - R where R is a reduction limit */
		score = -e.negamax(-beta, -beta+1, depth-1-2)

		// decrement ply
		e.ply--

		// decrement repetition index
		e.RepetitionTable = e.RepetitionTable[:len(e.RepetitionTable)-1]

		// restore board state
		e.Board = boardCopy

		// reutrn 0 if time is up
		if e.Stopped == 1 {
			return 0
		}

		// fail-hard beta cutoff
		if score >= beta {
			// node (position) fails high
			return beta
		}
	}

	// create move list instance
	var moves Moves

	// generate moves
	e.generateMoves(&moves)

	// if we are now following PV line
	if e.follow_pv != 0 {
		// enable PV move scoring
		e.enable_pv_scoring(&moves)
	}

	// sort moves
	e.sort_moves(&moves, bestMove)

	// number of moves searched in a move list
	var moves_searched = 0

	// loop over moves within a movelist
	for count := 0; count < moves.count; count++ {
		// preserve board state
		boardCopy := e.Board

		// increment ply
		e.ply++

		// increment repetition index & store hash key
		e.RepetitionTable = append(e.RepetitionTable, e.HashKey)

		// make sure to make only legal moves
		if !e.MakeMove(moves.moves[count], false) {
			// decrement ply
			e.ply--

			// decrement repetition index
			e.RepetitionTable = e.RepetitionTable[:len(e.RepetitionTable)-1]

			// skip to next move
			continue
		}

		// increment legal moves
		legal_moves++

		// full depth search
		if moves_searched == 0 {
			// do normal alpha beta search
			score = -e.negamax(-beta, -alpha, depth-1)
		} else { // late move reduction (LMR)
			// condition to consider LMR
			if moves_searched >= full_depth_moves && depth >= reduction_limit && !in_check && !moves.moves[count].isCapture() && moves.moves[count].getPromotionPiece() == no_piece {
				// search current move with reduced depth:
				score = -e.negamax(-alpha-1, -alpha, depth-2)
			} else { // hack to ensure that full-depth search is done
				score = alpha + 1
			}

			// principle variation search PVS
			if score > alpha {
				score = -e.negamax(-alpha-1, -alpha, depth-1)

				if (score > alpha) && (score < beta) {
					score = -e.negamax(-beta, -alpha, depth-1)
				}
			}
		}

		// decrement ply
		e.ply--

		// decrement repetition index
		e.RepetitionTable = e.RepetitionTable[:len(e.RepetitionTable)-1]

		// take move back
		e.Board = boardCopy

		// reutrn 0 if time is up
		if e.Stopped == 1 {
			return 0
		}

		// increment the counter of moves searched so far
		moves_searched++

		// found a better move
		if score > alpha {
			// switch hash flag from storing score for fail-low node
			// to the one storing score for PV node
			hash_flag = hash_flag_exact

			// store best move (for TT)
			bestMove = moves.moves[count]

			// on quiet moves
			if !moves.moves[count].isCapture() {
				// store history moves
				e.history_moves[moves.moves[count].getPiece()][moves.moves[count].getTarget()] += depth
			}

			// PV node (position)
			alpha = score

			// write PV move
			e.pv_table[e.ply][e.ply] = moves.moves[count]

			// loop over the next ply
			for next_ply := e.ply + 1; next_ply < e.pv_length[e.ply+1]; next_ply++ {
				// copy move from deeper ply into a current ply's line
				e.pv_table[e.ply][next_ply] = e.pv_table[e.ply+1][next_ply]
			}

			// adjust PV length
			e.pv_length[e.ply] = e.pv_length[e.ply+1]

			// fail-hard beta cutoff
			if score >= beta {
				// store hash entry with the score equal to beta
				e.write_hash_entry(beta, depth, hash_flag_beta, bestMove)

				// on quiet moves
				if !moves.moves[count].isCapture() {
					// store killer moves
					e.killer_moves[1][e.ply] = e.killer_moves[0][e.ply]
					e.killer_moves[0][e.ply] = moves.moves[count]
				}

				// node (position) fails high
				return beta
			}
		}
	}

	// we don't have any legal moves to make in the current postion
	if legal_moves == 0 {
		// king is in check
		if in_check {
			// return mating score (assuming closest distance to mating position)
			return -MATE_VALUE + e.ply
		} else { // king is not in check
			// return stalemate score
			return 0
		}
	}

	// store hash entry with the score equal to alpha
	e.write_hash_entry(alpha, depth, hash_flag, bestMove)

	// node (position) fails low
	return alpha
}

// search position for the best move
func (e *Engine) searchPosition(depth int) {
	// define best score variable
	score := 0

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
	for current_depth := 1; current_depth <= depth; current_depth++ {
		// if time is up
		if e.Stopped == 1 {
			// stop calculating and return best move so far
			break
		}

		// enable follow PV flag
		e.follow_pv = 1

		// find best move within a given position
		score = e.negamax(-INFINITY, INFINITY, current_depth)

		// Might have bugs if we dont do this!
		if e.Stopped == 1 {
			// print best move
			fmt.Printf("bestmove %s\n", moveToString(safeBestMove))
			return
		} else {
			safeBestMove = e.pv_table[0][0]
		}

		// if PV is available
		if e.pv_length[0] != 0 {
			// print search info
			if score > -MATE_VALUE && score < -MATE_SCORE {
				fmt.Printf("info score mate %d depth %d nodes %d time %d pv ", -(score+MATE_VALUE)/2-1, current_depth, e.nodes, GetTimeMs()-e.StartTime)
			} else if score > MATE_SCORE && score < MATE_VALUE {
				fmt.Printf("info score mate %d depth %d nodes %d time %d pv ", (MATE_VALUE-score)/2+1, current_depth, e.nodes, GetTimeMs()-e.StartTime)
			} else {
				fmt.Printf("info score cp %d depth %d nodes %d time %d pv ", score, current_depth, e.nodes, GetTimeMs()-e.StartTime)
			}

			// loop over the moves within a PV line
			for count := 0; count < e.pv_length[0]; count++ {
				// print PV move
				fmt.Printf("%s ", pvMoveToString(e.pv_table[0][count]))
			}

			// print new line
			fmt.Println()
		}
	}

	// print best move
	fmt.Printf("bestmove %s\n", moveToString(e.pv_table[0][0]))
}
