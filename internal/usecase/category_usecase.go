package usecase

import (
	"context"
	"errors"
	"fmt"
	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"
	"strconv"
)

type CategoryUsecase interface {
	CreateCategory(ctx context.Context, name, description string) (*entity.Category, error)
	GetCategory(ctx context.Context, slug string) (*entity.Category, error)
	GetAllCategories(ctx context.Context) ([]*entity.Category, error)
	UpdateCategory(ctx context.Context, id uint, name, description string) (*entity.Category, error)
	DeleteCategory(ctx context.Context, id uint) error
	GetCategoriesWithCount(ctx context.Context) ([]*entity.CategoryWithCount, error)
}

type categoryUsecase struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryUsecase(categoryRepo repository.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{
		categoryRepo: categoryRepo,
	}
}

func (u *categoryUsecase) CreateCategory(ctx context.Context, name, description string) (*entity.Category, error) {
	category := &entity.Category{
		Name:        name,
		Description: description,
	}

	// Validate category
	if err := category.Validate(); err != nil {
		return nil, err
	}

	// Generate slug
	category.GenerateSlug()

	// Check if name already exists
	exists, err := u.categoryRepo.ExistsByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to check if category name exists: %w", err)
	}
	if exists {
		return nil, errors.New("category with this name already exists")
	}

	// Handle slug conflicts
	originalSlug := category.Slug
	counter := 1
	for {
		exists, err := u.categoryRepo.ExistsBySlug(ctx, category.Slug)
		if err != nil {
			return nil, fmt.Errorf("failed to check if category slug exists: %w", err)
		}
		if !exists {
			break
		}
		category.Slug = originalSlug + "-" + strconv.Itoa(counter)
		counter++
	}

	// Create category
	if err := u.categoryRepo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

func (u *categoryUsecase) GetCategory(ctx context.Context, slug string) (*entity.Category, error) {
	if slug == "" {
		return nil, errors.New("slug is required")
	}

	category, err := u.categoryRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	return category, nil
}

func (u *categoryUsecase) GetAllCategories(ctx context.Context) ([]*entity.Category, error) {
	categories, err := u.categoryRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return categories, nil
}

func (u *categoryUsecase) UpdateCategory(ctx context.Context, id uint, name, description string) (*entity.Category, error) {
	if id == 0 {
		return nil, errors.New("category ID is required")
	}

	// Get existing category
	category, err := u.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	// Update fields
	if name != "" {
		category.Name = name
		category.GenerateSlug()
	}
	if description != "" {
		category.Description = description
	}

	// Validate updated category
	if err := category.Validate(); err != nil {
		return nil, err
	}

	// Check for name conflicts (excluding current category)
	if name != "" {
		exists, err := u.categoryRepo.ExistsByName(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("failed to check if category name exists: %w", err)
		}
		if exists {
			// Check if it's a different category with the same name
			existingCategory, err := u.categoryRepo.GetBySlug(ctx, category.Slug)
			if err == nil && existingCategory.ID != id {
				return nil, errors.New("category with this name already exists")
			}
		}
	}

	// Handle slug conflicts (excluding current category)
	if name != "" {
		originalSlug := category.Slug
		counter := 1
		for {
			exists, err := u.categoryRepo.ExistsBySlug(ctx, category.Slug)
			if err != nil {
				return nil, fmt.Errorf("failed to check if category slug exists: %w", err)
			}
			if !exists {
				break
			}
			// Check if it's the same category
			existingCategory, err := u.categoryRepo.GetBySlug(ctx, category.Slug)
			if err == nil && existingCategory.ID == id {
				break
			}
			category.Slug = originalSlug + "-" + strconv.Itoa(counter)
			counter++
		}
	}

	// Update category
	if err := u.categoryRepo.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

func (u *categoryUsecase) DeleteCategory(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("category ID is required")
	}

	// Check if category exists
	_, err := u.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("category not found: %w", err)
	}

	// Delete category (articles will have their category_id set to NULL due to foreign key constraints)
	if err := u.categoryRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

func (u *categoryUsecase) GetCategoriesWithCount(ctx context.Context) ([]*entity.CategoryWithCount, error) {
	categories, err := u.categoryRepo.GetWithArticleCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories with count: %w", err)
	}

	return categories, nil
}