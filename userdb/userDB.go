package main
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type Tag struct {
	ID int `json:ID`
}
//"SELECT Lkey from weather where userid='1234'"
func SelectRows(db *sql.DB, tarCol string, database string, desCol string, value string ) (error){
	results, err := db.Query("SELECT " + tarCol + " from " + database + " where " + desCol+ "='" + value +"'")
	if err != nil {
		log.Println(err)
	}
	var x int
	for results.Next() {
		var tag Tag

		err = results.Scan(&tag.ID)
		if err != nil {
			log.Println(err)
			return err
		}
		x = tag.ID
	}
	log.Println(x)

	defer results.Close()
	return err
}

//delete duplicate rows leaving ID to the be lowest int
func DeleteDupes(db *sql.DB, table string, col1 string, ID string) (error){
	_, err := db.Exec("delete foo from " + table + " foo inner join " + table + " bar on bar." + col1 + " = foo." + col1 + " and bar." + ID + "< foo." + ID)

	if err != nil {
		return err
	}

	return err
}

func Insert(db *sql.DB, database string, col1 string, col2 string, userid string, Lkey string) (error){
	_, err := db.Exec("insert into " + database + "(" + col1 + ", " + col2 + ") values('" + userid + "', " + Lkey + ")")
	if err != nil {
		return err
	}
	return err 
}

func main() {
	db, err := sql.Open("mysql", "killer:toor@tcp(127.0.0.1:3306)/discord")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()


	err = Insert(db, "weather", "userid", "lkey", "157959951501885440", "335315")

	if err != nil{
		log.Println(err)
		return
	}
	//"SELECT Lkey from weather where userid='1234'"
	//SelectRows(db, "Lkey", "weather", "userid", "1234")

	


}