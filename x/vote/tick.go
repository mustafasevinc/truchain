package vote

import (
	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/game"
	queue "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// NewResponseEndBlock is called at the end of every block tick
func (k Keeper) NewResponseEndBlock(ctx sdk.Context) sdk.Tags {
	store := ctx.KVStore(k.activeGamesQueueKey)
	q := queue.NewQueue(k.GetCodec(), store)

	err := checkGames(ctx, k, q)
	if err != nil {
		panic(err)
	}

	// TODO: maybe tags should return err?

	return sdk.NewTags()
}

// ============================================================================

// checkGames checks to see if a validation game has ended.
// It calls itself recursively until all games have been processed.
func checkGames(ctx sdk.Context, k Keeper, q queue.Queue) (err sdk.Error) {
	// check the head of the queue
	var gameID int64
	if err := q.Peek(&gameID); err != nil {
		return nil
	}

	// retrieve the game
	game, err := k.gameKeeper.Get(ctx, gameID)
	if err != nil {
		return err
	}

	// terminate recursion on finding the first non-ended game
	if game.Ended(ctx.BlockHeader().Time) {
		return nil
	}

	// remove ended game from queue
	q.Pop()

	// process ended game
	err = processGame(ctx, k, game)
	if err != nil {
		return err
	}

	// check next game
	return checkGames(ctx, k, q)
}

// tally votes and distribute rewards
func processGame(ctx sdk.Context, k Keeper, game game.Game) sdk.Error {
	// tally backings, challenges, and votes
	trueVotes, falseVotes, err := tally(ctx, k, game)
	if err != nil {
		return err
	}

	// check if story was confirmed
	confirmed := confirmStory(trueVotes, falseVotes)

	// calculate reward pool
	rewardPool, err := rewardPool(ctx, k.bankKeeper, trueVotes, falseVotes, confirmed)
	if err != nil {
		return err
	}

	// distribute rewards
	err = distributeRewards(ctx, k.bankKeeper, rewardPool, trueVotes, falseVotes, confirmed)
	if err != nil {
		return err
	}

	// update story state
	err = k.storyKeeper.EndGame(ctx, game.StoryID, confirmed)
	if err != nil {
		return err
	}

	return nil
}

func rewardPool(
	ctx sdk.Context,
	bankKeeper bank.Keeper,
	trueVotes []interface{},
	falseVotes []interface{},
	confirmed bool) (rewardPool sdk.Coin, err sdk.Error) {

	if confirmed {
		rewardPool, err = confirmedStoryRewardPool(ctx, bankKeeper, falseVotes)
	} else {
		rewardPool, err = rejectedStoryRewardPool(ctx, bankKeeper, trueVotes, falseVotes)
	}
	if err != nil {
		return
	}

	return
}

func distributeRewards(
	ctx sdk.Context,
	bankKeeper bank.Keeper,
	rewardPool sdk.Coin,
	trueVotes []interface{},
	falseVotes []interface{},
	confirmed bool) (err sdk.Error) {

	if confirmed {
		err = distributeConfirmedStoryRewards(
			ctx, bankKeeper, trueVotes, falseVotes, rewardPool)
	} else {
		err = distributeRejectedStoryRewards(
			ctx, bankKeeper, falseVotes, rewardPool)
	}
	if err != nil {
		return
	}

	return
}

// tally backings, challenges, and token votes into two true and false vote arrays
func tally(
	ctx sdk.Context,
	k Keeper,
	game game.Game) (trueVotes []interface{}, falseVotes []interface{}, err sdk.Error) {

	// tally backings
	trueBackings, falseBackings, err := k.backingKeeper.Tally(ctx, game.StoryID)
	if err != nil {
		return
	}
	trueVotes = append(trueVotes, trueBackings)
	falseVotes = append(falseVotes, falseBackings)

	// tally challenges
	trueChallenges, falseChallenges, err := k.challengeKeeper.Tally(ctx, game.ID)
	if err != nil {
		return
	}
	trueVotes = append(trueVotes, trueChallenges)
	falseVotes = append(falseVotes, falseChallenges)

	// tally token votes
	trueTokenVotes, falseTokenVotes, err := k.Tally(ctx, game.ID)
	if err != nil {
		return
	}
	trueVotes = append(trueVotes, trueTokenVotes)
	falseVotes = append(falseVotes, falseTokenVotes)

	return trueVotes, falseVotes, nil
}

// determine if a story is confirmed or rejected
func confirmStory(trueVotes []interface{}, falseVotes []interface{}) (confirmed bool) {
	// calculate weighted votes
	trueWeight := weightedVote(trueVotes)
	falseWeight := weightedVote(falseVotes)

	// majority wins
	if trueWeight.GT(falseWeight) {
		// story confirmed
		return true
	}

	// story rejected
	return false
}

// calculate weighted vote based on user's total category coin balance
func weightedVote(votes []interface{}) sdk.Int {
	weightedAmount := sdk.ZeroInt()
	for _, vote := range votes {
		v := vote.(app.Vote)
		user := auth.NewBaseAccountWithAddress(v.Creator)
		categoryCoins := user.Coins.AmountOf(v.Amount.Denom)
		weightedAmount = weightedAmount.Add(categoryCoins)
	}

	return weightedAmount
}