package usecase

import (
	"context"
	"layered-architecture-template/internal/domain/entity"
)

// MockTagRepository implements repository.TagRepository for testing
type MockTagRepository struct {
	tags   []*entity.Tag
	nextID uint
}

func NewMockTagRepository() *MockTagRepository {
	return &MockTagRepository{
		tags:   []*entity.Tag{},
		nextID: 1,
	}
}

func (m *MockTagRepository) Create(ctx context.Context, tag *entity.Tag) error {
	tag.ID = m.nextID
	m.nextID++
	m.tags = append(m.tags, tag)
	return nil
}

func (m *MockTagRepository) GetByID(ctx context.Context, id uint) (*entity.Tag, error) {
	for _, tag := range m.tags {
		if tag.ID == id {
			return tag, nil
		}
	}
	return nil, &entity.ValidationError{Message: "tag not found"}
}

func (m *MockTagRepository) GetBySlug(ctx context.Context, slug string) (*entity.Tag, error) {
	for _, tag := range m.tags {
		if tag.Slug == slug {
			return tag, nil
		}
	}
	return nil, &entity.ValidationError{Message: "tag not found"}
}

func (m *MockTagRepository) GetByNames(ctx context.Context, names []string) ([]*entity.Tag, error) {
	var result []*entity.Tag
	for _, tag := range m.tags {
		for _, name := range names {
			if tag.Name == name {
				result = append(result, tag)
				break
			}
		}
	}
	return result, nil
}

func (m *MockTagRepository) GetAll(ctx context.Context) ([]*entity.Tag, error) {
	return m.tags, nil
}

func (m *MockTagRepository) Update(ctx context.Context, tag *entity.Tag) error {
	for i, t := range m.tags {
		if t.ID == tag.ID {
			m.tags[i] = tag
			return nil
		}
	}
	return &entity.ValidationError{Message: "tag not found"}
}

func (m *MockTagRepository) Delete(ctx context.Context, id uint) error {
	for i, tag := range m.tags {
		if tag.ID == id {
			m.tags = append(m.tags[:i], m.tags[i+1:]...)
			return nil
		}
	}
	return &entity.ValidationError{Message: "tag not found"}
}

func (m *MockTagRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	for _, tag := range m.tags {
		if tag.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func (m *MockTagRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	for _, tag := range m.tags {
		if tag.Slug == slug {
			return true, nil
		}
	}
	return false, nil
}

func (m *MockTagRepository) CreateMultiple(ctx context.Context, tags []*entity.Tag) error {
	for _, tag := range tags {
		if err := m.Create(ctx, tag); err != nil {
			return err
		}
	}
	return nil
}

func (m *MockTagRepository) GetPopular(ctx context.Context, limit int) ([]*entity.TagWithCount, error) {
	// Mock implementation - return empty result
	return []*entity.TagWithCount{}, nil
}