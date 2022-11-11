package mapReduceDir

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
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

type wcExit struct {
	WC   map[string]int
	Exit int
}

const (
	exitSig = 1
	exitOk  = 2
)

func PrintWC(dir string) {
	libRegEx, e := regexp.Compile("^.+\\.txt$")
	if e != nil {
		log.Fatal(e)
	}

	files := make([]string, 0)
	e = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && libRegEx.MatchString(info.Name()) {
			files = append(files, path)
		}
		return nil
	})

	if e != nil {
		log.Fatal(e)
	}

	//for _, f := range files {
	//	fmt.Println(f)
	//}

	// run workers (waiting for requests)
	resp := make(chan wcExit, N)
	wg := sync.WaitGroup{}

	// reduce goroutine
	wc := make(map[string]int)
	total := 0

	go func() {
		defer func() { resp <- wcExit{nil, exitOk} }()

		for {
			m := <-resp
			if m.Exit == exitSig {
				break
			}

			for k, v := range m.WC {
				wc[k] += v
				total += v
			}
			wg.Done()
		}
	}()

	// workers
	guard := make(chan struct{}, N)
	for _, fPath := range files {
		guard <- struct{}{}
		file, err := os.Open(fPath)
		if err != nil {
			log.Fatalf("can't open file \"%s\"", fPath)
		}
		defer file.Close()

		r := bufio.NewReaderSize(file, bufSize)
		wg.Add(1)
		go func(path string) {
			m := make(map[string]int)
			for {
				w, err := readWord(r)
				if err != nil {
					if len(w) > 0 {
						m[string(w)]++
					}
					break
				}

				m[string(w)]++
			}

			resp <- wcExit{WC: m, Exit: 0}
			<-guard
		}(fPath)
	}

	wg.Wait()

	resp <- wcExit{nil, exitSig}
	<-resp

	close(resp)

	log.Printf("total=%d", total)
	log.Printf("total_uniq=%d", len(wc))

	for w, c := range wc {
		//fmt.Printf("%s %d\n", w, c)
		fmt.Printf("%d %s\n", c, w)
	}
}
