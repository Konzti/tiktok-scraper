import {API, BLOB_VIDEO, HOST, HOST_M, HOST_M2} from "../constants";

export const reloadPage = () => {
    window.location.href = '/'
}

export const isValidUrl = (urlString: string) => {
    let a  = document.createElement('a');
    a.href = urlString;
    return (a.host && a.host !== window.location.host && (a.host === HOST || a.host === HOST_M || a.host === HOST_M2));
}

const downloadFile = async (response: Response) => {
    let fileName
    const header = response.headers.get('Content-Disposition');
    const parts = header?.split(';');
    fileName = parts![1].split('=')[1].replaceAll("\"", "");
    let blob = await response.blob()
    if (blob !== null) {
        let url = window.URL.createObjectURL(blob);
        let a = document.createElement('a');
        a.href = url;
        a.download = fileName;
        document.body.appendChild(a);
        a.click();
        a.remove();
    }
}

export const requestDownload = async (videoId: string) => {
    let res = await fetch(API + BLOB_VIDEO + videoId, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/octet-stream'
        },
    })
    if (!res.ok) {
        alert("Error: " + res.status)
        return
    }
    await downloadFile(res)
}