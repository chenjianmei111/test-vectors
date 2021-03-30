package main

import (
	"github.com/chenjianmei111/go-address"
	"github.com/chenjianmei111/go-state-types/abi"
	"github.com/chenjianmei111/go-state-types/big"
	"github.com/chenjianmei111/go-state-types/exitcode"
	"github.com/chenjianmei111/lotus/conformance/chaos"

	. "github.com/chenjianmei111/test-vectors/gen/builders"
)

func callerValidation(args chaos.CallerValidationArgs, expectedCode exitcode.ExitCode) func(v *MessageVectorBuilder) {
	return func(v *MessageVectorBuilder) {
		v.Messages.SetDefaults(GasLimit(1_000_000_000), GasPremium(1), GasFeeCap(200))

		alice := v.Actors.Account(address.SECP256K1, abi.NewTokenAmount(1_000_000_000_000))
		v.CommitPreconditions()

		v.Messages.Raw(alice.ID, chaos.Address, chaos.MethodCallerValidation, MustSerialize(&args), Nonce(0), Value(big.Zero()))
		v.CommitApplies()

		v.Assert.EveryMessageResultSatisfies(ExitCode(expectedCode))
	}
}
