package database

type Model interface {
	TableName() string
	// Scan(any) error
	FetchQuery() string
}
