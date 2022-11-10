import Header from "./components/header/header";
import Footer from "./components/footer/footer";
import MainSection from "./mainSection/mainSection";
import {ThemeProvider} from "./theme/useTheme";

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
