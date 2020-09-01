package cst

import (
	"miller/containers"
	"miller/lib"
)

// ----------------------------------------------------------------
// There are three CST roots: begin-block, body-block, and end-block.
//
// Next-level items are:
// * srec assignments
// * oosvar assignments
// * localvar assignments
// * emit et al.
// * bare-boolean
// * break/continue/return
// * statement block (if-body, for-body, etc)
// ----------------------------------------------------------------

// ----------------------------------------------------------------
// AST nodes (TNodeType) at the moment:
//
// NodeTypeStringLiteral
// NodeTypeIntLiteral
// NodeTypeFloatLiteral
// NodeTypeBoolLiteral
//
// NodeTypeDirectFieldName
// NodeTypeIndirectFieldName
//
// NodeTypeStatementBlock
// NodeTypeAssignment
// NodeTypeOperator
// NodeTypeContextVariable
// ----------------------------------------------------------------

// ----------------------------------------------------------------
type State struct {
	Inrec   *containers.Lrec
	Context *containers.Context
	// oosvars too
	// stack frames will go into individual statement-block nodes
}

func NewState(
	inrec *containers.Lrec,
	context *containers.Context,
) *State {
	return &State{
		Inrec:   inrec,
		Context: context,
	}
}

// ----------------------------------------------------------------
type IExecutable interface {
	Execute(state *State)
}

// ----------------------------------------------------------------
type Root struct {
	// array of statements/blocks
	executables []IExecutable
}

func NewRoot() *Root {
	return &Root{
		make([]IExecutable, 0),
	}
}
func (this *Root) AppendStatement(executable IExecutable) {
	this.executables = append(this.executables, executable)
}

type DirectSrecFieldAssignment struct {
	lhsFieldName string
	rhs          IEvaluable
}

func NewDirectSrecFieldAssignment(
	lhsFieldName string,
	rhs IEvaluable,
) *DirectSrecFieldAssignment {
	return &DirectSrecFieldAssignment{
		lhsFieldName: lhsFieldName,
		rhs:          rhs,
	}
}

// xxx separate dsl/cst/execute.go ?
// xxx implement IExecutable

type IndirectSrecFieldAssignment struct {
	lhsFieldName IEvaluable
	rhs          IEvaluable
}

// xxx implement IExecutable

type StatementBlock struct {
	// list of statement
}

// xxx implement IExecutable

// ================================================================
type IEvaluable interface {
	Evaluate(state *State) lib.Mlrval
	// Needs an Evaluate which takes context and produces a mlrval
}
