package gormwrap

import (
	"gorm.io/gorm/clause"
)

// BuildClauseExpressionEQWithOR used to construct the conditions that the specified column
// equal to any value in "items".
// Deprecated: use gormwrap.BuildConditionClauseInWithString instead.
func BuildClauseExpressionEQWithOR(column string, items []string) clause.Expression {
	if len(items) == 0 {
		panic("gormwrap BuildClauseExpressionEQWithOR: items cannot be empty")
	}

	if column == "" {
		panic("gormwrap BuildClauseExpressionEQWithOR: column cannot be empty")
	}

	eqExpr := make([]clause.Expression, len(items))
	for i := 0; i < len(items); i++ {
		value := items[i]
		if value == "" {
			panic("gormwrap BuildClauseExpressionEQWithOR: the item value cannot be empty")
		}
		eqExpr[i] = clause.Eq{Column: column, Value: value}
	}
	var expr clause.Expression
	if len(eqExpr) == 1 {
		expr = eqExpr[0]
	} else {
		expr = clause.Or(eqExpr...)
	}
	return expr
}

// BuildConditionClauseInWithString used build the conditions of `Where IN` with string slice.
func BuildConditionClauseInWithString(column string, items []string) clause.IN {
	if len(items) == 0 {
		panic("gormwrap BuildConditionClauseInWithString: items cannot be empty")
	}

	if column == "" {
		panic("gormwrap BuildConditionClauseInWithString: column cannot be empty")
	}

	values := make([]interface{}, len(items))
	for i := 0; i < len(items); i++ {
		values[i] = items[i]
	}
	return clause.IN{Column: column, Values: values}
}
