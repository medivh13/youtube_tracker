package video

import (
	"context"
	"math"
	"sort"
	"strings"
	"youtube_tracker/src/infra/youtube"
	"youtube_tracker/src/model"
	repository "youtube_tracker/src/repository/video"
)

// VideoService defines the interface for video-related operations.
type VideoService interface {
	SearchVideos(ctx context.Context, keywords []string) error
	GetVideos(ctx context.Context, keywords []string) ([]*model.Video, error)
	CalculatePopularityScore(keywords []string) (map[string]float64, error)
}

// videoServiceImpl is an implementation of the VideoService interface.
type videoServiceImpl struct {
	videoRepo     repository.VideoRepository
	youtubeClient youtube.Client
}

// NewVideoService creates a new instance of VideoService.
func NewVideoService(videoRepo repository.VideoRepository, youtubeClient youtube.Client) VideoService {
	return &videoServiceImpl{
		videoRepo, youtubeClient,
	}
}

// SearchVideos searches for videos on YouTube using the provided keywords and saves them in the repository.
func (s *videoServiceImpl) SearchVideos(ctx context.Context, keywords []string) error {
	for _, keyword := range keywords {
		videos, err := s.youtubeClient.SearchVideos(keyword)
		if err != nil {
			return err
		}
		for _, video := range videos {
			if err := s.videoRepo.SaveVideo(video); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetVideos retrieves videos from the repository based on the provided keywords.
func (s *videoServiceImpl) GetVideos(ctx context.Context, keywords []string) ([]*model.Video, error) {
	var allVideos []*model.Video
	for _, keyword := range keywords {
		videos, err := s.videoRepo.GetVideosByKeyword(keyword)
		if err != nil {
			return nil, err
		}
		allVideos = append(allVideos, videos...)
	}
	return allVideos, nil
}

// CalculatePopularityScore calculates the popularity score for the provided keywords based on video analytics.
func (s *videoServiceImpl) CalculatePopularityScore(keywords []string) (map[string]float64, error) {
	scores := make(map[string]float64)

	videos, err := s.videoRepo.GetVideoAnalytics(keywords)
	if err != nil {
		return nil, err
	}

	if len(videos) == 0 {
		for _, keyword := range keywords {
			scores[keyword] = 0
		}
		return scores, nil
	}

	// Calculate the median of log-transformed view counts, like counts, and comment counts.
	medianViews := calculateMedian(logTransform(getViewCounts(videos)))
	medianLikes := calculateMedian(logTransform(getLikeCounts(videos)))
	medianComments := calculateMedian(logTransform(getCommentCounts(videos)))

	totalScores := make(map[string]float64)
	videoCounts := make(map[string]int)

	for _, video := range videos {
		logViewCount := math.Log(float64(video.ViewCount) + 1)
		logLikeCount := math.Log(float64(video.LikeCount) + 1)
		logCommentCount := math.Log(float64(video.CommentCount) + 1)

		zScoreViews := (logViewCount - medianViews) / medianViews
		zScoreLikes := (logLikeCount - medianLikes) / medianLikes
		zScoreComments := (logCommentCount - medianComments) / medianComments

		// Calculate the weighted score using Z-scores.
		score := 0.6*zScoreViews + 0.3*zScoreLikes + 0.1*zScoreComments

		for _, keyword := range keywords {
			if strings.Contains(video.ChannelName, keyword) || strings.Contains(video.Description, keyword) {
				totalScores[keyword] += score
				videoCounts[keyword]++
			}
		}
	}

	for _, keyword := range keywords {
		if videoCounts[keyword] == 0 {
			scores[keyword] = 0
		} else {
			scores[keyword] = totalScores[keyword] / float64(videoCounts[keyword])
		}
	}

	return scores, nil
}

// logTransform applies logarithmic transformation to a slice of int64 values.
// Used to reduce data scale and resolve large differences in the number of impressions, likes, or comments.
func logTransform(values []int64) []float64 {
	result := make([]float64, len(values))
	for i, value := range values {
		result[i] = math.Log(float64(value) + 1)
	}
	return result
}

// calculateMedian calculates the median of a slice of float64 values.
// The median is the middle value of sorted data and is used to reduce the influence of outliers
func calculateMedian(values []float64) float64 {
	n := len(values)
	if n == 0 {
		return 0
	}
	sorted := make([]float64, n)
	copy(sorted, values)
	sort.Float64s(sorted)
	if n%2 == 1 {
		return sorted[n/2]
	}
	return (sorted[n/2-1] + sorted[n/2]) / 2
}

// getViewCounts extracts the view counts from a slice of Video objects.
func getViewCounts(videos []*model.Video) []int64 {
	var counts []int64
	for _, video := range videos {
		counts = append(counts, video.ViewCount)
	}
	return counts
}

// getLikeCounts extracts the like counts from a slice of Video objects.
func getLikeCounts(videos []*model.Video) []int64 {
	var counts []int64
	for _, video := range videos {
		counts = append(counts, video.LikeCount)
	}
	return counts
}

// getCommentCounts extracts the comment counts from a slice of Video objects.
func getCommentCounts(videos []*model.Video) []int64 {
	var counts []int64
	for _, video := range videos {
		counts = append(counts, video.CommentCount)
	}
	return counts
}
