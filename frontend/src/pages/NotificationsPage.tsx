import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { AppDispatch } from '../store/store';
import { 
  fetchNotifications, 
  setActiveFilter, 
  markAsRead,
  selectNotifications,
  selectNotificationsLoading,
  selectNotificationsError,
  selectActiveFilter
} from '../store/features/notifications/notificationsSlice';
import { selectCurrentUser } from '../store/features/auth/authSlice';
import { BellIcon, ExclamationCircleIcon, InboxIcon } from '@heroicons/react/24/outline';

const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  const now = new Date();
  const diffInHours = Math.abs(now.getTime() - date.getTime()) / 36e5;

  if (diffInHours < 24) {
    if (diffInHours < 1) {
      const minutes = Math.round(diffInHours * 60);
      return `${minutes} minute${minutes === 1 ? '' : 's'} ago`;
    }
    return `${Math.round(diffInHours)} hour${Math.round(diffInHours) === 1 ? '' : 's'} ago`;
  } else {
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
};

const NotificationsPage: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
  const notifications = useSelector(selectNotifications);
  const isLoading = useSelector(selectNotificationsLoading);
  const error = useSelector(selectNotificationsError);
  const activeTab = useSelector(selectActiveFilter);
  const currentUser = useSelector(selectCurrentUser);

  useEffect(() => {
    if (currentUser?.id) {
      dispatch(fetchNotifications(currentUser.id.toString()));
    }
  }, [dispatch, currentUser?.id, activeTab]);

  const filteredNotifications = notifications.filter(notification => {
    if (activeTab === 'all') return true;
    if (activeTab === 'requests') return notification.message.toLowerCase().includes('request');
    if (activeTab === 'alerts') return notification.message.toLowerCase().includes('expired') || 
                                   notification.message.toLowerCase().includes('declined');
    return true;
  });

  const handleTabChange = (tab: 'all' | 'requests' | 'alerts') => {
    dispatch(setActiveFilter(tab));
  };

  const handleNotificationClick = (notificationId: string) => {
    dispatch(markAsRead(notificationId));
  };

  const getNotificationIcon = (type: string) => {
    if (type.toLowerCase().includes('request')) {
      return <InboxIcon className="w-5 h-5 text-blue-500" />;
    }
    if (type.toLowerCase().includes('expired') || type.toLowerCase().includes('declined')) {
      return <ExclamationCircleIcon className="w-5 h-5 text-red-500" />;
    }
    return <BellIcon className="w-5 h-5 text-gray-500" />;
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-4xl mx-auto p-6">
        <div className="bg-white rounded-xl shadow-sm border border-gray-100">
          <div className="p-6">
            <div className="flex items-center justify-between mb-8">
              <div>
                <h1 className="text-2xl font-bold text-gray-900">Notifications</h1>
                <p className="text-sm text-gray-500 mt-1">Stay updated with your latest activities</p>
              </div>
              <div className="flex items-center space-x-2">
                <span className="text-sm text-gray-500">
                  {filteredNotifications.length} {filteredNotifications.length === 1 ? 'notification' : 'notifications'}
                </span>
              </div>
            </div>

            {/* Tabs */}
            <div className="flex space-x-1 bg-gray-50 p-1 rounded-lg mb-6">
              <button
                className={`flex-1 px-4 py-2.5 text-sm font-medium rounded-md transition-all duration-200 ${
                  activeTab === 'all' 
                    ? 'bg-white text-blue-600 shadow-sm' 
                    : 'text-gray-600 hover:text-gray-900 hover:bg-gray-100'
                }`}
                onClick={() => handleTabChange('all')}
              >
                All
              </button>
              <button
                className={`flex-1 px-4 py-2.5 text-sm font-medium rounded-md transition-all duration-200 ${
                  activeTab === 'requests' 
                    ? 'bg-white text-blue-600 shadow-sm' 
                    : 'text-gray-600 hover:text-gray-900 hover:bg-gray-100'
                }`}
                onClick={() => handleTabChange('requests')}
              >
                Requests
              </button>
              <button
                className={`flex-1 px-4 py-2.5 text-sm font-medium rounded-md transition-all duration-200 ${
                  activeTab === 'alerts' 
                    ? 'bg-white text-blue-600 shadow-sm' 
                    : 'text-gray-600 hover:text-gray-900 hover:bg-gray-100'
                }`}
                onClick={() => handleTabChange('alerts')}
              >
                Alerts
              </button>
            </div>

            {/* Notification List */}
            <div className="notification-list">
              {isLoading ? (
                <div className="text-center py-12">
                  <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-600 mx-auto"></div>
                  <p className="mt-4 text-gray-600">Loading your notifications...</p>
                </div>
              ) : error ? (
                <div className="text-center py-12">
                  <ExclamationCircleIcon className="w-12 h-12 text-red-500 mx-auto mb-4" />
                  <p className="text-red-500">{error}</p>
                </div>
              ) : filteredNotifications.length === 0 ? (
                <div className="text-center py-12">
                  <BellIcon className="w-12 h-12 text-gray-400 mx-auto mb-4" />
                  <p className="text-gray-500">No notifications found</p>
                </div>
              ) : (
                <div className="space-y-3">
                  {filteredNotifications.map((notification) => (
                    <div 
                      key={notification.id} 
                      className={`group p-4 rounded-lg border transition-all duration-200 ${
                        notification.isNew 
                          ? 'bg-blue-50 border-blue-200 hover:bg-blue-100' 
                          : 'bg-white border-gray-200 hover:bg-gray-50'
                      } cursor-pointer`}
                      onClick={() => handleNotificationClick(notification.id)}
                    >
                      <div className="flex items-start space-x-4">
                        <div className="flex-shrink-0 mt-1">
                          {getNotificationIcon(notification.message)}
                        </div>
                        <div className="flex-1 min-w-0">
                          <div className="flex items-center justify-between">
                            <p className={`text-sm font-medium ${
                              notification.isNew ? 'text-gray-900' : 'text-gray-700'
                            }`}>
                              {notification.message}
                            </p>
                            {notification.isNew && (
                              <span className="ml-2 px-2 py-0.5 text-xs font-medium bg-blue-100 text-blue-800 rounded-full">
                                New
                              </span>
                            )}
                          </div>
                          <p className="mt-1 text-xs text-gray-500">
                            {formatDate(notification.createdAt)}
                          </p>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default NotificationsPage; 