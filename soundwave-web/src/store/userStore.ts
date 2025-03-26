import { create } from 'zustand';
import { message } from 'antd';
import axios from '../utils/axios';

interface User {
  id: string;
  username: string;
  role: string;
  permissions: string[];
}

interface UserState {
  user: User | null;
  token: string | null;
  setUser: (user: User | null) => void;
  setToken: (token: string | null) => void;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
}

export const useUserStore = create<UserState>((set) => ({
  user: null,
  token: localStorage.getItem('token'),
  
  setUser: (user) => set({ user }),
  setToken: (token) => set({ token }),
  
  login: async (username: string, password: string) => {
    try {
      const response = await axios.post('/auth/login', { username, password });
      const { token, user } = response.data;
      
      localStorage.setItem('token', token);
      localStorage.setItem('user', JSON.stringify(user));
      
      set({ user, token });
      message.success('登录成功');
    } catch (error) {
      message.error('登录失败');
      throw error;
    }
  },
  
  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    set({ user: null, token: null });
    message.success('已退出登录');
  },
})); 