import { useEffect, useRef, useState } from 'react'
import { Card, Alert, Button, Divider, List, Row, Col, Space, message } from 'antd'

export default function DevicePanel({ title = 'Подключение к устройству' }) {
  const [wsStatus, setWsStatus] = useState('disconnected')
  const [photos, setPhotos] = useState([])
  const [deviceLogs, setDeviceLogs] = useState([])
  const socketRef = useRef(null)

  useEffect(() => {
    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
    const socket = new WebSocket(`${protocol}://${window.location.host}/ws`)
    socketRef.current = socket

    socket.onopen = () => {
      setWsStatus('connected')
      pushLog('Подключено к устройству')
      socket.send(JSON.stringify({ type: 'get_photos' }))
    }

    socket.onclose = () => {
      setWsStatus('disconnected')
      pushLog('Подключение закрыто')
    }

    socket.onerror = () => {
      pushLog('Ошибка подключения к устройству')
    }

    socket.onmessage = (event) => {
      try {
        const payload = JSON.parse(event.data)
        handleDeviceMessage(payload)
      } catch (error) {
        pushLog('Не удалось обработать сообщение устройства')
      }
    }

    return () => {
      socket.close()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const pushLog = (text) => {
    setDeviceLogs((prev) => [text, ...prev].slice(0, 25))
  }

  const handleDeviceMessage = (payload) => {
    switch (payload?.type) {
      case 'photo_list':
        if (Array.isArray(payload.data)) {
          setPhotos(payload.data)
        }
        pushLog('Получен список фотографий')
        break
      case 'new_photo':
        if (payload.data) {
          setPhotos((prev) => [payload.data, ...prev])
        }
        pushLog('Добавлено новое фото')
        break
      case 'photo_saved':
        pushLog('Фото сохранено')
        break
      case 'stream_started':
        pushLog('Трансляция запущена')
        break
      case 'stream_stopped':
        pushLog('Трансляция остановлена')
        break
      case 'error':
        pushLog(payload?.data?.message || 'Ошибка устройства')
        message.error(payload?.data?.message || 'Ошибка устройства')
        break
      default:
        if (payload?.type) {
          pushLog(`Событие: ${payload.type}`)
        }
    }
  }

  const sendDeviceCommand = (type, data = {}) => {
    if (!socketRef.current || socketRef.current.readyState !== WebSocket.OPEN) {
      message.warning('Нет подключения к устройству')
      return
    }
    socketRef.current.send(JSON.stringify({ type, data }))
  }

  return (
    <Card className="mb-6">
      <Row gutter={16}>
        <Col span={24}>
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold">{title}</h2>
            <Alert
              message={wsStatus === 'connected' ? 'Устройство подключено' : 'Нет соединения'}
              type={wsStatus === 'connected' ? 'success' : 'warning'}
              showIcon
            />
          </div>
          <Space wrap>
            <Button onClick={() => sendDeviceCommand('start_stream')}>Начать трансляцию</Button>
            <Button onClick={() => sendDeviceCommand('stop_stream')}>Остановить трансляцию</Button>
            <Button onClick={() => sendDeviceCommand('get_photos')}>Обновить список фото</Button>
          </Space>
        </Col>
      </Row>
      <Divider />
      <Row gutter={16}>
        <Col xs={24} md={12}>
          <Card title="Фотографии" size="small">
            <List
              dataSource={photos}
              locale={{ emptyText: 'Нет данных' }}
              renderItem={(item) => (
                <List.Item actions={[
                  <a key="open" href={item?.path || '#'} target="_blank" rel="noreferrer">
                    Открыть
                  </a>,
                ]}>
                  <List.Item.Meta
                    title={item?.filename || 'Фото'}
                    description={item?.timestamp ? new Date(item.timestamp).toLocaleString() : ''}
                  />
                </List.Item>
              )}
            />
          </Card>
        </Col>
        <Col xs={24} md={12}>
          <Card title="Журнал событий" size="small">
            <List
              dataSource={deviceLogs}
              locale={{ emptyText: 'Пока нет событий' }}
              renderItem={(item, index) => <List.Item key={index}>{item}</List.Item>}
            />
          </Card>
        </Col>
      </Row>
    </Card>
  )
}
