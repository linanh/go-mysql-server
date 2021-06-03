package analyzer

import (
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/plan"
)

func resolveCreateSelect(ctx *sql.Context, a *Analyzer, n sql.Node, scope *Scope) (sql.Node, error) {
	planCreate, ok := n.(*plan.CreateTable)
	if !ok || planCreate.Select() == nil {
		return n, nil
	}

	analyzedSelect, err := a.Analyze(ctx, planCreate.Select(), scope)
	if err != nil {
		return nil, err
	}

	// Get the correct schema of the CREATE TABLE based on the select query
	inputSpec := planCreate.TableSpec()
	selectSchema := analyzedSelect.Schema()
	mergedSchema := mergeSchemas(inputSpec.Schema, selectSchema)
	newSch := make(sql.Schema, len(mergedSchema))

	for i, col := range mergedSchema {
		tempCol := *col
		tempCol.Source = planCreate.Name()
		newSch[i] = &tempCol
	}

	newSpec := inputSpec.WithSchema(newSch)

	newCreateTable := plan.NewCreateTable(planCreate.Database(), planCreate.Name(), planCreate.IfNotExists(), planCreate.Temporary(), newSpec)

	return newCreateTable.WithSelect(stripQueryProcess(analyzedSelect)), nil
}

// mergeSchemas takes in the table spec of the CREATE TABLE and merges it with the schema used by the
// select query
func mergeSchemas(inputSchema sql.Schema, selectSchema sql.Schema) sql.Schema {
	// Get the matching columns between the two via name
	matchingColumns := make([]*sql.Column, 0)
	for _, col := range inputSchema {
		for _, col2 := range selectSchema {
			if col.Name == col2.Name {
				matchingColumns = append(matchingColumns, col)
			}
		}
	}

	// Get inputSchema exclusive columns
	leftExclusive := make([]*sql.Column, 0)
	for _, col := range inputSchema {
		found := false
		for _, col2 := range selectSchema {
			if col.Name == col2.Name {
				found = true
				break
			}
		}

		if !found {
			leftExclusive = append(leftExclusive, col)
		}
	}

	rightExclusive := make([]*sql.Column, 0)
	for _, col := range selectSchema {
		found := false
		for _, col2 := range inputSchema {
			if col.Name == col2.Name {
				found = true
				break
			}
		}

		if !found {
			rightExclusive = append(rightExclusive, col)
		}
	}

	return append(append(leftExclusive, matchingColumns...), rightExclusive...)
}