package targettypes

const (
	DB   = "db"
	SQL  = "sql"
	JSON = "json"
	CSV  = "csv"
)

func List() []string {
	return []string{DB, SQL, JSON, CSV}
}
