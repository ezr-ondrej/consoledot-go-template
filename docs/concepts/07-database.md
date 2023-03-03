# Database

When we want to keep state, we most likely want a relational database for our service.
ConsoleDot platform has settled on using PostgreSQL for all apps.
This give apps clear choice for drivers, we can build our service with only single database in mind.

When building a go application we could use [GORM](https://gorm.io/) that deals with any database and has many ORM features.
In this service tho, we expect to the database model being quite simple and use low level drivers.
This gives us much more power to use pure go objects and strong typing.

## Migrations

We use simple tool tern for migrations and migrations written in pure SQL.
It allows to lock database during migrations.
We use go embedding for migrations, so migrations are embedded into the migration binary.

## Code structure

### DB package
We keep database initialization in a separate package.
Here we just initialize the database connection global pool that can be accessed from other packages.
Initialization `db.Initialize` accepts schema, we will use this for integration testing later on.

### DAO
To organize our code, we want some structure and this service uses Data access objects.
These are objects that abstract the data fetching process and define data access interfaces.

This allows for easy database abstraction.
We can use this during unit testing, to stub away database.

### Model
Model is simple data structure that represents an entity of state.

## DAO methods

Example dao method, that accepts a model and saves it in a database.

```go
func (x *helloDaoPgx) RecordHello(ctx context.Context, hello *models.Hello) error {
	query := `
		INSERT INTO hellos (from, to, message)
		VALUES ($1, $2, $3) RETURNING id`

	err := db.Pool.QueryRow(ctx, query, hello.From, hello.To, hello.Message).Scan(&hello.ID)
	if err != nil {
		return fmt.Errorf("pgx error: %w", err)
	}
	return nil
}
```