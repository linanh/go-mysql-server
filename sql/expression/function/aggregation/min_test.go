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

package aggregation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/linanh/go-mysql-server/sql"
	"github.com/linanh/go-mysql-server/sql/expression"
)

func TestMin_Name(t *testing.T) {
	assert := require.New(t)

	m := NewMin(sql.NewEmptyContext(), expression.NewGetField(0, sql.Int32, "field", true))
	assert.Equal("MIN(field)", m.String())
}

func TestMin_Eval_Int32(t *testing.T) {
	assert := require.New(t)
	ctx := sql.NewEmptyContext()

	m := NewMin(sql.NewEmptyContext(), expression.NewGetField(0, sql.Int32, "field", true))
	b := m.NewBuffer()

	m.Update(ctx, b, sql.NewRow(int32(7)))
	m.Update(ctx, b, sql.NewRow(int32(2)))
	m.Update(ctx, b, sql.NewRow(nil))

	v, err := m.Eval(ctx, b)
	assert.NoError(err)
	assert.Equal(int32(2), v)
}

func TestMin_Eval_Text(t *testing.T) {
	assert := require.New(t)
	ctx := sql.NewEmptyContext()

	m := NewMin(sql.NewEmptyContext(), expression.NewGetField(0, sql.Text, "field", true))
	b := m.NewBuffer()

	m.Update(ctx, b, sql.NewRow("a"))
	m.Update(ctx, b, sql.NewRow("A"))
	m.Update(ctx, b, sql.NewRow("b"))

	v, err := m.Eval(ctx, b)
	assert.NoError(err)
	assert.Equal("A", v)
}

func TestMin_Eval_Timestamp(t *testing.T) {
	assert := require.New(t)
	ctx := sql.NewEmptyContext()

	m := NewMin(sql.NewEmptyContext(), expression.NewGetField(0, sql.Timestamp, "field", true))
	b := m.NewBuffer()

	expected, _ := time.Parse(sql.TimestampDatetimeLayout, "2006-01-02 15:04:05")
	someTime, _ := time.Parse(sql.TimestampDatetimeLayout, "2007-01-02 15:04:05")
	otherTime, _ := time.Parse(sql.TimestampDatetimeLayout, "2008-01-02 15:04:05")

	m.Update(ctx, b, sql.NewRow(someTime))
	m.Update(ctx, b, sql.NewRow(expected))
	m.Update(ctx, b, sql.NewRow(otherTime))

	v, err := m.Eval(ctx, b)
	assert.NoError(err)
	assert.Equal(expected, v)
}

func TestMin_Eval_NULL(t *testing.T) {
	assert := require.New(t)
	ctx := sql.NewEmptyContext()

	m := NewMin(sql.NewEmptyContext(), expression.NewGetField(0, sql.Int32, "field", true))
	b := m.NewBuffer()

	m.Update(ctx, b, sql.NewRow(nil))
	m.Update(ctx, b, sql.NewRow(nil))
	m.Update(ctx, b, sql.NewRow(nil))

	v, err := m.Eval(ctx, b)
	assert.NoError(err)
	assert.Equal(nil, v)
}

func TestMin_Eval_Empty(t *testing.T) {
	assert := require.New(t)
	ctx := sql.NewEmptyContext()

	m := NewMin(sql.NewEmptyContext(), expression.NewGetField(0, sql.Int32, "field", true))
	b := m.NewBuffer()

	v, err := m.Eval(ctx, b)
	assert.NoError(err)
	assert.Equal(nil, v)
}

func TestMin_Distinct(t *testing.T) {
	assert := require.New(t)
	ctx := sql.NewEmptyContext()

	m := NewMin(sql.NewEmptyContext(), expression.NewDistinctExpression(expression.NewGetField(0, sql.Int32, "field", true)))
	b := m.NewBuffer()

	require.Equal(t, "MIN(DISTINCT field)", m.String())

	require.NoError(t, m.Update(ctx, b, sql.Row{1}))
	require.NoError(t, m.Update(ctx, b, sql.Row{1}))
	require.NoError(t, m.Update(ctx, b, sql.Row{2}))
	require.NoError(t, m.Update(ctx, b, sql.Row{3}))
	require.NoError(t, m.Update(ctx, b, sql.Row{3}))

	v, err := m.Eval(ctx, b)
	assert.NoError(err)
	assert.Equal(1, v)

	m = NewMin(sql.NewEmptyContext(), expression.NewDistinctExpression(expression.NewGetField(0, sql.Int32, "field", true)))
	b = m.NewBuffer()

	require.Equal(t, "MIN(DISTINCT field)", m.String())

	require.NoError(t, m.Update(ctx, b, sql.Row{1}))
	require.NoError(t, m.Update(ctx, b, sql.Row{1}))
	require.NoError(t, m.Update(ctx, b, sql.Row{2}))
	require.NoError(t, m.Update(ctx, b, sql.Row{nil}))
	require.NoError(t, m.Update(ctx, b, sql.Row{nil}))
	require.NoError(t, m.Update(ctx, b, sql.Row{3}))

	v, err = m.Eval(ctx, b)
	assert.NoError(err)
	assert.Equal(1, v)
}
