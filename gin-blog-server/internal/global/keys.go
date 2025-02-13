package global

// Redis key
const (
	ONLINE_USER  = "online_user:"  // 在线用户
	OFFLINE_USER = "offline_user:" // 强制下线用户

	ARTICLE_USER_LIKE_SET = "article_user_like:" // 文章点赞 Set
	ARTICLE_LIKE_COUNT    = "article_like_count" // 文章点赞数
	ARTICLE_VIEW_COUNT    = "article_view_count" // 文章查看数

	COMMENT_USER_LIKE_SET = "comment_user_like:" // 评论点赞 Set
	COMMENT_LIKE_COUNT    = "comment_like_count" // 评论点赞数
)

// Gin Context Key | Session Key
const (
	CTX_DB        = "_db_field"
	CTX_RDB       = "_rdb_field"
	CTX_USER_AUTH = "_user_auth_field"
)
