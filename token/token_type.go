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

package token

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	// Identifiers + literals
	IDENT
	INT
	// Operators
	ASSIGN
	PLUS
	MINUS
	BANG
	ASTERISK
	SLASH
	LT
	GT
	EQ
	NOT_EQ
	// Delimiters
	COMMA
	SEMICOLON
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET
	// Keywords
	FUNCTION
	LET
	TRUE
	FALSE
	IF
	ELSE
	RETURN
	STRING
)

var tokenString = [...]string{
	ILLEGAL:   "illegal",
	EOF:       "eof",
	IDENT:     "identifier",
	INT:       "integer",
	ASSIGN:    "=",
	PLUS:      "+",
	MINUS:     "-",
	BANG:      "!",
	ASTERISK:  "*",
	SLASH:     "/",
	LT:        "<",
	GT:        ">",
	EQ:        "==",
	NOT_EQ:    "!=",
	COMMA:     ",",
	SEMICOLON: ";",
	LPAREN:    "(",
	RPAREN:    ")",
	LBRACE:    "{",
	RBRACE:    "}",
	LBRACKET:  "[",
	RBRACKET:  "]",
	FUNCTION:  "function",
	LET:       "let",
	TRUE:      "true",
	FALSE:     "false",
	IF:        "if",
	ELSE:      "else",
	RETURN:    "return",
	STRING:    "string",
}

func (t TokenType) String() string {
	if t < 0 || int(t) > len(tokenString) {
		return "invalid token"
	}
	return tokenString[t]
}

//TODO: define the language and add new token type
