# GraphQL


## Query

* Users:
     ```bash
     curl -X POST http://localhost:8081/api/v1/graphql \
          -H "Content-Type: application/json" \
          -d '{"query": "query { users(page: 1, limit: 10) { users { id firstName lastName nickname email country createdAt updatedAt } totalCount nextPage } }"}'
     ```

* User: insert user UUID from the `users` or `create` query:
     ```bash
     curl -X POST http://localhost:8081/api/v1/graphql \
     -H "Content-Type: application/json" \
     -d '{ "query": "query { user(id: \"3dc87204-a3fb-48d6-89be-a5ba85200462\") { id firstName lastName nickname email country createdAt updatedAt } }"}'
     ```


## Mutation

* Create
     ```bash
     curl -X POST http://localhost:8081/api/v1/graphql \
          -H "Content-Type: application/json" \
          -d '{"query": "mutation { create(input: { firstName: \"John\", lastName: \"Doe\", nickname: \"johnd1001\", email: \"john.doe.1001@example.com\", country: \"US\", password: \"securepassword\" }) { id firstName lastName nickname email country createdAt updatedAt } }"}'
     ```

* Update
```bash
curl -X POST http://localhost:8081/api/v1/graphql \
     -H "Content-Type: application/json" \
     -d '{"query": "mutation { update(input: { id: \"123e4567-e89b-12d3-a456-426614174000\", firstName: \"Johnny\" }) { id firstName lastName nickname email country createdAt updatedAt } }"}'
```

* Delete: insert user UUID from the `users` or `create` query:
     ```bash
     curl -X POST http://localhost:8081/api/v1/graphql \
     -H "Content-Type: application/json" \
     -d '{"query": "mutation { delete(id: \"dd60af0a-b9f6-4867-bf4c-b2d3b0658d8b\") }"}'
     ```
