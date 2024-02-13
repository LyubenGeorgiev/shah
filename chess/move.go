package chess

type Move uint32

func encode_move(source, target square, piece, promoted piece, capture, double, enpassant, castling uint32) Move {
	return Move(uint32(source) | uint32(target)<<6 | uint32(piece)<<12 | uint32(promoted)<<16 | uint32(capture)<<20 | uint32(double)<<21 | uint32(enpassant)<<22 | uint32(castling)<<23)
}

// extract source square
func (m Move) getSource() square {
	return square(m & 0x3f)
}

// extract source square
func (m Move) GetSource() int {
	return int(m & 0x3f)
}

// extract target square
func (m Move) getTarget() square {
	return square(m >> 6 & 0x3f)
}

// extract target square
func (m Move) GetTarget() int {
	return int(m >> 6 & 0x3f)
}

// extract piece
func (m Move) getPiece() piece {
	return piece(m >> 12 & 0xf)
}

// extract promoted piece
func (m Move) getPromotionPiece() piece {
	return piece(m >> 16 & 0xf)
}

// extract capture flag
func (m Move) isCapture() bool {
	return m>>20&1 > 0
}

// extract capture flag
func (m Move) IsCapture() bool {
	return m>>20&1 > 0
}

// extract double pawn push flag
func (m Move) isDoublePawnPush() bool {
	return m>>21&1 > 0
}

// extract enpassant flag
func (m Move) isEnpassant() bool {
	return m>>22&1 > 0
}

// extract castling flag
func (m Move) isCastling() bool {
	return m>>23&1 > 0
}

type Moves struct {
	moves [256]Move
	count int
}

func (m *Moves) addMove(move Move) {
	m.moves[m.count] = move
	m.count++
}

func (m *Moves) Count() int {
	return m.count
}

func (m *Moves) At(i int) Move {
	return m.moves[i]
}

/*  =======================
         Move ordering
    =======================

    1. PV move
    2. Captures in MVV/LVA
    3. 1st killer move
    4. 2nd killer move
    5. History moves
    6. Unsorted moves
*/

// score moves
func (e *Engine) score_move(move Move) int {
	// if PV move scoring is allowed
	if e.score_pv != 0 {
		// make sure we are dealing with PV move
		if e.pv_table[0][e.ply] == move {
			// disable score PV flag
			e.score_pv = 0

			// give PV move the highest score to search it first
			return 20000
		}
	}

	// score capture move
	if move.isCapture() {
		// pick up bitboard piece index ranges depending on side
		start_piece, end_piece := P, K

		// pick up side to move
		if e.Side == white {
			start_piece, end_piece = p, k
		}

		// loop over bitboards opposite to the current side to move
		for bb_piece := start_piece; bb_piece <= end_piece; bb_piece++ {
			// if there's a piece on the target square
			if e.Bitboards[bb_piece].getBit(move.getTarget()) {
				// score move by MVV LVA lookup [source piece][target piece]
				return mvv_lva[move.getPiece()][bb_piece] + 10000
			}
		}
	}

	// score quiet move
	if e.killer_moves[0][e.ply] == move { // score 1st killer move
		return 9000
	} else if e.killer_moves[1][e.ply] == move { // score 2nd killer move
		return 8000
	}

	// score history move
	return e.history_moves[move.getPiece()][move.getTarget()]
}

var sortScores = [256]int{}

// sort moves in descending order
func (e *Engine) sort_moves(moves *Moves, bestMove Move) {
	for i := 0; i < moves.count; i++ {
		if moves.moves[i] != bestMove {
			sortScores[i] = e.score_move(moves.moves[i])
		} else {
			sortScores[i] = 30000
		}
	}

	for i := 1; i < moves.count; i++ {
		cur, curScore := moves.moves[i], sortScores[i]

		j := i - 1
		for ; j >= 0 && curScore > sortScores[j]; j-- {
			moves.moves[j+1], sortScores[j+1] = moves.moves[j], sortScores[j]
		}

		// Move the current element to its correct position
		moves.moves[j+1], sortScores[j+1] = cur, curScore
	}
}

// position repetition detection
func (e *Engine) is_repetition() bool {
	// loop over repetition indicies range
	for _, repetition := range e.RepetitionTable {
		if repetition == e.HashKey {
			return true
		}
	}

	return false
}
