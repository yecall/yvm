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

package object

import (
	"fmt"
)

const (
	NULL_OBJ = "NULL"
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {return INTEGER_OBJ}
func (i *Integer) Inspect() string {return fmt.Sprintf("%d", i.Value)}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {return BOOLEAN_OBJ}
func (b *Boolean) Inspect() string {return fmt.Sprintf("%t", b.Value)}

type Null struct {}

func (n *Null) Type() ObjectType {return NULL_OBJ}
func (n *Null) Inspect() string {return "null"}
