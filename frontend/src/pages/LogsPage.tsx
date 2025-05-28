import React, { useEffect, useState } from 'react';
import { getLogs, Log } from '../services/logsService';

interface LogsResponse {
    data?: Log[];
    logs?: Log[];
    [key: string]: any;
}

const LogsPage: React.FC = () => {
    const [logs, setLogs] = useState<Log[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchLogs = async () => {
            try {
                const data = await getLogs();

                // Ensure data is an array before setting it
                if (Array.isArray(data)) {
                    setLogs(data);
                } else if (data && typeof data === 'object') {
                    // If the response is wrapped in an object, try to find the array
                    const response = data as LogsResponse;
                    
                    // Try to find the logs array in the response
                    let logsArray;
                    if (response.data) {
                        logsArray = response.data;
                    } else if (response.logs) {
                        logsArray = response.logs;
                    } else {
                        // Try to find any array in the response
                        const possibleArrays = Object.values(response).filter(Array.isArray);
                        logsArray = possibleArrays[0];
                    }

                    if (Array.isArray(logsArray)) {
                        setLogs(logsArray);
                    } else {
                        throw new Error('Invalid response format: logs data is not an array');
                    }
                } else {
                    throw new Error('Invalid response format: expected an array of logs');
                }
                setError(null);
            } catch (err) {
                console.error('Error fetching logs:', err);
                setError('Failed to fetch logs');
            } finally {
                setLoading(false);
            }
        };

        fetchLogs();
    }, []);

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleString();
    };

    if (loading) {
        return (
            <div className="flex justify-center items-center min-h-[400px]">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="flex justify-center items-center min-h-[400px]">
                <p className="text-red-500">{error}</p>
            </div>
        );
    }

    if (!Array.isArray(logs) || logs.length === 0) {
        return (
            <div className="p-6">
                <h1 className="text-2xl font-bold mb-6">System Logs</h1>
                <div className="text-center text-gray-500">No logs available</div>
            </div>
        );
    }

    return (
        <div className="p-6">
            <h1 className="text-2xl font-bold mb-6">System Logs</h1>
            <div className="overflow-x-auto rounded-lg shadow">
                <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                        <tr>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">User</th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Account</th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Action</th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Table</th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Description</th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Timestamp</th>
                        </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                        {logs.map((log, index) => (
                            <tr key={index} className="hover:bg-gray-50">
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{log.FullName}</td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{log.AccountNumber}</td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{log.ActionType}</td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{log.TableName}</td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{log.Description}</td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{formatDate(log.ActionTimestamp)}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default LogsPage; 