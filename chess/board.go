package chess

import "strconv"

type Board struct {
	Bitboards   [12]Bitboard
	Occupancies [3]Bitboard
	Side        Side
	Enpassant   square
	Castle      castle
	HashKey     Bitboard
	Fifty       int
}

// parse FEN string
func NewBoadFromFen(fen []byte) *Board {
	var b Board

	// reset game state variables
	b.Enpassant = no_sq

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
				b.Bitboards[piece].setBit(square)

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
					if b.Bitboards[bb_piece].getBit(square) {
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
		b.Side = white
	} else {
		b.Side = black
	}

	// go to parsing castling rights
	fen = fen[2:]

	// parse castling rights
	for fen[0] != ' ' {
		switch fen[0] {
		case 'K':
			b.Castle |= wk
		case 'Q':
			b.Castle |= wq
		case 'k':
			b.Castle |= bk
		case 'q':
			b.Castle |= bq
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
		b.Enpassant = square(rank*8 + file)
	} else { // no enpassant square
		b.Enpassant = no_sq
	}

	// go to parsing half move counter (increment pointer to FEN string)
	fen = fen[1:]

	// parse half move counter to init fifty move counter
	halfmove, err := strconv.Atoi(string(fen))
	if err == nil {
		b.Fifty = halfmove
	}

	// loop over white pieces bitboards
	for piece := P; piece <= K; piece++ {
		// populate white occupancy bitboard
		b.Occupancies[white] |= b.Bitboards[piece]
	}
	// loop over black pieces bitboards
	for piece := p; piece <= k; piece++ {
		// populate white occupancy bitboard
		b.Occupancies[black] |= b.Bitboards[piece]
	}

	// init all occupancies
	b.Occupancies[both] = b.Occupancies[white] | b.Occupancies[black]

	b.HashKey = b.generate_hash_key()

	return &b
}

// get bishop attacks
func (b *Board) get_bishop_attacks(square square) Bitboard {
	// get bishop attacks assuming current board occupancy
	occupancy := b.Occupancies[both]

	occupancy &= bishop_masks[square]
	occupancy *= bishop_magic_numbers[square]
	occupancy >>= 64 - bishop_relevant_bits[square]

	// return bishop attacks
	return bishop_attacks[square][occupancy]
}

// get rook attacks
func (b *Board) get_rook_attacks(square square) Bitboard {
	// get bishop attacks assuming current board occupancy
	occupancy := b.Occupancies[both]

	occupancy &= rook_masks[square]
	occupancy *= rook_magic_numbers[square]
	occupancy >>= 64 - rook_relevant_bits[square]

	// return rook attacks
	return rook_attacks[square][occupancy]
}

// get queen attacks
func (b *Board) get_queen_attacks(square square) Bitboard {
	// get bishop attacks assuming current board occupancy
	bishop_occupancy := b.Occupancies[both]
	rook_occupancy := b.Occupancies[both]

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
	return (pawn_attacks[black][square]&b.Bitboards[P] > 0) ||
		(knight_attacks[square]&b.Bitboards[N] > 0) ||
		(b.get_bishop_attacks(square)&b.Bitboards[B] > 0) ||
		(b.get_rook_attacks(square)&b.Bitboards[R] > 0) ||
		(b.get_queen_attacks(square)&b.Bitboards[Q] > 0) ||
		(king_attacks[square]&b.Bitboards[K] > 0)
}

func (bb Board) isBlackSquareAttacked(square square) bool {
	return (pawn_attacks[white][square]&bb.Bitboards[p] > 0) ||
		(knight_attacks[square]&bb.Bitboards[n] > 0) ||
		(bb.get_bishop_attacks(square)&bb.Bitboards[b] > 0) ||
		(bb.get_rook_attacks(square)&bb.Bitboards[r] > 0) ||
		(bb.get_queen_attacks(square)&bb.Bitboards[q] > 0) ||
		(king_attacks[square]&bb.Bitboards[k] > 0)
}

func (b *Board) isSquareAttacked(square square, side Side) bool {
	if side == white {
		return b.isWhiteSquareAttacked(square)
	} else {
		return b.isBlackSquareAttacked(square)
	}
}

func (b *Board) isOccupied(square square) bool {
	return b.Occupancies[both].getBit(square)
}

func (b *Board) isEmpty(square square) bool {
	return !b.Occupancies[both].getBit(square)
}

func (board *Board) generateCaptureMoves(moves *Moves) {
	// define source & target squares
	var source_square, target_square square

	// define current piece's bitboard copy & it's attacks
	var bitboard, attacks Bitboard

	if board.Side == white {
		// genarate pawn moves
		bitboard = board.Bitboards[P]

		// loop over white pawns within white pawn bitboard
		for ; bitboard > 0; bitboard.popBit(source_square) {
			// init source square
			source_square = square(bitboard.GetLs1bIndex())

			// init target square
			target_square = source_square - 8

			// generate pawn captures
			for attacks = pawn_attacks[board.Side][source_square] & board.Occupancies[black]; attacks > 0; attacks.popBit(target_square) {
				// init target square
				target_square = square(attacks.GetLs1bIndex())

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
			if board.Enpassant != no_sq {
				// lookup pawn attacks and bitwise AND with enpassant square (bit)
				enpassant_attacks := pawn_attacks[board.Side][source_square] & (Bitboard(1) << board.Enpassant)

				// make sure enpassant capture available
				if enpassant_attacks > 0 {
					// init enpassant capture target square
					target_enpassant := square(enpassant_attacks.GetLs1bIndex())
					moves.addMove(encode_move(source_square, target_enpassant, P, no_piece, 1, 0, 1, 0))
				}
			}
		}
	} else {
		// genarate pawn moves
		bitboard = board.Bitboards[p]

		// loop over white pawns within white pawn bitboard
		for ; bitboard > 0; bitboard.popBit(source_square) {
			// init source square
			source_square = square(bitboard.GetLs1bIndex())

			// init target square
			target_square = source_square + 8

			// generate pawn captures
			for attacks = pawn_attacks[board.Side][source_square] & board.Occupancies[white]; attacks > 0; attacks.popBit(target_square) {
				// init target square
				target_square = square(attacks.GetLs1bIndex())

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
			if board.Enpassant != no_sq {
				// lookup pawn attacks and bitwise AND with enpassant square (bit)
				enpassant_attacks := pawn_attacks[board.Side][source_square] & (Bitboard(1) << board.Enpassant)

				// make sure enpassant capture available
				if enpassant_attacks > 0 {
					// init enpassant capture target square
					target_enpassant := square(enpassant_attacks.GetLs1bIndex())
					moves.addMove(encode_move(source_square, target_enpassant, p, no_piece, 1, 0, 1, 0))
				}
			}
		}
	}

	// Init curPiece
	curPiece := N
	if board.Side == black {
		curPiece = n
	}

	// genarate knight moves
	bitboard = board.Bitboards[curPiece]

	// loop over source squares of piece bitboard copy
	for ; bitboard > 0; bitboard.popBit(source_square) {
		// init source square
		source_square = square(bitboard.GetLs1bIndex())

		// // loop over target squares available from generated attacks
		for attacks = knight_attacks[source_square] & (^board.Occupancies[board.Side]); attacks > 0; attacks.popBit(target_square) {
			// init target square
			target_square = square(attacks.GetLs1bIndex())

			if board.Occupancies[board.Side.opposite()].getBit(target_square) { // capture move
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 1, 0, 0, 0))
			}
		}
	}

	// generate bishop moves
	curPiece++
	bitboard = board.Bitboards[curPiece]

	// loop over source squares of piece bitboard copy
	for ; bitboard > 0; bitboard.popBit(source_square) {
		// init source square
		source_square = square(bitboard.GetLs1bIndex())

		// loop over target squares available from generated attacks
		for attacks = board.get_bishop_attacks(source_square) & (^board.Occupancies[board.Side]); attacks > 0; attacks.popBit(target_square) {
			// init target square
			target_square = square(attacks.GetLs1bIndex())

			if board.Occupancies[board.Side.opposite()].getBit(target_square) { // capture move
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 1, 0, 0, 0))
			}
		}
	}

	// generate rook moves
	curPiece++
	bitboard = board.Bitboards[curPiece]

	// loop over source squares of piece bitboard copy
	for ; bitboard > 0; bitboard.popBit(source_square) {
		// init source square
		source_square = square(bitboard.GetLs1bIndex())

		// loop over target squares available from generated attacks
		for attacks = board.get_rook_attacks(source_square) & (^board.Occupancies[board.Side]); attacks > 0; attacks.popBit(target_square) {
			// init target square
			target_square = square(attacks.GetLs1bIndex())

			if board.Occupancies[board.Side.opposite()].getBit(target_square) { // capture move
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 1, 0, 0, 0))
			}
		}
	}

	// generate queen moves
	curPiece++
	bitboard = board.Bitboards[curPiece]

	// loop over source squares of piece bitboard copy
	for ; bitboard > 0; bitboard.popBit(source_square) {
		// init source square
		source_square = square(bitboard.GetLs1bIndex())

		// loop over target squares available from generated attacks
		for attacks = board.get_queen_attacks(source_square) & (^board.Occupancies[board.Side]); attacks > 0; attacks.popBit(target_square) {
			// init target square
			target_square = square(attacks.GetLs1bIndex())

			if board.Occupancies[board.Side.opposite()].getBit(target_square) { // capture move
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 1, 0, 0, 0))
			}
		}
	}

	// genarate king moves
	curPiece++
	bitboard = board.Bitboards[curPiece]

	// loop over source squares of piece bitboard copy
	for ; bitboard > 0; bitboard.popBit(source_square) {
		// init source square
		source_square = square(bitboard.GetLs1bIndex())

		// init piece attacks in order to get set of target squares

		// // loop over target squares available from generated attacks
		for attacks = king_attacks[source_square] & (^board.Occupancies[board.Side]); attacks > 0; attacks.popBit(target_square) {
			// init target square
			target_square = square(attacks.GetLs1bIndex())

			if board.Occupancies[board.Side.opposite()].getBit(target_square) { // capture move
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 1, 0, 0, 0))
			}
		}
	}
}

// generate all moves
func (board *Board) generateMoves(moves *Moves) {
	// define source & target squares
	var source_square, target_square square

	// define current piece's bitboard copy & it's attacks
	var bitboard, attacks Bitboard

	if board.Side == white {
		// genarate pawn moves
		bitboard = board.Bitboards[P]

		// loop over white pawns within white pawn bitboard
		for ; bitboard > 0; bitboard.popBit(source_square) {
			// init source square
			source_square = square(bitboard.GetLs1bIndex())

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
					if (source_square >= a2 && source_square <= h2) && !board.Occupancies[both].getBit(target_square-8) {
						moves.addMove(encode_move(source_square, target_square-8, P, no_piece, 0, 1, 0, 0))
					}
				}
			}

			// generate pawn captures
			for attacks = pawn_attacks[board.Side][source_square] & board.Occupancies[black]; attacks > 0; attacks.popBit(target_square) {
				// init target square
				target_square = square(attacks.GetLs1bIndex())

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
			if board.Enpassant != no_sq {
				// lookup pawn attacks and bitwise AND with enpassant square (bit)
				enpassant_attacks := pawn_attacks[board.Side][source_square] & (Bitboard(1) << board.Enpassant)

				// make sure enpassant capture available
				if enpassant_attacks > 0 {
					// init enpassant capture target square
					target_enpassant := square(enpassant_attacks.GetLs1bIndex())
					moves.addMove(encode_move(source_square, target_enpassant, P, no_piece, 1, 0, 1, 0))
				}
			}
		}

		// genarate castling moves
		bitboard = board.Bitboards[K]

		// king side castling is available
		if board.Castle&wk > 0 {
			// make sure square between king and king's rook are empty
			if board.isEmpty(f1) && board.isEmpty(g1) {
				// make sure king and the f1 squares are not under attacks
				if !board.isBlackSquareAttacked(e1) && !board.isBlackSquareAttacked(f1) {
					moves.addMove(encode_move(e1, g1, K, no_piece, 0, 0, 0, 1))
				}
			}
		}

		// queen side castling is available
		if board.Castle&wq > 0 {
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
		bitboard = board.Bitboards[p]

		// loop over white pawns within white pawn bitboard
		for ; bitboard > 0; bitboard.popBit(source_square) {
			// init source square
			source_square = square(bitboard.GetLs1bIndex())

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
					if (source_square >= a7 && source_square <= h7) && !board.Occupancies[both].getBit(target_square+8) {
						moves.addMove(encode_move(source_square, target_square+8, p, no_piece, 0, 1, 0, 0))
					}
				}
			}

			// generate pawn captures
			for attacks = pawn_attacks[board.Side][source_square] & board.Occupancies[white]; attacks > 0; attacks.popBit(target_square) {
				// init target square
				target_square = square(attacks.GetLs1bIndex())

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
			if board.Enpassant != no_sq {
				// lookup pawn attacks and bitwise AND with enpassant square (bit)
				enpassant_attacks := pawn_attacks[board.Side][source_square] & (Bitboard(1) << board.Enpassant)

				// make sure enpassant capture available
				if enpassant_attacks > 0 {
					// init enpassant capture target square
					target_enpassant := square(enpassant_attacks.GetLs1bIndex())
					moves.addMove(encode_move(source_square, target_enpassant, p, no_piece, 1, 0, 1, 0))
				}
			}
		}

		// genarate castling moves
		bitboard = board.Bitboards[k]

		// king side castling is available
		if board.Castle&bk > 0 {
			// make sure square between king and king's rook are empty
			if board.isEmpty(f8) && board.isEmpty(g8) {
				// make sure king and the f8 squares are not under attacks
				if !board.isWhiteSquareAttacked(e8) && !board.isWhiteSquareAttacked(f8) {
					moves.addMove(encode_move(e8, g8, k, no_piece, 0, 0, 0, 1))
				}
			}
		}

		// queen side castling is available
		if board.Castle&bq > 0 {
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
	if board.Side == black {
		curPiece = n
	}

	// genarate knight moves
	bitboard = board.Bitboards[curPiece]

	// loop over source squares of piece bitboard copy
	for ; bitboard > 0; bitboard.popBit(source_square) {
		// init source square
		source_square = square(bitboard.GetLs1bIndex())

		// // loop over target squares available from generated attacks
		for attacks = knight_attacks[source_square] & (^board.Occupancies[board.Side]); attacks > 0; attacks.popBit(target_square) {
			// init target square
			target_square = square(attacks.GetLs1bIndex())

			// quiet move
			if !board.Occupancies[board.Side.opposite()].getBit(target_square) {
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 0, 0, 0, 0))
			} else { // capture move
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 1, 0, 0, 0))
			}
		}
	}

	// generate bishop moves
	curPiece++
	bitboard = board.Bitboards[curPiece]

	// loop over source squares of piece bitboard copy
	for ; bitboard > 0; bitboard.popBit(source_square) {
		// init source square
		source_square = square(bitboard.GetLs1bIndex())

		// loop over target squares available from generated attacks
		for attacks = board.get_bishop_attacks(source_square) & (^board.Occupancies[board.Side]); attacks > 0; attacks.popBit(target_square) {
			// init target square
			target_square = square(attacks.GetLs1bIndex())

			// quiet move
			if !board.Occupancies[board.Side.opposite()].getBit(target_square) {
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 0, 0, 0, 0))
			} else { // capture move
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 1, 0, 0, 0))
			}
		}
	}

	// generate rook moves
	curPiece++
	bitboard = board.Bitboards[curPiece]

	// loop over source squares of piece bitboard copy
	for ; bitboard > 0; bitboard.popBit(source_square) {
		// init source square
		source_square = square(bitboard.GetLs1bIndex())

		// loop over target squares available from generated attacks
		for attacks = board.get_rook_attacks(source_square) & (^board.Occupancies[board.Side]); attacks > 0; attacks.popBit(target_square) {
			// init target square
			target_square = square(attacks.GetLs1bIndex())

			// quiet move
			if !board.Occupancies[board.Side.opposite()].getBit(target_square) {
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 0, 0, 0, 0))
			} else { // capture move
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 1, 0, 0, 0))
			}
		}
	}

	// generate queen moves
	curPiece++
	bitboard = board.Bitboards[curPiece]

	// loop over source squares of piece bitboard copy
	for ; bitboard > 0; bitboard.popBit(source_square) {
		// init source square
		source_square = square(bitboard.GetLs1bIndex())

		// loop over target squares available from generated attacks
		for attacks = board.get_queen_attacks(source_square) & (^board.Occupancies[board.Side]); attacks > 0; attacks.popBit(target_square) {
			// init target square
			target_square = square(attacks.GetLs1bIndex())

			// quiet move
			if !board.Occupancies[board.Side.opposite()].getBit(target_square) {
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 0, 0, 0, 0))
			} else { // capture move
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 1, 0, 0, 0))
			}
		}
	}

	// genarate king moves
	curPiece++
	bitboard = board.Bitboards[curPiece]

	// loop over source squares of piece bitboard copy
	for ; bitboard > 0; bitboard.popBit(source_square) {
		// init source square
		source_square = square(bitboard.GetLs1bIndex())

		// init piece attacks in order to get set of target squares

		// // loop over target squares available from generated attacks
		for attacks = king_attacks[source_square] & (^board.Occupancies[board.Side]); attacks > 0; attacks.popBit(target_square) {
			// init target square
			target_square = square(attacks.GetLs1bIndex())

			// quiet move
			if !board.Occupancies[board.Side.opposite()].getBit(target_square) {
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 0, 0, 0, 0))
			} else { // capture move
				moves.addMove(encode_move(source_square, target_square, curPiece, no_piece, 1, 0, 0, 0))
			}
		}
	}
}

// Tries to make a move
func (board *Board) MakeMove(move Move, capturesOnly bool) bool {

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
		board.Bitboards[curPiece].popBit(sourceSquare)
		board.Bitboards[curPiece].setBit(targrtSquare)

		// hash piece
		board.HashKey ^= piece_keys[curPiece][sourceSquare] // remove piece from source square in hash key
		board.HashKey ^= piece_keys[curPiece][targrtSquare] // set piece to the target square in hash key

		// increment fifty move rule counter
		board.Fifty++

		// if pawn moved
		if curPiece == P || curPiece == p {
			// reset fifty move rule counter
			board.Fifty = 0
		}

		// Handling capture moves
		if capture {
			// reset fifty move rule counter
			board.Fifty = 0

			startPiece, endPiece := P, K
			if board.Side == white {
				startPiece, endPiece = p, k
			}

			for bb_piece := startPiece; bb_piece <= endPiece; bb_piece++ {
				if board.Bitboards[bb_piece].getBit(targrtSquare) {
					board.Bitboards[bb_piece].popBit(targrtSquare)

					// remove the piece from hash key
					board.HashKey ^= piece_keys[bb_piece][targrtSquare]
					break
				}
			}
		}

		// Handling pawn promotions
		if promotionPiece != no_piece {
			board.Bitboards[curPiece].popBit(targrtSquare)
			board.HashKey ^= piece_keys[curPiece][targrtSquare]

			board.Bitboards[promotionPiece].setBit(targrtSquare)
			board.HashKey ^= piece_keys[promotionPiece][targrtSquare]
		}

		// Handling enpassant captures
		if enpassant {
			if board.Side == white {
				board.Bitboards[p].popBit(targrtSquare + 8)
				board.HashKey ^= piece_keys[p][targrtSquare+8]
			} else {
				board.Bitboards[P].popBit(targrtSquare - 8)
				board.HashKey ^= piece_keys[P][targrtSquare-8]
			}
		}

		// hash enpassant if available (remove enpassant square from hash key )
		if board.Enpassant != no_sq {
			board.HashKey ^= enpassant_keys[board.Enpassant]
		}

		// Reset enpassant square
		board.Enpassant = no_sq

		// Handle double pawn push
		if double {
			if board.Side == white {
				board.Enpassant = targrtSquare + 8

				// hash enpassant
				board.HashKey ^= enpassant_keys[targrtSquare+8]
			} else {
				board.Enpassant = targrtSquare - 8

				// hash enpassant
				board.HashKey ^= enpassant_keys[targrtSquare-8]
			}
		}

		// Handle castling
		if castling {
			switch targrtSquare {
			case g1:
				board.Bitboards[R].popBit(h1)
				board.Bitboards[R].setBit(f1)

				// hash rook
				board.HashKey ^= piece_keys[R][h1] // remove rook from h1 from hash key
				board.HashKey ^= piece_keys[R][f1] // put rook on f1 into a hash key
			case c1:
				board.Bitboards[R].popBit(a1)
				board.Bitboards[R].setBit(d1)

				// hash rook
				board.HashKey ^= piece_keys[R][a1] // remove rook from a1 from hash key
				board.HashKey ^= piece_keys[R][d1] // put rook on d1 into a hash key
			case g8:
				board.Bitboards[r].popBit(h8)
				board.Bitboards[r].setBit(f8)

				// hash rook
				board.HashKey ^= piece_keys[r][h8] // remove rook from h8 from hash key
				board.HashKey ^= piece_keys[r][f8] // put rook on f8 into a hash key
			case c8:
				board.Bitboards[r].popBit(a8)
				board.Bitboards[r].setBit(d8)

				// hash rook
				board.HashKey ^= piece_keys[r][a8] // remove rook from a8 from hash key
				board.HashKey ^= piece_keys[r][d8] // put rook on d8 into a hash key
			}
		}

		// hash castling
		board.HashKey ^= castle_keys[board.Castle]

		// Update castling rights
		board.Castle &= castlingRights[sourceSquare]
		board.Castle &= castlingRights[targrtSquare]

		// hash castling
		board.HashKey ^= castle_keys[board.Castle]

		// Update occupancies
		board.Occupancies[white] = board.Bitboards[P] | board.Bitboards[N] | board.Bitboards[B] | board.Bitboards[R] | board.Bitboards[Q] | board.Bitboards[K]
		board.Occupancies[black] = board.Bitboards[p] | board.Bitboards[n] | board.Bitboards[b] | board.Bitboards[r] | board.Bitboards[q] | board.Bitboards[k]
		board.Occupancies[both] = board.Occupancies[white] | board.Occupancies[black]

		king := k
		if board.Side == white {
			king = K
		}

		// change side
		board.Side = board.Side.opposite()

		// hash side
		board.HashKey ^= side_key

		// Make sure king is not in check
		if board.isSquareAttacked(square(board.Bitboards[king].GetLs1bIndex()), board.Side) {
			*board = boardCopy

			return false
		}

		return true
	} else if move.isCapture() { // Capture moves only
		return board.MakeMove(move, false)
	}

	return false
}

func (board *Board) GetLegalMoves(capturesOnly bool) *Moves {
	var moves Moves

	board.generateMoves(&moves)

	legalCount := 0
	for i := 0; i < moves.count; i++ {
		boardCopy := *board

		if board.MakeMove(moves.moves[i], capturesOnly) {
			moves.moves[legalCount] = moves.moves[i]
			legalCount++
			*board = boardCopy
		}
	}

	moves.count = legalCount

	return &moves
}

func (m *Moves) FilterSelected(source int) *Moves {
	filteredCount := 0
	for i := 0; i < m.count; i++ {
		if m.moves[i].getSource() == square(source) {
			m.moves[filteredCount] = m.moves[i]
			filteredCount++
		}
	}

	m.count = filteredCount

	return m
}

func (board *Board) Gameover() bool {
	return board.GetLegalMoves(false).count == 0
}

func Startpos() Board {
	return *NewBoadFromFen(STARTPOS_FEN)
}
