package chess

type Move uint32

func encode_move(source, target square, piece, promoted piece, capture, double, enpassant, castling uint32) Move {
	return Move(uint32(source) | uint32(target)<<6 | uint32(piece)<<12 | uint32(promoted)<<16 | uint32(capture)<<20 | uint32(double)<<21 | uint32(enpassant)<<22 | uint32(castling)<<23)
}

// extract source square
func (m Move) getSource() square {
	return square(m & 0x3f)
}

// extract target square
func (m Move) getTarget() square {
	return square(m >> 6 & 0x3f)
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

type Moves [256]Move
