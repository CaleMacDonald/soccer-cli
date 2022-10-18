package iostreams

import (
	"golang.org/x/term"
	"os"
)

// ttySize measures the size of the controlling terminal for the current process
func ttySize() (int, int, error) {
	f, err := os.Open("/dev/tty")
	if err != nil {
		return -1, -1, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	return term.GetSize(int(f.Fd()))
}
