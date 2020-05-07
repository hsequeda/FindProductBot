package main

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type User struct {
	Id       string
	Province string
}

type store struct {
	name    string
	rawName string
}

var provinces = map[string][]store{
	"Pinar del Rio": {
		{
			name:    "Pinar del Rio",
			rawName: "pinar",
		},
	},
	"Artemisa": {
		{
			name:    "Artemisa",
			rawName: "artemisa",
		},
	},
	"Mayabeque": {
		{
			name:    "Mayabeque",
			rawName: "mayabeque-tv",
		},
	},
	"Matanzas": {
		{
			name:    "Matanzas",
			rawName: "matanzas",
		},
	},
	"Cienfuegos": {
		{
			name:    "Cienfuegos",
			rawName: "cienfuegos",
		},
	},
	"Villa Clara": {
		{
			name:    "Villa Clara",
			rawName: "villaclara",
		},
	},
	"Sancti Spiritus": {
		{
			name:    "Sancti Spiritus",
			rawName: "sancti",
		},
	},
	"Ciego de Avila": {
		{
			name:    "Ciego de Avila",
			rawName: "ciego",
		},
	},
	"Camaguey": {
		{
			name:    "Camaguey",
			rawName: "camaguey",
		},
	},
	"Las Tunas": {
		{
			name:    "Las Tunas",
			rawName: "tunas",
		},
	},
	"Holguin": {
		{
			name:    "Holguin",
			rawName: "holguin",
		},
	},
	"Granma": {
		{
			name:    "Granma",
			rawName: "granma",
		},
	},
	"Santiago de Cuba": {
		{
			name:    "Santiago de Cuba",
			rawName: "santiago",
		},
	},
	"Guantanamo": {
		{
			name:    "Guantanamo",
			rawName: "guantanamo",
		},
	},
	"La Isla": {
		{
			name:    "La Isla",
			rawName: "isla",
		},
	},
	"La Habana": {
		{
			name:    "Carlos III",
			rawName: "carlos3",
		}, {
			name:    "Cuatro Caminos",
			rawName: "4caminos",
		}, {
			name:    "5ta y 42",
			rawName: "5taY42",
		}, {
			name:    "Pedregal",
			rawName: "tvpedregal",
		}, {
			name:    "Villa Diana",
			rawName: "caribehabana",
		}},
}

func InsertUser(id, province string) error {
	user := User{
		Id:       id,
		Province: province,
	}
	b, err := json.Marshal(user)
	if err != nil {
		return err
	}

	if err = dbWrite([]byte(id), b); err != nil {
		return err
	}
	return nil
}

func GetUser(id string) (*User, error) {
	logrus.Println(id)
	bb, err := dbRead(id)
	var u User
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(bb, &u); err != nil {
		return nil, err
	}
	return &u, nil
}
