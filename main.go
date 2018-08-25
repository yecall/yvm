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
	"fmt"
	"os"
	"os/user"

	"flag"
	"github.com/yeeco/yvm/repl"
)

var engine = flag.String("engine", "eval", "use 'eval' or 'vm'")
var verbose = flag.Bool("v", false, "verbose")

func main() {
	flag.Parse()

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s!\n", user.Username)
	fmt.Printf("Feel free to type in commands to %v\n", *engine)

	repl.Start(os.Stdout, *engine, *verbose)
}
