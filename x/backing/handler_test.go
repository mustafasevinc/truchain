package backing

import (
	"encoding/binary"
	"strings"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestBackStoryMsg_FailBasicValidation(t *testing.T) {
	ctx, bk, _, _, _, _ := mockDB()

	h := NewHandler(bk)
	assert.NotNil(t, h)

	storyID := int64(1)
	amount, _ := sdk.ParseCoin("5trushane")
	creator := sdk.AccAddress([]byte{1, 2})
	duration := 5 * time.Hour
	msg := NewBackStoryMsg(storyID, amount, creator, duration)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	hasInvalidBackingPeriod := strings.Contains(res.Log, "901")
	assert.True(t, hasInvalidBackingPeriod, "should return err code")
}

func TestBackStoryMsg_FailInsufficientFunds(t *testing.T) {
	ctx, bk, _, _, _, _ := mockDB()

	h := NewHandler(bk)
	assert.NotNil(t, h)

	storyID := int64(1)
	amount, _ := sdk.ParseCoin("5trushane")
	creator := sdk.AccAddress([]byte{1, 2})
	duration := 99 * time.Hour
	msg := NewBackStoryMsg(storyID, amount, creator, duration)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	hasInsufficientFunds := strings.Contains(res.Log, "65541")
	assert.True(t, hasInsufficientFunds, "should return err code")
}

func TestBackStoryMsg(t *testing.T) {
	ctx, bk, sk, ck, _, am := mockDB()

	h := NewHandler(bk)
	assert.NotNil(t, h)

	storyID := createFakeStory(ctx, sk, ck)
	amount, _ := sdk.ParseCoin("5trudex")
	creator := createFakeFundedAccount(ctx, am, sdk.Coins{amount})
	duration := 99 * time.Hour
	msg := NewBackStoryMsg(storyID, amount, creator, duration)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	x, _ := binary.Varint(res.Data)
	assert.Equal(t, int64(1), x, "incorrect result backing id")
}