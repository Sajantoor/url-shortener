package database

import (
	"log"

	"github.com/gocql/gocql"
)

type Cassandra struct {
	cluster *gocql.ClusterConfig
	session *gocql.Session
}

func New() *Cassandra {
	log.Println("Connecting to Cassandra...")

	cluster := gocql.NewCluster("127.0.0.1")
	cluster.ProtoVersion = 4
	cluster.Consistency = gocql.Quorum
	cluster.Keyspace = "url_shortener"

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Failed to connect to Cassandra", err)
	}

	log.Println("Connected to Cassandra")

	return &Cassandra{
		cluster: cluster,
		session: session,
	}
}

func (c *Cassandra) Close() {
	c.session.Close()
}

func (c *Cassandra) Client() *gocql.Session {
	return c.session
}

func (c *Cassandra) ExecuteQuery(query string, args ...interface{}) error {
	return c.session.Query(query, args...).Exec()
}
