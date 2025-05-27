import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import axiosInstance from '../../../services/axios';

interface Notification {
  id: string;
  message: string;
  isNew?: boolean;
  createdAt: string;
}

interface NotificationsState {
  dbNotifications: Notification[];
  websocketNotifications: Notification[];
  isLoading: boolean;
  error: string | null;
  activeFilter: 'all' | 'requests' | 'alerts';
  newNotificationsCount: number;
}

const initialState: NotificationsState = {
  dbNotifications: [],
  websocketNotifications: [],
  isLoading: false,
  error: null,
  activeFilter: 'all',
  newNotificationsCount: 0,
};

export const fetchNotifications = createAsyncThunk(
  'notifications/fetchNotifications',
  async (userId: string) => {
    try {
      const response = await axiosInstance.get(`/api/user/notifications?user_id=${userId}`);
      // Transform the API response to include id and map created_at to createdAt
      // Assuming 'isNew' status comes from the backend for fetched notifications
      const notifications = response.data.map((notification: any) => ({
        id: notification.id || Date.now().toString() + Math.random().toString(16).slice(2), // Ensure unique ID
        message: notification.message,
        isNew: notification.is_new, // Assuming backend provides is_new
        createdAt: notification.created_at,
      }));
      return notifications;
    } catch (error: any) {
      throw new Error(error.response?.data?.message || 'Failed to fetch notifications');
    }
  }
);

const notificationsSlice = createSlice({
  name: 'notifications',
  initialState,
  reducers: {
    addNotification: (state, action: PayloadAction<Notification>) => {
      // Add websocket notification with isNew: true by default
      state.websocketNotifications.unshift({...action.payload, isNew: true, id: action.payload.id || Date.now().toString() + Math.random().toString(16).slice(2) }); // Ensure unique ID
      state.newNotificationsCount++; // Increment count for new websocket messages
    },
    setActiveFilter: (state, action: PayloadAction<'all' | 'requests' | 'alerts'>) => {
      state.activeFilter = action.payload;
    },
    markAsRead: (state, action: PayloadAction<string>) => {
      const notificationDb = state.dbNotifications.find(n => n.id === action.payload);
      const notificationWs = state.websocketNotifications.find(n => n.id === action.payload);

      if (notificationDb && notificationDb.isNew) {
        notificationDb.isNew = false;
        state.newNotificationsCount--;
      } else if (notificationWs && notificationWs.isNew) {
        notificationWs.isNew = false;
        state.newNotificationsCount--;
      }
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchNotifications.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchNotifications.fulfilled, (state, action: PayloadAction<Notification[]>) => {
        state.isLoading = false;
        state.dbNotifications = action.payload; // Replace DB notifications with fetched ones
        // Recalculate new count based on both DB and websocket notifications
        state.newNotificationsCount = state.dbNotifications.filter(n => n.isNew).length + 
                                      state.websocketNotifications.filter(n => n.isNew).length;
      })
      .addCase(fetchNotifications.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'Failed to fetch notifications';
      });
  },
});

export const { addNotification, setActiveFilter, markAsRead } = notificationsSlice.actions;
export default notificationsSlice.reducer;

// Selectors
export const selectNotifications = (state: { notifications: NotificationsState }) =>
  [...state.notifications.websocketNotifications, ...state.notifications.dbNotifications]; // Combine both lists
export const selectNotificationsLoading = (state: { notifications: NotificationsState }) =>
  state.notifications.isLoading;
export const selectNotificationsError = (state: { notifications: NotificationsState }) =>
  state.notifications.error;
export const selectActiveFilter = (state: { notifications: NotificationsState }) =>
  state.notifications.activeFilter;
export const selectNewNotificationsCount = (state: { notifications: NotificationsState }) =>
  state.notifications.newNotificationsCount; 