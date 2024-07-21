package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	err = dk.AddReverseMappingOwnerToOwnedDymName(ctx, owner1, dymName1.Name)
	require.NoError(t, err)

	dymName2 := dymnstypes.DymName{
		Name:       "not-bonded-pool",
		Owner:      owner2,
		Controller: owner2,
		ExpireAt:   time.Now().UTC().Add(time.Hour).Unix(),
	}
	require.NoError(t, dk.SetDymName(ctx, dymName2))
	err = dk.AddReverseMappingOwnerToOwnedDymName(ctx, owner2, dymName2.Name)
	require.NoError(t, err)

	dymName3 := dymnstypes.DymName{
		Name:       "not-bonded-pool2",
		Owner:      owner2,
		Controller: owner2,
		ExpireAt:   time.Now().UTC().Add(time.Hour).Unix(),
	}
	require.NoError(t, dk.SetDymName(ctx, dymName3))
	err = dk.AddReverseMappingOwnerToOwnedDymName(ctx, owner2, dymName3.Name)
	require.NoError(t, err)

	require.NoError(
		t,
		dk.AddReverseMappingOwnerToOwnedDymName(ctx, owner2, "not-exists"),
		"no check non-existing dym-name",
	)

	t.Run("no error if duplicated name", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			require.NoError(t,
				dk.AddReverseMappingOwnerToOwnedDymName(ctx, owner2, dymName2.Name),
			)
		}
	})

	ownedBy1, err1 := dk.GetDymNamesOwnedBy(ctx, owner1, 0)
	require.NoError(t, err1)
	require.Len(t, ownedBy1, 1)

	ownedBy2, err2 := dk.GetDymNamesOwnedBy(ctx, owner2, 0)
	require.NoError(t, err2)
	require.NotEqual(t, 3, len(ownedBy2), "should not include non-existing dym-name")
	require.Len(t, ownedBy2, 2)

	ownedByNonExists, err3 := dk.GetDymNamesOwnedBy(ctx, "dym1zg69v7yszg69v7yszg69v7yszg69v7ys8xdv96", 0)
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

	const owner = "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"

	dymName1 := dymnstypes.DymName{
		Name:       "bonded-pool",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   time.Now().UTC().Add(time.Hour).Unix(),
	}
	setDymNameWithFunctionsAfter(ctx, dymName1, t, dk)

	dymName2 := dymnstypes.DymName{
		Name:       "a",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   time.Now().UTC().Add(time.Hour).Unix(),
	}
	setDymNameWithFunctionsAfter(ctx, dymName2, t, dk)

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

//goland:noinspection SpellCheckingInspection
func TestKeeper_GetAddReverseMappingConfiguredAddressToDymName(t *testing.T) {
	dk, _, _, ctx := testkeeper.DymNSKeeper(t)

	require.Error(
		t,
		dk.AddReverseMappingConfiguredAddressToDymName(ctx, " ", "a"),
		"should not allow blank address",
	)

	_, err := dk.GetDymNamesContainsConfiguredAddress(ctx, " ", 0)
	require.Error(
		t,
		err,
		"should not allow invalid blank address",
	)

	owner1 := "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"
	owner2 := "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d"
	anotherAccount := "dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4"
	interchainAccount := "dym1zg69v7yszg69v7yszg69v7yszg69v7yszg69v7yszg69v7yszg6qrz80ul"

	dymName1 := dymnstypes.DymName{
		Name:       "bonded-pool",
		Owner:      owner1,
		Controller: owner1,
		ExpireAt:   time.Now().UTC().Add(time.Hour).Unix(),
		Configs: []dymnstypes.DymNameConfig{{
			Type:  dymnstypes.DymNameConfigType_NAME,
			Value: anotherAccount,
		}},
	}
	require.NoError(t, dk.SetDymName(ctx, dymName1))
	err = dk.AddReverseMappingConfiguredAddressToDymName(ctx, anotherAccount, dymName1.Name)
	require.NoError(t, err)

	dymName2 := dymnstypes.DymName{
		Name:       "not-bonded-pool",
		Owner:      owner2,
		Controller: owner2,
		ExpireAt:   time.Now().UTC().Add(time.Hour).Unix(),
	}
	require.NoError(t, dk.SetDymName(ctx, dymName2))
	err = dk.AddReverseMappingConfiguredAddressToDymName(ctx, owner2, dymName2.Name)
	require.NoError(t, err)

	dymName3 := dymnstypes.DymName{
		Name:       "not-bonded-pool2",
		Owner:      owner2,
		Controller: owner2,
		ExpireAt:   time.Now().UTC().Add(time.Hour).Unix(),
		Configs: []dymnstypes.DymNameConfig{{
			Type:  dymnstypes.DymNameConfigType_NAME,
			Value: anotherAccount,
		}},
	}
	require.NoError(t, dk.SetDymName(ctx, dymName3))
	err = dk.AddReverseMappingConfiguredAddressToDymName(ctx, anotherAccount, dymName3.Name)
	require.NoError(t, err)

	require.NoError(
		t,
		dk.AddReverseMappingConfiguredAddressToDymName(ctx, anotherAccount, "not-exists"),
		"no check non-existing dym-name",
	)

	t.Run("no error if duplicated name", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			require.NoError(t,
				dk.AddReverseMappingConfiguredAddressToDymName(ctx, owner2, dymName2.Name),
			)
		}
	})

	linked1, err1 := dk.GetDymNamesContainsConfiguredAddress(ctx, anotherAccount, 0)
	require.NoError(t, err1)
	require.Len(t, linked1, 2)
	requireEqualsStrings(t,
		[]string{dymName1.Name, dymName3.Name},
		[]string{linked1[0].Name, linked1[1].Name},
	)

	linked2, err2 := dk.GetDymNamesContainsConfiguredAddress(ctx, owner2, 0)
	require.NoError(t, err2)
	require.NotEqual(t, 2, len(linked2), "should not include non-existing dym-name")
	require.Len(t, linked2, 1)
	requireEqualsStrings(t,
		[]string{dymName2.Name},
		[]string{linked2[0].Name},
	)

	linkedByNotExists, err3 := dk.GetDymNamesContainsConfiguredAddress(ctx, "dym1zg69v7yszg69v7yszg69v7yszg69v7ys8xdv96", 0)
	require.NoError(t, err3)
	require.Len(t, linkedByNotExists, 0)

	t.Run("allow Interchain Account (32 bytes)", func(t *testing.T) {
		require.NoError(
			t,
			dk.AddReverseMappingConfiguredAddressToDymName(ctx, interchainAccount, dymName1.Name),
		)

		linked3, err := dk.GetDymNamesContainsConfiguredAddress(ctx, interchainAccount, 0)
		require.NoError(t, err)
		require.Len(t, linked3, 1)
		require.Equal(t, dymName1.Name, linked3[0].Name)
	})
}

//goland:noinspection SpellCheckingInspection
func TestKeeper_RemoveReverseMappingConfiguredAddressToDymName(t *testing.T) {
	dk, _, _, ctx := testkeeper.DymNSKeeper(t)

	require.Error(
		t,
		dk.RemoveReverseMappingConfiguredAddressToDymName(ctx, " ", "a"),
		"should not allow blank address",
	)

	const owner = "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"
	const anotherAccount = "dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4"
	const interchainAccount = "dym1zg69v7yszg69v7yszg69v7yszg69v7yszg69v7yszg69v7yszg6qrz80ul"

	dymName1 := dymnstypes.DymName{
		Name:       "bonded-pool",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   time.Now().UTC().Add(time.Hour).Unix(),
		Configs: []dymnstypes.DymNameConfig{{
			Type:  dymnstypes.DymNameConfigType_NAME,
			Value: anotherAccount,
		}},
	}
	require.NoError(t, dk.SetDymName(ctx, dymName1))
	err := dk.AddReverseMappingConfiguredAddressToDymName(ctx, anotherAccount, dymName1.Name)
	require.NoError(t, err)

	dymName2 := dymnstypes.DymName{
		Name:       "a",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   time.Now().UTC().Add(time.Hour).Unix(),
		Configs: []dymnstypes.DymNameConfig{{
			Type:  dymnstypes.DymNameConfigType_NAME,
			Value: anotherAccount,
		}},
	}
	require.NoError(t, dk.SetDymName(ctx, dymName2))
	err = dk.AddReverseMappingConfiguredAddressToDymName(ctx, anotherAccount, dymName2.Name)
	require.NoError(t, err)

	require.NoError(
		t,
		dk.RemoveReverseMappingConfiguredAddressToDymName(ctx, "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d", "a"),
		"no error if record not exists",
	)

	linked, err := dk.GetDymNamesContainsConfiguredAddress(ctx, anotherAccount, 0)
	require.NoError(t, err)
	require.Len(t, linked, 2, "existing data must be kept")

	t.Run("no error if element is not in the list", func(t *testing.T) {
		require.NoError(
			t,
			dk.RemoveReverseMappingConfiguredAddressToDymName(ctx, anotherAccount, "aaaaa"),
		)
		linked, err = dk.GetDymNamesContainsConfiguredAddress(ctx, anotherAccount, 0)
		require.NoError(t, err)
		require.Len(t, linked, 2, "existing data must be kept")
	})

	t.Run("remove correctly", func(t *testing.T) {
		require.NoError(
			t,
			dk.RemoveReverseMappingConfiguredAddressToDymName(ctx, anotherAccount, dymName1.Name),
		)

		linked, err = dk.GetDymNamesContainsConfiguredAddress(ctx, anotherAccount, 0)
		require.NoError(t, err)
		require.Len(t, linked, 1)
		require.Equal(t, dymName2.Name, linked[0].Name)

		require.NoError(
			t,
			dk.RemoveReverseMappingConfiguredAddressToDymName(ctx, anotherAccount, dymName2.Name),
		)

		linked, err = dk.GetDymNamesContainsConfiguredAddress(ctx, anotherAccount, 0)
		require.NoError(t, err)
		require.Empty(t, linked)
	})

	t.Run("remove correctly with Interchain Account (32 bytes)", func(t *testing.T) {
		require.NoError(
			t,
			dk.AddReverseMappingConfiguredAddressToDymName(ctx, interchainAccount, dymName1.Name),
		)

		linked3, err := dk.GetDymNamesContainsConfiguredAddress(ctx, interchainAccount, 0)
		require.NoError(t, err)
		require.Len(t, linked3, 1)

		require.NoError(
			t,
			dk.RemoveReverseMappingConfiguredAddressToDymName(ctx, interchainAccount, dymName1.Name),
		)

		linked, err = dk.GetDymNamesContainsConfiguredAddress(ctx, interchainAccount, 0)
		require.NoError(t, err)
		require.Empty(t, linked)
	})
}

//goland:noinspection SpellCheckingInspection
func TestKeeper_GetAddReverseMapping0xAddressToDymName(t *testing.T) {
	dk, _, _, ctx := testkeeper.DymNSKeeper(t)

	for size := 0; size <= 128; size++ {
		if size == 20 || size == 32 {
			continue // two valid size
		}

		var addr []byte

		addr = make([]byte, size)

		require.Errorf(
			t,
			dk.AddReverseMapping0xAddressToDymName(ctx, addr, "a"),
			"should not allow %d bytes address", size,
		)

		_, err := dk.GetDymNamesContains0xAddress(ctx, addr, 0)
		require.Errorf(
			t,
			err,
			"should not allow %d bytes address", size,
		)
	}

	owner1 := "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"
	owner2 := "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d"
	owner2AccAddr := sdk.MustAccAddressFromBech32(owner2)
	anotherAccount := "dym1zg69v7yszg69v7yszg69v7yszg69v7ys8xdv96"
	anotherAcc0xAddr := sdk.MustAccAddressFromBech32(anotherAccount)
	require.Len(t, anotherAcc0xAddr.Bytes(), 20)

	interchainAccountBech32Addr := "dym1zg69v7yszg69v7yszg69v7yszg69v7yszg69v7yszg69v7yszg6qrz80ul"
	interchainAccount0xAddr := sdk.MustAccAddressFromBech32(interchainAccountBech32Addr)
	require.Len(t, interchainAccount0xAddr.Bytes(), 32)

	dymName1 := dymnstypes.DymName{
		Name:       "bonded-pool",
		Owner:      owner1,
		Controller: owner1,
		ExpireAt:   time.Now().UTC().Add(time.Hour).Unix(),
		Configs: []dymnstypes.DymNameConfig{{
			Type:  dymnstypes.DymNameConfigType_NAME,
			Value: anotherAccount,
		}},
	}
	require.NoError(t, dk.SetDymName(ctx, dymName1))
	err := dk.AddReverseMapping0xAddressToDymName(ctx, anotherAcc0xAddr, dymName1.Name)
	require.NoError(t, err)

	dymName2 := dymnstypes.DymName{
		Name:       "not-bonded-pool",
		Owner:      owner2,
		Controller: owner2,
		ExpireAt:   time.Now().UTC().Add(time.Hour).Unix(),
	}
	require.NoError(t, dk.SetDymName(ctx, dymName2))
	err = dk.AddReverseMapping0xAddressToDymName(
		ctx,
		owner2AccAddr,
		dymName2.Name,
	)
	require.NoError(t, err)

	dymName3 := dymnstypes.DymName{
		Name:       "not-bonded-pool2",
		Owner:      owner2,
		Controller: owner2,
		ExpireAt:   time.Now().UTC().Add(time.Hour).Unix(),
		Configs: []dymnstypes.DymNameConfig{{
			Type:  dymnstypes.DymNameConfigType_NAME,
			Value: anotherAccount,
		}},
	}
	require.NoError(t, dk.SetDymName(ctx, dymName3))
	err = dk.AddReverseMapping0xAddressToDymName(ctx, anotherAcc0xAddr, dymName3.Name)
	require.NoError(t, err)

	require.NoError(
		t,
		dk.AddReverseMapping0xAddressToDymName(ctx, anotherAcc0xAddr, "not-exists"),
		"no check non-existing dym-name",
	)

	t.Run("no error if duplicated name", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			require.NoError(t,
				dk.AddReverseMapping0xAddressToDymName(ctx, owner2AccAddr, dymName2.Name),
			)
		}
	})

	linked1, err1 := dk.GetDymNamesContains0xAddress(ctx, anotherAcc0xAddr, 0)
	require.NoError(t, err1)
	require.Len(t, linked1, 2)
	requireEqualsStrings(t,
		[]string{dymName1.Name, dymName3.Name},
		[]string{linked1[0].Name, linked1[1].Name},
	)

	linked2, err2 := dk.GetDymNamesContains0xAddress(ctx, owner2AccAddr, 0)
	require.NoError(t, err2)
	require.NotEqual(t, 2, len(linked2), "should not include non-existing dym-name")
	require.Len(t, linked2, 1)
	requireEqualsStrings(t,
		[]string{dymName2.Name},
		[]string{linked2[0].Name},
	)

	linkedByNotExists, err3 := dk.GetDymNamesContains0xAddress(
		ctx,
		make([]byte, 20),
		0,
	)
	require.NoError(t, err3)
	require.Len(t, linkedByNotExists, 0)

	t.Run("allow Interchain Account (32 bytes)", func(t *testing.T) {
		require.NoError(
			t,
			dk.AddReverseMapping0xAddressToDymName(ctx, interchainAccount0xAddr.Bytes(), dymName1.Name),
		)

		linked3, err := dk.GetDymNamesContains0xAddress(ctx, interchainAccount0xAddr.Bytes(), 0)
		require.NoError(t, err)
		require.Len(t, linked3, 1)
		require.Equal(t, dymName1.Name, linked3[0].Name)
	})
}

//goland:noinspection SpellCheckingInspection
func TestKeeper_RemoveReverseMapping0xAddressToDymName(t *testing.T) {
	dk, _, _, ctx := testkeeper.DymNSKeeper(t)

	for size := 0; size <= 128; size++ {
		if size == 20 || size == 32 {
			continue // two valid size
		}

		var addr []byte

		addr = make([]byte, size)

		require.Errorf(
			t,
			dk.RemoveReverseMapping0xAddressToDymName(ctx, addr, "a"),
			"should not allow %d bytes address", size,
		)
	}

	const owner = "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"
	const anotherAccount = "dym1zg69v7yszg69v7yszg69v7yszg69v7ys8xdv96"
	anotherAcc0xAddr := sdk.MustAccAddressFromBech32(anotherAccount)

	interchainAccountBech32Addr := "dym1zg69v7yszg69v7yszg69v7yszg69v7yszg69v7yszg69v7yszg6qrz80ul"
	interchainAccount0xAddr := sdk.MustAccAddressFromBech32(interchainAccountBech32Addr)
	require.Len(t, interchainAccount0xAddr.Bytes(), 32)

	dymName1 := dymnstypes.DymName{
		Name:       "bonded-pool",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   time.Now().UTC().Add(time.Hour).Unix(),
		Configs: []dymnstypes.DymNameConfig{{
			Type:  dymnstypes.DymNameConfigType_NAME,
			Value: anotherAccount,
		}},
	}
	require.NoError(t, dk.SetDymName(ctx, dymName1))
	err := dk.AddReverseMapping0xAddressToDymName(ctx, anotherAcc0xAddr, dymName1.Name)
	require.NoError(t, err)

	dymName2 := dymnstypes.DymName{
		Name:       "a",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   time.Now().UTC().Add(time.Hour).Unix(),
		Configs: []dymnstypes.DymNameConfig{{
			Type:  dymnstypes.DymNameConfigType_NAME,
			Value: anotherAccount,
		}},
	}
	require.NoError(t, dk.SetDymName(ctx, dymName2))
	err = dk.AddReverseMapping0xAddressToDymName(ctx, anotherAcc0xAddr, dymName2.Name)
	require.NoError(t, err)

	require.NoError(
		t,
		dk.RemoveReverseMapping0xAddressToDymName(ctx,
			sdk.MustAccAddressFromBech32("dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d"),
			"a",
		),
		"no error if record not exists",
	)

	linked, err := dk.GetDymNamesContains0xAddress(ctx, anotherAcc0xAddr, 0)
	require.NoError(t, err)
	require.Len(t, linked, 2, "existing data must be kept")

	t.Run("no error if element is not in the list", func(t *testing.T) {
		require.NoError(
			t,
			dk.RemoveReverseMapping0xAddressToDymName(ctx, anotherAcc0xAddr, "aaaaa"),
		)
		linked, err = dk.GetDymNamesContains0xAddress(ctx, anotherAcc0xAddr, 0)
		require.NoError(t, err)
		require.Len(t, linked, 2, "existing data must be kept")
	})

	t.Run("remove correctly", func(t *testing.T) {
		require.NoError(
			t,
			dk.RemoveReverseMapping0xAddressToDymName(ctx, anotherAcc0xAddr, dymName1.Name),
		)

		linked, err = dk.GetDymNamesContains0xAddress(ctx, anotherAcc0xAddr, 0)
		require.NoError(t, err)
		require.Len(t, linked, 1)
		require.Equal(t, dymName2.Name, linked[0].Name)

		require.NoError(
			t,
			dk.RemoveReverseMapping0xAddressToDymName(ctx, anotherAcc0xAddr, dymName2.Name),
		)

		linked, err = dk.GetDymNamesContains0xAddress(ctx, anotherAcc0xAddr, 0)
		require.NoError(t, err)
		require.Empty(t, linked)
	})

	t.Run("allow Interchain Account (32 bytes)", func(t *testing.T) {
		require.NoError(
			t,
			dk.AddReverseMapping0xAddressToDymName(ctx, interchainAccount0xAddr, dymName1.Name),
		)

		linked3, err := dk.GetDymNamesContains0xAddress(ctx, interchainAccount0xAddr, 0)
		require.NoError(t, err)
		require.Len(t, linked3, 1)

		require.NoError(
			t,
			dk.RemoveReverseMapping0xAddressToDymName(ctx, interchainAccount0xAddr, dymName1.Name),
		)
		linked3, err = dk.GetDymNamesContains0xAddress(ctx, interchainAccount0xAddr, 0)
		require.NoError(t, err)
		require.Empty(t, linked3)
	})
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
