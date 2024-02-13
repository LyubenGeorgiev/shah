package chess

import (
	"bytes"
	"fmt"
	"unicode"
)

// parse user/GUI move string input (e.g. "e7e8q")
func (board *Board) parseMove(moveString []byte) Move {
	// create move list instance
	var moves Moves

	// generate moves
	board.generateMoves(&moves)

	// parse source square
	sourceSquare := square((moveString[0] - 'a') + (8-(moveString[1]-'0'))*8)

	// parse target square
	targetSquare := square((moveString[2] - 'a') + (8-(moveString[3]-'0'))*8)

	// loop over the moves within a move list
	for i := 0; i < moves.count; i++ {
		// init move
		move := moves.moves[i]

		// make sure source & target squares are available within the generated move
		if sourceSquare == move.getSource() && targetSquare == move.getTarget() {
			// init promoted piece
			promotedPiece := move.getPromotionPiece()

			// promoted piece is available
			if promotedPiece != no_piece {
				// promoted to queen
				if (promotedPiece == Q || promotedPiece == q) && moveString[4] == 'q' {
					// return legal move
					return move
				} else if (promotedPiece == R || promotedPiece == r) && moveString[4] == 'r' {
					// return legal move
					return move
				} else if (promotedPiece == B || promotedPiece == b) && moveString[4] == 'b' {
					// return legal move
					return move
				} else if (promotedPiece == N || promotedPiece == n) && moveString[4] == 'n' {
					// return legal move
					return move
				}

				// continue the loop on possible wrong promotions (e.g. "e7e8f")
				continue
			}

			// return legal move
			return move
		}
	}

	// return illegal move
	return 0
}

// parse UCI "position" command
func (e *Engine) parsePosition(command []byte) {
	e.RepetitionTable = e.RepetitionTable[:0]

	// parse UCI "startpos" command
	if bytes.HasPrefix(command, STARTPOS) {
		e.Board = *NewBoadFromFen(STARTPOS_FEN)
	} else {
		// if no "fen" command is available within command string
		if !bytes.Contains(command, FEN) {
			// init chess board with start position
			e.Board = *NewBoadFromFen(STARTPOS_FEN)
		} else {
			// init chess board with position from FEN string
			e.Board = *NewBoadFromFen(command[4:])
		}
	}

	// parse moves after position
	index := bytes.Index(command, MOVES)

	// moves available
	if index != -1 {
		// shift pointer to the right where next token begins
		command = command[index+6:]

		// loop over moves within a move string
		for i := 0; i < len(command); i++ {
			// parse next move
			move := e.Board.parseMove(command[i:])

			// if no more moves
			if move == 0 {
				// break out of the loop
				break
			}

			e.RepetitionTable = append(e.RepetitionTable, e.HashKey)

			// make move on the chess board
			if !e.Board.MakeMove(move, false) {
				fmt.Println("Failed at", MoveToString(move))
			}

			// move current character mointer to the end of current move
			for i < len(command) && command[i] != ' ' {
				i++
			}
		}
	}
}

func parseOption(command []byte, opt []byte) int {
	if i := bytes.Index(command, opt); i != -1 {
		var val int
		_, err := fmt.Sscan(string(command[i+len(opt)+1:]), &val)
		if err != nil {
			return -1
		}

		return val
	}

	return -1
}

// parse UCI "go" command
func (e *Engine) parseGo(command []byte) {
	e.TimeController = NewTimeController()

	var inc, time int
	if e.Board.Side == white {
		inc = parseOption(command, WINC)
		time = parseOption(command, WTIME)
	} else {
		inc = parseOption(command, BINC)
		time = parseOption(command, BTIME)
	}

	if inc != -1 {
		e.Inc = inc
	}
	if time != -1 {
		e.Time = time
	}

	movestogo := parseOption(command, MOVESTOGO)
	if movestogo != -1 {
		e.MovesToGo = movestogo
	}

	movetime := parseOption(command, MOVETIME)
	if movetime != -1 {
		e.MoveTime = movetime
	}

	depth := parseOption(command, DEPTH)

	// if move time is not available
	if e.MoveTime != -1 {
		// set time equal to move time
		e.Time = e.MoveTime

		// set moves to go to 1
		e.MovesToGo = 1
	}

	// init start time
	e.StartTime = GetTimeMs()

	// if time control is available
	if e.Time != -1 {
		// flag we're playing with time control
		e.Timeset = 1

		// set up timing
		e.Time /= e.MovesToGo

		duration := int64(e.Time) + int64(e.Inc) - 1000
		if duration < 0 {
			e.Timeset = 0
			depth = 5
		}

		e.StopTime = e.StartTime + duration
	}

	// if depth is not available
	if depth == -1 {
		// set depth to 64 plies (takes ages to complete...)
		depth = 64
	}

	// fmt.Println("Side:", e.Board.Side == white)
	// fmt.Printf("%#v\n", e.TimeController)
	// fmt.Printf("time: %d  start: %d  stop: %d  depth: %d  timeset:%d\n", e.Time, e.StartTime, e.StopTime, depth, e.Timeset)

	// search position
	e.searchPosition(depth)
}

// print move (for UCI purposes)
func MoveToString(move Move) string {
	if move.getPromotionPiece() != no_piece {
		return fmt.Sprintf("%s%s%c", squareToString[move.getSource()], squareToString[move.getTarget()], pieceToChar[move.getPromotionPiece()])
	}

	return fmt.Sprintf("%s%s", squareToString[move.getSource()], squareToString[move.getTarget()])
}

// print move (for UCI purposes)
func pvMoveToString(move Move) string {
	if move.getPromotionPiece() != no_piece {
		return fmt.Sprintf("%s%s%c", squareToString[move.getSource()], squareToString[move.getTarget()], byte(unicode.ToLower(rune(pieceToChar[move.getPromotionPiece()]))))
	}

	return fmt.Sprintf("%s%s", squareToString[move.getSource()], squareToString[move.getTarget()])
}

var (
	UCI          = []byte("uci")
	ISREADY      = []byte("isready")
	UCINEWGAME   = []byte("ucinewgame")
	POSITION     = []byte("position")
	GO           = []byte("go")
	STARTPOS     = []byte("startpos")
	FEN          = []byte("fen")
	STARTPOS_FEN = []byte(start_position)
	MOVES        = []byte("moves")
	DEPTH        = []byte("depth")
	WINC         = []byte("winc")
	BINC         = []byte("binc")
	WTIME        = []byte("wtime")
	BTIME        = []byte("btime")
	MOVETIME     = []byte("movetime")
	MOVESTOGO    = []byte("movestogo")
)

func (e *Engine) ReceiveCommand(message []byte) {
	message = bytes.TrimSpace(message)
	messageType := bytes.ToLower(bytes.Split(message, []byte{' '})[0])

	if bytes.Equal(messageType, UCI) {
		fmt.Println("uciok")
	} else if bytes.Equal(messageType, ISREADY) {
		fmt.Println("readyok")
	} else if bytes.Equal(messageType, UCINEWGAME) {
		e.parsePosition(STARTPOS)
		clear_hash_table()
	} else if bytes.Equal(messageType, POSITION) {
		e.parsePosition(message[9:])
		clear_hash_table()
	} else if bytes.Equal(messageType, GO) {
		e.parseGo(message[3:])
	}
}
