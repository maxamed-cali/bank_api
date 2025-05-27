import axiosInstance from "./axios";

export interface User {
    id: number;
    name: string;
    phone: string;
    role: string;
    status: boolean;
    created: string;
}

export const usersService = {
    getUsers: async (): Promise<User[]> => {
        const { data } = await axiosInstance.get('/api/admin/users');
        return data;
    },

    updateUserStatus: async (userId: number, isActive: boolean): Promise<void> => {
        await axiosInstance.put(`/api/admin/users/${userId}/status`, {
            is_active: isActive
        });
    }
}; 