package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgRegisterName_ValidateBasic(t *testing.T) {
	//goland:noinspection ALL
	tests := []struct {
		name            string
		DymName         string
		Duration        int32
		Owner           string
		wantErr         bool
		wantErrContains string
	}{
		{
			name:     "valid 1 yr",
			DymName:  "a",
			Duration: 1,
			Owner:    "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:     "valid 1+ yrs",
			DymName:  "a",
			Duration: 5,
			Owner:    "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:            "missing name",
			DymName:         "",
			Duration:        5,
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "name is not a valid dym name",
		},
		{
			name:            "name is too long",
			DymName:         "123456789012345678901",
			Duration:        5,
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: ErrDymNameTooLong.Error(),
		},
		{
			name:            "invalid name",
			DymName:         "-a",
			Duration:        5,
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "name is not a valid dym name",
		},
		{
			name:            "zero duration",
			DymName:         "a",
			Duration:        0,
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "duration must be at least 1 year",
		},
		{
			name:            "negative duration",
			DymName:         "a",
			Duration:        -1,
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "duration must be at least 1 year",
		},
		{
			name:            "empty owner",
			DymName:         "a",
			Duration:        1,
			Owner:           "",
			wantErr:         true,
			wantErrContains: "owner is not a valid bech32 account address",
		},
		{
			name:            "invalid owner",
			DymName:         "a",
			Duration:        1,
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4",
			wantErr:         true,
			wantErrContains: "owner is not a valid bech32 account address",
		},
		{
			name:            "owner must be dym1",
			DymName:         "a",
			Duration:        1,
			Owner:           "nim1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3pklgjx",
			wantErr:         true,
			wantErrContains: "owner is not a valid bech32 account address",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MsgRegisterName{
				Name:     tt.DymName,
				Duration: tt.Duration,
				Owner:    tt.Owner,
			}

			err := m.ValidateBasic()
			if tt.wantErr {
				require.NotEmpty(t, tt.wantErrContains, "mis-configured test case")
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
