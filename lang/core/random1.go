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

package core // TODO: should this be in its own individual package?

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/purpleidea/mgmt/lang/funcs"
	"github.com/purpleidea/mgmt/lang/interfaces"
	"github.com/purpleidea/mgmt/lang/types"
	"github.com/purpleidea/mgmt/util/errwrap"
)

const (
	// Random1FuncName is the name this function is registered as.
	Random1FuncName = "random1"

	// arg names...
	random1ArgNameLength = "length"

	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func init() {
	funcs.Register(Random1FuncName, func() interfaces.Func { return &Random1Func{} })
}

// Random1Func returns one random string of a certain length.
// XXX: return a stream instead, and combine this with a first(?) function which
// takes the first value and then puts backpressure on the stream. This should
// notify parent functions somehow that their values are no longer required so
// that they can shutdown if possible. Maybe it should be returning a stream of
// floats [0,1] as well, which someone can later map to the alphabet that they
// want. Should random() take an interval to know how often to spit out values?
// It could also just do it once per second, and we could filter for less. If we
// want something high precision, we could add that in the future... We could
// name that "random" and this one can be "random1" until we deprecate it.
type Random1Func struct {
	init *interfaces.Init

	finished bool // did we send the random string?
}

// String returns a simple name for this function. This is needed so this struct
// can satisfy the pgraph.Vertex interface.
func (obj *Random1Func) String() string {
	return Random1FuncName
}

// ArgGen returns the Nth arg name for this function.
func (obj *Random1Func) ArgGen(index int) (string, error) {
	seq := []string{random1ArgNameLength}
	if l := len(seq); index >= l {
		return "", fmt.Errorf("index %d exceeds arg length of %d", index, l)
	}
	return seq[index], nil
}

// Validate makes sure we've built our struct properly. It is usually unused for
// normal functions that users can use directly.
func (obj *Random1Func) Validate() error {
	return nil
}

// Info returns some static info about itself.
func (obj *Random1Func) Info() *interfaces.Info {
	return &interfaces.Info{
		Pure: false,
		Memo: false,
		Fast: false,
		Spec: false,
		Sig:  types.NewType(fmt.Sprintf("func(%s int) str", random1ArgNameLength)),
		Err:  obj.Validate(),
	}
}

// generate generates a random string.
func generate(length uint16) (string, error) {
	max := len(alphabet) - 1 // last index
	output := ""

	// FIXME: have someone verify this is cryptographically secure & correct
	for i := uint16(0); i < length; i++ {
		big, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
		if err != nil {
			return "", errwrap.Wrapf(err, "could not generate random string")
		}
		ix := big.Int64()
		output += string(alphabet[ix])
	}

	if length != 0 && output == "" { // safety against empty strings
		return "", fmt.Errorf("string is empty")
	}

	if uint16(len(output)) != length { // safety against weird bugs
		return "", fmt.Errorf("random string is too short") // bug!
	}

	return output, nil
}

// Init runs some startup code for this function.
func (obj *Random1Func) Init(init *interfaces.Init) error {
	obj.init = init
	return nil
}

// Stream returns the single value that was generated and then closes.
func (obj *Random1Func) Stream(ctx context.Context) error {
	defer close(obj.init.Output) // the sender closes
	var result string
	for {
		select {
		case input, ok := <-obj.init.Input:
			if !ok {
				return nil // can't output any more
			}
			//if err := input.Type().Cmp(obj.Info().Sig.Input); err != nil {
			//	return errwrap.Wrapf(err, "wrong function input")
			//}

			if obj.finished {
				// TODO: continue instead?
				return fmt.Errorf("you can only pass a single input to random")
			}

			length := input.Struct()[random1ArgNameLength].Int()
			// TODO: if negative, randomly pick a length ?
			if length < 0 {
				return fmt.Errorf("can't generate a negative length")
			}

			var err error
			if result, err = generate(uint16(length)); err != nil {
				return err // no errwrap needed b/c helper func
			}

		case <-ctx.Done():
			return nil
		}

		select {
		case obj.init.Output <- &types.StrValue{
			V: result,
		}:
			// we only send one value, then wait for input to close
			obj.finished = true

		case <-ctx.Done():
			return nil
		}
	}
}
