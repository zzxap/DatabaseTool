module DataBaseManage

go 1.18

require (
	DataBaseManage/HTTPBusiness v1.0.0
	DataBaseManage/asset v1.0.0
	DataBaseManage/dal v1.0.0
	DataBaseManage/public v1.0.0
	github.com/elazarl/go-bindata-assetfs v1.0.1
	github.com/gorilla/mux v1.8.0
	github.com/lib/pq v1.10.5
	github.com/mattn/go-sqlite3 v1.14.12
	github.com/webview/webview v0.0.0-20220507210603-42b68f4dda20
)

require (
	github.com/bwmarrin/snowflake v0.3.0 // indirect
	github.com/coreos/etcd v2.3.8+incompatible // indirect
	github.com/dchest/captcha v0.0.0-20200903113550-03f5f0333e1f // indirect
	github.com/denisenkom/go-mssqldb v0.12.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/etcd-io/etcd v2.3.8+incompatible // indirect
	github.com/go-ini/ini v1.66.4 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-sql/civil v0.0.0-20190719163853-cb61b32ac6fe // indirect
	github.com/golang-sql/sqlexp v0.0.0-20170517235910-f1bb20e5a188 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/otiai10/copy v1.7.0 // indirect
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e // indirect
	github.com/tealeg/xlsx v1.0.5 // indirect
	github.com/zheng-ji/goSnowFlake v0.0.0-20180906112711-fc763800eec9 // indirect
	golang.org/x/crypto v0.0.0-20220507011949-2cf3adece122 // indirect
	golang.org/x/image v0.0.0-20191009234506-e7c1f5e7dbb8 // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
)

replace DataBaseManage/HTTPBusiness v1.0.0 => ./HTTPBusiness

replace DataBaseManage/dal v1.0.0 => ./dal

replace DataBaseManage/public v1.0.0 => ./public

replace DataBaseManage/asset v1.0.0 => ./asset
