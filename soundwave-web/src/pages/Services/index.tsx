import React, { useEffect, useState } from 'react';
import { Table, Tag, Card, Row, Col, Statistic, Button, Space, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { ReloadOutlined, ApiOutlined, CloudServerOutlined, CheckCircleOutlined } from '@ant-design/icons';
import axios from 'axios';

const { Title } = Typography;

interface Service {
  name: string;
  id: string;
  hostname: string;
  ip: string;
  port: number;
  status: string;
  metadata: Record<string, string>;
  last_heartbeat: string;
  version: string;
}

const ServicesPage: React.FC = () => {
  const [services, setServices] = useState<Service[]>([]);
  const [loading, setLoading] = useState(false);

  const columns: ColumnsType<Service> = [
    {
      title: '服务名称',
      dataIndex: 'name',
      key: 'name',
      render: (text: string) => (
        <Space>
          <ApiOutlined style={{ color: '#1890ff' }} />
          <span>{text}</span>
        </Space>
      ),
    },
    {
      title: '实例ID',
      dataIndex: 'id',
      key: 'id',
      width: 200,
      ellipsis: true,
    },
    {
      title: '主机名',
      dataIndex: 'hostname',
      key: 'hostname',
    },
    {
      title: 'IP地址',
      dataIndex: 'ip',
      key: 'ip',
      render: (text: string, record: Service) => (
        <Tag color="blue">{`${text}:${record.port}`}</Tag>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'UP' ? 'success' : 'error'}>
          {status === 'UP' ? <CheckCircleOutlined /> : null} {status}
        </Tag>
      ),
    },
    {
      title: '版本',
      dataIndex: 'version',
      key: 'version',
      render: (version: string) => (
        <Tag color="purple">{version || '-'}</Tag>
      ),
    },
    {
      title: '最后心跳时间',
      dataIndex: 'last_heartbeat',
      key: 'last_heartbeat',
      render: (time: string) => new Date(time).toLocaleString(),
    },
  ];

  const fetchServices = async () => {
    try {
      setLoading(true);
      const response = await axios.get<{ services: Record<string, Service[]> }>('http://localhost:7777/services');
      const servicesList = Object.values(response.data.services).flat();
      setServices(servicesList);
    } catch (error) {
      console.error('获取服务列表失败:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchServices();
    const interval = setInterval(fetchServices, 10000);
    return () => clearInterval(interval);
  }, []);

  const getStatistics = () => {
    const totalServices = new Set(services.map(s => s.name)).size;
    const totalInstances = services.length;
    const healthyInstances = services.filter(s => s.status === 'UP').length;
    
    return { totalServices, totalInstances, healthyInstances };
  };

  const stats = getStatistics();

  return (
    <div>
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={8}>
          <Card bordered={false} className="overview-card">
            <Statistic
              title="服务总数"
              value={stats.totalServices}
              prefix={<ApiOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card bordered={false} className="overview-card">
            <Statistic
              title="实例总数"
              value={stats.totalInstances}
              prefix={<CloudServerOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card bordered={false} className="overview-card">
            <Statistic
              title="健康实例"
              value={stats.healthyInstances}
              prefix={<CheckCircleOutlined />}
              valueStyle={{ color: '#13c2c2' }}
            />
          </Card>
        </Col>
      </Row>

      <Card
        style={{ marginTop: 16 }}
        title={
          <Space>
            <CloudServerOutlined />
            <span>服务列表</span>
          </Space>
        }
        extra={
          <Button
            type="primary"
            icon={<ReloadOutlined />}
            onClick={fetchServices}
            loading={loading}
          >
            刷新
          </Button>
        }
      >
        <Table
          columns={columns}
          dataSource={services}
          rowKey={(record) => `${record.name}-${record.id}`}
          loading={loading}
          pagination={{
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => `共 ${total} 条`,
          }}
        />
      </Card>
    </div>
  );
};

export default ServicesPage; 