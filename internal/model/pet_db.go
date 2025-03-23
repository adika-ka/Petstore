package model

type PetDB struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	Status    string `db:"status"`
	PhotoUrls string `db:"photo_urls"`
	Tags      string `db:"tags"`
	Category  string `db:"category"`
}
