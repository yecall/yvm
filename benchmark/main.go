/*
 * // Copyright (C) 2017 gyee authors
 * //
 * // This file is part of the gyee library.
 * //
 * // the gyee library is free software: you can redistribute it and/or modify
 * // it under the terms of the GNU General Public License as published by
 * // the Free Software Foundation, either version 3 of the License, or
 * // (at your option) any later version.
 * //
 * // the gyee library is distributed in the hope that it will be useful,
 * // but WITHOUT ANY WARRANTY; without even the implied warranty of
 * // MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * // GNU General Public License for more details.
 * //
 * // You should have received a copy of the GNU General Public License
 * // along with the gyee library.  If not, see <http://www.gnu.org/licenses/>.
 *
 *
 */

package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/yeeco/yvm/compiler"
	"github.com/yeeco/yvm/evaluator"
	"github.com/yeeco/yvm/lexer"
	"github.com/yeeco/yvm/object"
	"github.com/yeeco/yvm/parser"
	"github.com/yeeco/yvm/vm"
)

var engine = flag.String("engine", "vm", "use 'vm' or 'eval'")
var input = `
let fibonacci = fn(x) {
if (x == 0) { 0
     } else {
       if (x == 1) {
         return 1;
       } else {
         fibonacci(x - 1) + fibonacci(x - 2);
       }
} };
   fibonacci(35);
   `

func main() {
	flag.Parse()

	var duration time.Duration
	var result object.Object

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if *engine == "vm" {
		comp := compiler.NewCompiler()
		err := comp.Compile(program)
		if err != nil {
			fmt.Printf("compiler error: %s", err)
			return
		}

		machine := vm.NewVM(comp.Bytecode())

		start := time.Now()
		err = machine.Run()
		if err != nil {
			fmt.Printf("vm error: %s", err)
			return
		}

		duration = time.Since(start)
		result = machine.LastPoppedStackElem()
	} else if *engine == "eval" {
		env := object.NewEnvironment()
		start := time.Now()
		result = evaluator.Eval(program, env)
		duration = time.Since(start)
	} else if *engine == "native" {
		start := time.Now()
		res := fibonacci(35)
		result = &object.Integer{Value: int64(res)}

		duration = time.Since(start)
	}
	fmt.Printf(
		"engine=%s, result=%s, duration=%s\n",
		*engine,
		result.Inspect(),
		duration)
}

func fibonacci(x int) int {
	if x == 0 {
		return 0
	} else {
		if x == 1 {
			return 1
		} else {
			return fibonacci(x-1) + fibonacci(x-2)
		}
	}
}
