import React from 'react';
import { Card, Table, Typography } from 'antd';
import { useQuery } from '@tanstack/react-query';
import axios from '../utils/axios';
import type { ColumnsType } from 'antd/es/table';

const { Title } = Typography;

interface ServiceStats {
  name: string;
  instanceCount: number;
  averageResponseTime: number;
  successRate: number;
  lastUpdated: string;
}

const Stats: React.FC = () => {
  const { data: stats, isLoading } = useQuery<ServiceStats[]>({
    queryKey: ['serviceStats'],
    queryFn: async () => {
      const response = await axios.get<{ stats: ServiceStats[] }>('/api/services/stats');
      return response.data.stats ?? [];
    },
  });

  const columns: ColumnsType<ServiceStats> = [
    {
      title: '服务名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '实例数量',
      dataIndex: 'instanceCount',
      key: 'instanceCount',
      sorter: (a: ServiceStats, b: ServiceStats) => a.instanceCount - b.instanceCount,
    },
    {
      title: '平均响应时间(ms)',
      dataIndex: 'averageResponseTime',
      key: 'averageResponseTime',
      sorter: (a: ServiceStats, b: ServiceStats) => a.averageResponseTime - b.averageResponseTime,
    },
    {
      title: '成功率',
      dataIndex: 'successRate',
      key: 'successRate',
      render: (rate: number) => `${(rate * 100).toFixed(2)}%`,
      sorter: (a: ServiceStats, b: ServiceStats) => a.successRate - b.successRate,
    },
    {
      title: '最后更新时间',
      dataIndex: 'lastUpdated',
      key: 'lastUpdated',
      render: (time: string) => new Date(time).toLocaleString(),
    },
  ];

  return (
    <div>
      <Title level={2}>服务统计</Title>
      <Card>
        <Table<ServiceStats>
          columns={columns}
          dataSource={stats ?? []}
          loading={isLoading}
          rowKey="name"
          pagination={false}
        />
      </Card>
    </div>
  );
};

export default Stats; 