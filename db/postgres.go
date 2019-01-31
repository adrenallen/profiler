package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresConn struct {
	user              string
	password          string
	host              string
	port              int
	defaultDB         string //required for create db call
	additionalOptions map[string]string
}

//Creates a new postgres connection object
func NewPostgresConn(user string, password string, host string, port int, defaultDB string, additionalOptions map[string]string) PostgresConn {
	return PostgresConn{
		user:              user,
		password:          password,
		host:              host,
		port:              port,
		defaultDB:         defaultDB,
		additionalOptions: additionalOptions,
	}
}

//Connect to specified database
func (p PostgresConn) GetConnectionToDatabase(dbName string) (*sql.DB, error) {
	connString := fmt.Sprintf(`user=%s password=%s host=%s port=%d dbname=%s`,
		p.user,
		p.password,
		p.host,
		p.port,
		dbName)

	for key, val := range p.additionalOptions {
		connString = fmt.Sprintf(`%s %s=%s`, connString, key, val)
	}

	return sql.Open(`postgres`, connString)
}

//Connect to default database
func (p PostgresConn) GetConnection() (*sql.DB, error) {
	return p.GetConnectionToDatabase(p.defaultDB)
}

func (p PostgresConn) CheckAndCreateDB(dbName string) error {
	conn, err := p.GetConnection()
	if err != nil {
		return err
	}

	exists, _ := p.dbExists(dbName)
	//if db exists then just return nil we are done
	if exists {
		return nil
	}

	//if we get here then try to create and return the error
	_, err = conn.Exec(fmt.Sprintf(`CREATE DATABASE %s`, dbName))
	return err
}

func (p PostgresConn) dbExists(dbName string) (bool, error) {
	conn, err := p.GetConnection()
	if err != nil {
		return false, err
	}
	row := conn.QueryRow(
		`SELECT datname FROM pg_catalog.pg_database WHERE datname = $1;`,
		dbName,
	)

	var name string
	err = row.Scan(&name)
	if err != nil {
		return false, err
	}

	return name == dbName, nil
}
