import React from 'react';
import { Form, Input, Button, Card, message } from 'antd';
import { UserOutlined, LockOutlined, CloudServerOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useUserStore } from '../../store/userStore';
import styles from './Login.module.css';

interface LoginForm {
  username: string;
  password: string;
}

const Login: React.FC = () => {
  const navigate = useNavigate();
  const login = useUserStore((state) => state.login);
  const [loading, setLoading] = React.useState(false);

  const onFinish = async (values: LoginForm) => {
    setLoading(true);
    try {
      await login(values.username, values.password);
      navigate('/');
    } catch (error) {
      // 错误处理已在 store 中完成
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.container}>
      <Card className={styles.card} bordered={false}>
        <div className={styles.logo}>
          <CloudServerOutlined className={styles.logoIcon} />
          <h1 className={styles.title}>SOUNDWAVE</h1>
          <p className={styles.subtitle}>服务注册与发现中心</p>
        </div>
        <Form
          name="login"
          onFinish={onFinish}
          autoComplete="off"
          className={styles.form}
          size="large"
        >
          <Form.Item
            name="username"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input
              prefix={<UserOutlined style={{ color: '#1890ff' }} />}
              placeholder="用户名"
              className={styles.input}
              disabled={loading}
            />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password
              prefix={<LockOutlined style={{ color: '#1890ff' }} />}
              placeholder="密码"
              className={styles.input}
              disabled={loading}
            />
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              block
              loading={loading}
              className={styles.button}
            >
              登录
            </Button>
          </Form.Item>
        </Form>
        <div className={styles.footer}>
          <p>© 2024 Soundwave. All rights reserved.</p>
          <p>
            Need help? <a href="#">Contact support</a>
          </p>
        </div>
      </Card>
    </div>
  );
};

export default Login; 