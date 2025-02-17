// Copyright 2020-2021 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plan

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/linanh/go-mysql-server/memory"
	"github.com/linanh/go-mysql-server/sql"
	"github.com/linanh/go-mysql-server/sql/expression"
)

func TestSubqueryAliasSchema(t *testing.T) {
	require := require.New(t)

	tableSchema := sql.Schema{
		{Name: "foo", Type: sql.Text, Nullable: false, Source: "bar"},
		{Name: "baz", Type: sql.Text, Nullable: false, Source: "bar"},
	}

	subquerySchema := sql.Schema{
		{Name: "foo", Type: sql.Text, Nullable: false, Source: "alias"},
		{Name: "baz", Type: sql.Text, Nullable: false, Source: "alias"},
	}

	table := memory.NewTable("bar", tableSchema)

	subquery := NewProject(
		[]sql.Expression{
			expression.NewGetField(0, sql.Text, "foo", false),
			expression.NewGetField(1, sql.Text, "baz", false),
		},
		NewResolvedTable(table, nil, nil),
	)

	require.Equal(
		subquerySchema,
		NewSubqueryAlias("alias", "", subquery).Schema(),
	)
}
