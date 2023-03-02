// A package for miscellaneous functions that are used throughout the project.
package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
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
			f.Seek(0, io.SeekStart)
			return count, nil

		case err != nil:
			f.Seek(0, io.SeekStart)
			return count, err
		}
	}
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
         / /\         ` + "`" + `) )    / / ` + "`" + `"".` + "`" + `\
   , _.-'.'\ \        / /    ( (     / /
    ` + "`" + `--~` + "`" + `   ) )    .-'.'      '.'.  | (
           (/` + "`" + `    ( (` + "`" + `          ) )  '-;
            ` + "`" + `      '-;         (-'`
	fmt.Println(ascii)
}
