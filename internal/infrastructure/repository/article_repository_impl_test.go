package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"layered-architecture-template/internal/domain/entity"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
)

func TestArticleRepositoryImpl_Create(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewArticleRepository(db)
	ctx := context.Background()

	t.Run("successful article creation", func(t *testing.T) {
		article := &entity.Article{
			Title:    "Test Article",
			Content:  "Test Content",
			AuthorID: 1,
			Status:   "draft",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "articles" ("title","content","author_id","status","favorite_count","view_count","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)).
			WithArgs("Test Article", "Test Content", 1, "draft", 0, 0, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err := repo.Create(ctx, article)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("database error during creation", func(t *testing.T) {
		article := &entity.Article{
			Title:    "Test Article",
			Content:  "Test Content",
			AuthorID: 1,
			Status:   "draft",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "articles" ("title","content","author_id","status","favorite_count","view_count","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)).
			WithArgs("Test Article", "Test Content", 1, "draft", 0, 0, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := repo.Create(ctx, article)
		if err == nil {
			t.Error("expected error, got nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

func TestArticleRepositoryImpl_GetByID(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewArticleRepository(db)
	ctx := context.Background()

	t.Run("successful get by ID with author", func(t *testing.T) {
		now := time.Now()
		expectedArticle := &entity.Article{
			ID:        1,
			Title:     "Test Article",
			Content:   "Test Content",
			AuthorID:  1,
			Status:    "draft",
			CreatedAt: now,
			UpdatedAt: now,
		}

		// Mock for article query with Preload("Author")
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "articles" WHERE "articles"."id" = $1 ORDER BY "articles"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "author_id", "status", "created_at", "updated_at"}).
				AddRow(expectedArticle.ID, expectedArticle.Title, expectedArticle.Content, expectedArticle.AuthorID, expectedArticle.Status, expectedArticle.CreatedAt, expectedArticle.UpdatedAt))

		// Mock for author query (Preload)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
				AddRow(1, "Test Author", "author@example.com", now, now))

		article, err := repo.GetByID(ctx, 1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if article == nil {
			t.Error("expected article, got nil")
			return
		}
		if article.ID != expectedArticle.ID {
			t.Errorf("expected ID %d, got %d", expectedArticle.ID, article.ID)
		}
		if article.Title != expectedArticle.Title {
			t.Errorf("expected title %s, got %s", expectedArticle.Title, article.Title)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("article not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "articles" WHERE "articles"."id" = $1 ORDER BY "articles"."id" LIMIT $2`)).
			WithArgs(99, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		article, err := repo.GetByID(ctx, 99)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if article != nil {
			t.Error("expected nil article, got article")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

func TestArticleRepositoryImpl_GetAll(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewArticleRepository(db)
	ctx := context.Background()

	t.Run("successful get all articles", func(t *testing.T) {
		now := time.Now()

		// Mock for articles query
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "articles"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "author_id", "status", "created_at", "updated_at"}).
				AddRow(1, "Article 1", "Content 1", 1, "draft", now, now).
				AddRow(2, "Article 2", "Content 2", 1, "published", now, now))

		// Mock for authors query (Preload) - GORM will optimize and query unique IDs only
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
				AddRow(1, "Test Author", "author@example.com", now, now))

		articles, err := repo.GetAll(ctx)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(articles) != 2 {
			t.Errorf("expected 2 articles, got %d", len(articles))
		}
		if articles[0].Title != "Article 1" {
			t.Errorf("expected first article title 'Article 1', got %s", articles[0].Title)
		}
		if articles[1].Title != "Article 2" {
			t.Errorf("expected second article title 'Article 2', got %s", articles[1].Title)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("empty result", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "articles"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "author_id", "status", "created_at", "updated_at"}))
		// No preload query expected for empty result

		articles, err := repo.GetAll(ctx)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(articles) != 0 {
			t.Errorf("expected 0 articles, got %d", len(articles))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

func TestArticleRepositoryImpl_GetByStatus(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewArticleRepository(db)
	ctx := context.Background()

	t.Run("successful get by status", func(t *testing.T) {
		now := time.Now()

		// Mock for articles query with status filter
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "articles" WHERE status = $1`)).
			WithArgs("published").
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "author_id", "status", "created_at", "updated_at"}).
				AddRow(1, "Published Article", "Content", 1, "published", now, now))

		// Mock for author query (Preload)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
				AddRow(1, "Test Author", "author@example.com", now, now))

		articles, err := repo.GetByStatus(ctx, "published")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(articles) != 1 {
			t.Errorf("expected 1 article, got %d", len(articles))
		}
		if articles[0].Status != "published" {
			t.Errorf("expected status 'published', got %s", articles[0].Status)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

func TestArticleRepositoryImpl_Update(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewArticleRepository(db)
	ctx := context.Background()

	t.Run("successful article update", func(t *testing.T) {
		article := &entity.Article{
			ID:       1,
			Title:    "Updated Title",
			Content:  "Updated Content",
			AuthorID: 1,
			Status:   "published",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "articles" SET "title"=$1,"content"=$2,"author_id"=$3,"status"=$4,"favorite_count"=$5,"view_count"=$6,"created_at"=$7,"updated_at"=$8 WHERE "id" = $9`)).
			WithArgs("Updated Title", "Updated Content", 1, "published", 0, 0, sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(ctx, article)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("database error during update", func(t *testing.T) {
		article := &entity.Article{
			ID:       1,
			Title:    "Updated Title",
			Content:  "Updated Content",
			AuthorID: 1,
			Status:   "published",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "articles" SET "title"=$1,"content"=$2,"author_id"=$3,"status"=$4,"favorite_count"=$5,"view_count"=$6,"created_at"=$7,"updated_at"=$8 WHERE "id" = $9`)).
			WithArgs("Updated Title", "Updated Content", 1, "published", 0, 0, sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := repo.Update(ctx, article)
		if err == nil {
			t.Error("expected error, got nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

func TestArticleRepositoryImpl_Delete(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewArticleRepository(db)
	ctx := context.Background()

	t.Run("successful article deletion", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "articles" WHERE "articles"."id" = $1`)).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.Delete(ctx, 1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("database error during deletion", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "articles" WHERE "articles"."id" = $1`)).
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := repo.Delete(ctx, 1)
		if err == nil {
			t.Error("expected error, got nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("delete non-existent article", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "articles" WHERE "articles"."id" = $1`)).
			WithArgs(999).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := repo.Delete(ctx, 999)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

// MockTransactionImpl for testing transaction methods
type MockTransactionImpl struct {
	db *gorm.DB
}

func (m *MockTransactionImpl) GetDB() any {
	return m.db
}

func (m *MockTransactionImpl) Commit() error   { return nil }
func (m *MockTransactionImpl) Rollback() error { return nil }

func TestArticleRepositoryImpl_CreateWithTx(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewArticleRepository(db)
	ctx := context.Background()

	t.Run("successful article creation with transaction", func(t *testing.T) {
		article := &entity.Article{
			Title:    "Test Article",
			Content:  "Test Content",
			AuthorID: 1,
			Status:   "draft",
		}

		// Create a transaction mock that uses the same mocked db
		mock.ExpectBegin()
		txDB := db.Begin()
		txMock := &MockTransactionImpl{db: txDB}

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "articles" ("title","content","author_id","status","favorite_count","view_count","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)).
			WithArgs("Test Article", "Test Content", 1, "draft", 0, 0, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		err := repo.CreateWithTx(ctx, txMock, article)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

func TestArticleRepositoryImpl_UpdateWithTx(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewArticleRepository(db)
	ctx := context.Background()

	t.Run("successful article update with transaction", func(t *testing.T) {
		article := &entity.Article{
			ID:       1,
			Title:    "Updated Title",
			Content:  "Updated Content",
			AuthorID: 1,
			Status:   "published",
		}

		// Create a transaction mock that uses the same mocked db
		mock.ExpectBegin()
		txDB := db.Begin()
		txMock := &MockTransactionImpl{db: txDB}

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "articles" SET "title"=$1,"content"=$2,"author_id"=$3,"status"=$4,"favorite_count"=$5,"view_count"=$6,"created_at"=$7,"updated_at"=$8 WHERE "id" = $9`)).
			WithArgs("Updated Title", "Updated Content", 1, "published", 0, 0, sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateWithTx(ctx, txMock, article)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

func TestArticleRepositoryImpl_DeleteWithTx(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewArticleRepository(db)
	ctx := context.Background()

	t.Run("successful article deletion with transaction", func(t *testing.T) {
		// Create a transaction mock that uses the same mocked db
		mock.ExpectBegin()
		txDB := db.Begin()
		txMock := &MockTransactionImpl{db: txDB}

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "articles" WHERE "articles"."id" = $1`)).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteWithTx(ctx, txMock, 1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}