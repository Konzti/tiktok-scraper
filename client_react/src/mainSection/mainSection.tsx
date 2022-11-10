import Form from "../components/form/form";
import InfoSection from "../components/info/infoSection";
import {useState} from "react";
import VideoContainer from "../components/videoContainer/videoContainer";
import {ResponseProps} from "../components/videoContainer/videoContainer";

const MainSection = () => {
    const [response, setResponse] = useState<ResponseProps | null>(null);

    return (
        <main className="main">
            <Form setResponse={setResponse}/>
            {response !== null && <VideoContainer {...response}/>}
            <InfoSection/>
        </main>
    )
}

export default MainSection
