package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	//"github.com/jackc/pgx/v5/internal/anynil"
	//"github.com/jackc/pgx/v5/internal/sanitize"
	//"github.com/jackc/pgx/v5/internal/stmtcache"
	"github.com/jackc/pgx/v5/pgconn"
	//"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5"
	//"github.com/jackc/pgx/v5/pgxpool"
)

//this interface contains pgx method, which i would use :D maybe
type client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOption pgx.TxOptions) (pgx.Tx, error)
	BeginTxFunc(ctx context.Context, txOption pgx.TxOptions, F func(pgx.Tx) error) error
}

//TODO: move username , password, host, port, name to separately .go file or struct
/*func newClient(ctx context.Context, maxAttempts int, username , password, host, port, name string) {
	// login:password@host:port/name
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", username , password, host, port, name)
	doWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		//pool, err := pgxpool.Connect(ctx, dsn)
		return nil
	}, maxAttempts, 5*time.Second)

}*/
//move to separately .go file(optional)
func doWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts < 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--

			continue
		}
		return nil
	}
	return
}

func M() {
	//urlExample := "postgres://username:password@localhost:5432/database_name"
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var name string
	var weight int64
	err = conn.QueryRow(context.Background(), "select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(name, weight)
}
