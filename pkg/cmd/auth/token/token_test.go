package token

import (
	"github.com/CaleMacDonald/soccer-cli/internal/config"
	"github.com/CaleMacDonald/soccer-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCmdToken(t *testing.T) {
	tests := []struct {
		name       string
		opts       Options
		wantStdout string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "token",
			opts: Options{
				Config: func() (config.Config, error) {
					cfg := config.NewTestConfig()
					_ = cfg.SetAuthToken("token_123")
					return cfg, nil
				},
			},
			wantStdout: "token_123\n",
		},
		{
			name: "no token",
			opts: Options{
				Config: func() (config.Config, error) {
					cfg := config.NewTestConfig()
					return cfg, nil
				},
			},
			wantErr:    true,
			wantErrMsg: "no token found",
		},
	}

	for _, tt := range tests {
		ios, _, stdout, _ := iostreams.Test()
		tt.opts.IO = ios

		t.Run(tt.name, func(t *testing.T) {
			err := tokenRun(&tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErrMsg)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantStdout, stdout.String())
		})
	}
}
