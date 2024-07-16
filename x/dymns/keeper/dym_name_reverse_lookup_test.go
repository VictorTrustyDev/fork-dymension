package keeper_test

import (
	testkeeper "github.com/dymensionxyz/dymension/v3/testutil/keeper"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
	"time"
)

//goland:noinspection SpellCheckingInspection
func TestKeeper_GetAddReverseMappingOwnerToOwnedDymName(t *testing.T) {
	dk, _, _, ctx := testkeeper.DymNSKeeper(t)

	require.Error(
		t,
		dk.AddReverseMappingOwnerToOwnedDymName(ctx, "0x", "a"),
		"should not allow invalid owner address",
	)

	_, err := dk.GetDymNamesOwnedBy(ctx, "0x", 0)
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
		dk.AddReverseMappingOwnerToOwnedDymName(ctx, owner2, "not-exists"),
		"no check non-existing dym-name",
	)

	require.NoError(
		t,
		dk.AddReverseMappingOwnerToOwnedDymName(ctx, owner2, dymName2.Name),
		"no error if duplicated name",
	)

	ownedBy1, err1 := dk.GetDymNamesOwnedBy(ctx, owner1, 0)
	require.NoError(t, err1)
	require.Len(t, ownedBy1, 1)

	ownedBy2, err2 := dk.GetDymNamesOwnedBy(ctx, owner2, 0)
	require.NoError(t, err2)
	require.NotEqual(t, 3, len(ownedBy2), "should not include non-existing dym-name")
	require.Len(t, ownedBy2, 2)

	ownedByNonExists, err3 := dk.GetDymNamesOwnedBy(ctx, "dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4", 0)
	require.NoError(t, err3)
	require.Len(t, ownedByNonExists, 0)

	require.NoError(
		t,
		dk.AddReverseMappingOwnerToOwnedDymName(ctx, owner2, dymName1.Name),
		"no error if dym-name owned by another owner",
	)
	ownedBy2, err2 = dk.GetDymNamesOwnedBy(ctx, owner2, 0)
	require.NoError(t, err2)
	require.Len(t, ownedBy2, 2, "should not include dym-name owned by another owner")
}

//goland:noinspection SpellCheckingInspection
func TestKeeper_RemoveReverseMappingOwnerToOwnedDymName(t *testing.T) {
	dk, _, _, ctx := testkeeper.DymNSKeeper(t)

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
	ownedBy, err := dk.GetDymNamesOwnedBy(ctx, owner, 0)
	require.NoError(t, err)
	require.Len(t, ownedBy, 2, "existing data must be kept")

	require.NoError(
		t,
		dk.RemoveReverseMappingOwnerToOwnedDymName(ctx, owner, "not-exists"),
		"no error if not owned dym-name",
	)
	ownedBy, err = dk.GetDymNamesOwnedBy(ctx, owner, 0)
	require.NoError(t, err)
	require.Len(t, ownedBy, 2, "existing data must be kept")

	require.NoError(
		t,
		dk.RemoveReverseMappingOwnerToOwnedDymName(ctx, owner, dymName1.Name),
		"no error if not owned dym-name",
	)
	ownedBy, err = dk.GetDymNamesOwnedBy(ctx, owner, 0)
	require.NoError(t, err)
	require.Len(t, ownedBy, 1)

	require.NoError(
		t,
		dk.RemoveReverseMappingOwnerToOwnedDymName(ctx, owner, dymName2.Name),
		"no error if not owned dym-name",
	)
	ownedBy, err = dk.GetDymNamesOwnedBy(ctx, owner, 0)
	require.NoError(t, err)
	require.Len(t, ownedBy, 0)
}

var keyTestReverseLookupDymNames = append(
	dymnstypes.KeyPrefixRvlConfiguredAddressToDymNamesInclude,
	0x1, 0x2,
)

func TestKeeper_GenericAddReverseLookupDymNamesRecord(t *testing.T) {
	dk, _, _, ctx := testkeeper.DymNSKeeper(t)

	t.Run("should able to create new record for non-existing", func(t *testing.T) {
		record := dk.GenericGetReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames)
		require.Empty(t, record.DymNames)

		err := dk.GenericAddReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames, "test")
		require.NoError(t, err)

		record = dk.GenericGetReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames)
		require.Equal(t, 1, len(record.DymNames))
		require.Equal(t, "test", record.DymNames[0])
	})

	t.Run("should able to add more new record for existing", func(t *testing.T) {
		record := dk.GenericGetReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames)
		require.Equal(t, 1, len(record.DymNames))
		require.Equal(t, "test", record.DymNames[0])

		err := dk.GenericAddReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames, "test2")
		require.NoError(t, err)

		record = dk.GenericGetReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames)
		requireEqualsStrings(t, []string{"test", "test2"}, record.DymNames)
	})

	t.Run("should able to add duplicated record but not persist into store", func(t *testing.T) {
		record := dk.GenericGetReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames)
		requireEqualsStrings(t, []string{"test", "test2"}, record.DymNames)

		err := dk.GenericAddReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames, "test2")
		require.NoError(t, err)

		record = dk.GenericGetReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames)
		requireEqualsStrings(t, []string{"test", "test2"}, record.DymNames)
	})
}

func TestKeeper_GenericGetReverseLookupDymNamesRecord(t *testing.T) {
	dk, _, _, ctx := testkeeper.DymNSKeeper(t)

	t.Run("should able to get non-exist record", func(t *testing.T) {
		record := dk.GenericGetReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames)
		require.Empty(t, record.DymNames)
	})

	err := dk.GenericAddReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames, "test")
	require.NoError(t, err)

	record := dk.GenericGetReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames)
	require.Equal(t, 1, len(record.DymNames))
	require.Equal(t, "test", record.DymNames[0])

	err = dk.GenericAddReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames, "test2")
	require.NoError(t, err)

	record = dk.GenericGetReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames)
	requireEqualsStrings(t, []string{"test", "test2"}, record.DymNames)
}

func TestKeeper_GenericRemoveReverseLookupDymNamesRecord(t *testing.T) {
	t.Run("should able to remove non-existing record without error", func(t *testing.T) {
		dk, _, _, ctx := testkeeper.DymNSKeeper(t)

		err := dk.GenericRemoveReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames, "test")
		require.NoError(t, err)

		record := dk.GenericGetReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames)
		require.Empty(t, record.DymNames)
	})

	t.Run("should able to remove existing record, single name, without error", func(t *testing.T) {
		dk, _, _, ctx := testkeeper.DymNSKeeper(t)

		err := dk.GenericAddReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames, "test")
		require.NoError(t, err)

		record := dk.GenericGetReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames)
		require.Equal(t, 1, len(record.DymNames))
		require.Equal(t, "test", record.DymNames[0])

		err = dk.GenericRemoveReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames, "test")
		require.NoError(t, err)

		record = dk.GenericGetReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames)
		require.Empty(t, record.DymNames)
	})

	t.Run("should able to remove existing record, multiple names, without error", func(t *testing.T) {
		dk, _, _, ctx := testkeeper.DymNSKeeper(t)

		err := dk.GenericAddReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames, "test")
		require.NoError(t, err)

		err = dk.GenericAddReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames, "test2")
		require.NoError(t, err)

		record := dk.GenericGetReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames)
		requireEqualsStrings(t, []string{"test", "test2"}, record.DymNames)

		err = dk.GenericRemoveReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames, "test")
		require.NoError(t, err)

		record = dk.GenericGetReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames)
		require.Equal(t, 1, len(record.DymNames))
		require.Equal(t, "test2", record.DymNames[0])
	})

	t.Run("should able to remove in existing record, but name not in existing list, without error", func(t *testing.T) {
		dk, _, _, ctx := testkeeper.DymNSKeeper(t)

		err := dk.GenericAddReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames, "test")
		require.NoError(t, err)

		record := dk.GenericGetReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames)
		require.Equal(t, 1, len(record.DymNames))

		err = dk.GenericRemoveReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames, "test2")
		require.NoError(t, err)

		record = dk.GenericGetReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames)
		require.Equal(t, 1, len(record.DymNames))
		require.Equal(t, "test", record.DymNames[0])

		err = dk.GenericAddReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames, "test3")
		require.NoError(t, err)

		record = dk.GenericGetReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames)
		requireEqualsStrings(t, []string{"test", "test3"}, record.DymNames)

		err = dk.GenericRemoveReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames, "test4")
		require.NoError(t, err)

		record = dk.GenericGetReverseLookupDymNamesRecord(ctx, keyTestReverseLookupDymNames)
		requireEqualsStrings(t, []string{"test", "test3"}, record.DymNames)
	})
}

func requireEqualsStrings(t *testing.T, expected, actual []string) {
	t.Helper()

	require.Equal(t, len(expected), len(actual))

	sort.Strings(expected)
	sort.Strings(actual)

	require.Equal(t, expected, actual)
}
