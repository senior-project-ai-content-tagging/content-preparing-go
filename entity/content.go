package entity

type Content struct {
	Id        int64  `db:"id"`
	TitleTH   string `db:"title_th"`
	ContentTH string `db:"content_th"`
	TitleEN   string `db:"title_en"`
	ContentEN string `db:"content_en"`
}
