import React, { useState, useEffect } from 'react';
import { useSelector } from 'react-redux';
import toast from 'react-hot-toast';
import { useNavigate } from 'react-router-dom';
import { moneyRequestService, MoneyRequest } from '../services/moneyRequestService';
import { accountService, AccountRegistration, AccountDetails } from '../services/accountService';
import { RootState } from '../store/store';
import { selectIsAuthenticated } from '../store/features/auth/authSlice';
import { Menu } from '@headlessui/react';
import { EllipsisVerticalIcon } from '@heroicons/react/24/solid';
import useDebounce from '../hooks/useDebounce';

const MoneyRequestPage: React.FC = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [amount, setAmount] = useState('');
  const [requester_id, setRequesterId] = useState('');
  const [recipient_id, setRecipientId] = useState('');
  const [description, setDescription] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [accounts, setAccounts] = useState<AccountRegistration[]>([]);
  const [requests, setRequests] = useState<MoneyRequest[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [recipientDetails, setRecipientDetails] = useState<AccountDetails | null>(null);
  const [isLoadingRecipient, setIsLoadingRecipient] = useState(false);
  const debouncedRecipient = useDebounce(recipient_id, 1000);

  const user = useSelector((state: RootState) => state.auth.user);
  const isAuthenticated = useSelector(selectIsAuthenticated);
  const navigate = useNavigate();

  // Fetch accounts and requests when component mounts
  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/auth');
      return;
    }
    fetchAccounts();
    fetchRequests();
  }, [isAuthenticated, navigate]);

  // Add effect to fetch recipient details when debounced value changes
  useEffect(() => {
    const fetchRecipientDetails = async () => {
      if (!debouncedRecipient || debouncedRecipient.length < 6) {
        setRecipientDetails(null);
        return;
      }
      
      setIsLoadingRecipient(true);
      try {
        const details = await accountService.getAccountDetails(debouncedRecipient);
        setRecipientDetails(details);
      } catch (error: any) {
        toast.error('Account not found');
        setRecipientDetails(null);
      } finally {
        setIsLoadingRecipient(false);
      }
    };

    fetchRecipientDetails();
  }, [debouncedRecipient]);

  const fetchAccounts = async () => {
    try {
      const data = await accountService.getAll();
      setAccounts(data);
      if (data.length > 0 && !requester_id) {
        setRequesterId(data[0].account_number);
      }
    } catch (error: any) {
      toast.error('Failed to fetch accounts');
    }
  };

  const fetchRequests = async () => {
    if (!user?.id) return;
    setIsLoading(true);
    try {
      const data = await moneyRequestService.getUserRequests(user.id);
      setRequests(data);
    } catch (error: any) {
      toast.error('Failed to fetch money requests');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!amount || !requester_id || !recipient_id) {
      toast.error('Amount, From Account, and To Account are required.');
      return;
    }
    if (!user) {
      toast.error('User not found. Please login again.');
      return;
    }
    setIsSubmitting(true);
    try {
      await moneyRequestService.createRequest({
        user_id: user.id,
        requester_id: requester_id,
        recipient_id: recipient_id,
        amount: parseFloat(amount),
        description,
      });
      toast.success('Money request created successfully!');
      setIsModalOpen(false);
      setAmount('');
      setRecipientId('');
      setDescription('');
      fetchRequests();
    } catch (error: any) {
      toast.error(error?.response?.data?.error || 'Failed to create money request');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleAccept = async (requestId: number) => {
    try {
      await moneyRequestService.acceptRequest(requestId);
      toast.success('Money request accepted!');
      fetchRequests();
    } catch (error: any) {
      toast.error(error?.response?.data?.error || 'Failed to accept request');
    }
  };

  const handleReject = async (requestId: number) => {
    try {
      await moneyRequestService.rejectRequest(requestId);
      toast.success('Money request rejected!');
      fetchRequests();
    } catch (error: any) {
      toast.error(error?.response?.data?.error || 'Failed to reject request');
    }
  };

  if (!isAuthenticated) {
    return null;
  }

  return (
    <div className="container mx-auto p-4">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-800">Money Requests</h1>
        <button
          onClick={() => {
            setIsModalOpen(true);
            fetchAccounts();
          }}
          className="bg-blue-600 text-white px-4 py-2 rounded-lg flex items-center gap-2 hover:bg-blue-700 transition-colors"
        >
          New Money Request
        </button>
      </div>

      {/* Money Requests List */}
      <div className="mt-8">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-bold text-gray-800">Request History</h2>
          <button
            onClick={fetchRequests}
            className="text-blue-600 hover:text-blue-800 text-sm font-medium"
          >
            Refresh
          </button>
        </div>
        {isLoading ? (
          <div className="text-center">Loading requests...</div>
        ) : requests.length === 0 ? (
          <div className="text-center text-gray-500">No money requests found</div>
        ) : (
         <div className="bg-white rounded-lg shadow-md overflow-hidden">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-100">
          <tr>
            {['Date', 'Amount', 'From Account', 'To Account', 'Status', 'Actions'].map((header) => (
              <th
                key={header}
                className="px-6 py-3 text-left text-sm font-semibold text-gray-700 uppercase tracking-wide"
              >
                {header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {requests.map((request) => (
            <tr key={request.ID} className="hover:bg-gray-50">
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                {new Date(request.RequesteAt).toLocaleDateString()}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                ${request.Amount ? request.Amount.toFixed(2) : '0.00'}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">{request.requester_id}</td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">{request.recipient_id}</td>
              <td className="px-6 py-4 whitespace-nowrap">
                <span
                  className={`px-3 py-1 inline-flex text-xs font-medium rounded-full ${
                    request.Status === 'ACCEPTED'
                      ? 'bg-green-100 text-green-800'
                      : request.Status === 'REJECTED'
                      ? 'bg-red-100 text-red-800'
                      : 'bg-yellow-100 text-yellow-800'
                  }`}
                >
                  {request.Status}
                </span>
              </td>
              <td className="px-6 py-4 text-sm text-center font-medium z-[99999999]">
                {request.Status === 'PENDING' && (
                  <Menu as="div" className="relative inline-block text-center">
                    <Menu.Button className="inline-flex justify-center w-full text-sm font-medium text-gray-700 hover:text-gray-900">
                      <EllipsisVerticalIcon className="h-5 w-5" aria-hidden="true" />
                    </Menu.Button>
                    <Menu.Items className="absolute right-0 z-10 mt-2 w-28 origin-top-right bg-white border border-gray-200 rounded-md shadow-lg focus:outline-none">
                      <div className="py-1">
                        <Menu.Item>
                          {({ active }) => (
                            <button
                              onClick={() => handleAccept(request.ID)}
                              className={`${
                                active ? 'bg-gray-100' : ''
                              } w-full text-left px-4 py-2 text-sm text-green-700`}
                            >
                              Accept
                            </button>
                          )}
                        </Menu.Item>
                        <Menu.Item>
                          {({ active }) => (
                            <button
                              onClick={() => handleReject(request.ID)}
                              className={`${
                                active ? 'bg-gray-100' : ''
                              } w-full text-left px-4 py-2 text-sm text-red-700`}
                            >
                              Decline
                            </button>
                          )}
                        </Menu.Item>
                      </div>
                    </Menu.Items>
                  </Menu>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
        )}
      </div>

      {/* New Money Request Modal */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 bg-black bg-opacity-50 flex items-center justify-center px-4 sm:px-6">
        <div className="bg-white rounded-xl shadow-lg p-8 w-full max-w-lg transform transition-all">
          <h2 className="text-2xl font-semibold text-gray-800 mb-6 text-center">New Money Request</h2>
          <form onSubmit={handleSubmit} className="space-y-6">
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">From Account</label>
              <select
                value={requester_id}
                onChange={e => setRequesterId(e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 focus:outline-none"
                required
              >
                {accounts.length === 0 && <option value="">No accounts found</option>}
                {accounts.map(acc => (
                  <option key={acc.account_number} value={acc.account_number}>
                    {acc.account_number} {acc.type_name ? `(${acc.type_name})` : ''}
                  </option>
                ))}
              </select>
            </div>
      
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">To Account Number</label>
              <input
                type="text"
                value={recipient_id}
                onChange={e => setRecipientId(e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 focus:outline-none"
                required
                placeholder="Enter recipient account number"
              />
              {isLoadingRecipient && (
                <div className="mt-2 text-sm text-gray-500 flex items-center">
                  <svg className="animate-spin -ml-1 mr-3 h-4 w-4 text-blue-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Searching for account...
                </div>
              )}
              {recipientDetails && !isLoadingRecipient && (
                <div className="mt-2 p-3 bg-gray-50 rounded-lg border border-gray-200">
                  <p className="text-sm text-gray-700">
                    <span className="font-medium">Account Holder:</span> {recipientDetails.full_name}
                  </p>
                  <p className="text-sm text-gray-700">
                    <span className="font-medium">Account Number:</span> {recipientDetails.account_number}
                  </p>
                </div>
              )}
            </div>
      
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Amount</label>
              <input
                type="number"
                min="0.01"
                step="0.01"
                value={amount}
                onChange={e => setAmount(e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 focus:outline-none"
                required
                placeholder="Enter amount"
              />
            </div>
      
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Description (Optional)</label>
              <input
                type="text"
                value={description}
                onChange={e => setDescription(e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 focus:outline-none"
                placeholder="Add a description"
              />
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
                className="px-5 py-2 text-sm font-semibold text-white bg-blue-600 rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-60"
                disabled={isSubmitting}
              >
                {isSubmitting ? 'Submitting...' : 'Submit'}
              </button>
            </div>
          </form>
        </div>
      </div>
      
      )}
    </div>
  );
};

export default MoneyRequestPage; 