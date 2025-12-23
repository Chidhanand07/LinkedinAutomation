package stealth

import (
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

// HumanType simulates realistic human typing with variable intervals,
// occasional typos with corrections, and natural rhythm variations
func HumanType(el *rod.Element, text string) {
	runes := []rune(text)
	typoChance := 0.05 // 5% chance of typo per character

	for i := 0; i < len(runes); i++ {
		c := runes[i]

		// Simulate occasional typos
		if rand.Float64() < typoChance && i > 0 {
			// Type wrong character
			wrongChar := getAdjacentKey(c)
			el.MustInput(string(wrongChar))
			time.Sleep(time.Duration(rand.Intn(150)+100) * time.Millisecond)

			// Realize mistake and backspace
			el.MustType(input.Backspace)
			time.Sleep(time.Duration(rand.Intn(200)+150) * time.Millisecond)

			// Correct with right character
			el.MustInput(string(c))
		} else {
			// Normal typing
			el.MustInput(string(c))
		}

		// Variable keystroke intervals mimicking human rhythm
		delay := getTypingDelay(c, i, len(runes))
		time.Sleep(delay)

		// Occasional brief pauses (thinking)
		if rand.Float64() < 0.15 && i < len(runes)-1 {
			time.Sleep(time.Duration(rand.Intn(400)+300) * time.Millisecond)
		}
	}
}

// TypeWithBackspace occasionally types extra and backspaces for realism
func TypeWithBackspace(el *rod.Element, text string) {
	runes := []rune(text)

	for i := 0; i < len(runes); i++ {
		c := runes[i]

		// 10% chance to type an extra character and backspace
		if rand.Float64() < 0.1 && i > 0 {
			extraChar := getAdjacentKey(c)
			el.MustInput(string(extraChar))
			time.Sleep(time.Duration(rand.Intn(120)+80) * time.Millisecond)
			el.MustType(input.Backspace)
			time.Sleep(time.Duration(rand.Intn(150)+100) * time.Millisecond)
		}

		el.MustInput(string(c))
		delay := getTypingDelay(c, i, len(runes))
		time.Sleep(delay)
	}
}

// getTypingDelay returns realistic variable delay between keystrokes
func getTypingDelay(char rune, index, totalLength int) time.Duration {
	baseDelay := rand.Intn(100) + 60 // 60-160ms base

	// Slower at start (thinking/positioning)
	if index < 3 {
		baseDelay += rand.Intn(80) + 40
	}

	// Slower before spaces (word boundaries)
	if char == ' ' {
		baseDelay += rand.Intn(60) + 30
	}

	// Slightly faster in middle of word (flow state)
	if index > 3 && index < totalLength-3 && char != ' ' {
		baseDelay -= rand.Intn(30)
	}

	// Special characters take longer (searching on keyboard)
	if isSpecialChar(char) {
		baseDelay += rand.Intn(150) + 100
	}

	return time.Duration(baseDelay) * time.Millisecond
}

// getAdjacentKey returns a keyboard-adjacent character for typo simulation
func getAdjacentKey(char rune) rune {
	// QWERTY keyboard layout adjacency map
	adjacency := map[rune][]rune{
		'a': {'q', 's', 'z'},
		'b': {'v', 'g', 'h', 'n'},
		'c': {'x', 'd', 'f', 'v'},
		'd': {'s', 'e', 'r', 'f', 'c', 'x'},
		'e': {'w', 'r', 'd', 's'},
		'f': {'d', 'r', 't', 'g', 'v', 'c'},
		'g': {'f', 't', 'y', 'h', 'b', 'v'},
		'h': {'g', 'y', 'u', 'j', 'n', 'b'},
		'i': {'u', 'o', 'k', 'j'},
		'j': {'h', 'u', 'i', 'k', 'm', 'n'},
		'k': {'j', 'i', 'o', 'l', 'm'},
		'l': {'k', 'o', 'p'},
		'm': {'n', 'j', 'k'},
		'n': {'b', 'h', 'j', 'm'},
		'o': {'i', 'p', 'l', 'k'},
		'p': {'o', 'l'},
		'q': {'w', 'a'},
		'r': {'e', 't', 'f', 'd'},
		's': {'a', 'w', 'e', 'd', 'x', 'z'},
		't': {'r', 'y', 'g', 'f'},
		'u': {'y', 'i', 'j', 'h'},
		'v': {'c', 'f', 'g', 'b'},
		'w': {'q', 'e', 's', 'a'},
		'x': {'z', 's', 'd', 'c'},
		'y': {'t', 'u', 'h', 'g'},
		'z': {'a', 's', 'x'},
	}

	lower := []rune(string(char))
	if len(lower) > 0 {
		if adjacent, ok := adjacency[lower[0]]; ok && len(adjacent) > 0 {
			return adjacent[rand.Intn(len(adjacent))]
		}
	}

	// Fallback to random character
	alternatives := []rune("abcdefghijklmnopqrstuvwxyz")
	return alternatives[rand.Intn(len(alternatives))]
}

// isSpecialChar checks if character requires extra time to type
func isSpecialChar(char rune) bool {
	return char == '@' || char == '#' || char == '$' || char == '%' ||
		char == '&' || char == '*' || char == '(' || char == ')' ||
		char == '_' || char == '+' || char == '=' || char == '!' ||
		char == '?' || char == '.' || char == ','
}
