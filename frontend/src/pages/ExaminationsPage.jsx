import { useState, useEffect } from 'react'
import { Table, Button, Tag, Modal, Form, Select, Input, message, Space } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { examinationsAPI, patientsAPI } from '../api'
import DevicePanel from '../components/DevicePanel'

const statusColors = {
  pending: 'default',
  in_progress: 'processing',
  completed: 'success',
  failed: 'error',
}

export default function ExaminationsPage() {
  const [examinations, setExaminations] = useState([])
  const [patients, setPatients] = useState([])
  const [currentUser] = useState(() => {
    try {
      return JSON.parse(localStorage.getItem('user') || '{}')
    } catch (error) {
      return {}
    }
  })
  const [loading, setLoading] = useState(false)
  const [createLoading, setCreateLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [form] = Form.useForm()
  const navigate = useNavigate()

  useEffect(() => {
    loadExaminations()
    loadPatients()
  }, [])

  useEffect(() => {
    if (modalVisible) {
      form.setFieldsValue({ doctor_id: currentUser?.id || undefined })
    }
  }, [modalVisible, currentUser, form])

  const loadExaminations = async () => {
    setLoading(true)
    try {
      const response = await examinationsAPI.list()
      setExaminations(response.data.data || [])
    } catch (error) {
      message.error('Не удалось загрузить исследования')
    } finally {
      setLoading(false)
    }
  }

  const loadPatients = async () => {
    try {
      const response = await patientsAPI.list({ limit: 100 })
      setPatients(response.data.data || [])
    } catch (error) {
      message.error('Не удалось загрузить пациентов')
    }
  }

  const getPatientName = (id) => {
    const patient = patients.find((item) => item.id === id)
    if (!patient) {
      return id
    }
    return `${patient.last_name || ''} ${patient.first_name || ''}`.trim()
  }

  const handleCreate = async (values) => {
    const doctorId = values.doctor_id || currentUser?.id
    if (!doctorId) {
      message.error('Укажите врача')
      return
    }

    setCreateLoading(true)
    try {
      await examinationsAPI.create({
        patient_id: values.patient_id,
        doctor_id: doctorId,
        description: values.description || '',
      })
      message.success('Исследование создано')
      setModalVisible(false)
      form.resetFields()
      loadExaminations()
    } catch (error) {
      message.error('Ошибка создания исследования')
    } finally {
      setCreateLoading(false)
    }
  }

  const columns = [
    {
      title: 'Пациент',
      dataIndex: 'patient_id',
      key: 'patient_id',
      render: (value) => getPatientName(value),
    },
    {
      title: 'Статус',
      dataIndex: 'status',
      key: 'status',
      render: (status) => <Tag color={statusColors[status] || 'default'}>{status}</Tag>,
    },
    {
      title: 'Описание',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: 'Создано',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (value) => (value ? new Date(value).toLocaleString('ru-RU') : ''),
    },
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
      <div className="flex justify-between mb-4">
        <h1 className="text-2xl font-bold">Исследования</h1>
        <Space>
          <Button onClick={loadExaminations}>Обновить</Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalVisible(true)}>
            Новое исследование
          </Button>
        </Space>
      </div>
      <DevicePanel />
      <Table dataSource={examinations} columns={columns} loading={loading} rowKey="id" />
      <Modal
        title="Создать исследование"
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        onOk={() => form.submit()}
        confirmLoading={createLoading}
      >
        <Form form={form} layout="vertical" onFinish={handleCreate}>
          <Form.Item name="patient_id" label="Пациент" rules={[{ required: true, message: 'Выберите пациента' }]}>
            <Select
              showSearch
              placeholder="Выберите пациента"
              options={patients.map((patient) => ({
                value: patient.id,
                label: `${patient.last_name || ''} ${patient.first_name || ''}`.trim() || patient.id,
              }))}
              optionFilterProp="label"
            />
          </Form.Item>
          <Form.Item
            name="doctor_id"
            label="Врач"
            initialValue={currentUser?.id}
            extra={currentUser?.id ? 'Определено автоматически из профиля' : 'Введите ID врача вручную'}
            rules={[{ required: !currentUser?.id, message: 'Укажите врача' }]}
          >
            <Input placeholder="ID врача" />
          </Form.Item>
          <Form.Item name="description" label="Описание">
            <Input.TextArea rows={4} placeholder="Например, предварительный диагноз" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
