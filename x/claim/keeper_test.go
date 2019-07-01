package claim

import (
	"net/url"
	"testing"
	"time"

	// sdk "github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestAddGetClaim(t *testing.T) {
	ctx, keeper := mockDB()

	// test getting a non-existent claim
	claim, ok := keeper.Claim(ctx, uint64(5))
	assert.False(t, ok)
	assert.Equal(t, Claim{}, claim)

	claim = createFakeClaim(ctx, keeper)

	// test getting an existing claim
	_, ok = keeper.Claim(ctx, claim.ID)
	assert.True(t, ok)

	// test incrementing id by adding another story
	claim = createFakeClaim(ctx, keeper)
	assert.Equal(t, uint64(2), claim.ID)

	claims := keeper.CreatorClaims(ctx, claim.Creator)
	assert.Len(t, claims, 2)

	claims = keeper.CommunityClaims(ctx, claim.CommunityID)
	assert.Len(t, claims, 2)

	claims = keeper.ClaimsBetweenIDs(ctx, 0, 100)
	assert.Len(t, claims, 2)

	claims = keeper.ClaimsBetweenIDs(ctx, 2, 100)
	assert.Len(t, claims, 1)

	tt := time.Now().UTC()
	claims = keeper.ClaimsAfterTime(ctx, tt)
	assert.Len(t, claims, 0)

	tt = tt.Add(-60 * time.Minute)
	claims = keeper.ClaimsAfterTime(ctx, tt)
	assert.Len(t, claims, 2)

	claims = keeper.ClaimsBeforeTime(ctx, tt)
	assert.Len(t, claims, 0)

	tt = tt.Add(60 * 60 * time.Minute)
	claims = keeper.ClaimsBeforeTime(ctx, tt)
	assert.Len(t, claims, 2)
}

func createFakeClaim(ctx sdk.Context, keeper Keeper) Claim {
	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Now().UTC()})
	body := "Preethi can handle liquor better than Aamir."
	communityID := "crypto"
	creator := sdk.AccAddress([]byte{1, 2})
	source := url.URL{}

	claim, err := keeper.SubmitClaim(ctx, body, communityID, creator, source)
	if err != nil {
		panic(err)
	}

	return claim
}