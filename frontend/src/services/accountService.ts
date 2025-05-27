import axiosInstance from './axios';

export interface AccountRegistration {
  id?: string;
  user_id: number;
  account_number: string;
  balance: number;
  account_type_id?: number;
  type_name?: string;
  Description?: string;
  currency?: string;
  is_active?: boolean;
  created_at?: string;
  updatedAt?: string;
  name?: string;
}

export interface AccountType {
  ID: number;
  type_name: string;
  Description: string;
  currency: string;
}

export interface AccountDetails {
  account_number: string;
  full_name: string;
}

export const accountService = {
  // Create new account registration
  async create(account: Omit<AccountRegistration, 'id' | 'status' | 'createdAt' | 'updatedAt'>) {
    try {
      console.log('Creating account with data:', account);
      const response = await axiosInstance.post('/api/user/accounts', account);
      console.log('Create account response:', response.data);
      return response.data;
    } catch (error: any) {
      console.error('Error creating account:', error.response?.data || error.message);
      throw error;
    }
  },

  // Get all account registrations
  async getAll(): Promise<AccountRegistration[]> {
    const user = JSON.parse(localStorage.getItem('user') || '{}');
    const response = await axiosInstance.get(`/api/user/accounts/${user.id}`);
    return response.data;
  },

  // Get single account registration
  async getById(id: string) {
    try {
      const response = await axiosInstance.get(`/api/user/accounts/${id}`);
      return response.data;
    } catch (error: any) {
      console.error('Error fetching account:', error.response?.data || error.message);
      throw error;
    }
  },

  // Update account registration
  async update(id: string, account: Partial<AccountRegistration>) {
    try {
      const response = await axiosInstance.put(`/api/user/accounts/${id}`, account);
      return response.data;
    } catch (error: any) {
      console.error('Error updating account:', error.response?.data || error.message);
      throw error;
    }
  },

  // Delete account registration
  async delete(id: string) {
    try {
      const response = await axiosInstance.delete(`/api/user/accounts/${id}`);
      return response.data;
    } catch (error: any) {
      console.error('Error deleting account:', error.response?.data || error.message);
      throw error;
    }
  },

  // Get all account types
  async getAccountTypes() {
    try {
      const response = await axiosInstance.get('/api/user/account-types');
      return response.data;
    } catch (error: any) {
      console.error('Error fetching account types:', error.response?.data || error.message);
      throw error;
    }
  },

  getAccountDetails: async (accountNumber: string): Promise<AccountDetails> => {
   try {
    const response = await axiosInstance.get(`/api/user/account-details?account_number=${accountNumber}`);
    return response.data;
   } catch (error: any) {
    console.error('Error fetching account details:', error.response?.data || error.message);
    throw error;
   }
  },
}; 