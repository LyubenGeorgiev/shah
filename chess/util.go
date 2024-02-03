package chess

import "math/bits"

const (
	not_a_file      Bitboard = 18374403900871474942
	not_h_file      Bitboard = 9187201950435737471
	not_hg_file     Bitboard = 4557430888798830399
	not_ab_file     Bitboard = 18229723555195321596
	empty_board     string   = "8/8/8/8/8/8/8/8 w - - "
	start_position  string   = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1 "
	tricky_position string   = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1 "
	killer_position string   = "rnbqkb1r/pp1p1pPp/8/2p1pP2/1P1P4/3P3P/P1P1P3/RNBQKBNR w KQkq e6 0 1"
	cmk_position    string   = "r2q1rk1/ppp2ppp/2n1bn2/2b1p3/3pP3/3P1NPP/PPP1NPB1/R1BQ1RK1 b - - 0 9 "
)

// Convert ASCII character pieces to encoded constants
var charToPiece = [...]piece{'P': P, 'N': N, 'B': B, 'R': R, 'Q': Q, 'K': K, 'p': p, 'n': n, 'b': b, 'r': r, 'q': q, 'k': k}

// Convert encoded constants to ASCII character pieces
var pieceToChar = [...]byte{P: 'P', N: 'N', B: 'B', R: 'R', Q: 'Q', K: 'K', p: 'p', n: 'n', b: 'b', r: 'r', q: 'q', k: 'k', no_piece: ' '}

// Used to determine what the rights should be if a piece from here has moved
var castlingRights = [...]castle{
	7, 15, 15, 15, 3, 15, 15, 11,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	13, 15, 15, 15, 12, 15, 15, 14,
}

// bishop relevant occupancy bit count for every square on board
var bishop_relevant_bits = [64]int{
	6, 5, 5, 5, 5, 5, 5, 6,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	6, 5, 5, 5, 5, 5, 5, 6,
}

// rook relevant occupancy bit count for every square on board
var rook_relevant_bits = [64]int{
	12, 11, 11, 11, 11, 11, 11, 12,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	12, 11, 11, 11, 11, 11, 11, 12,
}

// bishop magic numbers
var bishop_magic_numbers = [64]Bitboard{
	0x40040844404084,
	0x2004208a004208,
	0x10190041080202,
	0x108060845042010,
	0x581104180800210,
	0x2112080446200010,
	0x1080820820060210,
	0x3c0808410220200,
	0x4050404440404,
	0x21001420088,
	0x24d0080801082102,
	0x1020a0a020400,
	0x40308200402,
	0x4011002100800,
	0x401484104104005,
	0x801010402020200,
	0x400210c3880100,
	0x404022024108200,
	0x810018200204102,
	0x4002801a02003,
	0x85040820080400,
	0x810102c808880400,
	0xe900410884800,
	0x8002020480840102,
	0x220200865090201,
	0x2010100a02021202,
	0x152048408022401,
	0x20080002081110,
	0x4001001021004000,
	0x800040400a011002,
	0xe4004081011002,
	0x1c004001012080,
	0x8004200962a00220,
	0x8422100208500202,
	0x2000402200300c08,
	0x8646020080080080,
	0x80020a0200100808,
	0x2010004880111000,
	0x623000a080011400,
	0x42008c0340209202,
	0x209188240001000,
	0x400408a884001800,
	0x110400a6080400,
	0x1840060a44020800,
	0x90080104000041,
	0x201011000808101,
	0x1a2208080504f080,
	0x8012020600211212,
	0x500861011240000,
	0x180806108200800,
	0x4000020e01040044,
	0x300000261044000a,
	0x802241102020002,
	0x20906061210001,
	0x5a84841004010310,
	0x4010801011c04,
	0xa010109502200,
	0x4a02012000,
	0x500201010098b028,
	0x8040002811040900,
	0x28000010020204,
	0x6000020202d0240,
	0x8918844842082200,
	0x4010011029020020,
}

// rook magic numbers
var rook_magic_numbers = [64]Bitboard{
	0x8a80104000800020,
	0x140002000100040,
	0x2801880a0017001,
	0x100081001000420,
	0x200020010080420,
	0x3001c0002010008,
	0x8480008002000100,
	0x2080088004402900,
	0x800098204000,
	0x2024401000200040,
	0x100802000801000,
	0x120800800801000,
	0x208808088000400,
	0x2802200800400,
	0x2200800100020080,
	0x801000060821100,
	0x80044006422000,
	0x100808020004000,
	0x12108a0010204200,
	0x140848010000802,
	0x481828014002800,
	0x8094004002004100,
	0x4010040010010802,
	0x20008806104,
	0x100400080208000,
	0x2040002120081000,
	0x21200680100081,
	0x20100080080080,
	0x2000a00200410,
	0x20080800400,
	0x80088400100102,
	0x80004600042881,
	0x4040008040800020,
	0x440003000200801,
	0x4200011004500,
	0x188020010100100,
	0x14800401802800,
	0x2080040080800200,
	0x124080204001001,
	0x200046502000484,
	0x480400080088020,
	0x1000422010034000,
	0x30200100110040,
	0x100021010009,
	0x2002080100110004,
	0x202008004008002,
	0x20020004010100,
	0x2048440040820001,
	0x101002200408200,
	0x40802000401080,
	0x4008142004410100,
	0x2060820c0120200,
	0x1001004080100,
	0x20c020080040080,
	0x2935610830022400,
	0x44440041009200,
	0x280001040802101,
	0x2100190040002085,
	0x80c0084100102001,
	0x4024081001000421,
	0x20030a0244872,
	0x12001008414402,
	0x2006104900a0804,
	0x1004081002402,
}

// pawn attacks table [side][square]
var pawn_attacks [2][64]Bitboard

// knight attacks table [square]
var knight_attacks [64]Bitboard

// king attacks table [square]
var king_attacks [64]Bitboard

// bishop attack masks
var bishop_masks [64]Bitboard

// rook attack masks
var rook_masks [64]Bitboard

// bishop attacks table [square][occupancies]
var bishop_attacks [64][512]Bitboard

// rook attacks rable [square][occupancies]
var rook_attacks [64][4096]Bitboard

// generate pawn attacks
func mask_pawn_attacks(side Side, square square) Bitboard {
	// result attacks bitboard
	var attacks Bitboard = 0

	// piece bitboard
	var bitboard Bitboard = 0

	// set piece on board
	bitboard.setBit(square)

	// white pawns
	if side == white {
		if ((bitboard >> 7) & not_a_file) > 0 {
			attacks |= (bitboard >> 7)
		}
		if ((bitboard >> 9) & not_h_file) > 0 {
			attacks |= (bitboard >> 9)
		}
	} else {
		if ((bitboard << 7) & not_h_file) > 0 {
			attacks |= (bitboard << 7)
		}
		if ((bitboard << 9) & not_a_file) > 0 {
			attacks |= (bitboard << 9)
		}
	}

	return attacks
}

// generate knight attacks
func mask_knight_attacks(square square) Bitboard {
	// result attacks bitboard
	var attacks Bitboard = 0

	// piece bitboard
	var bitboard Bitboard = 0

	// set piece on board
	bitboard.setBit(square)

	// generate knight attacks
	if ((bitboard >> 17) & not_h_file) > 0 {
		attacks |= (bitboard >> 17)
	}
	if ((bitboard >> 15) & not_a_file) > 0 {
		attacks |= (bitboard >> 15)
	}
	if ((bitboard >> 10) & not_hg_file) > 0 {
		attacks |= (bitboard >> 10)
	}
	if ((bitboard >> 6) & not_ab_file) > 0 {
		attacks |= (bitboard >> 6)
	}
	if ((bitboard << 17) & not_a_file) > 0 {
		attacks |= (bitboard << 17)
	}
	if ((bitboard << 15) & not_h_file) > 0 {
		attacks |= (bitboard << 15)
	}
	if ((bitboard << 10) & not_ab_file) > 0 {
		attacks |= (bitboard << 10)
	}
	if ((bitboard << 6) & not_hg_file) > 0 {
		attacks |= (bitboard << 6)
	}

	return attacks
}

// generate king attacks
func mask_king_attacks(square square) Bitboard {
	// result attacks bitboard
	var attacks Bitboard = 0

	// piece bitboard
	var bitboard Bitboard = 0

	// set piece on board
	bitboard.setBit(square)

	// generate king attacks
	if (bitboard >> 8) > 0 {
		attacks |= (bitboard >> 8)
	}
	if ((bitboard >> 9) & not_h_file) > 0 {
		attacks |= (bitboard >> 9)
	}
	if ((bitboard >> 7) & not_a_file) > 0 {
		attacks |= (bitboard >> 7)
	}
	if ((bitboard >> 1) & not_h_file) > 0 {
		attacks |= (bitboard >> 1)
	}
	if (bitboard << 8) > 0 {
		attacks |= (bitboard << 8)
	}
	if ((bitboard << 9) & not_a_file) > 0 {
		attacks |= (bitboard << 9)
	}
	if ((bitboard << 7) & not_h_file) > 0 {
		attacks |= (bitboard << 7)
	}
	if ((bitboard << 1) & not_a_file) > 0 {
		attacks |= (bitboard << 1)
	}

	return attacks
}

// mask bishop attacks
func mask_bishop_attacks(square square) Bitboard {
	// result attacks bitboard
	var attacks Bitboard = 0

	// init target rank & files
	tr := square / 8
	tf := square % 8

	// mask relevant bishop occupancy bits
	for r, f := tr+1, tf+1; r <= 6 && f <= 6; r, f = r+1, f+1 {
		attacks |= (Bitboard(1) << (r*8 + f))
	}
	for r, f := tr-1, tf+1; r >= 1 && f <= 6; r, f = r-1, f+1 {
		attacks |= (Bitboard(1) << (r*8 + f))
	}
	for r, f := tr+1, tf-1; r <= 6 && f >= 1; r, f = r+1, f-1 {
		attacks |= (Bitboard(1) << (r*8 + f))
	}
	for r, f := tr-1, tf-1; r >= 1 && f >= 1; r, f = r-1, f-1 {
		attacks |= (Bitboard(1) << (r*8 + f))
	}

	return attacks
}

// mask rook attacks
func mask_rook_attacks(square square) Bitboard {
	// result attacks bitboard
	var attacks Bitboard = 0

	// init target rank & files
	tr := square / 8
	tf := square % 8

	// mask relevant rook occupancy bits
	for r := tr + 1; r <= 6; r++ {
		attacks |= (Bitboard(1) << (r*8 + tf))
	}
	for r := tr - 1; r >= 1; r-- {
		attacks |= (Bitboard(1) << (r*8 + tf))
	}
	for f := tf + 1; f <= 6; f++ {
		attacks |= (Bitboard(1) << (tr*8 + f))
	}
	for f := tf - 1; f >= 1; f-- {
		attacks |= (Bitboard(1) << (tr*8 + f))
	}

	return attacks
}

// generate bishop attacks on the fly
func bishop_attacks_on_the_fly(square square, block Bitboard) Bitboard {
	// result attacks bitboard
	var attacks Bitboard = 0

	// init target rank & files
	tr := square / 8
	tf := square % 8

	// generate bishop atacks
	for r, f := tr+1, tf+1; r <= 7 && f <= 7; r, f = r+1, f+1 {
		attacks |= (Bitboard(1) << (r*8 + f))
		if ((Bitboard(1) << (r*8 + f)) & block) > 0 {
			break
		}
	}
	for r, f := tr-1, tf+1; r >= 0 && f <= 7; r, f = r-1, f+1 {
		attacks |= (Bitboard(1) << (r*8 + f))
		if ((Bitboard(1) << (r*8 + f)) & block) > 0 {
			break
		}
	}
	for r, f := tr+1, tf-1; r <= 7 && f >= 0; r, f = r+1, f-1 {
		attacks |= (Bitboard(1) << (r*8 + f))
		if ((Bitboard(1) << (r*8 + f)) & block) > 0 {
			break
		}
	}
	for r, f := tr-1, tf-1; r >= 0 && f >= 0; r, f = r-1, f-1 {
		attacks |= (Bitboard(1) << (r*8 + f))
		if ((Bitboard(1) << (r*8 + f)) & block) > 0 {
			break
		}
	}

	return attacks
}

// generate rook attacks on the fly
func rook_attacks_on_the_fly(square square, block Bitboard) Bitboard {
	// result attacks bitboard
	var attacks Bitboard = 0

	// init target rank & files
	tr := square / 8
	tf := square % 8

	// generate rook attacks
	for r := tr + 1; r <= 7; r++ {
		attacks |= (Bitboard(1) << (r*8 + tf))
		if ((Bitboard(1) << (r*8 + tf)) & block) > 0 {
			break
		}
	}
	for r := tr - 1; r >= 0; r-- {
		attacks |= (Bitboard(1) << (r*8 + tf))
		if ((Bitboard(1) << (r*8 + tf)) & block) > 0 {
			break
		}
	}
	for f := tf + 1; f <= 7; f++ {
		attacks |= (Bitboard(1) << (tr*8 + f))
		if ((Bitboard(1) << (tr*8 + f)) & block) > 0 {
			break
		}
	}
	for f := tf - 1; f >= 0; f-- {
		attacks |= (Bitboard(1) << (tr*8 + f))
		if ((Bitboard(1) << (tr*8 + f)) & block) > 0 {
			break
		}
	}

	return attacks
}

// init leaper pieces attacks
func init_leapers_attacks() {
	// loop over 64 board squares
	for square := square(0); square < 64; square++ {
		// init pawn attacks
		pawn_attacks[white][square] = mask_pawn_attacks(white, square)
		pawn_attacks[black][square] = mask_pawn_attacks(black, square)

		// init knight attacks
		knight_attacks[square] = mask_knight_attacks(square)

		// init king attacks
		king_attacks[square] = mask_king_attacks(square)
	}
}

// set occupancies
func set_occupancy(index int, bits_in_mask int, attack_mask Bitboard) Bitboard {
	// occupancy map
	var occupancy Bitboard = 0

	// loop over the range of bits within attack mask
	for count := 0; count < bits_in_mask; count++ {
		// get LS1B index of attacks mask
		square := square(attack_mask.GetLs1bIndex())

		// pop LS1B in attack map
		attack_mask.popBit(square)

		// make sure occupancy is on board
		if (index & (1 << count)) > 0 {
			occupancy |= (Bitboard(1) << square)
		}
	}

	// return occupancy map
	return occupancy
}

// init slider piece's attack tables
func init_sliders_attacks(slider slider) {
	// loop over 64 board squares
	for square := square(0); square < 64; square++ {
		// init bishop & rook masks
		bishop_masks[square] = mask_bishop_attacks(square)
		rook_masks[square] = mask_rook_attacks(square)

		// init current mask
		attack_mask := bishop_masks[square]
		if slider == rook {
			attack_mask = rook_masks[square]
		}

		// init relevant occupancy bit count
		relevant_bits_count := attack_mask.countBits()

		// init occupancy indicies
		occupancy_indicies := (1 << relevant_bits_count)

		// loop over occupancy indicies
		for index := 0; index < occupancy_indicies; index++ {
			if slider == bishop {
				// init current occupancy variation
				occupancy := set_occupancy(index, relevant_bits_count, attack_mask)

				// init magic index
				magic_index := (occupancy * bishop_magic_numbers[square]) >> (64 - bishop_relevant_bits[square])

				// init bishop attacks
				bishop_attacks[square][magic_index] = bishop_attacks_on_the_fly(square, occupancy)
			} else {
				// init current occupancy variation
				occupancy := set_occupancy(index, relevant_bits_count, attack_mask)

				// init magic index
				magic_index := (occupancy * rook_magic_numbers[square]) >> (64 - rook_relevant_bits[square])

				// init bishop attacks
				rook_attacks[square][magic_index] = rook_attacks_on_the_fly(square, occupancy)

			}
		}
	}
}

func (b *Bitboard) setBit(s square) {
	*b |= (Bitboard(1) << s)
}

func (b Bitboard) getBit(s square) bool {
	return (b & (Bitboard(1) << s)) > 0
}

func (b *Bitboard) popBit(s square) {
	*b &= ^(Bitboard(1) << s)
}

// count bits within a bitboard (Brian Kernighan's way)
func (b Bitboard) countBits() int {
	return bits.OnesCount64(uint64(b))
}

// get least significant 1st bit index
func (b Bitboard) GetLs1bIndex() int {
	// make sure bitboard is not 0
	if uint64(b) > 0 {
		return bits.TrailingZeros64(uint64(b))
	} else {
		return -1
	}
}

func (b *Bitboard) PopLs1b() {
	b.popBit(square(b.GetLs1bIndex()))
}
