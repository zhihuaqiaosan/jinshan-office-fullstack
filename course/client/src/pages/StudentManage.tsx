import { useEffect, useState } from 'react';
import {
  Table, Button, Input, Select, Space, Modal, Form, message, Popconfirm, Tag, Card, Row, Col, Checkbox,
} from 'antd';
import { PlusOutlined, SearchOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { getStudents, getStudentClasses, createStudent, updateStudent, deleteStudent, getCourses } from '../api';
import { Student, Course } from '../types';

const { Option } = Select;

const StudentManage = () => {
  const [loading, setLoading] = useState(false);
  const [students, setStudents] = useState<Student[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [keyword, setKeyword] = useState('');
  const [classFilter, setClassFilter] = useState<string>('');
  const [statusFilter, setStatusFilter] = useState<string>('');
  const [classes, setClasses] = useState<string[]>([]);
  const [courses, setCourses] = useState<Course[]>([]);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingStudent, setEditingStudent] = useState<Student | null>(null);
  const [form] = Form.useForm();

  useEffect(() => {
    fetchStudents();
    fetchClasses();
    fetchCourses();
  }, [page, pageSize, keyword, classFilter, statusFilter]);

  const fetchStudents = async () => {
    setLoading(true);
    try {
      const res = await getStudents({
        page, pageSize, keyword, className: classFilter, status: statusFilter,
      });
      if (res.data.code === 0) {
        setStudents(res.data.data.list);
        setTotal(res.data.data.total);
      }
    } catch (error) {
      message.error('获取学生列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchClasses = async () => {
    try {
      const res = await getStudentClasses();
      if (res.data.code === 0) {
        setClasses(res.data.data);
      }
    } catch (error) {
      console.error('获取班级失败', error);
    }
  };

  const fetchCourses = async () => {
    try {
      const res = await getCourses({ pageSize: 100 });
      if (res.data.code === 0) {
        setCourses(res.data.data.list);
      }
    } catch (error) {
      console.error('获取课程失败', error);
    }
  };

  const handleAdd = () => {
    setEditingStudent(null);
    form.resetFields();
    form.setFieldsValue({ status: 'active', course_ids: [] });
    setModalVisible(true);
  };

  const handleEdit = (record: Student) => {
    setEditingStudent(record);
    form.setFieldsValue({ ...record, course_ids: record.course_ids });
    setModalVisible(true);
  };

  const handleDelete = async (id: number) => {
    try {
      const res = await deleteStudent(id);
      if (res.data.code === 0) {
        message.success('删除成功');
        fetchStudents();
      }
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editingStudent) {
        const res = await updateStudent(editingStudent.id, values);
        if (res.data.code === 0) {
          message.success('更新成功');
          setModalVisible(false);
          fetchStudents();
        }
      } else {
        const res = await createStudent(values);
        if (res.data.code === 0) {
          message.success('创建成功');
          setModalVisible(false);
          fetchStudents();
        }
      }
    } catch (error: any) {
      message.error(error.response?.data?.msg || '操作失败');
    }
  };

  const columns = [
    { title: '姓名', dataIndex: 'name', key: 'name', width: 100 },
    { title: '学号', dataIndex: 'student_no', key: 'student_no', width: 120 },
    { title: '班级', dataIndex: 'class_name', key: 'class_name', width: 120 },
    { title: '手机号', dataIndex: 'phone', key: 'phone', width: 130 },
    { title: '邮箱', dataIndex: 'email', key: 'email', width: 180, ellipsis: true },
    {
      title: '已选课程',
      dataIndex: 'course_ids',
      key: 'course_ids',
      width: 150,
      render: (ids: number[]) => {
        const selectedCourses = courses.filter(c => ids.includes(c.id));
        return <span>{selectedCourses.map(c => c.name).join('、').slice(0, 30)}</span>;
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 80,
      render: (status: string) => (
        <Tag color={status === 'active' ? 'green' : 'default'}>
          {status === 'active' ? '活跃' : '非活跃'}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 120,
      render: (_: any, record: Student) => (
        <Space>
          <Button size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>编辑</Button>
          <Popconfirm title="确定删除该学生吗？" onConfirm={() => handleDelete(record.id)} okText="确定" cancelText="取消">
            <Button size="small" danger icon={<DeleteOutlined />}>删除</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Card>
      <Row gutter={16} className="mb-4">
        <Col span={6}>
          <Input
            placeholder="搜索姓名/学号"
            prefix={<SearchOutlined />}
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            onPressEnter={fetchStudents}
          />
        </Col>
        <Col span={4}>
          <Select placeholder="全部班级" allowClear value={classFilter} onChange={setClassFilter} className="w-full">
            {classes.map(c => <Option key={c} value={c}>{c}</Option>)}
          </Select>
        </Col>
        <Col span={4}>
          <Select placeholder="全部状态" allowClear value={statusFilter} onChange={setStatusFilter} className="w-full">
            <Option value="active">活跃</Option>
            <Option value="inactive">非活跃</Option>
          </Select>
        </Col>
        <Col span={4}>
          <Button type="primary" onClick={fetchStudents}>搜索</Button>
        </Col>
        <Col span={6} className="text-right">
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>新增学生</Button>
        </Col>
      </Row>

      <Table
        columns={columns}
        dataSource={students}
        rowKey="id"
        loading={loading}
        pagination={{ current: page, pageSize, total, showSizeChanger: true, pageSizeOptions: [10, 20, 50] }}
        onChange={(pagination) => {
          setPage(pagination.current || 1);
          setPageSize(pagination.pageSize || 10);
        }}
        scroll={{ x: 1000 }}
      />

      <Modal
        title={editingStudent ? '编辑学生' : '新增学生'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        width={700}
      >
        <Form form={form} layout="vertical">
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="name" label="姓名" rules={[{ required: true, message: '请输入姓名' }]}>
                <Input placeholder="请输入姓名" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="student_no" label="学号" rules={[{ required: true, message: '请输入学号' }]}>
                <Input placeholder="请输入学号" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="class_name" label="班级">
                <Select placeholder="请选择班级" allowClear>
                  {classes.map(c => <Option key={c} value={c}>{c}</Option>)}
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="status" label="状态">
                <Select>
                  <Option value="active">活跃</Option>
                  <Option value="inactive">非活跃</Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="phone" label="手机号">
                <Input placeholder="请输入手机号" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="email" label="邮箱">
                <Input placeholder="请输入邮箱" />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="course_ids" label="选课">
            <Checkbox.Group>
              <Row gutter={[16, 8]}>
                {courses.map(course => (
                  <Col span={8} key={course.id}>
                    <Checkbox value={course.id}>{course.name}</Checkbox>
                  </Col>
                ))}
              </Row>
            </Checkbox.Group>
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  );
};

export default StudentManage;