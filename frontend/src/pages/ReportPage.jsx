import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Spin, Descriptions, Button, message, Tag, Divider, Image, Modal, Input, Form } from 'antd'
import { ArrowLeftOutlined, FileTextOutlined, EditOutlined } from '@ant-design/icons'
import { reportsAPI, examinationsAPI, patientsAPI } from '../api'

const { TextArea } = Input

export default function ReportPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [report, setReport] = useState(null)
  const [examination, setExamination] = useState(null)
  const [patient, setPatient] = useState(null)
  const [loading, setLoading] = useState(false)
  const [editModalVisible, setEditModalVisible] = useState(false)
  const [form] = Form.useForm()

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

  const handleEdit = () => {
    form.setFieldsValue({
      content: report.content,
      summary: report.summary,
      diagnosis: report.diagnosis,
      recommendations: report.recommendations,
    })
    setEditModalVisible(true)
  }

  const handleSave = async () => {
    try {
      const values = await form.validateFields()
      await reportsAPI.update(id, values)
      message.success('Отчет обновлен')
      setEditModalVisible(false)
      loadReport()
    } catch (error) {
      message.error('Ошибка обновления отчета')
    }
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/reports')}>
          К списку отчетов
        </Button>
        <Button type="primary" icon={<EditOutlined />} onClick={handleEdit}>
          Редактировать
        </Button>
      </div>
      
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

      {report.images && report.images.length > 0 && (
        <Card title="Изображения исследования" style={{ marginTop: 24 }}>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(200px, 1fr))', gap: 16 }}>
            <Image.PreviewGroup>
              {report.images.map((img) => (
                <Image
                  key={img.id}
                  src={img.url}
                  alt={img.filename}
                  style={{ width: '100%', height: 200, objectFit: 'cover', borderRadius: 8 }}
                  fallback="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
                />
              ))}
            </Image.PreviewGroup>
          </div>
        </Card>
      )}

      <Modal
        title="Редактировать отчет"
        open={editModalVisible}
        onOk={handleSave}
        onCancel={() => setEditModalVisible(false)}
        width={800}
        okText="Сохранить"
        cancelText="Отмена"
      >
        <Form form={form} layout="vertical">
          <Form.Item name="summary" label="Краткое содержание">
            <TextArea rows={3} />
          </Form.Item>
          <Form.Item name="diagnosis" label="Диагноз">
            <TextArea rows={3} />
          </Form.Item>
          <Form.Item name="recommendations" label="Рекомендации">
            <TextArea rows={3} />
          </Form.Item>
          <Form.Item name="content" label="Подробное описание">
            <TextArea rows={6} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
