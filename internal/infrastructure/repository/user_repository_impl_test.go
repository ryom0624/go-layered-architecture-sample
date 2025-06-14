package repository

import (
	"context"
	"database/sql"
	_ "database/sql/driver"
	"regexp"
	"testing"
	"time"

	"layered-architecture-template/internal/domain/entity"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open gorm connection: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return gormDB, mock, cleanup
}

func TestUserRepositoryImpl_Create(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("successful user creation", func(t *testing.T) {
		user := &entity.User{
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("name","email","created_at","updated_at") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
			WithArgs("John Doe", "john@example.com", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err := repo.Create(ctx, user)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("database error during creation", func(t *testing.T) {
		user := &entity.User{
			Name:  "Jane Doe",
			Email: "jane@example.com",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("name","email","created_at","updated_at") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
			WithArgs("Jane Doe", "jane@example.com", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := repo.Create(ctx, user)
		if err == nil {
			t.Error("expected error, got nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

func TestUserRepositoryImpl_GetByID(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("successful get by ID", func(t *testing.T) {
		expectedUser := &entity.User{
			ID:        1,
			Name:      "John Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
			AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.CreatedAt, expectedUser.UpdatedAt)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(rows)

		user, err := repo.GetByID(ctx, 1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if user == nil {
			t.Error("expected user, got nil")
			return
		}
		if user.ID != expectedUser.ID {
			t.Errorf("expected ID %d, got %d", expectedUser.ID, user.ID)
		}
		if user.Name != expectedUser.Name {
			t.Errorf("expected name %s, got %s", expectedUser.Name, user.Name)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(99, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		user, err := repo.GetByID(ctx, 99)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if user != nil {
			t.Error("expected nil user, got user")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

func TestUserRepositoryImpl_GetByEmail(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("successful get by email", func(t *testing.T) {
		expectedUser := &entity.User{
			ID:        1,
			Name:      "John Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
			AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.CreatedAt, expectedUser.UpdatedAt)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs("john@example.com", 1).
			WillReturnRows(rows)

		user, err := repo.GetByEmail(ctx, "john@example.com")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if user == nil {
			t.Error("expected user, got nil")
			return
		}
		if user.Email != expectedUser.Email {
			t.Errorf("expected email %s, got %s", expectedUser.Email, user.Email)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("email not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs("notfound@example.com", 1).
			WillReturnError(gorm.ErrRecordNotFound)

		user, err := repo.GetByEmail(ctx, "notfound@example.com")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if user != nil {
			t.Error("expected nil user, got user")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

func TestUserRepositoryImpl_GetAll(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("successful get all users", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
			AddRow(1, "John Doe", "john@example.com", now, now).
			AddRow(2, "Jane Doe", "jane@example.com", now, now)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(rows)

		users, err := repo.GetAll(ctx)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(users) != 2 {
			t.Errorf("expected 2 users, got %d", len(users))
		}
		if users[0].Name != "John Doe" {
			t.Errorf("expected first user name 'John Doe', got %s", users[0].Name)
		}
		if users[1].Name != "Jane Doe" {
			t.Errorf("expected second user name 'Jane Doe', got %s", users[1].Name)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("empty result", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(rows)

		users, err := repo.GetAll(ctx)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(users) != 0 {
			t.Errorf("expected 0 users, got %d", len(users))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

func TestUserRepositoryImpl_Update(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("successful user update", func(t *testing.T) {
		user := &entity.User{
			ID:    1,
			Name:  "Updated Name",
			Email: "updated@example.com",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "name"=$1,"email"=$2,"created_at"=$3,"updated_at"=$4 WHERE "id" = $5`)).
			WithArgs("Updated Name", "updated@example.com", sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(ctx, user)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("database error during update", func(t *testing.T) {
		user := &entity.User{
			ID:    1,
			Name:  "Updated Name",
			Email: "updated@example.com",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "name"=$1,"email"=$2,"created_at"=$3,"updated_at"=$4 WHERE "id" = $5`)).
			WithArgs("Updated Name", "updated@example.com", sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := repo.Update(ctx, user)
		if err == nil {
			t.Error("expected error, got nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

func TestUserRepositoryImpl_Delete(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("successful user deletion", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE "users"."id" = $1`)).
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
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE "users"."id" = $1`)).
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

	t.Run("delete non-existent user", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE "users"."id" = $1`)).
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
