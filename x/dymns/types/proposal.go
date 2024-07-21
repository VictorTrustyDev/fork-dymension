package types

import (
	govcdc "github.com/cosmos/cosmos-sdk/x/gov/codec"
	v1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// constants

const (
	ProposalTypeMigrateChainIdsProposal string = "MigrateChainIdsProposal"
)

// Implements Proposal Interface
var (
	_ v1beta1.Content = &MigrateChainIdsProposal{}
)

func init() {
	v1beta1.RegisterProposalType(ProposalTypeMigrateChainIdsProposal)
	govcdc.ModuleCdc.Amino.RegisterConcrete(&MigrateChainIdsProposal{}, "dymns/"+ProposalTypeMigrateChainIdsProposal, nil)
}

// NewMigrateChainIdsProposal returns new instance of MigrateChainIdsProposal
func NewMigrateChainIdsProposal(title, description string, replacement ...MigrateChainId) v1beta1.Content {
	return &MigrateChainIdsProposal{
		Title:       title,
		Description: description,
		Replacement: replacement,
	}
}

// ProposalRoute returns router key for this proposal
func (*MigrateChainIdsProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for this proposal
func (*MigrateChainIdsProposal) ProposalType() string {
	return ProposalTypeMigrateChainIdsProposal
}
