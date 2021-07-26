package creatures

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func GetConnection() *sql.DB {
	var db *sql.DB
	var err error
	db, err = sql.Open("sqlite3", "./tibiawiki.db")
	fmt.Println("aqusss")
	if err != nil {
		panic(err)
	}
	return db
}

type Creature struct {
	CreatureID int    `json:"creature_id,omitempty"`
	Name       string `json:"name"`
	Drops      []Item `json:"drops"`
}
type Item struct {
	ItemID    int     `json:"item_id,omitempty"`
	Name      string  `json:"name"`
	SellPrice int     `json:"sell_price"`
	Chance    float32 `json:"chance"`
	MinCount  int8    `json:"min_count"`
	MaxCount  int8    `json:"max_count"`
}

func (n *Creature) GetAll() ([]Creature, error) {
	db := GetConnection()
	q := `SELECT
        article_id, name
                        FROM creature`
	// Ejecutamos la query
	rows, err := db.Query(q)
	if err != nil {
		return []Creature{}, err
	}
	// Cerramos el recurso
	defer rows.Close()

	creatures := []Creature{}

	for rows.Next() {
		rows.Scan(&n.CreatureID, &n.Name)
		creatures = append(creatures, *n)
	}

	//Drop queries

	w := `SELECT
        item_id, chance, min, max
        FROM creature_drop WHERE creature_id=?`
	e := `SELECT
        name, value
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
			itemRows.Scan(&item.ItemID, &item.Chance, &item.MinCount, &item.MaxCount)
			//Complete with item information
			err = db.QueryRow(e, item.ItemID).Scan(&item.Name, &item.SellPrice)
			if err != nil {
				continue
			}
			items = append(items, item)
		}

		defer itemRows.Close()
		element.Drops = items

	}

	return creatures, nil
}

/*
func (n *Note) GetByID(id int) (Note, error) {
        db := GetConnection()
        q := `SELECT
                id, title, description, created_at, updated_at
                FROM notes WHERE id=?`

        err := db.QueryRow(q, id).Scan(
                &n.ID, &n.Title, &n.Description, &n.CreatedAt, &n.UpdatedAt,
        )
        if err != nil {
                return Note{}, err
        }

        return *n, nil
}

func (n Note) Update() error {
        db := GetConnection()
        q := `UPDATE notes set title=?, description=?, updated_at=?
                WHERE id=?`
        stmt, err := db.Prepare(q)
        if err != nil {
                return err
        }
        defer stmt.Close()

        r, err := stmt.Exec(n.Title, n.Description, time.Now(), n.ID)
        if err != nil {
                return err
        }
        if i, err := r.RowsAffected(); err != nil || i != 1 {
                return errors.New("ERROR: Se esperaba una fila afectada")
        }
        return nil
}

func (n Note) Delete(id int) error {
        db := GetConnection()

        q := `DELETE FROM notes
                WHERE id=?`
        stmt, err := db.Prepare(q)
        if err != nil {
                return err
        }
        defer stmt.Close()

        r, err := stmt.Exec(id)
        if err != nil {
                return err
        }
        if i, err := r.RowsAffected(); err != nil || i != 1 {
                return errors.New("ERROR: Se esperaba una fila afectada")
        }
        return nil
}
*/
