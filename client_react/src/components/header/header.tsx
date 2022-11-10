import icon from '../../assets/icon.svg'
import dark from '../../assets/dark.svg'
import light from '../../assets/light.svg'
import {reloadPage} from "../../utils/utils";
import {useTheme} from "../../theme/useTheme";
import './header.css'


const Header = () => {
    const {theme, setTheme} = useTheme()
    return (
        <header className="header">
            <Logo/>
            <div className="heading">
                <h2>TikTok Downloader</h2>
            </div>
            <ThemeSwitcher theme={theme} setTheme={setTheme} />
        </header>
    )
}
export default Header

const Logo = () => {
    return (
        <div className="icon_wrap" onClick={reloadPage}>
            <img className="icon" alt="icon" src={icon}/>
        </div>
    )
}

type ThemeSwitcherProps = {
    theme: string
    setTheme: (theme: string) => void
}

const ThemeSwitcher = ({theme, setTheme}: ThemeSwitcherProps) => {
    function toggleTheme() {
        if (theme === "dark") {
            setTheme("light")
        } else {
            setTheme("dark")
        }
    }

    return (
        <div className="theme_btn_wrap" onClick={toggleTheme}>
            <span className="theme_btn">{theme === "dark" ? <img src={light} alt="theme"/> : <img src={dark} alt="theme"/> }</span>
        </div>
    )
}
