package targettypes

const (
	SQL  = "sql"
	JSON = "json"
	CSV  = "csv"
)

func List() []string {
	return []string{SQL, JSON, CSV}
}
