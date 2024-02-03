package chess

// material score [game phase][piece]
var materialScore = [2][12]int{
	// opening material score
	{82, 337, 365, 477, 1025, 12000, -82, -337, -365, -477, -1025, -12000},

	// endgame material score
	{94, 281, 297, 512, 936, 12000, -94, -281, -297, -512, -936, -12000},
}

// game phase scores
const openingPhaseScore = 6192
const endgamePhaseScore = 518

// game phases
const (
	opening    = 0
	endgame    = 1
	middlegame = 2
)

// piece types
const (
	PAWN   = 0
	KNIGHT = 1
	BISHOP = 2
	ROOK   = 3
	QUEEN  = 4
	KING   = 5
)

// positional piece scores [game phase][piece][square]
var positionalScore = [2][6][64]int{
	// Opening positional piece scores //
	{
		//pawn
		[64]int{
			0, 0, 0, 0, 0, 0, 0, 0,
			98, 134, 61, 95, 68, 126, 34, -11,
			-6, 7, 26, 31, 65, 56, 25, -20,
			-14, 13, 6, 21, 23, 12, 17, -23,
			-27, -2, -5, 12, 17, 6, 10, -25,
			-26, -4, -4, -10, 3, 3, 33, -12,
			-35, -1, -20, -23, -15, 24, 38, -22,
			0, 0, 0, 0, 0, 0, 0, 0,
		},

		// knight
		[64]int{
			-167, -89, -34, -49, 61, -97, -15, -107,
			-73, -41, 72, 36, 23, 62, 7, -17,
			-47, 60, 37, 65, 84, 129, 73, 44,
			-9, 17, 19, 53, 37, 69, 18, 22,
			-13, 4, 16, 13, 28, 19, 21, -8,
			-23, -9, 12, 10, 19, 17, 25, -16,
			-29, -53, -12, -3, -1, 18, -14, -19,
			-105, -21, -58, -33, -17, -28, -19, -23,
		},

		// bishop
		[64]int{
			-29, 4, -82, -37, -25, -42, 7, -8,
			-26, 16, -18, -13, 30, 59, 18, -47,
			-16, 37, 43, 40, 35, 50, 37, -2,
			-4, 5, 19, 50, 37, 37, 7, -2,
			-6, 13, 13, 26, 34, 12, 10, 4,
			0, 15, 15, 15, 14, 27, 18, 10,
			4, 15, 16, 0, 7, 21, 33, 1,
			-33, -3, -14, -21, -13, -12, -39, -21,
		},

		// rook
		[64]int{
			32, 42, 32, 51, 63, 9, 31, 43,
			27, 32, 58, 62, 80, 67, 26, 44,
			-5, 19, 26, 36, 17, 45, 61, 16,
			-24, -11, 7, 26, 24, 35, -8, -20,
			-36, -26, -12, -1, 9, -7, 6, -23,
			-45, -25, -16, -17, 3, 0, -5, -33,
			-44, -16, -20, -9, -1, 11, -6, -71,
			-19, -13, 1, 17, 16, 7, -37, -26,
		},

		// queen
		[64]int{
			-28, 0, 29, 12, 59, 44, 43, 45,
			-24, -39, -5, 1, -16, 57, 28, 54,
			-13, -17, 7, 8, 29, 56, 47, 57,
			-27, -27, -16, -16, -1, 17, -2, 1,
			-9, -26, -9, -10, -2, -4, 3, -3,
			-14, 2, -11, -2, -5, 2, 14, 5,
			-35, -8, 11, 2, 8, 15, -3, 1,
			-1, -18, -9, 10, -15, -25, -31, -50,
		},

		// king
		[64]int{
			-65, 23, 16, -15, -56, -34, 2, 13,
			29, -1, -20, -7, -8, -4, -38, -29,
			-9, 24, 2, -16, -20, 6, 22, -22,
			-17, -20, -12, -27, -30, -25, -14, -36,
			-49, -1, -27, -39, -46, -44, -33, -51,
			-14, -14, -22, -46, -44, -30, -15, -27,
			1, 7, -8, -64, -43, -16, 9, 8,
			-15, 36, 12, -54, 8, -28, 24, 14,
		},
	},

	// Endgame positional piece scores //
	{
		//pawn
		[64]int{
			0, 0, 0, 0, 0, 0, 0, 0,
			178, 173, 158, 134, 147, 132, 165, 187,
			94, 100, 85, 67, 56, 53, 82, 84,
			32, 24, 13, 5, -2, 4, 17, 17,
			13, 9, -3, -7, -7, -8, 3, -1,
			4, 7, -6, 1, 0, -5, -1, -8,
			13, 8, 8, 10, 13, 0, 2, -7,
			0, 0, 0, 0, 0, 0, 0, 0,
		},

		// knight
		[64]int{
			-58, -38, -13, -28, -31, -27, -63, -99,
			-25, -8, -25, -2, -9, -25, -24, -52,
			-24, -20, 10, 9, -1, -9, -19, -41,
			-17, 3, 22, 22, 22, 11, 8, -18,
			-18, -6, 16, 25, 16, 17, 4, -18,
			-23, -3, -1, 15, 10, -3, -20, -22,
			-42, -20, -10, -5, -2, -20, -23, -44,
			-29, -51, -23, -15, -22, -18, -50, -64,
		},

		// bishop
		[64]int{
			-14, -21, -11, -8, -7, -9, -17, -24,
			-8, -4, 7, -12, -3, -13, -4, -14,
			2, -8, 0, -1, -2, 6, 0, 4,
			-3, 9, 12, 9, 14, 10, 3, 2,
			-6, 3, 13, 19, 7, 10, -3, -9,
			-12, -3, 8, 10, 13, 3, -7, -15,
			-14, -18, -7, -1, 4, -9, -15, -27,
			-23, -9, -23, -5, -9, -16, -5, -17,
		},

		// rook
		[64]int{
			13, 10, 18, 15, 12, 12, 8, 5,
			11, 13, 13, 11, -3, 3, 8, 3,
			7, 7, 7, 5, 4, -3, -5, -3,
			4, 3, 13, 1, 2, 1, -1, 2,
			3, 5, 8, 4, -5, -6, -8, -11,
			-4, 0, -5, -1, -7, -12, -8, -16,
			-6, -6, 0, 2, -9, -9, -11, -3,
			-9, 2, 3, -1, -5, -13, 4, -20,
		},

		// queen
		[64]int{
			-9, 22, 22, 27, 27, 19, 10, 20,
			-17, 20, 32, 41, 58, 25, 30, 0,
			-20, 6, 9, 49, 47, 35, 19, 9,
			3, 22, 24, 45, 57, 40, 57, 36,
			-18, 28, 19, 47, 31, 34, 39, 23,
			-16, -27, 15, 6, 9, 17, 10, 5,
			-22, -23, -30, -16, -16, -23, -36, -32,
			-33, -28, -22, -43, -5, -32, -20, -41,
		},

		// king
		[64]int{
			-74, -35, -18, -18, -11, 15, 4, -17,
			-12, 17, 14, 17, 17, 38, 23, 11,
			10, 17, 23, 15, 20, 45, 44, 13,
			-8, 22, 24, 27, 26, 33, 26, 3,
			-18, -4, 21, 24, 27, 23, 9, -11,
			-19, -3, 11, 21, 23, 16, 7, -9,
			-27, -11, 4, 13, 14, 4, -5, -17,
			-53, -34, -21, -11, -28, -14, -24, -43,
		},
	},
}

// mirror positional score tables for opposite side
var mirrorScore = [...]square{
	a1, b1, c1, d1, e1, f1, g1, h1,
	a2, b2, c2, d2, e2, f2, g2, h2,
	a3, b3, c3, d3, e3, f3, g3, h3,
	a4, b4, c4, d4, e4, f4, g4, h4,
	a5, b5, c5, d5, e5, f5, g5, h5,
	a6, b6, c6, d6, e6, f6, g6, h6,
	a7, b7, c7, d7, e7, f7, g7, h7,
	a8, b8, c8, d8, e8, f8, g8, h8,
}

// get game phase score
func (b *Board) getGamePhaseScore() int {
	// white & black game phase scores
	whitePieceScores, blackPieceScores := 0, 0

	// loop over white pieces
	for piece := N; piece <= Q; piece++ {
		whitePieceScores += b.Bitboards[piece].countBits() * materialScore[opening][piece]
	}

	// loop over white pieces
	for piece := n; piece <= q; piece++ {
		blackPieceScores += b.Bitboards[piece].countBits() * -materialScore[opening][piece]
	}

	// return game phase score
	return whitePieceScores + blackPieceScores
}

// convert BBC piece code to Stockfish piece codes
var nnue_pieces = [12]int{6, 5, 4, 3, 2, 1, 12, 11, 10, 9, 8, 7}

// convert BBC square indices to Stockfish indices
var nnue_squares = [64]square{
	a1, b1, c1, d1, e1, f1, g1, h1,
	a2, b2, c2, d2, e2, f2, g2, h2,
	a3, b3, c3, d3, e3, f3, g3, h3,
	a4, b4, c4, d4, e4, f4, g4, h4,
	a5, b5, c5, d5, e5, f5, g5, h5,
	a6, b6, c6, d6, e6, f6, g6, h6,
	a7, b7, c7, d7, e7, f7, g7, h7,
	a8, b8, c8, d8, e8, f8, g8, h8,
}

func (board *Board) evaluate() int {
	// get game phase score
	gamePhaseScore := board.getGamePhaseScore()

	// game phase (opening, middle game, endgame)
	gamePhase := -1

	// pick up game phase based on game phase score
	if gamePhaseScore > openingPhaseScore {
		gamePhase = opening
	} else if gamePhaseScore < endgamePhaseScore {
		gamePhase = endgame
	} else {
		gamePhase = middlegame
	}

	// static evaluation score
	score, scoreOpening, scoreEndgame := 0, 0, 0

	// current pieces bitboard copy
	var bitboard Bitboard

	// // penalties
	// doublePawns := 0

	piece := P
	for bitboard = board.Bitboards[piece]; bitboard > 0; {
		// init square
		square := square(bitboard.GetLs1bIndex())

		// get opening/endgame material score
		scoreOpening += materialScore[opening][piece]
		scoreEndgame += materialScore[endgame][piece]

		// get opening/endgame positional score
		scoreOpening += positionalScore[opening][PAWN][square]
		scoreEndgame += positionalScore[endgame][PAWN][square]

		bitboard.popBit(square)
	}

	piece = N
	for bitboard = board.Bitboards[piece]; bitboard > 0; {
		// init square
		square := square(bitboard.GetLs1bIndex())

		// get opening/endgame material score
		scoreOpening += materialScore[opening][piece]
		scoreEndgame += materialScore[endgame][piece]

		// get opening/endgame positional score
		scoreOpening += positionalScore[opening][KNIGHT][square]
		scoreEndgame += positionalScore[endgame][KNIGHT][square]

		bitboard.popBit(square)
	}

	piece = B
	for bitboard = board.Bitboards[piece]; bitboard > 0; {
		// init square
		square := square(bitboard.GetLs1bIndex())

		// get opening/endgame material score
		scoreOpening += materialScore[opening][piece]
		scoreEndgame += materialScore[endgame][piece]

		// get opening/endgame positional score
		scoreOpening += positionalScore[opening][BISHOP][square]
		scoreEndgame += positionalScore[endgame][BISHOP][square]

		bitboard.popBit(square)
	}

	piece = R
	for bitboard = board.Bitboards[piece]; bitboard > 0; {
		// init square
		square := square(bitboard.GetLs1bIndex())

		// get opening/endgame material score
		scoreOpening += materialScore[opening][piece]
		scoreEndgame += materialScore[endgame][piece]

		// get opening/endgame positional score
		scoreOpening += positionalScore[opening][ROOK][square]
		scoreEndgame += positionalScore[endgame][ROOK][square]

		bitboard.popBit(square)
	}

	piece = Q
	for bitboard = board.Bitboards[piece]; bitboard > 0; {
		// init square
		square := square(bitboard.GetLs1bIndex())

		// get opening/endgame material score
		scoreOpening += materialScore[opening][piece]
		scoreEndgame += materialScore[endgame][piece]

		// get opening/endgame positional score
		scoreOpening += positionalScore[opening][QUEEN][square]
		scoreEndgame += positionalScore[endgame][QUEEN][square]

		bitboard.popBit(square)
	}

	piece = K
	for bitboard = board.Bitboards[piece]; bitboard > 0; {
		// init square
		square := square(bitboard.GetLs1bIndex())

		// get opening/endgame material score
		scoreOpening += materialScore[opening][piece]
		scoreEndgame += materialScore[endgame][piece]

		// get opening/endgame positional score
		scoreOpening += positionalScore[opening][KING][square]
		scoreEndgame += positionalScore[endgame][KING][square]

		bitboard.popBit(square)
	}

	piece = p
	for bitboard = board.Bitboards[piece]; bitboard > 0; {
		// init square
		square := square(bitboard.GetLs1bIndex())

		// get opening/endgame material score
		scoreOpening += materialScore[opening][piece]
		scoreEndgame += materialScore[endgame][piece]

		// get opening/endgame positional score
		scoreOpening -= positionalScore[opening][PAWN][mirrorScore[square]]
		scoreEndgame -= positionalScore[endgame][PAWN][mirrorScore[square]]

		bitboard.popBit(square)
	}

	piece = n
	for bitboard = board.Bitboards[piece]; bitboard > 0; {
		// init square
		square := square(bitboard.GetLs1bIndex())

		// get opening/endgame material score
		scoreOpening += materialScore[opening][piece]
		scoreEndgame += materialScore[endgame][piece]

		// get opening/endgame positional score
		scoreOpening -= positionalScore[opening][KNIGHT][mirrorScore[square]]
		scoreEndgame -= positionalScore[endgame][KNIGHT][mirrorScore[square]]

		bitboard.popBit(square)
	}

	piece = b
	for bitboard = board.Bitboards[piece]; bitboard > 0; {
		// init square
		square := square(bitboard.GetLs1bIndex())

		// get opening/endgame material score
		scoreOpening += materialScore[opening][piece]
		scoreEndgame += materialScore[endgame][piece]

		// get opening/endgame positional score
		scoreOpening -= positionalScore[opening][BISHOP][mirrorScore[square]]
		scoreEndgame -= positionalScore[endgame][BISHOP][mirrorScore[square]]

		bitboard.popBit(square)
	}

	piece = r
	for bitboard = board.Bitboards[piece]; bitboard > 0; {
		// init square
		square := square(bitboard.GetLs1bIndex())

		// get opening/endgame material score
		scoreOpening += materialScore[opening][piece]
		scoreEndgame += materialScore[endgame][piece]

		// get opening/endgame positional score
		scoreOpening -= positionalScore[opening][ROOK][mirrorScore[square]]
		scoreEndgame -= positionalScore[endgame][ROOK][mirrorScore[square]]

		bitboard.popBit(square)
	}

	piece = q
	for bitboard = board.Bitboards[piece]; bitboard > 0; {
		// init square
		square := square(bitboard.GetLs1bIndex())

		// get opening/endgame material score
		scoreOpening += materialScore[opening][piece]
		scoreEndgame += materialScore[endgame][piece]

		// get opening/endgame positional score
		scoreOpening -= positionalScore[opening][QUEEN][mirrorScore[square]]
		scoreEndgame -= positionalScore[endgame][QUEEN][mirrorScore[square]]

		bitboard.popBit(square)
	}

	piece = k
	for bitboard = board.Bitboards[piece]; bitboard > 0; {
		// init square
		square := square(bitboard.GetLs1bIndex())

		// get opening/endgame material score
		scoreOpening += materialScore[opening][piece]
		scoreEndgame += materialScore[endgame][piece]

		// get opening/endgame positional score
		scoreOpening -= positionalScore[opening][KING][mirrorScore[square]]
		scoreEndgame -= positionalScore[endgame][KING][mirrorScore[square]]

		bitboard.popBit(square)
	}

	if gamePhase == middlegame {
		score = (scoreOpening*gamePhaseScore + scoreEndgame*(openingPhaseScore-gamePhaseScore)) / openingPhaseScore
	} else if gamePhase == opening {
		score = scoreOpening
	} else if gamePhase == endgame {
		score = scoreEndgame
	}

	// return final evaluation based on side
	if board.Side == white {
		return score
	}

	return -score
}
