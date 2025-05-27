import ProductImage from "@/assets/product-image.jpg";
import ProfileImage from "@/assets/profile-image.jpg";
import { BanknoteIcon,DollarSign,BarChart, Home, LucideIcon, NotepadText, Settings, UserCheck, UserPlus, Users } from "lucide-react";


interface NavLink {
    label: string;
    icon: LucideIcon;
    path: string;
    role: string;
}

interface NavGroup {
    title: string;
    role: string;
    links: NavLink[];
}

interface OverviewData {
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

interface Product {
    number: number;
    name: string;
    image: string;
    description: string;
    price: number;
    status: string;
    rating: number;
}

export const navbarLinks: NavGroup[] = [
    {
        title: "Dashboard",
        role: "all",
        links: [
            {
                label: "Dashboard",
                icon: Home,
                path: "/",
                role: "all",
            },
          

           
        ],
    },
    {
        title: "Customers",
        role: "all",
        links: [
            
            {
                label: "Open New Account",
                icon: BanknoteIcon,
                path: "/account-registration",
                role: "User",
            },
            {
                label: "Send Money",
                icon: DollarSign,
                path: "/transfer-funds",
                role: "User",
            },
            {
                label: "Request Money",
                icon: UserCheck,
                path: "/money-request",
                role: "User",
            },
        ],
    },

    {
        title: "Admin",
        role: "all",
        links: [
            
            {
                label: "Users",
                icon: Users,
                path: "/users",
                role: "Admin",
            },
            {
                label: "Role Assign",
                icon: NotepadText,
                path: "/Assign",
                role: "Admin",
            },
            {
                label: "audit logs",
                icon: BarChart,
                path: "/customers",
                role: "Admin",
            },
        ],
    },
    
];



export const recentSalesData: Sale[] = [
    {
        id: 1,
        name: "Olivia Martin",
        email: "olivia.martin@email.com",
        image: ProfileImage,
        total: 1500,
    },
    {
        id: 2,
        name: "James Smith",
        email: "james.smith@email.com",
        image: ProfileImage,
        total: 2000,
    },
    {
        id: 3,
        name: "Sophia Brown",
        email: "sophia.brown@email.com",
        image: ProfileImage,
        total: 4000,
    },
    {
        id: 4,
        name: "Noah Wilson",
        email: "noah.wilson@email.com",
        image: ProfileImage,
        total: 3000,
    },
    {
        id: 5,
        name: "Emma Jones",
        email: "emma.jones@email.com",
        image: ProfileImage,
        total: 2500,
    },
    {
        id: 6,
        name: "William Taylor",
        email: "william.taylor@email.com",
        image: ProfileImage,
        total: 4500,
    },
    {
        id: 7,
        name: "Isabella Johnson",
        email: "isabella.johnson@email.com",
        image: ProfileImage,
        total: 5300,
    },
];

