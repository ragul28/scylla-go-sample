package main

import (
	"github.com/ragul28/scylla-go-sample/scylla"
	log "github.com/ragul28/scylla-go-sample/utils"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"github.com/scylladb/gocqlx/table"
	"go.uber.org/zap"
)

type Record struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Address   string `db:"address"`
}

var stmts = createStatements()

func main() {
	logger := log.CreateLogger("info")

	cluster := scylla.CreateCluster(gocql.Quorum, "sample", "localhost:9042")
	session, err := gocql.NewSession(*cluster)
	if err != nil {
		logger.Fatal("unable to connect to scylla", zap.Error(err))
	}
	defer session.Close()

	selectQuery(session, logger)
	insertQuery(session, "John", "Doe", "12 Foo Lane", logger)
	insertQuery(session, "Tony", "Hawk", "34 bar St", logger)
	selectQuery(session, logger)
	deleteQuery(session, "John", "Doe", logger)
	selectQuery(session, logger)
	deleteQuery(session, "Tony", "Hawk", logger)
	selectQuery(session, logger)
}

func deleteQuery(session *gocql.Session, firstName string, lastName string, logger *zap.Logger) {
	logger.Info("Deleting " + firstName + "......")
	r := Record{
		FirstName: firstName,
		LastName:  lastName,
	}
	err := gocqlx.Query(session.Query(stmts.del.stmt), stmts.del.names).BindStruct(r).ExecRelease()
	if err != nil {
		logger.Error("delete catalog.user_data", zap.Error(err))
	}
}

func insertQuery(session *gocql.Session, firstName, lastName, address string, logger *zap.Logger) {
	logger.Info("Inserting " + firstName + "......")
	r := Record{
		FirstName: firstName,
		LastName:  lastName,
		Address:   address,
	}
	err := gocqlx.Query(session.Query(stmts.ins.stmt), stmts.ins.names).BindStruct(r).ExecRelease()
	if err != nil {
		logger.Error("insert catalog.user_data", zap.Error(err))
	}
}

func selectQuery(session *gocql.Session, logger *zap.Logger) {
	logger.Info("Displaying Results:")
	var rs []Record
	err := gocqlx.Query(session.Query(stmts.sel.stmt), stmts.sel.names).SelectRelease(&rs)
	if err != nil {
		logger.Warn("select catalog.mutant", zap.Error(err))
		return
	}
	for _, r := range rs {
		logger.Info("\t" + r.FirstName + " " + r.LastName + ", " + r.Address)
	}
}

func createStatements() *statements {
	m := table.Metadata{
		Name:    "user_data",
		Columns: []string{"first_name", "last_name", "address"},
		PartKey: []string{"first_name", "last_name"},
	}
	tbl := table.New(m)

	deleteStmt, deleteNames := tbl.Delete()
	insertStmt, insertNames := tbl.Insert()
	// Normally a select statement such as this would use `tbl.Select()` to select by
	// primary key but now we just want to display all the records...
	selectStmt, selectNames := qb.Select(m.Name).Columns(m.Columns...).ToCql()
	return &statements{
		del: query{
			stmt:  deleteStmt,
			names: deleteNames,
		},
		ins: query{
			stmt:  insertStmt,
			names: insertNames,
		},
		sel: query{
			stmt:  selectStmt,
			names: selectNames,
		},
	}
}

type query struct {
	stmt  string
	names []string
}

type statements struct {
	del query
	ins query
	sel query
}
