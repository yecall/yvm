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
	"bufio"
	"fmt"
	"io"

	"github.com/yeeco/yvm/compiler"
	"github.com/yeeco/yvm/evaluator"
	"github.com/yeeco/yvm/lexer"
	"github.com/yeeco/yvm/object"
	"github.com/yeeco/yvm/parser"
	"github.com/yeeco/yvm/vm"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		if line == "exit" {
			break
		}

		l := lexer.NewLexer(line)
		p := parser.NewParser(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		//io.WriteString(out, program.String())
		//io.WriteString(out, "\n")

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func StartVM(in io.Reader, out io.Writer, verbose bool) {
	scanner := bufio.NewScanner(in)

	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalsSize)
	symbolTable := compiler.NewSymbolTable()
	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		if line == "exit" {
			break
		}

		l := lexer.NewLexer(line)
		p := parser.NewParser(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

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
			printGlobals(out, globals)
		}

		machine := vm.NewVMWithGlobalsStore(code, globals)
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)
			continue
		}

		lastPopped := machine.LastPoppedStackElem()
		io.WriteString(out, lastPopped.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func printByteCode(out io.Writer, code *compiler.Bytecode){
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

func printGlobals(out io.Writer, globals []object.Object){
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
