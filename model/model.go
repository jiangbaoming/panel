package model

// 用户
type User struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username" gorm:"uniqueIndex;not null"`
	Avatar   string `json:"avatar" gorm:"default:👤"`
	Role     string `json:"role" gorm:"default:guest"`
	Password string `json:"-" gorm:"not null"`
}

// 登录响应
type LoginResponse struct {
	Token    string `json:"token"`
	ID       int    `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
}

// 分组
type Group struct {
	ID        int        `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int        `json:"user_id" gorm:"not null;index"`
	Name      string     `json:"name" gorm:"not null"`
	Icon      string     `json:"icon" gorm:"default:📁"`
	Sort      int        `json:"sort" gorm:"default:0;index"`
	Bookmarks []Bookmark `json:"bookmarks" gorm:"foreignKey:GroupID;references:ID"`
}

// 书签
type Bookmark struct {
	ID      int    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID  int    `json:"user_id" gorm:"not null;index"`
	GroupID int    `json:"group_id" gorm:"not null;index"`
	Name    string `json:"name" gorm:"not null"`
	URL     string `json:"url" gorm:"not null"`
	Icon    string `json:"icon" gorm:"default:🔗"`
	Sort    int    `json:"sort" gorm:"default:0;index"`
	BgColor string `json:"bg_color" gorm:"default:''"`
	IconBg  string `json:"icon_bg" gorm:"default:''"`
	Pinned  int    `json:"pinned" gorm:"default:0;index"`
}

// 常驻书签
type PinnedBookmark struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	Icon      string `json:"icon"`
	BgColor   string `json:"bg_color"`
	IconBg    string `json:"icon_bg"`
	GroupName string `json:"group_name"`
}

// 图片
type ImageItem struct {
	ID           int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Filename     string `json:"filename" gorm:"not null"`
	OriginalName string `json:"original_name" gorm:"not null"`
	Category     string `json:"category" gorm:"default:icon;index"`
	URL          string `json:"url" gorm:"-:all"`
	CreatedAt    string `json:"created_at" gorm:"autoCreateTime"`
}

// 图片表结构（用于 GORM 操作）
type Image struct {
	ID           int    `gorm:"primaryKey;autoIncrement"`
	Filename     string `gorm:"not null"`
	OriginalName string `gorm:"not null"`
	Category     string `gorm:"default:icon;index"`
	UploadedBy   int    `gorm:"default:1"`
	CreatedAt    string `gorm:"autoCreateTime"`
}

func (ImageItem) TableName() string {
	return "images"
}

func (Image) TableName() string {
	return "images"
}

// 图片列表响应
type ImageListResponse struct {
	Data     []ImageItem `json:"data"`
	Total    int         `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// 用户设置
type UserSettings struct {
	UserID      int    `json:"user_id" gorm:"primaryKey"`
	BgImage     string `json:"bg_image" gorm:"default:''"`
	DisplayMode string `json:"display_mode" gorm:"default:both"`
	PageTitle   string `json:"page_title" gorm:"default:个人导航页"`
	PageFavicon string `json:"page_favicon" gorm:"default:''"`
}

func (UserSettings) TableName() string {
	return "user_settings"
}

// 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// 重排序请求
type ReorderRequest struct {
	IDs []string `json:"ids"`
}

// 搜索书签结果
type SearchResult struct {
	ID        string `json:"id"`
	GroupID   int    `json:"group_id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	Icon      string `json:"icon"`
	GroupName string `json:"groupName"`
}
