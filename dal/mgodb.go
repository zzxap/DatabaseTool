/*
一：官网下载
1：下载  https://www.mongodb.org/downloads  注意要选择32或者64位
我服务器选择的是：mongodb-win32-x86_64-3.0.3
2：新建目录“D:\MongoDB”，解压下载到的安装包，找到bin目录下面全部.exe文件，拷贝到刚创建的目录下。
3：在“D:\MongoDB”目录下新建“data”文件夹，它将会作为数据存放的根文件夹。
4：在“D:\MongoDB”目录下新建“log”文件夹，作为日志文件夹。

二：启动测试

配置Mongo服务端：

　　打开CMD窗口，按照如下方式输入命令：
　　> d:
　　> cd /d D:\MongoDB\bin
cd /d D:\MongoDB\Server\3.6\bin
　　> mongod --dbpath D:\MongoDB\data
　　配置成功后会看到如下画面：
在浏览器输入：http://localhost:27017/，可以看到如下提示：
You are trying to access MongoDB on the native driver port. For http diagnostic access, add 1000 to the port number
　　如此，MongoDB数据库服务已经成功启动了。

三：封装服务：

还是运行cmd，进入MongoDB目录，执行下列命令
控制台执行命令：D:\MongoDB\>
mongod -dbpath "D:\MongoDB\data" -logpath "D:\MongoDB\data\log\MongoDB.log" -install -serviceName "MongoDB"

这里--MongoDB.log就是开始建立的日志文件，--serviceName "MongoDB" 服务名为MongoDB。

　接着启动mongodb服务

　> D:\MongoDB>NET START MongoDB

*/

package dal

import (
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"DataBaseManage/public"

	"gopkg.in/mgo.v2/bson"
)

var _mongodb *mgo.Database
var dbsession *mgo.Session
var URL = ""

func InitMongoDB() bool {
	URL = Dbhost + ":" + Dbport
	_mongodb = getMongoDB()
	if _mongodb == nil {
		return false
	}
	return true
}

// get mongodb db
func getMongoDB() *mgo.Database {

	session, err := mgo.Dial(URL)

	if err != nil {
		return nil
	}

	session.SetMode(mgo.Monotonic, true)
	db := session.DB(Dbname)
	dbsession = session
	return db
}
func IsDBLive() bool {
	if dbsession == nil {
		return false
	}
	err := dbsession.Ping()
	if err != nil {
		public.Log(err)
		InitMongoDB()
		return false
	}
	return true
}

var (
	mgoSession *mgo.Session
	dataBase   = Dbname
)

/**
 * 公共方法，获取session，如果存在则拷贝一份
 */
func GetDb() *mgo.Session {

	if mgoSession == nil {

		var err error
		mgoSession, err = mgo.Dial(URL)
		if err != nil {
			public.Log("connect fail")
			panic(err) //直接终止程序运行
		}
	}
	//最大连接池默认为4096
	return mgoSession.Clone()
}

/*
func GetCollection(tbname string) (*Collection, error) {

	err := _mongodb.Ping()
	if err != nil {
		public.Log(err)
		InitMongoDB()
		return nil, err
	}
	return _mongodb.C(tbname), nil

}
*/

//公共方法，获取collection对象
func WitchCollection(collection string, s func(*mgo.Collection) error) error {
	session := GetDb()
	defer session.Close()
	c := session.DB(dataBase).C(collection)
	return s(c)
}

type Person struct {
	key   string "bson:'key'"
	value string "bson:'value'"
}

func AddPerson(key string, value string) bool {
	collection := _mongodb.C("person") //如果该集合已经存在的话，则直接返回
	count, _ := collection.Find(bson.M{"key": key, "value": value}).Count()

	if count == 0 {

		err := collection.Insert(&Person{key, value})
		if err != nil {
			panic(err)
			return false
		}
		return true
	}
	return false
}

//获取所有的person数据
func GetPersonList() []Person {
	var persons []Person
	query := func(c *mgo.Collection) error {
		//"-id" 按id倒序排列  "id"正序排列
		//public.Log(c.Find(nil).Sort("-id"))
		return c.Find(nil).Sort("-key").Skip(0).Limit(2).All(&persons)
	}
	err := WitchCollection("person", query)

	for i := 0; i < len(persons); i++ {
		per := persons[i]
		public.Log(per.key)
	}

	if err != nil {
		return persons
	}
	return persons
}
