import { useContext } from "react";
import { ThemeProviderContext } from "@/contexts/theme-context";

type Theme = "light" | "dark";

interface ThemeContextType {
    theme: Theme;
    setTheme: (theme: Theme) => void;
}

export const useTheme = (): ThemeContextType => {
    const context = useContext(ThemeProviderContext);

    if (context === undefined) throw new Error("useTheme must be used within a ThemeProvider");

    return context;
}; 