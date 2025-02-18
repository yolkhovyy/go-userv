package graphql

import (
	"github.com/graphql-go/graphql"
)

func (c *Controller) schemaConfig() graphql.SchemaConfig {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: c.user(),
			},
			"users": &graphql.Field{
				Type: usersType,
				Args: graphql.FieldConfigArgument{
					"page":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"limit":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"country": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: c.users(),
			},
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"create": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(userCreateType),
					},
				},
				Resolve: c.create(),
			},
			"update": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(userUpdateType),
					},
				},
				Resolve: c.update(),
			},
			"delete": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: c.delete(),
			},
		},
	})

	return graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	}
}
