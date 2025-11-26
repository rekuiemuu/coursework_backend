import { useState, useEffect } from 'react'
import { Table, Button, Tag } from 'antd'
import { useNavigate } from 'react-router-dom'
import { examinationsAPI } from '../api'

export default function ExaminationsPage() {
  const [examinations, setExaminations] = useState([])
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  useEffect(() => {
    loadExaminations()
  }, [])

  const loadExaminations = async () => {
    setLoading(true)
    try {
      const response = await examinationsAPI.list()
      setExaminations(response.data.data || [])
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const statusColors = {
    pending: 'default',
    in_progress: 'processing',
    completed: 'success',
    failed: 'error',
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id' },
    { title: 'Пациент', dataIndex: 'patient_id', key: 'patient_id' },
    {
      title: 'Статус',
      dataIndex: 'status',
      key: 'status',
      render: (status) => <Tag color={statusColors[status]}>{status}</Tag>,
    },
    { title: 'Дата', dataIndex: 'created_at', key: 'created_at' },
    {
      title: 'Действия',
      key: 'actions',
      render: (_, record) => (
        <Button type="link" onClick={() => navigate(`/examinations/${record.id}`)}>
          Открыть
        </Button>
      ),
    },
  ]

  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">Исследования</h1>
      <Table dataSource={examinations} columns={columns} loading={loading} rowKey="id" />
    </div>
  )
}
