package builders

import (
	"github.com/chenjianmei111/lotus/conformance/chaos"

	"github.com/chenjianmei111/go-state-types/abi"
	"github.com/chenjianmei111/go-state-types/big"

	"github.com/chenjianmei111/specs-actors/actors/builtin"
	"github.com/chenjianmei111/specs-actors/actors/builtin/cron"
	init_ "github.com/chenjianmei111/specs-actors/actors/builtin/init"
	"github.com/chenjianmei111/specs-actors/actors/builtin/market"
	"github.com/chenjianmei111/specs-actors/actors/builtin/miner"
	"github.com/chenjianmei111/specs-actors/actors/builtin/multisig"
	"github.com/chenjianmei111/specs-actors/actors/builtin/power"
	"github.com/chenjianmei111/specs-actors/actors/builtin/reward"
	"github.com/chenjianmei111/specs-actors/actors/runtime/proof"

	"github.com/chenjianmei111/go-address"
)

func Transfer() TypedCall {
	return func(_ *Messages) (method abi.MethodNum, params []byte) {
		return builtin.MethodSend, []byte{}
	}
}

// ----------------------------------------------------------------------------
// | ACCOUNT
// ----------------------------------------------------------------------------

func AccountConstructor(params *address.Address) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsAccount.Constructor, MustSerialize(params)
	}
}

func AccountPubkeyAddress(params *abi.EmptyValue) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsAccount.PubkeyAddress, MustSerialize(params)
	}
}

// ----------------------------------------------------------------------------
// | MARKET
// ----------------------------------------------------------------------------

func MarketConstructor(params *abi.EmptyValue) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMarket.Constructor, MustSerialize(params)
	}
}

func MarketAddBalance(params *address.Address) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMarket.AddBalance, MustSerialize(params)
	}
}
func MarketWithdrawBalance(params *market.WithdrawBalanceParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMarket.WithdrawBalance, MustSerialize(params)
	}
}
func MarketPublishStorageDeals(params *market.PublishStorageDealsParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMarket.PublishStorageDeals, MustSerialize(params)
	}
}
func MarketVerifyDealsForActivation(params *market.VerifyDealsForActivationParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMarket.VerifyDealsForActivation, MustSerialize(params)
	}
}
func MarketActivateDeals(params *market.ActivateDealsParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMarket.ActivateDeals, MustSerialize(params)
	}
}
func MarketOnMinerSectorsTerminate(params *market.OnMinerSectorsTerminateParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMarket.OnMinerSectorsTerminate, MustSerialize(params)
	}
}
func MarketComputeDataCommitment(params *market.ComputeDataCommitmentParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMarket.ComputeDataCommitment, MustSerialize(params)
	}
}
func MarketCronTick(params *abi.EmptyValue) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMarket.CronTick, MustSerialize(params)
	}
}

// ----------------------------------------------------------------------------
// | MINER
// ----------------------------------------------------------------------------

func MinerConstructor(params *power.MinerConstructorParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMiner.Constructor, MustSerialize(params)
	}
}
func MinerControlAddresses(params *abi.EmptyValue) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMiner.ControlAddresses, MustSerialize(params)
	}
}
func MinerChangeWorkerAddress(params *miner.ChangeWorkerAddressParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMiner.ChangeWorkerAddress, MustSerialize(params)
	}
}
func MinerChangePeerID(params *miner.ChangePeerIDParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMiner.ChangePeerID, MustSerialize(params)
	}
}
func MinerSubmitWindowedPoSt(params *miner.SubmitWindowedPoStParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMiner.SubmitWindowedPoSt, MustSerialize(params)
	}
}
func MinerPreCommitSector(params *miner.SectorPreCommitInfo) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMiner.PreCommitSector, MustSerialize(params)
	}
}
func MinerProveCommitSector(params *miner.ProveCommitSectorParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMiner.ProveCommitSector, MustSerialize(params)
	}
}
func MinerExtendSectorExpiration(params *miner.ExtendSectorExpirationParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMiner.ExtendSectorExpiration, MustSerialize(params)
	}
}
func MinerTerminateSectors(params *miner.TerminateSectorsParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMiner.TerminateSectors, MustSerialize(params)
	}
}
func MinerDeclareFaults(params *miner.DeclareFaultsParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMiner.DeclareFaults, MustSerialize(params)
	}
}
func MinerDeclareFaultsRecovered(params *miner.DeclareFaultsRecoveredParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMiner.DeclareFaultsRecovered, MustSerialize(params)
	}
}
func MinerOnDeferredCronEvent(params *miner.CronEventPayload) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMiner.OnDeferredCronEvent, MustSerialize(params)
	}
}
func MinerCheckSectorProven(params *miner.CheckSectorProvenParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMiner.CheckSectorProven, MustSerialize(params)
	}
}
func MinerAddLockedFund(params *big.Int) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMiner.AddLockedFund, MustSerialize(params)
	}
}
func MinerReportConsensusFault(params *miner.ReportConsensusFaultParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMiner.ReportConsensusFault, MustSerialize(params)
	}
}
func MinerWithdrawBalance(params *miner.WithdrawBalanceParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMiner.WithdrawBalance, MustSerialize(params)
	}
}
func MinerConfirmSectorProofsValid(params *builtin.ConfirmSectorProofsParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMiner.ConfirmSectorProofsValid, MustSerialize(params)
	}
}
func MinerChangeMultiaddrs(params *miner.ChangeMultiaddrsParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMiner.ChangeMultiaddrs, MustSerialize(params)
	}
}

// ----------------------------------------------------------------------------
// | MULTISIG -- use sugared methods for other operations.
// ----------------------------------------------------------------------------

func MultisigAddSigner(params *multisig.AddSignerParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMultisig.AddSigner, MustSerialize(params)
	}
}
func MultisigRemoveSigner(params *multisig.RemoveSignerParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMultisig.RemoveSigner, MustSerialize(params)
	}
}
func MultisigSwapSigner(params *multisig.SwapSignerParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMultisig.SwapSigner, MustSerialize(params)
	}
}
func MultisigChangeNumApprovalsThreshold(params *multisig.ChangeNumApprovalsThresholdParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsMultisig.ChangeNumApprovalsThreshold, MustSerialize(params)
	}
}

// ----------------------------------------------------------------------------
// | POWER
// ----------------------------------------------------------------------------

func PowerConstructor(params *abi.EmptyValue) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsPower.Constructor, MustSerialize(params)
	}
}
func PowerCreateMiner(params *power.CreateMinerParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsPower.CreateMiner, MustSerialize(params)
	}
}
func PowerUpdateClaimedPower(params *power.UpdateClaimedPowerParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsPower.UpdateClaimedPower, MustSerialize(params)
	}
}
func PowerEnrollCronEvent(params *power.EnrollCronEventParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsPower.EnrollCronEvent, MustSerialize(params)
	}
}
func PowerOnEpochTickEnd(params *abi.EmptyValue) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsPower.OnEpochTickEnd, MustSerialize(params)
	}
}
func PowerUpdatePledgeTotal(params *big.Int) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsPower.UpdatePledgeTotal, MustSerialize(params)
	}
}
func PowerOnConsensusFault(params *big.Int) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsPower.OnConsensusFault, MustSerialize(params)
	}
}
func PowerSubmitPoRepForBulkVerify(params *proof.SealVerifyInfo) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsPower.SubmitPoRepForBulkVerify, MustSerialize(params)
	}
}
func PowerCurrentTotalPower(params *abi.EmptyValue) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsPower.CurrentTotalPower, MustSerialize(params)
	}
}

// ----------------------------------------------------------------------------
// | REWARD
// ----------------------------------------------------------------------------

func RewardConstructor(params *abi.EmptyValue) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsReward.Constructor, MustSerialize(params)
	}
}
func RewardAwardBlockReward(params *reward.AwardBlockRewardParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsReward.AwardBlockReward, MustSerialize(params)
	}
}
func RewardThisEpochReward(params *abi.EmptyValue) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsReward.ThisEpochReward, MustSerialize(params)
	}
}
func RewardUpdateNetworkKPI(params *big.Int) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsReward.UpdateNetworkKPI, MustSerialize(params)
	}
}

// ----------------------------------------------------------------------------
// | PAYCH -- use sugared methods.
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// | CRON
// ----------------------------------------------------------------------------

func CronConstructor(params *cron.ConstructorParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsCron.Constructor, MustSerialize(params)
	}
}
func CronEpochTick(params *abi.EmptyValue) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsCron.EpochTick, MustSerialize(params)
	}
}

// ----------------------------------------------------------------------------
// | INIT
// ----------------------------------------------------------------------------

func InitConstructor(params *init_.ConstructorParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsInit.Constructor, MustSerialize(params)
	}
}
func InitExec(params *init_.ExecParams) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return builtin.MethodsInit.Exec, MustSerialize(params)
	}
}

// ----------------------------------------------------------------------------
// | CHAOS
// ----------------------------------------------------------------------------

// ChaosSend is a message to a chaos actor that sends a message to another actor.
func ChaosSend(args *chaos.SendArgs) TypedCall {
	return func(_ *Messages) (abi.MethodNum, []byte) {
		return chaos.MethodSend, MustSerialize(args)
	}
}
