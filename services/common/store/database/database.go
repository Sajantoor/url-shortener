package database

import (
	"github.com/gocql/gocql"
	"go.uber.org/zap"
)

type Cassandra struct {
	cluster *gocql.ClusterConfig
	session *gocql.Session
}

func New() *Cassandra {
	zap.L().Info("Connecting to Cassandra...")

	cluster := gocql.NewCluster("127.0.0.1")
	cluster.ProtoVersion = 4
	cluster.Consistency = gocql.Quorum
	cluster.Keyspace = "url_shortener"

	session, err := cluster.CreateSession()
	if err != nil {
		zap.L().Sugar().Fatal("Failed to connect to Cassandra", err)
	}

	zap.L().Info("Connected to Cassandra")

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
