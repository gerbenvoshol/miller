// Code generated by gocc; DO NOT EDIT.

package parser

import "miller/dsl/ast"

type (
	//TODO: change type and variable names to be consistent with other tables
	ProdTab      [numProductions]ProdTabEntry
	ProdTabEntry struct {
		String     string
		Id         string
		NTType     int
		Index      int
		NumSymbols int
		ReduceFunc func([]Attrib) (Attrib, error)
	}
	Attrib interface {
	}
)

var productionsTable = ProdTab{
	ProdTabEntry{
		String: `S' : StatementList	<<  >>`,
		Id:         "S'",
		NTType:     0,
		Index:      0,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `StatementList : Statement	<< ast.NewStatementList(X[0]) >>`,
		Id:         "StatementList",
		NTType:     1,
		Index:      1,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewStatementList(X[0])
		},
	},
	ProdTabEntry{
		String: `StatementList : StatementList Statement	<< ast.AppendStatement(X[0], X[1]) >>`,
		Id:         "StatementList",
		NTType:     1,
		Index:      2,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.AppendStatement(X[0], X[1])
		},
	},
	ProdTabEntry{
		String: `Statement : md_token_field_name md_token_assign md_token_number	<< ast.NewStatement(X[0]) >>`,
		Id:         "Statement",
		NTType:     2,
		Index:      3,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewStatement(X[0])
		},
	},
}
