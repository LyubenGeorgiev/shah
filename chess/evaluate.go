package chess

func (b *Board) evaluate() int {
	// static evaluation score
	score := 0

	// loop over piece bitboards
	for bb_piece := P; bb_piece <= k; bb_piece++ {
		// Material score
		score += b.bitboards[bb_piece].countBits() * materialScore[bb_piece]

		// Positional score
		for bitboard := b.bitboards[bb_piece]; bitboard > 1; {
			square := square(bitboard.getLs1bIndex())
			score += scores[bb_piece][square]
			bitboard.popBit(square)
		}
	}

	// return final evaluation based on side
	if b.side == white {
		return score
	}

	return -score
}
