package tablemodes

const (
	Append = "append"
	Merge  = "merge"
)

func List() []string {
	return []string{Append, Merge}
}
