package repository

import (
	"github.com/senior-project-ai-content-tagging/content-preparing/entity"
)

type ContentRepository interface {
	UpdateContent(content entity.Content) (string, *entity.Content)
	FindContentById(contentId int64) (string, int64)
}

type ContentRepositorySqlx struct {
}

func (r *ContentRepositorySqlx) UpdateContent(content entity.Content) (string, *entity.Content) {
	return "UPDATE contents SET title_th = :title_th, content_th = :content_th, title_en = :title_en, content_en = :content_en WHERE id = :id", &content
}

func (r *ContentRepositorySqlx) FindContentById(contentId int64) (string, int64) {
	return "SELECT id, title_th, content_th, title_en, content_en FROM contents WHERE id=$1", contentId
}

func NewContentRepository() ContentRepository {
	return &ContentRepositorySqlx{}
}
