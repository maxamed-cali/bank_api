import axiosInstance from "./axios";
import { MoneyRequest } from "./moneyRequestService";

// Common interfaces
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

// Admin specific interfaces
export interface AdminDashboardData extends DashboardData {
    total_users: number;
    active_users: number;
    total_accounts: number;
    system_balance: number;
}

export interface UserStats {
    total_users: number;
    active_users: number;
    new_users_today: number;
    new_users_this_week: number;
}

// User Dashboard Service
export const userDashboardService = {
    getTransactionsSummary: async (): Promise<DashboardData> => {
        const response = await axiosInstance.get("/api/user/dashboard/transactions-summary");
        return response.data.data;
    },

    getRecentRequests: async (userId: number): Promise<MoneyRequest[]> => {
        const response = await axiosInstance.get(`/api/user/transactions/money-request?user_id=${userId}`);
        console.log('API Response:', response.data);
        // If the response is an array, return it directly
        if (Array.isArray(response.data)) {
            return response.data;
        }
        // If the response has a data property that's an array, return that
        if (response.data && Array.isArray(response.data.data)) {
            return response.data.data;
        }
        // If neither, return empty array
        return [];
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

// Admin Dashboard Service
export const adminDashboardService = {
    getAdminDashboardSummary: async (): Promise<AdminDashboardData> => {
        const response = await axiosInstance.get("/api/admin/admindashboard/transactions-summary");
        return response.data.data;
    },

    getRecentRequests: async (): Promise<MoneyRequest[]> => {
        const response = await axiosInstance.get("/api/admin/dashboard/recent-requests");
        return response.data.data || [];
    },

    getMonthlyUserGrowth: async (): Promise<ChartData[]> => {
        const response = await axiosInstance.get("/api/admin/admindashboard/monthly-transactions");
        return Array.isArray(response.data) ? response.data : [];
    },

    getSystemTransactions: async (): Promise<Transaction[]> => {
        const response = await axiosInstance.get("/api/admin/transactions/history");
        return response.data.data || [];
    },

    // getPendingRequests: async (): Promise<any[]> => {
    //     const response = await axiosInstance.get("/api/admin/admindashboard/pending-requests");
    //     return response.data.data || [];
    // }
}; 