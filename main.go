package main

import (
  "Final-Project/config"
)

func main() {
  db := config.ConnectDataBase()
  sqlDB, _ := db.DB()
  defer sqlDB.Close()
}