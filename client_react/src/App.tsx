import Header from "./components/header/header";
import Footer from "./components/footer/footer";
import {ThemeProvider} from "./theme/useTheme";
import MainSection from "./mainSection/mainSection";

function App() {
    return (
        <ThemeProvider>
            <div className="app">
                <Header />
                <MainSection />
                <Footer />
            </div>
        </ThemeProvider>
  )
}

export default App
