package pgsql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"DataBaseManage/public"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const (
	// Initialize connection constants.
	HOST     = "mypgserver-20170401.postgres.database.azure.com"
	DATABASE = "mypgsqldb"
	USER     = "mylogin@mypgserver-20170401"
	PASSWORD = "<server_admin_password>"
	// var connectionString string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require", HOST, USER, PASSWORD, DATABASE)
	// db, err := sql.Open("postgres", connectionString)
)

var _db *sql.DB

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}
func InitDB() {
	fmt.Println("Init DB ...")
	//curdir := public.GetCurDir()
	var err error
	_db, err := sql.Open("postgres", "user=postgres password=ap123 dbname=dbtest sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}
	if _db == nil {
		log.Fatal(err)
	}
	//fmt.Printf("nDB: %v\n", ODB)
	fmt.Println("testing db connection...")
	_db.SetMaxOpenConns(2000) //设置最大打开连接数
	_db.SetMaxIdleConns(1000) //设置最大空闲连接数

	err2 := _db.Ping()
	if err2 != nil {
		log.Fatalf("Error on opening database connection: %s", err2.Error())
	} else {
		fmt.Println("connection.success")
	}

	// Drop previous table of same name if one exists.
	_, err = _db.Exec("DROP TABLE IF EXISTS inventory;")
	checkError(err)
	fmt.Println("Finished dropping table (if existed)")

	// Create table.
	_, err = _db.Exec("CREATE TABLE inventory (id serial PRIMARY KEY, name VARCHAR(50), quantity INTEGER);")
	checkError(err)
	fmt.Println("Finished creating table")

	// Insert some data into table.
	sql_statement := "INSERT INTO inventory (name, quantity) VALUES ($1, $2);"
	_, err = _db.Exec(sql_statement, "banana", 150)
	checkError(err)
	_, err = _db.Exec(sql_statement, "orange", 154)
	checkError(err)
	_, err = _db.Exec(sql_statement, "apple", 100)
	checkError(err)
	fmt.Println("Inserted 3 rows of data")

	/*
	   避免错误操作，例如LOCK TABLE后用 INSERT会死锁，因为两个操作不是同一个连接，insert的连接没有table lock。
	   当需要连接，且连接池中没有可用连接时，新的连接就会被创建。
	   默认没有连接上限，你可以设置一个，但这可能会导致数据库产生错误“too many connections”
	   db.SetMaxIdleConns(N)设置最大空闲连接数
	   db.SetMaxOpenConns(N)设置最大打开连接数
	   长时间保持空闲连接可能会导致db timeout
	*/
}

func Getdb() *sql.DB {

	//public.Log("成功打开数据库文件")

	err2 := _db.Ping()
	if err2 != nil {
		log.Fatalf("Error on opening database connection: %s", err2.Error())
	}
	// Ping验证连接到数据库是否还活着，
	//必要时建立连接。
	return _db

}
func PrintCurrentPath() {

	dir, errer := filepath.Abs(filepath.Dir(os.Args[0]))
	if errer != nil {
		log.Fatal(errer)

	}
	fmt.Println(dir)
}

func GetSingle(sqlStr string) string {
	var id string
	public.Log(sqlStr)
	rows, _ := Getdb().Query(sqlStr)
	//defer db.Close()  .Scan(&id)
	i := 0
	defer rows.Close()
	//defer db.Close()
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for i, _ := range columns {
		valuePtrs[i] = &values[i]
	}

	for rows.Next() {
		//public.Log("has row")
		i++
		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatal(err)
			public.Log("not Single result")
			return ""
		}

		for i, _ := range columns {
			v := values[i]
			if v == nil {
				id = ""
			} else {
				switch t := v.(type) {
				default:
					str = v.(string)
				case bool:
					cell.Value = fmt.Sprintf("%s", v) //v
				case int:
					id = fmt.Sprintf("%d", v) //v
				case int64:
					id = fmt.Sprintf("%d", v) //v
				case int32:
					id = fmt.Sprintf("%d", v) //v
				case float64:
					id = fmt.Sprintf("%1.2f", v) //v
				case float32:
					id = fmt.Sprintf("%1.2f", v) //v
				case string:
					id = fmt.Sprintf("%s", v) //v
				case []byte:
					id = string(v.([]byte))
				case time.Time:
					id = v.(string)
				}
			}

		}
	}
	if i == 0 {
		//public.Log("has no row")
		return ""
	}

	public.Log(id)
	return id
}

func ExecuteUpdate(sqlStr string) int {
	//public.Log("sqlStr=" + sqlStr)
	res, err := Getdb().Exec(sqlStr)
	if err != nil {
		fmt.Println("exec sql failed:", err.Error()+sqlStr)
		return 0
	}
	if strings.Contains(strings.ToLower(sqlStr), "insert") {

		rowId, err := res.LastInsertId()
		if err != nil {
			fmt.Println("fetch last insert id failed:", err.Error())
			return 0
		}

		//fmt.Println("lastInsertId=")
		//fmt.Println(rowId)
		//str := strconv.Itoa(lastInsertId)
		str := strconv.FormatInt(rowId, 10)
		ret, _ := strconv.Atoi(str)
		return ret
	} else {
		rowId, err := res.RowsAffected()
		if err != nil {
			fmt.Println("fetch RowsAffected failed:", err.Error())
			return 0
		}
		//fmt.Println("update idddd=")
		//fmt.Println(rowId)
		str := strconv.FormatInt(rowId, 10)
		ret, _ := strconv.Atoi(str)
		return ret
	}
	//return 1
}

func ExecuteQuery(sqlStr string) map[int]map[string]string {
	db := Getdb()
	rows, err := db.Query(sqlStr)

	public.Log("sql =" + sqlStr)
	if err != nil {
		CheckErr(err)
		return nil
	}
	defer rows.Close()
	//defer db.Close()
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	final_result := make(map[int]map[string]string)
	result_id := 0
	for i, _ := range columns {
		valuePtrs[i] = &values[i]
	}

	for rows.Next() {
		rows.Scan(valuePtrs...)
		m := make(map[string]string)
		//Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		for i, col := range columns {
			v := values[i]
			//tmp_struct[col] =  fmt.Sprintf("%s", v)
			//if columns[i] == "imagename" {
			//public.Log(columns[i], ": ", v)
			//}

			//if result_id == 0 {
			//public.Log(reflect.TypeOf(v))
			//}

			if v == nil {
				m[strings.ToLower(col)] = ""
			} else {
				switch t := v.(type) {
				default:
					m[strings.ToLower(col)] = fmt.Sprintf("%s", v) //v
				case bool:
					m[strings.ToLower(col)] = fmt.Sprintf("%s", v) //v
				case int:
					m[strings.ToLower(col)] = fmt.Sprintf("%d", v) //v
				case int64:
					m[strings.ToLower(col)] = fmt.Sprintf("%d", v) //v
				case int32:
					m[strings.ToLower(col)] = fmt.Sprintf("%d", v) //v
				case float64:
					m[strings.ToLower(col)] = fmt.Sprintf("%1.2f", v) //v
				case float32:
					m[strings.ToLower(col)] = fmt.Sprintf("%1.2f", v) //v
				case string:
					m[strings.ToLower(col)] = fmt.Sprintf("%s", v) //v
				case []byte:
					m[strings.ToLower(col)] = string(v.([]byte))
				case time.Time:
					m[strings.ToLower(col)] = v.(string)
				}
			}

		}
		//fmt.Print(m)
		final_result[result_id] = m
		result_id++
	}
	return final_result

}
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

/*
delete  FROM bom;
delete  FROM producttype;
delete  FROM producttypeTwo;
delete  FROM product;
delete  FROM invoicing;
delete  FROM location;
delete  FROM oafile;
delete  FROM oa_approval;
delete  FROM oa_askforleave;
delete  FROM oa_location;
delete  FROM oa_news;
delete  FROM oa_registration;
delete  FROM oa_workreport;
delete  FROM orders;
delete  FROM orderdesk;
delete  FROM supplier;
delete  FROM warehouse;
delete  FROM users;
delete  FROM users;
delete  FROM folder;
delete  FROM Recharge;
delete  FROM folder;
;
;
*/
