package main

import (
	"log"

	"aiforum/handlers"
	"aiforum/middleware"
	"aiforum/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	err := models.InitDB()
	if err != nil {
		log.Fatal("数据库初始化失败:", err)
	}

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建路由
	r := gin.Default()

	// 静态文件服务
	r.Static("/static", "./static")
	r.Static("/images", "./images")
	r.LoadHTMLGlob("templates/*")

	// 设置路由
	setupRoutes(r)

	// 启动服务器
	log.Println("AI论坛服务器启动在端口 8080...")
	log.Fatal(r.Run(":8080"))
}

func setupRoutes(r *gin.Engine) {
	// 首页
	r.GET("/", handlers.HomePage)

	// 用户认证相关路由
	auth := r.Group("/auth")
	{
		auth.GET("/register", handlers.RegisterPage)
		auth.POST("/register", handlers.Register)
		auth.GET("/login", handlers.LoginPage)
		auth.POST("/login", handlers.Login)
		auth.GET("/logout", handlers.Logout)
	}

	// 问答相关路由
	qa := r.Group("/qa")
	{
		qa.GET("", handlers.QAPage)
		qa.GET("/search", handlers.AdvancedSearch)
		qa.GET("/ask", handlers.AskQuestionPage)
		qa.POST("/ask", handlers.AskQuestion)
		qa.GET("/:id", handlers.ViewQuestion)
		qa.POST("/:id/answer", handlers.AnswerQuestion)
		qa.POST("/answer/:answer_id/accept", handlers.AcceptAnswer)
		qa.POST("/answer/:answer_id/like", handlers.LikeAnswer)
	}

	// 技术分享相关路由
	techShare := r.Group("/tech-share")
	{
		techShare.GET("", handlers.TechSharePage)
		techShare.GET("/:id", handlers.TechShareDetailPage)
		techShare.GET("/topic/:slug", handlers.TopicPage)
		techShare.GET("/publish", handlers.PublishTechSharePage)
		techShare.POST("/publish", handlers.PublishTechShare)
	}

	// API路由
	api := r.Group("/api")
	{
		api.GET("/posts", handlers.GetPosts)
		api.GET("/posts/:id", handlers.GetPost)
		api.GET("/categories", handlers.GetCategories)
		api.GET("/tags", handlers.GetTags)
		
		// 问答API
		api.POST("/answers", handlers.AnswerQuestion)
		api.POST("/answers/:answer_id/accept", handlers.AcceptAnswer)
		api.POST("/answers/:answer_id/like", handlers.LikeAnswer)
		api.POST("/questions/:id/favorite", handlers.FavoriteQuestion)
		api.POST("/questions/:id/report", handlers.ReportQuestion)
		
		// 技术分享API
		api.POST("/tech-share/publish", handlers.PublishTechShare)
		api.POST("/tech-share/:id/like", handlers.LikeTechArticle)
		api.POST("/authors/:author_id/follow", handlers.FollowAuthor)
	}
	
	// 个人中心API路由
	userAPI := r.Group("/api/user")
	userAPI.Use(middleware.AuthMiddleware())
	{
		userAPI.GET("/activity", handlers.GetUserActivity)
		userAPI.GET("/questions", handlers.GetUserQuestions)
		userAPI.GET("/answers", handlers.GetUserAnswers)
		userAPI.GET("/shares", handlers.GetUserShares)
		userAPI.GET("/resources", handlers.GetUserResources)
		userAPI.GET("/favorites", handlers.GetUserFavorites)
		userAPI.GET("/following", handlers.GetUserFollowing)
		userAPI.GET("/followers", handlers.GetUserFollowers)
		userAPI.GET("/messages", handlers.GetUserMessages)
		userAPI.POST("/avatar", handlers.UpdateUserAvatar)
		userAPI.PUT("/profile", handlers.UpdateUserProfile)
		userAPI.PUT("/password", handlers.ChangeUserPassword)
		userAPI.PUT("/notifications", handlers.SaveNotificationSettings)
		userAPI.PUT("/messages/read-all", handlers.MarkAllMessagesRead)
		userAPI.PUT("/messages/:id/read", handlers.MarkMessageRead)
		userAPI.DELETE("/messages/:id", handlers.DeleteMessage)
		userAPI.POST("/:id/follow", handlers.FollowUser)
		userAPI.DELETE("/:id/unfollow", handlers.UnfollowUser)
		userAPI.DELETE("/favorites/:id", handlers.RemoveFavorite)
		userAPI.DELETE("/questions/:id", handlers.DeleteUserQuestion)
		userAPI.DELETE("/answers/:id", handlers.DeleteUserAnswer)
		userAPI.DELETE("/shares/:id", handlers.DeleteUserShare)
		userAPI.DELETE("/resources/:id", handlers.DeleteUserResource)
	}

	// 需要认证的路由
	authenticated := r.Group("/")
	authenticated.Use(middleware.AuthMiddleware())
	{
		// 发帖相关
		authenticated.GET("/post/new", handlers.NewPostPage)
		authenticated.POST("/post/new", handlers.CreatePost)
		authenticated.GET("/post/:id", handlers.ViewPost)
		authenticated.POST("/post/:id/reply", handlers.CreateReply)

		// 用户相关
		authenticated.GET("/profile", handlers.ProfilePage)
		authenticated.POST("/profile/update", handlers.UpdateProfile)

		// 搜索
		authenticated.GET("/search", handlers.Search)
	}


} 