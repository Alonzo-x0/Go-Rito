package lib
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type Tag struct {
	ID int `json:ID`
}
//"SELECT Lkey from weather where userid='1234'"
func SelectRows(db *sql.DB, tarCol string, table string, desCol string, value string ) (int, error){
	results, err := db.Query("SELECT " + tarCol + " from " + table + " where " + desCol+ "='" + value +"'")
	if err != nil {
		return 0, err
	}
	var x int
	for results.Next() {
		var tag Tag

		err = results.Scan(&tag.ID)
		if err != nil {
			return 0, err
		}
		x = tag.ID
	}

	defer results.Close()
	return x, err
}

//delete duplicate rows leaving ID to the be lowest int
func DeleteDupes(db *sql.DB, table string, col1 string, col2 string) (error){
	quer := "delete foo from " + table + " foo inner join " + table + " bar on bar." + col1 + " = foo." + col1 + " and bar." + col2 + "< foo." + col2
	log.Println(quer)
	_, err := db.Exec(quer)

	if err != nil {
		return err
	}

	return err
}

func Insert(db *sql.DB, table string, col1 string, col2 string, userid string, Lkey string) (error){
	_, err := db.Exec("insert into " + table + "(" + col1 + ", " + col2 + ") values('" + userid + "', " + Lkey + ")")
	if err != nil {
		return err
	}
	return err 
}
