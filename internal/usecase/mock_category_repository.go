package usecase

import (
	"context"
	"layered-architecture-template/internal/domain/entity"
)

// MockCategoryRepository implements repository.CategoryRepository for testing
type MockCategoryRepository struct {
	categories []*entity.Category
	nextID     uint
}

func NewMockCategoryRepository() *MockCategoryRepository {
	return &MockCategoryRepository{
		categories: []*entity.Category{},
		nextID:     1,
	}
}

func (m *MockCategoryRepository) Create(ctx context.Context, category *entity.Category) error {
	category.ID = m.nextID
	m.nextID++
	m.categories = append(m.categories, category)
	return nil
}

func (m *MockCategoryRepository) GetByID(ctx context.Context, id uint) (*entity.Category, error) {
	for _, category := range m.categories {
		if category.ID == id {
			return category, nil
		}
	}
	return nil, &entity.ValidationError{Message: "category not found"}
}

func (m *MockCategoryRepository) GetBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	for _, category := range m.categories {
		if category.Slug == slug {
			return category, nil
		}
	}
	return nil, &entity.ValidationError{Message: "category not found"}
}

func (m *MockCategoryRepository) GetAll(ctx context.Context) ([]*entity.Category, error) {
	return m.categories, nil
}

func (m *MockCategoryRepository) Update(ctx context.Context, category *entity.Category) error {
	for i, c := range m.categories {
		if c.ID == category.ID {
			m.categories[i] = category
			return nil
		}
	}
	return &entity.ValidationError{Message: "category not found"}
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id uint) error {
	for i, category := range m.categories {
		if category.ID == id {
			m.categories = append(m.categories[:i], m.categories[i+1:]...)
			return nil
		}
	}
	return &entity.ValidationError{Message: "category not found"}
}

func (m *MockCategoryRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	for _, category := range m.categories {
		if category.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func (m *MockCategoryRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	for _, category := range m.categories {
		if category.Slug == slug {
			return true, nil
		}
	}
	return false, nil
}

func (m *MockCategoryRepository) GetWithArticleCount(ctx context.Context) ([]*entity.CategoryWithCount, error) {
	result := make([]*entity.CategoryWithCount, len(m.categories))
	for i, category := range m.categories {
		result[i] = &entity.CategoryWithCount{
			Category:     category,
			ArticleCount: 0, // Mock implementation - always return 0
		}
	}
	return result, nil
}