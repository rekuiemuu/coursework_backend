import { useState, useEffect } from 'react'
import { Table, Button, Card, Tag, message, Empty, Space } from 'antd'
import { FileTextOutlined, EyeOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { reportsAPI, examinationsAPI, patientsAPI } from '../api'

export default function ReportsPage() {
  const [reports, setReports] = useState([])
  const [examinations, setExaminations] = useState({})
  const [patients, setPatients] = useState({})
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    setLoading(true)
    try {
      const [reportsRes, examsRes, patientsRes] = await Promise.all([
        reportsAPI.list?.() || Promise.resolve({ data: { data: [] } }),
        examinationsAPI.list({ limit: 100 }),
        patientsAPI.list({ limit: 100 })
      ])

      setReports(reportsRes.data.data || [])
      
      const examsMap = {}
      ;(examsRes.data.data || []).forEach(exam => {
        examsMap[exam.id] = exam
      })
      setExaminations(examsMap)

      const patientsMap = {}
      ;(patientsRes.data.data || []).forEach(patient => {
        patientsMap[patient.id] = patient
      })
      setPatients(patientsMap)
    } catch (error) {
      message.error('Ошибка загрузки данных')
    } finally {
      setLoading(false)
    }
  }

  const getPatientName = (examination) => {
    if (!examination) return '-'
    const patient = patients[examination.patient_id]
    return patient ? `${patient.last_name} ${patient.first_name}` : examination.patient_id
  }

  const getDuration = (examination) => {
    if (!examination || !examination.created_at || !examination.completed_at) return '-'
    const start = new Date(examination.created_at)
    const end = new Date(examination.completed_at)
    const minutes = Math.floor((end - start) / 60000)
    return minutes > 0 ? `${minutes} мин` : '< 1 мин'
  }

  const columns = [
    {
      title: 'Название отчета',
      dataIndex: 'title',
      key: 'title',
      render: (text) => (
        <Space>
          <FileTextOutlined style={{ color: '#dc2626' }} />
          <span style={{ fontWeight: 500 }}>{text}</span>
        </Space>
      )
    },
    {
      title: 'Пациент',
      dataIndex: 'examination_id',
      key: 'patient',
      render: (examId) => getPatientName(examinations[examId])
    },
    {
      title: 'Дата исследования',
      dataIndex: 'examination_id',
      key: 'exam_date',
      render: (examId) => {
        const exam = examinations[examId]
        return exam ? new Date(exam.created_at).toLocaleString('ru-RU') : '-'
      }
    },
    {
      title: 'Длительность',
      dataIndex: 'examination_id',
      key: 'duration',
      render: (examId) => getDuration(examinations[examId])
    },
    {
      title: 'Статус',
      dataIndex: 'examination_id',
      key: 'status',
      render: (examId) => {
        const exam = examinations[examId]
        if (!exam) return <Tag>-</Tag>
        const colors = { pending: 'orange', in_progress: 'blue', completed: 'green', failed: 'red' }
        return <Tag color={colors[exam.status] || 'default'}>{exam.status}</Tag>
      }
    },
    {
      title: 'Создан',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date) => new Date(date).toLocaleString('ru-RU')
    },
    {
      title: 'Действия',
      key: 'actions',
      render: (_, record) => (
        <Button
          type="link"
          icon={<EyeOutlined />}
          onClick={() => navigate(`/reports/${record.id}`)}
        >
          Открыть
        </Button>
      )
    }
  ]

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <h1 style={{ fontSize: 24, fontWeight: 600, margin: 0 }}>Отчеты исследований</h1>
        <Button onClick={loadData}>Обновить</Button>
      </div>

      {reports.length === 0 && !loading ? (
        <Card>
          <Empty
            description="Нет отчетов"
            image={Empty.PRESENTED_IMAGE_SIMPLE}
          >
            <p style={{ color: '#999', marginTop: 16 }}>
              Отчеты создаются автоматически после завершения анализа исследований
            </p>
          </Empty>
        </Card>
      ) : (
        <Table
          dataSource={reports}
          columns={columns}
          loading={loading}
          rowKey="id"
          pagination={{ pageSize: 20 }}
        />
      )}
    </div>
  )
}
