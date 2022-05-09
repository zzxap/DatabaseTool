package dal

/*

func heartBeating(a *agent, bytes chan byte, timeout int) {
	var after <-chan time.Time
loop:
	after = time.After(time.Duration(timeout) * time.Second)
	for {
		select {
		case <-bytes:
			a.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
			goto loop
		case <-after:
			a.Destroy()
			close(bytes)
			goto end
			return
		}
	}
end:
	return
}
*/
import (
	"DataBaseManage/public"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

//Go 设计模式：https://blogtitle.github.io/some-useful-patterns/
var _db *sql.DB

var PthSep string
var ISmac int

var Dbtype = "postgres"
var Dbhost = ""
var Dbport = "5432"
var Dbuser = "postgres"
var Dbpassword = ""
var Dbname = ""

var SSHHost = ""
var SSHPort = ""
var SSHUser = ""
var SSHPass = ""

var SqliteFileName = ""

var oSingle sync.Once

var logsql = false

func GetTimeId() int64 {

	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetTimeIdStr() string {
	return strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)

}

func InitDB() (bool, string) {
	public.Log("InitDB")
	var testdb *sql.DB

	var err error

	public.Log(" Init DB type=" + Dbtype + " SqliteFileName=" + SqliteFileName)
	connectstr := ""
	if Dbtype == "postgres" {

		connectstr = fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", Dbhost, Dbport, Dbuser, Dbpassword)
		//connectstr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", Dbhost, Dbport, Dbuser, Dbpassword, Dbname)

		testdb, err = sql.Open("postgres", connectstr)

	} else if Dbtype == "mysql" {
		connectstr = "" + Dbuser + ":" + Dbpassword + "@tcp(" + Dbhost + ":" + Dbport + ")/?charset=utf8"
		testdb, err = sql.Open("mysql", connectstr)

	} else if Dbtype == "mssql" {

		connectstr = fmt.Sprintf("server=%s;port%s;user id=%s;password=%s", Dbhost, Dbport, Dbuser, Dbpassword)
		testdb, err = sql.Open("mssql", connectstr)

	} else if Dbtype == "etcd" {
		if InitETCD() {
			return true, "success"
		} else {
			return false, "fail"
		}

	} else if Dbtype == "mongodb" {
		if InitMongoDB() {
			return true, "success"
		} else {
			return false, "fail"
		}

	} else {
		testdb = initSqliteDb()
	}
	if err != nil {
		public.Log(err)
		return false, err.Error()
	}
	if testdb == nil {
		public.Log("Init DB fail")
		return false, "Init DB fail"
	}
	//fmt.Printf("nDB: %v\n", ODB)
	public.Log("testing db connection...")

	//defer _db.Close()
	err2 := testdb.Ping()
	public.Log("ping..." + Dbtype)
	if err2 != nil {
		fmt.Printf("Error on opening database connection: %s", err2.Error())
		return false, "Error on opening database connection " + err2.Error()
	} else {
		public.Log("connection.success")
		_db = testdb
		_db.SetMaxOpenConns(2000) //设置最大打开连接数
		_db.SetMaxIdleConns(100)  //设置最大空闲连接数

		return true, "success"
	}

}

type ViaSSHDialer struct {
	client *ssh.Client
}

func (self *ViaSSHDialer) Dial(network, address string) (net.Conn, error) {
	return self.client.Dial(network, address)
}
func (self *ViaSSHDialer) DialMysql(addr string) (net.Conn, error) {
	return self.client.Dial("tcp", addr)
}
func (self *ViaSSHDialer) Open(s string) (_ driver.Conn, err error) {
	return pq.DialOpen(self, s)
}

func (self *ViaSSHDialer) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	return self.client.Dial(network, address)
}

func InitDBBySSH() (bool, string) {
	public.Log("InitDBBySSH")
	sshHost := SSHHost // SSH Server Hostname/IP
	sshPort := SSHPort // SSH Port
	sshUser := SSHUser // SSH Username
	sshPass := SSHPass // Empty string for no password
	//dbUser := Dbuser                // DB username
	//dbPass := Dbpassword            // DB Password
	//dbHost := Dbhost + ":" + Dbport // DB Hostname/IP
	//dbName := Dbname                // Database name

	var agentClient agent.Agent
	// Establish a connection to the local ssh-agent
	if conn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		defer conn.Close()

		// Create a new instance of the ssh agent
		agentClient = agent.NewClient(conn)
	}

	// The client configuration with configuration option to use the ssh-agent
	sshConfig := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{},
	}

	// When the agentClient connection succeeded, add them as AuthMethod
	if agentClient != nil {
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeysCallback(agentClient.Signers))
	}
	// When there's a non empty password add the password AuthMethod
	if sshPass != "" {
		sshConfig.Auth = append(sshConfig.Auth, ssh.PasswordCallback(func() (string, error) {
			return sshPass, nil
		}))
	}

	// Connect to the SSH Server
	if sshcon, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshHost, sshPort), sshConfig); err == nil {

		if Dbtype == "mysql" {
			mysql.RegisterDial("mysql+tcp", (&ViaSSHDialer{sshcon}).DialMysql)
		} else {
			sql.Register("postgres+ssh", &ViaSSHDialer{sshcon})
		}

		//defer sshcon.Close()
		if len(Dbname) > 0 {
			return InitDB()
		} else {
			return ReInitDB(Dbtype, Dbhost, Dbport, Dbuser, Dbpassword, Dbname)
		}

	} else {
		return true, "fail"
	}
}

func initSqliteDb() *sql.DB {
	if SqliteFileName == "" {
		saveSqliteFileName := public.GetiniValueByKey("SqliteFileName")

		if saveSqliteFileName != "" {
			public.Log("saveSqliteFileName=" + saveSqliteFileName)
			SqliteFileName = saveSqliteFileName

		} else {
			if ISmac == 1 {
				public.Log("is mac")
				SqliteFileName = "InvoicingMac.sqlite"
				public.SetKeyValue("SqliteFileName", "InvoicingMac.sqlite")
			} else {
				public.Log("is not mac")
				SqliteFileName = "db.sqlite"
				public.SetKeyValue("SqliteFileName", "db.sqlite")
			}
		}
		public.Log("111")
		public.Log("SqliteFileNameeee=" + saveSqliteFileName)
	}
	if SqliteFileName == "" {
		return nil
	}
	public.Log("init SqliteFileName=" + SqliteFileName)

	var err error
	curdir := public.GetCurDir()
	PthSep = string(os.PathSeparator)
	dbpath := curdir + PthSep + "db" + PthSep + SqliteFileName

	public.Log("Init DB sqlite ..." + SqliteFileName)
	if !public.ExistsPath(dbpath) {
		public.Log("db not exists" + dbpath)
	}
	var testdb *sql.DB
	if public.ExistsPath(dbpath) {
		public.Log(dbpath + "  存在")
		testdb, err = sql.Open("sqlite3", dbpath)
		if err == nil {
			return testdb
		}
	} else {
		public.Log("db not exists" + dbpath)
	}
	return nil
}
func ReInitDB(dbtype, dbhost, dbport, dbuser, dbpassword, dbname string) (bool, string) {
	if len(SSHHost) > 0 && len(SSHPort) > 0 && len(SSHUser) > 0 && len(SSHPass) > 0 {
		return InitDBBySSH()
	}
	var testdb *sql.DB

	var err error

	public.Log("Re Init DB type=" + dbtype + " SqliteFileName=" + SqliteFileName)
	if dbtype == "postgres" {

		psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbhost, Dbport, dbuser, dbpassword, dbname)
		testdb, err = sql.Open("postgres", psqlInfo)

	} else if Dbtype == "mysql" {

		testdb, err = sql.Open("mysql", ""+dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8")

	} else if Dbtype == "mssql" {

		connString := fmt.Sprintf("server=%s;port%s;database=%s;user id=%s;password=%s", dbhost, dbport, dbname, dbuser, dbpassword)
		testdb, err = sql.Open("mssql", connString)

	} else {
		testdb = initSqliteDb()
	}
	if err != nil {
		public.Log(err)
		return false, err.Error()
	}
	if testdb == nil {
		public.Log("Init DB fail")
		return false, "Init DB fail"
	}
	//fmt.Printf("nDB: %v\n", ODB)
	public.Log("testing db connection...")

	//defer _db.Close()
	err2 := testdb.Ping()
	public.Log("ping..." + Dbtype)
	if err2 != nil {
		fmt.Printf("Error on opening database connection: %s", err2.Error())
		return false, "Error on opening database connection " + err2.Error()
	} else {
		public.Log("connection.success")
		_db = testdb
		_db.SetMaxOpenConns(2000) //设置最大打开连接数
		_db.SetMaxIdleConns(100)  //设置最大空闲连接数

		return true, "success"
	}

}

func CreateDb() bool {
	var testdb *sql.DB
	public.Log("Create Db ...")

	var err error
	//Dbhost = public.GetDBServer()
	fmt.Printf("Create Db ..." + Dbhost + "\n")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", Dbhost, Dbport, Dbuser, Dbpassword)
	testdb, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		public.Log(err)
		return false
	}
	if testdb == nil {
		public.Log("Init DB fail")
		return false
	}
	//fmt.Printf("nDB: %v\n", ODB)
	public.Log("testing db connection...")

	err2 := testdb.Ping()
	public.Log("ping...")
	if err2 != nil {
		fmt.Printf("Error on opening database connection: %s", err2.Error())
		return false
	} else {

		public.Log("connection.success")

		num := GetSingleInDb(testdb, "SELECT count(*) as num FROM pg_catalog.pg_database WHERE lower(datname) = lower('"+Dbname+"')")
		if num == "0" {

			ExecuteUpdateInDb(testdb, "create database "+Dbname+" OWNER "+Dbuser)
			num = GetSingleInDb(testdb, "SELECT count(*) as num FROM pg_catalog.pg_database WHERE lower(datname) = lower('"+Dbname+"')")
			public.Log("count=" + num)
			if num == "1" {
				public.Log("create database " + Dbname + " success")

				return true

			} else {
				public.Log("create database " + Dbname + " fail")
			}
		} else {

			public.Log("database " + Dbname + " exists")

			return true
		}
		return false
	}
}

func CreateTable() {

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
	public.Log(dir)
}
func GetSingleJson(sqlStr string) string {

	return GetSingleJsonInDb(Getdb(), sqlStr)

}

func GetSingle(sqlStr string) string {

	return GetSingleInDb(Getdb(), sqlStr)

}
func GetSingleInDb(db *sql.DB, sqlStr string) string {
	var id string
	if logsql {
		public.Log(sqlStr)
	}

	rows, errr := db.Query(sqlStr)
	if errr != nil {
		public.Log(errr)
		return ""
	}
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
		//public.Log("333")
		for i, _ := range columns {
			v := values[i]
			if v == nil {
				id = ""
			} else {

				switch v.(type) {
				default:
					id = fmt.Sprintf("%s", v)
				case bool:
					id = fmt.Sprintf("%s", v) //v
				case int:
					id = fmt.Sprintf("%d", v)
				case int64:
					id = fmt.Sprintf("%d", v)
				case int32:
					id = fmt.Sprintf("%d", v)
				case float64:
					id = fmt.Sprintf("%1.6f", v)
				case float32:
					id = fmt.Sprintf("%1.6f", v)
				case string:
					id = fmt.Sprintf("%s", v)
				case []byte: // -- all cases go HERE!
					id = string(v.([]byte))
				case time.Time:
					id = fmt.Sprintf("%s", v)
				}

			}

		}
	}
	if i == 0 {
		//public.Log("has no row")
		return ""
	}

	//public.Log(id)
	return id
}

func GetSingleJsonInDb(db *sql.DB, sqlStr string) string {
	rows, err := db.Query(sqlStr)
	if logsql {
		public.Log("sqlStr=" + sqlStr)
	}
	if err != nil {
		public.Log(err.Error())
		return ""
	}
	defer rows.Close()

	//defer db.Close()
	columns, _ := rows.Columns()
	count := len(columns)
	//public.Log(count)
	if count == 0 {
		return ""
	}
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	//public.Log("33333")

	final_result := make([]map[string]string, 0)

	result_id := 0
	for i, _ := range columns {
		valuePtrs[i] = &values[i]
	}
	//public.Log("4444")
	for rows.Next() {
		//public.Log("5555")
		rows.Scan(valuePtrs...)
		m := make(map[string]string)
		//Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		for i, col := range columns {
			//public.Log("666")

			v := values[i]

			key := col
			public.Log("key=" + key)
			if v == nil {
				m[key] = ""
			} else {

				switch v.(type) {
				default:
					m[key] = fmt.Sprintf("%s", v)
				case bool:

					m[key] = fmt.Sprintf("%s", v)
				case int:

					m[key] = fmt.Sprintf("%d", v)
				case int64:

					m[key] = fmt.Sprintf("%d", v)
				case float64:

					m[key] = fmt.Sprintf("%1.2f", v)
				case float32:

					m[key] = fmt.Sprintf("%1.2f", v)
				case string:

					m[key] = fmt.Sprintf("%s", v)
				case []byte: // -- all cases go HERE!

					m[key] = string(v.([]byte))
				case time.Time:
					m[key] = fmt.Sprintf("%s", v)
				}
			}
		}
		//public.Log("777")
		//fmt.Print(m)
		final_result = append(final_result, m)

		result_id++
	}

	//public.Log("888")
	jsonData, err := json.Marshal(final_result)
	if err != nil {
		return ""
	}
	fmt.Println(string(jsonData))
	return string(jsonData)
}

func ExecuteUpdate(sqlStr string) (int, error) {

	return ExecuteUpdateInDb(Getdb(), sqlStr)

}

func ExecuteUpdateInDb(db *sql.DB, sqlStr string) (int, error) {

	if logsql {
		public.Log("sqlStr=" + sqlStr)
	}
	//res, err := Getdb().Exec(sqlStr)

	if strings.Contains(strings.ToLower(sqlStr), "insert") {
		sqlStr = strings.Replace(sqlStr, "'00:00:00'", "null", -1)

		if Dbtype == "postgres" {
			rowId := 0
			sqlStr += " RETURNING id "

			err := db.QueryRow(sqlStr).Scan(&rowId)
			if err != nil {
				public.Log("exec sql failed:", err.Error()+" "+sqlStr)
				return 0, err
			} else {
				//public.Log("exec Update sql success")
			}
			public.Log("lastInsertId=")
			public.Log(rowId)
			return rowId, nil

		} else {
			res, err := db.Exec(sqlStr)
			if err != nil {
				public.Log("exec sql failed:", err.Error()+" "+sqlStr)
				return 0, err
			} else {
				//public.Log("exec Update sql success")
			}

			rowId, err := res.LastInsertId()
			if err != nil {
				public.Log("fetch last insert id failed:", err.Error())
				return 0, err
			}

			//public.Log("lastInsertId=")
			//public.Log(rowId)

			str := strconv.FormatInt(rowId, 10)
			//public.Log(str)
			ret, _ := strconv.Atoi(str)
			//public.Log(ret)
			return ret, nil
		}

	} else {
		res, err := db.Exec(sqlStr)
		if err != nil {
			public.Log("exec sql failed:", err.Error()+" "+sqlStr)
			return 0, err
		} else {
			//public.Log("exec Update sql success")
		}
		rowId, err := res.RowsAffected()
		if err != nil {
			public.Log("fetch RowsAffected failed:", err.Error())
			return 0, err
		}
		//public.Log("update idddd=")
		//public.Log(rowId)
		str := strconv.FormatInt(rowId, 10)
		ret, _ := strconv.Atoi(str)
		return ret, nil
	}
	//return 1
}

func ExecuteQuery(sqlStr string) map[int]map[string]string {

	return ExecuteQueryInDb(Getdb(), sqlStr)
}

func ExecuteQueryJson(sqlStr string) (string, error) {

	return ExecuteQueryJsonInDb(Getdb(), sqlStr)

}
func GetRows(sqlStr string) *sql.Rows {
	db := Getdb()
	if db == nil {
		return nil
	}
	rows, err := db.Query(sqlStr)
	if logsql {
		public.Log("sqlStr=" + sqlStr)
	}
	if err != nil {
		public.Log(err.Error())
		return nil
	}
	//defer rows.Close()

	return rows

}
func ExecuteQueryInDb(db *sql.DB, sqlStr string) map[int]map[string]string {

	rows, err := db.Query(sqlStr)
	if logsql {
		public.Log("sqlStr=" + sqlStr)
	}
	if err != nil {
		public.Log(err.Error())
		return nil
	}
	defer rows.Close()

	//defer db.Close()
	columns, _ := rows.Columns()
	count := len(columns)
	//public.Log(count)
	if count == 0 {
		return nil
	}
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	//public.Log("33333")
	final_result := make(map[int]map[string]string)
	result_id := 0
	for i, _ := range columns {
		valuePtrs[i] = &values[i]
	}
	public.Log("4444")
	for rows.Next() {
		public.Log("rows.Next")
		rows.Scan(valuePtrs...)
		m := make(map[string]string)
		//Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		for i, col := range columns {
			//public.Log("666")

			v := values[i]

			key := col
			if v == nil {
				m[key] = ""
			} else {

				switch v.(type) {
				default:
					m[key] = fmt.Sprintf("%s", v)
				case bool:

					m[key] = fmt.Sprintf("%s", v)
				case int:

					m[key] = fmt.Sprintf("%d", v)
				case int64:

					m[key] = fmt.Sprintf("%d", v)
				case float64:

					m[key] = fmt.Sprintf("%1.2f", v)
				case float32:

					m[key] = fmt.Sprintf("%1.2f", v)
				case string:

					m[key] = fmt.Sprintf("%s", v)
				case []byte: // -- all cases go HERE!

					m[key] = string(v.([]byte))
				case time.Time:
					m[key] = fmt.Sprintf("%s", v)
				}
			}
		}
		//public.Log("777")
		//fmt.Print(m)
		final_result[result_id] = m
		result_id++
	}
	//public.Log("888")
	return final_result
}

func ExecuteQueryJsonInDb(db *sql.DB, sqlStr string) (string, error) {

	rows, err := db.Query(sqlStr)
	if logsql {
		public.Log("sqlStr=" + sqlStr)
	}
	if err != nil {
		public.Log(err.Error())
		return "", err
	}
	defer rows.Close()

	//defer db.Close()
	columns, _ := rows.Columns()
	count := len(columns)
	//public.Log(count)
	if count == 0 {
		return "", nil
	}
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	//public.Log("33333")

	final_result := make([]map[string]string, 0)

	result_id := 0
	for i, _ := range columns {
		valuePtrs[i] = &values[i]
	}
	//public.Log("4444")
	//public.Log(count)
	for rows.Next() {
		//public.Log("rows.Next")
		rows.Scan(valuePtrs...)
		m := make(map[string]string)
		//Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		for i, col := range columns {
			v := values[i]
			key := col
			if v == nil {
				m[key] = ""
			} else {
				switch v.(type) {
				default:
					m[key] = fmt.Sprintf("%s", v)
				case bool:
					m[key] = fmt.Sprintf("%s", v)
				case int:
					m[key] = fmt.Sprintf("%d", v)
				case int64:
					m[key] = fmt.Sprintf("%d", v)
				case float64:
					m[key] = fmt.Sprintf("%1.2f", v)
				case float32:
					m[key] = fmt.Sprintf("%1.2f", v)
				case string:
					m[key] = fmt.Sprintf("%s", v)
				case []byte: // -- all cases go HERE!
					m[key] = string(v.([]byte))
				case time.Time:
					m[key] = fmt.Sprintf("%s", v)
				}

			}
		}

		final_result = append(final_result, m)

		result_id++
	}

	jsonData, err := json.Marshal(final_result)
	if err != nil {
		return "", err
	}
	//fmt.Println(string(jsonData))
	return string(jsonData), nil

}

func ExecuteQueryJsonInDb2222(db *sql.DB, sqlStr string) map[int]map[string]string {

	rows, err := db.Query(sqlStr)
	if logsql {
		public.Log("sqlStr=" + sqlStr)
	}
	if err != nil {
		public.Log(err.Error())
		return nil
	}
	defer rows.Close()

	//defer db.Close()
	columns, _ := rows.Columns()
	count := len(columns)
	//public.Log(count)
	if count == 0 {
		return nil
	}
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	//public.Log("33333")
	final_result := make(map[int]map[string]string)
	result_id := 0
	for i, _ := range columns {
		valuePtrs[i] = &values[i]
	}
	//public.Log("4444")
	for rows.Next() {
		//public.Log("5555")
		rows.Scan(valuePtrs...)
		m := make(map[string]string)
		//Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		for i, col := range columns {
			//public.Log("666")
			v := values[i]

			key := col
			if v == nil {
				m[key] = ""
			} else if reflect.TypeOf(v).Kind() == reflect.String {
				/*
					value := fmt.Sprintf("%s", v)
					if strings.Contains(key, "time") {
						public.Log("timevalue=" + value)
						//value = strings.Replace(value, " +0000", "", -1)
					}
				*/

				m[key] = fmt.Sprintf("%s", v) //v
			} else if reflect.TypeOf(v).Kind() == reflect.Int64 {
				m[key] = fmt.Sprintf("%d", v) //v
			} else if reflect.TypeOf(v).Kind() == reflect.Int32 {
				m[key] = fmt.Sprintf("%d", v) //v
			} else if reflect.TypeOf(v).Kind() == reflect.Float64 {
				m[key] = fmt.Sprintf("%1.2f", v) //v
			} else if reflect.TypeOf(v).Kind() == reflect.Float32 {
				m[key] = fmt.Sprintf("%1.2f", v) //v
			} else {
				m[key] = fmt.Sprintf("%s", v) //v
			}
		}
		//public.Log("777")
		//fmt.Print(m)
		final_result[result_id] = m
		result_id++
	}
	//public.Log("888")
	return final_result

}

func GetColumnName(sqlStr string) []string {
	return GetColumnNameInDb(Getdb(), sqlStr)
}

func GetColumnNameInDb(db *sql.DB, sqlStr string) []string {

	rows, err := db.Query(sqlStr)

	if logsql {
		public.Log("sqlStr=" + sqlStr)
	}
	if err != nil {
		CheckErr(err)
		return nil
	}
	defer rows.Close()
	//defer db.Close()
	columns, _ := rows.Columns()

	names := make([]string, len(columns))
	for _, col := range columns {
		names = append(names, col)
	}

	return names

}

func GetColumnValueList(sqlStr string) []string {

	return GetColumnValueListInDb(Getdb(), sqlStr)

}
func GetColumnValueListInDb(db *sql.DB, sqlStr string) []string {

	rows, err := db.Query(sqlStr)

	if logsql {
		public.Log("sqlStr=" + sqlStr)
	}
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

	for i, _ := range columns {
		valuePtrs[i] = &values[i]
	}
	names := make([]string, count)
	for rows.Next() {
		rows.Scan(valuePtrs...)

		var v interface{}
		val := values[0]
		b, err := val.([]byte)
		if err {
			v = string(b)
		} else {
			v = val
		}
		names = append(names, fmt.Sprintf("%s", v))

	}
	return names

}

func CheckErr(err error) {
	if err != nil {
		public.Log(err)
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
