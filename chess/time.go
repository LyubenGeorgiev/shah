package chess

import "time"

type TimeController struct {
	// exit from engine flag
	Quit int

	// UCI "movestogo" command moves counter
	MovesToGo int

	// UCI "movetime" command time counter
	MoveTime int

	// UCI "time" command holder (ms)
	Time int

	// UCI "inc" command's time increment holder
	Inc int

	// UCI "starttime" command time holder
	StartTime int64

	// UCI "stoptime" command time holder
	StopTime int64

	// variable to flag time control availability
	Timeset int

	// variable to flag when the time is up
	Stopped int
}

func NewTimeController() TimeController {
	return TimeController{
		Quit:      0,
		MovesToGo: 30,
		MoveTime:  -1,
		Time:      -1,
		Inc:       0,
		StartTime: 0,
		StopTime:  0,
		Timeset:   0,
		Stopped:   0,
	}
}

// get time in milliseconds
func GetTimeMs() int64 {
	return time.Now().UnixMilli()
}

// a bridge function to interact between search and GUI input
func (e *Engine) timeout() {
	// if time is up break here
	if e.Timeset == 1 && GetTimeMs() > e.StopTime {
		// tell engine to stop calculating
		e.Stopped = 1
	}
}
