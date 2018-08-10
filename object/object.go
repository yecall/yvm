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
	ERROR_OBJ = "ERROR"
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Error struct {
	Message string
}
//TODO: 需要增加stack trace， line number，column number，lexer中增加后这儿也要增加

func (e *Error) Type() ObjectType {return ERROR_OBJ}
func (e *Error) Inspect() string {return "ERROR:" + e.Message}

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

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {return RETURN_VALUE_OBJ}
func (rv *ReturnValue) Inspect() string {return rv.Value.Inspect()}


type Environment struct {
	store map[string] Object
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

//TODO: go里的primitive类型也可以实现接口的，可以用来提高性能？
//TODO: 这里object type也可以用按token的实现不用string
