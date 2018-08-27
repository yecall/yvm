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

package repl

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterh/liner"
	"github.com/yeeco/yvm/compiler"
	"github.com/yeeco/yvm/evaluator"
	"github.com/yeeco/yvm/lexer"
	"github.com/yeeco/yvm/object"
	"github.com/yeeco/yvm/parser"
	"github.com/yeeco/yvm/vm"
)

const PROMPT = ">> "

func init() {
	//gob.Register(&object.Boolean{}) //?这个不注册貌似可以？因为bool值不会进constant
	gob.Register(&object.Integer{})
	gob.Register(&object.String{})
	gob.Register(&object.CompiledFunction{})
}

func Start(out io.Writer, engine string, verbose bool) {
	l := liner.NewLiner()
	defer l.Close()

	var lines []string
	l.SetCtrlCAborts(true)
	l.SetMultiLineMode(true)

	history := filepath.Join(os.TempDir(), ".yvm_history")
	if f, err := os.Open(history); err == nil {
		l.ReadHistory(f)
		f.Close()
	}

	// for eval
	var env *object.Environment

	// for vm
	var constants []object.Object
	var globals []object.Object
	var symbolTable *compiler.SymbolTable

	if engine == "eval" {
		env = object.NewEnvironment()
	} else if engine == "vm" {
		constants = []object.Object{}
		globals = make([]object.Object, vm.GlobalsSize)
		symbolTable = compiler.NewSymbolTable()
		for i, v := range object.Builtins {
			symbolTable.DefineBuiltin(i, v.Name)
		}
	} else {
		fmt.Println("Error engine")
		return
	}

	for {
		if line, err := l.Prompt(PROMPT); err == nil {
			if line == "exit" || line == "quit" {
				if f, err := os.Create(history); err == nil {
					l.WriteHistory(f)
					f.Close()
				}
				break
			}

			tmp := strings.TrimSpace(line)
			if len(tmp) == 0 {
				continue
			}

			if tmp[len(tmp)-1:] == "\\" {
				lines = append(lines, strings.TrimRight(tmp, "\\"))
				continue
			} else {
				lines = append(lines, tmp)
			}

			input := strings.Join(lines, "")
			l.AppendHistory(input)
			lines = nil

			l := lexer.NewLexer(input)
			p := parser.NewParser(l)

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				continue
			}

			if engine == "eval" {
				evaluated := evaluator.Eval(program, env)
				if evaluated != nil {
					io.WriteString(out, evaluated.Inspect())
					io.WriteString(out, "\n")
				}
			} else if engine == "vm" {
				comp := compiler.NewCompilerWithState(symbolTable, constants)
				err := comp.Compile(program)
				if err != nil {
					fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
					continue
				}

				code := comp.Bytecode()
				constants = code.Constants

				if verbose {
					printByteCode(out, code)
				}

				code = Deserialize(Serialize(code))

				machine := vm.NewVMWithGlobalsStore(code, globals)
				err = machine.Run()
				if err != nil {
					fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)
					continue
				}

				if verbose {
					printGlobals(out, globals)
				}

				lastPopped := machine.LastPoppedStackElem()
				io.WriteString(out, lastPopped.Inspect())
				io.WriteString(out, "\n")
			}
		}
	}
}

func Serialize(code *compiler.Bytecode) []byte {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(code); err != nil {
		fmt.Println("encode error:", err)
	}

	return buf.Bytes()
}

func Deserialize(c []byte) *compiler.Bytecode {
	var buf bytes.Buffer

	buf.Write(c)
	dec := gob.NewDecoder(&buf)

	bc := compiler.Bytecode{}
	if err := dec.Decode(&bc); err != nil {
		fmt.Println("decode error:", err)
	}

	return &bc
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func printByteCode(out io.Writer, code *compiler.Bytecode) {
	fmt.Fprintln(out, "Constants:")
	for i, obj := range code.Constants {
		fmt.Fprintf(out, "%4d: %s", i, obj.Inspect())
		fmt.Fprintln(out)

		if obj.Type() == object.COMPILED_FUNCTION_OBJ {
			o := obj.(*object.CompiledFunction)
			fmt.Fprintf(out, o.Instructions.String())
		}
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "Instructions:")
	fmt.Fprintf(out, code.Instructions.String())
}

func printGlobals(out io.Writer, globals []object.Object) {
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Globals:")
	for i, obj := range globals {
		if obj == nil {
			break
		}
		fmt.Fprintf(out, "%4d: %s", i, obj.Inspect())
		fmt.Fprintln(out)

		if obj.Type() == object.COMPILED_FUNCTION_OBJ {
			o := obj.(*object.CompiledFunction)
			fmt.Fprintf(out, o.Instructions.String())
		}
	}
	fmt.Fprintln(out)
}

//TODO: 加一个开关能够显示op code，把字节码print出来
//TODO: Bytecode的序列化和反序列化
//TODO: repl支持多行输入
