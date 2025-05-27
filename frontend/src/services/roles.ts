import axiosInstance from "./axios";

export interface User {
    id: number;
    name: string;
    phone: string;
    role: string;
    status: boolean;
    created: string;
}

export const rolesService = {
    getUsers: async (): Promise<User[]> => {
        const { data } = await axiosInstance.get('/api/admin/users');
        return data;
    },

    assignRole: async (userId: number, roles: string[]): Promise<void> => {
        await axiosInstance.post('/api/admin/assign-roles', {
            user_id: userId,
            roles
        });
    }
}; 