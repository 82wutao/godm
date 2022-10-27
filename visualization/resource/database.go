package resource

import (
	"dm.net/datamine/orm"
	_ "github.com/lib/pq"
)

var globalConn orm.DatabaseConnection

func InitDBConnection() error {
	//TODO read properties from config
	var prop orm.ConnectionProperties
	prop.Dialect = "postgres"      //      string
	prop.User = "postgres"         // string
	prop.Password = "123456"       //string
	prop.ConnectionType = ""       //string
	prop.Host = "k8s.66ssyl.com"   // string
	prop.Port = 5432               //uint16
	prop.Database = "saas_account" //string
	prop.Charset = ""              //string

	conn, err := orm.OpenConnection(prop)
	if err != nil {
		return err
	}
	globalConn = conn
	return nil
}

func GetDBConnection() orm.DatabaseConnection {
	return globalConn
}
