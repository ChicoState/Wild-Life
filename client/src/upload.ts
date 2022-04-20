import axios, {AxiosResponse} from 'axios'
import {Websocket, WebsocketBuilder} from 'websocket-ts';

let env = import.meta.env

export interface UploadState {
    time: string
    state: string
    message: string
    data: any
}

// The expected data format for a response from the server
export interface UploadResponse {
    data: {
        token: string
    }
}

export class Upload {
    token?: string
    socket?: Websocket
    update: (res: UploadState) => void
    error: (res: any) => void
    form: FormData
    state: UploadState

    // Initializes the Upload class
    constructor() {
        this.form = new FormData();
        this.state = {time: "", state: "uploading"} as UploadState
        this.token = ""
        this.update = () => {
        }
        this.error = () => {
        }
    }

    localURL(): string {
        return `http://localhost:5069/upload`
    }

    // Returns the relevant rest endpoint url
    restURL(): string {
        let url = `http://localhost:5069/upload`
        if (env.VITE_ENDPOINT || env.ENDPOINT || env.PROD || env.VITE_PROD) {
            url = `https://${env.VITE_ENDPOINT}/api/upload`
        }
        return url
    }

    // Returns the relevant websocket endpoint url
    socketURL(): string {
        let url = `ws://localhost:5069/sockets/${this.token}`
        if (env.VITE_ENDPOINT || env.ENDPOINT || env.PROD || env.VITE_PROD) {
            url = `wss://${env.VITE_ENDPOINT}/api/sockets/${this.token}`
        }
        return url
    }

    // Add file to the submission
    addFile(file: File) {
        this.form.append("file", file)
    }

    // Upload and submit the image to the servers
    submit(): void {
        // Define the request parameters
        const headers = {'Content-Type': 'multipart/form-data'};
        // Make the post request to the servers
        axios.post(this.restURL(), this.form, {headers})
            .then((res: AxiosResponse) => {
                // Handle successful requests
                this.uploadSuccess(res.data)
            })
            .catch((res: any) => {
                // Handle unsuccessful requests
                this.uploadFailure(res)
            })
    }

    // Called when a file is successfully uploaded
    uploadSuccess(res: UploadResponse) {
        this.token = res.data.token
        let ref = this
        // Build the WebSocket instance
        const builder = new WebsocketBuilder(this.socketURL())
            .onOpen((instance: Websocket, ev: Event): any => {
                console.log("Opened: " + ev)
            })
            .onClose((instance: Websocket, ev: CloseEvent): any => {
                console.log("Closed: " + ev)
            })
            .onMessage((instance: Websocket, ev: MessageEvent): any => {
                let parsed = JSON.parse(ev.data)
                ref.state = parsed as UploadState
                ref.update(this.state)
            })
            .onError((instance: Websocket, ev: Event): any => {
                ref.error(ev)
            })
        // Assign the socket to a local variable
        this.socket = builder.build()
    }

    // Called when a file is unsuccessfully uploaded
    uploadFailure(res: any) {
        this.error(JSON.stringify(res))
    }

}
