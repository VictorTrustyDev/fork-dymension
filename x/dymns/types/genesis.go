package types

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

func (m GenesisState) Validate() error {
	if err := (&m.Params).Validate(); err != nil {
		return ErrValidationFailed.Wrapf("params: %v", err)
	}

	for _, dymName := range m.DymNames {
		if err := dymName.Validate(); err != nil {
			return ErrValidationFailed.Wrapf("dym name '%s': %v", dymName.Name, err)
		}
	}

	for _, opoBid := range m.OpenPurchaseOrderBids {
		if err := opoBid.Validate(); err != nil {
			return ErrValidationFailed.Wrapf("open purchase order bid by '%s': %v", opoBid.Bidder, err)
		}
	}

	return nil
}
