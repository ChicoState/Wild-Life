import axios from 'axios'

export class Upload {

    onsuccess(data:any) {}
    onfail(error:any) {}
    form: FormData;

    constructor() {
        this.form = new FormData();

    }

    addFile(file: File) {
        this.form.append("file", file)
    }

    submit() {
        console.log("Uploading...")
        const headers = {'Content-Type': 'multipart/form-data'};
        return axios.post('http://localhost:5069/upload', this.form, {headers})
    }


}

export default {
    Upload
}