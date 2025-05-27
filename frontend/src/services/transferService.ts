import axiosInstance from './axios';

export interface TransferFundsPayload {
  user_id: number;
  account_id: string;
  to_account_id: string;
  amount: number;
  description?: string;
}

export const transferService = {
  async transferFunds(payload: TransferFundsPayload) {
    const response = await axiosInstance.post('/api/user/money-transer', payload);
    return response.data;
  }
}; 