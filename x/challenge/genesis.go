package challenge

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all story state that must be provided at genesis
type GenesisState struct {
	Challenges []Challenge `json:"challenges"`
	Params     Params      `json:"params"`
}

// DefaultGenesisState for tests
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

// InitGenesis initializes story state from genesis file
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, challenge := range data.Challenges {
		keeper.setChallenge(ctx, challenge)
	}
	keeper.SetLen(ctx, int64(len(data.Challenges)))
	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	var challenges []Challenge
	keeper.Each(ctx, func(bz []byte) bool {
		var c Challenge
		keeper.GetCodec().MustUnmarshalBinaryLengthPrefixed(bz, &c)
		challenges = append(challenges, c)
		return true
	})

	params := keeper.GetParams(ctx)

	return GenesisState{
		Challenges: challenges,
		Params:     params,
	}
}
