import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { Card, Button, Spin, message, Descriptions, Tag } from 'antd'
import { examinationsAPI } from '../api'

export default function ExaminationDetailPage() {
  const { id } = useParams()
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
      <h1 className="text-2xl font-bold mb-4">Исследование {examination.id}</h1>
      <Card className="mb-4">
        <Descriptions bordered column={2}>
          <Descriptions.Item label="ID">{examination.id}</Descriptions.Item>
          <Descriptions.Item label="Пациент">{examination.patient_id}</Descriptions.Item>
          <Descriptions.Item label="Врач">{examination.doctor_id}</Descriptions.Item>
          <Descriptions.Item label="Статус">
            <Tag>{examination.status}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="Описание" span={2}>
            {examination.description}
          </Descriptions.Item>
        </Descriptions>
      </Card>
      <Button type="primary" onClick={startAnalysis}>
        Запустить анализ
      </Button>
    </div>
  )
}
