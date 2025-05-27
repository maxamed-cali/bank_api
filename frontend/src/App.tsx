import { Provider } from 'react-redux';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { ThemeProvider } from './contexts/theme-context';
import { store } from './store/store';
import MainLayout from './layouts/main-layout';
import DashboardPage from './pages/Dashboard';
import Auth from './pages/Auth';
import { ProtectedRoute } from './components/ProtectedRoute';
import { Toaster } from 'react-hot-toast';
import NotificationsPage from './pages/NotificationsPage';
import WebSocketInitializer from './components/WebSocketInitializer';
import UserManagement from './pages/UserManagement';

// Import your other components here
// import Dashboard from './pages/Dashboard';
import AccountRegistration from './pages/AccountRegistration';
import TransferFundsPage from './pages/TransferFundsPage';
import MoneyRequestPage from './pages/MoneyRequestPage';
import TransactionHistoryPage from './pages/TransactionHistory';
import RoleAssignment from './pages/RoleAssignment';

function App() {
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
          element: <DashboardPage />,
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
          element: <RoleAssignment /> ,
        },
        {
          path: "reports",
          element: <h1 className="title">Role Assing</h1>,
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
        {
          path: "settings",
          element: <h1 className="title">Settings</h1>,
        },
      ],
    },
  ]);


  return (
    <Provider store={store}>
      <ThemeProvider storageKey="theme">
        <Toaster position="top-right" />
        <WebSocketInitializer />
        <RouterProvider router={router} />
      </ThemeProvider> 
    </Provider>
  );
}

export default App; 