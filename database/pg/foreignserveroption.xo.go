// Package pg contains the types for schema 'information_schema'.
package pg

// GENERATED BY XO. DO NOT EDIT.

import (
	"github.com/pkg/errors"
)

// ForeignServerOptionTable is the database name for the table.
const ForeignServerOptionTable = "information_schema.foreign_server_options"

// ForeignServerOption represents a row from 'information_schema.foreign_server_options'.
type ForeignServerOption struct {
	ForeignServerCatalog SQLIdentifier `json:"foreign_server_catalog"` // foreign_server_catalog
	ForeignServerName    SQLIdentifier `json:"foreign_server_name"`    // foreign_server_name
	OptionName           SQLIdentifier `json:"option_name"`            // option_name
	OptionValue          CharacterData `json:"option_value"`           // option_value
}

// Constants defining each column in the table.
const (
	ForeignServerOptionForeignServerCatalogField = "foreign_server_catalog"
	ForeignServerOptionForeignServerNameField    = "foreign_server_name"
	ForeignServerOptionOptionNameField           = "option_name"
	ForeignServerOptionOptionValueField          = "option_value"
)

// WhereClauses for every type in ForeignServerOption.
var (
	ForeignServerOptionForeignServerCatalogWhere SQLIdentifierField = "foreign_server_catalog"
	ForeignServerOptionForeignServerNameWhere    SQLIdentifierField = "foreign_server_name"
	ForeignServerOptionOptionNameWhere           SQLIdentifierField = "option_name"
	ForeignServerOptionOptionValueWhere          CharacterDataField = "option_value"
)

// QueryOneForeignServerOption retrieves a row from 'information_schema.foreign_server_options' as a ForeignServerOption.
func QueryOneForeignServerOption(db XODB, where WhereClause, order OrderBy) (*ForeignServerOption, error) {
	const origsqlstr = `SELECT ` +
		`foreign_server_catalog, foreign_server_name, option_name, option_value ` +
		`FROM information_schema.foreign_server_options WHERE (`

	idx := 1
	sqlstr := origsqlstr + where.String(&idx) + ") " + order.String() + " LIMIT 1"

	fso := &ForeignServerOption{}
	err := db.QueryRow(sqlstr, where.Values()...).Scan(&fso.ForeignServerCatalog, &fso.ForeignServerName, &fso.OptionName, &fso.OptionValue)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return fso, nil
}

// QueryForeignServerOption retrieves rows from 'information_schema.foreign_server_options' as a slice of ForeignServerOption.
func QueryForeignServerOption(db XODB, where WhereClause, order OrderBy) ([]*ForeignServerOption, error) {
	const origsqlstr = `SELECT ` +
		`foreign_server_catalog, foreign_server_name, option_name, option_value ` +
		`FROM information_schema.foreign_server_options WHERE (`

	idx := 1
	sqlstr := origsqlstr + where.String(&idx) + ") " + order.String()

	var vals []*ForeignServerOption
	q, err := db.Query(sqlstr, where.Values()...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for q.Next() {
		fso := ForeignServerOption{}

		err = q.Scan(&fso.ForeignServerCatalog, &fso.ForeignServerName, &fso.OptionName, &fso.OptionValue)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		vals = append(vals, &fso)
	}
	return vals, nil
}

// AllForeignServerOption retrieves all rows from 'information_schema.foreign_server_options' as a slice of ForeignServerOption.
func AllForeignServerOption(db XODB, order OrderBy) ([]*ForeignServerOption, error) {
	const origsqlstr = `SELECT ` +
		`foreign_server_catalog, foreign_server_name, option_name, option_value ` +
		`FROM information_schema.foreign_server_options`

	sqlstr := origsqlstr + order.String()

	var vals []*ForeignServerOption
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for q.Next() {
		fso := ForeignServerOption{}

		err = q.Scan(&fso.ForeignServerCatalog, &fso.ForeignServerName, &fso.OptionName, &fso.OptionValue)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		vals = append(vals, &fso)
	}
	return vals, nil
}