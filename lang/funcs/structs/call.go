// Mgmt
// Copyright (C) James Shubin and the project contributors
// Written by James Shubin <james@shubin.ca> and the project contributors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//
// Additional permission under GNU GPL version 3 section 7
//
// If you modify this program, or any covered work, by linking or combining it
// with embedded mcl code and modules (and that the embedded mcl code and
// modules which link with this program, contain a copy of their source code in
// the authoritative form) containing parts covered by the terms of any other
// license, the licensors of this program grant you additional permission to
// convey the resulting work. Furthermore, the licensors of this program grant
// the original author, James Shubin, additional permission to update this
// additional permission if he deems it necessary to achieve the goals of this
// additional permission.

package structs

import (
	"context"
	"fmt"
	"sync"

	"github.com/purpleidea/mgmt/lang/interfaces"
	"github.com/purpleidea/mgmt/lang/types"
	"github.com/purpleidea/mgmt/lang/types/full"
	"github.com/purpleidea/mgmt/util/errwrap"
)

const (
	// CallFuncName is the unique name identifier for this function.
	CallFuncName = "call"

	// CallFuncArgNameFunction is the name for the edge which connects the
	// input function to CallFunc.
	CallFuncArgNameFunction = "fn"
)

// CallFunc receives a function from upstream, but not the arguments. Instead,
// the Funcs which emit those arguments must be specified at construction time.
// The arguments are connected to the received FuncValues in such a way that
// CallFunc emits the result of applying the function to the arguments.
type CallFunc struct {
	interfaces.Textarea

	Type     *types.Type // the type of the result of applying the function
	FuncType *types.Type // the type of the function
	EdgeName string      // name of the edge used

	ArgVertices []interfaces.Func

	OutputVertex interfaces.Func

	init *interfaces.Init

	lastFuncValue *full.FuncValue // remember the last function value
}

// String returns a simple name for this function. This is needed so this struct
// can satisfy the pgraph.Vertex interface.
func (obj *CallFunc) String() string {
	return CallFuncName
}

// Validate makes sure we've built our struct properly.
func (obj *CallFunc) Validate() error {
	if obj.Type == nil {
		return fmt.Errorf("must specify a type")
	}
	if obj.FuncType == nil {
		return fmt.Errorf("must specify a func type")
	}
	// TODO: maybe we can remove this if we use this for core functions...
	if obj.EdgeName == "" {
		return fmt.Errorf("must specify an edge name")
	}
	typ := obj.FuncType
	// we only care about the output type of calling our func
	if err := obj.Type.Cmp(typ.Out); err != nil {
		return errwrap.Wrapf(err, "call expr type must match func out type")
	}
	if len(obj.ArgVertices) != len(typ.Ord) {
		return fmt.Errorf("number of arg Funcs must match number of func args in the type")
	}

	if obj.OutputVertex == nil {
		return fmt.Errorf("the output vertex is missing")
	}

	return nil
}

// Info returns some static info about itself.
func (obj *CallFunc) Info() *interfaces.Info {
	var typ *types.Type
	if obj.Type != nil && obj.FuncType != nil { // don't panic if called speculatively
		typ = types.NewType(fmt.Sprintf("func(%s %s) %s", obj.EdgeName, obj.FuncType, obj.Type))
	}

	return &interfaces.Info{
		Pure: true,
		Memo: false, // TODO: ???
		Sig:  typ,
		Err:  obj.Validate(),
	}
}

// Init runs some startup code for this composite function.
func (obj *CallFunc) Init(init *interfaces.Init) error {
	obj.init = init
	obj.lastFuncValue = nil
	return nil
}

// Stream takes an input struct in the format as described in the Func and Graph
// methods of the Expr, and returns the actual expected value as a stream based
// on the changing inputs to that value.
func (obj *CallFunc) Stream(ctx context.Context) error {
	// XXX: is there a sync.Once sort of solution that would be more elegant here?
	mutex := &sync.Mutex{}
	done := false
	send := func(ctx context.Context, b bool) error {
		mutex.Lock()
		defer mutex.Unlock()
		if done {
			return nil
		}
		done = true
		defer close(obj.init.Output) // the sender closes

		if !b {
			return nil
		}

		// send dummy value to the output
		select {
		case obj.init.Output <- types.NewFloat(): // XXX: dummy value
		case <-ctx.Done():
			return ctx.Err()
		}

		return nil
	}
	defer send(ctx, false) // just close

	defer func() {
		obj.init.Txn.Reverse()
	}()

	for {
		select {
		case input, ok := <-obj.init.Input:
			if !ok {
				obj.init.Input = nil // block looping back here
				if !done {
					return fmt.Errorf("input closed without ever sending anything")
				}
				return nil
			}

			value, exists := input.Struct()[obj.EdgeName]
			if !exists {
				return fmt.Errorf("programming error, can't find edge")
			}

			newFuncValue, ok := value.(*full.FuncValue)
			if !ok {
				return fmt.Errorf("programming error, can't convert to *FuncValue")
			}

			// It's important to have this compare step to avoid
			// redundant graph replacements which slow things down,
			// but also cause the engine to lock, which can preempt
			// the process scheduler, which can cause duplicate or
			// unnecessary re-sending of values here, which causes
			// the whole process to repeat ad-nauseum.
			if newFuncValue == obj.lastFuncValue {
				continue
			}
			// If we have a new function, then we need to replace
			// the subgraph with a new one that uses the new
			// function.
			obj.lastFuncValue = newFuncValue

			if err := obj.replaceSubGraph(newFuncValue); err != nil {
				return errwrap.Wrapf(err, "could not replace subgraph")
			}

			send(ctx, true) // send dummy and then close

			continue

		case <-ctx.Done():
			return nil
		}
	}
}

func (obj *CallFunc) replaceSubGraph(newFuncValue *full.FuncValue) error {
	// Create a subgraph which looks as follows.
	//
	// digraph {
	//   ArgVertices[0] -> "fn" # spawned
	//   ArgVertices[1] -> "fn" # spawned
	//   ArgVertices[2] -> "fn" # spawned
	//
	//   "fn" -> "callSubgraphOutput" # spawned
	//
	//   "FuncValue" -> "CallFunc" # fn
	//   "CallFunc" -> "callSubgraphOutput" # dummy
	//   "callSubgraphOutput" -> downstream
	// }

	// delete the old subgraph
	if err := obj.init.Txn.Reverse(); err != nil {
		return errwrap.Wrapf(err, "could not Reverse")
	}

	// This passed in Txn has AddVertex, AddEdge, and possibly AddGraph
	// methods called on it. Nothing else. It will _not_ call Commit or
	// Reverse. It adds to the graph, and our Commit and Reverse operations
	// are the ones that actually make the change.
	outputFunc, err := newFuncValue.CallWithFuncs(obj.init.Txn, obj.ArgVertices)
	if err != nil {
		return errwrap.Wrapf(err, "could not call newFuncValue.Call()")
	}

	// create the new subgraph
	edgeName := OutputFuncArgName
	edge := &interfaces.FuncEdge{Args: []string{edgeName}}
	obj.init.Txn.AddVertex(outputFunc)
	obj.init.Txn.AddEdge(outputFunc, obj.OutputVertex, edge)

	return obj.init.Txn.Commit()
}
