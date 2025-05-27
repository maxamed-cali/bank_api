import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { AppDispatch } from '../store/store';
import { 
  markAsRead,
  selectNotifications,
  selectNewNotificationsCount
} from '../store/features/notifications/notificationsSlice';

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

const NotificationsPanel: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
  const navigate = useNavigate();
  const notifications = useSelector(selectNotifications);
  const newNotificationsCount = useSelector(selectNewNotificationsCount);

  const handleNotificationClick = (notificationId: string) => {
    dispatch(markAsRead(notificationId));
  };

  const handleViewAll = () => {
    navigate('/notifications');
  };

  return (
    <div className="notifications-panel p-4 bg-white rounded-lg shadow">
      <h2 className="text-xl font-bold mb-4">Notifications</h2>

      {/* Notification List */}
      <div className="notification-list max-h-[50vh] overflow-y-auto">
        {notifications.length === 0 ? (
          <div className="text-center text-gray-500 py-4">No notifications found.</div>
        ) : (
          <div className="space-y-2">
            {notifications.slice(0, 5).map((notification) => (
              <div 
                key={notification.id} 
                className={`p-3 rounded-lg border ${
                  notification.isNew ? 'bg-blue-50 border-blue-200' : 'bg-white border-gray-200'
                }`}
                onClick={() => handleNotificationClick(notification.id)}
              >
                <div className="flex justify-between items-center">
                  <div className="flex-1">
                    <div className="flex items-center gap-2">
                      <span className="font-medium">
                        {notification.message}
                      </span>
                      {notification.isNew && (
                        <span className="px-2 py-0.5 text-xs bg-blue-100 text-blue-800 rounded-full">
                          New
                        </span>
                      )}
                    </div>
                    <div className="text-xs text-gray-500 mt-2">
                      {formatDate(notification.createdAt)}
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
      
      {/* View All Button */}
      <div className="mt-4 text-center">
        <button
          onClick={handleViewAll}
          className="px-4 py-2 text-sm font-medium text-blue-600 hover:text-blue-800 hover:bg-blue-50 rounded-md transition-colors"
        >
          View All Notifications
        </button>
      </div>

      <div className="mt-4 text-center">
        <p>You have {newNotificationsCount} new notifications.</p>
      </div>
    </div>
  );
};

export default NotificationsPanel;
