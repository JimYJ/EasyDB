package mysql

func (self *MysqlDB) getVal(query string) (string, error) {
	row := self.dbConn.QueryRow(query)
	var str string
	err := row.Scan(&str)
	// common.ErrHendle("Scan Warning:", err)
	return str, err
}

func (self *MysqlDB) getValByStmt(query string, param ...interface{}) (string, error) {
	stmt, err := self.dbConn.Prepare(query)
	if err != nil {
		return "", nil
	}
	defer stmt.Close()
	row := stmt.QueryRow(param...)
	var str string
	err2 := row.Scan(&str)
	return str, err2
}

func (self *MysqlDB) GetVal(qtype int, query string, param ...interface{}) (string, error) {
	lastQuery = query
	if qtype == Statement {
		return self.getValByStmt(query, param...)
	} else {
		return self.getVal(query)
	}
}
