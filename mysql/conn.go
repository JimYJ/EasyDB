package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"sync"
)

var (
	Statement              = 1
	Normal                 = 0
	insert                 = 0
	update                 = 1
	delete                 = 2
	charset                = "utf8"
	customColumns []string = nil
	mysqldb       *MysqlDB
	once          sync.Once
	dBHost        string
	dBPort        int
	dBName        string
	dBuser        string
	dBpass        string
	dbcharset     string
	isinit        bool
	maxIdleConns  int
	maxOpenConns  int
	lastQuery     string
	showErrors    = true
)

var (
	errorInit     = "DB param is not initialize!"
	errorSetField = "Field List is Error!"
	errorTxInit   = "Transaction didn't initializtion!"
)

func Init(MysqlDBHost string, MysqlDBPort int, MysqlDBName, MysqlDBuser, MysqlDBpass, MysqlDBcharset string, MaxIdleConns, MaxOpenConns int) {
	dBHost = MysqlDBHost
	dBPort = MysqlDBPort
	dBName = MysqlDBName
	dBuser = MysqlDBuser
	dBpass = MysqlDBpass
	isinit = true
	maxIdleConns = MaxIdleConns
	maxOpenConns = MaxOpenConns
	if len(MysqlDBcharset) == 0 {
		dbcharset = charset
	} else {
		dbcharset = MysqlDBcharset
	}
}

func GetMysqlConn() (*MysqlDB, error) {
	if isinit == false {
		return nil, errors.New(errorInit)
	}
	var err error
	once.Do(func() {
		mysqldb = &MysqlDB{dBHost, dBuser, dBName, dBpass, dbcharset, dBPort, nil, nil, nil}
		err = mysqldb.Conn(maxIdleConns, maxOpenConns)
		printErrors(err)
	})
	return mysqldb, err
}

func NewMysqlConn(MysqlDBHost string, MysqlDBPort int, MysqlDBName string, MysqlDBuser string, MysqlDBpass string, MysqlDBcharset string, MaxIdleConns int, MaxOpenConns int) (*MysqlDB, error) {
	var err error
	var DBcharset string
	if len(MysqlDBcharset) == 0 {
		DBcharset = charset
	} else {
		DBcharset = MysqlDBcharset
	}
	mysqldb = &MysqlDB{MysqlDBHost, MysqlDBuser, MysqlDBName, MysqlDBpass, DBcharset, MysqlDBPort, nil, nil, nil}
	err = mysqldb.Conn(MaxIdleConns, MaxOpenConns)
	printErrors(err)
	return mysqldb, err
}

//MysqlDB include mysqldb all params
type MysqlDB struct {
	host, user, dbname, pass, charset string
	port                              int
	dbConn                            *sql.DB
	fieldlist                         []string
	tx                                *sql.Tx
}

func (mdb *MysqlDB) Conn(MaxIdleConns int, MaxOpenConns int) error {
	if mdb.host == "" || mdb.pass == "" || mdb.user == "" || mdb.dbname == "" {
		errs := errors.New(errorInit)
		printErrors(errs)
		return errs
	}
	lastQuery = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", mdb.user, mdb.pass, mdb.host, mdb.port, mdb.dbname, mdb.charset)
	db, err := sql.Open("mysql", lastQuery)
	if err != nil {
		printErrors(err)
		return err
	}
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetMaxOpenConns(MaxOpenConns)
	err2 := db.Ping()
	if err2 != nil {
		printErrors(err2)
		return err2
	}
	mdb.dbConn = db
	return nil
}

func (mdb *MysqlDB) Close() {
	mdb.dbConn.Close()
}

func (mdb *MysqlDB) SetFields(fieldlist []string) {
	mdb.fieldlist = fieldlist
}
