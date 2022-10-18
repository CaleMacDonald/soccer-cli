package iostreams

import (
	"bytes"
	"errors"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"io"
	"os"
	"sync"

	"github.com/briandowns/spinner"
	"github.com/muesli/termenv"
)

type ErrClosedPagerPipe struct {
	error
}

type IOStreams struct {
	In     io.ReadCloser
	Out    io.Writer
	ErrOut io.Writer

	// the original (non-colorable) output stream
	originalOut   io.Writer
	colorEnabled  bool
	is256enabled  bool
	hasTrueColor  bool
	terminalTheme string

	progressIndicatorEnabled bool
	progressIndicator        *spinner.Spinner
	progressIndicatorMu      sync.Mutex

	stdinTTYOverride  bool
	stdinIsTTY        bool
	stdoutTTYOverride bool
	stdoutIsTTY       bool
	stderrTTYOverride bool
	stderrIsTTY       bool
	termWidthOverride int
	ttySize           func() (int, int, error)

	pagerCommand string
	pagerProcess *os.Process
}

func (s *IOStreams) ColourEnabled() bool {
	return s.colorEnabled
}

func (s *IOStreams) ColourSupport256() bool {
	return s.is256enabled
}

func (s *IOStreams) HasTrueColour() bool {
	return s.hasTrueColor
}

// DetectTerminalTheme is a utility to call before starting the output pager so that the terminal background
// can be reliably detected.
func (s *IOStreams) DetectTerminalTheme() {
	if !s.ColourEnabled() {
		s.terminalTheme = "none"
		return
	}

	if s.pagerProcess != nil {
		s.terminalTheme = "none"
		return
	}

	style := os.Getenv("GLAMOUR_STYLE")
	if style != "" && style != "auto" {
		s.terminalTheme = "none"
		return
	}

	if termenv.HasDarkBackground() {
		s.terminalTheme = "dark"
		return
	}

	s.terminalTheme = "light"
}

// TerminalTheme returns "light", "dark", or "none" depending on the background color of the terminal.
func (s *IOStreams) TerminalTheme() string {
	if s.terminalTheme == "" {
		s.DetectTerminalTheme()
	}

	return s.terminalTheme
}

func (s *IOStreams) SetColorEnabled(colorEnabled bool) {
	s.colorEnabled = colorEnabled
}

func (s *IOStreams) ColorScheme() *ColourScheme {
	return NewColourScheme(s.ColourEnabled(), s.ColourSupport256(), s.HasTrueColour())
}

func (s *IOStreams) SetStdinTTY(isTTY bool) {
	s.stdinTTYOverride = true
	s.stdinIsTTY = isTTY
}

func System() *IOStreams {
	stdoutIsTTY := isTerminal(os.Stdout)
	stderrIsTTY := isTerminal(os.Stderr)

	assumeTrueColor := false
	if stdoutIsTTY {
		if err := enableVirtualTerminalProcessing(os.Stdout); err == nil {
			assumeTrueColor = true
		}
	}

	io := &IOStreams{
		In:           os.Stdin,
		originalOut:  os.Stdout,
		Out:          colorable.NewColorable(os.Stdout),
		ErrOut:       colorable.NewColorable(os.Stderr),
		colorEnabled: EnvColourForced() || (!EnvColourDisabled() && stdoutIsTTY),
		is256enabled: assumeTrueColor || Is256ColourSupported(),
		hasTrueColor: assumeTrueColor || IsTrueColourSupported(),
		pagerCommand: os.Getenv("PAGER"),
		ttySize:      ttySize,
	}

	if stdoutIsTTY && stderrIsTTY {
		io.progressIndicatorEnabled = true
	}

	// prevent duplicate isTerminal queries now that we know the answer
	io.SetStdoutTTY(stdoutIsTTY)
	io.SetStderrTTY(stderrIsTTY)
	return io
}

func Test() (ios *IOStreams, in *bytes.Buffer, out *bytes.Buffer, errOut *bytes.Buffer) {
	in = &bytes.Buffer{}
	out = &bytes.Buffer{}
	errOut = &bytes.Buffer{}
	ios = &IOStreams{
		In: &fdReader{
			fd:         0,
			ReadCloser: io.NopCloser(in),
		},
		Out:    &fdWriter{fd: 1, Writer: out},
		ErrOut: errOut,
		ttySize: func() (int, int, error) {
			return -1, -1, errors.New("ttySize not implemented in tests")
		},
	}
	ios.SetStdinTTY(false)
	ios.SetStdoutTTY(false)
	ios.SetStderrTTY(false)
	return ios, in, out, errOut
}

func isTerminal(f *os.File) bool {
	return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
}

func (s *IOStreams) SetStdoutTTY(isTTY bool) {
	s.stdoutTTYOverride = true
	s.stdoutIsTTY = isTTY
}

func (s *IOStreams) IsStdoutTTY() bool {
	if s.stdoutTTYOverride {
		return s.stdoutIsTTY
	}
	if stdout, ok := s.Out.(*os.File); ok {
		return isTerminal(stdout)
	}
	return false
}

func (s *IOStreams) SetStderrTTY(isTTY bool) {
	s.stderrTTYOverride = true
	s.stderrIsTTY = isTTY
}

func (s *IOStreams) IsStderrTTY() bool {
	if s.stderrTTYOverride {
		return s.stderrIsTTY
	}
	if stderr, ok := s.ErrOut.(*os.File); ok {
		return isTerminal(stderr)
	}
	return false
}

// fdWriter represents a wrapped stdout Writer that preserves the original file descriptor
type fdWriter struct {
	io.Writer
	fd uintptr
}

func (w *fdWriter) Fd() uintptr {
	return w.fd
}

// fdWriteCloser represents a wrapped stdout Writer that preserves the original file descriptor
type fdWriteCloser struct {
	io.WriteCloser
	fd uintptr
}

func (w *fdWriteCloser) Fd() uintptr {
	return w.fd
}

// fdWriter represents a wrapped stdin ReadCloser that preserves the original file descriptor
type fdReader struct {
	io.ReadCloser
	fd uintptr
}

func (r *fdReader) Fd() uintptr {
	return r.fd
}
