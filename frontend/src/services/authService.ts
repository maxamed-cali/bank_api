import axiosInstance from './axios';

interface LoginCredentials {
  email: string;
  password: string;
}

interface RegisterCredentials {
  password: string;
  full_name: string;
  email: string;
  phone_number: string;
}

// Mock API calls
export const authService = {
  async login({ email, password }: LoginCredentials) {
    const response = await axiosInstance.post('/login', { email, password });
   
    return response.data;
  },

  async register(credentials: RegisterCredentials) {
    const response = await axiosInstance.post('/register', credentials);
   
    return response.data;
  },

  async logout() {
    await axiosInstance.post('/auth/logout');
    localStorage.removeItem('token');
    localStorage.removeItem('user')
  },

  async changePassword(old_password: string, new_password: string): Promise<void> {
    // Call the real backend API for password reset
    await axiosInstance.post('/api/user/password-reset', {
      old_password,
      new_password,
    });
  }
}; 