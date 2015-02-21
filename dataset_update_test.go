package gql

import (
	"database/sql/driver"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func (me *datasetTest) TestUpdateSqlWithStructs() {
	t := me.T()
	ds1 := From("items")
	type item struct {
		Address string `db:"address"`
		Name    string `db:"name"`
	}
	sql, err := ds1.UpdateSql(item{Name: "Test", Address: "111 Test Addr"})
	assert.NoError(t, err)
	assert.Equal(t, sql, `UPDATE "items" SET "address"='111 Test Addr',"name"='Test'`)
}

func (me *datasetTest) TestUpdateSqlWithMaps() {
	t := me.T()
	ds1 := From("items")
	sql, err := ds1.UpdateSql(map[string]interface{}{"name": "Test", "address": "111 Test Addr"})
	assert.NoError(t, err)
	assert.Equal(t, sql, `UPDATE "items" SET "address"='111 Test Addr',"name"='Test'`)

}

func (me *datasetTest) TestUpdateSqlWithByteSlice() {
	t := me.T()
	ds1 := From("items")
	type item struct {
		Name string `db:"name"`
		Data []byte `db:"data"`
	}
	sql, err := ds1.Returning(I("items").All()).UpdateSql(item{Name: "Test", Data: []byte(`{"someJson":"data"}`)})
	assert.NoError(t, err)
	assert.Equal(t, sql, `UPDATE "items" SET "name"='Test',"data"='{"someJson":"data"}' RETURNING "items".*`)
}

type valuerType []byte

func (j valuerType) Value() (driver.Value, error) {
	return []byte(fmt.Sprintf("%s World", string(j))), nil
}

func (me *datasetTest) TestUpdateSqlWithValuer() {
	t := me.T()
	ds1 := From("items")
	type item struct {
		Name string     `db:"name"`
		Data valuerType `db:"data"`
	}
	sql, err := ds1.Returning(I("items").All()).UpdateSql(item{Name: "Test", Data: []byte(`Hello`)})
	assert.NoError(t, err)
	assert.Equal(t, sql, `UPDATE "items" SET "name"='Test',"data"='Hello World' RETURNING "items".*`)
}

func (me *datasetTest) TestUpdateSqlWithUnsupportedType() {
	t := me.T()
	ds1 := From("items")
	_, err := ds1.UpdateSql([]string{"HELLO"})
	assert.EqualError(t, err, "gql: Unsupported update interface type []string")
}

func (me *datasetTest) TestUpdateSqlWithSkipupdateTag() {
	t := me.T()
	ds1 := From("items")
	type item struct {
		Address string `db:"address" gql:"skipupdate"`
		Name    string `db:"name"`
	}
	sql, err := ds1.UpdateSql(item{Name: "Test", Address: "111 Test Addr"})
	assert.NoError(t, err)
	assert.Equal(t, sql, `UPDATE "items" SET "name"='Test'`)
}

func (me *datasetTest) TestUpdateSqlWithWhere() {
	t := me.T()
	ds1 := From("items")
	type item struct {
		Address string `db:"address"`
		Name    string `db:"name"`
	}
	sql, err := ds1.Where(I("name").IsNull()).UpdateSql(item{Name: "Test", Address: "111 Test Addr"})
	assert.NoError(t, err)
	assert.Equal(t, sql, `UPDATE "items" SET "address"='111 Test Addr',"name"='Test' WHERE ("name" IS NULL)`)

	sql, err = ds1.Where(I("name").IsNull()).UpdateSql(map[string]interface{}{"name": "Test", "address": "111 Test Addr"})
	assert.NoError(t, err)
	assert.Equal(t, sql, `UPDATE "items" SET "address"='111 Test Addr',"name"='Test' WHERE ("name" IS NULL)`)
}

func (me *datasetTest) TestUpdateSqlWithReturning() {
	t := me.T()
	ds1 := From("items")
	type item struct {
		Address string `db:"address"`
		Name    string `db:"name"`
	}
	sql, err := ds1.Returning(I("items").All()).UpdateSql(item{Name: "Test", Address: "111 Test Addr"})
	assert.NoError(t, err)
	assert.Equal(t, sql, `UPDATE "items" SET "address"='111 Test Addr',"name"='Test' RETURNING "items".*`)

	sql, err = ds1.Where(I("name").IsNull()).Returning(Literal(`"items".*`)).UpdateSql(map[string]interface{}{"name": "Test", "address": "111 Test Addr"})
	assert.NoError(t, err)
	assert.Equal(t, sql, `UPDATE "items" SET "address"='111 Test Addr',"name"='Test' WHERE ("name" IS NULL) RETURNING "items".*`)
}