package keeper_test

import (
	"sort"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dymensionxyz/dymension/v3/app/params"
	dymnskeeper "github.com/dymensionxyz/dymension/v3/x/dymns/keeper"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var dymNsModuleAccAddr = authtypes.NewModuleAddress(dymnstypes.ModuleName)

func setDymNameWithFunctionsAfter(ctx sdk.Context, dymName dymnstypes.DymName, t *testing.T, dk dymnskeeper.Keeper) {
	require.NoError(t, dk.SetDymName(ctx, dymName))
	require.NoError(t, dk.AfterDymNameOwnerChanged(ctx, dymName.Name))
	require.NoError(t, dk.AfterDymNameConfigChanged(ctx, dymName.Name))
}

func requireDymNameList(dymNames []dymnstypes.DymName, wantNames []string, t *testing.T, msgAndArgs ...any) {
	var gotNames []string
	for _, dymName := range dymNames {
		gotNames = append(gotNames, dymName.Name)
	}

	sort.Strings(gotNames)
	sort.Strings(wantNames)

	if len(wantNames) == 0 {
		wantNames = nil
	}

	require.Equal(t, wantNames, gotNames, msgAndArgs...)
}

func requireErrorContains(t *testing.T, err error, contains string) {
	require.Error(t, err)
	require.NotEmpty(t, contains, "mis-configured test")
	require.Contains(t, err.Error(), contains)
}

func requireErrorFContains(t *testing.T, f func() error, contains string) {
	requireErrorContains(t, f(), contains)
}

// ta stands for test-address, a simple wrapper for generating account for testing purpose.
// Usage is short, memorable, easy to type.
// The generated address is predictable, deterministic, supports output multiple formats.
type ta struct {
	bz []byte
}

// testAddr creates a general 20-bytes address from seed.
func testAddr(no uint64) ta {
	bz1 := sdk.Uint64ToBigEndian(no)
	bz2 := make([]byte, 20)
	copy(bz2, bz1)
	return ta{bz: bz2}
}

// testICAddr creates a 32-bytes address of Interchain Account from seed.
func testICAddr(no uint64) ta {
	bz1 := sdk.Uint64ToBigEndian(no)
	bz2 := make([]byte, 32)
	copy(bz2, bz1)
	return ta{bz: bz2}
}

func (a ta) bytes() []byte {
	return a.bz
}

func (a ta) bech32() string {
	return a.bech32C(params.AccountAddressPrefix)
}

func (a ta) bech32Valoper() string {
	return a.bech32C(params.AccountAddressPrefix + "valoper")
}

func (a ta) bech32C(customHrp string) string {
	return sdk.MustBech32ifyAddressBytes(customHrp, a.bz)
}

func (a ta) hexStr() string {
	if len(a.bz) == 20 {
		return common.BytesToAddress(a.bz).String()
	} else if len(a.bz) == 32 {
		return common.BytesToHash(a.bz).String()
	} else {
		panic("invalid length")
	}
}

type dymNameBuilder struct {
	name       string
	owner      string
	controller string
	expireAt   int64
	configs    []dymnstypes.DymNameConfig
}

func newDN(name, owner string) *dymNameBuilder {
	return &dymNameBuilder{
		name:       name,
		owner:      owner,
		controller: owner,
		expireAt:   time.Now().Unix() + 10,
		configs:    nil,
	}
}

func (m *dymNameBuilder) exp(now time.Time, offset int64) *dymNameBuilder {
	m.expireAt = now.Unix() + offset
	return m
}

func (m *dymNameBuilder) cfgN(chainId, path, resolveTo string) *dymNameBuilder {
	m.configs = append(m.configs, dymnstypes.DymNameConfig{
		Type:    dymnstypes.DymNameConfigType_NAME,
		ChainId: chainId,
		Path:    path,
		Value:   resolveTo,
	})
	return m
}

func (m *dymNameBuilder) build() dymnstypes.DymName {
	return dymnstypes.DymName{
		Name:       m.name,
		Owner:      m.owner,
		Controller: m.controller,
		ExpireAt:   m.expireAt,
		Configs:    m.configs,
	}
}

func (m *dymNameBuilder) buildSlice() []dymnstypes.DymName {
	return []dymnstypes.DymName{m.build()}
}
