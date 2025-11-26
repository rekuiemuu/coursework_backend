import { Layout, Menu } from 'antd'
import { Link, Outlet, useNavigate } from 'react-router-dom'
import { HomeOutlined, UserOutlined, FileTextOutlined, ExperimentOutlined } from '@ant-design/icons'

const { Header, Content, Footer } = Layout

export default function MainLayout() {
  const navigate = useNavigate()

  const menuItems = [
    { key: '/', label: 'Главная', icon: <HomeOutlined /> },
    { key: '/patients', label: 'Пациенты', icon: <UserOutlined /> },
    { key: '/examinations', label: 'Исследования', icon: <ExperimentOutlined /> },
  ]

  return (
    <Layout className="min-h-screen">
      <Header className="flex items-center">
        <div className="text-white text-xl font-bold mr-8">Капилляроскопия</div>
        <Menu
          theme="dark"
          mode="horizontal"
          items={menuItems}
          onClick={({ key }) => navigate(key)}
          className="flex-1"
        />
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
