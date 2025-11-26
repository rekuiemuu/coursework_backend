import { useState, useEffect } from 'react'
import { Table, Button, Modal, Form, Input, DatePicker, Select, message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { patientsAPI } from '../api'

export default function PatientsPage() {
  const [patients, setPatients] = useState([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [form] = Form.useForm()

  useEffect(() => {
    loadPatients()
  }, [])

  const loadPatients = async () => {
    setLoading(true)
    try {
      const response = await patientsAPI.list()
      setPatients(response.data.data || [])
    } catch (error) {
      message.error('Ошибка загрузки пациентов')
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = async (values) => {
    try {
      await patientsAPI.create(values)
      message.success('Пациент создан')
      setModalVisible(false)
      form.resetFields()
      loadPatients()
    } catch (error) {
      message.error('Ошибка создания пациента')
    }
  }

  const columns = [
    { title: 'Фамилия', dataIndex: 'last_name', key: 'last_name' },
    { title: 'Имя', dataIndex: 'first_name', key: 'first_name' },
    { title: 'Отчество', dataIndex: 'middle_name', key: 'middle_name' },
    { title: 'Пол', dataIndex: 'gender', key: 'gender', render: (g) => g === 'male' ? 'Муж' : 'Жен' },
    { title: 'Телефон', dataIndex: 'phone', key: 'phone' },
    { title: 'Email', dataIndex: 'email', key: 'email' },
  ]

  return (
    <div>
      <div className="flex justify-between mb-4">
        <h1 className="text-2xl font-bold">Пациенты</h1>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalVisible(true)}>
          Добавить
        </Button>
      </div>
      <Table dataSource={patients} columns={columns} loading={loading} rowKey="id" />
      <Modal
        title="Новый пациент"
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        onOk={() => form.submit()}
      >
        <Form form={form} onFinish={handleCreate} layout="vertical">
          <Form.Item name="last_name" label="Фамилия" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="first_name" label="Имя" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="middle_name" label="Отчество">
            <Input />
          </Form.Item>
          <Form.Item name="date_of_birth" label="Дата рождения" rules={[{ required: true }]}>
            <DatePicker style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="gender" label="Пол" rules={[{ required: true }]}>
            <Select>
              <Select.Option value="male">Мужской</Select.Option>
              <Select.Option value="female">Женский</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="phone" label="Телефон">
            <Input />
          </Form.Item>
          <Form.Item name="email" label="Email">
            <Input />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
