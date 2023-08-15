// Package stuttering provides functions to check for and identify stuttering words in a string.
// TODO(Idelchi): Read https://gophersnippets.com/how-to-parse-comments-from-go-code
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
	for i := 0; i < len(words)-1; i++ {
		firstWord := words[i]
		secondWord := words[i+1]

		if strings.ToLower(firstWord) == normalize(secondWord) {
			stutter := fmt.Sprintf("%s %s", firstWord, secondWord)
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
		var result []string

		i := 0
		for i < len(words) {
			// Check if the next word exists and if it forms a stutter with the current word.
			if i < len(words)-1 && strings.ToLower(words[i]) == normalize(words[i+1]) {
				// If a stutter is found, add the second word (with its non-alphabet characters) to the result.
				result = append(result, words[i+1])
				i += 2 // Move to the word after the stutter.
			} else {
				// If no stutter is found, add the current word to the result.
				result = append(result, words[i])
				i++
			}
		}

		line = strings.Join(result, " ")
	}

	return line
}

// tokenize converts a string into a slice of words.
// This is a helper function and assumes words are separated by whitespace.
func tokenize(line string) []string {
	return strings.FieldsFunc(line, func(c rune) bool {
		return unicode.IsSpace(c)
	})
}

// normalize removes non-alphabet characters from the end of a word and converts to lower-case.
func normalize(word string) string {
	for len(word) > 0 && !unicode.IsLetter(rune(word[len(word)-1])) {
		word = word[:len(word)-1]
	}
	return strings.ToLower(word)
}
