package mq

import (
	"encoding/json"
)

type AnalysisTaskMessage struct {
	AnalysisID    string `json:"analysis_id"`
	ExaminationID string `json:"examination_id"`
	ImageID       string `json:"image_id"`
	ImagePath     string `json:"image_path"`
}

func NewAnalysisTaskMessage(analysisID, examinationID, imageID, imagePath string) *AnalysisTaskMessage {
	return &AnalysisTaskMessage{
		AnalysisID:    analysisID,
		ExaminationID: examinationID,
		ImageID:       imageID,
		ImagePath:     imagePath,
	}
}

func (a *AnalysisTaskMessage) ToJSON() ([]byte, error) {
	return json.Marshal(a)
}

func ParseAnalysisTaskMessage(data []byte) (*AnalysisTaskMessage, error) {
	var msg AnalysisTaskMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
