package keeper

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"flmainchain/x/fedlearning/types"
)

func (k Keeper) AdvanceRoundState(ctx sdk.Context) {
	currentRound, err := k.CurrentRound.Get(ctx)
	if err != nil { return }
	round, err := k.Round.Get(ctx, currentRound.RoundId)
	if err != nil { return }

	if round.Status == "WeightSubmissionOpen" && len(round.SubmittedLNodes) >= len(round.RequiredLNodes) {
		round.Status = "ScoreSubmissionOpen"
		k.Round.Set(ctx, round.RoundId, round)
		ctx.Logger().Info("Round advanced to ScoreSubmissionOpen", "round", currentRound.RoundId)
	}

	if round.Status == "ScoreSubmissionOpen" && len(round.SubmittedCNodes) >= len(round.RequiredCNodes) {
		round.Status = "AggregationReady"
		k.Round.Set(ctx, round.RoundId, round)
		ctx.Logger().Info("Round advanced to AggregationReady", "round", currentRound.RoundId)
	}
	if round.Status == "AggregationReady" {
        k.AggregateScoresAndCreateATT(ctx)
    }
}

func (k Keeper) AggregateScoresAndCreateATT(ctx sdk.Context) {
    start := time.Now()
    // 소수 4자리 타임스탬프 포맷터(로컬)
    ts4 := func(t time.Time) string {
        return t.Format("15:04:05.") + fmt.Sprintf("%04d", t.Nanosecond()/1e5)
    }

    currentRound, err := k.CurrentRound.Get(ctx)
    if err != nil { return }
    round, err := k.Round.Get(ctx, currentRound.RoundId)
    if err != nil || round.Status != "AggregationReady" { return }

    // 시작 로그 추가
    //ctx.Logger().Info(
    //    "PERF Start",
    //    "Phase", "AggregateScores-ATT",
    //    "round_id", round.RoundId,
    //    "start_time", ts4(start),
    //)

	scoreMap := make(map[string][]uint64)
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.SubmittedScoreKey)
	iteratorPrefix := []byte(fmt.Sprintf("%d-", round.RoundId))
	iterator := store.Iterator(iteratorPrefix, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var submittedScore types.SubmittedScore
		k.cdc.MustUnmarshal(iterator.Value(), &submittedScore)
		for i, lNodeAddr := range submittedScore.LnodeAddresses {
			score, err := strconv.ParseUint(submittedScore.Scores[i], 10, 64)
			if err != nil { continue }
			scoreMap[lNodeAddr] = append(scoreMap[lNodeAddr], score)
		}
	}

	type lnodeScorePair struct { Addr string; Score uint64 }
	var sortedScores []lnodeScorePair
	for lNode, scores := range scoreMap {
		var sum uint64 = 0
		if len(scores) == 0 { continue }
		for _, s := range scores { sum += s }
		avgScore := sum / uint64(len(scores))
		sortedScores = append(sortedScores, lnodeScorePair{Addr: lNode, Score: avgScore})
	}
	sort.Slice(sortedScores, func(i, j int) bool { return sortedScores[i].Score > sortedScores[j].Score })

	finalATT := types.FinalAtt{RoundId: round.RoundId}
	for _, pair := range sortedScores {
		finalATT.LnodeAddresses = append(finalATT.LnodeAddresses, pair.Addr)
		finalATT.Scores = append(finalATT.Scores, strconv.FormatUint(pair.Score, 10))
	}
	k.FinalAtt.Set(ctx, finalATT.RoundId, finalATT)

	round.Status = "AggregationComplete"
	k.Round.Set(ctx, round.RoundId, round)

    end := time.Now()
    elapsed := end.Sub(start)

    // 종료 로그(시작/종료/소요시간)
    ctx.Logger().Info(
        "PERF Done",
        "Phase", "AggregateScores-ATT",
        "Round", round.RoundId,
        "start_time", ts4(start),
        "end_time", ts4(end),
        "elapsed_time", fmt.Sprintf("%.5f", elapsed.Seconds()),
    )
}

func (k Keeper) ElectNextCommittee(ctx sdk.Context) {
    start := time.Now()
    // 소수 4자리 타임스탬프 포맷터(로컬)
    ts4 := func(t time.Time) string {
        return t.Format("15:04:05.") + fmt.Sprintf("%04d", t.Nanosecond()/1e5)
    }

    currentRound, err := k.CurrentRound.Get(ctx)
    if err != nil { return }
    round, err := k.Round.Get(ctx, currentRound.RoundId)
    if err != nil || round.Status != "AggregationComplete" { return }

    att, err := k.FinalAtt.Get(ctx, round.RoundId)
    if err != nil { return }

	var nextCommitteeMembers []string
	if len(att.LnodeAddresses) > 5 {
		nextCommitteeMembers = att.LnodeAddresses[:5]
	} else {
		nextCommitteeMembers = att.LnodeAddresses
	}
	if len(nextCommitteeMembers) == 0 { return }

	nextRoundID := round.RoundId + 1
	// Roundcommittee -> RoundCommittee 로 수정
	k.RoundCommittee.Set(ctx, nextRoundID, types.RoundCommittee{RoundId: nextRoundID, Members: nextCommitteeMembers})
	k.Round.Set(ctx, nextRoundID, types.Round{
		RoundId:         nextRoundID,
		Status:          "WeightSubmissionOpen",
		RequiredLNodes:  round.RequiredLNodes,
		SubmittedLNodes: []string{},
		RequiredCNodes:  nextCommitteeMembers,
		SubmittedCNodes: []string{},
	})
	currentRound.RoundId = nextRoundID
	k.CurrentRound.Set(ctx, currentRound)

    end := time.Now()
    elapsed := end.Sub(start)

    // 종료 로그(시작/종료/소요시간)
    ctx.Logger().Info(
        "PERF Done",
        "Phase", "ElectNextCommittee",
        "Round", nextRoundID,
        "start_time", ts4(start),
        "end_time", ts4(end),
        "elapsed_time", fmt.Sprintf("%.5f", elapsed.Seconds()),
    )
}