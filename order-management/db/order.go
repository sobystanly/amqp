package db

type (
	orderDB struct {
		db *db
	}
)

func NewOrderDB(db *db) *orderDB {
	return &orderDB{db: db}
}
