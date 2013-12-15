package main

import (
  "path"
  "os/user"
  "strconv"
  "github.com/syndtr/goleveldb/leveldb"
  "fmt"
)

type Database struct {
  db        *leveldb.DB
}

func (db *Database) Bootstrap() error {

  lastBlock := db.LastBlock();
  if lastBlock == 0 {
    fmt.Println("No last block. Setting genesis")
    // If no last block is set, set it to the genesis block defined in server.go
    s := strconv.FormatUint(uint64(GenesisBlock), 10)
    err := db.db.Put([]byte("lastBlock"), []byte(s), nil)

    if err != nil {
      return err
    }
  }

  fmt.Println("Last block", lastBlock)

  return nil
}

func (db *Database) Close() {
  defer db.db.Close()
}

func (db *Database) LastBlock() uint32 {
  data, _ := db.db.Get([]byte("lastBlock"), nil)
  if len(data) == 0 {
    return 0
  }

  lastBlock, _ := strconv.ParseUint(string(data), 10, 32)

  return uint32(lastBlock)
}

func NewDatabase() (*Database, error) {
  usr, _ := user.Current()
  dbPath := path.Join(usr.HomeDir, ".mastercoin", "database")

  // Open the db
  db, err := leveldb.OpenFile(dbPath, nil)
  if err != nil {
    return nil, err
  }

  database := &Database{db: db}

  // Bootstrap database. Sets a few defaults; such as the last block
  database.Bootstrap()

  return database, nil
}
