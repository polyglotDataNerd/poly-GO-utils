package database

import (
	p "github.com/magiconair/properties"
	"testing"
)

type DatabaseMock struct{}

//func (obj DatabaseMock) CreateTunnelDB() (*sql.DB, error) {
//	return nil, nil
//}

func TestConnection(t *testing.T) {
	config := p.MustLoadFile("../config.properties", p.UTF8)
	dbMock, _ := TunnelDB{
		ProxyHost:  config.GetString("ProxyHost", ""),
		Sshuser:    config.GetString("Sshuser", ""),
		SshPort:    config.GetInt("SshPort", 22),
		PrivateKey: config.GetString("PrivateKey", ""),
		LocalHost:  config.GetString("LocalHost", ""),
		DBConnect:  config.GetString("Database", ""),
		Port:       config.GetInt("Port", 5432),
	}.CreateTunnelDB()

	dbCheckErr := dbMock.Ping()
	if dbCheckErr != nil {
		t.Error("Could not ping tunneled db", dbCheckErr.Error())
	}

}
