import { useState } from 'react';
import toast from 'react-hot-toast';
import { authService } from '../services/authService'; // Import authService

interface PasswordResetFormProps {
  onClose: () => void;
}

export default function PasswordResetForm({ onClose }: PasswordResetFormProps) {
  const [OldPassword, setOldPassword] = useState('');
  const [NewPassword, setNewPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    
    try {
      await authService.changePassword(OldPassword, NewPassword);
      toast.success('Password reset successful!');
      onClose();
    } catch (error: any) {
      const errorMessage = error.message || 'Failed to reset password';OldPassword
      toast.error(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 px-4 sm:px-6">
    <div className="bg-white rounded-xl shadow-lg p-8 w-full max-w-md">
      <h2 className="text-2xl font-semibold text-gray-800 mb-6 text-center">Reset Password</h2>
  
      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Old Password */}
        <div>
          <label htmlFor="old-password" className="block text-sm font-medium text-gray-700 mb-1">
            Old Password
          </label>
          <input
            type="password"
            id="old-password"
            name="OldPassword"
            value={OldPassword}
            onChange={(e) => setOldPassword(e.target.value)}
            required
            className="w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 transition duration-150 ease-in-out"
          />
        </div>
  
        {/* New Password */}
        <div>
          <label htmlFor="new-password" className="block text-sm font-medium text-gray-700 mb-1">
            New Password
          </label>
          <input
            type="password"
            id="new-password"
            name="NewPassword"
            value={NewPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            required
            className="w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 transition duration-150 ease-in-out"
          />
        </div>
  
        {/* Buttons */}
        <div className="flex justify-end gap-4 pt-4 border-t border-gray-200">
          <button
            type="button"
            onClick={onClose}
            className="px-4 py-2 text-sm font-medium text-gray-600 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={isLoading}
            className="px-5 py-2 text-sm font-semibold text-white bg-blue-600 rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-60"
          >
            {isLoading ? 'Resetting...' : 'Reset Password'}
          </button>
        </div>
      </form>
    </div>
  </div>
  
  

  );
} 