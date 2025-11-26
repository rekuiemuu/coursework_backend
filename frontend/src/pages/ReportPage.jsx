import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { Card, Spin, Descriptions } from 'antd'
import { reportsAPI } from '../api'

export default function ReportPage() {
  const { id } = useParams()
  const [report, setReport] = useState(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    loadReport()
  }, [id])

  const loadReport = async () => {
    setLoading(true)
    try {
      const response = await reportsAPI.get(id)
      setReport(response.data.data)
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  if (loading || !report) {
    return <Spin size="large" />
  }

  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">{report.title}</h1>
      <Card>
        <Descriptions bordered column={1}>
          <Descriptions.Item label="Краткое содержание">{report.summary}</Descriptions.Item>
          <Descriptions.Item label="Диагноз">{report.diagnosis}</Descriptions.Item>
          <Descriptions.Item label="Рекомендации">{report.recommendations}</Descriptions.Item>
          <Descriptions.Item label="Содержание">{report.content}</Descriptions.Item>
        </Descriptions>
      </Card>
    </div>
  )
}
