package video

import (
	"fmt"
	"strings"
	"youtube_tracker/src/model"

	"github.com/jmoiron/sqlx"
)

type VideoRepository interface {
	SaveVideo(video *model.Video) error
	GetVideosByKeyword(keyword string) ([]*model.Video, error)
	GetVideoAnalytics(artists []string) ([]*model.Video, error)
}

type videoRepositoryImpl struct {
	db *sqlx.DB
}

func NewVideoRepository(db *sqlx.DB) VideoRepository {
	return &videoRepositoryImpl{db}
}

func (r *videoRepositoryImpl) SaveVideo(video *model.Video) error {
	_, err := r.db.NamedExec(`INSERT INTO videos (video_id, title, channel_name, description, view_count, like_count, comment_count)
                              VALUES (:video_id, :title, :channel_name, :description, :view_count, :like_count, :comment_count)
							  ON CONFLICT (video_id) DO NOTHING`,
		video)
	return err
}

func (r *videoRepositoryImpl) GetVideosByKeyword(keyword string) ([]*model.Video, error) {
	var videos []*model.Video
	err := r.db.Select(&videos, `SELECT * FROM videos WHERE title LIKE $1 OR description LIKE $2`, "%"+keyword+"%", "%"+keyword+"%")
	return videos, err
}

func (r *videoRepositoryImpl) GetVideoAnalytics(artists []string) ([]*model.Video, error) {
	var videos []*model.Video
	// Construct the query to include all keywords
	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT video_id, title, channel_name, description, view_count, like_count, comment_count FROM videos WHERE ")

	conditions := []string{}
	args := []interface{}{}
	for i, artist := range artists {
		conditions = append(conditions, fmt.Sprintf("(title ILIKE '%%' || $%d || '%%' OR description ILIKE '%%' || $%d || '%%')", i+1, i+1))
		args = append(args, artist)
	}

	queryBuilder.WriteString(strings.Join(conditions, " OR "))

	query := queryBuilder.String()
	err := r.db.Select(&videos, query, args...)
	if err != nil {
		return nil, err
	}
	return videos, nil
}
