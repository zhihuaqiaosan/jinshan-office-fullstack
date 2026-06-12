import axios from 'axios';

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// 认证
export const login = (username: string, password: string) =>
  api.post('/auth/login', { username, password });

export const getCurrentUser = () => api.get('/auth/me');

// 工作台
export const getDashboard = () => api.get('/dashboard');

// 课程管理
export const getCourses = (params: any) => api.get('/courses', { params });
export const getCourseCategories = () => api.get('/courses/categories');
export const getCourseDetail = (id: number) => api.get(`/courses/${id}`);
export const createCourse = (data: any) => api.post('/courses', data);
export const updateCourse = (id: number, data: any) => api.put(`/courses/${id}`, data);
export const deleteCourse = (id: number) => api.delete(`/courses/${id}`);
export const toggleCourseStatus = (id: number) => api.patch(`/courses/${id}/status`);

// 学生管理
export const getStudents = (params: any) => api.get('/students', { params });
export const getStudentClasses = () => api.get('/students/classes');
export const getStudentDetail = (id: number) => api.get(`/students/${id}`);
export const createStudent = (data: any) => api.post('/students', data);
export const updateStudent = (id: number, data: any) => api.put(`/students/${id}`, data);
export const deleteStudent = (id: number) => api.delete(`/students/${id}`);

// 学习总结
export const getSummary = () => api.get('/summary');

export default api;