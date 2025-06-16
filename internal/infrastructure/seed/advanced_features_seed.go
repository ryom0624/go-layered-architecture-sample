package seed

import (
	"layered-architecture-template/internal/domain/entity"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

func SeedReadingLists(db *gorm.DB) error {
	var users []entity.User
	var articles []entity.Article
	
	if err := db.Limit(25).Find(&users).Error; err != nil {
		return err
	}
	if err := db.Where("status = ?", "published").Find(&articles).Error; err != nil {
		return err
	}

	if len(users) == 0 || len(articles) == 0 {
		return nil
	}

	readingLists := []entity.ReadingList{}
	readingListItems := []entity.ReadingListItem{}
	listID := uint(1)
	itemID := uint(1)

	listTemplates := []struct {
		Name        string
		Description string
		IsPublic    bool
	}{
		{"AI/ML Learning Path", "Essential articles for learning machine learning and artificial intelligence", true},
		{"DevOps Best Practices", "Comprehensive guide to modern DevOps practices and tools", true},
		{"Web Development Fundamentals", "Core concepts every web developer should know", true},
		{"Security Essentials", "Critical security practices for modern applications", false},
		{"Frontend Technologies", "Latest trends and best practices in frontend development", true},
		{"Backend Architecture", "Scalable backend design patterns and practices", false},
		{"Mobile Development", "Cross-platform mobile development techniques", true},
		{"Data Engineering", "Big data processing and pipeline management", false},
		{"Cloud Computing", "Cloud-native development and deployment strategies", true},
		{"Blockchain & Web3", "Decentralized technologies and blockchain development", true},
		{"My Reading Queue", "Articles I want to read later", false},
		{"Team Recommendations", "Articles recommended by my team", false},
		{"Conference Prep", "Articles to review before the upcoming conference", false},
		{"Weekend Reading", "Light technical reading for weekends", false},
		{"Advanced Topics", "Complex topics requiring deep understanding", false},
	}

	for i, user := range users {
		// Each user creates 1-4 reading lists
		numLists := rand.Intn(4) + 1
		
		for j := 0; j < numLists && j < len(listTemplates); j++ {
			template := listTemplates[(i*4+j)%len(listTemplates)]
			
			readingList := entity.ReadingList{
				ID:          listID,
				Name:        template.Name,
				Description: template.Description,
				UserID:      user.ID,
				IsPublic:    template.IsPublic,
				CreatedAt:   time.Now().AddDate(0, 0, -rand.Intn(60)),
				UpdatedAt:   time.Now().AddDate(0, 0, -rand.Intn(30)),
			}
			readingLists = append(readingLists, readingList)
			
			// Add 3-10 articles to each reading list
			numArticles := rand.Intn(8) + 3
			selectedArticles := make(map[uint]bool)
			
			for k := 0; k < numArticles && k < len(articles); k++ {
				var articleID uint
				for {
					articleID = articles[rand.Intn(len(articles))].ID
					if !selectedArticles[articleID] {
						selectedArticles[articleID] = true
						break
					}
				}

				var notes string
				noteOptions := []string{
					"",
					"Need to review this for the project",
					"Great example for the team",
					"Important for certification",
					"Reference for implementation",
					"Follow up on this approach",
					"Discuss with senior dev",
					"Try this in sandbox environment",
				}
				notes = noteOptions[rand.Intn(len(noteOptions))]

				item := entity.ReadingListItem{
					ID:            itemID,
					ReadingListID: listID,
					ArticleID:     articleID,
					AddedAt:       time.Now().AddDate(0, 0, -rand.Intn(45)),
					Notes:         notes,
				}
				readingListItems = append(readingListItems, item)
				itemID++
			}
			
			listID++
		}
	}

	// Insert reading lists
	if err := db.Create(&readingLists).Error; err != nil {
		return err
	}

	// Insert reading list items
	batchSize := 100
	for i := 0; i < len(readingListItems); i += batchSize {
		end := i + batchSize
		if end > len(readingListItems) {
			end = len(readingListItems)
		}
		
		if err := db.Create(readingListItems[i:end]).Error; err != nil {
			return err
		}
	}

	return nil
}

func SeedViewHistory(db *gorm.DB) error {
	var users []entity.User
	var articles []entity.Article
	
	if err := db.Limit(30).Find(&users).Error; err != nil {
		return err
	}
	if err := db.Where("status = ?", "published").Find(&articles).Error; err != nil {
		return err
	}

	if len(users) == 0 || len(articles) == 0 {
		return nil
	}

	articleViews := []entity.ArticleView{}
	readingHistories := []entity.UserReadingHistory{}
	viewID := uint(1)
	historyID := uint(1)

	// Generate realistic IP addresses
	ipAddresses := []string{
		"192.168.1.100", "10.0.0.50", "172.16.0.25", "203.0.113.45",
		"198.51.100.30", "192.0.2.75", "169.254.1.100", "10.10.10.10",
	}

	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_7_1 like Mac OS X)",
		"Mozilla/5.0 (Android 11; Mobile; rv:68.0) Gecko/68.0",
	}

	for _, user := range users {
		// Each user views 10-40 articles
		numViews := rand.Intn(31) + 10
		viewedArticles := make(map[uint]bool)
		
		for i := 0; i < numViews && i < len(articles); i++ {
			var articleID uint
			for {
				articleID = articles[rand.Intn(len(articles))].ID
				if !viewedArticles[articleID] {
					viewedArticles[articleID] = true
					break
				}
			}

			viewTime := time.Now().AddDate(0, 0, -rand.Intn(90)) // Random time within last 90 days
			
			// Create article view
			view := entity.ArticleView{
				ID:        viewID,
				ArticleID: articleID,
				UserID:    &user.ID,
				IPAddress: ipAddresses[rand.Intn(len(ipAddresses))],
				UserAgent: userAgents[rand.Intn(len(userAgents))],
				ViewedAt:  viewTime,
			}
			articleViews = append(articleViews, view)
			
			// Create reading history with progress
			progress := rand.Float32() * 100 // 0-100%
			readingTime := rand.Intn(1800) + 30 // 30 seconds to 30 minutes
			isCompleted := progress >= 90.0
			
			history := entity.UserReadingHistory{
				ID:               historyID,
				UserID:           user.ID,
				ArticleID:        articleID,
				ReadingProgress:  progress,
				TotalReadingTime: readingTime,
				FirstViewedAt:    viewTime.AddDate(0, 0, -rand.Intn(30)),
				LastViewedAt:     viewTime,
				IsCompleted:      isCompleted,
			}
			readingHistories = append(readingHistories, history)
			
			viewID++
			historyID++
		}
	}

	// Insert article views
	batchSize := 100
	for i := 0; i < len(articleViews); i += batchSize {
		end := i + batchSize
		if end > len(articleViews) {
			end = len(articleViews)
		}
		
		if err := db.Create(articleViews[i:end]).Error; err != nil {
			return err
		}
	}

	// Insert reading histories
	for i := 0; i < len(readingHistories); i += batchSize {
		end := i + batchSize
		if end > len(readingHistories) {
			end = len(readingHistories)
		}
		
		if err := db.Create(readingHistories[i:end]).Error; err != nil {
			return err
		}
	}

	return nil
}

func SeedStatistics(db *gorm.DB) error {
	var articles []entity.Article
	if err := db.Find(&articles).Error; err != nil {
		return err
	}

	if len(articles) == 0 {
		return nil
	}

	articleStats := []entity.ArticleStatistics{}
	dailyStats := []entity.DailyStatistics{}
	statsID := uint(1)
	dailyID := uint(1)

	// Create article statistics
	for _, article := range articles {
		views := rand.Intn(10000) + 100
		
		stats := entity.ArticleStatistics{
			ID:               statsID,
			ArticleID:        article.ID,
			TotalViews:       views,
			UniqueViews:      views - rand.Intn(views/3),
			AuthenticatedViews: rand.Intn(views/2),
			AnonymousViews:   views - rand.Intn(views/2),
			AverageReadingTime: float32(rand.Intn(600) + 60), // 1-10 minutes
			CompletionRate:   float32(rand.Intn(80) + 20), // 20-100%
			TotalReadingTime: views * (rand.Intn(300) + 60), // Total reading time
			LastCalculatedAt: time.Now().AddDate(0, 0, -rand.Intn(7)),
		}
		articleStats = append(articleStats, stats)
		statsID++
	}

	// Create daily statistics for the last 30 days
	for i := 0; i < 30; i++ {
		date := time.Now().AddDate(0, 0, -i)
		
		daily := entity.DailyStatistics{
			ID:                dailyID,
			Date:              date,
			TotalViews:        rand.Intn(50000) + 10000,
			UniqueUsers:       rand.Intn(5000) + 1000,
			NewUsers:          rand.Intn(100) + 10,
			TotalReadingTime:  rand.Intn(500000) + 100000,
			AverageReadingTime: float32(rand.Intn(300) + 60),
			TotalArticles:     len(articles),
			NewArticles:       rand.Intn(10),
			PublishedArticles: rand.Intn(len(articles)) + len(articles)/2,
		}
		dailyStats = append(dailyStats, daily)
		dailyID++
	}

	// Insert statistics
	if err := db.Create(&articleStats).Error; err != nil {
		return err
	}

	if err := db.Create(&dailyStats).Error; err != nil {
		return err
	}

	return nil
}

func SeedAuthTokens(db *gorm.DB) error {
	var users []entity.User
	if err := db.Limit(20).Find(&users).Error; err != nil {
		return err
	}

	if len(users) == 0 {
		return nil
	}

	tokens := []entity.RefreshToken{}
	tokenID := uint(1)

	// Create some active refresh tokens
	for i, user := range users {
		if i >= 15 { // Only 15 users have active tokens
			break
		}
		
		// 1-3 devices per user
		numTokens := rand.Intn(3) + 1
		
		for j := 0; j < numTokens; j++ {
			expiresAt := time.Now().AddDate(0, 0, 7) // 7 days from now
			
			token := entity.RefreshToken{
				ID:        tokenID,
				UserID:    user.ID,
				Token:     generateRandomToken(),
				ExpiresAt: expiresAt,
				CreatedAt: time.Now().AddDate(0, 0, -rand.Intn(7)),
			}
			tokens = append(tokens, token)
			tokenID++
		}
	}

	if err := db.Create(&tokens).Error; err != nil {
		return err
	}

	return nil
}

func generateRandomToken() string {
	// Generate a realistic-looking refresh token
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	token := make([]byte, 64)
	for i := range token {
		token[i] = chars[rand.Intn(len(chars))]
	}
	return string(token)
}