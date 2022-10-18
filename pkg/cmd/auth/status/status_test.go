package status

import (
	"bytes"
	"github.com/CaleMacDonald/soccer-cli/internal/config"
	"github.com/CaleMacDonald/soccer-cli/pkg/cmdutil"
	"github.com/CaleMacDonald/soccer-cli/pkg/iostreams"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func Test_NewCmdStatus(t *testing.T) {
	tests := []struct {
		name  string
		cli   string
		wants Options
	}{
		{
			name:  "no arguments",
			cli:   "",
			wants: Options{},
		},
		{
			name: "show token",
			cli:  "--show-token",
			wants: Options{
				ShowToken: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &cmdutil.Factory{}

			argv, err := shlex.Split(tt.cli)
			assert.NoError(t, err)

			var gotOpts *Options
			cmd := NewCmdStatus(f, func(opts *Options) error {
				gotOpts = opts
				return nil
			})

			// TODO cobra hack-around
			cmd.Flags().BoolP("help", "x", false, "")

			cmd.SetArgs(argv)
			cmd.SetIn(&bytes.Buffer{})
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})

			_, err = cmd.ExecuteC()
			assert.NoError(t, err)

			assert.Equal(t, tt.wants.ShowToken, gotOpts.ShowToken)
		})
	}
}

func Test_statusRun(t *testing.T) {
	tests := []struct {
		name       string
		opts       *Options
		cfgStubs   func(*config.ConfigMock)
		wantErr    string
		wantErrOut *regexp.Regexp
		wantStdOut *regexp.Regexp
	}{
		{
			name: "no auth token",
			opts: &Options{},
			cfgStubs: func(c *config.ConfigMock) {
				_ = c.SetAuthToken("")
			},
			wantErr:    "SilentError",
			wantErrOut: regexp.MustCompile("no authentication token is set\n"),
		},
		{
			name: "all good",
			opts: &Options{},
			cfgStubs: func(c *config.ConfigMock) {
				_ = c.SetAuthToken("abc123")
			},
			wantStdOut: regexp.MustCompile(`authentication token is present\n`),
		},
		{
			name: "show token",
			opts: &Options{
				ShowToken: true,
			},
			cfgStubs: func(c *config.ConfigMock) {
				_ = c.SetAuthToken("abc123")
			},
			wantStdOut: regexp.MustCompile(`authentication token is present and set to abc123\n`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.opts == nil {
				tt.opts = &Options{}
			}

			ios, _, stdout, stderr := iostreams.Test()

			ios.SetStdinTTY(true)
			ios.SetStderrTTY(true)
			ios.SetStdoutTTY(true)
			tt.opts.IO = ios

			cfg := config.NewTestConfig()
			if tt.cfgStubs != nil {
				tt.cfgStubs(cfg)
			}
			tt.opts.Config = func() (config.Config, error) {
				return cfg, nil
			}

			err := statusRun(tt.opts)
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				return
			} else {
				assert.NoError(t, err)
			}

			if tt.wantErrOut == nil {
				assert.Equal(t, "", stderr.String())
			} else {
				assert.True(t, tt.wantErrOut.MatchString(stderr.String()))
			}

			if tt.wantStdOut != nil {
				assert.True(t, tt.wantStdOut.MatchString(stdout.String()))
			}
		})
	}
}
