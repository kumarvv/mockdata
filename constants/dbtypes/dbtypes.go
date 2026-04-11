package dbtypes

const (
	Sqlite     = "sqlite"
	Postgresql = "postgresql"
	Mysql      = "mysql"
)

func List() []string {
	return []string{Sqlite, Postgresql, Mysql}
}
