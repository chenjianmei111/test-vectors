package main

import (
	"github.com/chenjianmei111/go-state-types/exitcode"

	"github.com/chenjianmei111/lotus/chain/actors/builtin/paych"
	"github.com/chenjianmei111/lotus/chain/types"

	. "github.com/chenjianmei111/test-vectors/gen/builders"

	"github.com/chenjianmei111/go-state-types/big"
	init_ "github.com/chenjianmei111/specs-actors/actors/builtin/init"

	"github.com/chenjianmei111/go-address"
)

func sequentialAddresses(v *MessageVectorBuilder) {
	v.Messages.SetDefaults(GasLimit(1_000_000_000), GasPremium(1), GasFeeCap(200))

	initial := big.NewInt(1_000_000_000_000_000)

	// Set up sender and receiver accounts.
	var sender, receiver AddressHandle
	v.Actors.AccountN(address.SECP256K1, initial, &sender, &receiver)
	v.CommitPreconditions()

	// Create 10 payment channels.
	for i := uint64(0); i < 10; i++ {
		v.Messages.Sugar().PaychMessage(sender.Robust, func(b paych.MessageBuilder) (*types.Message, error) {
			return b.Create(receiver.Robust, big.NewInt(1000))
		}, Value(big.NewInt(1000)), Nonce(i))
	}
	v.CommitApplies()

	v.Assert.EveryMessageResultSatisfies(ExitCode(exitcode.Ok))

	for i, am := range v.Messages.All() {
		expectedActorAddr := AddressHandle{
			ID:     MustNewIDAddr(MustIDFromAddress(receiver.ID) + uint64(i) + 1),
			Robust: sender.NextActorAddress(am.Message.Nonce, 0),
		}

		// Verify that the return contains the expected addresses.
		var ret init_.ExecReturn
		MustDeserialize(am.Result.Return, &ret)
		v.Assert.Equal(expectedActorAddr.Robust, ret.RobustAddress)
		v.Assert.Equal(expectedActorAddr.ID, ret.IDAddress)
	}

	v.Assert.EveryMessageSenderSatisfies(BalanceUpdated(big.Zero()))
	v.Assert.EveryMessageSenderSatisfies(NonceUpdated())

}
