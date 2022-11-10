import {createContext, FunctionComponent, ReactNode, useContext, useEffect, useState} from "react";

type ThemeProviderProps = {
    children: ReactNode
}

const getTheme = (): string => {
    let theme: string
    let storedTheme: string | null = localStorage.getItem('theme')
    if ( storedTheme !== null ) {
        theme = storedTheme
    } else {
        if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
            theme = 'dark'
        } else {
            theme = 'light'
        }
    }
    if (theme === 'dark') {
        document.documentElement.setAttribute('data-theme', 'dark')
    } else {
        document.documentElement.setAttribute('data-theme', 'light')

    }
    return theme
}

let currentTheme = getTheme()

const ThemeState = {
    theme: currentTheme,
    setTheme: (theme: string) => {}
}

const ThemeContext = createContext<typeof ThemeState>({
    theme: currentTheme,
    setTheme: (theme: string) => {}
});

export const ThemeProvider = ({children}: ThemeProviderProps) => {
    const [theme, setTheme] = useState(currentTheme)

    useEffect(() => {
        localStorage.setItem("theme", theme);
        if (theme === 'dark') {
            document.documentElement.setAttribute('data-theme', 'dark')
        } else {
            document.documentElement.setAttribute('data-theme', 'light')
        }
    }, [theme]);

    return (
        <ThemeContext.Provider value={{theme, setTheme}}>
            {children}
        </ThemeContext.Provider>
    );
}

export const useTheme = () => useContext(ThemeContext);