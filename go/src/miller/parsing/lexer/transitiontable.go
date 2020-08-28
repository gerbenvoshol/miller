// Code generated by gocc; DO NOT EDIT.

package lexer

/*
Let s be the current state
Let r be the current input rune
transitionTable[s](r) returns the next state.
*/
type TransitionTable [NumStates]func(rune) int

var TransTab = TransitionTable{
	// S0
	func(r rune) int {
		switch {
		case r == 9: // ['\t','\t']
			return 1
		case r == 10: // ['\n','\n']
			return 1
		case r == 13: // ['\r','\r']
			return 1
		case r == 32: // [' ',' ']
			return 1
		case r == 36: // ['$','$']
			return 2
		case r == 45: // ['-','-']
			return 3
		case 48 <= r && r <= 57: // ['0','9']
			return 4
		case r == 61: // ['=','=']
			return 5
		}
		return NoState
	},
	// S1
	func(r rune) int {
		switch {
		}
		return NoState
	},
	// S2
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 6
		case 65 <= r && r <= 90: // ['A','Z']
			return 7
		case r == 95: // ['_','_']
			return 8
		case 97 <= r && r <= 122: // ['a','z']
			return 7
		}
		return NoState
	},
	// S3
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 9
		}
		return NoState
	},
	// S4
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 4
		}
		return NoState
	},
	// S5
	func(r rune) int {
		switch {
		}
		return NoState
	},
	// S6
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 6
		case 65 <= r && r <= 90: // ['A','Z']
			return 7
		case r == 95: // ['_','_']
			return 8
		case 97 <= r && r <= 122: // ['a','z']
			return 7
		}
		return NoState
	},
	// S7
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 6
		case 65 <= r && r <= 90: // ['A','Z']
			return 7
		case r == 95: // ['_','_']
			return 8
		case 97 <= r && r <= 122: // ['a','z']
			return 7
		}
		return NoState
	},
	// S8
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 6
		case 65 <= r && r <= 90: // ['A','Z']
			return 7
		case r == 95: // ['_','_']
			return 8
		case 97 <= r && r <= 122: // ['a','z']
			return 7
		}
		return NoState
	},
	// S9
	func(r rune) int {
		switch {
		case 48 <= r && r <= 57: // ['0','9']
			return 9
		}
		return NoState
	},
}
