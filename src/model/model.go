package model

type Video struct {
	VideoID      string `db:"video_id" json:"video_id"`
	Title        string `db:"title" json:"title"`
	ChannelName  string `db:"channel_name" json:"channel_name"`
	Description  string `db:"description" json:"description"`
	ViewCount    int64  `db:"view_count" json:"view_count"`
	LikeCount    int64  `db:"like_count" json:"like_count"`
	CommentCount int64  `db:"comment_count" json:"comment_count"`
}
