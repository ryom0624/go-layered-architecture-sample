package usecase

import (
	"context"
	"errors"
	"fmt"
	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"
	"strconv"
)

type TagUsecase interface {
	CreateTag(ctx context.Context, name, color string) (*entity.Tag, error)
	GetTag(ctx context.Context, slug string) (*entity.Tag, error)
	GetAllTags(ctx context.Context) ([]*entity.Tag, error)
	UpdateTag(ctx context.Context, id uint, name, color string) (*entity.Tag, error)
	DeleteTag(ctx context.Context, id uint) error
	GetOrCreateTags(ctx context.Context, tagNames []string) ([]*entity.Tag, error)
	GetPopularTags(ctx context.Context, limit int) ([]*entity.TagWithCount, error)
}

type tagUsecase struct {
	tagRepo repository.TagRepository
}

func NewTagUsecase(tagRepo repository.TagRepository) TagUsecase {
	return &tagUsecase{
		tagRepo: tagRepo,
	}
}

func (u *tagUsecase) CreateTag(ctx context.Context, name, color string) (*entity.Tag, error) {
	tag := &entity.Tag{
		Name:  name,
		Color: color,
	}

	// Validate tag
	if err := tag.Validate(); err != nil {
		return nil, err
	}

	// Generate slug
	tag.GenerateSlug()

	// Check if name already exists
	exists, err := u.tagRepo.ExistsByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to check if tag name exists: %w", err)
	}
	if exists {
		return nil, errors.New("tag with this name already exists")
	}

	// Handle slug conflicts
	originalSlug := tag.Slug
	counter := 1
	for {
		exists, err := u.tagRepo.ExistsBySlug(ctx, tag.Slug)
		if err != nil {
			return nil, fmt.Errorf("failed to check if tag slug exists: %w", err)
		}
		if !exists {
			break
		}
		tag.Slug = originalSlug + "-" + strconv.Itoa(counter)
		counter++
	}

	// Create tag
	if err := u.tagRepo.Create(ctx, tag); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return tag, nil
}

func (u *tagUsecase) GetTag(ctx context.Context, slug string) (*entity.Tag, error) {
	if slug == "" {
		return nil, errors.New("slug is required")
	}

	tag, err := u.tagRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("tag not found: %w", err)
	}

	return tag, nil
}

func (u *tagUsecase) GetAllTags(ctx context.Context) ([]*entity.Tag, error) {
	tags, err := u.tagRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	return tags, nil
}

func (u *tagUsecase) UpdateTag(ctx context.Context, id uint, name, color string) (*entity.Tag, error) {
	if id == 0 {
		return nil, errors.New("tag ID is required")
	}

	// Get existing tag
	tag, err := u.tagRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("tag not found: %w", err)
	}

	// Update fields
	if name != "" {
		tag.Name = name
		tag.GenerateSlug()
	}
	if color != "" {
		tag.Color = color
	}

	// Validate updated tag
	if err := tag.Validate(); err != nil {
		return nil, err
	}

	// Check for name conflicts (excluding current tag)
	if name != "" {
		exists, err := u.tagRepo.ExistsByName(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("failed to check if tag name exists: %w", err)
		}
		if exists {
			// Check if it's a different tag with the same name
			existingTag, err := u.tagRepo.GetBySlug(ctx, tag.Slug)
			if err == nil && existingTag.ID != id {
				return nil, errors.New("tag with this name already exists")
			}
		}
	}

	// Handle slug conflicts (excluding current tag)
	if name != "" {
		originalSlug := tag.Slug
		counter := 1
		for {
			exists, err := u.tagRepo.ExistsBySlug(ctx, tag.Slug)
			if err != nil {
				return nil, fmt.Errorf("failed to check if tag slug exists: %w", err)
			}
			if !exists {
				break
			}
			// Check if it's the same tag
			existingTag, err := u.tagRepo.GetBySlug(ctx, tag.Slug)
			if err == nil && existingTag.ID == id {
				break
			}
			tag.Slug = originalSlug + "-" + strconv.Itoa(counter)
			counter++
		}
	}

	// Update tag
	if err := u.tagRepo.Update(ctx, tag); err != nil {
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	return tag, nil
}

func (u *tagUsecase) DeleteTag(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("tag ID is required")
	}

	// Check if tag exists
	_, err := u.tagRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("tag not found: %w", err)
	}

	// Delete tag (article_tags relationships will be deleted due to foreign key constraints)
	if err := u.tagRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

func (u *tagUsecase) GetOrCreateTags(ctx context.Context, tagNames []string) ([]*entity.Tag, error) {
	if len(tagNames) == 0 {
		return []*entity.Tag{}, nil
	}

	// Get existing tags
	existingTags, err := u.tagRepo.GetByNames(ctx, tagNames)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing tags: %w", err)
	}

	// Create map of existing tag names
	existingTagNames := make(map[string]*entity.Tag)
	for _, tag := range existingTags {
		existingTagNames[tag.Name] = tag
	}

	// Identify missing tags
	var missingTagNames []string
	for _, name := range tagNames {
		if _, exists := existingTagNames[name]; !exists {
			missingTagNames = append(missingTagNames, name)
		}
	}

	// Create missing tags
	var newTags []*entity.Tag
	for _, name := range missingTagNames {
		tag := &entity.Tag{
			Name:  name,
			Color: "#3B82F6", // Default blue color
		}

		// Validate tag
		if err := tag.Validate(); err != nil {
			return nil, fmt.Errorf("invalid tag '%s': %w", name, err)
		}

		// Generate slug
		tag.GenerateSlug()

		// Handle slug conflicts
		originalSlug := tag.Slug
		counter := 1
		for {
			exists, err := u.tagRepo.ExistsBySlug(ctx, tag.Slug)
			if err != nil {
				return nil, fmt.Errorf("failed to check if tag slug exists: %w", err)
			}
			if !exists {
				break
			}
			tag.Slug = originalSlug + "-" + strconv.Itoa(counter)
			counter++
		}

		newTags = append(newTags, tag)
	}

	// Create new tags in batch
	if len(newTags) > 0 {
		if err := u.tagRepo.CreateMultiple(ctx, newTags); err != nil {
			return nil, fmt.Errorf("failed to create new tags: %w", err)
		}
	}

	// Combine existing and new tags
	var allTags []*entity.Tag
	for _, name := range tagNames {
		if existingTag, exists := existingTagNames[name]; exists {
			allTags = append(allTags, existingTag)
		} else {
			// Find the new tag
			for _, newTag := range newTags {
				if newTag.Name == name {
					allTags = append(allTags, newTag)
					break
				}
			}
		}
	}

	return allTags, nil
}

func (u *tagUsecase) GetPopularTags(ctx context.Context, limit int) ([]*entity.TagWithCount, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	tags, err := u.tagRepo.GetPopular(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular tags: %w", err)
	}

	return tags, nil
}