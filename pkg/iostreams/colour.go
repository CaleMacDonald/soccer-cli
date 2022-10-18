package iostreams

import (
	"fmt"
	"github.com/mgutz/ansi"
	"os"
	"strconv"
	"strings"
)

var (
	magenta  = ansi.ColorFunc("magenta")
	cyan     = ansi.ColorFunc("cyan")
	red      = ansi.ColorFunc("red")
	yellow   = ansi.ColorFunc("yellow")
	blue     = ansi.ColorFunc("blue")
	green    = ansi.ColorFunc("green")
	gray     = ansi.ColorFunc("black+h")
	bold     = ansi.ColorFunc("default+b")
	cyanBold = ansi.ColorFunc("cyan+b")

	gray256 = func(t string) string {
		return fmt.Sprintf("\x1b[%d;5;%dm%s\x1b[m", 38, 242, t)
	}
)

func EnvColourDisabled() bool {
	return os.Getenv("NO_COLOR") != "" || os.Getenv("CLICOLOR") == "0"
}

func EnvColourForced() bool {
	return os.Getenv("CLICOLOR_FORCE") != "" && os.Getenv("CLICOLOR_FORCE") != "0"
}

func Is256ColourSupported() bool {
	return IsTrueColourSupported() ||
		strings.Contains(os.Getenv("TERM"), "256") ||
		strings.Contains(os.Getenv("COLORTERM"), "256")
}

func IsTrueColourSupported() bool {
	term := os.Getenv("TERM")
	colorterm := os.Getenv("COLORTERM")

	return strings.Contains(term, "24bit") ||
		strings.Contains(term, "truecolor") ||
		strings.Contains(colorterm, "24bit") ||
		strings.Contains(colorterm, "truecolor")
}

func NewColourScheme(enabled, is256enabled bool, trueColor bool) *ColourScheme {
	return &ColourScheme{
		enabled:      enabled,
		is256enabled: is256enabled,
		hasTrueColor: trueColor,
	}
}

type ColourScheme struct {
	enabled      bool
	is256enabled bool
	hasTrueColor bool
}

func (c *ColourScheme) Bold(t string) string {
	if !c.enabled {
		return t
	}
	return bold(t)
}

func (c *ColourScheme) Boldf(t string, args ...interface{}) string {
	return c.Bold(fmt.Sprintf(t, args...))
}

func (c *ColourScheme) Red(t string) string {
	if !c.enabled {
		return t
	}
	return red(t)
}

func (c *ColourScheme) Redf(t string, args ...interface{}) string {
	return c.Red(fmt.Sprintf(t, args...))
}

func (c *ColourScheme) Yellow(t string) string {
	if !c.enabled {
		return t
	}
	return yellow(t)
}

func (c *ColourScheme) Yellowf(t string, args ...interface{}) string {
	return c.Yellow(fmt.Sprintf(t, args...))
}

func (c *ColourScheme) Green(t string) string {
	if !c.enabled {
		return t
	}
	return green(t)
}

func (c *ColourScheme) Greenf(t string, args ...interface{}) string {
	return c.Green(fmt.Sprintf(t, args...))
}

func (c *ColourScheme) Gray(t string) string {
	if !c.enabled {
		return t
	}
	if c.is256enabled {
		return gray256(t)
	}
	return gray(t)
}

func (c *ColourScheme) Grayf(t string, args ...interface{}) string {
	return c.Gray(fmt.Sprintf(t, args...))
}

func (c *ColourScheme) Magenta(t string) string {
	if !c.enabled {
		return t
	}
	return magenta(t)
}

func (c *ColourScheme) Magentaf(t string, args ...interface{}) string {
	return c.Magenta(fmt.Sprintf(t, args...))
}

func (c *ColourScheme) Cyan(t string) string {
	if !c.enabled {
		return t
	}
	return cyan(t)
}

func (c *ColourScheme) Cyanf(t string, args ...interface{}) string {
	return c.Cyan(fmt.Sprintf(t, args...))
}

func (c *ColourScheme) CyanBold(t string) string {
	if !c.enabled {
		return t
	}
	return cyanBold(t)
}

func (c *ColourScheme) Blue(t string) string {
	if !c.enabled {
		return t
	}
	return blue(t)
}

func (c *ColourScheme) Bluef(t string, args ...interface{}) string {
	return c.Blue(fmt.Sprintf(t, args...))
}

func (c *ColourScheme) SuccessIcon() string {
	return c.SuccessIconWithColor(c.Green)
}

func (c *ColourScheme) SuccessIconWithColor(colo func(string) string) string {
	return colo("✓")
}

func (c *ColourScheme) WarningIcon() string {
	return c.Yellow("!")
}

func (c *ColourScheme) FailureIcon() string {
	return c.FailureIconWithColor(c.Red)
}

func (c *ColourScheme) FailureIconWithColor(colo func(string) string) string {
	return colo("X")
}

func (c *ColourScheme) ColorFromString(s string) func(string) string {
	s = strings.ToLower(s)
	var fn func(string) string
	switch s {
	case "bold":
		fn = c.Bold
	case "red":
		fn = c.Red
	case "yellow":
		fn = c.Yellow
	case "green":
		fn = c.Green
	case "gray":
		fn = c.Gray
	case "magenta":
		fn = c.Magenta
	case "cyan":
		fn = c.Cyan
	case "blue":
		fn = c.Blue
	default:
		fn = func(s string) string {
			return s
		}
	}

	return fn
}

// ColorFromRGB returns a function suitable for TablePrinter.AddField
// that calls HexToRGB, coloring text if supported by the terminal.
func (c *ColourScheme) ColorFromRGB(hex string) func(string) string {
	return func(s string) string {
		return c.HexToRGB(hex, s)
	}
}

// HexToRGB uses the given hex to color x if supported by the terminal.
func (c *ColourScheme) HexToRGB(hex string, x string) string {
	if !c.enabled || !c.hasTrueColor || len(hex) != 6 {
		return x
	}

	r, _ := strconv.ParseInt(hex[0:2], 16, 64)
	g, _ := strconv.ParseInt(hex[2:4], 16, 64)
	b, _ := strconv.ParseInt(hex[4:6], 16, 64)
	return fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", r, g, b, x)
}
