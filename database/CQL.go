package database

//BASIC CRUD and CONNECTION FUNCTONS/METHODS FOR APACHE CASSANDRA

import (
	c "github.com/gocql/gocql"
	read "github.com/polyglotDataNerd/zib-Go-utils/reader"
	"github.com/polyglotDataNerd/zib-Go-utils/scanner"
	aws "github.com/polyglotDataNerd/zib-Go-utils/aws"
	utils "github.com/polyglotDataNerd/zib-Go-utils/utils"
	"strings"
	"sync"
	"time"
)

type OrderHistoryTable struct {
	Gid        string
	OrderDate  time.Time
	OrderId    string
	Entree     string
	CustomerID string
}

type CQLProps struct {
	utils.Mutator
	TableName string
	Keyspace  string
}

type CQL struct {
	CQLProps
	ChannelLine chan string
	ChannelOut  chan string
	S3Bucket    string
	S3key       string
	Wg          sync.WaitGroup
	SSLPath     string
}

//http://www.code2succeed.com/go-cassandra-crud-example/
func (t *CQL) CassandraSession() *c.Session {
	/*hosts*/
	t.SetEnv()
	env := strings.Split(t.GetEnv(), ",")
	/*userID*/
	t.SetUser()
	user := t.GetUser()
	/*userPW*/
	t.SetPW()
	pw := t.GetPW()

	/*loads env variables*/
	sslopts := &c.SslOptions{
		CaPath: t.SSLPath,
	}

	clusterConfig := c.NewCluster(env[0])
	clusterConfig.Port = 9142
	clusterConfig.CQLVersion = "3.4.4"
	clusterConfig.ProtoVersion = 4
	clusterConfig.Authenticator = c.PasswordAuthenticator{
		Username: user,
		Password: pw,
	}
	clusterConfig.SslOpts = sslopts
	clusterConfig.Consistency = c.LocalOne
	clusterConfig.Keyspace = t.Keyspace
	clusterConfig.MaxWaitSchemaAgreement = time.Duration(5) * time.Minute
	clusterConfig.Timeout = time.Duration(5) * time.Minute
	clusterConfig.ConnectTimeout = time.Duration(5) * time.Minute
	clusterConfig.Compressor = c.SnappyCompressor{}

	session, sessionErr := clusterConfig.CreateSession()
	if sessionErr != nil {
		utils.Error.Fatalln("Connection Failed ", sessionErr)
	}

	return session

}

func (t *CQL) CassReadOrderHistory(queryString string, session *c.Session) ([]OrderHistoryTable, error) {
	start := time.Now()
	result := map[string]interface{}{}
	var table []OrderHistoryTable
	/* read from Cassandra */
	t.Wg.Add(1)
	time.Sleep(1 * time.Millisecond)
	go func() {
		defer t.Wg.Done()
		iter := session.Query(queryString).Iter()
		for iter.MapScan(result) {
			table = append(table, OrderHistoryTable{
				Gid:        result["gid"].(string),
				OrderDate:  result["order_date"].(time.Time),
				OrderId:    result["order_id"].(string),
				Entree:     result["entree"].(string),
				CustomerID: result["customer_id"].(string),})
			result = map[string]interface{}{}
		}
	}()
	t.Wg.Wait()
	utils.Info.Println("response time CassReadOrderHistory: ", time.Since(start))
	return table, nil
}

func (t *CQL) CassS3Write(insertStatment string, fields []string, session *c.Session) {
	start := time.Now()
	/* producer */
	go scanner.ProcessDir(t.ChannelLine, t.S3Bucket, t.S3key, "gzip")
	/* consumer */
	go read.ReadObj(t.ChannelLine, t.ChannelOut)

	for line := range t.ChannelOut {
		/* builds insert index: valueString must match fields in the number of length and indexes
		var valueString []string
		for range fields {
			valueString = append(valueString, "?")
		}*/
		/* insert into Cassandra */
		t.Wg.Add(1)
		time.Sleep(1 * time.Millisecond)
		go func() {
			defer t.Wg.Done()
			values := strings.Join(strings.Split(strings.ReplaceAll(line, "\"", "'"), "\t"), ",")
			insert := insertStatment + " " + t.Keyspace + "." + t.TableName + "(" + strings.Join(fields, ", ") + ") VALUES(" + values + ") IF NOT EXISTS;"
			insertErr := session.Query(insert).Exec()
			if insertErr != nil {
				utils.Error.Println(insertErr.Error(), ":", insert)
			}
			utils.Info.Println("insert success: ", insert)
		}()
	}
	t.Wg.Wait()
	utils.Info.Println("Runtime took ", time.Since(start))
}

func (t *CQL) CassCopy(copyStament string, fields []string, delimiter string, header string, session *c.Session) {
	defer session.Close()
	obj, objerr := aws.S3Obj{
		t.S3Bucket,
		t.S3key}.S3ReadObjGZIPDir(aws.SessionGenerator("default", "us-west-2"))

	if objerr != nil {
		utils.Error.Fatalln(objerr.Error())
	}
	for _, v := range obj {
		copy := copyStament + " " + t.TableName + " (" + strings.Join(fields, ", ") + " FROM " + "\\'" + v + "\\'" + " WITH delimiter= " + delimiter + " AND header=" + header
		session.Query(copy)
	}
}
