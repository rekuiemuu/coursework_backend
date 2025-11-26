import { useState } from 'react'
import { Form, Input, Button, Card, message, Tabs } from 'antd'
import { UserOutlined, LockOutlined, MailOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { authAPI } from '../api'

export default function LoginPage() {
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  const onLogin = async (values) => {
    setLoading(true)
    try {
      const response = await authAPI.login(values)
      localStorage.setItem('token', response.data.data.token)
      localStorage.setItem('user', JSON.stringify(response.data.data.user))
      message.success('Вход выполнен успешно')
      navigate('/')
    } catch (error) {
      message.error('Ошибка входа')
    } finally {
      setLoading(false)
    }
  }

  const onRegister = async (values) => {
    setLoading(true)
    try {
      await authAPI.register(values)
      message.success('Регистрация успешна. Теперь войдите в систему')
      const loginData = { username: values.username, password: values.password }
      const response = await authAPI.login(loginData)
      localStorage.setItem('token', response.data.data.token)
      localStorage.setItem('user', JSON.stringify(response.data.data.user))
      navigate('/')
    } catch (error) {
      message.error('Ошибка регистрации')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <Card title="Капилляроскопия" className="w-96">
        <Tabs
          items={[
            {
              key: 'login',
              label: 'Вход',
              children: (
                <Form onFinish={onLogin} layout="vertical">
                  <Form.Item name="username" rules={[{ required: true, message: 'Введите логин' }]}>
                    <Input prefix={<UserOutlined />} placeholder="Логин" />
                  </Form.Item>
                  <Form.Item name="password" rules={[{ required: true, message: 'Введите пароль' }]}>
                    <Input.Password prefix={<LockOutlined />} placeholder="Пароль" />
                  </Form.Item>
                  <Form.Item>
                    <Button type="primary" htmlType="submit" loading={loading} block>
                      Войти
                    </Button>
                  </Form.Item>
                </Form>
              ),
            },
            {
              key: 'register',
              label: 'Регистрация',
              children: (
                <Form onFinish={onRegister} layout="vertical">
                  <Form.Item name="username" rules={[{ required: true, message: 'Введите логин' }]}>
                    <Input prefix={<UserOutlined />} placeholder="Логин" />
                  </Form.Item>
                  <Form.Item name="email" rules={[{ required: true, type: 'email', message: 'Введите email' }]}>
                    <Input prefix={<MailOutlined />} placeholder="Email" />
                  </Form.Item>
                  <Form.Item name="password" rules={[{ required: true, min: 6, message: 'Минимум 6 символов' }]}>
                    <Input.Password prefix={<LockOutlined />} placeholder="Пароль" />
                  </Form.Item>
                  <Form.Item name="first_name" rules={[{ required: true, message: 'Введите имя' }]}>
                    <Input placeholder="Имя" />
                  </Form.Item>
                  <Form.Item name="last_name" rules={[{ required: true, message: 'Введите фамилию' }]}>
                    <Input placeholder="Фамилия" />
                  </Form.Item>
                  <Form.Item name="role" initialValue="doctor" hidden>
                    <Input />
                  </Form.Item>
                  <Form.Item>
                    <Button type="primary" htmlType="submit" loading={loading} block>
                      Зарегистрироваться
                    </Button>
                  </Form.Item>
                </Form>
              ),
            },
          ]}
        />
      </Card>
    </div>
  )
}
