import axiosInstance from "./axios";

export interface DashboardData {
    total_transactions: number;
    pending_requests: number;
    total_transfers: number;
    total_sent_amount: number;
    total_received_amount: number;
}

export interface ChartData {
    name: string;
    total: number;
}

export interface Transaction {
    transaction_date: string;
    transaction_type: 'DEBIT' | 'CREDIT';
    amount: number;
    account_id: string;
    to_account_id: string;
    description: string;
}

export const dashboardService = {
    getTransactionsSummary: async (): Promise<DashboardData> => {
        const response = await axiosInstance.get("/api/user/dashboard/transactions-summary");
        return response.data.data;
    },

    getMonthlyTransactions: async (): Promise<ChartData[]> => {
        const response = await axiosInstance.get("/api/user/dashboard/monthly-transactions");
        return Array.isArray(response.data) ? response.data : [];
    },

    getTransactionHistory: async (userId: string | number): Promise<Transaction[]> => {
        const response = await axiosInstance.get(`/api/user/transactions/history?user_id=${userId}`);
        return response.data.data || [];
    }
}; 