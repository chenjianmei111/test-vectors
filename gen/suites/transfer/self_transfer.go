package main

import (
	"github.com/chenjianmei111/go-address"
	"github.com/chenjianmei111/go-state-types/abi"
	"github.com/chenjianmei111/go-state-types/exitcode"

	. "github.com/chenjianmei111/test-vectors/gen/builders"
)

func selfTransfer(from, to func(h AddressHandle) address.Address) func(v *MessageVectorBuilder) {
	return func(v *MessageVectorBuilder) {
		initial := abi.NewTokenAmount(1_000_000_000_000)
		transfer := abi.NewTokenAmount(10)
		v.Messages.SetDefaults(GasLimit(1_000_000_000), GasPremium(1), GasFeeCap(200))

		// Set up sender account.
		account := v.Actors.Account(address.SECP256K1, initial)
		v.CommitPreconditions()

		// Perform the transfer.
		msg := v.Messages.Sugar().Transfer(from(account), to(account), Value(transfer), Nonce(0))
		v.CommitApplies()

		v.Assert.Equal(exitcode.Ok, msg.Result.ExitCode)

		// the transfer balance comes back to us.
		v.Assert.EveryMessageSenderSatisfies(BalanceUpdated(transfer))
	}
}
