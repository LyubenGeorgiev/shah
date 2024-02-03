package chess

// random piece keys [piece][square]
var piece_keys [12][64]Bitboard

// random enpassant keys [square]
var enpassant_keys [64]Bitboard

// random castling keys
var castle_keys [16]Bitboard

// random side key
var side_key Bitboard

func get_random_U32_number(random_state *uint32) uint32 {
	// get current state
	number := *random_state

	// XOR shift algorithm
	number ^= number << 13
	number ^= number >> 17
	number ^= number << 5

	// update random number state
	*random_state = number

	// return random number
	return number
}

// generate 64-bit pseudo legal numbers
func get_random_U64_number(random_state *uint32) Bitboard {
	// define 4 random numbers
	var n1, n2, n3, n4 uint64

	// init random numbers slicing 16 bits from MS1B side
	n1 = uint64(get_random_U32_number(random_state)) & 0xFFFF
	n2 = uint64(get_random_U32_number(random_state)) & 0xFFFF
	n3 = uint64(get_random_U32_number(random_state)) & 0xFFFF
	n4 = uint64(get_random_U32_number(random_state)) & 0xFFFF

	// return random number
	return Bitboard(n1 | (n2 << 16) | (n3 << 32) | (n4 << 48))
}

// init random hash keys
func init_random_keys() {
	// update pseudo random number state
	var random_state uint32 = 1804289383

	// loop over piece codes
	for piece := P; piece <= k; piece++ {
		// loop over board squares
		for square := 0; square < 64; square++ {
			// init random piece keys
			piece_keys[piece][square] = get_random_U64_number(&random_state)
		}
	}

	// loop over board squares
	for square := 0; square < 64; square++ {
		// init random enpassant keys
		enpassant_keys[square] = get_random_U64_number(&random_state)
	}

	// loop over castling keys
	for index := 0; index < 16; index++ {
		// init castling keys
		castle_keys[index] = get_random_U64_number(&random_state)
	}

	// init random side key
	side_key = get_random_U64_number(&random_state)
}

// generate "almost" unique position ID aka hash key from scratch
func (b *Board) generate_hash_key() Bitboard {
	// final hash key
	final_key := Bitboard(0)

	// loop over piece bitboards
	for piece := P; piece <= k; piece++ {
		// init piece bitboard copy
		bitboard := b.Bitboards[piece]

		// loop over the pieces within a bitboard
		for bitboard > 0 {
			// init square occupied by the piece
			square := square(bitboard.GetLs1bIndex())

			// hash piece
			final_key ^= piece_keys[piece][square]

			// pop LS1B
			bitboard.popBit(square)
		}
	}

	// if enpassant square is on board
	if b.Enpassant != no_sq {
		// hash enpassant
		final_key ^= enpassant_keys[b.Enpassant]
	}

	// hash castling rights
	final_key ^= castle_keys[b.Castle]

	// hash the side only if black is to move
	if b.Side == black {
		final_key ^= side_key
	}

	// return generated hash key
	return final_key
}
