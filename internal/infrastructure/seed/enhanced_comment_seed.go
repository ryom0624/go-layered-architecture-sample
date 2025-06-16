package seed

import (
	"layered-architecture-template/internal/domain/entity"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

func SeedEnhancedComments(db *gorm.DB) error {
	// Get all articles to distribute comments
	var articles []entity.Article
	if err := db.Find(&articles).Error; err != nil {
		return err
	}

	// Get all users for comment authors
	var users []entity.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}

	if len(articles) == 0 || len(users) == 0 {
		return nil
	}

	// Get the current max comment ID to avoid conflicts
	var maxID uint
	db.Model(&entity.Comment{}).Select("COALESCE(MAX(id), 0)").Scan(&maxID)
	
	comments := []entity.Comment{}
	commentID := maxID + 1

	// Technical discussions and responses
	techComments := []string{
		"Excellent implementation! I particularly appreciate the clean separation of concerns.",
		"This approach reminds me of the patterns we use in production at Google. Very scalable.",
		"Have you considered the performance implications at scale? Might want to add caching.",
		"Great tutorial! Could you add examples for error handling in concurrent scenarios?",
		"This is exactly what I was looking for. The code examples are very clear.",
		"I've been struggling with this concept for weeks. Your explanation finally made it click!",
		"Interesting approach. How does this compare to the microservices pattern?",
		"Love the practical examples. Do you have a GitHub repo with the complete code?",
		"This would be perfect for our current project. Thanks for sharing!",
		"I implemented this in our system and saw a 40% performance improvement.",
		"Question: How would you handle database migrations in this architecture?",
		"The dependency injection pattern shown here is really elegant.",
		"This solves a major pain point we've been having. Brilliant solution!",
		"Could you expand on the testing strategy for this pattern?",
		"We use a similar approach but with Redis for caching. Great minds think alike!",
		"The security considerations section is particularly valuable. Thank you!",
		"I'd love to see a follow-up article on monitoring and observability.",
		"This pattern fits perfectly with our event-driven architecture.",
		"The performance benchmarks would be interesting to see.",
		"How does this handle failure scenarios and circuit breaking?",
		"Brilliant article! The step-by-step breakdown makes it easy to follow.",
		"I'm curious about the resource usage implications of this approach.",
		"This is production-ready code. Love the attention to detail.",
		"Have you tested this with Kubernetes deployments?",
		"The configuration management aspect is handled very well here.",
		"This would integrate nicely with our CI/CD pipeline.",
		"I appreciate the focus on maintainability and code readability.",
		"The async patterns here are exactly what modern applications need.",
		"Question: How would you handle versioning with this API design?",
		"This reminds me of the patterns Martin Fowler discusses. Excellent work!",
	}

	replies := []string{
		"Thanks! I spent a lot of time refining this pattern in production.",
		"Good point about caching. I'll add that to the follow-up article.",
		"The GitHub repo is linked at the bottom of the article.",
		"You raise an interesting question about migrations. Let me think about that.",
		"Redis is a great choice for caching in this scenario!",
		"Circuit breaking is definitely important. I'll cover that in part 2.",
		"Kubernetes deployment works great with this pattern. No issues so far.",
		"Performance benchmarks are a great idea. I'll add those to the repo.",
		"Martin Fowler's work definitely influenced this approach.",
		"Error handling in concurrent scenarios deserves its own dedicated article.",
		"The testing strategy gets complex, but it's worth the investment.",
		"Monitoring is crucial. I'm planning a whole series on observability.",
		"Version handling through headers works well with this API design.",
		"Resource usage is minimal due to the efficient data structures used.",
		"I'm glad it helped you understand the concept better!",
	}

	// Seed comments for each article
	for i, article := range articles {
		if i >= 60 { // Limit to first 60 articles to keep it manageable
			break
		}

		// 3-8 comments per article
		numComments := rand.Intn(6) + 3
		
		for j := 0; j < numComments; j++ {
			// Select random user (avoid article author commenting on their own article frequently)
			var authorID uint
			for {
				authorID = users[rand.Intn(len(users))].ID
				if authorID != article.AuthorID || rand.Float32() < 0.1 { // 10% chance author comments on own article
					break
				}
			}

			// Determine comment status
			var status string
			statusRand := rand.Float32()
			if statusRand < 0.85 {
				status = "approved"
			} else if statusRand < 0.95 {
				status = "pending"
			} else {
				status = "rejected"
			}

			comment := entity.Comment{
				ID:        commentID,
				Content:   techComments[rand.Intn(len(techComments))],
				AuthorID:  authorID,
				ArticleID: article.ID,
				Status:    status,
			}

			comments = append(comments, comment)
			commentID++

			// 30% chance of having a reply to this comment
			if rand.Float32() < 0.3 && status == "approved" {
				// Reply by article author or another user
				var replyAuthorID uint
				if rand.Float32() < 0.6 {
					replyAuthorID = article.AuthorID // Article author replies
				} else {
					replyAuthorID = users[rand.Intn(len(users))].ID
				}

				reply := entity.Comment{
					ID:        commentID,
					Content:   replies[rand.Intn(len(replies))],
					AuthorID:  replyAuthorID,
					ArticleID: article.ID,
					ParentID:  &comment.ID,
					Status:    "approved",
				}

				comments = append(comments, reply)
				commentID++

				// 15% chance of a second-level reply
				if rand.Float32() < 0.15 {
					var secondReplyAuthorID uint
					for {
						secondReplyAuthorID = users[rand.Intn(len(users))].ID
						if secondReplyAuthorID != replyAuthorID {
							break
						}
					}

					secondReply := entity.Comment{
						ID:        commentID,
						Content:   replies[rand.Intn(len(replies))],
						AuthorID:  secondReplyAuthorID,
						ArticleID: article.ID,
						ParentID:  &reply.ID,
						Status:    "approved",
					}

					comments = append(comments, secondReply)
					commentID++
				}
			}
		}
	}

	// Add some specific technical discussions for popular topics
	specificDiscussions := []struct {
		ArticleTitle string
		Comments     []string
	}{
		{
			ArticleTitle: "Deep Learning with TensorFlow 2.0",
			Comments: []string{
				"The eager execution feature in TF 2.0 is a game-changer for debugging.",
				"Have you tried this with TPUs? The performance gains are incredible.",
				"I'm curious about memory usage compared to PyTorch. Any benchmarks?",
				"The Keras integration makes this so much more accessible.",
				"This tutorial helped me migrate our legacy TF 1.x models. Thank you!",
			},
		},
		{
			ArticleTitle: "Kubernetes Cluster Management",
			Comments: []string{
				"The RBAC section is particularly valuable for production deployments.",
				"We use this exact setup in our production cluster serving 10M+ requests daily.",
				"The network policy examples are spot on. Security is often overlooked.",
				"How do you handle cluster autoscaling with this configuration?",
				"The monitoring setup with Prometheus is exactly what we needed.",
			},
		},
		{
			ArticleTitle: "Smart Contract Development with Solidity",
			Comments: []string{
				"The security patterns shown here are essential for avoiding exploits.",
				"Gas optimization techniques saved us thousands in deployment costs.",
				"The testing framework setup is comprehensive. Love the edge case coverage.",
				"Have you considered the implications of EIP-1559 on these contracts?",
				"The upgrade pattern using proxies is handled very well here.",
			},
		},
	}

	// Add specific discussions
	for _, discussion := range specificDiscussions {
		var targetArticle entity.Article
		if err := db.Where("title = ?", discussion.ArticleTitle).First(&targetArticle).Error; err != nil {
			continue // Skip if article not found
		}

		for _, commentText := range discussion.Comments {
			comment := entity.Comment{
				ID:        commentID,
				Content:   commentText,
				AuthorID:  users[rand.Intn(len(users))].ID,
				ArticleID: targetArticle.ID,
				Status:    "approved",
			}
			comments = append(comments, comment)
			commentID++
		}
	}

	// Batch insert comments
	batchSize := 50
	for i := 0; i < len(comments); i += batchSize {
		end := i + batchSize
		if end > len(comments) {
			end = len(comments)
		}
		
		if err := db.Create(comments[i:end]).Error; err != nil {
			return err
		}
	}

	return nil
}


// Helper function to seed additional data for testing new features
func SeedAdvancedFeatureData(db *gorm.DB) error {
	// This will be used for seeding favorites, reading lists, view history, etc.
	// Implementation will be added as we expand the seeding
	
	return SeedFavorites(db)
}

func SeedFavorites(db *gorm.DB) error {
	// Get sample users and articles
	var users []entity.User
	var articles []entity.Article
	
	if err := db.Limit(20).Find(&users).Error; err != nil {
		return err
	}
	if err := db.Where("status = ?", "published").Limit(30).Find(&articles).Error; err != nil {
		return err
	}

	if len(users) == 0 || len(articles) == 0 {
		return nil
	}

	// Create random favorite relationships
	favorites := []entity.Favorite{}
	favoriteID := uint(1)

	for _, user := range users {
		// Each user favorites 3-8 random articles
		numFavorites := rand.Intn(6) + 3
		selectedArticles := make(map[uint]bool)
		
		for i := 0; i < numFavorites && i < len(articles); i++ {
			// Select random article that hasn't been favorited by this user yet
			var articleID uint
			for {
				articleID = articles[rand.Intn(len(articles))].ID
				if !selectedArticles[articleID] {
					selectedArticles[articleID] = true
					break
				}
			}

			favorite := entity.Favorite{
				ID:        favoriteID,
				UserID:    user.ID,
				ArticleID: articleID,
				CreatedAt: time.Now().AddDate(0, 0, -rand.Intn(30)), // Random date within last 30 days
			}
			favorites = append(favorites, favorite)
			favoriteID++
		}
	}

	// Batch insert favorites
	batchSize := 50
	for i := 0; i < len(favorites); i += batchSize {
		end := i + batchSize
		if end > len(favorites) {
			end = len(favorites)
		}
		
		if err := db.Create(favorites[i:end]).Error; err != nil {
			return err
		}
	}

	return nil
}