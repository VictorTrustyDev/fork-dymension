package keeper_test

import (
	testkeeper "github.com/dymensionxyz/dymension/v3/testutil/keeper"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
)

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
