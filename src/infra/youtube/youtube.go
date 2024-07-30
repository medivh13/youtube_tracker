package youtube

import (
	"context"
	"youtube_tracker/src/model"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Client interface {
	SearchVideos(keyword string) ([]*model.Video, error)
}

type clientImpl struct {
	service *youtube.Service
}

func NewClient(apiKey string) (Client, error) {
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &clientImpl{service}, nil
}

func (c *clientImpl) SearchVideos(keyword string) ([]*model.Video, error) {
	call := c.service.Search.List([]string{"id", "snippet"}).Q(keyword).MaxResults(100)
	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	var videos []*model.Video
	for _, item := range response.Items {
		if item.Id.Kind == "youtube#video" {
			videoID := item.Id.VideoId
			statsCall := c.service.Videos.List([]string{"statistics"}).Id(videoID)
			statsResponse, err := statsCall.Do()
			if err != nil {
				return nil, err
			}

			if len(statsResponse.Items) > 0 {
				stats := statsResponse.Items[0].Statistics
				videos = append(videos, &model.Video{
					VideoID:      videoID,
					Title:        item.Snippet.Title,
					ChannelName:  item.Snippet.ChannelTitle,
					Description:  item.Snippet.Description,
					ViewCount:    int64(stats.ViewCount),
					LikeCount:    int64(stats.LikeCount),
					CommentCount: int64(stats.CommentCount),
				})
			}
		}
	}

	return videos, nil
}
