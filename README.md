# Generating a GraphQL Server with gqlgen

![](https://i.imgur.com/TlmffK0.png)

GraphQL is quickly supplanting REST as the front-end API standard, allowing clients to query exactly the data they require with end-to-end type safety. This workshop will show you how to generate a functioning GraphQL server starting from a base schema, and integrating local and external APIs into a single graph.

## Setup

For this workshop you will need:
- Go 1.9
- A `$GOPATH` setup in your environment
- `$GOTPATH/bin/` in your `$PATH`
- Run `gofmt` on save

### Dependencies

Install dependencies:

```shell
$ go get -v -u github.com/99designs/gqlgen/... \
    github.com/99designs/gqlgen-workshop/...
```

Create a folder for our project:

```shell
$ mkdir -p $GOPATH/src/github.com/[your-github-username]/gqlgen-workshop
$ cd $GOPATH/src/github.com/[your-github-username]/gqlgen-workshop
```

## What is GraphQL?

- A query language for an API
- Ask for your exact data requirements
- Schemas and type system

### How Does This Compare to REST?

- Multiple resources in a single query
- Reduce overfetching
- Schema provides a complete description of an endpoint


### How do we Write a GraphQL Server?

- Write Resolver functions that resolve the nodes and edges of a query
- Schema first — end to end types

## What Are We Building?

- In this workshop we're going to build a graph that exposes:
    - a local mock database
    - a remove Movie API that we can search for results on
    - a combination of both services into a single graph that allows users to like Movies
- This graph could be the starting point of a Movie liking application

## Local Database

So first imagine we have a local app with a User entity:

```go
type User struct {
    ID   int
    Name string
}
```

For this workshop we've provided this in the `db` package which also exposes these methods:

```go
func GetUser(id int) *User
func AddUser(name string)
```

### GraphQL Schema

To generate our GraphQL server, we first need schema.

Create a file named `schema.graphql`:

```go
type Query {
    user(id: ID!): User
}

type User {
    id: ID!
    name: String!
}
```

### Initialise gqlgen

```shell
$ gqlgen init
```

After this your project directory should look like:

```shell
$ ls
generated.go   gqlgen.yml     models_gen.go  resolver.go    schema.graphql server
```

### Type Mapping

`gqlgen` has generated a model for us in `models_gen.go`, but we have a `User` model already, so let's map to that.  Edit `gqlgen.yml`:

```yaml
models:
  User:
    model: github.com/99designs/gqlgen-workshop/db.User
```

And generate.  We're also going to remove resolver so that it's updated to point to our new package.

```shell
$ rm resolver.go
$ gqlgen
```

### Implement User Resolver

Now we can connect to our backend through our generated resolver.  Add the following implementation to `resolver.go`:

```go
func (r *queryResolver) User(ctx context.Context, id string) (*db.User, error) {
    user := db.GetUser(id)
    if user == nil {
        return nil, errors.New("User not found")
    }
    return user, nil
}
```

### Boot the Server!

```shell
$ go run server/server.go
```

Open `http://localhost:8080/` and try a query:

```go
query {
    user(id:"1") {
        name
    }
}

```

## External Services

The provided package `github.com/99designs/gqlgen-workshop/db` exposes an API client we can use to query a public movie database:

```go
type Movie struct {
	ID    string
	Title string
}

func Search(term string) ([]Movie, error)
```

We're now going to expose this API through our graph and integrate it with our `User` type.

### Movie — Schema

First we should update our schema.  Add a new Type to `schema.graphql`:

```go
type Movie {
    id: ID!
    title: String!
}
```

Then add a movie edge to the root query:

```go
type Query {
    // ... 
    movies(search: String!): [Movie!]!
}
```

### Movie — Type Mapping

Since we have a third party type, instead of having `gqlgen` generate a model for us, we can instead configure the type mapping in `gelgen.yml`:

```yaml
models:
  Movie:
    model: github.com/99designs/gqlgen-workshop/omdb.Movie
```

And run generate:

```shell
$ gqlgen
```

### Movie — Resolver Implementation

Now we need to implement the new resolver that has been generated for us:

```go
func (r *queryResolver) Movies(ctx context.Context, search string) ([]omdb.Movie, error) {
	return omdb.Search(search)
}
```

Running the server now, we can query for Movies!

```go
query {
    movies(search:"Star Wars") {
        id
        title
    }
}
```

## Likes

Now let's say we want to enable users to like movies — we need a mutation for the like action, and we need to expose likes for a user.  First we'll add our mutation; update `schema.graphql`:

```go
type Mutation {
    like(userId: ID!, movieId: ID!): User
}
```

And regenerate with `gqlgen`.

### Likes — Mutation Resolver

`gqlgen` has now generated a `MutationResolver` interface for us.  We can implement this like so in `resolver.go`:

```go
type mutationResolver struct{ *Resolver }

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *mutationResolver) Like(ctx context.Context, userId string, movieId string) (*db.User, error) {
	user := db.GetUser(userId)
	if user == nil {
		return nil, errors.New("User not found")
	}
	user.Like(movieId)
	return user, nil
}
```

### Likes — Retrieving For User

If we can like a movie for a user, we should then be able to retrieve all movies they have previously liked.  Let's add this ability to our `schema.graphql`:

```go
type User {
    // ...
    likes: [Movie!]!
}
```

And regenerate

### Likes — Retrieving Resolver

By default, `gqlgen` will generate us resolvers for any fields in our schema that it does not know about.   If you look in `generated.go` you will now see a new resolver for our `likes` field.

```go
type UserResolver interface {
	Likes(ctx context.Context, obj *db.User) ([]omdb.Movie, error)
}
```

Add the following to `resolver.go`:

```go
type userResolver struct{ *Resolver }

func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}

func (r *userResolver) Likes(ctx context.Context, u *db.User) ([]omdb.Movie, error) {
	return omdb.GetAll(u.Likes)
}

```

### Give it a Shot!

With everything in place, we can now test it out.  Start the server again and run this query:

```go
mutation {
    like(userId:"1", movieId:"tt0076759") {
        name
        likes {
            title
        }
    }
}
```

You should now be able to search for a movie title and add it to the list of liked movies for a user.

## Advanced Topics

- Testing
- Directives
- Subscriptions

