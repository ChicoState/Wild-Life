import axios from 'axios'


export interface UploadState {
    time: string
    state: string
    message: string
    data: any
}

export interface UploadResponse {
    data: {
        data: {
            token: string
        }
    }
}

export class Upload {
    token?: string
    form: FormData
    update: (res: UploadState) => void
    error: (res: UploadState) => void
    state: UploadState
    history: UploadState[]
    socket?: WebSocket

    constructor() {
        this.form = new FormData();
        this.state = {time: "", state: "uploading"} as UploadState
        this.token = ""
        this.update = () => {
        }
        this.error = () => {
        }
        this.history = []
    }

    // Upload

    addFile(file: File) {
        this.form.append("file", file)
    }

    submit() {
        const headers = {'Content-Type': 'multipart/form-data'};
        const self = this
        axios.post('http://localhost:5069/upload', this.form, {headers})
            .then((res) => {
                self.uploadSuccess(res)
            }).catch((res) => {
            self.uploadFailure(res)
        })

    }

    uploadSuccess(res: UploadResponse) {
        this.token = res.data.data.token
        console.log("Upload Success")
        let conn = `ws://localhost:5069/sockets/${this.token}`
        console.log(conn)
        this.socket = new WebSocket(conn)
        let ref = this
        this.socket.onmessage = (res: MessageEvent) => {
            let parsed = JSON.parse(res.data)
            ref.history.push(parsed as UploadState)
            ref.state = parsed as UploadState
            ref.update(ref.state)
        }
        this.socket.onopen = this.wsOpen
        this.socket.onerror = this.wsError
        this.socket.onclose = this.wsClose
    }

    uploadFailure(res: any) {
        console.log(res.data)
    }

    // Websockets

    wsOpen(res: Event) {
        console.log(res)
    }

    wsError(res: Event) {
        console.log(res)
    }

    wsClose(res: CloseEvent) {
        console.log(res)
    }


}
