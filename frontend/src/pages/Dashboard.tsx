import { recentSalesData } from "@/constants";
import { useTheme } from "@/hooks/use-theme";
import { Footer } from "@/layouts/footer";
import { RootState } from '@/store/store';
import {
    ArrowDownCircle,
    ArrowUpCircle,
    Clock,
    Repeat,
    Scale,
    TrendingUp
} from "lucide-react";
import React, { useEffect, useState } from 'react';
import { useSelector } from 'react-redux';
import { useNavigate } from "react-router-dom";
import { Area, AreaChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import { AccountRegistration, accountService } from '../services/accountService';
import { dashboardService, Transaction } from '../services/dashboard';

interface DashboardCardProps {
    icon: React.ReactNode;
    title: string;
    value: string;
    percentage: string;
    onClick?: () => void;
}

const DashboardCard: React.FC<DashboardCardProps> = ({ icon, title, value, percentage, onClick }) => (
    <div className={`card ${onClick ? 'cursor-pointer' : ''}`} onClick={onClick}>
        <div className="card-header">
            <div className="w-fit rounded-lg bg-blue-500/20 p-2 text-blue-500 transition-colors dark:bg-blue-600/20 dark:text-blue-600">
                {icon}
            </div>
            <p className="card-title">{title}</p>
        </div>
        <div className="card-body bg-slate-100 transition-colors dark:bg-slate-950">
            <p className="text-3xl font-bold text-slate-900 dark:text-slate-50">{value}</p>
            <span className="flex w-fit items-center gap-x-2 rounded-full border border-blue-500 px-2 py-1 font-medium text-blue-500 dark:border-blue-600 dark:text-blue-600">
                <TrendingUp size={18} />
                {percentage}
            </span>
        </div>
    </div>
);

interface ChartData {
    name: string;
    total: number;
}

interface Sale {
    id: number;
    name: string;
    email: string;
    image: string;
    total: number;
}

interface Wallet {
    id: string;
    name: string;
    balance: number;
}

const DashboardPage = () => {
    const { theme } = useTheme();
    const navigate = useNavigate();
    const user = useSelector((state: RootState) => state.auth.user);

    const formatCurrency = (value: number): string => `$${value}`;

    const calculateTotalBalance = (): number => {
        return accounts.reduce((total, account) => total + (account.balance || 0), 0);
    };

    const [accounts, setAccounts] = useState<AccountRegistration[]>([]);
    const [selectedAccount, setSelectedAccount] = useState<AccountRegistration | undefined>();
    const [dashboardData, setDashboardData] = useState<any>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [chartData, setChartData] = useState<ChartData[]>([]);
    const [recentTransactions, setRecentTransactions] = useState<Transaction[]>([]);
    const [transactionsLoading, setTransactionsLoading] = useState(true);
    const [transactionsError, setTransactionsError] = useState<string | null>(null);

    useEffect(() => {
        fetchAccounts();
    }, []);

    const fetchAccounts = async () => {
        try {
            const data = await accountService.getAll();
            setAccounts(data);
        } catch (error: any) {
            console.error('Error fetching accounts:', error);
            setError('Failed to load accounts');
        }
    };

    const handleAccountChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
        const account = accounts.find(acc => acc.account_number === event.target.value);
        setSelectedAccount(account);
    };

    useEffect(() => {
        async function fetchDashboard() {
            try {
                setLoading(true);
                const data = await dashboardService.getTransactionsSummary();
                setDashboardData(data);
            } catch {
                setError("Failed to load dashboard data");
            } finally {
                setLoading(false);
            }
        }
        fetchDashboard();
    }, []);

    useEffect(() => {
        async function fetchChartData() {
            try {
                const data = await dashboardService.getMonthlyTransactions();
                setChartData(data);
            } catch (err) {
                console.error("Error fetching chart data:", err);
                setChartData([]);
            }
        }
        fetchChartData();
    }, []);

    useEffect(() => {
        async function fetchTransactions() {
            try {
                setTransactionsLoading(true);
                const data = await dashboardService.getTransactionHistory(user?.id || '');
                setRecentTransactions(data);
            } catch (err) {
                setTransactionsError("Failed to load transactions");
            } finally {
                setTransactionsLoading(false);
            }
        }
        if (user?.id) {
            fetchTransactions();
        }
    }, [user?.id]);

    return (
        <div className="flex flex-col gap-y-4">
            <h1 className="title">Dashboard</h1>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
                <div className={`card ${selectedAccount ? 'cursor-pointer' : ''}`} >
                    <div className="card-header">
                        <div className="w-fit rounded-lg bg-blue-500/20 p-2 text-blue-500 transition-colors dark:bg-blue-600/20 dark:text-blue-600">
                            <Scale size={26} />
                        </div>
                        <div>
                            <p className="card-title">Account Balance</p>
                            <div className="flex items-center gap-2">
                                <select
                                    className="ml-2 p-1 border rounded"
                                    value={selectedAccount?.account_number || ''}
                                    onChange={handleAccountChange}
                                >
                                    <option value="">Select an account</option>
                                    {accounts.map(account => (
                                        <option key={account.account_number} value={account.account_number}>
                                            {account.account_number}
                                        </option>
                                    ))}
                                </select>
                        
                            </div>
                        </div>
                    </div>
                    
                    <div className="flex items-center gap-4">
                        {selectedAccount && (
                            <p 
                                className="text-2xl font-bold text-blue-600 cursor-pointer hover:underline"
                                onClick={() => navigate(`/transaction-history?account_id=${selectedAccount.account_number}`)}
                            >
                                {selectedAccount.account_number}
                            </p>
                        )}
                        <p className="text-3xl font-bold text-slate-900 dark:text-slate-50">
                            {selectedAccount ? formatCurrency(selectedAccount.balance) : formatCurrency(calculateTotalBalance())}
                        </p>
                    </div>
                    <span className="flex w-fit items-center gap-x-2 rounded-full border border-blue-500 px-2 py-1 font-medium text-blue-500 dark:border-blue-600 dark:text-blue-600">
                        <TrendingUp size={18} />
                        {selectedAccount ? formatCurrency(selectedAccount.balance) : formatCurrency(calculateTotalBalance())}
                    </span>
                </div>

                {loading ? (
                    <div>Loading dashboard...</div>
                ) : error ? (
                    <div className="text-red-500">{error}</div>
                ) : (
                    <>
                        <DashboardCard
                            icon={<Repeat size={26} />}
                            title="Total Transactions"
                            value={dashboardData?.total_transactions ?? "—"}
                            percentage={dashboardData?.total_transactions ?? "—"}
                        />
                        <DashboardCard
                            icon={<Clock size={26} />}
                            title="Pending Requests"
                            value={dashboardData?.pending_requests ?? "—"}
                            percentage={dashboardData?.pending_requests ?? "—"}
                        />
                        <DashboardCard
                            icon={<Repeat size={26} />}
                            title="Total Transfers"
                            value={dashboardData?.total_transfers ?? "—"}
                            percentage={dashboardData?.total_transfers ?? "—"}
                        />
                        <DashboardCard
                            icon={<ArrowUpCircle size={26} className="text-red-500 dark:text-red-400" />}
                            title="Total Sent Amount"
                            value={dashboardData ? formatCurrency(dashboardData.total_sent_amount) : "—"}
                            percentage={dashboardData ? formatCurrency(dashboardData.total_sent_amount) : "—"}
                        />
                        <DashboardCard
                            icon={<ArrowDownCircle size={26} className="text-green-500 dark:text-green-400" />}
                            title="Total Received Amount"
                            value={dashboardData ? formatCurrency(dashboardData.total_received_amount) : "—"}
                            percentage={dashboardData ? formatCurrency(dashboardData.total_received_amount) : "—"}
                        />
                    </>
                )}
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-7">
                <div className="card col-span-1 md:col-span-2 lg:col-span-4">
                    <div className="card-header">
                        <p className="card-title">Monthly Transaction Summary</p>
                    </div>
                    <div className="card-body p-0">
                        <ResponsiveContainer width="100%" height={300}>
                            <AreaChart data={chartData} margin={{ top: 0, right: 0, left: 0, bottom: 0 }}>
                                <defs>
                                    <linearGradient id="colorTotal" x1="0" y1="0" x2="0" y2="1">
                                        <stop offset="5%" stopColor="#2563eb" stopOpacity={0.8} />
                                        <stop offset="95%" stopColor="#2563eb" stopOpacity={0} />
                                    </linearGradient>
                                </defs>
                                <Tooltip cursor={false} formatter={(value: any) => formatCurrency(Number(value))} />
                                <XAxis dataKey="name" strokeWidth={0} stroke={theme === "light" ? "#475569" : "#94a3b8"} tickMargin={6} />
                                <YAxis dataKey="total" strokeWidth={0} stroke={theme === "light" ? "#475569" : "#94a3b8"} tickFormatter={formatCurrency} tickMargin={6} />
                                <Area type="monotone" dataKey="total" stroke="#2563eb" fillOpacity={1} fill="url(#colorTotal)" />
                            </AreaChart>
                        </ResponsiveContainer>
                    </div>
                </div>

                <div className="card col-span-1 md:col-span-2 lg:col-span-3">
                    <div className="card-header">
                        <p className="card-title">Recent Sales</p>
                    </div>
                    <div className="card-body h-[300px] overflow-auto p-0">
                        {recentSalesData?.length ? recentSalesData.map((sale: Sale) => (
                            <div key={sale.id} className="flex items-center justify-between gap-x-4 py-2 pr-2">
                                <div className="flex items-center gap-x-4">
                                    <img src={sale.image} alt={sale.name} className="size-10 flex-shrink-0 rounded-full object-cover" />
                                    <div className="flex flex-col gap-y-2">
                                        <p className="font-medium text-slate-900 dark:text-slate-50">{sale.name}</p>
                                        <p className="text-sm text-slate-600 dark:text-slate-400">{sale.email}</p>
                                    </div>
                                </div>
                                <p className="font-medium text-slate-900 dark:text-slate-50">{formatCurrency(sale.total)}</p>
                            </div>
                        )) : <p className="p-4 text-sm text-slate-500">No recent sales available.</p>}
                    </div>
                </div>
            </div>

            <div className="card">
                <div className="card-header">
                    <p className="card-title">Recent Transactions</p>
                </div>
                <div className="card-body p-0">
                    <div className="relative h-[500px] overflow-auto">
                        <table className="table">
                            <thead className="table-header">
                                <tr className="table-row">
                                    <th className="table-head">Date</th>
                                    <th className="table-head">Type</th>
                                    <th className="table-head">Amount</th>
                                    <th className="table-head">From Account</th>
                                    <th className="table-head">To Account</th>
                                    <th className="table-head">Description</th>
                                </tr>
                            </thead>
                            <tbody className="table-body">
                                {transactionsLoading ? (
                                    <tr>
                                        <td colSpan={6} className="p-4 text-center text-sm text-slate-500">
                                            Loading transactions...
                                        </td>
                                    </tr>
                                ) : transactionsError ? (
                                    <tr>
                                        <td colSpan={6} className="p-4 text-center text-sm text-red-500">
                                            {transactionsError}
                                        </td>
                                    </tr>
                                ) : recentTransactions.length ? (
                                    recentTransactions.map((tx: Transaction, idx: number) => (
                                        <tr key={idx} className="table-row">
                                            <td className="table-cell">{tx.transaction_date.split('T')[0]}</td>
                                            <td className="table-cell">
                                                <span className={`px-2 py-1 rounded-full text-xs font-bold ${tx.transaction_type === 'DEBIT' ? 'bg-red-100 text-red-500' : 'bg-green-100 text-green-500'}`}>
                                                    {tx.transaction_type}
                                                </span>
                                            </td>
                                            <td className="table-cell">${tx.amount}</td>
                                            <td className="table-cell">{tx.account_id}</td>
                                            <td className="table-cell">{tx.to_account_id}</td>
    
                                            <td className="table-cell">{tx.description}</td>
                                        </tr>
                                    ))
                                ) : (
                                    <tr>
                                        <td colSpan={6} className="p-4 text-center text-sm text-slate-500">
                                            No recent transactions available.
                                        </td>
                                    </tr>
                                )}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>

            <Footer />
        </div>
    );
};

export default DashboardPage;
