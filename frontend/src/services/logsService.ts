import axiosInstance from "./axios";

export interface Log {
    ID: number;
    UserID: number;
    FullName: string;
    AccountNumber: string;
    ActionType: string;
    TableName: string;
    RecordID: number;
    Description: string;
    ActionTimestamp: string;
}

interface ApiResponse {
    data: {
        audit_logs: Log[];
    };
}

export const getLogs = async (): Promise<Log[]> => {
    try {
        const response = await axiosInstance.get('/api/admin/audit-logs');
        console.log('API Response:', response);
        
        // Check if the response has the expected structure
        if (response.data && response.data.audit_logs && Array.isArray(response.data.audit_logs)) {
            return response.data.audit_logs;
        }
        
        // If we get here, the response format is unexpected
        console.error('Unexpected response format:', response);
        throw new Error('Unexpected response format from server');
    } catch (error) {
        console.error('Error fetching logs:', error);
        throw error;
    }
}; 