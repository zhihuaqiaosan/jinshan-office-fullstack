import { useEffect, useState } from 'react';
import { Card, Row, Col, Statistic, Spin } from 'antd';
import ReactECharts from 'echarts-for-react';
import { getDashboard } from '../api';

const Dashboard = () => {
  const [loading, setLoading] = useState(true);
  const [data, setData] = useState<any>(null);

  useEffect(() => {
    fetchDashboard();
  }, []);

  const fetchDashboard = async () => {
    try {
      const res = await getDashboard();
      if (res.data.code === 0) {
        setData(res.data.data);
      }
    } catch (error) {
      console.error('获取工作台数据失败', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div className="flex justify-center items-center h-96"><Spin size="large" /></div>;
  }

  if (!data) return null;

  const { stats, charts } = data;
  const publishRate = stats.totalCourses ? Math.round((stats.publishedCourses / stats.totalCourses) * 100) : 0;
  const activeRate = stats.totalStudents ? Math.round((stats.activeStudents / stats.totalStudents) * 100) : 0;

  // 柱状图配置
  const enrollmentOption = {
    title: { text: '课程选课人数 TOP 8', left: 'center' },
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    xAxis: { type: 'category', data: charts.enrollment.map((item: any) => item.name), axisLabel: { rotate: 30 } },
    yAxis: { type: 'value', name: '选课人数' },
    series: [{ type: 'bar', data: charts.enrollment.map((item: any) => item.value), itemStyle: { color: '#5470c6' } }],
  };

  // 折线图配置
  const activityOption = {
    title: { text: '近7天学习活跃度', left: 'center' },
    tooltip: { trigger: 'axis' },
    legend: { data: ['学习人数', '学习时长(小时)'], bottom: 0 },
    xAxis: { type: 'category', data: charts.activity.map((item: any) => item.label) },
    yAxis: [{ type: 'value', name: '学习人数' }, { type: 'value', name: '学习时长(小时)' }],
    series: [
      { name: '学习人数', type: 'line', data: charts.activity.map((item: any) => item.students), smooth: true, color: '#5470c6' },
      { name: '学习时长(小时)', type: 'line', data: charts.activity.map((item: any) => item.duration), smooth: true, color: '#fac858', yAxisIndex: 1 },
    ],
  };

  // 学生状态饼图
  const statusPieOption = {
    title: { text: '学生状态分布', left: 'center' },
    tooltip: { trigger: 'item' },
    series: [{ type: 'pie', radius: '50%', data: charts.statusDist, label: { show: true, formatter: '{b}: {d}%' } }],
  };

  // 课程分类饼图
  const categoryPieOption = {
    title: { text: '课程分类分布', left: 'center' },
    tooltip: { trigger: 'item' },
    series: [{ type: 'pie', radius: '50%', data: charts.categoryDist, label: { show: true, formatter: '{b}: {d}%' } }],
  };

  return (
    <div>
      <Row gutter={16} className="mb-6">
        <Col span={6}>
          <Card>
            <Statistic title="课程总数" value={stats.totalCourses} suffix={`/ 已发布 ${stats.publishedCourses}`} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="学生总数" value={stats.totalStudents} suffix={`/ 活跃 ${stats.activeStudents}`} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="课程发布率" value={publishRate} suffix="%" />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="学生活跃率" value={activeRate} suffix="%" />
          </Card>
        </Col>
      </Row>

      <Row gutter={16} className="mb-6">
        <Col span={12}>
          <Card>
            <ReactECharts option={enrollmentOption} style={{ height: 350 }} />
          </Card>
        </Col>
        <Col span={12}>
          <Card>
            <ReactECharts option={activityOption} style={{ height: 350 }} />
          </Card>
        </Col>
      </Row>

      <Row gutter={16}>
        <Col span={12}>
          <Card>
            <ReactECharts option={statusPieOption} style={{ height: 350 }} />
          </Card>
        </Col>
        <Col span={12}>
          <Card>
            <ReactECharts option={categoryPieOption} style={{ height: 350 }} />
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Dashboard;