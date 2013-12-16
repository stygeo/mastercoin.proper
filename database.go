package main

import (
  "path"
  "os/user"
  "encoding/binary"
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

    buf := make([]byte, 8) // 64 bit buffer
    binary.PutUvarint(buf, uint64(GenesisBlock))

    err := db.db.Put([]byte("lastBlock"), buf, nil)

    if err != nil {
      db.Close()

      panic("Bootstrapping database failed. Unable to set lastBlock")
    }
  }

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

  lastBlock, _ := binary.Uvarint(data)

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
