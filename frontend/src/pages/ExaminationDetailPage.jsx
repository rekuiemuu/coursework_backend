import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Button, Spin, message, Descriptions, Tag, Divider, Space } from 'antd'
import { ArrowLeftOutlined } from '@ant-design/icons'
import { examinationsAPI } from '../api'
import DevicePanel from '../components/DevicePanel'

export default function ExaminationDetailPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [examination, setExamination] = useState(null)
  const [loading, setLoading] = useState(false)
  useEffect(() => {
    loadExamination()
  }, [id])

  const loadExamination = async () => {
    setLoading(true)
    try {
      const response = await examinationsAPI.get(id)
      setExamination(response.data.data)
    } catch (error) {
      message.error('Ошибка загрузки исследования')
    } finally {
      setLoading(false)
    }
  }

  const startAnalysis = async () => {
    try {
      await examinationsAPI.startAnalysis(id)
      message.success('Анализ запущен')
      loadExamination()
    } catch (error) {
      message.error('Ошибка запуска анализа')
    }
  }

  if (loading || !examination) {
    return <Spin size="large" />
  }

  return (
    <div>
      <div className="flex justify-between mb-4">
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/examinations')}>
          Назад
        </Button>
        <Space>
          <Button onClick={loadExamination}>Обновить</Button>
          <Button type="primary" onClick={startAnalysis}>
            Запустить анализ
          </Button>
        </Space>
      </div>
      <h1 className="text-2xl font-bold mb-4">Детали исследования</h1>
      <Card className="mb-4">
        <Descriptions bordered column={2}>
          <Descriptions.Item label="ID">{examination.id}</Descriptions.Item>
          <Descriptions.Item label="Пациент">{examination.patient_id}</Descriptions.Item>
          <Descriptions.Item label="Врач">{examination.doctor_id}</Descriptions.Item>
          <Descriptions.Item label="Статус">
            <Tag>{examination.status}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="Создано">{new Date(examination.created_at).toLocaleString('ru-RU')}</Descriptions.Item>
          <Descriptions.Item label="Обновлено">{new Date(examination.updated_at).toLocaleString('ru-RU')}</Descriptions.Item>
          <Descriptions.Item label="Описание" span={2}>
            {examination.description || 'Не указано'}
          </Descriptions.Item>
        </Descriptions>
      </Card>
      <Divider />
      <DevicePanel title="Управление камерой для исследования" />
    </div>
  )
}
