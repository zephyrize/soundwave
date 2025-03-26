import React from 'react';
import { Card, Form, Input, Button, message, Switch, InputNumber, Space } from 'antd';
import { useQuery, useMutation } from '@tanstack/react-query';
import axios from '../../utils/axios';

interface SystemSettings {
  registryEnabled: boolean;
  heartbeatInterval: number;
  maxInstances: number;
  logLevel: string;
}

const Settings: React.FC = () => {
  const [form] = Form.useForm();

  const { data: settings, isLoading } = useQuery<SystemSettings>({
    queryKey: ['settings'],
    queryFn: async () => {
      const response = await axios.get('/api/settings');
      return response.data.settings;
    },
  });

  const updateSettings = useMutation({
    mutationFn: (values: SystemSettings) => axios.put('/api/settings', values),
    onSuccess: () => {
      message.success('设置更新成功');
    },
  });

  React.useEffect(() => {
    if (settings) {
      form.setFieldsValue(settings);
    }
  }, [settings, form]);

  return (
    <Card title="系统设置" loading={isLoading}>
      <Form
        form={form}
        layout="vertical"
        onFinish={values => updateSettings.mutate(values)}
      >
        <Form.Item
          name="registryEnabled"
          label="启用服务注册"
          valuePropName="checked"
        >
          <Switch />
        </Form.Item>

        <Form.Item
          name="heartbeatInterval"
          label="心跳间隔(秒)"
          rules={[{ required: true, message: '请输入心跳间隔' }]}
        >
          <InputNumber min={1} max={60} />
        </Form.Item>

        <Form.Item
          name="maxInstances"
          label="最大实例数"
          rules={[{ required: true, message: '请输入最大实例数' }]}
        >
          <InputNumber min={1} max={100} />
        </Form.Item>

        <Form.Item
          name="logLevel"
          label="日志级别"
          rules={[{ required: true, message: '请选择日志级别' }]}
        >
          <Input />
        </Form.Item>

        <Form.Item>
          <Space>
            <Button type="primary" htmlType="submit" loading={updateSettings.isPending}>
              保存设置
            </Button>
            <Button onClick={() => form.resetFields()}>
              重置
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Card>
  );
};

export default Settings; 