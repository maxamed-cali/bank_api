import React, { useState, useEffect } from 'react';
import { accountService, AccountRegistration, AccountType } from '../services/accountService';
import toast from 'react-hot-toast';
import { FiEdit2, FiTrash2, FiPlus } from 'react-icons/fi';
import { useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';

const AccountRegistrationPage: React.FC = () => {
  const [accounts, setAccounts] = useState<AccountRegistration[]>([]);
  const [accountTypes, setAccountTypes] = useState<AccountType[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingAccount, setEditingAccount] = useState<AccountRegistration | null>(null);
  const [formData, setFormData] = useState({
    account_number: '',
    balance: 0,
    account_type_id: undefined as number | undefined,
  });

  const navigate = useNavigate();
  // Get auth state from Redux store
  const { user, isAuthenticated } = useSelector((state: any) => state.auth);

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/auth');
      return;
    }
    fetchAccounts();
    fetchAccountTypes();
  }, [isAuthenticated, navigate]);

  const fetchAccounts = async () => {
    try {
      setIsLoading(true);
      const data = await accountService.getAll();
      setAccounts(data);
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || 'Failed to fetch accounts';
      toast.error(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  const fetchAccountTypes = async () => {
    try {
      const data = await accountService.getAccountTypes();
      setAccountTypes(data);
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || 'Failed to fetch account types';
      toast.error(errorMessage);
    }
  };

  const validateForm = () => {
    if (!formData.account_number) {
      toast.error('Account number is required');
      return false;
    }
    if (formData.balance < 0) {
      toast.error('Balance cannot be negative');
      return false;
    }
    return true;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    if (!isAuthenticated || !user?.id) {
      toast.error('Please login to continue');
      navigate('/auth');
      return;
    }

    try {
      const accountData = {
        ...formData,
        user_id: user.id,
      };

      console.log('Submitting account data:', accountData);

      if (editingAccount) {
        await accountService.update(editingAccount.id!, accountData);
        toast.success('Account updated successfully');
      } else {
        await accountService.create(accountData);
        toast.success('Account created successfully');
      }
      setIsModalOpen(false);
      setEditingAccount(null);
      setFormData({
        account_number: '',
        balance: 0,
        account_type_id: undefined,
      });
      fetchAccounts();
    } catch (error: any) {
      const errorMessage = error?.response?.data?.error || (editingAccount ? 'Failed to update account' : 'Failed to create account');
      toast.error(errorMessage);
      console.error('Error submitting account:', error);
    }
  };

  const handleEdit = (account: AccountRegistration) => {
    setEditingAccount(account);
    setFormData({
      account_number: account.account_number,
      balance: account.balance,
      account_type_id: account.account_type_id,
    });
    setIsModalOpen(true);
  };

  const handleDelete = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this account?')) {
      try {
        await accountService.delete(id);
        toast.success('Account deleted successfully');
        fetchAccounts();
      } catch (error: any) {
        const errorMessage = error.response?.data?.message || 'Failed to delete account';
        toast.error(errorMessage);
      }
    }
  };

  if (!isAuthenticated) {
    return null; // Don't render anything if not authenticated
  }

  return (
    <div className="container mx-auto p-4">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-800">Account Registration</h1>
        <button
          onClick={() => {
            setEditingAccount(null);
            setFormData({
              account_number: '',
              balance: 0,
              account_type_id: undefined,
            });
            setIsModalOpen(true);
          }}
          className="bg-blue-600 text-white px-4 py-2 rounded-lg flex items-center gap-2 hover:bg-blue-700 transition-colors"
        >
          <FiPlus /> Add New Account
        </button>
      </div>

      {isLoading ? (
        <div className="text-center">Loading...</div>
      ) : (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Account Number</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Balance</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Currency</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Created</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {(accounts || []).map((account) => (
                <tr key={account.id}>
                  <td className="px-6 py-4 whitespace-nowrap">{account.name || 'N/A'}</td>
                  <td className="px-6 py-4 whitespace-nowrap">{account.account_number}</td>
                  <td className="px-6 py-4 whitespace-nowrap">{account.type_name || 'N/A'}</td>
                  <td className="px-6 py-4 whitespace-nowrap">${account.balance.toFixed(2)}</td>  
                  <td className="px-6 py-4 whitespace-nowrap">{account.currency || 'N/A'}</td>
                  <td className="px-6 py-4 whitespace-nowrap">{account.created_at ? new Date(account.created_at).toLocaleDateString() : 'N/A'}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <button
                      onClick={() => handleEdit(account)}
                      className="text-indigo-600 hover:text-indigo-900 mr-4"
                    >
                      <FiEdit2 className="inline" />
                    </button>
                    <button
                      onClick={() => handleDelete(account.id!)}
                      className="text-red-600 hover:text-red-900"
                    >
                      <FiTrash2 className="inline" />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {isModalOpen && (
        <div className="fixed inset-0 z-50 bg-black bg-opacity-50 flex items-center justify-center px-4 sm:px-6 lg:px-8">
        <div className="bg-white rounded-xl shadow-xl transform transition-all p-8 w-full max-w-lg">
          <h2 className="text-2xl font-semibold text-gray-800 mb-6 text-center">
            {editingAccount ? 'Edit Account' : 'Create New Account'}
          </h2>
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Account Number</label>
              <input
                type="text"
                value={formData.account_number}
                onChange={(e) => setFormData({ ...formData, account_number: e.target.value })}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 focus:outline-none"
                required
                placeholder="Enter account number"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Initial Balance</label>
              <input
                type="number"
                step="0.01"
                min="0"
                value={formData.balance}
                onChange={(e) => setFormData({ ...formData, balance: parseFloat(e.target.value) || 0 })}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 focus:outline-none"
                required
                placeholder="Enter initial balance"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Account Type</label>
              <select
                value={formData.account_type_id || ''}
                onChange={(e) => setFormData({ ...formData, account_type_id: e.target.value ? parseInt(e.target.value) : undefined })}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 focus:outline-none"
              >
                <option value="">Select an account type</option>
                {accountTypes.map((type) => (
                  <option key={type.ID} value={type.ID}>
                    {type.type_name} ({type.currency})
                  </option>
                ))}
              </select>
            </div>
            <div className="flex justify-end gap-4 pt-4 border-t border-gray-200">
              <button
                type="button"
                onClick={() => setIsModalOpen(false)}
                className="px-4 py-2 text-sm font-medium text-gray-600 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                className="px-5 py-2 text-sm font-semibold text-white bg-blue-600 rounded-lg hover:bg-blue-700 transition-colors"
              >
                {editingAccount ? 'Update' : 'Create'}
              </button>
            </div>
          </form>
        </div>
      </div>
      
      )}
    </div>
  );
};

export default AccountRegistrationPage;