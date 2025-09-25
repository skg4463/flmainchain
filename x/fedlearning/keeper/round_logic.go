package keeper

import (
	"fmt"
	"sort"
	"strconv"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"flmainchain/x/fedlearning/types"
)

// AdvanceRoundState. BeginBlocker에서 호출됨.
// 모든 참여자의 제출이 완료되면 라운드 단계를 자동으로 진행.
func (k Keeper) AdvanceRoundState(ctx sdk.Context) {
	currentRound, err := k.CurrentRound.Get(ctx)
	if err != nil {
		// 현재 라운드 정보가 없으면 아무것도 하지 않음.
		return
	}
	round, err := k.Round.Get(ctx, currentRound.RoundId)
	if err != nil {
		return
	}

	// L-node 제출이 모두 완료되었는지 확인.
	if round.Status == "WeightSubmissionOpen" && len(round.SubmittedLNodes) >= len(round.RequiredLNodes) {
		round.Status = "ScoreSubmissionOpen"
		k.Round.Set(ctx, round.RoundId, round)
		ctx.Logger().Info("Round advanced to ScoreSubmissionOpen", "round", currentRound.RoundId)
	}

	// C-node 제출이 모두 완료되었는지 확인.
	if round.Status == "ScoreSubmissionOpen" && len(round.SubmittedCNodes) >= len(round.RequiredCNodes) {
		round.Status = "AggregationReady"
		k.Round.Set(ctx, round.RoundId, round)
		ctx.Logger().Info("Round advanced to AggregationReady", "round", currentRound.RoundId)
	}
}

// AggregateScoresAndCreateATT. 3N-1 블록의 EndBlocker에서 호출됨.
func (k Keeper) AggregateScoresAndCreateATT(ctx sdk.Context) {
	currentRound, err := k.CurrentRound.Get(ctx)
	if err != nil { return }
	round, err := k.Round.Get(ctx, currentRound.RoundId)
	if err != nil || round.Status != "AggregationReady" { return }

	scoreMap := make(map[string][]uint64)
	
	// --- 이 부분만 수정 ---
	// k.storeService에서 먼저 KVStore를 열어줍니다.
	kvStore := k.storeService.OpenKVStore(ctx)
	// 열린 KVStore를 어댑터에 전달합니다.
	storeAdapter := runtime.KVStoreAdapter(kvStore)
	// --- 수정 끝 ---
	
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
	
	// ... (이하 평균 계산 및 저장 로직은 이전과 동일) ...
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
	ctx.Logger().Info("ATT for round aggregated and saved.", "round", round.RoundId)
}


// ElectNextCommittee. 3N 블록의 EndBlocker에서 호출됨.
// 최종 점수표를 바탕으로 다음 라운드 위원회를 선출.
func (k Keeper) ElectNextCommittee(ctx sdk.Context) {
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
	k.RoundCommittee.Set(ctx, nextRoundID, types.RoundCommittee{Members: nextCommitteeMembers})
	k.Round.Set(ctx, nextRoundID, types.Round{
		Status:          "WeightSubmissionOpen",
		RequiredLNodes:  nextCommitteeMembers, 
		SubmittedLNodes: []string{},
		RequiredCNodes:  nextCommitteeMembers, 
		SubmittedCNodes: []string{},
	})
	currentRound.RoundId = nextRoundID
	k.CurrentRound.Set(ctx, currentRound)

	ctx.Logger().Info("Committee for next round elected.", "round", nextRoundID, "members", nextCommitteeMembers)
}