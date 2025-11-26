import { Layout, Menu } from 'antd'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { UserOutlined, ExperimentOutlined, LogoutOutlined } from '@ant-design/icons'

const { Sider, Content } = Layout

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
    { type: 'divider' },
    { key: 'logout', label: 'Выход', icon: <LogoutOutlined />, danger: true },
  ]

  const handleMenuClick = ({ key }) => {
    if (key === 'logout') {
      handleLogout()
    } else {
      navigate(key)
    }
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider
        width={240}
        style={{
          background: '#7f1d1d',
          boxShadow: '2px 0 8px rgba(0,0,0,0.1)',
        }}
      >
        <div style={{
          height: 64,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          color: 'white',
          fontSize: 18,
          fontWeight: 600,
          borderBottom: '1px solid rgba(255,255,255,0.1)',
        }}>
          Капилляроскопия
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={handleMenuClick}
          style={{ background: 'transparent', border: 'none' }}
        />
      </Sider>
      <Layout>
        <Content style={{ padding: 24, background: '#f5f5f5' }}>
          <div style={{ background: 'white', padding: 24, borderRadius: 8, minHeight: '100%' }}>
            <Outlet />
          </div>
        </Content>
      </Layout>
    </Layout>
  )
}
