package db

import (
	"panel/config"
	"panel/logger"
	"panel/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(sqlite.Open(config.DBPath+"?_journal_mode=WAL&_foreign_keys=ON"), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		logger.Fatal("打开数据库失败", "error", err)
	}

	// 自动迁移表结构
	runMigrations()

	// 初始化默认数据
	initDefaultData()
}

func runMigrations() {
	// 迁移所有模型
	err := DB.AutoMigrate(
		&model.User{},
		&model.Group{},
		&model.Bookmark{},
		&model.Image{},
		&model.UserSettings{},
	)
	if err != nil {
		logger.Error("自动迁移失败", "error", err)
	}
	logger.Info("表结构迁移完成")
}

func initDefaultData() {
	var count int64
	DB.Model(&model.User{}).Count(&count)
	if count == 0 {
		logger.Info("初始化默认用户和书签数据...")

		adminHash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		guestHash, _ := bcrypt.GenerateFromPassword([]byte("guest"), bcrypt.DefaultCost)

		admin := model.User{Username: "admin", Password: string(adminHash), Avatar: "👤", Role: "admin"}
		guest := model.User{Username: "guest", Password: string(guestHash), Avatar: "👥", Role: "guest"}
		DB.Create(&admin)
		DB.Create(&guest)

		// 默认数据归属管理员（user_id=1）
		const adminID = 1

		// 常用工具
		group1 := model.Group{ID: 1, UserID: adminID, Name: "常用工具", Icon: "🛠️", Sort: 1}
		DB.Create(&group1)
		bookmarks1 := []model.Bookmark{
			{UserID: adminID, GroupID: 1, Name: "GitHub", URL: "https://github.com", Icon: "🐙", Sort: 1},
			{UserID: adminID, GroupID: 1, Name: "Google", URL: "https://www.google.com", Icon: "🔍", Sort: 2},
			{UserID: adminID, GroupID: 1, Name: "Gmail", URL: "https://mail.google.com", Icon: "📧", Sort: 3},
			{UserID: adminID, GroupID: 1, Name: "ChatGPT", URL: "https://chatgpt.com", Icon: "🤖", Sort: 4},
			{UserID: adminID, GroupID: 1, Name: "DeepSeek", URL: "https://chat.deepseek.com", Icon: "🧠", Sort: 5},
			{UserID: adminID, GroupID: 1, Name: "百度", URL: "https://www.baidu.com", Icon: "🌐", Sort: 6},
			{UserID: adminID, GroupID: 1, Name: "Bing", URL: "https://www.bing.com", Icon: "🔎", Sort: 7},
			{UserID: adminID, GroupID: 1, Name: "维基百科", URL: "https://en.wikipedia.org", Icon: "📖", Sort: 8},
			{UserID: adminID, GroupID: 1, Name: "Stack Overflow", URL: "https://stackoverflow.com", Icon: "💻", Sort: 9},
			{UserID: adminID, GroupID: 1, Name: "NPM", URL: "https://www.npmjs.com", Icon: "📦", Sort: 10},
		}
		DB.Create(&bookmarks1)

		// 社交媒体
		group2 := model.Group{ID: 2, UserID: adminID, Name: "社交媒体", Icon: "💬", Sort: 2}
		DB.Create(&group2)
		bookmarks2 := []model.Bookmark{
			{UserID: adminID, GroupID: 2, Name: "Bilibili", URL: "https://www.bilibili.com", Icon: "📺", Sort: 1},
			{UserID: adminID, GroupID: 2, Name: "微博", URL: "https://weibo.com", Icon: "📱", Sort: 2},
			{UserID: adminID, GroupID: 2, Name: "知乎", URL: "https://www.zhihu.com", Icon: "💡", Sort: 3},
			{UserID: adminID, GroupID: 2, Name: "小红书", URL: "https://www.xiaohongshu.com", Icon: "📕", Sort: 4},
			{UserID: adminID, GroupID: 2, Name: "抖音", URL: "https://www.douyin.com", Icon: "🎵", Sort: 5},
			{UserID: adminID, GroupID: 2, Name: "微信公众号", URL: "https://mp.weixin.qq.com", Icon: "💬", Sort: 6},
			{UserID: adminID, GroupID: 2, Name: "豆瓣", URL: "https://www.douban.com", Icon: "📚", Sort: 7},
			{UserID: adminID, GroupID: 2, Name: "Twitter / X", URL: "https://x.com", Icon: "🐦", Sort: 8},
			{UserID: adminID, GroupID: 2, Name: "Reddit", URL: "https://www.reddit.com", Icon: "👽", Sort: 9},
		}
		DB.Create(&bookmarks2)

		// 影音娱乐
		group3 := model.Group{ID: 3, UserID: adminID, Name: "影音娱乐", Icon: "🎬", Sort: 3}
		DB.Create(&group3)
		bookmarks3 := []model.Bookmark{
			{UserID: adminID, GroupID: 3, Name: "YouTube", URL: "https://www.youtube.com", Icon: "▶️", Sort: 1},
			{UserID: adminID, GroupID: 3, Name: "Netflix", URL: "https://www.netflix.com", Icon: "🎬", Sort: 2},
			{UserID: adminID, GroupID: 3, Name: "Spotify", URL: "https://open.spotify.com", Icon: "🎵", Sort: 3},
			{UserID: adminID, GroupID: 3, Name: "网易云音乐", URL: "https://music.163.com", Icon: "🎶", Sort: 4},
			{UserID: adminID, GroupID: 3, Name: "腾讯视频", URL: "https://v.qq.com", Icon: "🎥", Sort: 5},
			{UserID: adminID, GroupID: 3, Name: "爱奇艺", URL: "https://www.iqiyi.com", Icon: "📺", Sort: 6},
		}
		DB.Create(&bookmarks3)

		// 开发技术
		group4 := model.Group{ID: 4, UserID: adminID, Name: "开发技术", Icon: "⚙️", Sort: 4}
		DB.Create(&group4)
		bookmarks4 := []model.Bookmark{
			{UserID: adminID, GroupID: 4, Name: "MDN Web Docs", URL: "https://developer.mozilla.org", Icon: "📝", Sort: 1},
			{UserID: adminID, GroupID: 4, Name: "Vue 文档", URL: "https://cn.vuejs.org", Icon: "💚", Sort: 2},
			{UserID: adminID, GroupID: 4, Name: "Node.js 文档", URL: "https://nodejs.org/docs", Icon: "🟢", Sort: 3},
			{UserID: adminID, GroupID: 4, Name: "Docker Hub", URL: "https://hub.docker.com", Icon: "🐳", Sort: 4},
			{UserID: adminID, GroupID: 4, Name: "掘金", URL: "https://juejin.cn", Icon: "⛏️", Sort: 5},
			{UserID: adminID, GroupID: 4, Name: "CSDN", URL: "https://www.csdn.net", Icon: "📘", Sort: 6},
			{UserID: adminID, GroupID: 4, Name: "V2EX", URL: "https://www.v2ex.com", Icon: "🧑‍💻", Sort: 7},
			{UserID: adminID, GroupID: 4, Name: "Hacker News", URL: "https://news.ycombinator.com", Icon: "📰", Sort: 8},
			{UserID: adminID, GroupID: 4, Name: "Product Hunt", URL: "https://www.producthunt.com", Icon: "🚀", Sort: 9},
		}
		DB.Create(&bookmarks4)

		// 设计灵感
		group5 := model.Group{ID: 5, UserID: adminID, Name: "设计灵感", Icon: "🎨", Sort: 5}
		DB.Create(&group5)
		bookmarks5 := []model.Bookmark{
			{UserID: adminID, GroupID: 5, Name: "Dribbble", URL: "https://dribbble.com", Icon: "🏀", Sort: 1},
			{UserID: adminID, GroupID: 5, Name: "Behance", URL: "https://www.behance.net", Icon: "✨", Sort: 2},
			{UserID: adminID, GroupID: 5, Name: "Pinterest", URL: "https://www.pinterest.com", Icon: "📌", Sort: 3},
			{UserID: adminID, GroupID: 5, Name: "UI中国", URL: "https://www.ui.cn", Icon: "🎯", Sort: 4},
			{UserID: adminID, GroupID: 5, Name: "站酷", URL: "https://www.zcool.com.cn", Icon: "🎨", Sort: 5},
			{UserID: adminID, GroupID: 5, Name: "Unsplash", URL: "https://unsplash.com", Icon: "📷", Sort: 6},
			{UserID: adminID, GroupID: 5, Name: "IconFont", URL: "https://www.iconfont.cn", Icon: "🔤", Sort: 7},
		}
		DB.Create(&bookmarks5)

		// 效率工具
		group6 := model.Group{ID: 6, UserID: adminID, Name: "效率工具", Icon: "⏱️", Sort: 6}
		DB.Create(&group6)
		bookmarks6 := []model.Bookmark{
			{UserID: adminID, GroupID: 6, Name: "Notion", URL: "https://www.notion.so", Icon: "📋", Sort: 1},
			{UserID: adminID, GroupID: 6, Name: "飞书", URL: "https://www.feishu.cn", Icon: "📄", Sort: 2},
			{UserID: adminID, GroupID: 6, Name: "语雀", URL: "https://www.yuque.com", Icon: "📗", Sort: 3},
			{UserID: adminID, GroupID: 6, Name: "腾讯文档", URL: "https://docs.qq.com", Icon: "📑", Sort: 4},
			{UserID: adminID, GroupID: 6, Name: "石墨文档", URL: "https://shimo.im", Icon: "✏️", Sort: 5},
			{UserID: adminID, GroupID: 6, Name: "ProcessOn", URL: "https://www.processon.com", Icon: "📊", Sort: 6},
			{UserID: adminID, GroupID: 6, Name: "Canva", URL: "https://www.canva.com", Icon: "🖼️", Sort: 7},
			{UserID: adminID, GroupID: 6, Name: "DeepL", URL: "https://www.deepl.com", Icon: "🌍", Sort: 8},
			{UserID: adminID, GroupID: 6, Name: "百度翻译", URL: "https://fanyi.baidu.com", Icon: "🔤", Sort: 9},
			{UserID: adminID, GroupID: 6, Name: "腾讯云", URL: "https://cloud.tencent.com", Icon: "☁️", Sort: 10},
		}
		DB.Create(&bookmarks6)

		logger.Info("默认数据初始化完成")
	}
}
