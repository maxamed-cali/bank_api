import axiosInstance from './axios';

export interface MoneyRequest {
  ID: number;
  user_id: number;
  requester_id: string;
  recipient_id: string;
  Amount: number;
  Status: 'PENDING' | 'ACCEPTED' | 'REJECTED';
  ExpiresAt: string;
  RequesteAt: string;
  description?: string;
}

export interface CreateMoneyRequestPayload {
  user_id: number;
  requester_id: string;
  recipient_id: string;
  amount: number;
  description: string;
}

export const moneyRequestService = {
  // Create a new money request
  async createRequest(payload: CreateMoneyRequestPayload) {
    try {
      const response = await axiosInstance.post('/api/user/money-request', payload);
      return response.data;
    } catch (error: any) {
      console.error('Error creating money request:', error.response?.data || error.message);
      throw error;
    }
  },

  // Get all money requests for a user
  async getUserRequests(userId: number) {
    try {
      const response = await axiosInstance.get(`/api/user/transactions/money-request?user_id=${userId}`);
      return response.data;
    } catch (error: any) {
      console.error('Error fetching money requests:', error.response?.data || error.message);
      throw error;
    }
  },

  // Accept a money request
  async acceptRequest(requestId: number) {
    try {
      const response = await axiosInstance.put(`/api/user/accept-money-request/${requestId}`);
      return response.data;
    } catch (error: any) {
      console.error('Error accepting money request:', error.response?.data || error.message);
      throw error;
    }
  },

  // Reject a money request
  async rejectRequest(requestId: number) {
    try {
      const response = await axiosInstance.put(`/api/user/decline-money-request/${requestId}`);
      return response.data;
    } catch (error: any) {
      console.error('Error rejecting money request:', error.response?.data || error.message);
      throw error;
    }
  }
}; 