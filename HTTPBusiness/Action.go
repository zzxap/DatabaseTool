package HTTPBusiness

import (
	"DataBaseManage/dal"
	"DataBaseManage/public"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var SerVersion = "1.68"
var myUserid string
var isneedbuy = true
var isInitDBSuccess = false

func ActionTask(w http.ResponseWriter, r *http.Request) {

	if r.Method == "OPTIONS" {
		return
	}

	table := r.FormValue("table")
	//public.Log("table= " + table)
	if table != "connectdb" && table != "uploadsqlite" {
		if isInitDBSuccess == false {
			WriteValue(w, r, "-1", "db disconnect", "", "")
			return
		}

		token := r.Header.Get("Authorization")
		if len(token) == 0 {
			token = r.FormValue("token")
		}

		if !ValidateJwtToken(token, "1000") {
			public.Log("add data token error token=" + token)
			public.Log(r.URL.Path)

			ret := "{\"code\":2,\"userid\":1000,\"message\":\"token error\",\"data\":\"[]\"}"
			w.Write([]byte(ret))
			return
		}

	}
	if table == "connectdb" {
		InitDbTask(w, r)
	} else if table == "table" {
		GetTableTask(w, r)
	} else if table == "database" {
		GetDBTask(w, r)
	} else if table == "column" {
		GetColumnTask(w, r)
	} else if table == "runsql" {
		RunSqlnTask(w, r)
	} else if table == "uploadsqlite" {
		uploadSqliteTask(w, r)
	} else if table == "etcd" {
		getEtcdDataTask(w, r)
	}

}
func AddKVTask(w http.ResponseWriter, r *http.Request) {

	if r.Method == "OPTIONS" {
		return
	}
	key := r.FormValue("key")
	value := r.FormValue("value")
	mode := r.FormValue("mode")
	ret := false
	if mode == "add" {
		ret = dal.PutETCD(key, value)
	} else if mode == "edit" {
		ret = dal.PutETCD(key, value)
	} else if mode == "delete" {
		ret = dal.Delete(key)
	}

	if ret == true {
		WriteValue(w, r, "1", "success", "", "")
	} else {
		WriteValue(w, r, "0", "fail", "", "")
	}

}

func InitDbTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	osname := public.GetOSName()
	if osname != "darwin" && osname != "windows" {

		WriteValue(w, r, "0", "fail", "", "")
		return
	}

	ip := r.FormValue("ip")
	port := r.FormValue("port")
	username := r.FormValue("username")
	password := r.FormValue("password")
	dbtype := r.FormValue("dbtype")
	dbname := r.FormValue("dbname")
	sshHost := r.FormValue("sshHost")
	sshPort := r.FormValue("sshPort")
	sshUser := r.FormValue("sshUser")
	sshPass := r.FormValue("sshPass")

	public.Log("InitDbTask " + dbtype)
	errmessage := ""
	if dbtype == "sqlite" {
		dal.Dbtype = "sqlite"
		dal.InitDB()

		if dbname == "demo.sqlite" {
			if dal.ISmac == 1 {
				public.Log("is mac")
				dbname = "InvoicingMac.sqlite"
				//public.SetKeyValue("SqliteFileName", "InvoicingMac.sqlite")
			} else {
				public.Log("is not mac")
				dbname = "db.sqlite"
				//public.SetKeyValue("SqliteFileName", "db.sqlite")
			}

		} else {
			dal.SqliteFileName = dbname
		}

		public.SetKeyValue("SqliteFileName", dbname)

		isInitDBSuccess = true
		token, _ := GetJwtToken("1000")
		ret := "{\"code\":1,\"token\":\"" + token + "\",\"total\":0,\"pagesize\":50,\"message\":\"success\",\"data\":[]}"
		w.Write([]byte(ret))
		return
	} else {
		if len(ip) > 0 && len(port) > 0 && len(username) > 0 && len(password) > 0 {
			json := "{\"ip\":\"" + ip + "\",\"port\":\"" + port + "\",\"username\":\"" + username + "\",\"password\":\"" + password + "\"}"
			public.StoreValue(json)
			var postRead = ""
			public.LoadValue(&postRead)
			public.Log("init db \n")
			public.Log(postRead)

			dal.Dbhost = ip
			dal.Dbport = port
			dal.Dbuser = username
			dal.Dbpassword = password
			dal.Dbname = dbname
			dal.Dbtype = dbtype
			ret := false
			message := ""
			if len(sshHost) > 0 && len(sshPort) > 0 && len(sshUser) > 0 && len(sshPass) > 0 {
				dal.SSHHost = sshHost
				dal.SSHPort = sshPort
				dal.SSHUser = sshUser
				dal.SSHPass = sshPass
				ret, message = dal.InitDBBySSH()
			} else {
				ret, message = dal.InitDB()
			}

			errmessage = message
			if ret {

				token, _ := GetJwtToken("1000")

				isInitDBSuccess = true

				ret := "{\"code\":1,\"token\":\"" + token + "\",\"total\":0,\"pagesize\":50,\"message\":\"success\",\"data\":[]}"
				w.Write([]byte(ret))

				return
			}

		}
		public.Log("init fail  \n")
		isInitDBSuccess = false
		WriteValue(w, r, "0", errmessage, "", "")

	}
}
func getEtcdDataTask(w http.ResponseWriter, r *http.Request) {
	mp := dal.GetMapArray("")
	//WriteValue(w, r, "1", "success", "", "")
	lenth := strconv.Itoa(len(mp))
	jsonString, err := json.Marshal(mp)
	if err == nil && mp != nil {
		WriteValue(w, r, "1", "success", string(jsonString[:]), lenth)
	} else {
		WriteValue(w, r, "0", "fail", "", "")
	}

}

func uploadSqliteTask(w http.ResponseWriter, r *http.Request) {
	public.Log("Add uploadSqliteTask ")
	PthSep := string(os.PathSeparator)
	fileNames := public.SaveUploadFileask(w, r, public.GetCurDir()+PthSep+"db"+PthSep)

	if len(fileNames) > 0 {

		isInitDBSuccess = true
		token, _ := GetJwtToken("1000")
		dal.SqliteFileName = fileNames
		public.Log("Add uploadSqliteTask fileNames= " + fileNames)

		dal.Dbtype = "sqlite"
		dal.InitDB()

		ret := "{\"code\":1,\"token\":\"" + token + "\",\"total\":0,\"pagesize\":50,\"message\":\"success\",\"data\":[]}"
		w.Write([]byte(ret))

	} else {
		WriteValue(w, r, "0", "fail", "", "")
	}
}

func RunSqlnTask(w http.ResponseWriter, r *http.Request) {
	sql := r.FormValue("sql")
	//public.Log("runsql=" + sql)
	lowsql := strings.ToLower(sql)

	sql = strings.Replace(sql, "‘", "'", -1)
	sql = strings.Replace(sql, "“", "\"", -1)
	sql = strings.Replace(sql, "\\n", "", -1)
	if strings.Contains(lowsql, "select") {
		jsonstr, err := dal.ExecuteQueryJson(sql)
		if err == nil {
			WriteValue(w, r, "1", "success", jsonstr, "")
		} else {
			WriteValue(w, r, "0", string(err.Error()), "", "")
		}
	} else {
		_, err := dal.ExecuteUpdate(sql)

		if err == nil {
			WriteValue(w, r, "1", "success", "", "")
		} else {
			WriteValue(w, r, "0", replacemessage(string(err.Error())), "", "")
		}
	}

}

func replacemessage(message string) string {
	message = strings.Replace(message, "'", "", -1)
	message = strings.Replace(message, ":", "", -1)
	message = strings.Replace(message, "\\n", "", -1)
	message = strings.Replace(message, "\"", "", -1)
	return message
}

/*
SELECT
  m.name AS table_name,
  p.cid AS col_id,
  p.name AS col_name,
  p.type AS col_type,
  p.pk AS col_is_pk,
  p.dflt_value AS col_default_val,
  p.[notnull] AS col_is_not_null
FROM sqlite_master m
LEFT OUTER JOIN pragma_table_info((m.name)) p
  ON m.name <> p.name
WHERE m.type = 'table'
and m.name='Users'




*/
func GetColumnTask(w http.ResponseWriter, r *http.Request) {
	tablename := r.FormValue("tablename")
	total := "50"

	sql := ""
	totalSql := ""
	if dal.Dbtype == "sqlite" {

		sql = `SELECT p.name AS column_name,  p.type AS data_type
FROM sqlite_master m
LEFT OUTER JOIN pragma_table_info((m.name)) p
  ON m.name <> p.name
WHERE m.type = 'table'
and m.name=` + "'" + tablename + "'"

		totalSql = `SELECT  count(*) as num
FROM sqlite_master m
LEFT OUTER JOIN pragma_table_info((m.name)) p
  ON m.name <> p.name
WHERE m.type = 'table'
and m.name=` + "'" + tablename + "'"

	} else if dal.Dbtype == "postgres" {

		sql = "select column_name,data_type from information_schema.columns  where table_name = '" + tablename + "';"

		totalSql = "select count(*) as num  from information_schema.columns  where table_name = '" + tablename + "';"

	} else if dal.Dbtype == "mysql" {

		//sql = "SELECT column_name,data_type FROM INFORMATION_SCHEMA.COLUMNS WHERE  table_schema='" + dal.Dbname + "' and TABLE_NAME = '" + tablename + "';"
		//totalSql = "SELECT  count(*) as num   FROM INFORMATION_SCHEMA.COLUMNS WHERE   table_schema='" + dal.Dbname + "' and  TABLE_NAME = '" + tablename + "';"
		sql = "SELECT column_name,data_type FROM INFORMATION_SCHEMA.COLUMNS WHERE  TABLE_NAME = '" + tablename + "';"
		totalSql = "SELECT  count(*) as num   FROM INFORMATION_SCHEMA.COLUMNS WHERE    TABLE_NAME = '" + tablename + "';"

	} else if dal.Dbtype == "mssql" {

		sql = "SELECT column_name,data_type FROM INFORMATION_SCHEMA.COLUMNS WHERE  TABLE_NAME = '" + tablename + "'"
		totalSql = "SELECT  count(*) as num   FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = '" + tablename + "'"
	}
	total = dal.GetSingle(totalSql)
	jsonstr, err := dal.ExecuteQueryJson(sql)
	//public.Log(jsonstr)
	if err != nil {
		WriteValue(w, r, "0", err.Error(), "", "")
	} else {
		WriteValue(w, r, "1", "success", jsonstr, total)
	}

}

func GetTableTask(w http.ResponseWriter, r *http.Request) {

	dbname := r.FormValue("dbname")
	//public.Log("GetTableTask " + dal.Dbtype)
	total := "0"
	if dal.Dbtype != "sqlite" {
		dal.ReInitDB(dal.Dbtype, dal.Dbhost, dal.Dbport, dal.Dbuser, dal.Dbpassword, dbname)
	}

	sql := ""
	totalSql := ""
	if dal.Dbtype == "sqlite" {
		sql = "SELECT name as table_name FROM sqlite_master WHERE type='table' ORDER BY table_name "
		totalSql = "SELECT count(*) as num FROM sqlite_master WHERE type='table' "
	} else if dal.Dbtype == "postgres" {
		sql = "SELECT table_name FROM information_schema.tables where table_schema='public' AND table_type='BASE TABLE' ORDER BY table_name   "

		totalSql = "SELECT count(*) as num FROM information_schema.tables where table_schema='public' AND table_type='BASE TABLE' "
	} else if dal.Dbtype == "mysql" {

		sql = "SELECT lower(table_name) as table_name FROM information_schema.tables WHERE table_schema = '" + dbname + "' ORDER BY table_name  "
		totalSql = "SELECT  count(*) as num  FROM information_schema.tables WHERE table_schema = '" + dbname + "'"
	} else if dal.Dbtype == "mssql" {

		sql = "SELECT  lower(table_name) as table_name  FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE='BASE TABLE'  AND TABLE_CATALOG='" + dbname + "'"
		totalSql = "SELECT  count(*) as num  FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE='BASE TABLE' AND TABLE_CATALOG='" + dbname + "'"
	}

	//public.Log(sql)

	total = dal.GetSingle(totalSql)
	jsonstr, err := dal.ExecuteQueryJson(sql)
	//public.Log(jsonstr)
	if jsonstr == "[]" && dal.Dbtype == "sqlite" {
		addDefaulttable()
		jsonstr, err = dal.ExecuteQueryJson(sql)
		total = dal.GetSingle(totalSql)
	}

	if err != nil {
		WriteValue(w, r, "0", err.Error(), "", "")
	} else {
		WriteValue(w, r, "1", "success", jsonstr, total)
	}

}

func addDefaulttable() {
	sql := `CREATE TABLE "users"  
  ("ID" INTEGER ,
 "UserId" INTEGER,"UserName" VARCHAR,
 "IdNum" VARCHAR,
 "BirthDay" datetime,
 "AddTime" datetime)`
	dal.ExecuteUpdate(sql)
	sql = `CREATE TABLE "product"  
  ("ID" INTEGER ,
 "UserId" INTEGER,"productName" VARCHAR,
 "Barcode" VARCHAR,
 "price" float,
 "AddTime" datetime)`
	dal.ExecuteUpdate(sql)

}
func GetDBTask(w http.ResponseWriter, r *http.Request) {
	sql := ""
	total := "0"
	totalSql := ""
	if dal.Dbtype == "sqlite" {
		sql = "SELECT name as dbname FROM sqlite_master WHERE type='database' ORDER BY name"
		totalSql = "SELECT  count(*) as num  FROM sqlite_master WHERE type='database' ORDER BY name"

		WriteValue(w, r, "1", "success", `[{"dbname":"`+dal.SqliteFileName+`"}]`, "1")
		return

	} else if dal.Dbtype == "postgres" {

		sql = "SELECT datname as dbname FROM pg_database WHERE datistemplate = false "
		totalSql = "SELECT   count(*) as num  FROM pg_database WHERE datistemplate = false "
	} else if dal.Dbtype == "mysql" {
		//SELECT SCHEMA_NAME AS `dbname` FROM INFORMATION_SCHEMA.SCHEMATA;
		sql = "SELECT SCHEMA_NAME AS `dbname` FROM INFORMATION_SCHEMA.SCHEMATA;"
		totalSql = "SELECT     count(*) as num  FROM INFORMATION_SCHEMA.SCHEMATA;"

	} else if dal.Dbtype == "mssql" {

		sql = "　SELECT Name as dbname FROM Master..SysDatabases ORDER BY Name"
		totalSql = "　SELECT   count(*) as num  FROM Master..SysDatabases ORDER BY Name"
	} else {

	}
	total = dal.GetSingle(totalSql)
	jsonstr, err := dal.ExecuteQueryJson(sql)
	//public.Log(jsonstr)
	if err != nil {
		WriteValue(w, r, "0", err.Error(), "", "")
	} else {
		WriteValue(w, r, "1", "success", jsonstr, total)
	}

}

func GetHistoryPriceTask(w http.ResponseWriter, r *http.Request) {

}

func WriteValue(w http.ResponseWriter, r *http.Request, code, message, data, total string) {

	//w.Header().Set("Access-Control-Allow-Origin", "*")
	if data == "" {
		data = "[]"
	}
	if total == "" {
		total = "0"
	}
	message = strings.Replace(message, "\"", "\\\"", -1)
	message = strings.Replace(message, ":", "", -1)
	message = strings.Replace(message, ",", "", -1)
	ret := "{\"code\":" + code + ",\"total\":" + total + ",\"pagesize\":50,\"message\":\"" + message + "\",\"data\":" + data + "}"

	w.Write([]byte(ret))

}
