// Copyright 2020 Alexey Yanchenko <mail@yanchenko.me>
//
// This file is part of the Gufo library.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gufodao

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	UID           string `gorm:"column:uid;type:varchar(60);UNIQUE;NOT NULL;"`
	Name          string `gorm:"column:name;type:varchar(60);NOT NULL;DEFAULT '';UNIQUE"`
	Pass          string `gorm:"column:pass;type:varchar(128);NOT NULL;DEFAULT ''"`
	Mail          string `gorm:"column:mail;type:varchar(254);DEFAULT '';UNIQUE"`
	Mailsent      int    `gorm:"column:mailsent;type:int;DEFAULT '0'"`
	Mailconfirmed int    `gorm:"column:mailconfirmed;:int;DEFAULT '0'"`
	Created       int    `gorm:"column:created;type:int;DEFAULT '0'"`
	Access        int    `gorm:"column:access;type:int;DEFAULT '0'"`
	Login         int    `gorm:"column:login;type:int;DEFAULT '0'"`
	Status        bool   `gorm:"column:status;type:bool;DEFAULT 'false'"`
	Completed     bool   `gorm:"column:completed;type:bool;DEFAULT 'false'"`
	IsAdmin       bool   `gorm:"column:is_admin;type:bool;DEFAULT 'false'"`
	Readonly      bool   `gorm:"column:readonly;type:bool;DEFAULT 'false'"`
}

type APITokens struct {
	gorm.Model
	TokenId    string `gorm:"column:tokenid;type:varchar(60);UNIQUE;NOT NULL;"`
	Token      string `gorm:"column:token;type:varchar(254);UNIQUE;NOT NULL;"`
	TokenName  string `gorm:"column:tokenname;type:varchar(60);DEFAULT '';"`
	UID        string `gorm:"column:uid;type:varchar(60);NOT NULL;"`
	Created    int    `gorm:"column:created;type:int;DEFAULT '0'"`
	Expiration int    `gorm:"column:expiration;type:int;DEFAULT '0'"`    //if 0 - no expiration time
	Status     bool   `gorm:"column:status;type:bool;DEFAULT 'true'"`    // if true - active, if false - deactivated
	IsAdmin    bool   `gorm:"column:is_admin;type:bool;DEFAULT 'false'"` //only if generated by admin
	Readonly   bool   `gorm:"column:readonly;type:bool;DEFAULT 'false'"`
	Comment    string `gorm:"column:comment;type:varchar(60);DEFAULT '';"`
}

type ImpersonateTokens struct {
	gorm.Model
	TokenId   string `gorm:"column:tokenid;type:varchar(60);UNIQUE;NOT NULL;"`
	Token     string `gorm:"column:token;type:varchar(254);UNIQUE;NOT NULL;"`
	UID       string `gorm:"column:uid;type:varchar(60);NOT NULL;"`
	Created   int    `gorm:"column:created;type:int;DEFAULT '0'"`
	CreatedBy string `gorm:"column:createdby;type:varchar(254);DEFAULT '';"`
}

type Entrypoint struct {
	gorm.Model
	ID      string `gorm:"column:entrypointid;type:varchar(60);NOT NULL;"`
	Status  bool   `gorm:"column:status;type:bool;DEFAULT 'false'"`
	Version string `gorm:"column:version;type:varchar(254);UNIQUE;NOT NULL;"`
}

/*
type Roles struct {
	gorm.Model
	UID   string `gorm:"type:varchar(60);UNIQUE;NOT NULL;"`
	Admin bool   `gorm:"type:double;DEFAULT 'false'"`
}


type Settings struct {
	gorm.Model
	Email_Confirmation bool `gorm:"type:double;DEFAULT 'false'"`
	Registration       bool `gorm:"type:double;DEFAULT 'false'"`
}
*/
/*
Timehash table structure:
Uid - users hash
email - users email
hash - 64 hash
param - Which function create this record. We need confirm email in signup and change current password
created - Where does record was created
livetime - hash life time
*/
type TimeHash struct {
	gorm.Model
	UID      string `gorm:"column:uid;type:varchar(60);NOT NULL;"`
	Mail     string `gorm:"column:mail;type:varchar(254);DEFAULT '';"`
	Hash     string `gorm:"column:hash;type:varchar(254);DEFAULT '';"`
	Param    string `gorm:"column:param;type:varchar(254);DEFAULT '';"`
	Created  int    `gorm:"column:created;type:int;DEFAULT '0'"`
	Livetime int    `gorm:"column:livetime;type:int;DEFAULT '0'"`
}

func CheckDBStructure() {
	//Check DB and table config
	db, err := ConnectDBv2()
	if err != nil {
		SetErrorLog("dbstructure.go:81: " + err.Error())
		//return "error with db"
	}

	dbtype := viper.GetString("database.type")

	/*
		if !db.Conn.Migrator().HasTable(&Roles{}) {
			//Create roles table
			db.Conn.Set("gorm:table_options", "ENGINE=InnoDB;").Migrator().CreateTable(&Roles{})
		}
	*/
	//Check if table users and roles exist
	if !db.Conn.Migrator().HasTable(&Users{}) {
		SetErrorLog("dbstructure.go:94: " + "Table users do not exist. Create table Users")
		//db.Conn.Debug().AutoMigrate(&Users{})
		//Create users table
		if dbtype == "mysql" {
			db.Conn.Set("gorm:table_options", "ENGINE=InnoDB;").Migrator().CreateTable(&Users{})
		} else {
			db.Conn.Migrator().CreateTable(&Users{})
		}

		//Add admin user
		//1. generate user hash
		userhash := Hashgen(8)
		//2. generate users Password
		userpass := RandomString(12)
		//2.1 generete pass passhash
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userpass), 8)
		if err != nil {

			SetErrorLog("dbstructure.go:108: " + err.Error())
		}

		//3. Admin User email
		useremail := viper.GetString("email.address")

		user := Users{
			UID:           userhash,
			Name:          "admin",
			Pass:          string(hashedPassword),
			Mail:          useremail,
			Mailsent:      int(time.Now().Unix()),
			Mailconfirmed: int(time.Now().Unix()),
			Created:       int(time.Now().Unix()),
			Status:        true,
			Completed:     true,
			IsAdmin:       true,
		}
		/*
			role := Roles{
				UID:   userhash,
				Admin: true,
			}
		*/
		db.Conn.Create(&user)
		//db.Conn.Create(&role)

		ans := fmt.Sprintf("Admin User created!\t\nname: admin\t\npass: %s\t\nemail: %s \t\n", userpass, useremail)

		//Check for email settings
		if viper.GetString("email.address") != "" {
			//Send email with password
			str := fmt.Sprintf("Your Gufo admin account was created with password: %v", userpass)
			linkarray := []string{str}
			ms := &MailSettings{}
			ms.Custom = false
			go SendHTMLEmail(user.Mail, "Hi, admin", linkarray, "New account", "email.html", nil, ms)
		}

		fmt.Printf(ans)

	}

	//Create timehash table
	if !db.Conn.Migrator().HasTable(&TimeHash{}) {
		if dbtype == "mysql" {
			db.Conn.Set("gorm:table_options", "ENGINE=InnoDB;").Migrator().CreateTable(&TimeHash{})
		} else {
			db.Conn.Migrator().CreateTable(&TimeHash{})
		}
	}

	//Create timehash table
	if !db.Conn.Migrator().HasTable(&Entrypoint{}) {
		if dbtype == "mysql" {
			db.Conn.Set("gorm:table_options", "ENGINE=InnoDB;").Migrator().CreateTable(&Entrypoint{})
		} else {
			db.Conn.Migrator().CreateTable(&Entrypoint{})
		}
	}

	/*
		if !db.Conn.Migrator().HasTable(&Settings{}) {
			//Create settings table
			db.Conn.Set("gorm:table_options", "ENGINE=InnoDB;").Migrator().CreateTable(&Settings{})
			setting := Settings{
				Email_Confirmation: false,
				Registration:       true,
			}
			db.Conn.Create(&setting)
		}
	*/
}