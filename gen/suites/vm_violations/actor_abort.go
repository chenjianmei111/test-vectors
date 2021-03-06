package main

import (
	"github.com/chenjianmei111/go-address"
	"github.com/chenjianmei111/go-state-types/abi"
	"github.com/chenjianmei111/go-state-types/big"
	"github.com/chenjianmei111/go-state-types/exitcode"
	"github.com/chenjianmei111/lotus/conformance/chaos"

	. "github.com/chenjianmei111/test-vectors/gen/builders"
)

func actorAbort(abortCode exitcode.ExitCode, msg string, expectedCode exitcode.ExitCode) func(*MessageVectorBuilder) {
	return func(v *MessageVectorBuilder) {
		v.Messages.SetDefaults(GasLimit(1_000_000_000), GasPremium(1), GasFeeCap(200))

		sender := v.Actors.Account(address.SECP256K1, abi.NewTokenAmount(1_000_000_000_000))
		v.CommitPreconditions()

		v.Messages.Raw(
			sender.ID,
			chaos.Address,
			chaos.MethodAbortWith,
			MustSerialize(&chaos.AbortWithArgs{Code: abortCode, Message: msg}),
			Value(big.Zero()),
			Nonce(0),
		)
		v.CommitApplies()

		v.Assert.LastMessageResultSatisfies(ExitCode(expectedCode))
	}
}

func actorPanic(msg string) func(*MessageVectorBuilder) {
	return func(v *MessageVectorBuilder) {
		v.Messages.SetDefaults(GasLimit(1_000_000_000), GasPremium(1), GasFeeCap(200))

		sender := v.Actors.Account(address.SECP256K1, abi.NewTokenAmount(1_000_000_000_000))
		v.CommitPreconditions()

		v.Messages.Raw(
			sender.ID,
			chaos.Address,
			chaos.MethodAbortWith,
			MustSerialize(&chaos.AbortWithArgs{Uncontrolled: true, Message: msg}),
			Value(big.Zero()),
			Nonce(0),
		)
		v.CommitApplies()

		v.Assert.LastMessageResultSatisfies(Failed())
	}
}
