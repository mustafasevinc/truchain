package backing

import (
	"time"

	t "github.com/TruStory/truchain/types"
	s "github.com/TruStory/truchain/x/story"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NativeTokenName is the coin  denom for Trustory's native reserve token.
const NativeTokenName = "trustake"

// BackStoryMsg defines a message to back a story. It implements the
// `Msg` interface which is required for transactions on Cosmos blockchains.
type BackStoryMsg struct {
	StoryID  int64          `json:"story_id"`
	Amount   sdk.Coin       `json:"amount"`
	Creator  sdk.AccAddress `json:"creator"`
	Duration time.Duration  `json:"duration"`
}

// NewBackStoryMsg creates a message to back a story
func NewBackStoryMsg(
	storyID int64,
	amount sdk.Coin,
	creator sdk.AccAddress,
	duration time.Duration) BackStoryMsg {
	return BackStoryMsg{
		StoryID:  storyID,
		Amount:   amount,
		Creator:  creator,
		Duration: duration,
	}
}

// Type implements Msg
func (msg BackStoryMsg) Type() string { return t.GetType(msg) }

// Name implements Msg
func (msg BackStoryMsg) Name() string { return t.GetName(msg) }

// GetSignBytes implements Msg
func (msg BackStoryMsg) GetSignBytes() []byte {
	return t.MustGetSignBytes(msg)
}

// ValidateBasic implements Msg
func (msg BackStoryMsg) ValidateBasic() sdk.Error {

	params := NewParams()

	if msg.StoryID <= 0 {
		return s.ErrInvalidStoryID(msg.StoryID)
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	if msg.Amount.IsZero() == true {
		return sdk.ErrInsufficientFunds("Invalid backing amount" + msg.Amount.String())
	}
	if msg.Duration < params.MinPeriod || msg.Duration > params.MaxPeriod {
		return ErrInvalidPeriod(msg.Duration)
	}
	return nil
}

// GetSigners implements Msg
func (msg BackStoryMsg) GetSigners() []sdk.AccAddress {
	return t.GetSigners(msg.Creator)
}