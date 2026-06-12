import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { Layout as AntLayout, Menu, Button, Avatar, Dropdown, Space } from 'antd';
import {
  DashboardOutlined,
  BookOutlined,
  UserOutlined,
  FileTextOutlined,
  LogoutOutlined,
} from '@ant-design/icons';

const { Header, Content, Sider } = AntLayout;

interface LayoutProps {
  setIsAuthenticated: (value: boolean) => void;
}

const Layout = ({ setIsAuthenticated }: LayoutProps) => {
  const navigate = useNavigate();
  const location = useLocation();
  const user = JSON.parse(localStorage.getItem('user') || '{}');

  const menuItems = [
    { key: '/', icon: <DashboardOutlined />, label: '工作台' },
    { key: '/courses', icon: <BookOutlined />, label: '课程管理' },
    { key: '/students', icon: <UserOutlined />, label: '学生管理' },
    { key: '/summary', icon: <FileTextOutlined />, label: '学习总结' },
  ];

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setIsAuthenticated(false);
    navigate('/login');
  };

  const userMenu = {
    items: [
      {
        key: 'logout',
        icon: <LogoutOutlined />,
        label: '退出登录',
        onClick: handleLogout,
      },
    ],
  };

  return (
    <AntLayout className="min-h-screen">
      <Header className="bg-white shadow-sm flex items-center justify-between px-6">
        <h1 className="text-xl font-bold text-blue-600 m-0">学习管理平台</h1>
        <Dropdown menu={userMenu} placement="bottomRight">
          <Space className="cursor-pointer">
            <Avatar icon={<UserOutlined />} />
            <span>{user.name || '管理员'}</span>
          </Space>
        </Dropdown>
      </Header>
      <AntLayout>
        <Sider className="bg-white shadow-md" width={200}>
          <Menu
            mode="inline"
            selectedKeys={[location.pathname]}
            items={menuItems}
            onClick={({ key }) => navigate(key)}
            className="h-full border-r-0"
          />
        </Sider>
        <Content className="m-4 p-6 bg-gray-50 rounded-lg overflow-auto">
          <Outlet />
        </Content>
      </AntLayout>
    </AntLayout>
  );
};

export default Layout;