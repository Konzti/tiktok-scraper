import {IMG_RATIO} from "../../constants";
import { requestDownload } from "../../utils/utils";
import './container.css'

export type ResponseProps = {
    videoURL: string,
    img: string,
    videoId: string,
}

const VideoContainer = (props: ResponseProps) => {
    const openVideo = () => {
        window.open(props.videoURL, '_blank');
    }
    const downloadVideo = async () => {
        await requestDownload(props.videoId)
    }

    return (
        <div className="video_container" style={{backgroundImage: `url(${props.img})`, aspectRatio: `${IMG_RATIO}`}}>
            <div className="inner_container">
                <button onClick={openVideo} className="download_btn">Open video in new tab</button>
                <button onClick={downloadVideo} className="download_btn">Download video</button>
            </div>
        </div>
    )
}

export default VideoContainer