import React, { useState, useEffect } from 'react';
import { useSelector } from 'react-redux';
import toast from 'react-hot-toast';
import { useNavigate } from 'react-router-dom';
import { transferService } from '../services/transferService';
import { accountService, AccountRegistration, AccountDetails } from '../services/accountService';
import { transactionService, Transaction } from '../services/transactionService';
import { RootState } from '../store/store';
import { selectIsAuthenticated } from '../store/features/auth/authSlice';
import useDebounce from '../hooks/useDebounce';

const TransferFundsPage: React.FC = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [amount, setAmount] = useState('');
  const [recipient, setRecipient] = useState('');
  const [note, setNote] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [accounts, setAccounts] = useState<AccountRegistration[]>([]);
  const [fromAccount, setFromAccount] = useState('');
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [isLoadingTransactions, setIsLoadingTransactions] = useState(false);
  const [recipientDetails, setRecipientDetails] = useState<AccountDetails | null>(null);
  const [isLoadingRecipient, setIsLoadingRecipient] = useState(false);
  const debouncedRecipient = useDebounce(recipient, 1000); // 500ms delay

  const user = useSelector((state: RootState) => state.auth.user);
  const isAuthenticated = useSelector(selectIsAuthenticated);
  const navigate = useNavigate();

  // Fetch accounts when component mounts
  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/auth');
      return;
    }
    fetchAccounts();
  }, [isAuthenticated, navigate]);

  // Fetch transactions when component mounts and when modal is closed
  useEffect(() => {
    if (isAuthenticated && user?.id) {
      fetchTransactions();
    }
  }, [isAuthenticated, user?.id, isModalOpen]);

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

  const fetchTransactions = async () => {
    if (!user?.id) return;
    setIsLoadingTransactions(true);
    try {
      const data = await transactionService.getUserTransactions(user.id);
      setTransactions(data);
    } catch (error: any) {
      toast.error('Failed to fetch transactions');
    } finally {
      setIsLoadingTransactions(false);
    }
  };

  const fetchAccounts = async () => {
    try {
      const data = await accountService.getAll();
      setAccounts(data);
      if (data.length > 0 && !fromAccount) {
        setFromAccount(data[0].account_number);
      }
    } catch (error: any) {
      toast.error('Failed to fetch accounts');
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!amount || !recipient || !fromAccount) {
      toast.error('Amount, From Account, and Recipient are required.');
      return;
    }
    if (!user) {
      toast.error('User not found. Please login again.');
      return;
    }
    setIsSubmitting(true);
    try {
      await transferService.transferFunds({
        user_id: user.id,
        account_id: fromAccount,
        to_account_id: recipient,
        amount: parseFloat(amount),
        description: note,
      });
      toast.success('Transfer successful!');
      setIsModalOpen(false);
      setAmount('');
      setRecipient('');
      setNote('');
      // Refresh transactions after successful transfer
      await fetchTransactions();
    } catch (error: any) {
      toast.error(error?.response?.data?.error || 'Transfer failed');
    } finally {
      setIsSubmitting(false);
    }
  };

  if (!isAuthenticated) {
    return null;
  }

  return (
    <div className="container mx-auto p-4">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-800">Transfer Funds</h1>
        <button
          onClick={() => {
            setIsModalOpen(true);
            fetchAccounts(); // Fetch accounts when opening modal
          }}
          className="bg-blue-600 text-white px-4 py-2 rounded-lg flex items-center gap-2 hover:bg-blue-700 transition-colors"
        >
          Transfer Funds
        </button>
      </div>

      {/* Transaction History Section */}
      <div className="mt-8">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-bold text-gray-800">Transaction History</h2>
          <button
            onClick={fetchTransactions}
            className="text-blue-600 hover:text-blue-800 text-sm font-medium"
          >
            Refresh
          </button>
        </div>
        {isLoadingTransactions ? (
          <div className="text-center">Loading transactions...</div>
        ) : transactions.length === 0 ? (
          <div className="text-center text-gray-500">No transactions found</div>
        ) : (
          <div className="bg-white rounded-lg shadow overflow-hidden">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Amount</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">From Account</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">To Account</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Description</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {transactions.map((transaction) => (
                  <tr key={transaction.id}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {new Date(transaction.transaction_date).toLocaleDateString()}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                        transaction.transaction_type === 'CREDIT' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                      }`}>
                        {transaction.transaction_type}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      ${transaction.amount.toFixed(2)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {transaction.account_id}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {transaction.to_account_id}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {transaction.description}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {isModalOpen && (
        <div className="fixed inset-0 z-50 bg-black bg-opacity-50 flex items-center justify-center px-4 sm:px-6">
        <div className="bg-white rounded-xl shadow-lg p-8 w-full max-w-lg transform transition-all">
          <h2 className="text-2xl font-semibold text-gray-800 mb-6 text-center">Transfer Funds</h2>
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">From Account</label>
              <select
                value={fromAccount}
                onChange={e => setFromAccount(e.target.value)}
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
              <label className="block text-sm font-medium text-gray-700 mb-1">Recipient Account Number</label>
              <input
                type="text"
                value={recipient}
                onChange={e => setRecipient(e.target.value)}
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
              <label className="block text-sm font-medium text-gray-700 mb-1">Note (Optional)</label>
              <input
                type="text"
                value={note}
                onChange={e => setNote(e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 focus:outline-none"
                placeholder="Add a note"
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

export default TransferFundsPage; 