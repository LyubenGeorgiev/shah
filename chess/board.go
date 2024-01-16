package chess

type Board struct {
	bitboards   [12]Bitboard
	occupancies [3]Bitboard
	side        side
	enpassant   square
	castle      castle
}

// init all variables
func Init() {
	// init leaper pieces attacks
	init_leapers_attacks()

	// init slider pieces attacks
	init_sliders_attacks(bishop)
	init_sliders_attacks(rook)
}

// parse FEN string
func NewBoadFromFen(fen []byte) *Board {
	var b Board

	// reset game state variables
	b.enpassant = no_sq

	// loop over board ranks
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			// init current square
			square := square(rank*8 + file)

			// match ascii pieces within FEN string
			if (fen[0] >= 'a' && fen[0] <= 'z') || (fen[0] >= 'A' && fen[0] <= 'Z') {
				// init piece type
				piece := charToPiece[fen[0]]

				// set piece on corresponding bitboard
				b.bitboards[piece].setBit(square)

				// increment pointer to FEN string
				fen = fen[1:]
			}

			// match empty square numbers within FEN string
			if fen[0] >= '0' && fen[0] <= '9' {
				// init offset (convert char 0 to int 0)
				offset := fen[0] - '0'

				// define piece variable
				piece := no_piece

				// loop over all piece bitboards
				for bb_piece := P; bb_piece <= k; bb_piece++ {
					// if there is a piece on current square
					if b.bitboards[bb_piece].getBit(square) {
						// get piece code
						piece = bb_piece
					}
				}

				// on empty current square
				if piece == no_piece {
					// decrement file
					file--
				}

				// adjust file counter
				file += int(offset)

				// increment pointer to FEN string
				fen = fen[1:]
			}

			// match rank separator
			if fen[0] == '/' {
				// increment pointer to FEN string
				fen = fen[1:]
			}
		}
	}

	// got to parsing side to move (increment pointer to FEN string)
	fen = fen[1:]

	// parse side to move
	if fen[0] == 'w' {
		b.side = white
	} else {
		b.side = black
	}

	// go to parsing castling rights
	fen = fen[2:]

	// parse castling rights
	for fen[0] != ' ' {
		switch fen[0] {
		case 'K':
			b.castle |= wk
		case 'Q':
			b.castle |= wq
		case 'k':
			b.castle |= bk
		case 'q':
			b.castle |= bq
		case '-':

		}

		// increment pointer to FEN string
		fen = fen[1:]
	}

	// got to parsing enpassant square (increment pointer to FEN string)
	fen = fen[1:]

	// parse enpassant square
	if fen[0] != '-' {
		// parse enpassant file & rank
		file := fen[0] - 'a'
		rank := 8 - (fen[1] - '0')

		// init enpassant square
		b.enpassant = square(rank*8 + file)
	} else { // no enpassant square
		b.enpassant = no_sq
	}

	// loop over white pieces bitboards
	for piece := P; piece <= K; piece++ {
		// populate white occupancy bitboard
		b.occupancies[white] |= b.bitboards[piece]
	}
	// loop over black pieces bitboards
	for piece := p; piece <= k; piece++ {
		// populate white occupancy bitboard
		b.occupancies[black] |= b.bitboards[piece]
	}

	// init all occupancies
	b.occupancies[both] = b.occupancies[white] | b.occupancies[black]

	return &b
}

// get bishop attacks
func (b *Board) get_bishop_attacks(square square) Bitboard {
	// get bishop attacks assuming current board occupancy
	occupancy := b.occupancies[both]

	occupancy &= bishop_masks[square]
	occupancy *= bishop_magic_numbers[square]
	occupancy >>= 64 - bishop_relevant_bits[square]

	// return bishop attacks
	return bishop_attacks[square][occupancy]
}

// get rook attacks
func (b *Board) get_rook_attacks(square square) Bitboard {
	// get bishop attacks assuming current board occupancy
	occupancy := b.occupancies[both]

	occupancy &= rook_masks[square]
	occupancy *= rook_magic_numbers[square]
	occupancy >>= 64 - rook_relevant_bits[square]

	// return rook attacks
	return rook_attacks[square][occupancy]
}

// get queen attacks
func (b *Board) get_queen_attacks(square square) Bitboard {
	// get bishop attacks assuming current board occupancy
	bishop_occupancy := b.occupancies[both]
	rook_occupancy := b.occupancies[both]

	bishop_occupancy &= bishop_masks[square]
	bishop_occupancy *= bishop_magic_numbers[square]
	bishop_occupancy >>= 64 - bishop_relevant_bits[square]

	rook_occupancy &= rook_masks[square]
	rook_occupancy *= rook_magic_numbers[square]
	rook_occupancy >>= 64 - rook_relevant_bits[square]

	// return rook attacks
	return bishop_attacks[square][bishop_occupancy] | rook_attacks[square][rook_occupancy]
}

func (b Board) isWhiteSquareAttacked(square square) bool {
	return (pawn_attacks[black][square]&b.bitboards[P] > 0) ||
		(knight_attacks[square]&b.bitboards[N] > 0) ||
		(b.get_bishop_attacks(square)&b.bitboards[B] > 0) ||
		(b.get_rook_attacks(square)&b.bitboards[R] > 0) ||
		(b.get_queen_attacks(square)&b.bitboards[Q] > 0) ||
		(king_attacks[square]&b.bitboards[K] > 0)
}

func (bb Board) isBlackSquareAttacked(square square) bool {
	return (pawn_attacks[white][square]&bb.bitboards[p] > 0) ||
		(knight_attacks[square]&bb.bitboards[n] > 0) ||
		(bb.get_bishop_attacks(square)&bb.bitboards[b] > 0) ||
		(bb.get_rook_attacks(square)&bb.bitboards[r] > 0) ||
		(bb.get_queen_attacks(square)&bb.bitboards[q] > 0) ||
		(king_attacks[square]&bb.bitboards[k] > 0)
}

func (b *Board) isSquareAttacked(square square, side side) bool {
	if side == white {
		return b.isWhiteSquareAttacked(square)
	} else {
		return b.isBlackSquareAttacked(square)
	}
}

func (b *Board) isOccupied(square square) bool {
	return b.occupancies[both].getBit(square)
}

func (b *Board) isEmpty(square square) bool {
	return !b.occupancies[both].getBit(square)
}

// generate all moves
func (board Board) generateMoves(moves *Moves) {
	// define source & target squares
	var source_square, target_square square

	// define current piece's bitboard copy & it's attacks
	var bitboard, attacks Bitboard

	if board.side == white {
		// genarate pawn moves
		bitboard = board.bitboards[P]

		// loop over white pawns within white pawn bitboard
		for ; bitboard > 0; bitboard.popBit(source_square) {
			// init source square
			source_square = square(bitboard.getLs1bIndex())

			// init target square
			target_square = source_square - 8

			// generate quiet pawn moves
			if !(target_square < a8) && board.isEmpty(target_square) {
				// pawn promotion
				if source_square >= a7 && source_square <= h7 {
					moves.addMove(encode_move(source_square, target_square, P, Q, 0, 0, 0, 0))
					moves.addMove(encode_move(source_square, target_square, P, R, 0, 0, 0, 0))
					moves.addMove(encode_move(source_square, target_square, P, B, 0, 0, 0, 0))
					moves.addMove(encode_move(source_square, target_square, P, N, 0, 0, 0, 0))
				} else {
					// one square ahead pawn move
					moves.addMove(encode_move(source_square, target_square, P, no_piece, 0, 0, 0, 0))

					// two squares ahead pawn move
					if (source_square >= a2 && source_square <= h2) && !board.occupancies[both].getBit(target_square-8) {
						moves.addMove(encode_move(source_square, target_square-8, P, no_piece, 0, 1, 0, 0))
					}
				}
			}

			// generate pawn captures
			for attacks = pawn_attacks[board.side][source_square] & board.occupancies[black]; attacks > 0; attacks.popBit(target_square) {
				// init target square
				target_square = square(attacks.getLs1bIndex())

				// pawn promotion
				if source_square >= a7 && source_square <= h7 {
					moves.addMove(encode_move(source_square, target_square, P, Q, 1, 0, 0, 0))
					moves.addMove(encode_move(source_square, target_square, P, R, 1, 0, 0, 0))
					moves.addMove(encode_move(source_square, target_square, P, B, 1, 0, 0, 0))
					moves.addMove(encode_move(source_square, target_square, P, N, 1, 0, 0, 0))
				} else {
					// one square ahead pawn move
					moves.addMove(encode_move(source_square, target_square, P, no_piece, 1, 0, 0, 0))
				}
			}

			// generate enpassant captures
			if board.enpassant != no_sq {
				// lookup pawn attacks and bitwise AND with enpassant square (bit)
				enpassant_attacks := pawn_attacks[board.side][source_square] & (Bitboard(1) << board.enpassant)

				// make sure enpassant capture available
				if enpassant_attacks > 0 {
					// init enpassant capture target square
					target_enpassant := square(enpassant_attacks.getLs1bIndex())
					moves.addMove(encode_move(source_square, target_enpassant, P, no_piece, 1, 0, 1, 0))
				}
			}
		}

		// genarate castling moves
		bitboard = board.bitboards[K]

		// king side castling is available
		if board.castle&wk > 0 {
			// make sure square between king and king's rook are empty
			if board.isEmpty(f1) && board.isEmpty(g1) {
				// make sure king and the f1 squares are not under attacks
				if !board.isBlackSquareAttacked(e1) && !board.isBlackSquareAttacked(f1) {
					moves.addMove(encode_move(e1, g1, K, no_piece, 0, 0, 0, 1))
				}
			}
		}

		// queen side castling is available
		if board.castle&wq > 0 {
			// make sure square between king and queen's rook are empty
			if board.isEmpty(d1) && board.isEmpty(c1) && board.isEmpty(b1) {
				// make sure king and the d1 squares are not under attacks
				if !board.isBlackSquareAttacked(e1) && !board.isBlackSquareAttacked(d1) {
					moves.addMove(encode_move(e1, c1, K, no_piece, 0, 0, 0, 1))
				}
			}
		}
	} else {
		// genarate pawn moves
		bitboard = board.bitboards[p]

		// loop over white pawns within white pawn bitboard
		for ; bitboard > 0; bitboard.popBit(source_square) {
			// init source square
			source_square = square(bitboard.getLs1bIndex())

			// init target square
			target_square = source_square + 8

			// generate quiet pawn moves
			if !(target_square > h1) && board.isEmpty(target_square) {
				// pawn promotion
				if source_square >= a2 && source_square <= h2 {
					moves.addMove(encode_move(source_square, target_square, p, q, 0, 0, 0, 0))
					moves.addMove(encode_move(source_square, target_square, p, r, 0, 0, 0, 0))
					moves.addMove(encode_move(source_square, target_square, p, b, 0, 0, 0, 0))
					moves.addMove(encode_move(source_square, target_square, p, n, 0, 0, 0, 0))
				} else {
					// one square ahead pawn move
					moves.addMove(encode_move(source_square, target_square, p, no_piece, 0, 0, 0, 0))

					// two squares ahead pawn move
					if (source_square >= a7 && source_square <= h7) && !board.occupancies[both].getBit(target_square+8) {
						moves.addMove(encode_move(source_square, target_square+8, p, no_piece, 0, 1, 0, 0))
					}
				}
			}

			// generate pawn captures
			for attacks = pawn_attacks[board.side][source_square] & board.occupancies[white]; attacks > 0; attacks.popBit(target_square) {
				// init target square
				target_square = square(attacks.getLs1bIndex())

				// pawn promotion
				if source_square >= a2 && source_square <= h2 {
					moves.addMove(encode_move(source_square, target_square, p, q, 1, 0, 0, 0))
					moves.addMove(encode_move(source_square, target_square, p, r, 1, 0, 0, 0))
					moves.addMove(encode_move(source_square, target_square, p, b, 1, 0, 0, 0))
					moves.addMove(encode_move(source_square, target_square, p, n, 1, 0, 0, 0))
				} else {
					// one square ahead pawn move
					moves.addMove(encode_move(source_square, target_square, p, no_piece, 1, 0, 0, 0))
				}
			}

			// generate enpassant captures
			if board.enpassant != no_sq {
				// lookup pawn attacks and bitwise AND with enpassant square (bit)
				enpassant_attacks := pawn_attacks[board.side][source_square] & (Bitboard(1) << board.enpassant)

				// make sure enpassant capture available
				if enpassant_attacks > 0 {
					// init enpassant capture target square
					target_enpassant := square(enpassant_attacks.getLs1bIndex())
					moves.addMove(encode_move(source_square, target_enpassant, p, no_piece, 1, 0, 1, 0))
				}
			}
		}

		// genarate castling moves
		bitboard = board.bitboards[k]

		// king side castling is available
		if board.castle&bk > 0 {
			// make sure square between king and king's rook are empty
			if board.isEmpty(f8) && board.isEmpty(g8) {
				// make sure king and the f8 squares are not under attacks
				if !board.isWhiteSquareAttacked(e8) && !board.isWhiteSquareAttacked(f8) {
					moves.addMove(encode_move(e8, g8, k, no_piece, 0, 0, 0, 1))
				}
			}
		}

		// queen side castling is available
		if board.castle&bq > 0 {
			// make sure square between king and queen's rook are empty
			if board.isEmpty(d8) && board.isEmpty(c8) && board.isEmpty(b8) {
				// make sure king and the d8 squares are not under attacks
				if !board.isWhiteSquareAttacked(e8) && !board.isWhiteSquareAttacked(d8) {
					moves.addMove(encode_move(e8, c8, k, no_piece, 0, 0, 0, 1))
				}
			}
		}
	}

	// Init curPiece
	curPiece := N
	if board.side == black {
		curPiece = n
	}

	// genarate knight moves
	bitboard = board.bitboards[curPiece]

	// loop over source squares of piece bitboard copy
	for ; bitboard > 0; bitboard.popBit(source_square) {
		// init source square
		source_square = square(bitboard.getLs1bIndex())

		// // loop over target squares available from generated attacks
		for attacks = knight_attacks[source_square] & (^board.occupancies[board.side]); attacks > 0; attacks.popBit(target_square) {
			// init target square
			target_square = square(attacks.getLs1bIndex())

			// quiet move
			if !board.occupancies[board.side.opposite()].getBit(target_square) {
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 0, 0, 0, 0))
			} else { // capture move
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 1, 0, 0, 0))
			}
		}
	}

	// generate bishop moves
	curPiece++
	bitboard = board.bitboards[curPiece]

	// loop over source squares of piece bitboard copy
	for ; bitboard > 0; bitboard.popBit(source_square) {
		// init source square
		source_square = square(bitboard.getLs1bIndex())

		// loop over target squares available from generated attacks
		for attacks = board.get_bishop_attacks(source_square) & (^board.occupancies[board.side]); attacks > 0; attacks.popBit(target_square) {
			// init target square
			target_square = square(attacks.getLs1bIndex())

			// quiet move
			if !board.occupancies[board.side.opposite()].getBit(target_square) {
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 0, 0, 0, 0))
			} else { // capture move
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 1, 0, 0, 0))
			}
		}
	}

	// generate rook moves
	curPiece++
	bitboard = board.bitboards[curPiece]

	// loop over source squares of piece bitboard copy
	for ; bitboard > 0; bitboard.popBit(source_square) {
		// init source square
		source_square = square(bitboard.getLs1bIndex())

		// loop over target squares available from generated attacks
		for attacks = board.get_rook_attacks(source_square) & (^board.occupancies[board.side]); attacks > 0; attacks.popBit(target_square) {
			// init target square
			target_square = square(attacks.getLs1bIndex())

			// quiet move
			if !board.occupancies[board.side.opposite()].getBit(target_square) {
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 0, 0, 0, 0))
			} else { // capture move
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 1, 0, 0, 0))
			}
		}
	}

	// generate queen moves
	curPiece++
	bitboard = board.bitboards[curPiece]

	// loop over source squares of piece bitboard copy
	for ; bitboard > 0; bitboard.popBit(source_square) {
		// init source square
		source_square = square(bitboard.getLs1bIndex())

		// loop over target squares available from generated attacks
		for attacks = board.get_queen_attacks(source_square) & (^board.occupancies[board.side]); attacks > 0; attacks.popBit(target_square) {
			// init target square
			target_square = square(attacks.getLs1bIndex())

			// quiet move
			if !board.occupancies[board.side.opposite()].getBit(target_square) {
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 0, 0, 0, 0))
			} else { // capture move
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 1, 0, 0, 0))
			}
		}
	}

	// genarate king moves
	curPiece++
	bitboard = board.bitboards[curPiece]

	// loop over source squares of piece bitboard copy
	for ; bitboard > 0; bitboard.popBit(source_square) {
		// init source square
		source_square = square(bitboard.getLs1bIndex())

		// init piece attacks in order to get set of target squares

		// // loop over target squares available from generated attacks
		for attacks = king_attacks[source_square] & (^board.occupancies[board.side]); attacks > 0; attacks.popBit(target_square) {
			// init target square
			target_square = square(attacks.getLs1bIndex())

			// quiet move
			if !board.occupancies[board.side.opposite()].getBit(target_square) {
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 0, 0, 0, 0))
			} else { // capture move
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 1, 0, 0, 0))
			}
		}
	}
}

// generate all moves
func (board *Board) makeMove(move Move, capturesOnly bool) bool {

	if !capturesOnly { // All moves
		// Save the board
		boardCopy := *board

		// Parse move
		sourceSquare := move.getSource()
		targrtSquare := move.getTarget()
		curPiece := move.getPiece()
		promotionPiece := move.getPromotionPiece()
		capture := move.isCapture()
		double := move.isDoublePawnPush()
		enpassant := move.isEnpassant()
		castling := move.isCastling()

		// Move piece
		board.bitboards[curPiece].popBit(sourceSquare)
		board.bitboards[curPiece].setBit(targrtSquare)

		// Handling capture moves
		if capture {
			startPiece, endPiece := P, K
			if board.side == white {
				startPiece, endPiece = p, k
			}

			for bb_piece := startPiece; bb_piece <= endPiece; bb_piece++ {
				board.bitboards[bb_piece].popBit(targrtSquare)
			}
		}

		// Handling pawn promotions
		if promotionPiece != no_piece {
			board.bitboards[curPiece].popBit(targrtSquare)
			board.bitboards[promotionPiece].setBit(targrtSquare)
		}

		// Handling enpassant captures
		if enpassant {
			if board.side == white {
				board.bitboards[p].popBit(targrtSquare + 8)
			} else {
				board.bitboards[P].popBit(targrtSquare - 8)
			}
		}

		// Reset enpassant square
		board.enpassant = no_sq

		// Handle double pawn push
		if double {
			if board.side == white {
				board.enpassant = targrtSquare + 8
			} else {
				board.enpassant = targrtSquare - 8
			}
		}

		// Handle castling
		if castling {
			switch targrtSquare {
			case g1:
				board.bitboards[R].popBit(h1)
				board.bitboards[R].setBit(f1)
			case c1:
				board.bitboards[R].popBit(a1)
				board.bitboards[R].setBit(d1)
			case g8:
				board.bitboards[r].popBit(h8)
				board.bitboards[r].setBit(f8)
			case c8:
				board.bitboards[r].popBit(a8)
				board.bitboards[r].setBit(d8)
			}
		}

		// Update castling rights
		board.castle &= castlingRights[sourceSquare]
		board.castle &= castlingRights[targrtSquare]

		// Update occupancies
		board.occupancies[white] = board.bitboards[P] | board.bitboards[N] | board.bitboards[B] | board.bitboards[R] | board.bitboards[Q] | board.bitboards[K]
		board.occupancies[black] = board.bitboards[p] | board.bitboards[n] | board.bitboards[b] | board.bitboards[r] | board.bitboards[q] | board.bitboards[k]
		board.occupancies[both] = board.occupancies[white] | board.occupancies[black]

		king := k
		if board.side == white {
			king = K
		}
		// change side
		board.side = board.side.opposite()

		// Make sure king is not in check
		if board.isSquareAttacked(square(board.bitboards[king].getLs1bIndex()), board.side) {
			*board = boardCopy

			return false
		}

		return true
	} else if move.isCapture() { // Capture moves only
		return board.makeMove(move, false)
	}

	return false
}
