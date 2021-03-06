// Copyright 2017 The WL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate golex -o input.go input.l
//go:generate golex -o scanner.go scanner.l
//go:generate yy -kind Case -node Node -astImport "\"fmt\"\n\n\"go/token\"" parser.yy
//go:generate goyacc -xegen tmp -o parser.go parser.y
//go:generate touch xerrors
//go:generate sh -c "cat xerrors tmp > xegen"
//go:generate goyacc -dlval "__yyfmt__.Sprint(yylex.(*lexer).sdump())+` `+prettyString(lval.Token)" -dlvalf %v -ex -fs -xe xegen -o parser.go parser.y
//go:generate rm -f tmp xegen
//go:generate sh -c "go test -run ^Example | fe"

package wl // import "modernc.org/wl"
