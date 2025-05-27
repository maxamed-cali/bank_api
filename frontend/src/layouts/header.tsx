import { useTheme } from "@/hooks/use-theme";
import { Bell, ChevronsLeft, LogOut, Moon, Search, Sun, User, Key } from "lucide-react";
import profileImg from "@/assets/profile-image.jpg";
import { useDispatch, useSelector } from "react-redux";
import { useNavigate } from "react-router-dom";
import { logout } from "@/store/features/auth/authSlice";
import NotificationsPanel from "@/components/NotificationsPanel";
import PasswordResetForm from "@/components/PasswordResetForm";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@radix-ui/react-dropdown-menu";
import { selectNewNotificationsCount } from "@/store/features/notifications/notificationsSlice";
import { useState } from "react";

interface HeaderProps {
    collapsed: boolean;
    setCollapsed: (collapsed: boolean) => void;
}

export const Header = ({ collapsed, setCollapsed }: HeaderProps): JSX.Element => {
    const { theme, setTheme } = useTheme();
    const dispatch = useDispatch();
    const navigate = useNavigate();
    const newNotificationsCount = useSelector(selectNewNotificationsCount);
    const [showPasswordReset, setShowPasswordReset] = useState(false);

    const handleLogout = () => {
        dispatch(logout());
        navigate("/auth");
    };

    return (
        <>
        <header className="relative z-10 flex h-[60px] items-center justify-between bg-white px-4 shadow-md transition-colors dark:bg-slate-900">
            <div className="flex items-center gap-x-3">
                <button
                    className="btn-ghost size-10"
                    onClick={() => setCollapsed(!collapsed)}
                >
                    <ChevronsLeft className={collapsed ? "rotate-180" : ""} />
                </button>
                <div className="input">
                    <Search
                        size={20}
                        className="text-slate-300"
                    />
                    <input
                        type="text"
                        name="search"
                        id="search"
                        placeholder="Search..."
                        className="w-full bg-transparent text-slate-900 outline-0 placeholder:text-slate-300 dark:text-slate-50"
                    />
                </div>
            </div>
            <div className="flex items-center gap-x-3">
                <button
                    className="btn-ghost size-10"
                    onClick={() => setTheme(theme === "light" ? "dark" : "light")}
                >
                    <Sun
                        size={20}
                        className="dark:hidden"
                    />
                    <Moon
                        size={20}
                        className="hidden dark:block"
                    />
                </button>
                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                        <button className="btn-ghost size-10 relative">
                            <Bell size={20} />
                            {newNotificationsCount > 0 && (
                                <span className="absolute -top-1 -right-1 flex h-4 w-4 items-center justify-center rounded-full bg-red-500 text-[10px] text-white">
                                    {newNotificationsCount}
                                </span>
                            )}
                        </button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent 
                        className="w-[400px] p-0 rounded-md bg-white shadow-md dark:bg-slate-900"
                        align="end"
                    >
                        <NotificationsPanel />
                    </DropdownMenuContent>
                </DropdownMenu>
                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                        <button className="size-10 overflow-hidden rounded-full">
                            <img
                                src={profileImg}
                                alt="profile image"
                                className="size-full object-cover"
                            />
                        </button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent className="w-56 rounded-md bg-white p-1 shadow-md dark:bg-slate-900">
                        <DropdownMenuLabel className="px-2 py-1.5 text-sm font-medium text-slate-900 dark:text-slate-50">
                            My Account
                        </DropdownMenuLabel>
                        <DropdownMenuSeparator className="my-1 h-px bg-slate-200 dark:bg-slate-700" />
                        {/*
                        <DropdownMenuItem
                            className="flex cursor-pointer items-center gap-x-2 rounded-sm px-2 py-1.5 text-sm text-slate-900 outline-none hover:bg-slate-100 dark:text-slate-50 dark:hover:bg-slate-800"
                        >
                            <User size={16} />
                            Profile
                        </DropdownMenuItem>
                        */}
                        <DropdownMenuItem
                            className="flex cursor-pointer items-center gap-x-2 rounded-sm px-2 py-1.5 text-sm text-slate-900 outline-none hover:bg-slate-100 dark:text-slate-50 dark:hover:bg-slate-800"
                            onSelect={() => setShowPasswordReset(true)}
                        >
                            <Key size={16} />
                            Password Reset
                        </DropdownMenuItem>
                        <DropdownMenuItem
                            onClick={handleLogout}
                            className="flex cursor-pointer items-center gap-x-2 rounded-sm px-2 py-1.5 text-sm text-slate-900 outline-none hover:bg-slate-100 dark:text-slate-50 dark:hover:bg-slate-800"
                        >
                            <LogOut size={16} />
                            Logout
                        </DropdownMenuItem>
                    </DropdownMenuContent>
                </DropdownMenu>
            </div>
        </header>
            {showPasswordReset && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
                    <div className="rounded-md bg-white dark:bg-slate-900">
                        <PasswordResetForm onClose={() => setShowPasswordReset(false)} />
                    </div>
                </div>
            )}
        </>
    );
}; 