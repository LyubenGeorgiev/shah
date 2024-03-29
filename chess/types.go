package chess

type Side int

const (
	white Side = iota
	black
	both
)

func (s Side) opposite() Side {
	switch s {
	case white:
		return black
	case black:
		return white
	default:
		return both
	}
}

type castle int

// encode castling rights
const (
	wk castle = 1 << iota
	wq
	bk
	bq
)

type slider int

// encode sliding pieces
const (
	rook slider = iota
	bishop
)

type piece int

// encode pieces
const (
	P piece = iota
	N
	B
	R
	Q
	K
	p
	n
	b
	r
	q
	k
	piecesSize
	no_piece
)

type square int

const (
	a8 square = iota
	b8
	c8
	d8
	e8
	f8
	g8
	h8
	a7
	b7
	c7
	d7
	e7
	f7
	g7
	h7
	a6
	b6
	c6
	d6
	e6
	f6
	g6
	h6
	a5
	b5
	c5
	d5
	e5
	f5
	g5
	h5
	a4
	b4
	c4
	d4
	e4
	f4
	g4
	h4
	a3
	b3
	c3
	d3
	e3
	f3
	g3
	h3
	a2
	b2
	c2
	d2
	e2
	f2
	g2
	h2
	a1
	b1
	c1
	d1
	e1
	f1
	g1
	h1
	no_sq
)

type Bitboard uint64
