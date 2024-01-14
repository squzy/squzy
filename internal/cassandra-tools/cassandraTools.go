package cassandra_tools

import (
	"github.com/gocql/gocql"
	"time"
)

type CassandraTools interface {
	CreateSession() (*gocql.Session, error)
	ExecuteBatch(session *gocql.Session, batch *gocql.Batch) error
	NewBatch(session *gocql.Session) *gocql.Batch
	Close(session *gocql.Session)
}
type cassandraTools struct {
	cluster *gocql.ClusterConfig
}

func NewCassandraTools(clusterName, user, password string, timeout int32) CassandraTools {
	cluster := gocql.NewCluster(clusterName, clusterName, clusterName)
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Duration(timeout) * time.Second
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: user, Password: password}
	return &cassandraTools{
		cluster: cluster,
	}
}

func (c *cassandraTools) CreateSession() (*gocql.Session, error) {
	return c.cluster.CreateSession()
}

func (c *cassandraTools) ExecuteBatch(session *gocql.Session, batch *gocql.Batch) error {
	return session.ExecuteBatch(batch)
}

func (c *cassandraTools) NewBatch(session *gocql.Session) *gocql.Batch {
	return session.NewBatch(gocql.UnloggedBatch)
}

func (c *cassandraTools) Close(session *gocql.Session) {
	session.Close()
}
