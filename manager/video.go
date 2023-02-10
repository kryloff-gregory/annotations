package manager

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"main/annotation"
	"main/model"
	"main/services"
	"main/video"
)

type VideoManager struct {
	videoDAO      video.DAO
	annotationDAO annotation.DAO
	lockService   *services.LockService
}

func NewVideoManager(videoDAO video.DAO, lockService *services.LockService) *VideoManager {
	return &VideoManager{videoDAO: videoDAO, lockService: lockService}
}

func (m *VideoManager) CreateVideo(video *model.Video) error {
	//We are blocking here as we have nothing provided by LevelDB to atomically "write if not exist"
	m.lockService.LockItem(video.URL)
	defer m.lockService.UnlockItem(video.URL)

	exists, err := m.videoDAO.VideoExists(video.URL)
	if err != nil {
		return err
	}
	if exists {
		return errors.New(fmt.Sprintf("Video with url %v already exists", video.URL))
	}
	return m.videoDAO.UpsertVideo(video)
}

func (m *VideoManager) GetVideo(videoURL string) (*model.Video, error) {
	m.lockService.RLockItem(videoURL)
	defer m.lockService.RUnlockItem(videoURL)

	return m.videoDAO.FetchVideo(videoURL)
}

func (m *VideoManager) DeleteVideo(videoURL string, userID uuid.UUID) error {
	m.lockService.LockItem(videoURL)
	defer m.lockService.UnlockItem(videoURL)

	vid, err := m.videoDAO.FetchVideo(videoURL)
	if err != nil {
		return err
	}
	if vid.UserCreated != userID {
		return errors.New("user is not authorized to delete the video")
	}

	anns, err := m.annotationDAO.FetchAnnotationsForVideo(videoURL)
	if err != nil {
		return err
	}

	if len(anns) == 0 {
		return errors.New("there are annotations for the video to be deleted")
	}

	return m.videoDAO.RemoveVideo(videoURL)
}
