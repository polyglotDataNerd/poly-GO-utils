package database

//https://gist.github.com/vinzenz/7b6b1bf8d0c2b2b1e0d69a15ba9f02c7
//https://blog.alexellis.io/golang-writing-unit-tests/
import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	goutils "github.com/polyglotDataNerd/poly-Go-utils/utils"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"time"

	//postgres driver for go lang
	"github.com/lib/pq"
)

type Database interface {
	CreateDB() *sql.DB
}

type TunnelDB struct {
	ProxyHost  string
	Sshuser    string
	SshPort    int
	PrivateKey string
	LocalHost  string
	DBConnect  string
	Port       int
}

type DBCon struct {
	ConnectionString string
}

type SSHDialer struct {
	client *ssh.Client
}

func (obj *SSHDialer) Open(s string) (driver.Conn, error) {
	return pq.DialOpen(obj, s)
}

func (obj *SSHDialer) Dial(network, address string) (net.Conn, error) {
	return obj.client.Dial(network, address)
}

func (obj *SSHDialer) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	return obj.client.Dial(network, address)
}

func pemauthDB(pempath string) ssh.Signer {

	pemBytes, err := ioutil.ReadFile(pempath)
	if err != nil {
		goutils.Error.Fatal(err)
	}

	signer, err := ssh.ParsePrivateKey(pemBytes)
	fmt.Println("public key passed")
	if err != nil {
		goutils.Error.Fatalf("parse key failed: %s", err.Error())
	}
	return signer
}

func (obj TunnelDB) CreateDB() (*sql.DB, error) {
	authmethod := []ssh.AuthMethod{ssh.PublicKeys(pemauthDB(obj.PrivateKey))}
	config := &ssh.ClientConfig{
		User:            obj.Sshuser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            authmethod,
	}

	// Establish a connection to the local ssh-agent
	sshCon, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", obj.ProxyHost, obj.SshPort), config)
	if err != nil {
		goutils.Error.Fatalf("ssh.Dial failed: %s", err.Error())
	}

	sql.Register("postgres+ssh", &SSHDialer{sshCon})
	db, err := sql.Open("postgres+ssh", obj.DBConnect)
	if err != nil {
		goutils.Error.Fatalf("database connection failed: %s", err.Error())
	}
	return db, nil

}

func (obj DBCon) CreateDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", obj.ConnectionString)
	if err != nil {
		goutils.Error.Fatalf("database connection failed: %s", err.Error())
	}
	return db, nil

}
