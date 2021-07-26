package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

//middleware
func Cors(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "OPTIONS" {
		return
	}
	next(w, r)
}

//End middelware

// connection
func GetConnection() *sql.DB {
	var db *sql.DB
	var err error
	db, err = sql.Open("sqlite3", "./tibiawiki.db")
	if err != nil {
		panic(err)
	}
	return db
}

//end connection

//Api
type CreatureApi struct {
}

func NewCreatureApi() *CreatureApi {
	return &CreatureApi{}
}

type Creature struct {
	CreatureID int     `json:"creature_id,omitempty"`
	Name       string  `json:"name"`
	Drops      *[]Item `json:"drops"`
}
type Item struct {
	ItemID    int     `json:"item_id,omitempty"`
	Name      string  `json:"name"`
	SellPrice int     `json:"sell_price"`
	Chance    float32 `json:"chance"`
	MinCount  int8    `json:"min_count"`
	MaxCount  int8    `json:"max_count"`
}

func (n *CreatureApi) GetTable(world string, npc string) ([]*Creature, error) {
	db := GetConnection()

	//items filter query
	var npcCode int
	filterItemList := []int{}
	filterQuery := `SELECT
	item_id
	FROM npc_offer_buy WHERE npc_id=?`

	if npc == "Yasir" {
		npcCode = 56612
	}
	if npc == "Telas" {
		npcCode = 37031
	} else {
		npcCode = 56612
	}
	rows, err := db.Query(filterQuery, npcCode)
	if err != nil {
		return []*Creature{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var item int
		rows.Scan(&item)
		filterItemList = append(filterItemList, item)
	}
	//fin parte querry list

	// begin Querry creatures
	q := `SELECT
        article_id, name
                        FROM creature`
	rows, err = db.Query(q)
	if err != nil {
		return []*Creature{}, err
	}
	defer rows.Close()

	creatures := []*Creature{}

	for rows.Next() {
		creature := Creature{}
		rows.Scan(&creature.CreatureID, &creature.Name)
		creatures = append(creatures, &creature)
	}

	//Drop queries

	w := `SELECT
        item_id, chance, min, max
        FROM creature_drop WHERE creature_id=?`
	e := `SELECT
        name, value_sell
        FROM item WHERE article_id=?`

	for _, element := range creatures {
		// Ejecutamos la query
		items := []Item{}

		itemRows, err := db.Query(w, element.CreatureID)
		if err != nil {
			continue
		}

		for itemRows.Next() {
			item := Item{}
			err = itemRows.Scan(&item.ItemID, &item.Chance, &item.MinCount, &item.MaxCount)
			if err != nil {
				continue
			}

			if !contains(filterItemList, item.ItemID) {
				continue
			}

			//Complete with item information
			err = db.QueryRow(e, item.ItemID).Scan(&item.Name, &item.SellPrice)
			if err != nil {
				continue
			}

			items = append(items, item)
		}

		element.Drops = &items
		defer itemRows.Close()

	}

	return creatures, nil
}

//End api

//Main
func main() {

	r := mux.NewRouter()

	n := negroni.New(
		negroni.HandlerFunc(Cors),
		negroni.NewLogger(),
	)
	//create handlers
	http.Handle("/", r)
	CreaturesHandler(r, *n)

	//loggers and serve
	logger := log.New(os.Stderr, "logger: ", log.Lshortfile)
	srv := &http.Server{
		ReadTimeout:  180 * time.Second,
		WriteTimeout: 180 * time.Second,
		Addr:         ":" + strconv.Itoa(8085),
		Handler:      context.ClearHandler(http.DefaultServeMux),
		ErrorLog:     logger,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func GetCreaturesRanking() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Check world
		vars := mux.Vars(r)
		world, errPath := vars["world"]
		if errPath == false {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("world not assigned"))
			return
		}

		world = strings.Title(strings.ToLower(world))
		fmt.Println(world) //sacar

		//Check NPC
		npc := r.URL.Query().Get("npc")
		if npc == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("npc not assigned"))
			return
		} else {
			npc = strings.Title(strings.ToLower(npc))
		}
		fmt.Println(npc) //sacar

		n := NewCreatureApi()
		creatures, err := n.GetTable(world, npc)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		// Conviertiendo el slice de Note a formato JSON, retorna un []byte
		j, err := json.Marshal(creatures)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(j)

	})
}
func CreaturesHandler(r *mux.Router, n negroni.Negroni) {
	r.Handle("/creatures/{world}", n.With(
		negroni.Wrap(GetCreaturesRanking()),
	)).Methods("GET", "OPTIONS").Name("getCreaturesRanking")

}

func contains(s []int, str int) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
