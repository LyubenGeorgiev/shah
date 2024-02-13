package chess

var stringToSquare = map[string]square{
	"a1": 56, "a2": 48, "a3": 40, "a4": 32, "a5": 24, "a6": 16, "a7": 8, "a8": 0,
	"b1": 57, "b2": 49, "b3": 41, "b4": 33, "b5": 25, "b6": 17, "b7": 9, "b8": 1,
	"c1": 58, "c2": 50, "c3": 42, "c4": 34, "c5": 26, "c6": 18, "c7": 10, "c8": 2,
	"d1": 59, "d2": 51, "d3": 43, "d4": 35, "d5": 27, "d6": 19, "d7": 11, "d8": 3,
	"e1": 60, "e2": 52, "e3": 44, "e4": 36, "e5": 28, "e6": 20, "e7": 12, "e8": 4,
	"f1": 61, "f2": 53, "f3": 45, "f4": 37, "f5": 29, "f6": 21, "f7": 13, "f8": 5,
	"g1": 62, "g2": 54, "g3": 46, "g4": 38, "g5": 30, "g6": 22, "g7": 14, "g8": 6,
	"h1": 63, "h2": 55, "h3": 47, "h4": 39, "h5": 31, "h6": 23, "h7": 15, "h8": 7,
}

var StringToSquare = map[string]int{
	"a1": 56, "a2": 48, "a3": 40, "a4": 32, "a5": 24, "a6": 16, "a7": 8, "a8": 0,
	"b1": 57, "b2": 49, "b3": 41, "b4": 33, "b5": 25, "b6": 17, "b7": 9, "b8": 1,
	"c1": 58, "c2": 50, "c3": 42, "c4": 34, "c5": 26, "c6": 18, "c7": 10, "c8": 2,
	"d1": 59, "d2": 51, "d3": 43, "d4": 35, "d5": 27, "d6": 19, "d7": 11, "d8": 3,
	"e1": 60, "e2": 52, "e3": 44, "e4": 36, "e5": 28, "e6": 20, "e7": 12, "e8": 4,
	"f1": 61, "f2": 53, "f3": 45, "f4": 37, "f5": 29, "f6": 21, "f7": 13, "f8": 5,
	"g1": 62, "g2": 54, "g3": 46, "g4": 38, "g5": 30, "g6": 22, "g7": 14, "g8": 6,
	"h1": 63, "h2": 55, "h3": 47, "h4": 39, "h5": 31, "h6": 23, "h7": 15, "h8": 7,
}

var squareToString = [...]string{
	"a8", "b8", "c8", "d8", "e8", "f8", "g8", "h8",
	"a7", "b7", "c7", "d7", "e7", "f7", "g7", "h7",
	"a6", "b6", "c6", "d6", "e6", "f6", "g6", "h6",
	"a5", "b5", "c5", "d5", "e5", "f5", "g5", "h5",
	"a4", "b4", "c4", "d4", "e4", "f4", "g4", "h4",
	"a3", "b3", "c3", "d3", "e3", "f3", "g3", "h3",
	"a2", "b2", "c2", "d2", "e2", "f2", "g2", "h2",
	"a1", "b1", "c1", "d1", "e1", "f1", "g1", "h1",
}
