package ws

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type PhotoInfo struct {
	Filename  string    `json:"filename"`
	Timestamp time.Time `json:"timestamp"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
}

type DeviceControlData struct {
	Type  string  `json:"type"`
	Value float64 `json:"value"`
	Delta float64 `json:"delta"`
}

type DeviceManager struct {
	clients          map[*websocket.Conn]bool
	broadcast        chan []byte
	mutex            sync.Mutex
	streaming        bool
	photoStoragePath string
	upgrader         websocket.Upgrader
	onPhotoSaved     func(filename string)
	onControlChange  func(controlData DeviceControlData)
}

func NewDeviceManager(photoStoragePath string) *DeviceManager {
	return &DeviceManager{
		clients:          make(map[*websocket.Conn]bool),
		broadcast:        make(chan []byte, 256),
		streaming:        false,
		photoStoragePath: photoStoragePath,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (dm *DeviceManager) SetPhotoSavedCallback(callback func(filename string)) {
	dm.onPhotoSaved = callback
}

func (dm *DeviceManager) SetControlChangeCallback(callback func(controlData DeviceControlData)) {
	dm.onControlChange = callback
}

func (dm *DeviceManager) AddClient(conn *websocket.Conn) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()
	dm.clients[conn] = true
	log.Printf("Client connected. Total clients: %d", len(dm.clients))
}

func (dm *DeviceManager) RemoveClient(conn *websocket.Conn) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()
	if _, ok := dm.clients[conn]; ok {
		delete(dm.clients, conn)
		conn.Close()
		log.Printf("Client disconnected. Total clients: %d", len(dm.clients))
	}
}

func (dm *DeviceManager) BroadcastMessage(message []byte) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	for client := range dm.clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Error sending to client: %v", err)
			client.Close()
			delete(dm.clients, client)
		}
	}
}

func (dm *DeviceManager) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := dm.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	dm.AddClient(conn)
	defer dm.RemoveClient(conn)

	photos, err := dm.GetPhotoList()
	if err == nil {
		dm.SendResponse(conn, "photo_list", photos)
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("JSON parse error: %v", err)
			continue
		}

		dm.HandleMessage(conn, msg)
	}
}

func (dm *DeviceManager) HandleMessage(conn *websocket.Conn, msg Message) {
	switch msg.Type {
	case "start_stream":
		dm.SendResponse(conn, "stream_started", "Use browser camera access")

	case "stop_stream":
		dm.StopStreaming()
		dm.SendResponse(conn, "stream_stopped", "Streaming stopped")

	case "save_photo":
		dataMap, ok := msg.Data.(map[string]interface{})
		if !ok {
			dm.SendError(conn, "Invalid photo data")
			return
		}

		imageData, ok := dataMap["image"].(string)
		if !ok {
			dm.SendError(conn, "Missing image data")
			return
		}

		filename, err := dm.SavePhotoFromBase64(imageData)
		if err != nil {
			dm.SendError(conn, fmt.Sprintf("Failed to save photo: %v", err))
		} else {
			dm.SendResponse(conn, "photo_saved", map[string]string{
				"filename": filename,
				"url":      fmt.Sprintf("/api/photos/%s", filename),
			})
			dm.BroadcastNewPhoto(filename)

			if dm.onPhotoSaved != nil {
				dm.onPhotoSaved(filename)
			}
		}

	case "get_photos":
		photos, err := dm.GetPhotoList()
		if err != nil {
			dm.SendError(conn, fmt.Sprintf("Failed to get photos: %v", err))
		} else {
			dm.SendResponse(conn, "photo_list", photos)
		}

	case "delete_photo":
		dataMap, ok := msg.Data.(map[string]interface{})
		if !ok {
			dm.SendError(conn, "Invalid data")
			return
		}

		filename, ok := dataMap["filename"].(string)
		if !ok {
			dm.SendError(conn, "Missing filename")
			return
		}

		err := dm.DeletePhoto(filename)
		if err != nil {
			dm.SendError(conn, fmt.Sprintf("Failed to delete photo: %v", err))
		} else {
			dm.SendResponse(conn, "photo_deleted", map[string]string{"filename": filename})
		}

	case "control_change":
		dataMap, ok := msg.Data.(map[string]interface{})
		if ok {
			controlData := DeviceControlData{
				Type:  dataMap["type"].(string),
				Value: dataMap["value"].(float64),
			}
			if delta, exists := dataMap["delta"]; exists {
				controlData.Delta = delta.(float64)
			}

			dm.BroadcastControlChange(controlData)

			if dm.onControlChange != nil {
				dm.onControlChange(controlData)
			}
		}

	default:
		dm.SendError(conn, "Unknown message type")
	}
}

func (dm *DeviceManager) SendResponse(conn *websocket.Conn, msgType string, data interface{}) {
	response := Message{
		Type: msgType,
		Data: data,
	}
	jsonData, _ := json.Marshal(response)
	conn.WriteMessage(websocket.TextMessage, jsonData)
}

func (dm *DeviceManager) SendError(conn *websocket.Conn, errMsg string) {
	dm.SendResponse(conn, "error", map[string]string{"message": errMsg})
}

func (dm *DeviceManager) StopStreaming() {
	dm.streaming = false
	log.Println("Streaming stopped")
}

func (dm *DeviceManager) SavePhotoFromBase64(base64Data string) (string, error) {
	if err := os.MkdirAll(dm.photoStoragePath, 0755); err != nil {
		return "", err
	}

	if len(base64Data) > 23 && base64Data[:23] == "data:image/jpeg;base64," {
		base64Data = base64Data[23:]
	} else if len(base64Data) > 22 && base64Data[:22] == "data:image/png;base64," {
		base64Data = base64Data[22:]
	}

	imageBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	_, _, err = image.Decode(io.NopCloser(bytes.NewReader(imageBytes)))
	if err != nil {
		return "", fmt.Errorf("invalid image data: %w", err)
	}

	filename := fmt.Sprintf("photo_%s.jpg", time.Now().Format("20060102_150405"))
	filePath := filepath.Join(dm.photoStoragePath, filename)

	err = os.WriteFile(filePath, imageBytes, 0644)
	if err != nil {
		return "", err
	}

	log.Printf("Photo saved: %s (size: %d bytes)", filename, len(imageBytes))
	return filename, nil
}

func (dm *DeviceManager) DeletePhoto(filename string) error {
	filename = filepath.Base(filename)
	filePath := filepath.Join(dm.photoStoragePath, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found")
	}

	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	log.Printf("Photo deleted: %s", filename)
	return nil
}

func (dm *DeviceManager) BroadcastNewPhoto(filename string) {
	msg := Message{
		Type: "new_photo",
		Data: map[string]string{
			"filename": filename,
			"url":      fmt.Sprintf("/api/photos/%s", filename),
		},
	}
	jsonData, _ := json.Marshal(msg)
	dm.BroadcastMessage(jsonData)
}

func (dm *DeviceManager) BroadcastControlChange(data DeviceControlData) {
	msg := Message{
		Type: "control_change",
		Data: data,
	}
	jsonData, _ := json.Marshal(msg)
	dm.BroadcastMessage(jsonData)
	log.Printf("Broadcasting control change: %v", data)
}

func (dm *DeviceManager) GetPhotoList() ([]PhotoInfo, error) {
	var photos []PhotoInfo

	if _, err := os.Stat(dm.photoStoragePath); os.IsNotExist(err) {
		return photos, nil
	}

	files, err := os.ReadDir(dm.photoStoragePath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			info, err := file.Info()
			if err != nil {
				continue
			}

			ext := filepath.Ext(file.Name())
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
				photos = append(photos, PhotoInfo{
					Filename:  file.Name(),
					Timestamp: info.ModTime(),
					Path:      fmt.Sprintf("/api/photos/%s", file.Name()),
					Size:      info.Size(),
				})
			}
		}
	}

	return photos, nil
}

func (dm *DeviceManager) ServePhoto(w http.ResponseWriter, r *http.Request, filename string) {
	if filename == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}

	filepath := filepath.Join(dm.photoStoragePath, filepath.Base(filename))

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filepath)
}
