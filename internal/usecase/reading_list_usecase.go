package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"
)

type ReadingListUsecase interface {
	CreateReadingList(ctx context.Context, name, description string, userID uint, isPublic bool) (*entity.ReadingList, error)
	UpdateReadingList(ctx context.Context, id uint, name, description string, userID uint, isPublic bool) (*entity.ReadingList, error)
	DeleteReadingList(ctx context.Context, id, userID uint) error
	GetUserReadingLists(ctx context.Context, userID uint) ([]*entity.ReadingList, error)
	GetReadingList(ctx context.Context, id, userID uint) (*entity.ReadingList, error)
	GetPublicReadingLists(ctx context.Context) ([]*entity.ReadingList, error)
	AddArticleToList(ctx context.Context, readingListID, articleID, userID uint, notes string) (*entity.ReadingListItem, error)
	RemoveArticleFromList(ctx context.Context, readingListID, articleID, userID uint) error
	UpdateReadingListItem(ctx context.Context, readingListID, articleID, userID uint, notes string) (*entity.ReadingListItem, error)
	GetReadingListItems(ctx context.Context, readingListID, userID uint) ([]*entity.ReadingListItem, error)
}

type readingListUsecaseImpl struct {
	readingListRepo repository.ReadingListRepository
	articleRepo     repository.ArticleRepository
	userRepo        repository.UserRepository
}

func NewReadingListUsecase(
	readingListRepo repository.ReadingListRepository,
	articleRepo repository.ArticleRepository,
	userRepo repository.UserRepository,
) ReadingListUsecase {
	return &readingListUsecaseImpl{
		readingListRepo: readingListRepo,
		articleRepo:     articleRepo,
		userRepo:        userRepo,
	}
}

func (u *readingListUsecaseImpl) CreateReadingList(ctx context.Context, name, description string, userID uint, isPublic bool) (*entity.ReadingList, error) {
	// Validate input
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("reading list name is required")
	}
	if len(name) > 100 {
		return nil, errors.New("reading list name must be 100 characters or less")
	}
	if len(description) > 500 {
		return nil, errors.New("description must be 500 characters or less")
	}

	// Check if user exists
	if _, err := u.userRepo.GetByID(ctx, userID); err != nil {
		return nil, errors.New("user not found")
	}

	readingList := &entity.ReadingList{
		Name:        name,
		Description: description,
		UserID:      userID,
		IsPublic:    isPublic,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := u.readingListRepo.Create(ctx, readingList); err != nil {
		return nil, err
	}

	return readingList, nil
}

func (u *readingListUsecaseImpl) UpdateReadingList(ctx context.Context, id uint, name, description string, userID uint, isPublic bool) (*entity.ReadingList, error) {
	// Validate input
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("reading list name is required")
	}
	if len(name) > 100 {
		return nil, errors.New("reading list name must be 100 characters or less")
	}
	if len(description) > 500 {
		return nil, errors.New("description must be 500 characters or less")
	}

	// Get existing reading list
	readingList, err := u.readingListRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("reading list not found")
	}

	// Check ownership
	if readingList.UserID != userID {
		return nil, errors.New("unauthorized: you can only update your own reading lists")
	}

	// Update fields
	readingList.Name = name
	readingList.Description = description
	readingList.IsPublic = isPublic
	readingList.UpdatedAt = time.Now()

	if err := u.readingListRepo.Update(ctx, readingList); err != nil {
		return nil, err
	}

	return readingList, nil
}

func (u *readingListUsecaseImpl) DeleteReadingList(ctx context.Context, id, userID uint) error {
	// Get existing reading list
	readingList, err := u.readingListRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("reading list not found")
	}

	// Check ownership
	if readingList.UserID != userID {
		return errors.New("unauthorized: you can only delete your own reading lists")
	}

	return u.readingListRepo.Delete(ctx, id)
}

func (u *readingListUsecaseImpl) GetUserReadingLists(ctx context.Context, userID uint) ([]*entity.ReadingList, error) {
	// Check if user exists
	if _, err := u.userRepo.GetByID(ctx, userID); err != nil {
		return nil, errors.New("user not found")
	}

	return u.readingListRepo.GetByUserID(ctx, userID)
}

func (u *readingListUsecaseImpl) GetReadingList(ctx context.Context, id, userID uint) (*entity.ReadingList, error) {
	readingList, err := u.readingListRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("reading list not found")
	}

	// Check access permission
	if !readingList.IsPublic && readingList.UserID != userID {
		return nil, errors.New("unauthorized: this reading list is private")
	}

	return readingList, nil
}

func (u *readingListUsecaseImpl) GetPublicReadingLists(ctx context.Context) ([]*entity.ReadingList, error) {
	return u.readingListRepo.GetPublic(ctx)
}

func (u *readingListUsecaseImpl) AddArticleToList(ctx context.Context, readingListID, articleID, userID uint, notes string) (*entity.ReadingListItem, error) {
	// Validate notes length
	if len(notes) > 1000 {
		return nil, errors.New("notes must be 1000 characters or less")
	}

	// Get reading list and check ownership
	readingList, err := u.readingListRepo.GetByID(ctx, readingListID)
	if err != nil {
		return nil, errors.New("reading list not found")
	}
	if readingList.UserID != userID {
		return nil, errors.New("unauthorized: you can only add articles to your own reading lists")
	}

	// Check if article exists
	if _, err := u.articleRepo.GetByID(ctx, articleID); err != nil {
		return nil, errors.New("article not found")
	}

	// Check if article is already in the reading list
	items, err := u.readingListRepo.GetItems(ctx, readingListID)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.ArticleID == articleID {
			return nil, errors.New("article is already in this reading list")
		}
	}

	item := &entity.ReadingListItem{
		ReadingListID: readingListID,
		ArticleID:     articleID,
		AddedAt:       time.Now(),
		Notes:         notes,
	}

	if err := u.readingListRepo.AddItem(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (u *readingListUsecaseImpl) RemoveArticleFromList(ctx context.Context, readingListID, articleID, userID uint) error {
	// Get reading list and check ownership
	readingList, err := u.readingListRepo.GetByID(ctx, readingListID)
	if err != nil {
		return errors.New("reading list not found")
	}
	if readingList.UserID != userID {
		return errors.New("unauthorized: you can only remove articles from your own reading lists")
	}

	// Check if article is in the reading list
	items, err := u.readingListRepo.GetItems(ctx, readingListID)
	if err != nil {
		return err
	}
	found := false
	for _, item := range items {
		if item.ArticleID == articleID {
			found = true
			break
		}
	}
	if !found {
		return errors.New("article not found in this reading list")
	}

	return u.readingListRepo.RemoveItem(ctx, readingListID, articleID)
}

func (u *readingListUsecaseImpl) UpdateReadingListItem(ctx context.Context, readingListID, articleID, userID uint, notes string) (*entity.ReadingListItem, error) {
	// Validate notes length
	if len(notes) > 1000 {
		return nil, errors.New("notes must be 1000 characters or less")
	}

	// Get reading list and check ownership
	readingList, err := u.readingListRepo.GetByID(ctx, readingListID)
	if err != nil {
		return nil, errors.New("reading list not found")
	}
	if readingList.UserID != userID {
		return nil, errors.New("unauthorized: you can only update items in your own reading lists")
	}

	// Get the specific item
	items, err := u.readingListRepo.GetItems(ctx, readingListID)
	if err != nil {
		return nil, err
	}
	var targetItem *entity.ReadingListItem
	for _, item := range items {
		if item.ArticleID == articleID {
			targetItem = item
			break
		}
	}
	if targetItem == nil {
		return nil, errors.New("article not found in this reading list")
	}

	// Update notes
	targetItem.Notes = notes

	if err := u.readingListRepo.UpdateItem(ctx, targetItem); err != nil {
		return nil, err
	}

	return targetItem, nil
}

func (u *readingListUsecaseImpl) GetReadingListItems(ctx context.Context, readingListID, userID uint) ([]*entity.ReadingListItem, error) {
	// Get reading list and check access permission
	readingList, err := u.readingListRepo.GetByID(ctx, readingListID)
	if err != nil {
		return nil, errors.New("reading list not found")
	}
	if !readingList.IsPublic && readingList.UserID != userID {
		return nil, errors.New("unauthorized: this reading list is private")
	}

	return u.readingListRepo.GetItems(ctx, readingListID)
}