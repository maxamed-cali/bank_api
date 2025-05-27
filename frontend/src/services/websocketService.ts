// services/websocketService.ts
import toast from 'react-hot-toast';
import { store } from '../store/store';
import { addNotification } from '../store/features/notifications/notificationsSlice';

let socket: WebSocket | null = null;
let reconnectAttempts = 0;
const MAX_RECONNECT_ATTEMPTS = 5;
const RECONNECT_DELAY = 3000;

// We will manage message count outside this service, perhaps in a Redux slice

export const websocketService = {
  connect: (userId: string) => {
    const url = `ws://localhost:8080/ws?user_id=${userId}`;
    console.log('Attempting to connect to WebSocket at:', url);
    
    try {
      socket = new WebSocket(url);

      socket.onopen = () => {
        console.log('WebSocket connected successfully');
        reconnectAttempts = 0;
      };

      socket.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          console.log('WebSocket message received:', data);
          
          if (data && data.message) {
            // Show toast notification
            toast.success(data.message);
            
            // Add to Redux store
            store.dispatch(addNotification({
              id: Date.now().toString(), // Generate unique ID
              message: data.message,
              createdAt: new Date().toISOString(),
              isNew: true,
              userId: data.user_id
            }));
          }

        } catch (error) {
          console.error('Error parsing WebSocket message:', error);
        }
      };

      socket.onerror = (event) => {
        console.error('WebSocket error details:', {
          readyState: socket?.readyState,
          url: url,
          timestamp: new Date().toISOString()
        });
      };

      socket.onclose = (event) => {
        console.warn('WebSocket closed with details:', {
          code: event.code,
          reason: event.reason,
          wasClean: event.wasClean,
          timestamp: new Date().toISOString()
        });
        
        if (reconnectAttempts < MAX_RECONNECT_ATTEMPTS) {
          reconnectAttempts++;
          console.log(`Attempting to reconnect (${reconnectAttempts}/${MAX_RECONNECT_ATTEMPTS})...`);
          setTimeout(() => {
            websocketService.connect(userId);
          }, RECONNECT_DELAY);
        } else {
          console.error('Max reconnection attempts reached');
        }
      };
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error);
    }
  },

  disconnect: () => {
    if (socket) {
      socket.close();
      socket = null;
      reconnectAttempts = 0;
    }
  },

  // Send a message
  send: (message: any) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(typeof message === 'string' ? message : JSON.stringify(message));
    } else {
      console.error('WebSocket is not connected. Current state:', socket?.readyState);
    }
  },

  // Check connection status
  isConnected: () => {
    return socket?.readyState === WebSocket.OPEN;
  }
};
