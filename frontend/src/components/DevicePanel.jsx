import { useEffect, useRef, useState } from 'react'
import { Card, Button, List, message, Select, Tag, Checkbox } from 'antd'
import { CameraOutlined, PlayCircleOutlined, StopOutlined, ReloadOutlined } from '@ant-design/icons'
import { examinationsAPI } from '../api'

export default function DevicePanel({ title = 'Управление камерой', examinationId }) {
  const [wsStatus, setWsStatus] = useState('disconnected')
  const [photos, setPhotos] = useState([])
  const [selectedPhotos, setSelectedPhotos] = useState([])
  const [cameras, setCameras] = useState([])
  const [selectedCamera, setSelectedCamera] = useState('')
  const [streaming, setStreaming] = useState(false)
  const videoRef = useRef(null)
  const streamRef = useRef(null)
  const socketRef = useRef(null)

  useEffect(() => {
    requestCameraPermission()
    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
    const socket = new WebSocket(`${protocol}://${window.location.host}/ws`)
    socketRef.current = socket

    socket.onopen = () => {
      setWsStatus('connected')
      socket.send(JSON.stringify({ type: 'get_photos' }))
    }

    socket.onclose = () => {
      setWsStatus('disconnected')
    }

    socket.onmessage = (event) => {
      try {
        const payload = JSON.parse(event.data)
        handleDeviceMessage(payload)
      } catch (error) {
        console.error('WebSocket message error:', error)
      }
    }

    return () => {
      stopCamera()
      socket.close()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const requestCameraPermission = async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ video: true })
      stream.getTracks().forEach(track => track.stop())
      await loadCameras()
    } catch (error) {
      message.error('Нет доступа к камере. Разрешите доступ в браузере.')
    }
  }

  const loadCameras = async () => {
    try {
      const devices = await navigator.mediaDevices.enumerateDevices()
      const videoDevices = devices.filter(d => d.kind === 'videoinput')
      setCameras(videoDevices)
      
      const microscope = videoDevices.find(d => {
        const label = d.label.toLowerCase()
        return label.includes('microsope') || label.includes('04f2') || label.includes('3008')
      })
      
      if (microscope) {
        setSelectedCamera(microscope.deviceId)
        message.info(`Найден капилляроскоп: ${microscope.label}`)
      } else if (videoDevices.length > 0) {
        setSelectedCamera(videoDevices[0].deviceId)
      }
    } catch (error) {
      message.error('Ошибка загрузки камер')
    }
  }

  const startCamera = async () => {
    if (!selectedCamera) {
      message.warning('Выберите камеру')
      return
    }
    try {
      if (streamRef.current) {
        streamRef.current.getTracks().forEach(track => track.stop())
      }
      const stream = await navigator.mediaDevices.getUserMedia({
        video: { deviceId: { exact: selectedCamera }, width: 1280, height: 720 }
      })
      streamRef.current = stream
      if (videoRef.current) {
        videoRef.current.srcObject = stream
      }
      setStreaming(true)
      message.success('Камера запущена')
    } catch (error) {
      message.error('Ошибка доступа к камере')
    }
  }

  const stopCamera = () => {
    if (streamRef.current) {
      streamRef.current.getTracks().forEach(track => track.stop())
      streamRef.current = null
      if (videoRef.current) {
        videoRef.current.srcObject = null
      }
    }
    setStreaming(false)
  }

  const takePhoto = () => {
    if (!streamRef.current || !videoRef.current) {
      message.warning('Сначала запустите камеру')
      return
    }
    const canvas = document.createElement('canvas')
    canvas.width = videoRef.current.videoWidth
    canvas.height = videoRef.current.videoHeight
    canvas.getContext('2d').drawImage(videoRef.current, 0, 0)
    const imageData = canvas.toDataURL('image/jpeg', 0.95)
    sendDeviceCommand('save_photo', { image: imageData })
  }

  const handleDeviceMessage = (payload) => {
    switch (payload?.type) {
      case 'photo_list':
        if (Array.isArray(payload.data)) {
          setPhotos(payload.data)
        }
        break
      case 'new_photo':
      case 'photo_saved':
        sendDeviceCommand('get_photos')
        message.success('Фото сохранено')
        break
      case 'error':
        message.error(payload?.data?.message || 'Ошибка устройства')
        break
      default:
        break
    }
  }

  const sendDeviceCommand = (type, data = {}) => {
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      socketRef.current.send(JSON.stringify({ type, data }))
    }
  }

  const attachPhotosToExamination = async () => {
    try {
      await examinationsAPI.attachPhotos(examinationId, selectedPhotos)
      message.success(`Прикреплено ${selectedPhotos.length} фото к исследованию`)
      
      // Запускаем анализ
      await examinationsAPI.startAnalysis(examinationId)
      message.success('Анализ запущен! Результаты появятся в отчётах через несколько секунд')
      setSelectedPhotos([])
    } catch (error) {
      message.error('Ошибка прикрепления фото: ' + (error.response?.data?.message || error.message))
    }
  }

  return (
    <Card 
      title={<span style={{ fontSize: 16, fontWeight: 600 }}>{title}</span>}
      extra={
        <Tag color={wsStatus === 'connected' ? 'success' : 'error'}>
          {wsStatus === 'connected' ? 'Подключено' : 'Отключено'}
        </Tag>
      }
      style={{ marginBottom: 24 }}
    >
      <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: 24 }}>
        <div>
          <video
            ref={videoRef}
            autoPlay
            playsInline
            muted
            style={{
              width: '100%',
              height: 400,
              background: '#000',
              borderRadius: 8,
              objectFit: 'contain'
            }}
          />
          <div style={{ marginTop: 16, display: 'flex', gap: 8, flexWrap: 'wrap' }}>
            <Select
              value={selectedCamera}
              onChange={setSelectedCamera}
              style={{ flex: 1, minWidth: 200 }}
              placeholder="Выберите камеру"
              notFoundContent="Камеры не найдены"
            >
              {cameras.map(cam => {
                const isMicroscope = cam.label.toLowerCase().includes('microsope') || 
                                     cam.label.toLowerCase().includes('04f2') ||
                                     cam.label.toLowerCase().includes('3008')
                return (
                  <Select.Option key={cam.deviceId} value={cam.deviceId}>
                    {isMicroscope ? '🔬 ' : '🎥 '}
                    {cam.label || `Камера ${cam.deviceId.substring(0, 8)}`}
                  </Select.Option>
                )
              })}
            </Select>
            {!streaming ? (
              <Button type="primary" danger icon={<PlayCircleOutlined />} onClick={startCamera}>
                Запустить
              </Button>
            ) : (
              <Button danger icon={<StopOutlined />} onClick={stopCamera}>
                Остановить
              </Button>
            )}
            <Button icon={<CameraOutlined />} onClick={takePhoto} disabled={!streaming}>
              Фото
            </Button>
            <Button icon={<ReloadOutlined />} onClick={() => sendDeviceCommand('get_photos')}>
              Обновить
            </Button>
          </div>
        </div>
        <div>
          <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 12 }}>
            Сохраненные фото {photos.length > 0 && `(${photos.length})`}
          </div>
          <List
            size="small"
            dataSource={photos}
            locale={{ emptyText: 'Нет фотографий' }}
            style={{ maxHeight: 450, overflow: 'auto' }}
            renderItem={(item) => (
              <List.Item
                actions={[
                  <a key="view" href={item?.path} target="_blank" rel="noreferrer">
                    Открыть
                  </a>
                ]}
              >
                {examinationId && (
                  <Checkbox
                    checked={selectedPhotos.includes(item.filename)}
                    onChange={(e) => {
                      if (e.target.checked) {
                        setSelectedPhotos([...selectedPhotos, item.filename])
                      } else {
                        setSelectedPhotos(selectedPhotos.filter(p => p !== item.filename))
                      }
                    }}
                    style={{ marginRight: 8 }}
                  />
                )}
                <List.Item.Meta
                  title={<span style={{ fontSize: 13 }}>{item?.filename}</span>}
                  description={
                    <span style={{ fontSize: 12 }}>
                      {item?.timestamp ? new Date(item.timestamp).toLocaleString('ru-RU') : ''}
                    </span>
                  }
                />
              </List.Item>
            )}
          />
          {examinationId && selectedPhotos.length > 0 && (
            <Button 
              type="primary" 
              danger
              block 
              style={{ marginTop: 12 }}
              onClick={attachPhotosToExamination}
            >
              Прикрепить {selectedPhotos.length} фото к анализу
            </Button>
          )}
        </div>
      </div>
    </Card>
  )
}
