import {useState, useRef} from "react";

import {isValidUrl} from "../../utils/utils";
import Spinner from "../spinner/spinner";
import Warning from "../warning/warning";
import {API, ENDPOINT} from "../../constants";
import paste from "../../assets/paste.svg";
import clear from "../../assets/clear.svg";
import './form.css';
import {ResponseProps} from "../videoContainer/videoContainer";

type FormProps = {
    setResponse: (response: ResponseProps) => void
}

const Form = ({setResponse}: FormProps) => {
    const [inputUrl, setInputUrl] = useState<string>("");
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string>("");
    const inputRef = useRef<HTMLInputElement>(null);
    const currentUrlRef = useRef<string | null>(null);

    const submitHandler = async () => {
        if (!isValidUrl(inputUrl)) {
            setError("Please enter a valid TikTok URL.");
            inputRef.current?.focus();
            return;
        }
        if (inputUrl === currentUrlRef.current) {
            inputRef.current?.focus();
            return;
        }
        currentUrlRef.current = inputUrl;
        setLoading(true);
        try {
            let body = JSON.stringify({url: inputUrl});
            const response = await fetch(API + ENDPOINT, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body
            })
            if (response.ok) {
                let data = await response.json()
                console.log(data.data)
                setResponse(data.data)

            } else if (response.status === 404) {
                setError("Whoops, no video found")
            }
            else {
                setError("Please enter a valid TikTok URL.")
                inputRef.current?.focus()
            }
        }
        catch (e: any) {
            console.log("Error: ", e)
            setError(e.message);
        }
        finally {
            setLoading(false);
        }
    }

    const clearOrPasteInput = () => {
        setError("");
        if (inputUrl.trim() !== "") {
            setInputUrl("")
            inputRef.current?.focus();
        } else {
            navigator.clipboard.readText().then((text) => {
                setInputUrl(text);
            }).catch((err) => {
                setError("Please allow clipboard access in your browser settings.");
            });
        }

    }

    return (
        <>
            <div className="form_container">
                <h2>Enter TikTok Video URL:</h2>
                <form>
                    <div className="form_group">
                        <input
                            ref={inputRef}
                            className="input"
                            type="url"
                            name="url"
                            placeholder="Enter Tiktok Video URL"
                            value={inputUrl}
                            onChange={({target}) => setInputUrl(target.value)}/>
                        <span className="clear_btn" onClick={clearOrPasteInput}>
                        {inputUrl.trim().length > 0 ? <ClearButton/> : <PasteButton/>}
                    </span>
                    </div>
                    <input
                        disabled={loading}
                        className="btn_submit"
                        type="button"
                        value="Get Video"
                        onClick={submitHandler}/>
                </form>
            </div>
            { loading ? <Spinner/> : null }
            { error !== "" ? <Warning error={error} /> : null }
        </>

    )
}

export default Form

const PasteButton = () => {
    return (
        <>
            <img src={paste} alt="paste"/>
            <span>Paste</span>
        </>
    )
}
const ClearButton = () => {
    return (
        <>
            <img src={clear} alt="paste"/>
            <span>Clear</span>
        </>
    )
}


