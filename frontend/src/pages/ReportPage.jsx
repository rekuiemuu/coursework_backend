import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Spin, Descriptions, Button, message, Tag, Divider } from 'antd'
import { ArrowLeftOutlined, FileTextOutlined } from '@ant-design/icons'
import { reportsAPI, examinationsAPI, patientsAPI } from '../api'

export default function ReportPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [report, setReport] = useState(null)
  const [examination, setExamination] = useState(null)
  const [patient, setPatient] = useState(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    loadReport()
  }, [id])

  const loadReport = async () => {
    setLoading(true)
    try {
      const response = await reportsAPI.get(id)
      const reportData = response.data.data
      setReport(reportData)

      if (reportData.examination_id) {
        try {
          const examRes = await examinationsAPI.get(reportData.examination_id)
          const examData = examRes.data.data
          setExamination(examData)

          if (examData.patient_id) {
            const patientRes = await patientsAPI.get(examData.patient_id)
            setPatient(patientRes.data.data)
          }
        } catch (err) {
          console.error('Error loading related data:', err)
        }
      }
    } catch (error) {
      message.error('Ошибка загрузки отчета')
    } finally {
      setLoading(false)
    }
  }

  if (loading || !report) {
    return <Spin size="large" />
  }

  const getDuration = () => {
    if (!examination || !examination.created_at || !examination.completed_at) return '-'
    const start = new Date(examination.created_at)
    const end = new Date(examination.completed_at)
    const minutes = Math.floor((end - start) / 60000)
    return minutes > 0 ? `${minutes} мин` : '< 1 мин'
  }

  return (
    <div>
      <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/reports')} style={{ marginBottom: 16 }}>
        К списку отчетов
      </Button>
      
      <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 24 }}>
        <FileTextOutlined style={{ fontSize: 32, color: '#dc2626' }} />
        <h1 style={{ fontSize: 28, fontWeight: 600, margin: 0 }}>{report.title}</h1>
      </div>

      {patient && examination && (
        <Card style={{ marginBottom: 24 }} title="Информация об исследовании">
          <Descriptions bordered column={2}>
            <Descriptions.Item label="Пациент">
              {patient.last_name} {patient.first_name} {patient.middle_name}
            </Descriptions.Item>
            <Descriptions.Item label="Дата рождения">
              {new Date(patient.date_of_birth).toLocaleDateString('ru-RU')}
            </Descriptions.Item>
            <Descriptions.Item label="Дата исследования">
              {new Date(examination.created_at).toLocaleString('ru-RU')}
            </Descriptions.Item>
            <Descriptions.Item label="Длительность">
              {getDuration()}
            </Descriptions.Item>
            <Descriptions.Item label="Статус исследования">
              <Tag color="green">{examination.status}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label="Количество изображений">
              {examination.images?.length || 0}
            </Descriptions.Item>
          </Descriptions>
        </Card>
      )}

      <Card title="Результаты анализа">
        <Descriptions bordered column={1}>
          <Descriptions.Item label="Краткое содержание">
            {report.summary || 'Не указано'}
          </Descriptions.Item>
          <Descriptions.Item label="Диагноз">
            {report.diagnosis || 'Не указан'}
          </Descriptions.Item>
          <Descriptions.Item label="Рекомендации">
            {report.recommendations || 'Не указаны'}
          </Descriptions.Item>
        </Descriptions>
        
        <Divider />
        
        <div>
          <h3 style={{ fontSize: 16, fontWeight: 600, marginBottom: 12 }}>Подробное описание</h3>
          <div style={{ whiteSpace: 'pre-wrap', lineHeight: 1.6 }}>
            {report.content || 'Не указано'}
          </div>
        </div>

        <Divider />

        <Descriptions bordered column={2} size="small">
          <Descriptions.Item label="Создан">
            {new Date(report.created_at).toLocaleString('ru-RU')}
          </Descriptions.Item>
          <Descriptions.Item label="Обновлен">
            {new Date(report.updated_at).toLocaleString('ru-RU')}
          </Descriptions.Item>
        </Descriptions>
      </Card>
    </div>
  )
}
