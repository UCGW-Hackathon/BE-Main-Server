package services

import (
	"context"
	"errors"

	"situkang/models/entity"

	"gorm.io/gorm"
)

type KnowledgeService interface {
	ListArticles(ctx context.Context) (any, error)
	GetArticle(ctx context.Context, articleID string) (any, error)
	ListFAQ(ctx context.Context) (any, error)
}

type knowledgeService struct {
	db *gorm.DB
}

func NewKnowledgeService(db *gorm.DB) KnowledgeService {
	return &knowledgeService{db: db}
}

func (s *knowledgeService) ListArticles(ctx context.Context) (any, error) {
	var articles []entity.Article
	if err := s.db.WithContext(ctx).
		Where("is_published = TRUE").
		Order("published_at DESC NULLS LAST, created_at DESC").
		Find(&articles).Error; err != nil {
		return nil, err
	}

	data := make([]map[string]any, 0, len(articles))
	for _, article := range articles {
		data = append(data, articleListResponse(article))
	}
	return data, nil
}

func (s *knowledgeService) GetArticle(ctx context.Context, articleID string) (any, error) {
	var article entity.Article
	query := s.db.WithContext(ctx).Where("is_published = TRUE")
	if id, err := parseUUID(articleID); err == nil {
		query = query.Where("id = ?", id)
	} else {
		query = query.Where("slug = ?", articleID)
	}

	if err := query.First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return articleDetailResponse(article), nil
}

func (s *knowledgeService) ListFAQ(ctx context.Context) (any, error) {
	var faqs []entity.FAQ
	if err := s.db.WithContext(ctx).
		Where("is_active = TRUE").
		Order("display_order ASC, created_at DESC").
		Find(&faqs).Error; err != nil {
		return nil, err
	}

	data := make([]map[string]any, 0, len(faqs))
	for _, faq := range faqs {
		data = append(data, map[string]any{
			"faq_id":        faq.ID.String(),
			"question":      faq.Question,
			"answer":        faq.Answer,
			"category":      faq.Category,
			"display_order": faq.DisplayOrder,
		})
	}
	return data, nil
}

func articleListResponse(article entity.Article) map[string]any {
	return map[string]any{
		"article_id":        article.ID.String(),
		"title":             article.Title,
		"slug":              article.Slug,
		"category":          article.Category,
		"thumbnail_url":     article.ThumbnailURL,
		"excerpt":           article.Excerpt,
		"read_time_minutes": article.ReadTimeMinutes,
		"author":            article.Author,
		"tags":              article.Tags,
		"published_at":      article.PublishedAt,
	}
}

func articleDetailResponse(article entity.Article) map[string]any {
	data := articleListResponse(article)
	data["content_html"] = article.ContentHTML
	return data
}
