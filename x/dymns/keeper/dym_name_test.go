package keeper_test

import (
	testkeeper "github.com/dymensionxyz/dymension/v3/testutil/keeper"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

//goland:noinspection SpellCheckingInspection
func TestKeeper_GetSetDymName(t *testing.T) {
	dk, _, ctx := testkeeper.DymNSKeeper(t)

	owner := "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"

	dymName := dymnstypes.DymName{
		Name:       "bonded-pool",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   1,
		Configs: []dymnstypes.DymNameConfig{{
			Type:  dymnstypes.DymNameConfigType_NAME,
			Path:  "www",
			Value: owner,
		}},
	}

	err := dk.SetDymName(ctx, dymName)
	require.NoError(t, err)

	t.Run("reverse mapping must be set", func(t *testing.T) {
		ownedBy, err := dk.GetDymNamesOwnedBy(ctx, owner)
		require.NoError(t, err)
		require.Len(t, ownedBy, 1)
		require.EqualValues(t, dymName, ownedBy[0])
	})

	t.Run("event should be fired", func(t *testing.T) {
		events := ctx.EventManager().Events()
		require.NotEmpty(t, events)

		for _, event := range events {
			if event.Type == dymnstypes.EventTypeSetDymName {
				return
			}
		}

		t.Errorf("event %s not found", dymnstypes.EventTypeSetDymName)
	})

	t.Run("Dym-Name should be equals to original", func(t *testing.T) {
		require.EqualValues(t, dymName, *dk.GetDymName(ctx, dymName.Name))
	})

	t.Run("can not set invalid Dym-Name", func(t *testing.T) {
		require.Error(t, dk.SetDymName(ctx, dymnstypes.DymName{}))
	})

	t.Run("returns nil if non-exists", func(t *testing.T) {
		require.Nil(t, dk.GetDymName(ctx, "non-exists"))
	})
}

func TestKeeper_GetAllNonExpiredDymNames(t *testing.T) {
	dk, _, ctx := testkeeper.DymNSKeeper(t)

	//goland:noinspection SpellCheckingInspection
	dymName1 := dymnstypes.DymName{
		Name:       "bonded-pool",
		Owner:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		ExpireAt:   time.Now().UTC().Add(time.Hour).Unix(),
		Configs: []dymnstypes.DymNameConfig{{
			Type:  dymnstypes.DymNameConfigType_NAME,
			Path:  "www",
			Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		}},
	}
	require.NoError(t, dk.SetDymName(ctx, dymName1))

	//goland:noinspection SpellCheckingInspection
	dymName2 := dymnstypes.DymName{
		Name:       "not-bonded-pool",
		Owner:      "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d",
		Controller: "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d",
		ExpireAt:   time.Now().UTC().Add(time.Hour).Unix(),
		Configs: []dymnstypes.DymNameConfig{{
			Type:  dymnstypes.DymNameConfigType_NAME,
			Path:  "www",
			Value: "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d",
		}},
	}
	require.NoError(t, dk.SetDymName(ctx, dymName2))

	//goland:noinspection SpellCheckingInspection
	dymName3 := dymnstypes.DymName{
		Name:       "streamer",
		Owner:      "dym1ysjlrjcankjpmpxxzk27mvzhv25e266r80p5pv",
		Controller: "dym1ysjlrjcankjpmpxxzk27mvzhv25e266r80p5pv",
		ExpireAt:   time.Now().UTC().Add(-time.Hour).Unix(),
		Configs: []dymnstypes.DymNameConfig{{
			Type:  dymnstypes.DymNameConfigType_NAME,
			Path:  "www",
			Value: "dym1ysjlrjcankjpmpxxzk27mvzhv25e266r80p5pv",
		}},
	}
	require.NoError(t, dk.SetDymName(ctx, dymName3))

	list := dk.GetAllNonExpiredDymNames(ctx, time.Now().UTC().Unix())
	require.Len(t, list, 2)
	require.Contains(t, list, dymName1)
	require.Contains(t, list, dymName2)
	require.NotContains(t, list, dymName3, "should not include expired Dym-Name")
}

//goland:noinspection SpellCheckingInspection
func TestKeeper_GetSetReverseMappingOwnerToOwnedDymName(t *testing.T) {
	dk, _, ctx := testkeeper.DymNSKeeper(t)

	require.Error(
		t,
		dk.SetReverseMappingOwnerToOwnedDymName(ctx, "0x", "a"),
		"should not allow invalid owner address",
	)

	_, err := dk.GetDymNamesOwnedBy(ctx, "0x")
	require.Error(
		t,
		err,
		"should not allow invalid owner address",
	)

	owner1 := "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"
	owner2 := "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d"

	dymName1 := dymnstypes.DymName{
		Name:       "bonded-pool",
		Owner:      owner1,
		Controller: owner1,
		ExpireAt:   time.Now().UTC().Add(time.Hour).Unix(),
	}
	require.NoError(t, dk.SetDymName(ctx, dymName1))

	dymName2 := dymnstypes.DymName{
		Name:       "not-bonded-pool",
		Owner:      owner2,
		Controller: owner2,
		ExpireAt:   time.Now().UTC().Add(time.Hour).Unix(),
	}
	require.NoError(t, dk.SetDymName(ctx, dymName2))

	dymName3 := dymnstypes.DymName{
		Name:       "not-bonded-pool2",
		Owner:      owner2,
		Controller: owner2,
		ExpireAt:   time.Now().UTC().Add(time.Hour).Unix(),
	}
	require.NoError(t, dk.SetDymName(ctx, dymName3))

	require.NoError(
		t,
		dk.SetReverseMappingOwnerToOwnedDymName(ctx, owner2, "not-exists"),
		"no check non-existing dym-name",
	)

	require.NoError(
		t,
		dk.SetReverseMappingOwnerToOwnedDymName(ctx, owner2, dymName2.Name),
		"no error if duplicated name",
	)

	ownedBy1, err1 := dk.GetDymNamesOwnedBy(ctx, owner1)
	require.NoError(t, err1)
	require.Len(t, ownedBy1, 1)

	ownedBy2, err2 := dk.GetDymNamesOwnedBy(ctx, owner2)
	require.NoError(t, err2)
	require.NotEqual(t, 3, len(ownedBy2), "should not include non-existing dym-name")
	require.Len(t, ownedBy2, 2)

	ownedByNonExists, err3 := dk.GetDymNamesOwnedBy(ctx, "dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4")
	require.NoError(t, err3)
	require.Len(t, ownedByNonExists, 0)

	require.NoError(
		t,
		dk.SetReverseMappingOwnerToOwnedDymName(ctx, owner2, dymName1.Name),
		"no error if dym-name owned by another owner",
	)
	ownedBy2, err2 = dk.GetDymNamesOwnedBy(ctx, owner2)
	require.NoError(t, err2)
	require.Len(t, ownedBy2, 2, "should not include dym-name owned by another owner")
}

//goland:noinspection SpellCheckingInspection
func TestKeeper_RemoveReverseMappingOwnerToOwnedDymName(t *testing.T) {
	dk, _, ctx := testkeeper.DymNSKeeper(t)

	require.Error(
		t,
		dk.RemoveReverseMappingOwnerToOwnedDymName(ctx, "0x", "a"),
		"should not allow invalid owner address",
	)

	owner := "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"

	dymName1 := dymnstypes.DymName{
		Name:       "bonded-pool",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   time.Now().UTC().Add(time.Hour).Unix(),
	}
	require.NoError(t, dk.SetDymName(ctx, dymName1))

	dymName2 := dymnstypes.DymName{
		Name:       "a",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   time.Now().UTC().Add(time.Hour).Unix(),
	}
	require.NoError(t, dk.SetDymName(ctx, dymName2))

	require.NoError(
		t,
		dk.RemoveReverseMappingOwnerToOwnedDymName(ctx, "dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4", "a"),
		"no error if owner non-exists",
	)

	require.NoError(
		t,
		dk.RemoveReverseMappingOwnerToOwnedDymName(ctx, owner, "aaaaa"),
		"no error if not owned dym-name",
	)
	ownedBy, err := dk.GetDymNamesOwnedBy(ctx, owner)
	require.NoError(t, err)
	require.Len(t, ownedBy, 2, "existing data must be kept")

	require.NoError(
		t,
		dk.RemoveReverseMappingOwnerToOwnedDymName(ctx, owner, "not-exists"),
		"no error if not owned dym-name",
	)
	ownedBy, err = dk.GetDymNamesOwnedBy(ctx, owner)
	require.NoError(t, err)
	require.Len(t, ownedBy, 2, "existing data must be kept")

	require.NoError(
		t,
		dk.RemoveReverseMappingOwnerToOwnedDymName(ctx, owner, dymName1.Name),
		"no error if not owned dym-name",
	)
	ownedBy, err = dk.GetDymNamesOwnedBy(ctx, owner)
	require.NoError(t, err)
	require.Len(t, ownedBy, 1)

	require.NoError(
		t,
		dk.RemoveReverseMappingOwnerToOwnedDymName(ctx, owner, dymName2.Name),
		"no error if not owned dym-name",
	)
	ownedBy, err = dk.GetDymNamesOwnedBy(ctx, owner)
	require.NoError(t, err)
	require.Len(t, ownedBy, 0)
}
