package iostreams

import (
	"errors"
	"os"
)

func (s *IOStreams) EnableVirtualTerminalProcessing() error {
	return nil
}

func enableVirtualTerminalProcessing(_ *os.File) error {
	return errors.New("not implemented")
}
