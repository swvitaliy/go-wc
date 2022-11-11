package mapReduceBigFile

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"unicode"
)

const bufSize int = 4096
const N = 5

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
	fileInfo, _ := os.Stat("words.txt")
	totalSize := fileInfo.Size()

	chunkSize := totalSize / N
	lastChunkSize := totalSize - (N-1)*chunkSize

	resp := make(chan map[string]int, N)
	wg := sync.WaitGroup{}
	for i := 0; i < N; i++ {
		file, err := os.Open("words.txt")
		if err != nil {
			log.Fatalf("can't open file \"words.txt\"")
		}
		defer file.Close()

		offset := int64(i) * chunkSize
		file.Seek(offset, 0)

		var cs int64
		if i == N-1 {
			cs = lastChunkSize
		} else {
			cs = chunkSize
		}

		r0 := io.LimitedReader{R: file, N: cs}
		r := bufio.NewReaderSize(&r0, bufSize)

		wg.Add(1)
		go func(r *bufio.Reader) {
			defer wg.Done()

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

			resp <- wc
		}(r)
	}

	wg.Wait()
	close(resp)

	wc := make(map[string]int)

	total := 0
	for m := range resp {
		for k, v := range m {
			wc[k] += v
			total += v
		}
	}

	// log.Printf("total=%d", total)
	// log.Printf("total_uniq=%d", len(wc))

	for w, c := range wc {
		fmt.Printf("%s %d\n", w, c)
		//fmt.Printf("%d %s\n", c, w)
	}
}
