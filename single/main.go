package single

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

const bufSize int = 4096

func isAlphaDigit(r rune) bool {
	return r >= '0' && r <= '9' || unicode.IsLetter(r)
}

func skipWhitespaces(r *bufio.Reader) (rune, error) {
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			return 0, err
		}
		if isAlphaDigit(c) {
			return c, nil
		}
	}
}

func readWord(r *bufio.Reader) ([]rune, error) {
	c, err := skipWhitespaces(r)
	if err != nil {
		return nil, err
	}

	word := make([]rune, 1)
	word[0] = c
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			return word, err
		}
		if !isAlphaDigit(c) {
			break
		}
		word = append(word, c)
	}

	return word, nil
}

func PrintWC() {
	file, err := os.Open("words.txt")
	if err != nil {
		panic("can't open file \"words.txt\"")
	}
	defer file.Close()

	var r = bufio.NewReaderSize(file, bufSize)
	wc := make(map[string]int)
	for {
		w, err := readWord(r)
		if err != nil {
			if len(w) > 0 {
				wc[string(w)]++
			}
			break
		}

		wc[string(w)]++
	}

	for w, c := range wc {
		fmt.Printf("%s %d\n", w, c)
		//fmt.Printf("%d %s\n", c, w)
	}
}
