package cmdutil

import (
	"github.com/CaleMacDonald/soccer-cli/internal/config"
	"github.com/CaleMacDonald/soccer-cli/pkg/iostreams"
)

type Factory struct {
	IOStreams *iostreams.IOStreams
	Config    func() (config.Config, error)

	ExecutableName string
}
