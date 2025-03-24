package model

// ArangoConnProperties is used for storing connection properties for database
type ArangoConnProperties struct {
	Host     string
	Port     int
	DbName   string
	Username string
	Password string
}
