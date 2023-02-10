package controller

import (
	"github.com/google/uuid"
	"main/helpers"
	"main/manager"
	"main/model"
	"time"

	"github.com/gin-gonic/gin"
)

type videoAddRequest struct {
	URL      string `json:"url" binding:"required"`
	Duration string `json:"duration" binding:"required"`
}

type videoGetResponse struct {
	URL         string    `json:"url" binding:"required"`
	Duration    string    `json:"duration" binding:"required"`
	UserCreated uuid.UUID `json:"userCreated" binding:"required"`
}

type VideoController struct {
	videoManager *manager.VideoManager
}

func NewVideoController(videoManager *manager.VideoManager) *VideoController {
	return &VideoController{videoManager: videoManager}
}

func (v *VideoController) GetVideo(c *gin.Context) {
	vid, err := v.videoManager.GetVideo(c.Request.URL.Query().Get("video_url"))
	if err != nil {
		responseWithError(c, 500, "Could not get video: "+err.Error())
		return
	}

	data := &videoGetResponse{
		URL:         vid.URL,
		Duration:    vid.Duration.String(),
		UserCreated: vid.UserCreated,
	}

	c.JSON(200, data)
}

func (v *VideoController) CreateVideo(c *gin.Context) {
	var data videoAddRequest

	if err := c.BindJSON(&data); err != nil {
		responseWithError(c, 406, "Please provide valid videoAddRequest metadata: "+err.Error())
		return
	}

	if !helpers.IsValidURL(data.URL) {
		responseWithError(c, 406, "Invalid URL")
		return
	}

	duration, err := time.ParseDuration(data.Duration)
	if err != nil {
		responseWithError(c, 400, "Invalid duration: "+err.Error())
		return
	}
	if duration.Nanoseconds() < 0 {
		responseWithError(c, 400, "Negative duration: "+duration.String())
		return
	}

	usr, ok := c.Get("User")
	if !ok {
		responseWithError(c, 500, "No user in the context")
	}

	vid := &model.Video{
		URL:         data.URL,
		Duration:    duration,
		UserCreated: usr.(*model.User).ID,
	}

	if err := v.videoManager.CreateVideo(vid); err != nil {
		responseWithError(c, 500, "Problem saving your videoAddRequest: "+err.Error())
		return
	}

	c.JSON(200, gin.H{"message": "videoAddRequest saved"})
}

func (v *VideoController) DeleteVideo(c *gin.Context) {
	usr, ok := c.Get("User")
	if !ok {
		responseWithError(c, 500, "No user in the context")
	}

	if err := v.videoManager.DeleteVideo(c.Request.URL.Query().Get("video_url"), usr.(*model.User).ID); err != nil {
		responseWithError(c, 500, "Problem deleting your videoAddRequest: "+err.Error())
		return
	}

	c.JSON(200, gin.H{"message": "videoAddRequest deleted"})
}
