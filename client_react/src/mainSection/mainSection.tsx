import Form from "../components/form/form";
import InfoSection from "../components/info/infoSection";
import {useState} from "react";
import VideoContainer from "../components/videoContainer/videoContainer";
import {ScrapeResponse} from "../types";

const MainSection = () => {
    const [response, setResponse] = useState<ScrapeResponse | null>(null);

    return (
        <main className="main">
            <Form setResponse={setResponse}/>
            {response !== null && <VideoContainer {...response}/>}
            <InfoSection/>
        </main>
    )
}

export default MainSection
