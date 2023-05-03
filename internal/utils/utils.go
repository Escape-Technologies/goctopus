// A package for miscellaneous functions that are used throughout the project.
package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
)

const (
	VERSION = "v0.0.11"
)

func MinInt(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func CountLines(f *os.File) (int, error) {
	buf := make([]byte, 32*1024)
	count := 1
	lineSep := []byte{'\n'}

	for {
		c, err := f.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			// Seek is used to reset the file pointer to the beginning of the file
			if _, err := f.Seek(0, io.SeekStart); err != nil {
				panic(err)
			}
			return count, nil

		case err != nil:
			panic(err)
		}
	}
}

func IsUrl(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func DomainFromUrl(s string) string {
	u, _ := url.Parse(s)
	return u.Host
}

func PrintASCII() {
	ascii := `
                    .-'   ` + "`" + `'.
                   /         \
                   |         ;
                   |         |           ___.--,
          _.._     |0) ~ (0) |    _.---'` + "`" + `__.-( (_.
   __.--'` + "`" + `_.. '.__.\    '--. \_.-' ,.--'` + "`" + `     ` + "`" + `""` + "`" + `
  ( ,.--'` + "`" + `   ',__ /./;   ;, '.__.'` + "`" + `    __
  _` + "`" + `) )  .---.__.' / |   |\   \__..--""  """--.,_
 ` + "`" + `---' .'.''-._.-'` + "`" + `_./  /\ '.  \ _.-~~~` + "`" + `` + "`" + `` + "`" + `` + "`" + `~~~-._` + "`" + `-.__.'
       | |  .' _.-' |  |  \  \  '.               ` + "`" + `~---` + "`" + `
        \ \/ .'     \  \   '. '-._)
         \/ /        \  \    ` + "`" + `=.__` + "`" + `~-.
     jgs / /\         ` + "`" + `) )    / / ` + "`" + `"".` + "`" + `\
   , _.-'.'\ \        / /    ( (     / /
    ` + "`" + `--~` + "`" + `   ) )    .-'.'      '.'.  | (
           (/` + "`" + `    ( (` + "`" + `          ) )  '-;
            ` + "`" + `      '-;         (-'
                  _                        
  __ _  ___   ___| |_ ___  _ __  _   _ ___ 
 / _` + "`" + ` |/ _ \ / __| __/ _ \| '_ \| | | / __|
| (_| | (_) | (__| || (_) | |_) | |_| \__ \
 \__, |\___/ \___|\__\___/| .__/ \__,_|___/ ` + VERSION + `
 |___/                    |_|`
	fmt.Println(ascii)
}
