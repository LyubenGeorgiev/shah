package chess

// init all variables (Called automatically!)
func init() {
	// init leaper pieces attacks
	init_leapers_attacks()

	// init slider pieces attacks
	init_sliders_attacks(bishop)
	init_sliders_attacks(rook)

	// init random keys for hashing purposes
	init_random_keys()
}
