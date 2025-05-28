import axiosInstance from './axios';

export interface Transaction {
  id: number;
  user_id: number;
  account_id: string;
  to_account_id: string;
  amount: number;
  description: string;
  transaction_date: string;
  transaction_type: 'DEBIT' | 'CREDIT';
}

export interface TransactionResponse {
  data: Transaction[];
}

export interface TransactionFilters {
  account_id?: string;
  transaction_type?: 'DEBIT' | 'CREDIT';
}

function getCurrentUserRole(): string {
  const userStr = localStorage.getItem('user');
  if (!userStr) return 'User'; // default fallback
  try {
    const user = JSON.parse(userStr);
    return user.role || 'User';
  } catch {
    return 'User';
  }
}

export const transactionService = {
  async getUserTransactions(userId: number) {
    try {
      const response = await axiosInstance.get<TransactionResponse>(`/api/user/transactions/history?user_id=${userId}`);
      return response.data.data;
    } catch (error: any) {
      console.error('Error fetching transactions:', error.response?.data || error.message);
      throw error;
    }
  },

  async getAccountTransactions(filters: TransactionFilters) {
    try {
      const queryParams = new URLSearchParams();
      
      if (filters.account_id) {
        queryParams.append('account_id', filters.account_id);
      }
      
      if (filters.transaction_type) {
        queryParams.append('transaction_type', filters.transaction_type);
      }

      const data = await transactionService.getTransactionHistory(queryParams);
      return data;
    } catch (error: any) {
      console.error('Error fetching account transactions:', error.response?.data || error.message);
      throw error;
    }
  },

  getTransactionHistory: async (queryParams?: URLSearchParams) => {
    const role = getCurrentUserRole();
    let endpoint = "";

    if (role === "Admin") {
      endpoint = `/api/admin/transactions/history?${queryParams?.toString() || ""}`;
    } else {
      endpoint = `/api/user/transactions/history?${queryParams?.toString() || ""}`;
    }

    const response = await axiosInstance.get(endpoint);
    return response.data.data || [];
  }
}; 