package truchain

import (
	"encoding/binary"
	"reflect"

	db "github.com/TruStory/truchain/x/truchain/db"
	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for all TruStory messages
func NewHandler(k db.WriteKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case ts.SubmitStoryMsg:
			return handleSubmitStoryMsg(ctx, k, msg)
		case ts.BackStoryMsg:
			return handleBackStoryMsg(ctx, k, msg)
		default:
			errMsg := "Unrecognized Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// ============================================================================

// handleSubmitStoryMsg handles the logic of a SubmitStoryMsg
func handleSubmitStoryMsg(ctx sdk.Context, k db.WriteKeeper, msg ts.SubmitStoryMsg) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	storyID, err := k.NewStory(ctx, msg.Body, msg.Category, msg.Creator, msg.StoryType)
	if err != nil {
		return sdk.ErrInternal("Submit story error").Result()
	}

	return sdk.Result{Data: i2b(storyID)}
}

func handleBackStoryMsg(ctx sdk.Context, k db.WriteKeeper, msg ts.BackStoryMsg) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	backingID, err := k.NewBacking(
		ctx,
		msg.StoryID,
		msg.Amount,
		msg.Creator,
		msg.Duration)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Data: i2b(backingID)}
}

// i2b converts an int64 into a byte array
func i2b(x int64) []byte {
	var b [binary.MaxVarintLen64]byte
	return b[:binary.PutVarint(b[:], x)]
}