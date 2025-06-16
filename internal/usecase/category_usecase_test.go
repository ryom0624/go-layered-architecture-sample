package usecase

import (
	"context"
	"testing"

	"layered-architecture-template/internal/domain/entity"
)

func TestCategoryUsecase_CreateCategory(t *testing.T) {
	mockCategoryRepo := NewMockCategoryRepository()
	usecase := NewCategoryUsecase(mockCategoryRepo)
	ctx := context.Background()

	t.Run("successful category creation", func(t *testing.T) {
		name := "Technology"
		description := "Technology articles"

		category, err := usecase.CreateCategory(ctx, name, description)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if category.Name != name {
			t.Errorf("expected name %s, got %s", name, category.Name)
		}
		if category.Description != description {
			t.Errorf("expected description %s, got %s", description, category.Description)
		}
		if category.Slug != "technology" {
			t.Errorf("expected slug 'technology', got %s", category.Slug)
		}
	})

	t.Run("duplicate category name", func(t *testing.T) {
		name := "Technology"
		description := "Another tech category"

		// Try to create the same category again
		_, err := usecase.CreateCategory(ctx, name, description)
		if err == nil {
			t.Error("expected error for duplicate category name")
		}
		if err.Error() != "category with this name already exists" {
			t.Errorf("expected 'category with this name already exists', got %s", err.Error())
		}
	})

	t.Run("empty name validation", func(t *testing.T) {
		_, err := usecase.CreateCategory(ctx, "", "Description")
		if err == nil {
			t.Error("expected error for empty name")
		}
	})
}

func TestCategoryUsecase_GetCategory(t *testing.T) {
	mockCategoryRepo := NewMockCategoryRepository()
	usecase := NewCategoryUsecase(mockCategoryRepo)
	ctx := context.Background()

	// Create a test category
	testCategory := &entity.Category{
		Name:        "Science",
		Description: "Science articles",
	}
	testCategory.GenerateSlug()
	mockCategoryRepo.Create(ctx, testCategory)

	t.Run("get existing category", func(t *testing.T) {
		category, err := usecase.GetCategory(ctx, testCategory.Slug)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if category.Name != testCategory.Name {
			t.Errorf("expected name %s, got %s", testCategory.Name, category.Name)
		}
	})

	t.Run("get non-existent category", func(t *testing.T) {
		_, err := usecase.GetCategory(ctx, "non-existent")
		if err == nil {
			t.Error("expected error for non-existent category")
		}
	})

	t.Run("empty slug validation", func(t *testing.T) {
		_, err := usecase.GetCategory(ctx, "")
		if err == nil {
			t.Error("expected error for empty slug")
		}
	})
}

func TestCategoryUsecase_UpdateCategory(t *testing.T) {
	mockCategoryRepo := NewMockCategoryRepository()
	usecase := NewCategoryUsecase(mockCategoryRepo)
	ctx := context.Background()

	// Create a test category
	testCategory := &entity.Category{
		Name:        "Programming",
		Description: "Programming articles",
	}
	testCategory.GenerateSlug()
	mockCategoryRepo.Create(ctx, testCategory)

	t.Run("successful category update", func(t *testing.T) {
		updatedName := "Software Development"
		updatedDescription := "Software development articles"

		category, err := usecase.UpdateCategory(ctx, testCategory.ID, updatedName, updatedDescription)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if category.Name != updatedName {
			t.Errorf("expected name %s, got %s", updatedName, category.Name)
		}
		if category.Description != updatedDescription {
			t.Errorf("expected description %s, got %s", updatedDescription, category.Description)
		}
	})

	t.Run("update non-existent category", func(t *testing.T) {
		_, err := usecase.UpdateCategory(ctx, 999, "Updated Name", "Updated Description")
		if err == nil {
			t.Error("expected error for non-existent category")
		}
	})

	t.Run("update with zero ID", func(t *testing.T) {
		_, err := usecase.UpdateCategory(ctx, 0, "Updated Name", "Updated Description")
		if err == nil {
			t.Error("expected error for zero ID")
		}
	})
}

func TestCategoryUsecase_DeleteCategory(t *testing.T) {
	mockCategoryRepo := NewMockCategoryRepository()
	usecase := NewCategoryUsecase(mockCategoryRepo)
	ctx := context.Background()

	// Create a test category
	testCategory := &entity.Category{
		Name:        "Web Development",
		Description: "Web development articles",
	}
	testCategory.GenerateSlug()
	mockCategoryRepo.Create(ctx, testCategory)

	t.Run("successful category deletion", func(t *testing.T) {
		err := usecase.DeleteCategory(ctx, testCategory.ID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Verify category is deleted
		_, err = usecase.GetCategory(ctx, testCategory.Slug)
		if err == nil {
			t.Error("expected error when getting deleted category")
		}
	})

	t.Run("delete non-existent category", func(t *testing.T) {
		err := usecase.DeleteCategory(ctx, 999)
		if err == nil {
			t.Error("expected error for non-existent category")
		}
	})

	t.Run("delete with zero ID", func(t *testing.T) {
		err := usecase.DeleteCategory(ctx, 0)
		if err == nil {
			t.Error("expected error for zero ID")
		}
	})
}

func TestCategoryUsecase_GetAllCategories(t *testing.T) {
	mockCategoryRepo := NewMockCategoryRepository()
	usecase := NewCategoryUsecase(mockCategoryRepo)
	ctx := context.Background()

	t.Run("get all categories when empty", func(t *testing.T) {
		categories, err := usecase.GetAllCategories(ctx)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(categories) != 0 {
			t.Errorf("expected 0 categories, got %d", len(categories))
		}
	})

	t.Run("get all categories with data", func(t *testing.T) {
		// Create test categories
		category1 := &entity.Category{Name: "AI", Description: "AI articles"}
		category1.GenerateSlug()
		category2 := &entity.Category{Name: "Blockchain", Description: "Blockchain articles"}
		category2.GenerateSlug()

		mockCategoryRepo.Create(ctx, category1)
		mockCategoryRepo.Create(ctx, category2)

		categories, err := usecase.GetAllCategories(ctx)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(categories) != 2 {
			t.Errorf("expected 2 categories, got %d", len(categories))
		}
	})
}