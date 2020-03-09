package main

import (
	"fmt"
	read "github.com/polyglotDataNerd/zip-Go-Utils/reader"
	scan "github.com/polyglotDataNerd/zip-Go-Utils/scanner"
	database "github.com/polyglotDataNerd/zip-Go-Utils/sgdatabase"
	goutils "github.com/polyglotDataNerd/zip-Go-Utils/sgutils"
	"runtime"
	"strings"
	"sync"
	"time"
)

func main() {
	fmt.Println("Welcome to data infrastructure GOLang Utils")
	cassandraRead()
	//cassandra()
	//channels()
	//
	///*Examples for UTILS package*/
	//config := p.MustLoadFile("config.properties", p.UTF8)
	///*get DDB item*/
	//items := aws.DDBT{Attribute: "d29f0e57-f44e-11e9-97e1-0d091e77f50e"}.DDBGetQuery("sg-id-systemIdMapping-testing", "idx_uuid", "uuid")
	//goutils.Log.Println(*items)
	///*get DDB all items*/
	//aws.DDBT{}.DDBScanGetItems("sg-id-systemIdMapping-testing", "sourceSystemId", "uuid")
	//
	///*creates a session*/
	//accesskey := aws.SSMParams("/ubereats/ingest/accesskey", 0)
	//secretkey := aws.SSMParams("/ubereats/ingest/secretkey", 0)
	//aws.SessionGenerator(accesskey, secretkey, "us-west-2")
	//
	///*creates a session*/
	//sess := aws.SessionGenerator("default", "us-west-2")
	//println(sess.Config)
	//
	///*tests directory get*/
	//dataMap, _ := aws.S3Obj{
	//	Bucket: "sweetgreen-bigdata-application",
	//	Key:    "platform/raw/web/2019-08-25",
	//}.S3ReadObjDir(sess)
	//for k, v := range dataMap {
	//	goutils.Log.Println(k, v)
	//}
	//
	///*creates a postgres call using a tunnel and param store connection string to run locally*/
	//dbCon, _ := database.TunnelDB{
	//	ProxyHost:  config.GetString("ProxyHost", ""),
	//	Sshuser:    config.GetString("Sshuser", ""),
	//	SshPort:    config.GetInt("SshPort", 22),
	//	PrivateKey: config.GetString("PrivateKey", ""),
	//	LocalHost:  config.GetString("LocalHost", ""),
	//	DBConnect:  aws.SSMParams(config.GetString("Sweettouch", ""), 0),
	//	Port:       config.GetInt("Port", 5439),
	//}.CreateDB()
	//defer dbCon.Close()
	//rows, err := dbCon.Query("SELECT id from reporting.gravy_customers limit 10;")
	//if err != nil {
	//	goutils.Log.Fatalln(err)
	//}
	//defer rows.Close()
	//for rows.Next() {
	//	var date string
	//	rows.Scan(&date)
	//	fmt.Println(date)
	//}
	///*creates a postgres call using param store connection string to run in AWS*/
	//dbCon2, _ := database.DBCon{
	//	ConnectionString: aws.SSMParams(config.GetString("Sweettouch", ""), 0),
	//}.CreateDB()
	//defer dbCon2.Close()
	//rows2, err := dbCon2.Query("SELECT id from reporting.gravy_customers limit 10;")
	//if err != nil {
	//	goutils.Log.Fatalln(err)
	//}
	//defer rows2.Close()
	//for rows2.Next() {
	//	var date string
	//	rows2.Scan(&date)
	//	fmt.Println(date)
	//}
}

func channels() {
	/*
	   CHANNELS
	   fan out: Multiple functions can read from the same channel until that channel is closed;
	   this is called fan-out. This provides a way to distribute work amongst a group of workers to parallelize CPU use and I/O.
	*/
	goutils.Info.Println(runtime.NumCPU())
	chLine := make(chan string)
	chOut := make(chan string)
	start := time.Now()

	/* producer */
	go scan.ProcessDir(chLine, "sweetgreen-bigdata-unloads", "cassandra/2020-01-27/ml_service_008.gz")
	/* consumer */
	go read.ReadObj(chLine, chOut)

	for l := range chOut {
		fmt.Println(l)
	}

	fmt.Println("Runtime took ", time.Since(start))
	/* CHANNELS */

}

func cassandra() {

	chLine := make(chan string)
	chOut := make(chan string)
	var wg sync.WaitGroup

	props := goutils.Mutator{
		SetterKeyEnv:    "host",
		SetterValueEnv:  "cassandra.us-east-1.amazonaws.com",
		SetterKeyUser:   "",
		SetterValueUser: "",
		SetterKeyPW:     "",
		SetterValuePW:   "",
	}

	clientConfig := database.CQLProps{
		Mutator:   props,
		TableName: "order_history",
		Keyspace:  "sg_cass",
	}

	client := database.CQL{
		CQLProps:    clientConfig,
		ChannelLine: chLine,
		ChannelOut:  chOut,
		S3Bucket:    "sweetgreen-bigdata-unloads",
		S3key:       "cassandra/2020-02-16/",
		Wg:          wg,
		SSLPath:     "/Users/gerardbartolome/.mac-ca-roots",
	}

	session := client.CassandraSession()
	defer session.Close()
	client.CassS3Write("INSERT INTO", strings.Split("customer_id,gid,order_id,order_date,entree", ","), session)

}

func cassandraRead() {
	var wg sync.WaitGroup

	props := goutils.Mutator{
		SetterKeyEnv:    "host",
		SetterValueEnv:  "cassandra.us-east-1.amazonaws.com",
		SetterKeyUser:   "user",
		SetterValueUser: "",
		SetterKeyPW:     "pw",
		SetterValuePW:   "",
	}

	clientConfig := database.CQLProps{
		Mutator: props,
	}

	client := database.CQL{
		CQLProps: clientConfig,
		Wg:       wg,
		SSLPath:  "/Users/gerardbartolome/.mac-ca-roots",
	}

	session := client.CassandraSession()
	resultSet, rerr := client.CassReadOrderHistory("SELECT * FROM sg_cass.order_history where gid = 'd2618bd1-f44e-11e9-9d0b-0d091e77f50e' LIMIT 500 ALLOW FILTERING", session)
	if rerr != nil {
		goutils.Error.Fatalln(rerr)
	}
	for _, v := range resultSet {
		goutils.Info.Println(v.OrderDate, v.Gid, v.Entree)
	}
	session.Close()
}
