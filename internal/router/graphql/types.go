package graphql

import (
	"github.com/graphql-go/graphql"
)

//nolint:gochecknoglobals
var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.NewNonNull(graphql.ID)},
		"firstName": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"lastName":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"nickname":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"email":     &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"country":   &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"createdAt": &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
		"updatedAt": &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
	},
})

//nolint:gochecknoglobals
var userCreateType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "UserCreate",
	Fields: graphql.InputObjectConfigFieldMap{
		"firstName": &graphql.InputObjectFieldConfig{Type: graphql.String},
		"lastName":  &graphql.InputObjectFieldConfig{Type: graphql.String},
		"nickname":  &graphql.InputObjectFieldConfig{Type: graphql.String},
		"email":     &graphql.InputObjectFieldConfig{Type: graphql.String},
		"country":   &graphql.InputObjectFieldConfig{Type: graphql.String},
		"password":  &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

//nolint:gochecknoglobals
var userUpdateType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "UserUpdate",
	Fields: graphql.InputObjectConfigFieldMap{
		"id":        &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.ID)},
		"firstName": &graphql.InputObjectFieldConfig{Type: graphql.String},
		"lastName":  &graphql.InputObjectFieldConfig{Type: graphql.String},
		"nickname":  &graphql.InputObjectFieldConfig{Type: graphql.String},
		"email":     &graphql.InputObjectFieldConfig{Type: graphql.String},
		"country":   &graphql.InputObjectFieldConfig{Type: graphql.String},
		"password":  &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

//nolint:gochecknoglobals
var usersType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Users",
	Fields: graphql.Fields{
		"users":      &graphql.Field{Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(userType)))},
		"totalCount": &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
		"nextPage":   &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
	},
})
