import { Layout, Menu, Button } from 'antd'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { UserOutlined, ExperimentOutlined, LogoutOutlined } from '@ant-design/icons'

const { Header, Content, Footer } = Layout

export default function MainLayout() {
  const navigate = useNavigate()
  const location = useLocation()

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    navigate('/login')
  }

  const menuItems = [
    { key: '/patients', label: 'Пациенты', icon: <UserOutlined /> },
    { key: '/examinations', label: 'Исследования', icon: <ExperimentOutlined /> },
  ]

  return (
    <Layout className="min-h-screen">
      <Header className="flex items-center justify-between">
        <div className="flex items-center flex-1">
          <div className="text-white text-xl font-bold mr-8">Капилляроскопия</div>
          <Menu
            theme="dark"
            mode="horizontal"
            selectedKeys={[location.pathname]}
            items={menuItems}
            onClick={({ key }) => navigate(key)}
            className="flex-1"
          />
        </div>
        <Button
          type="text"
          icon={<LogoutOutlined />}
          onClick={handleLogout}
          className="text-white hover:text-gray-300"
        >
          Выход
        </Button>
      </Header>
      <Content className="p-8">
        <div className="bg-white p-6 rounded-lg min-h-full">
          <Outlet />
        </div>
      </Content>
      <Footer className="text-center">
        Капилляроскопия ©2025
      </Footer>
    </Layout>
  )
}
