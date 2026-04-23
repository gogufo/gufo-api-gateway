// Copyright 2019-2025 Alexey Yanchenko <mail@yanchenko.me>
//
// This file is part of the Gufo library.
//
// Licensed under the Business Source License 1.1 (the "License");
// you may not use this file except in compliance with the License.
//
// You may obtain a copy of the License in the LICENSE file at the root of this repository.
//
// As of the Change Date specified in that file, in accordance with the Business Source
// License, use of this software will be governed by the Apache License, Version 2.0.
//
// THIS SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
// PURPOSE AND NON-INFRINGEMENT.

package gufodao

import (
	"fmt"

	//	"gorm.io/driver/sqlite"
	viper "github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB struct
//Uncomment after remove v1 Librry

type DBv2 struct {
	Conn *gorm.DB
}

// Connection instance
var DBConnectionv2 = &DBv2{}

func DBConnectv2() (*DBv2, error) {
	dbtype := viper.GetString("database.type")
	user := viper.GetString("database.user")
	pass := viper.GetString("database.password")
	pass = DecryptConfigPasswords(pass)
	dbname := viper.GetString("database.dbname")
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	charset := viper.GetString("database.charset")
	sslmode := viper.GetString("database.sslmode")

	var request string

	var err error

	switch dbtype {
	case "mysql":
		request = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true", user, pass, host, port, dbname, charset)
		db, err := gorm.Open(mysql.Open(request), &gorm.Config{})
		DBConnectionv2.Conn = db
		return DBConnectionv2, err
	case "postgres":
		request = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbname, pass, sslmode)
		//	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
		db, err := gorm.Open(postgres.Open(request), &gorm.Config{})
		DBConnectionv2.Conn = db
		return DBConnectionv2, err
	default:
		return nil, err
	}

}

func DBCheck() bool {
	_, err := DBConnectv2()
	//defer db.Close()
	if err != nil {
		SetErrorLog("dbconnectv2.go:77: " + err.Error())
		return false
	} else {
		return true
	}
}

// GetConnection - connect to DB
func ConnectDBv2() (*DBv2, error) {
	db, err := DBConnectv2()
	if err != nil {
		return db, err
	}

	sqlDB, err := db.Conn.DB()
	if err != nil {
		return db, err
	}

	dbcon := viper.GetInt("database.connectionssize")
	dbpool := viper.GetInt("database.poolsize")

	sqlDB.SetMaxIdleConns(dbcon)
	sqlDB.SetMaxOpenConns(dbpool)

	return db, nil
}

func DBConnect(prefix string) (*gorm.DB, error) {

	dbtype := viper.GetString(prefix + ".type")
	user := viper.GetString(prefix + ".user")
	pass := DecryptConfigPasswords(viper.GetString(prefix + ".password"))
	dbname := viper.GetString(prefix + ".dbname")
	host := viper.GetString(prefix + ".host")
	port := viper.GetString(prefix + ".port")
	charset := viper.GetString(prefix + ".charset")
	sslmode := viper.GetString(prefix + ".sslmode")

	var db *gorm.DB
	var err error

	switch dbtype {

	case "mysql":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
			user, pass, host, port, dbname, charset,
		)

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			host, port, user, dbname, pass, sslmode,
		)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	default:
		return nil, fmt.Errorf("unsupported db type: %s", dbtype)
	}

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(viper.GetInt(prefix + ".connectionssize"))
	sqlDB.SetMaxOpenConns(viper.GetInt(prefix + ".poolsize"))

	return db, nil
}
