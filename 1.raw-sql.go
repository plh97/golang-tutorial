
import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main1() {
	dsn := "root:123456@tcp(127.0.0.1:3305)/gorm_db_new"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("dsn incorrect %s", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("db connect failed: %s", err)
	}

	//res, err := db.Exec("INSERT into users (name, age, email) values ('json', 18, 'json@gmail.com'), ('arron', 28, 'arron@gmail.com')")
	//if err != nil {
	//	log.Fatalf("insert errror: %s", err)
	//}
	//fmt.Println(res.LastInsertId())
	rows, err := db.Query("select id, name from users")
	if err != nil {
		log.Fatalf("database query failed: %s", err)
	}
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		fmt.Println(id, name, err)
	}

	var id int
	var name string
	err = db.QueryRow("select id, name from users").Scan(&id, &name)
	fmt.Println(id, name, err)
}
