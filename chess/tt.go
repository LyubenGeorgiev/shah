package chess

// transposition table hash flags
const (
	no_hash_entry   = 100000
	hash_flag_exact = 0
	hash_flag_alpha = 1
	hash_flag_beta  = 2
)

// transposition table data structure
type tt struct {
	hashKey  Bitboard
	depth    int // current search depth
	flag     int // flag the type of node (fail-low/fail-high/PV)
	score    int // score (alpha/beta/PV)
	bestMove Move
}

var hash_table [8 * 1024 * 1024]tt

var empty_hash_record = tt{}

// clear TT (hash table)
func clear_hash_table() {
	for i := range hash_table {
		hash_table[i].hashKey = 0
	}
}

// read hash entry data
func (e *Engine) read_hash_entry(alpha, beta, depth int, bestMove *Move) int {
	// hash, _ := hashstructure.Hash(e.Board, hashstructure.FormatV2, nil)
	hash_entry := hash_table[e.HashKey%Bitboard(len(hash_table))]

	if hash_entry.hashKey == e.HashKey {
		if hash_entry.depth == depth {
			// extract stored score from TT entry
			score := hash_entry.score

			// retrieve score independent from the actual path
			// from root node (position) to current node (position)
			if score < -MATE_SCORE {
				score += e.ply
			}
			if score > MATE_SCORE {
				score -= e.ply
			}

			// match the exact (PV node) score
			if hash_entry.flag == hash_flag_exact {
				// return exact (PV node) score
				return score
			}

			// match alpha (fail-low node) score
			if (hash_entry.flag == hash_flag_alpha) && (score <= alpha) {
				// return alpha (fail-low node) score
				return alpha
			}

			// match beta (fail-high node) score
			if (hash_entry.flag == hash_flag_beta) && (score >= beta) {
				// return beta (fail-high node) score
				return beta
			}
		}

		// Store best move
		*bestMove = hash_entry.bestMove

	}

	// if hash entry doesn't exist
	return no_hash_entry
}

// write hash entry data
func (e *Engine) write_hash_entry(score, depth, hash_flag int, bestMove Move) {
	// store score independent from the actual path
	// from root node (position) to current node (position)
	if score < -MATE_SCORE {
		score -= e.ply
	}
	if score > MATE_SCORE {
		score += e.ply
	}

	// hash, _ := hashstructure.Hash(e.Board, hashstructure.FormatV2, nil)
	hash_table[e.HashKey%Bitboard(len(hash_table))] = tt{hashKey: e.HashKey, flag: hash_flag, depth: depth, score: score, bestMove: bestMove}
}

// enable PV move scoring
func (e *Engine) enable_pv_scoring(moves *Moves) {
	// disable following PV
	e.follow_pv = 0

	// loop over the moves within a move list
	for count := 0; count < moves.count; count++ {
		// make sure we hit PV move
		if e.pv_table[0][e.ply] == moves.moves[count] {
			// enable move scoring
			e.score_pv = 1

			// enable following PV
			e.follow_pv = 1
		}
	}
}
