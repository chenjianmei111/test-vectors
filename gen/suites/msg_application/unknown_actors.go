package main

import (
	"github.com/chenjianmei111/go-address"
	"github.com/chenjianmei111/go-state-types/exitcode"

	. "github.com/chenjianmei111/test-vectors/gen/builders"
)

func failUnknownSender(v *MessageVectorBuilder) {
	v.Messages.SetDefaults(GasLimit(1_000_000_000), GasPremium(1), GasFeeCap(200))

	alice := v.Actors.Account(address.SECP256K1, balance1T)
	v.CommitPreconditions()

	v.Messages.Sugar().Transfer(unknown, alice.ID, Value(transferAmnt), Nonce(0))
	v.CommitApplies()

	v.Assert.EveryMessageResultSatisfies(ExitCode(exitcode.SysErrSenderInvalid))
}

func failUnknownReceiver(v *MessageVectorBuilder) {
	v.Messages.SetDefaults(GasLimit(1_000_000_000), GasPremium(1), GasFeeCap(200))

	alice := v.Actors.Account(address.SECP256K1, balance1T)
	v.CommitPreconditions()

	// Sending a message to non-existent ID address must produce an error.
	unknownID := MustNewIDAddr(10000000)
	v.Messages.Sugar().Transfer(alice.ID, unknownID, Value(transferAmnt), Nonce(0))

	unknownActor := MustNewActorAddr("1234")
	v.Messages.Sugar().Transfer(alice.ID, unknownActor, Value(transferAmnt), Nonce(1))
	v.CommitApplies()

	v.Assert.EveryMessageResultSatisfies(ExitCode(exitcode.SysErrInvalidReceiver))
}
