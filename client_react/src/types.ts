export type ScrapeResponse = {
    videoURL: string,
    img: string,
    videoId: string,
}

export type FormProps = {
    setResponse: (response: ScrapeResponse) => void
}

export type WarningProps = {
    error?: string
}