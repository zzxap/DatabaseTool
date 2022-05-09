package HTTPBusiness

import (
	"DataBaseManage/dal"
	"fmt"

	"DataBaseManage/public"
	"io"
	"net/http"

	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

func ExportExcelTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}
	//pageid := r.FormValue("pageid")
	//pagesize := r.FormValue("pagesize")
	//userid := r.FormValue("userid")
	table := r.FormValue("table")
	flag := r.FormValue("flag")
	belong := r.FormValue("belong")
	userid := r.FormValue("userid")
	//str := dal.GetLiuyanList(pageid, pagesize, userid, name)
	//WriteValue(w, r, str)SELECT id  FROM product
	sql := "SELECT * from  " + table + "  where userid='" + userid + "'"
	if flag != "" && len(flag) > 0 {
		sql += " and flag= " + flag
	}
	if belong != "" && len(belong) > 0 {
		sql += "  and   belong= '" + belong + "'"
	}

	sql += " order by id desc  limit 1000  OFFSET 0 "
	filename := table + "_" + public.GetRandom() // + flag + belong

	count := dal.GetSingleInDb(dal.Getdb(), " select count(*) as num from "+table)
	public.Log(table + " countt===")
	public.Log(count)
	public.Log("countt===")
	GetExcelFile(sql, filename)
	ResponseFile(w, filename)
}

func GetExcelFile(sqlStr string, filename string) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet(filename)
	if err != nil {
		fmt.Printf(err.Error())

	}
	firstrow := sheet.AddRow()
	rows := dal.GetRows(sqlStr)

	//defer db.Close()
	columns, _ := rows.Columns()
	count := len(columns)
	//public.Log(count)
	if count == 0 {
		cell = firstrow.AddCell()
		cell.Value = "no data"
	}
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	//public.Log("33333")

	for i, _ := range columns {
		valuePtrs[i] = &values[i]
	}

	index := 0
	//public.Log("4444")
	for rows.Next() {
		//public.Log("rows.Next")
		row = sheet.AddRow()
		rows.Scan(valuePtrs...)
		m := make(map[string]string)
		//Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		for i, col := range columns {
			//public.Log("666")

			v := values[i]

			key := strings.ToLower(col)

			if index == 0 {
				cell = firstrow.AddCell()
				cell.Value = key
			}

			cell = row.AddCell()

			if v == nil {
				m[key] = ""
			} else {

				switch v.(type) {
				default:
					cell.Value = fmt.Sprintf("%s", v)
				case bool:
					cell.Value = fmt.Sprintf("%s", v) //v
				case int:
					cell.Value = fmt.Sprintf("%d", v)
				case int64:
					cell.Value = fmt.Sprintf("%d", v)
				case int32:
					cell.Value = fmt.Sprintf("%d", v)
				case float64:
					cell.Value = fmt.Sprintf("%1.2f", v)
				case float32:
					cell.Value = fmt.Sprintf("%1.2f", v)
				case string:
					cell.Value = fmt.Sprintf("%s", v)
				case []byte:
					cell.Value = string(v.([]byte))
				case time.Time:
					cell.Value = fmt.Sprintf("%s", v)
				}
			}

		}
		index++
	}
	//curdir := public.GetCurDir()
	//PthSep := string(os.PathSeparator)
	//filename = curdir + PthSep + filename
	//public.Log(curdir)
	err = file.Save(filename + ".xlsx")
	if err == nil {

	}

}

func ResponseFile(writer http.ResponseWriter, filename string) {
	//curdir := public.GetCurDir()
	//PthSep := string(os.PathSeparator)
	//filename = curdir + PthSep + filename
	//First of check if Get is set in the URL
	Filename := filename + ".xlsx"
	if Filename == "" {
		//Get not set, send a 400 bad request
		http.Error(writer, "Get 'file' not specified in url.", 400)
		return
	}
	//fmt.Println("Client requests: " + Filename)

	//Check if file exists and open
	Openfile, err := os.Open(Filename)
	defer Openfile.Close() //Close after function return
	if err != nil {
		//File not found, send 404
		http.Error(writer, "File not found.", 404)
		return
	}

	//File is found, create and send the correct headers

	//Get the Content-Type of the file
	//Create a buffer to store the header of the file in
	FileHeader := make([]byte, 512)
	//Copy the headers into the FileHeader buffer
	Openfile.Read(FileHeader)
	//Get content type of file
	FileContentType := http.DetectContentType(FileHeader)

	//Get the file size
	FileStat, _ := Openfile.Stat()                     //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	//Send the headers
	writer.Header().Set("Content-Disposition", "attachment; filename="+Filename)
	writer.Header().Set("Content-Type", FileContentType)
	writer.Header().Set("Content-Length", FileSize)

	//Send the file
	//We read 512 bytes from the file already, so we reset the offset back to 0
	Openfile.Seek(0, 0)
	io.Copy(writer, Openfile) //'Copy' the file to the client
	return

}

func GetDataToExcel(sqlStr string) {

}
