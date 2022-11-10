import './warning.css'

type WarningProps = {
    error?: string
}

const Warning = ({error = "Please enter a valid TikTok URL."}: WarningProps) => {
    return (
        <p className="url_warn">{error}</p>
    )
}
export default Warning