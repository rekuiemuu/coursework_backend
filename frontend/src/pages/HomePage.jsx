import { Card, Statistic, Row, Col } from 'antd'
import { UserOutlined, ExperimentOutlined, FileTextOutlined } from '@ant-design/icons'

export default function HomePage() {
  return (
    <div>
      <h1 className="text-3xl font-bold mb-6">Система анализа капилляроскопии</h1>
      <Row gutter={16}>
        <Col span={8}>
          <Card>
            <Statistic
              title="Пациенты"
              value={0}
              prefix={<UserOutlined />}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="Исследования"
              value={0}
              prefix={<ExperimentOutlined />}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="Отчёты"
              value={0}
              prefix={<FileTextOutlined />}
            />
          </Card>
        </Col>
      </Row>
    </div>
  )
}
