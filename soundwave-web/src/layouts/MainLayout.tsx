import React, { useState, useEffect } from 'react';
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  UserOutlined,
  DesktopOutlined,
  FileOutlined,
  CloudServerOutlined,
  LogoutOutlined,
  LockOutlined,
} from '@ant-design/icons';
import { Layout, Menu, Button, theme, Typography, Dropdown, message } from 'antd';
import { Outlet, useNavigate } from 'react-router-dom';
import axios from '../utils/axios';
import { useUserStore } from '../store/userStore';
import ChangePasswordModal from '../components/ChangePasswordModal';

const { Header, Sider, Content } = Layout;
const { Title } = Typography;

interface User {
  username: string;
  role: string;
  permissions: string[];
}

const iconMap: Record<string, React.ReactNode> = {
  DesktopOutlined: <DesktopOutlined />,
  FileOutlined: <FileOutlined />,
  UserOutlined: <UserOutlined />,
  CloudServerOutlined: <CloudServerOutlined />,
};

const getIcon = (iconName: string): React.ReactNode => {
  return iconMap[iconName] || <CloudServerOutlined />;
};

const MainLayout: React.FC = () => {
  const [collapsed, setCollapsed] = useState(false);
  const [menus, setMenus] = useState<any[]>([]);
  const { user, setUser, logout } = useUserStore();
  const navigate = useNavigate();
  const [changePasswordVisible, setChangePasswordVisible] = useState(false);
  
  const {
    token: { colorBgContainer, borderRadiusLG },
  } = theme.useToken();

  useEffect(() => {
    // 从 localStorage 获取用户信息
    const userStr = localStorage.getItem('user');
    if (userStr) {
      setUser(JSON.parse(userStr));
    } else {
      navigate('/login');
      return;
    }

    // 获取菜单数据
    fetchMenus();
  }, [navigate, setUser]);

  const fetchMenus = async () => {
    try {
      const response = await axios.get('/api/menus');
      setMenus(response.data.menus);
    } catch (error) {
      message.error('获取菜单失败');
    }
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const userMenuItems = [
    {
      key: 'changePassword',
      icon: <LockOutlined />,
      label: '修改密码',
      onClick: () => setChangePasswordVisible(true),
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: handleLogout,
    },
  ];

  const menuItems = menus.map(menu => ({
    key: menu.path.substring(1),
    icon: getIcon(menu.icon),
    label: menu.name,
  }));

  return (
    <>
      <Layout style={{ minHeight: '100vh' }}>
        <Sider trigger={null} collapsible collapsed={collapsed}>
          <div style={{ 
            height: 32, 
            margin: 16, 
            background: 'rgba(255, 255, 255, 0.2)',
            borderRadius: 6
          }} />
          <Menu
            theme="dark"
            mode="inline"
            defaultSelectedKeys={['services']}
            items={menuItems}
            onClick={({ key }) => navigate(key)}
          />
        </Sider>
        <Layout>
          <Header style={{ 
            padding: 0, 
            background: colorBgContainer,
            boxShadow: '0 1px 4px rgba(0,21,41,.08)',
            position: 'sticky',
            top: 0,
            zIndex: 100,
            width: '100%',
            display: 'flex',
            alignItems: 'center'
          }}>
            <Button
              type="text"
              icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
              onClick={() => setCollapsed(!collapsed)}
              style={{
                fontSize: '16px',
                width: 64,
                height: 64,
              }}
            />
            <div style={{ flex: 1, display: 'flex', justifyContent: 'center' }}>
              <div className="header-title">
                <CloudServerOutlined style={{ fontSize: '24px', color: '#1890ff', marginRight: '8px' }} />
                <Title level={3} style={{ 
                  margin: 0,
                  background: 'linear-gradient(to right, #1890ff, #52c41a)',
                  WebkitBackgroundClip: 'text',
                  WebkitTextFillColor: 'transparent',
                  fontFamily: "'Montserrat', sans-serif",
                  letterSpacing: '1px',
                  fontWeight: 600
                }}>
                  SOUNDWAVE
                </Title>
              </div>
            </div>
            <div style={{ float: 'right', marginRight: 24 }}>
              <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
                <Button type="text" icon={<UserOutlined />}>
                  {user?.username}
                </Button>
              </Dropdown>
            </div>
          </Header>
          <Content
            style={{
              margin: '24px',
              minHeight: 280,
              background: colorBgContainer,
              borderRadius: borderRadiusLG,
              padding: 24,
            }}
          >
            <Outlet />
          </Content>
        </Layout>
      </Layout>
      <ChangePasswordModal
        visible={changePasswordVisible}
        onClose={() => setChangePasswordVisible(false)}
      />
    </>
  );
};

export default MainLayout; 