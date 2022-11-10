import {GITHUB} from "../../constants";
import github_icon from '../../assets/github.png'
import './footer.css'

function Footer () {
    return (
        <footer>
            <section className="disc_section">
                <p>Feel free to reach out to me:</p>
            </section>
            <div className="contact">
                <a href={GITHUB} target="_blank">
                    <div className="contact_item">
                    <img src={github_icon} alt="github"/><p>github.com/Konzti</p>
                    </div>
                </a>
            </div>
        </footer>
)
}
export default Footer