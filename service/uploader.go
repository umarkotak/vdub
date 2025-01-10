package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// Credentials structure for API keys and tokens

type (
	// VideoMetadata contains information about the video
	VideoMetadata struct {
		Title       string
		Description string
		Tags        []string
		Privacy     string // "private", "public", or "unlisted" for YouTube
	}

	// VideoUploader handles video uploads to different platforms
	VideoUploader struct {
		youtubeService *youtube.Service
		httpClient     *http.Client
	}

	GenericUploadParams struct {
		VideoPath   string
		Title       string
		Description string
	}

	GenericUploadData struct {
		YoutubeVideoID   string `json:"youtube_video_id"`
		YoutubeVideoLink string `json:"youtube_video_link"`
	}
)

// NewVideoUploader creates a new instance of VideoUploader
func NewVideoUploader() (*VideoUploader, error) {
	// Initialize YouTube client
	config := &oauth2.Config{
		ClientID:     config.Get().YoutubeClientID,
		ClientSecret: config.Get().YoutubeClientSecret,
		RedirectURL:  "https://vdubb.vercel.app/google/redir",
		Scopes: []string{
			youtube.YoutubeUploadScope,
		},
		Endpoint: google.Endpoint,
	}

	// For this example, we assume you have a token file saved
	// In practice, you'd implement the full OAuth2 flow
	token, err := getTokenFromConfig()
	if err != nil {
		err := fmt.Errorf("error getting YouTube token: %v", err)
		logrus.Error(err)
		return nil, err
	}

	// token := getTokenFromWeb(config)

	// tokenJson, _ := json.Marshal(token)

	// logrus.Infof("TOKEN: %+v", string(tokenJson))

	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		err := fmt.Errorf("error creating YouTube service: %v", err)
		logrus.Error(err)
		return nil, err
	}

	return &VideoUploader{
		youtubeService: youtubeService,
		httpClient:     &http.Client{},
	}, nil
}

// UploadToYouTube uploads a video to YouTube
func (vu *VideoUploader) UploadToYouTube(videoPath string, metadata VideoMetadata) (string, error) {
	file, err := os.Open(videoPath)
	if err != nil {
		return "", fmt.Errorf("error opening video file: %v", err)
	}
	defer file.Close()

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       metadata.Title,
			Description: metadata.Description,
			Tags:        metadata.Tags,
		},
		Status: &youtube.VideoStatus{
			PrivacyStatus: metadata.Privacy,
		},
	}

	call := vu.youtubeService.Videos.Insert([]string{"snippet", "status"}, upload)
	call.Media(file)

	response, err := call.Do()
	if err != nil {
		return "", fmt.Errorf("error uploading to YouTube: %v", err)
	}

	return response.Id, nil
}

// UploadToTikTok uploads a video to TikTok
func (vu *VideoUploader) UploadToTikTok(videoPath string, metadata VideoMetadata) (string, error) {
	// Note: This is a simplified version. TikTok's actual upload process requires multiple steps
	// including initiating an upload, getting an upload URL, and confirming the upload

	// TikTok API endpoint (this is a placeholder - use actual endpoint)
	endpoint := "https://open.tiktokapis.com/v2/video/upload/"

	file, err := os.Open(videoPath)
	if err != nil {
		return "", fmt.Errorf("error opening video file: %v", err)
	}
	defer file.Close()

	// Create multipart form request
	req, err := http.NewRequest("POST", endpoint, file)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Get().TiktokAccessToken))
	req.Header.Set("Content-Type", "video/mp4")

	resp, err := vu.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error uploading to TikTok: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("TikTok API error: %s", string(body))
	}

	// Parse response to get video ID
	var result struct {
		VideoID string `json:"video_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error parsing TikTok response: %v", err)
	}

	return result.VideoID, nil
}

// Helper function to get token from file
func getTokenFromConfig() (*oauth2.Token, error) {
	var token oauth2.Token
	err := json.Unmarshal([]byte(config.Get().YoutubeAccountOauthJson), &token)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"json": config.Get().YoutubeAccountOauthJson,
		}).Error(err)
		return nil, err
	}
	return &token, nil
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		logrus.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		logrus.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func Upload(ctx context.Context, params GenericUploadParams) (GenericUploadData, error) {
	uploader, err := NewVideoUploader()
	if err != nil {
		err = fmt.Errorf("error creating uploader: %v", err)
		logrus.WithContext(ctx).Error(err)
		return GenericUploadData{}, err
	}

	// videoPath := "shared/task-public-bapak para lantern/dubbed_video.mp4"
	metadata := VideoMetadata{
		Title:       params.Title,
		Description: params.Description,
		Tags:        []string{},
		Privacy:     "public",
	}

	// Upload to YouTube
	youtubeID, err := uploader.UploadToYouTube(params.VideoPath, metadata)
	if err != nil {
		err = fmt.Errorf("error uploading to YouTube: %v", err)
		logrus.WithContext(ctx).Error(err)
		return GenericUploadData{}, err
	}

	logrus.Infof("Successfully uploaded to YouTube. Video ID: %s", youtubeID)

	return GenericUploadData{
		YoutubeVideoID:   youtubeID,
		YoutubeVideoLink: fmt.Sprintf("https://www.youtube.com/watch?v=%s", youtubeID),
	}, nil
}
