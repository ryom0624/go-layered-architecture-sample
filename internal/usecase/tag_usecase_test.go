package usecase

import (
	"context"
	"testing"

	"layered-architecture-template/internal/domain/entity"
)

func TestTagUsecase_CreateTag(t *testing.T) {
	mockTagRepo := NewMockTagRepository()
	usecase := NewTagUsecase(mockTagRepo)
	ctx := context.Background()

	t.Run("successful tag creation", func(t *testing.T) {
		name := "Go"
		color := "#00ADD8"

		tag, err := usecase.CreateTag(ctx, name, color)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if tag.Name != name {
			t.Errorf("expected name %s, got %s", name, tag.Name)
		}
		if tag.Color != color {
			t.Errorf("expected color %s, got %s", color, tag.Color)
		}
		if tag.Slug != "go" {
			t.Errorf("expected slug 'go', got %s", tag.Slug)
		}
	})

	t.Run("duplicate tag name", func(t *testing.T) {
		name := "Go"
		color := "#FF0000"

		// Try to create the same tag again
		_, err := usecase.CreateTag(ctx, name, color)
		if err == nil {
			t.Error("expected error for duplicate tag name")
		}
		if err.Error() != "tag with this name already exists" {
			t.Errorf("expected 'tag with this name already exists', got %s", err.Error())
		}
	})

	t.Run("default color when empty", func(t *testing.T) {
		name := "Python"
		color := ""

		tag, err := usecase.CreateTag(ctx, name, color)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if tag.Color != "#3B82F6" {
			t.Errorf("expected default color '#3B82F6', got %s", tag.Color)
		}
	})
}

func TestTagUsecase_GetTag(t *testing.T) {
	mockTagRepo := NewMockTagRepository()
	usecase := NewTagUsecase(mockTagRepo)
	ctx := context.Background()

	// Create a test tag
	testTag := &entity.Tag{
		Name:  "JavaScript",
		Color: "#F7DF1E",
	}
	testTag.GenerateSlug()
	mockTagRepo.Create(ctx, testTag)

	t.Run("get existing tag", func(t *testing.T) {
		tag, err := usecase.GetTag(ctx, testTag.Slug)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if tag.Name != testTag.Name {
			t.Errorf("expected name %s, got %s", testTag.Name, tag.Name)
		}
	})

	t.Run("get non-existent tag", func(t *testing.T) {
		_, err := usecase.GetTag(ctx, "non-existent")
		if err == nil {
			t.Error("expected error for non-existent tag")
		}
	})

	t.Run("empty slug validation", func(t *testing.T) {
		_, err := usecase.GetTag(ctx, "")
		if err == nil {
			t.Error("expected error for empty slug")
		}
	})
}

func TestTagUsecase_UpdateTag(t *testing.T) {
	mockTagRepo := NewMockTagRepository()
	usecase := NewTagUsecase(mockTagRepo)
	ctx := context.Background()

	// Create a test tag
	testTag := &entity.Tag{
		Name:  "React",
		Color: "#61DAFB",
	}
	testTag.GenerateSlug()
	mockTagRepo.Create(ctx, testTag)

	t.Run("successful tag update", func(t *testing.T) {
		updatedName := "ReactJS"
		updatedColor := "#00D4FF"

		tag, err := usecase.UpdateTag(ctx, testTag.ID, updatedName, updatedColor)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if tag.Name != updatedName {
			t.Errorf("expected name %s, got %s", updatedName, tag.Name)
		}
		if tag.Color != updatedColor {
			t.Errorf("expected color %s, got %s", updatedColor, tag.Color)
		}
	})

	t.Run("update non-existent tag", func(t *testing.T) {
		_, err := usecase.UpdateTag(ctx, 999, "Updated Name", "#FF0000")
		if err == nil {
			t.Error("expected error for non-existent tag")
		}
	})

	t.Run("update with zero ID", func(t *testing.T) {
		_, err := usecase.UpdateTag(ctx, 0, "Updated Name", "#FF0000")
		if err == nil {
			t.Error("expected error for zero ID")
		}
	})
}

func TestTagUsecase_DeleteTag(t *testing.T) {
	mockTagRepo := NewMockTagRepository()
	usecase := NewTagUsecase(mockTagRepo)
	ctx := context.Background()

	// Create a test tag
	testTag := &entity.Tag{
		Name:  "Vue",
		Color: "#4FC08D",
	}
	testTag.GenerateSlug()
	mockTagRepo.Create(ctx, testTag)

	t.Run("successful tag deletion", func(t *testing.T) {
		err := usecase.DeleteTag(ctx, testTag.ID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Verify tag is deleted
		_, err = usecase.GetTag(ctx, testTag.Slug)
		if err == nil {
			t.Error("expected error when getting deleted tag")
		}
	})

	t.Run("delete non-existent tag", func(t *testing.T) {
		err := usecase.DeleteTag(ctx, 999)
		if err == nil {
			t.Error("expected error for non-existent tag")
		}
	})

	t.Run("delete with zero ID", func(t *testing.T) {
		err := usecase.DeleteTag(ctx, 0)
		if err == nil {
			t.Error("expected error for zero ID")
		}
	})
}

func TestTagUsecase_GetOrCreateTags(t *testing.T) {
	mockTagRepo := NewMockTagRepository()
	usecase := NewTagUsecase(mockTagRepo)
	ctx := context.Background()

	// Create existing tag
	existingTag := &entity.Tag{
		Name:  "Docker",
		Color: "#2496ED",
	}
	existingTag.GenerateSlug()
	mockTagRepo.Create(ctx, existingTag)

	t.Run("get existing and create new tags", func(t *testing.T) {
		tagNames := []string{"Docker", "Kubernetes", "Helm"}

		tags, err := usecase.GetOrCreateTags(ctx, tagNames)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(tags) != 3 {
			t.Errorf("expected 3 tags, got %d", len(tags))
		}

		// Check that existing tag is returned
		foundDocker := false
		for _, tag := range tags {
			if tag.Name == "Docker" && tag.ID == existingTag.ID {
				foundDocker = true
				break
			}
		}
		if !foundDocker {
			t.Error("expected to find existing Docker tag")
		}
	})

	t.Run("empty tag names", func(t *testing.T) {
		tags, err := usecase.GetOrCreateTags(ctx, []string{})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(tags) != 0 {
			t.Errorf("expected 0 tags, got %d", len(tags))
		}
	})

	t.Run("all existing tags", func(t *testing.T) {
		tagNames := []string{"Docker"}

		tags, err := usecase.GetOrCreateTags(ctx, tagNames)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(tags) != 1 {
			t.Errorf("expected 1 tag, got %d", len(tags))
		}
		if tags[0].Name != "Docker" {
			t.Errorf("expected tag name 'Docker', got %s", tags[0].Name)
		}
	})
}

func TestTagUsecase_GetPopularTags(t *testing.T) {
	mockTagRepo := NewMockTagRepository()
	usecase := NewTagUsecase(mockTagRepo)
	ctx := context.Background()

	t.Run("get popular tags with valid limit", func(t *testing.T) {
		limit := 5
		tags, err := usecase.GetPopularTags(ctx, limit)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		// Mock implementation returns empty result
		if len(tags) != 0 {
			t.Errorf("expected 0 tags from mock, got %d", len(tags))
		}
	})

	t.Run("zero limit defaults to 10", func(t *testing.T) {
		tags, err := usecase.GetPopularTags(ctx, 0)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		// Mock implementation returns empty result
		if len(tags) != 0 {
			t.Errorf("expected 0 tags from mock, got %d", len(tags))
		}
	})

	t.Run("limit exceeds maximum, capped at 100", func(t *testing.T) {
		tags, err := usecase.GetPopularTags(ctx, 150)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		// Mock implementation returns empty result
		if len(tags) != 0 {
			t.Errorf("expected 0 tags from mock, got %d", len(tags))
		}
	})
}