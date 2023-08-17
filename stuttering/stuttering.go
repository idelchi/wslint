// Package stuttering provides functions to check for and identify stuttering words in a string.
package stuttering

import (
	"fmt"
	"strings"
	"unicode"
)

// Has checks if a string has stuttering words.
// It returns true if the input string has stuttering words, false otherwise.
func Has(line string) bool {
	return len(Find(line)) > 0
}

// Find identifies and returns the stuttering words in a string.
// If there are no stuttering words, it returns an empty slice.
func Find(line string) []string {
	words := tokenize(line)
	stutters := []string{}

	// If there are less than two words, there can be no stuttering words.
	if len(words) <= 1 {
		return stutters
	}

	for i := 0; i < len(words)-1; i++ {
		firstWord := words[i]
		secondWord := words[i+1]

		if isStutteringPair(firstWord, secondWord) {
			stutter := fmt.Sprintf("(%s %s)", firstWord, secondWord)
			stutters = append(stutters, stutter)
		}
	}

	return stutters
}

// Trim removes the first word of all stuttering pairs from a string.
// It ensures that the second word of a stutter, including any trailing non-alphabetic characters, is retained.
// It returns a new string with the first words of stuttering pairs removed.
func Trim(line string) string {
	for Has(line) {
		words := tokenize(line)

		// If there are less than two words, there can be no stuttering words.
		if len(words) <= 1 {
			return line
		}

		var result []string

		for index := 0; index < len(words); index++ {
			// Check if the next word exists and if it forms a stutter with the current word.
			firstWord := words[index]

			if index < len(words)-1 {
				if secondWord := words[index+1]; isStutteringPair(firstWord, secondWord) {
					// If a stutter is found, add the second word (with its non-alphabet characters) to the result.
					result = append(result, secondWord)
					index++ // Move to the word after the stutter.

					continue
				}
			}
			// If no stutter is found, add the current word to the result.
			result = append(result, firstWord)
		}

		line = strings.Join(result, " ")
	}

	return line
}

// tokenize converts a string into a slice of words.
// This is a helper function and assumes words are separated by whitespace.
func tokenize(line string) []string {
	return strings.FieldsFunc(line, unicode.IsSpace)
}

// normalize removes non-alphabet characters from the end of a word and converts to lower-case.
func normalize2(word string) string {
	for len(word) > 0 && !unicode.IsLetter(rune(word[len(word)-1])) {
		word = word[:len(word)-1]
	}

	return strings.TrimSpace(strings.ToLower(word))
}

// normalize removes non-alphabet characters from the end of a word and converts to lower-case.
func normalize1(word string) string {
	return strings.TrimSpace(strings.ToLower(word))
}

// isStutteringPair checks if two words form a stuttering pair.
// At the moment, the first word is converted to lower-case and compared to the second word
// with non-alphabet characters removed.
func isStutteringPair(firstWord, secondWord string) bool {
	return normalize1(firstWord) == normalize2(secondWord)
}
