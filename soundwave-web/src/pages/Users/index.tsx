import React from 'react';
import { Card, Table, Tag, Button, Modal, Form, Input, Select, message } from 'antd';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import type { ColumnsType } from 'antd/es/table';
import axios from '../../utils/axios';

interface User {
  id: string;
  username: string;
  role: string;
  permissions: string[];
  createdAt: string;
}

const Users: React.FC = () => {
  const [form] = Form.useForm();
  const [modalVisible, setModalVisible] = React.useState(false);
  const queryClient = useQueryClient();

  const { data: users, isLoading } = useQuery<User[]>({
    queryKey: ['users'],
    queryFn: async () => {
      const response = await axios.get('/api/users');
      return response.data.users;
    },
  });

  const createUser = useMutation({
    mutationFn: (values: any) => axios.post('/api/users', values),
    onSuccess: () => {
      message.success('用户创建成功');
      setModalVisible(false);
      form.resetFields();
      queryClient.invalidateQueries({ queryKey: ['users'] });
    },
  });

  const columns: ColumnsType<User> = [
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
    },
    {
      title: '角色',
      dataIndex: 'role',
      key: 'role',
      render: (role: string) => (
        <Tag color={role === 'admin' ? 'red' : role === 'user' ? 'blue' : 'green'}>
          {role === 'admin' ? '管理员' : role === 'user' ? '普通用户' : '测试用户'}
        </Tag>
      ),
    },
    {
      title: '权限',
      dataIndex: 'permissions',
      key: 'permissions',
      render: (permissions: string[]) => (
        <>
          {permissions.map(perm => (
            <Tag key={perm} color="blue">{perm}</Tag>
          ))}
        </>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (time: string) => new Date(time).toLocaleString(),
    },
  ];

  return (
    <div>
      <Card
        title="用户管理"
        extra={
          <Button type="primary" onClick={() => setModalVisible(true)}>
            添加用户
          </Button>
        }
      >
        <Table<User>
          columns={columns}
          dataSource={users}
          loading={isLoading}
          rowKey="id"
        />
      </Card>

      <Modal
        title="添加用户"
        open={modalVisible}
        onOk={() => form.submit()}
        onCancel={() => setModalVisible(false)}
        confirmLoading={createUser.isPending}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={values => createUser.mutate(values)}
        >
          <Form.Item
            name="username"
            label="用户名"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="password"
            label="密码"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password />
          </Form.Item>
          <Form.Item
            name="role"
            label="角色"
            rules={[{ required: true, message: '请选择角色' }]}
          >
            <Select>
              <Select.Option value="user">普通用户</Select.Option>
              <Select.Option value="tester">测试用户</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default Users; 