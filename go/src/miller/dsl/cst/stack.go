package cst

import (
	"container/list"
	"fmt"

	"miller/types"
)

// ================================================================
// Stack frames for begin/end/if/for/function blocks
// ================================================================

// ----------------------------------------------------------------
type Stack struct {
	stackFrames *list.List // list of *StackFrame
}

func NewStack() *Stack {
	return &Stack{
		stackFrames: list.New(),
	}
}

func (this *Stack) PushStackFrame() {
	this.stackFrames.PushFront(NewStackFrame())
}

func (this *Stack) PopStackFrame() {
	this.stackFrames.Remove(this.stackFrames.Front())
}

// Sets the variable at the current frame whether it's defined outer from there
// or not.
//
// OK to use BindVariable:
//
//   k = 1                 <-- top-level -frame, k=1
//   for (k in $*) { ... } <-- another k is bound in the loop
//   $k = k                <-- k is still 1
//
// Not OK to use BindVariable:
//
//   z = 1         <-- top-level frame, z=1
//   if (NR < 2) {
//     z = 2       <-- this should adjust top-level z, not bind within if-block
//   } else {
//     z = 3       <-- this should adjust top-level z, not bind within else-block
//   }
//   $z = z        <-- z should be 2 or 3, not 1
//
func (this *Stack) BindVariable(name string, mlrval *types.Mlrval) {
	this.stackFrames.Front().Value.(*StackFrame).Set(name, mlrval)
}

// Used for the above BindVariable example where we look for outer-scope names,
// then set a new one only if not found in an outer scope.
func (this *Stack) SetVariable(name string, mlrval *types.Mlrval) {
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		if stackFrame.Has(name) {
			stackFrame.Set(name, mlrval)
			return
		}
	}
	this.BindVariable(name, mlrval)
}

// Returns nil on no-such
func (this *Stack) ReadVariable(name string) *types.Mlrval {

	// Scope-walk
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		mlrval := stackFrame.Get(name)
		if mlrval != nil {
			return mlrval
		}
	}
	return nil
}

// Returns nil on no-such
func (this *Stack) Dump() {
	fmt.Printf("STACK FRAMES (count %d):\n", this.stackFrames.Len())
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		fmt.Printf("  VARIABLES (count %d):\n", len(stackFrame.vars))
		for k, v := range stackFrame.vars {
			fmt.Printf("    %-16s %s\n", k, v.String())
		}
	}
}

// ----------------------------------------------------------------
type StackFrame struct {
	// TODO: just a map for now. In the C impl, pre-computation of
	// name-to-array-slot indices was an important optimization, especially for
	// compute-intensive scenarios.
	vars map[string]*types.Mlrval
}

func NewStackFrame() *StackFrame {
	return &StackFrame{
		vars: make(map[string]*types.Mlrval),
	}
}

func (this *StackFrame) Clear() {
	this.vars = make(map[string]*types.Mlrval)
}

func (this *StackFrame) Set(name string, mlrval *types.Mlrval) {
	this.vars[name] = mlrval.Copy()
}

// Returns nil on no such
func (this *StackFrame) Get(name string) *types.Mlrval {
	return this.vars[name]
}

// Returns nil on no such
func (this *StackFrame) Has(name string) bool {
	return this.vars[name] != nil
}
