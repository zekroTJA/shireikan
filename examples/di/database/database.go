package database

import "fmt"

type Database interface {
	GetData() string
}

type TestDB struct {
	i int
}

func (db *TestDB) GetData() string {
	db.i++
	return fmt.Sprintf("Data %d, yay!", db.i)
}
