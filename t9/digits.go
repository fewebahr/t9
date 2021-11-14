package t9

import (
	"errors"
	"fmt"
	"regexp"
	"unicode"
)

var (
	validDigits = regexp.MustCompile(`^[2-9]+$`)
)

// CheckDigits validates the designated digits and returns an error, if any.
func CheckDigits(digits string) error {
	if len(digits) == 0 {
		return errors.New(`digits are empty`)
	} else if !validDigits.MatchString(digits) {
		return fmt.Errorf(`digits are invalid (received :'%s')`, digits)
	}

	return nil
}

var digitToLetter = map[rune][]rune{
	'2': {'a', 'b', 'c'},
	'3': {'d', 'e', 'f'},
	'4': {'g', 'h', 'i'},
	'5': {'j', 'k', 'l'},
	'6': {'m', 'n', 'o'},
	'7': {'p', 'q', 'r', 's'},
	'8': {'t', 'u', 'v'},
	'9': {'w', 'x', 'y', 'z'},
}

var letterToDigit = func() map[rune]rune {
	// mutate letterToDigit so it goes in the reverse direction

	// size it for 26 letters of the alphabet, upper and lower case
	letterToDigit := make(map[rune]rune, 26*2)

	for digit, letters := range digitToLetter {

		for _, letter := range letters {
			// we want the map to be agnostic to case so we don't have to normalize
			// the case per-query. The map has constant lookup time regardless of
			// the size.
			letterToDigit[unicode.ToUpper(letter)] = digit
			letterToDigit[unicode.ToLower(letter)] = digit
		}
	}

	return letterToDigit
}()

func getDigits(word string) []rune {
	// Generally words will have the same number of digits as they have characters
	// ==> Exception case: words with punctuation. In those cases the character
	// has no digit equivalent. However we want those words to be find-able so we
	// will just ignore those characters when getting the digits.

	// So, as a rule, there are no more than len(w) digits
	digitsBuffer := make([]rune, len(word))

	// this gets incremented every time we add a rune to the digitsBuffer
	// critically, this is NOT incremented for non-alpha characters
	digitsBufferIndex := 0

	for _, r := range word {
		// r is type rune

		// case does not matter since letterToDigit includes both cases
		digit, ok := letterToDigit[r]
		if ok {
			// the rune in question
			digitsBuffer[digitsBufferIndex] = digit
			digitsBufferIndex++
		}
	}

	return digitsBuffer[:digitsBufferIndex]
}
