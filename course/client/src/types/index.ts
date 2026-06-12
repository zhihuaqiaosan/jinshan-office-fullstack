export interface User {
  id: number;
  username: string;
  name: string;
  role: string;
  avatar: string;
  created_at: string;
}

export interface Course {
  id: number;
  name: string;
  description: string;
  instructor: string;
  cover: string;
  category: string;
  status: 'published' | 'draft';
  student_count: number;
  lesson_count: number;
  created_at: string;
  updated_at: string;
}

export interface Student {
  id: number;
  name: string;
  student_no: string;
  class_name: string;
  phone: string;
  email: string;
  status: 'active' | 'inactive';
  course_ids: number[];
  created_at: string;
  updated_at: string;
}