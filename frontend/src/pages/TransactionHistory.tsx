import React, { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import { transactionService, Transaction } from '../services/transactionService';
// Assuming you will use react-router-dom to get walletId from URL
// import { useParams } from 'react-router-dom';

// Assuming lucide-react for icons, adjust if using different library
// import { Filter, CalendarDays, Tag } from 'lucide-react';

const TransactionHistoryPage: React.FC = () => {
    const [searchParams] = useSearchParams();
    const accountId = searchParams.get('account_id');
    
    const [transactions, setTransactions] = useState<Transaction[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [transactionType, setTransactionType] = useState<'DEBIT' | 'CREDIT' | ''>('');

    useEffect(() => {
        const fetchTransactions = async () => {
            if (!accountId) return;
            
            try {
                setLoading(true);
                const data = await transactionService.getAccountTransactions({
                    account_id: accountId,
                    transaction_type: transactionType || undefined
                });
                setTransactions(data);
                setError(null);
            } catch (err) {
                setError('Failed to fetch transactions');
                console.error(err);
            } finally {
                setLoading(false);
            }
        };

        fetchTransactions();
    }, [accountId, transactionType]);

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    };

    const formatAmount = (amount: number) => {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD'
        }).format(amount);
    };

    if (loading) {
        return <div className="flex justify-center items-center h-64">Loading...</div>;
    }

    if (error) {
        return <div className="text-red-500 text-center">{error}</div>;
    }

    return (
        <div className="container mx-auto px-4 py-8">
            <div className="bg-white rounded-lg shadow-lg p-6">
                <div className="flex justify-between items-center mb-6">
                    <h1 className="text-2xl font-bold text-gray-800">Transaction History</h1>
                    <div className="flex gap-4">
                        <select
                            className="p-2 border rounded-md bg-white"
                            value={transactionType}
                            onChange={(e) => setTransactionType(e.target.value as 'DEBIT' | 'CREDIT' | '')}
                        >
                            <option value="">All Transactions</option>
                            <option value="DEBIT">Debit</option>
                            <option value="CREDIT">Credit</option>
                        </select>
                    </div>
                </div>

                <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                            <tr>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Description</th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Amount</th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">To Account</th>
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                            {transactions.map((transaction) => (
                                <tr key={transaction.id} className="hover:bg-gray-50">
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                        {formatDate(transaction.transaction_date)}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                        {transaction.description}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap">
                                        <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                                            transaction.transaction_type === 'DEBIT' 
                                                ? 'bg-red-100 text-red-800' 
                                                : 'bg-green-100 text-green-800'
                                        }`}>
                                            {transaction.transaction_type}
                                        </span>
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                        {formatAmount(transaction.amount)}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                        {transaction.to_account_id}
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>

                {transactions.length === 0 && (
                    <div className="text-center py-8 text-gray-500">
                        No transactions found
                    </div>
                )}
            </div>
        </div>
    );
};

export default TransactionHistoryPage; 