package usecase

import (
	"context"
	"errors"
	"testing"

	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"
)

// Mock repositories
type mockCommentRepository struct {
	comments map[uint]*entity.Comment
	nextID   uint
}

func newMockCommentRepository() *mockCommentRepository {
	return &mockCommentRepository{
		comments: make(map[uint]*entity.Comment),
		nextID:   1,
	}
}

func (m *mockCommentRepository) Create(ctx context.Context, comment *entity.Comment) error {
	comment.ID = m.nextID
	m.nextID++
	m.comments[comment.ID] = comment
	return nil
}

func (m *mockCommentRepository) CreateWithTx(ctx context.Context, tx repository.Transaction, comment *entity.Comment) error {
	return m.Create(ctx, comment)
}

func (m *mockCommentRepository) GetByID(ctx context.Context, id uint) (*entity.Comment, error) {
	comment, exists := m.comments[id]
	if !exists {
		return nil, errors.New("comment not found")
	}
	return comment, nil
}

func (m *mockCommentRepository) Update(ctx context.Context, comment *entity.Comment) error {
	if _, exists := m.comments[comment.ID]; !exists {
		return errors.New("comment not found")
	}
	m.comments[comment.ID] = comment
	return nil
}

func (m *mockCommentRepository) UpdateWithTx(ctx context.Context, tx repository.Transaction, comment *entity.Comment) error {
	return m.Update(ctx, comment)
}

func (m *mockCommentRepository) Delete(ctx context.Context, id uint) error {
	if _, exists := m.comments[id]; !exists {
		return errors.New("comment not found")
	}
	delete(m.comments, id)
	return nil
}

func (m *mockCommentRepository) GetByArticleID(ctx context.Context, articleID uint) ([]*entity.Comment, error) {
	var comments []*entity.Comment
	for _, comment := range m.comments {
		if comment.ArticleID == articleID {
			comments = append(comments, comment)
		}
	}
	return comments, nil
}

func (m *mockCommentRepository) GetReplies(ctx context.Context, parentID uint) ([]*entity.Comment, error) {
	var replies []*entity.Comment
	for _, comment := range m.comments {
		if comment.ParentID != nil && *comment.ParentID == parentID {
			replies = append(replies, comment)
		}
	}
	return replies, nil
}

func (m *mockCommentRepository) GetByStatus(ctx context.Context, status string) ([]*entity.Comment, error) {
	var comments []*entity.Comment
	for _, comment := range m.comments {
		if comment.Status == status {
			comments = append(comments, comment)
		}
	}
	return comments, nil
}

func (m *mockCommentRepository) GetByUserID(ctx context.Context, userID uint) ([]*entity.Comment, error) {
	var comments []*entity.Comment
	for _, comment := range m.comments {
		if comment.AuthorID == userID {
			comments = append(comments, comment)
		}
	}
	return comments, nil
}

func (m *mockCommentRepository) GetCommentsWithReplies(ctx context.Context, articleID uint) ([]*entity.Comment, error) {
	return m.GetByArticleID(ctx, articleID)
}

func (m *mockCommentRepository) GetCommentDepth(ctx context.Context, commentID uint) (int, error) {
	return 0, nil // Simplified for testing
}

func (m *mockCommentRepository) CountByArticleID(ctx context.Context, articleID uint) (int64, error) {
	count := int64(0)
	for _, comment := range m.comments {
		if comment.ArticleID == articleID {
			count++
		}
	}
	return count, nil
}

func (m *mockCommentRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	count := int64(0)
	for _, comment := range m.comments {
		if comment.Status == status {
			count++
		}
	}
	return count, nil
}

type mockUserRepository struct {
	users map[uint]*entity.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: map[uint]*entity.User{
			1: {ID: 1, Name: "Test User", Email: "test@example.com"},
		},
	}
}

func (m *mockUserRepository) Create(ctx context.Context, user *entity.User) error {
	return nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	return nil, nil
}

func (m *mockUserRepository) GetAll(ctx context.Context) ([]*entity.User, error) {
	return nil, nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *entity.User) error {
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id uint) error {
	return nil
}

type mockArticleRepository struct {
	articles map[uint]*entity.Article
}

func newMockArticleRepository() *mockArticleRepository {
	return &mockArticleRepository{
		articles: map[uint]*entity.Article{
			1: {ID: 1, Title: "Test Article", Content: "Test Content", Status: "published"},
		},
	}
}

func (m *mockArticleRepository) Create(ctx context.Context, article *entity.Article) error {
	return nil
}

func (m *mockArticleRepository) CreateWithTx(ctx context.Context, tx repository.Transaction, article *entity.Article) error {
	return nil
}

func (m *mockArticleRepository) GetByID(ctx context.Context, id uint) (*entity.Article, error) {
	article, exists := m.articles[id]
	if !exists {
		return nil, errors.New("article not found")
	}
	return article, nil
}

func (m *mockArticleRepository) GetAll(ctx context.Context) ([]*entity.Article, error) {
	return nil, nil
}

func (m *mockArticleRepository) GetByStatus(ctx context.Context, status string) ([]*entity.Article, error) {
	return nil, nil
}

func (m *mockArticleRepository) Update(ctx context.Context, article *entity.Article) error {
	return nil
}

func (m *mockArticleRepository) UpdateWithTx(ctx context.Context, tx repository.Transaction, article *entity.Article) error {
	return nil
}

func (m *mockArticleRepository) Delete(ctx context.Context, id uint) error {
	return nil
}

func (m *mockArticleRepository) DeleteWithTx(ctx context.Context, tx repository.Transaction, id uint) error {
	return nil
}

type mockTransactionManager struct{}

func (m *mockTransactionManager) WithTransaction(ctx context.Context, fn func(tx repository.Transaction) error) error {
	return fn(&mockTransaction{})
}

type mockTransaction struct{}

func (m *mockTransaction) GetDB() interface{} {
	return nil
}

func (m *mockTransaction) Commit() error {
	return nil
}

func (m *mockTransaction) Rollback() error {
	return nil
}

func TestCommentUsecase_CreateComment(t *testing.T) {
	commentRepo := newMockCommentRepository()
	userRepo := newMockUserRepository()
	articleRepo := newMockArticleRepository()
	txManager := &mockTransactionManager{}

	usecase := NewCommentUsecase(commentRepo, userRepo, articleRepo, txManager)

	tests := []struct {
		name      string
		content   string
		authorID  uint
		articleID uint
		wantError bool
	}{
		{
			name:      "Valid comment creation",
			content:   "This is a test comment",
			authorID:  1,
			articleID: 1,
			wantError: false,
		},
		{
			name:      "Empty content should fail",
			content:   "",
			authorID:  1,
			articleID: 1,
			wantError: true,
		},
		{
			name:      "Invalid author should fail",
			content:   "This is a test comment",
			authorID:  999,
			articleID: 1,
			wantError: true,
		},
		{
			name:      "Invalid article should fail",
			content:   "This is a test comment",
			authorID:  1,
			articleID: 999,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comment, err := usecase.CreateComment(context.Background(), tt.content, tt.authorID, tt.articleID)

			if tt.wantError {
				if err == nil {
					t.Errorf("CreateComment() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("CreateComment() unexpected error: %v", err)
				}
				if comment == nil {
					t.Errorf("CreateComment() expected comment, got nil")
				}
				if comment.Content != tt.content {
					t.Errorf("CreateComment() content = %v, want %v", comment.Content, tt.content)
				}
				if comment.Status != "pending" {
					t.Errorf("CreateComment() status = %v, want 'pending'", comment.Status)
				}
			}
		})
	}
}

func TestCommentUsecase_ApproveComment(t *testing.T) {
	commentRepo := newMockCommentRepository()
	userRepo := newMockUserRepository()
	articleRepo := newMockArticleRepository()
	txManager := &mockTransactionManager{}

	usecase := NewCommentUsecase(commentRepo, userRepo, articleRepo, txManager)

	// Create a pending comment first
	comment := &entity.Comment{
		ID:        1,
		Content:   "Test comment",
		AuthorID:  1,
		ArticleID: 1,
		Status:    "pending",
	}
	commentRepo.comments[1] = comment

	approvedComment, err := usecase.ApproveComment(context.Background(), 1)
	if err != nil {
		t.Errorf("ApproveComment() unexpected error: %v", err)
	}

	if approvedComment.Status != "approved" {
		t.Errorf("ApproveComment() status = %v, want 'approved'", approvedComment.Status)
	}

	// Test approving non-existent comment
	_, err = usecase.ApproveComment(context.Background(), 999)
	if err == nil {
		t.Errorf("ApproveComment() expected error for non-existent comment")
	}
}