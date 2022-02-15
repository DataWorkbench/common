package gormwrap

import (
	"gorm.io/gorm/clause"
)

// BuildClauseExpressionEQWithOR used to construct the conditions that the specified column
// equal to any value in "items".
func BuildClauseExpressionEQWithOR(column string, items []string) clause.Expression {
	if len(items) == 0 {
		panic("gorm BuildClauseExpressionEQWithOR: items cannot be empty")
	}

	if column == "" {
		panic("gorm BuildClauseExpressionEQWithOR: column cannot be empty")
	}

	eqExpr := make([]clause.Expression, len(items))
	for i := 0; i < len(items); i++ {
		value := items[i]
		if value == "" {
			panic("gorm BuildClauseExpressionEQWithOR: the item value cannot be empty")
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
