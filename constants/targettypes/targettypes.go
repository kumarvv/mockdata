package targettypes

const (
	DB   = "db"
	SQL  = "sql"
	JSON = "json"
)

func List() []string {
	return []string{DB, SQL, JSON}
}
