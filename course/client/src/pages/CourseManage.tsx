import { useEffect, useState } from 'react';
import {
  Table, Button, Input, Select, Space, Modal, Form, message, Popconfirm, Tag, Card, Row, Col,
} from 'antd';
import { PlusOutlined, SearchOutlined, EditOutlined, DeleteOutlined, UnlockOutlined, LockOutlined } from '@ant-design/icons';
import { getCourses, getCourseCategories, createCourse, updateCourse, deleteCourse, toggleCourseStatus } from '../api';
import { Course } from '../types';

const { Option } = Select;

const CourseManage = () => {
  const [loading, setLoading] = useState(false);
  const [courses, setCourses] = useState<Course[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [keyword, setKeyword] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('');
  const [categoryFilter, setCategoryFilter] = useState<string>('');
  const [categories, setCategories] = useState<string[]>([]);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingCourse, setEditingCourse] = useState<Course | null>(null);
  const [form] = Form.useForm();
  const [sortField, setSortField] = useState<string>('');
  const [sortOrder, setSortOrder] = useState<string>('');

  useEffect(() => {
    fetchCourses();
    fetchCategories();
  }, [page, pageSize, keyword, statusFilter, categoryFilter, sortField, sortOrder]);

  const fetchCourses = async () => {
    setLoading(true);
    try {
      const res = await getCourses({
        page, pageSize, keyword, status: statusFilter, category: categoryFilter, sortField, sortOrder,
      });
      if (res.data.code === 0) {
        setCourses(res.data.data.list);
        setTotal(res.data.data.total);
      }
    } catch (error) {
      message.error('获取课程列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchCategories = async () => {
    try {
      const res = await getCourseCategories();
      if (res.data.code === 0) {
        setCategories(res.data.data);
      }
    } catch (error) {
      console.error('获取分类失败', error);
    }
  };

  const handleAdd = () => {
    setEditingCourse(null);
    form.resetFields();
    form.setFieldsValue({ status: 'draft', lesson_count: 0 });
    setModalVisible(true);
  };

  const handleEdit = (record: Course) => {
    setEditingCourse(record);
    form.setFieldsValue(record);
    setModalVisible(true);
  };

  const handleDelete = async (id: number) => {
    try {
      const res = await deleteCourse(id);
      if (res.data.code === 0) {
        message.success('删除成功');
        fetchCourses();
      }
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handleToggleStatus = async (record: Course) => {
    try {
      const res = await toggleCourseStatus(record.id);
      if (res.data.code === 0) {
        message.success(record.status === 'published' ? '已下架' : '已发布');
        fetchCourses();
      }
    } catch (error) {
      message.error('操作失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editingCourse) {
        const res = await updateCourse(editingCourse.id, values);
        if (res.data.code === 0) {
          message.success('更新成功');
          setModalVisible(false);
          fetchCourses();
        }
      } else {
        const res = await createCourse(values);
        if (res.data.code === 0) {
          message.success('创建成功');
          setModalVisible(false);
          fetchCourses();
        }
      }
    } catch (error: any) {
      message.error(error.response?.data?.msg || '操作失败');
    }
  };

  const columns = [
    { title: '课程名称', dataIndex: 'name', key: 'name', width: 200, ellipsis: true },
    { title: '描述', dataIndex: 'description', key: 'description', ellipsis: true },
    { title: '讲师', dataIndex: 'instructor', key: 'instructor', width: 100 },
    { title: '分类', dataIndex: 'category', key: 'category', width: 100 },
    { title: '课时', dataIndex: 'lesson_count', key: 'lesson_count', width: 80, sorter: true },
    { title: '选课人数', dataIndex: 'student_count', key: 'student_count', width: 100, sorter: true },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 80,
      render: (status: string) => (
        <Tag color={status === 'published' ? 'green' : 'default'}>
          {status === 'published' ? '已发布' : '草稿'}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_: any, record: Course) => (
        <Space>
          <Button size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>编辑</Button>
          <Button size="small" icon={record.status === 'published' ? <LockOutlined /> : <UnlockOutlined />} onClick={() => handleToggleStatus(record)}>
            {record.status === 'published' ? '下架' : '发布'}
          </Button>
          <Popconfirm title="确定删除该课程吗？" onConfirm={() => handleDelete(record.id)} okText="确定" cancelText="取消">
            <Button size="small" danger icon={<DeleteOutlined />}>删除</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const handleTableChange = (pagination: any, filters: any, sorter: any) => {
    setPage(pagination.current);
    setPageSize(pagination.pageSize);
    if (sorter.field) {
      setSortField(sorter.field);
      setSortOrder(sorter.order === 'ascend' ? 'ascend' : 'descend');
    } else {
      setSortField('');
      setSortOrder('');
    }
  };

  return (
    <Card>
      <Row gutter={16} className="mb-4">
        <Col span={6}>
          <Input
            placeholder="搜索课程名/讲师"
            prefix={<SearchOutlined />}
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            onPressEnter={fetchCourses}
          />
        </Col>
        <Col span={4}>
          <Select placeholder="全部状态" allowClear value={statusFilter} onChange={setStatusFilter} className="w-full">
            <Option value="published">已发布</Option>
            <Option value="draft">草稿</Option>
          </Select>
        </Col>
        <Col span={4}>
          <Select placeholder="全部分类" allowClear value={categoryFilter} onChange={setCategoryFilter} className="w-full">
            {categories.map(cat => <Option key={cat} value={cat}>{cat}</Option>)}
          </Select>
        </Col>
        <Col span={4}>
          <Button type="primary" onClick={fetchCourses}>搜索</Button>
        </Col>
        <Col span={6} className="text-right">
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>新增课程</Button>
        </Col>
      </Row>

      <Table
        columns={columns}
        dataSource={courses}
        rowKey="id"
        loading={loading}
        pagination={{ current: page, pageSize, total, showSizeChanger: true, pageSizeOptions: [10, 20, 50] }}
        onChange={handleTableChange}
        scroll={{ x: 1000 }}
      />

      <Modal
        title={editingCourse ? '编辑课程' : '新增课程'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="课程名称" rules={[{ required: true, message: '请输入课程名称' }]}>
            <Input placeholder="请输入课程名称" />
          </Form.Item>
          <Form.Item name="description" label="课程描述">
            <Input.TextArea rows={3} placeholder="请输入课程描述" />
          </Form.Item>
          <Form.Item name="instructor" label="讲师">
            <Input placeholder="请输入讲师" />
          </Form.Item>
          <Form.Item name="category" label="分类">
            <Select placeholder="请选择分类" allowClear>
              {categories.map(cat => <Option key={cat} value={cat}>{cat}</Option>)}
            </Select>
          </Form.Item>
          <Form.Item name="lesson_count" label="课时数">
            <Input type="number" min={0} />
          </Form.Item>
          <Form.Item name="status" label="状态">
            <Select>
              <Option value="draft">草稿</Option>
              <Option value="published">已发布</Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  );
};

export default CourseManage;