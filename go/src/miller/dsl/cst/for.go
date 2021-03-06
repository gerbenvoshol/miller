package cst

import (
	"errors"

	"miller/dsl"
	"miller/lib"
	"miller/types"
)

// ================================================================
// This is for various flavors of for-loop
// ================================================================

// ----------------------------------------------------------------
// Sample AST:

// mlr -n put -v 'for (k in $*) { emit { k : k } }'
// DSL EXPRESSION:
// for (k in $*) { emit { k : k} }
// RAW AST:
// * StatementBlock
//     * ForLoopKeyOnly "for"
//         * LocalVariable "k"
//         * FullSrec "$*"
//         * StatementBlock
//             * EmitStatement "emit"
//                 * MapLiteral "{}"
//                     * MapLiteralKeyValuePair ":"
//                         * LocalVariable "k"
//                         * LocalVariable "k"

func (this *RootNode) BuildForLoopKeyOnlyNode(astNode *dsl.ASTNode) (*ForLoopKeyValueNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeForLoopKeyOnly)
	lib.InternalCodingErrorIf(len(astNode.Children) != 3)

	keyVariableASTNode := astNode.Children[0]
	indexableASTNode := astNode.Children[1]
	blockASTNode := astNode.Children[2]

	lib.InternalCodingErrorIf(keyVariableASTNode.Type != dsl.NodeTypeLocalVariable)
	lib.InternalCodingErrorIf(keyVariableASTNode.Token == nil)
	keyVariableName := string(keyVariableASTNode.Token.Lit)

	// TODO: error if loop-over node isn't Mappable (inasmuch as can be
	// detected at CST-build time)
	// TODO: support arrays too.
	indexableNode, err := this.BuildEvaluableNode(indexableASTNode)
	if err != nil {
		return nil, err
	}

	statementBlockNode, err := this.BuildStatementBlockNode(blockASTNode)
	if err != nil {
		return nil, err
	}

	return NewForLoopKeyValueNode(
		keyVariableName,
		"", // No value-variable here, but re-use the code
		indexableNode,
		statementBlockNode,
	), nil
}

// ================================================================
type ForLoopKeyValueNode struct {
	keyVariableName    string
	valueVariableName  string
	indexableNode      IEvaluable
	statementBlockNode *StatementBlockNode
}

func NewForLoopKeyValueNode(
	keyVariableName string,
	valueVariableName string,
	indexableNode IEvaluable,
	statementBlockNode *StatementBlockNode,
) *ForLoopKeyValueNode {
	return &ForLoopKeyValueNode{
		keyVariableName,
		valueVariableName,
		indexableNode,
		statementBlockNode,
	}
}

// ----------------------------------------------------------------
// Sample AST:

// mlr -n put -v 'for (k, v in $*) { emit { k : v } }'
// DSL EXPRESSION:
// for (k, v in $*) { emit { k : v } }
// RAW AST:
// * StatementBlock
//     * ForLoopKeyValue "for"
//         * LocalVariable "k"
//         * LocalVariable "v"
//         * FullSrec "$*"
//         * StatementBlock
//             * EmitStatement "emit"
//                 * MapLiteral "{}"
//                     * MapLiteralKeyValuePair ":"
//                         * LocalVariable "k"
//                         * LocalVariable "v"

func (this *RootNode) BuildForLoopKeyValueNode(astNode *dsl.ASTNode) (*ForLoopKeyValueNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeForLoopKeyValue)
	lib.InternalCodingErrorIf(len(astNode.Children) != 4)

	keyVariableASTNode := astNode.Children[0]
	valueVariableASTNode := astNode.Children[1]
	indexableASTNode := astNode.Children[2]
	blockASTNode := astNode.Children[3]

	lib.InternalCodingErrorIf(keyVariableASTNode.Type != dsl.NodeTypeLocalVariable)
	lib.InternalCodingErrorIf(keyVariableASTNode.Token == nil)
	keyVariableName := string(keyVariableASTNode.Token.Lit)

	lib.InternalCodingErrorIf(valueVariableASTNode.Type != dsl.NodeTypeLocalVariable)
	lib.InternalCodingErrorIf(valueVariableASTNode.Token == nil)
	valueVariableName := string(valueVariableASTNode.Token.Lit)

	// TODO: error if loop-over node isn't Mappable (inasmuch as can be
	// detected at CST-build time)
	indexableNode, err := this.BuildEvaluableNode(indexableASTNode)
	if err != nil {
		return nil, err
	}

	statementBlockNode, err := this.BuildStatementBlockNode(blockASTNode)
	if err != nil {
		return nil, err
	}

	return NewForLoopKeyValueNode(
		keyVariableName,
		valueVariableName,
		indexableNode,
		statementBlockNode,
	), nil
}

// ----------------------------------------------------------------
// Note: The statement-block has its own push/pop for its localvars.
// Meanwhile we need to restrict scope of the bindvar to the for-loop.
//
// So we have:
//
//   mlr put '
//     x = 1;           <--- frame #1 main
//     for (k in $*) {  <--- frame #2 for for-loop bindvars (right here)
//       x = 2          <--- frame #3 for for-loop locals
//     }
//     x = 3;           <--- back in frame #1 main
//   '
//

func (this *ForLoopKeyValueNode) Execute(state *State) (*BlockExitPayload, error) {
	mlrval := this.indexableNode.Evaluate(state)

	if mlrval.IsMap() {

		mapval := mlrval.GetMap()

		// Make a frame for the loop variable(s)
		state.stack.PushStackFrame()
		defer state.stack.PopStackFrame()
		for pe := mapval.Head; pe != nil; pe = pe.Next {
			mapkey := types.MlrvalFromString(*pe.Key)

			state.stack.BindVariable(this.keyVariableName, &mapkey)
			if this.valueVariableName != "" { // 'for (k in ...)' not 'for (k,v in ...)'
				state.stack.BindVariable(this.valueVariableName, pe.Value)
			}
			// The loop body will push its own frame
			blockExitPayload, err := this.statementBlockNode.Execute(state)
			if err != nil {
				return nil, err
			}
			if blockExitPayload != nil {
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_BREAK {
					break
				}
				// If continue, keep going -- this means the body was exited
				// early but we keep going at this level
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VOID {
					return blockExitPayload, nil
				}
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VALUE {
					return blockExitPayload, nil
				}
			}
			// TODO: handle return statements
			// TODO: runtime errors for any other types
		}

	} else if mlrval.IsArray() {

		arrayval := mlrval.GetArray()

		// Note: Miller user-space array indices ("mindex") are 1-up. Internal
		// Go storage ("zindex") is 0-up.

		// Make a frame for the loop variable(s)
		state.stack.PushStackFrame()
		defer state.stack.PopStackFrame()
		for zindex, element := range arrayval {
			mindex := types.MlrvalFromInt64(int64(zindex + 1))

			state.stack.BindVariable(this.keyVariableName, &mindex)
			if this.valueVariableName != "" { // 'for (k in ...)' not 'for (k,v in ...)'
				state.stack.BindVariable(this.valueVariableName, &element)
			}
			// The loop body will push its own frame
			blockExitPayload, err := this.statementBlockNode.Execute(state)
			if err != nil {
				return nil, err
			}
			if blockExitPayload != nil {
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_BREAK {
					break
				}
				// If continue, keep going -- this means the body was exited
				// early but we keep going at this level
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VOID {
					return blockExitPayload, nil
				}
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VALUE {
					return blockExitPayload, nil
				}
			}
			// TODO: handle return statements
			// TODO: runtime errors for any other types
		}

	} else {
		// TODO: give more context
		return nil, errors.New("Miller: looped-over item is not a map.")
	}

	return nil, nil
}

// ================================================================
type TripleForLoopNode struct {
	startBlockNode             *StatementBlockNode
	continuationExpressionNode IEvaluable
	updateBlockNode            *StatementBlockNode
	bodyBlockNode              *StatementBlockNode
}

func NewTripleForLoopNode(
	startBlockNode *StatementBlockNode,
	continuationExpressionNode IEvaluable,
	updateBlockNode *StatementBlockNode,
	bodyBlockNode *StatementBlockNode,
) *TripleForLoopNode {
	return &TripleForLoopNode{
		startBlockNode,
		continuationExpressionNode,
		updateBlockNode,
		bodyBlockNode,
	}
}

// ----------------------------------------------------------------
// Sample ASTs:

// DSL EXPRESSION:
// for (;;) {}
// RAW AST:
// * StatementBlock
//     * TripleForLoop "for"
//         * StatementBlock
//         * StatementBlock
//         * StatementBlock
//         * StatementBlock

// mlr --from u/s.dkvp put -v for (i = 0; i < NR; i += 1) { $i += i }
// DSL EXPRESSION:
// for (i = 0; i < NR; i += 1) { $i += i }
// RAW AST:
// * StatementBlock
//     * TripleForLoop "for"
//         * StatementBlock
//             * Assignment "="
//                 * LocalVariable "i"
//                 * IntLiteral "0"
//         * StatementBlock
//             * BareBoolean
//                 * Operator "<"
//                     * LocalVariable "i"
//                     * ContextVariable "NR"
//         * StatementBlock
//             * Assignment "="
//                 * LocalVariable "i"
//                 * Operator "+"
//                     * LocalVariable "i"
//                     * IntLiteral "1"
//         * StatementBlock
//             * Assignment "="
//                 * DirectFieldValue "i"
//                 * Operator "+"
//                     * DirectFieldValue "i"
//                     * LocalVariable "i"

func (this *RootNode) BuildTripleForLoopNode(astNode *dsl.ASTNode) (*TripleForLoopNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeTripleForLoop)
	lib.InternalCodingErrorIf(len(astNode.Children) != 4)

	startBlockASTNode := astNode.Children[0]
	continuationExpressionASTNode := astNode.Children[1]
	updateBlockASTNode := astNode.Children[2]
	bodyBlockASTNode := astNode.Children[3]

	lib.InternalCodingErrorIf(startBlockASTNode.Type != dsl.NodeTypeStatementBlock)
	lib.InternalCodingErrorIf(continuationExpressionASTNode.Type != dsl.NodeTypeStatementBlock)
	lib.InternalCodingErrorIf(len(continuationExpressionASTNode.Children) > 1)
	lib.InternalCodingErrorIf(updateBlockASTNode.Type != dsl.NodeTypeStatementBlock)
	lib.InternalCodingErrorIf(bodyBlockASTNode.Type != dsl.NodeTypeStatementBlock)

	startBlockNode, err := this.BuildStatementBlockNode(startBlockASTNode)
	if err != nil {
		return nil, err
	}

	var continuationExpressionNode IEvaluable = nil
	// empty is true
	if len(continuationExpressionASTNode.Children) == 1 {
		bareBooleanASTNode := continuationExpressionASTNode.Children[0]
		lib.InternalCodingErrorIf(bareBooleanASTNode.Type != dsl.NodeTypeBareBoolean)
		lib.InternalCodingErrorIf(len(bareBooleanASTNode.Children) != 1)

		continuationExpressionNode, err = this.BuildEvaluableNode(bareBooleanASTNode.Children[0])
		if err != nil {
			return nil, err
		}
	}

	updateBlockNode, err := this.BuildStatementBlockNode(updateBlockASTNode)
	if err != nil {
		return nil, err
	}

	bodyBlockNode, err := this.BuildStatementBlockNode(bodyBlockASTNode)
	if err != nil {
		return nil, err
	}

	return NewTripleForLoopNode(
		startBlockNode,
		continuationExpressionNode,
		updateBlockNode,
		bodyBlockNode,
	), nil
}

// ----------------------------------------------------------------
// Note: The statement-block has its own push/pop for its localvars.
// Meanwhile we need to restrict scope of the bindvar to the for-loop.
//
// So we have:
//
//   mlr put '
//     x = 1;           <--- frame #1 main
//     for (k in $*) {  <--- frame #2 for for-loop bindvars (right here)
//       x = 2          <--- frame #3 for for-loop locals
//     }
//     x = 3;           <--- back in frame #1 main
//   '
//

func (this *TripleForLoopNode) Execute(state *State) (*BlockExitPayload, error) {
	// Make a frame for the loop variables.
	state.stack.PushStackFrame()
	defer state.stack.PopStackFrame()

	// Use ExecuteFrameless here, otherwise the start-statements would be
	// within an ephemeral, isolated frame and not accessible to the remaining
	// parts of the for-loop.
	_, err := this.startBlockNode.ExecuteFrameless(state)
	if err != nil {
		return nil, err
	}

	for {
		// state.stack.Dump()
		// empty is true
		if this.continuationExpressionNode != nil {
			continuationValue := this.continuationExpressionNode.Evaluate(state)
			boolValue, isBool := continuationValue.GetBoolValue()
			if !isBool {
				// TODO: propagate line-number context
				return nil, errors.New("Miller: for-loop continuation did not evaluate to boolean.")
			}
			if boolValue == false {
				break
			}
		}

		blockExitPayload, err := this.bodyBlockNode.Execute(state)
		if err != nil {
			return nil, err
		}
		if blockExitPayload != nil {
			if blockExitPayload.blockExitStatus == BLOCK_EXIT_BREAK {
				break
			}
			// If continue, keep going -- this means the body was exited
			// early but we keep going at this level. In particular we still
			// need to execute the update-block.
			if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VOID {
				return blockExitPayload, nil
			}
			if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VALUE {
				return blockExitPayload, nil
			}
		}
		// TODO: handle return statements
		// TODO: runtime errors for any other types

		// The loop body will push its own frame.
		state.stack.PushStackFrame()
		_, err = this.updateBlockNode.ExecuteFrameless(state)
		if err != nil {
			state.stack.PopStackFrame()
			return nil, err
		}
		state.stack.PopStackFrame()
	}

	return nil, nil
}
