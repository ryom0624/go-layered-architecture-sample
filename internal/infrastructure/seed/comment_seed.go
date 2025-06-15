package seed

import (
	"layered-architecture-template/internal/domain/entity"

	"gorm.io/gorm"
)

func SeedComments(db *gorm.DB) error {
	comments := []entity.Comment{
		{
			ID:        1,
			Content:   "素晴らしい記事でした！Go言語の基礎がよく理解できました。",
			AuthorID:  2,
			ArticleID: 1,
			Status:    "approved",
		},
		{
			ID:        2,
			Content:   "実装例がとても参考になります。ありがとうございます。",
			AuthorID:  3,
			ArticleID: 1,
			Status:    "approved",
		},
		{
			ID:        3,
			Content:   "特にエラーハンドリングの部分が勉強になりました。",
			AuthorID:  4,
			ArticleID: 1,
			ParentID:  uint32ToPointer(1),
			Status:    "approved",
		},
		{
			ID:        4,
			Content:   "Clean Architectureの理解が深まりました。次回の記事も楽しみです！",
			AuthorID:  5,
			ArticleID: 2,
			Status:    "approved",
		},
		{
			ID:        5,
			Content:   "依存性注入の実装がとても綺麗ですね。",
			AuthorID:  1,
			ArticleID: 2,
			Status:    "approved",
		},
		{
			ID:        6,
			Content:   "テストの書き方が参考になります。モックの使い方がよくわかりました。",
			AuthorID:  6,
			ArticleID: 3,
			Status:    "approved",
		},
		{
			ID:        7,
			Content:   "質問があります。パフォーマンスの観点から改善点はありますか？",
			AuthorID:  7,
			ArticleID: 3,
			Status:    "pending",
		},
		{
			ID:        8,
			Content:   "とても詳しい解説ありがとうございます。実際のプロジェクトで試してみます。",
			AuthorID:  8,
			ArticleID: 4,
			Status:    "approved",
		},
		{
			ID:        9,
			Content:   "レスポンシブデザインの実装例が素晴らしいです。",
			AuthorID:  9,
			ArticleID: 5,
			Status:    "approved",
		},
		{
			ID:        10,
			Content:   "CSS Gridの使い方が勉強になりました。",
			AuthorID:  10,
			ArticleID: 5,
			ParentID:  uint32ToPointer(9),
			Status:    "approved",
		},
		{
			ID:        11,
			Content:   "セキュリティの観点からとても参考になる記事でした。",
			AuthorID:  2,
			ArticleID: 6,
			Status:    "approved",
		},
		{
			ID:        12,
			Content:   "実装のベストプラクティスがよくまとまっていますね。",
			AuthorID:  3,
			ArticleID: 7,
			Status:    "approved",
		},
		{
			ID:        13,
			Content:   "この記事のおかげでパフォーマンス改善ができました！",
			AuthorID:  4,
			ArticleID: 8,
			Status:    "approved",
		},
		{
			ID:        14,
			Content:   "機械学習の基礎から応用まで幅広くカバーされていて素晴らしいです。",
			AuthorID:  5,
			ArticleID: 9,
			Status:    "approved",
		},
		{
			ID:        15,
			Content:   "Docker環境での開発が楽になりました。ありがとうございます。",
			AuthorID:  6,
			ArticleID: 10,
			Status:    "approved",
		},
	}

	for _, comment := range comments {
		var existingComment entity.Comment
		if err := db.Where("id = ?", comment.ID).First(&existingComment).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&comment).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}

func uint32ToPointer(val uint32) *uint {
	uintVal := uint(val)
	return &uintVal
}