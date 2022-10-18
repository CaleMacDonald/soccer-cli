package set

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

func Test_NewCmdSet(t *testing.T) {
	tests := []struct {
		name        string
		cli         string
		expectedErr string
		wants       Options
	}{
		{
			name:        "no arguments",
			cli:         "",
			expectedErr: "accepts 1 arg(s), received 0",
			wants:       Options{},
		},
		{
			name: "with token parameter",
			cli:  "abc123",
			wants: Options{
				ApiKey: "abc123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &cmdutil.Factory{}

			argv, err := shlex.Split(tt.cli)
			assert.NoError(t, err)

			var gotOpts *Options
			cmd := NewCmdSet(f, func(opts *Options) error {
				gotOpts = opts
				return nil
			})

			cmd.SetArgs(argv)
			cmd.SetIn(&bytes.Buffer{})
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})

			_, err = cmd.ExecuteC()
			if tt.expectedErr != "" {
				assert.Equal(t, tt.expectedErr, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wants.ApiKey, gotOpts.ApiKey)
			}
		})
	}
}

func Test_setRun(t *testing.T) {
	tests := []struct {
		name       string
		opts       *Options
		cfgStubs   func(*config.ConfigMock)
		wantErr    string
		wantErrOut *regexp.Regexp
		wantStdOut *regexp.Regexp
	}{
		{
			name:       "no auth token",
			opts:       &Options{},
			wantErr:    "the api key is not valid",
			wantErrOut: regexp.MustCompile("the api key is not valid\n"),
		},
		{
			name: "auth token in invalid format",
			opts: &Options{
				ApiKey: "abc123",
			},
			wantErr:    "the api key is not valid",
			wantErrOut: regexp.MustCompile("the api key is not valid\n"),
		},
		{
			name: "valid",
			opts: &Options{
				ApiKey: "a1b2c3",
			},
			wantStdOut: regexp.MustCompile("Auth token is now set"),
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

			err := setRun(tt.opts)
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
				assert.Regexp(t, tt.wantStdOut, stdout.String())
				//assert.True(t, tt.wantStdOut.MatchString(stdout.String()))
			}
		})
	}
}
