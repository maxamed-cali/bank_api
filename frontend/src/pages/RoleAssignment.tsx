import React, { useState, useEffect } from 'react';
import { useSelector } from 'react-redux';
import toast from 'react-hot-toast';
import { useNavigate } from 'react-router-dom';
import { FiUser, FiCheck } from 'react-icons/fi';
import { rolesService, User } from '../services/roles';

const AVAILABLE_ROLES = ["Admin", "User"]; // Add or modify roles as needed

const RoleAssignment: React.FC = () => {
    const [users, setUsers] = useState<User[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [selectedRoles, setSelectedRoles] = useState<{ [key: number]: string[] }>({});

    const { user, isAuthenticated } = useSelector((state: any) => state.auth);
    const navigate = useNavigate();

    useEffect(() => {
        if (!isAuthenticated) {
            navigate('/auth');
            return;
        }
        if (user?.role !== 'Admin') {
            toast.error('Unauthorized access');
            navigate('/');
            return;
        }
        fetchUsers();
    }, [isAuthenticated, navigate, user]);

    const fetchUsers = async () => {
        try {
            setIsLoading(true);
            const data = await rolesService.getUsers();
            setUsers(data);
            // Initialize selected roles
            const initialRoles: { [key: number]: string[] } = {};
            data.forEach((user: User) => {
                initialRoles[user.id] = [user.role];
            });
            setSelectedRoles(initialRoles);
        } catch (error: any) {
            toast.error(error.response?.data?.message || 'Failed to fetch users');
        } finally {
            setIsLoading(false);
        }
    };

    const handleRoleChange = (userId: number, role: string) => {
        setSelectedRoles(prev => ({
            ...prev,
            [userId]: [role] // Currently allowing only one role per user
        }));
    };

    const handleAssignRole = async (userId: number) => {
        try {
            await rolesService.assignRole(userId, selectedRoles[userId]);
            toast.success('Role assigned successfully');
            fetchUsers(); // Refresh the list
        } catch (error: any) {
            toast.error(error.response?.data?.message || 'Failed to assign role');
        }
    };

    return (
        <div className="container mx-auto p-4">
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-2xl font-bold text-gray-800">Role Assignment</h1>
            </div>

            {isLoading ? (
                <div className="text-center">Loading...</div>
            ) : (
                <div className="bg-white rounded-lg shadow overflow-hidden">
                    <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                            <tr>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Phone</th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Current Role</th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Assign Role</th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                            {users.map((user) => (
                                <tr key={user.id}>
                                    <td className="px-6 py-4 whitespace-nowrap">
                                        <div className="flex items-center">
                                            <div className="flex-shrink-0 h-10 w-10">
                                                <div className="h-10 w-10 rounded-full bg-gray-200 flex items-center justify-center">
                                                    <FiUser className="h-6 w-6 text-gray-500" />
                                                </div>
                                            </div>
                                            <div className="ml-4">
                                                <div className="text-sm font-medium text-gray-900">{user.name}</div>
                                            </div>
                                        </div>
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{user.phone}</td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{user.role}</td>
                                    <td className="px-6 py-4 whitespace-nowrap">
                                        <select
                                            value={selectedRoles[user.id]?.[0] || ''}
                                            onChange={(e) => handleRoleChange(user.id, e.target.value)}
                                            className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md"
                                        >
                                            <option value="">Select a role</option>
                                            {AVAILABLE_ROLES.map((role) => (
                                                <option key={role} value={role}>
                                                    {role}
                                                </option>
                                            ))}
                                        </select>
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                                        <button
                                            onClick={() => handleAssignRole(user.id)}
                                            disabled={!selectedRoles[user.id]?.length}
                                            className="inline-flex items-center px-3 py-1 rounded-md text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
                                        >
                                            <FiCheck className="mr-1" /> Assign
                                        </button>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            )}
        </div>
    );
};

export default RoleAssignment; 