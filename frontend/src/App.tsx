import { Provider } from 'react-redux';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { ThemeProvider } from './contexts/theme-context';
import { store } from './store/store';
import MainLayout from './layouts/main-layout';
import UserDashboard from './pages/UserDashboard';
import AdminDashboard from './pages/AdminDashboard';
import Auth from './pages/Auth';
import { ProtectedRoute } from './components/ProtectedRoute';
import { Toaster } from 'react-hot-toast';
import NotificationsPage from './pages/NotificationsPage';
import WebSocketInitializer from './components/WebSocketInitializer';
import UserManagement from './pages/UserManagement';
import LogsPage from './pages/LogsPage';
import AccountRegistration from './pages/AccountRegistration';
import TransferFundsPage from './pages/TransferFundsPage';
import MoneyRequestPage from './pages/MoneyRequestPage';
import TransactionHistoryPage from './pages/TransactionHistory';
import RoleAssignment from './pages/RoleAssignment';
import { useSelector } from 'react-redux';
import { RootState } from './store/store';

// Create a separate component for the router
const AppRouter = () => {
  const user = useSelector((state: RootState) => state.auth.user);
  const isAdmin = user?.role === 'Admin';

  const router = createBrowserRouter([
    {
      path: "/auth",
      element: <Auth />,
    },
    {
      path: "/",
      element: (
        <ProtectedRoute>
          <MainLayout />
        </ProtectedRoute>
      ),
      children: [
        {
          index: true,
          element: isAdmin ? <AdminDashboard /> : <UserDashboard />,
        },
        {
          path: "notifications",
          element: <NotificationsPage />,
        },
        {
          path: "users",
          element: <UserManagement />,
        },
        {
          path: "Assign",
          element: <RoleAssignment />,
        },
        {
          path: "logs",
          element: <LogsPage />,
        },
        {
          path: "account-registration",
          element: <AccountRegistration />,
        },
        {
          path: "transfer-funds",
          element: <TransferFundsPage/>,
        },
        {
          path: "money-request",
          element: <MoneyRequestPage/>,
        },
        {
          path: "transaction-history",
          element: <TransactionHistoryPage/>,
        },
        {
          path: "settings",
          element: <h1 className="title">Settings</h1>,
        },
      ],
    },
  ]);

  return <RouterProvider router={router} />;
};

function App() {
  return (
    <Provider store={store}>
      <ThemeProvider storageKey="theme">
        <Toaster position="top-right" />
        <WebSocketInitializer />
        <AppRouter />
      </ThemeProvider> 
    </Provider>
  );
}

export default App; 