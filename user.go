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

type province struct {
	name   string
	stores []store
}

var provinces = map[string]province{
	"pr": {
		name: "Pinar del Rio",
		stores: []store{{
			name:    "Pinar del Rio",
			rawName: "pinar",
		}},
	},
	"ar": {
		name: "Artemisa",
		stores: []store{{
			name:    "Artemisa",
			rawName: "artemisa",
		}},
	},
	"my": {
		name: "Mayabeque",
		stores: []store{{
			name:    "Mayabeque",
			rawName: "mayabeque-tv",
		}},
	},
	"mt": {
		name: "Matanzas",
		stores: []store{{
			name:    "Matanzas",
			rawName: "matanzas",
		}},
	},
	"cf": {
		name: "Cienfuegos",
		stores: []store{{
			name:    "Cienfuegos",
			rawName: "cienfuegos",
		}},
	},
	"vc": {
		name: "Villa Clara",
		stores: []store{{
			name:    "Villa Clara",
			rawName: "villaclara",
		}},
	},
	"ss": {
		name: "Sancti Spiritus",
		stores: []store{{
			name:    "Sancti Spiritus",
			rawName: "sancti",
		}},
	},
	"ca": {
		name: "Ciego de Avila",
		stores: []store{{
			name:    "Ciego de Avila",
			rawName: "ciego",
		}},
	},
	"cm": {
		name: "Camaguey",
		stores: []store{{
			name:    "Camaguey",
			rawName: "camaguey",
		}},
	},
	"lt": {
		name: "Las Tunas",
		stores: []store{{
			name:    "Las Tunas",
			rawName: "tunas",
		}},
	},
	"hg": {
		name: "Holguin",
		stores: []store{{
			name:    "Holguin",
			rawName: "holguin",
		}},
	},
	"gr": {
		name: "Granma",
		stores: []store{{
			name:    "Granma",
			rawName: "granma",
		}},
	},
	"st": {
		name: "Santiago de Cuba",
		stores: []store{{
			name:    "Santiago de Cuba",
			rawName: "santiago",
		}},
	},
	"gt": {
		name: "Guantanamo",
		stores: []store{{
			name:    "Guantanamo",
			rawName: "guantanamo",
		}},
	},
	"ij": {
		name: "La Isla",
		stores: []store{{
			name:    "La Isla",
			rawName: "isla",
		}},
	},
	"lh": {
		name: "La Habana",
		stores: []store{{
			name:    "Carlos III",
			rawName: "carlos3",
		}, {
			name:    "Cuatro Caminos",
			rawName: "4caminos",
		}},
	},
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
	logrus.Info("find ",string(bb))
	if err = json.Unmarshal(bb, &u); err != nil {
		return nil, err
	}
	return &u, nil
}
