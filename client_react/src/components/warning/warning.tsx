import {WarningProps} from "../../types";
import './warning.css'


const Warning = ({error = "Please enter a valid TikTok URL."}: WarningProps) => {
    return (
        <p className="url_warn">{error}</p>
    )
}
export default Warning