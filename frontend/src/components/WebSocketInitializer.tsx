import { useEffect } from 'react';
import { useSelector } from 'react-redux';
import { selectCurrentUser } from '../store/features/auth/authSlice';
import { websocketService } from '../services/websocketService';

const WebSocketInitializer = () => {
  const currentUser = useSelector(selectCurrentUser);

  useEffect(() => {
    if (currentUser?.id) {
      websocketService.connect(currentUser.id.toString());
    }

    return () => {
      websocketService.disconnect();
    };
  }, [currentUser?.id]);

  // This component doesn't render anything, it's just for the side effect
  return null;
};

export default WebSocketInitializer; 