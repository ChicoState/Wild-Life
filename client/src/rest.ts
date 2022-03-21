import axios from 'axios'

export class Upload {

    onsuccess(data) {}
    onfail(error) {}
    form: FormData;

    constructor() {
        this.form = new FormData();
        this.onsuccess = (data) => {
        }
        this.onfail = (error) => {
            console.log(error)
        }
    }

    addFile(file) {
        this.form.append("file", file)
    }

    submit() {
        console.log("Uploading...")
        const headers = {'Content-Type': 'multipart/form-data'};
        const request = axios.post('http://localhost:5069/upload', this.form, {headers})
        request.then(this.onsuccess).catch(this.onfail)
    }


}

export default {
    Upload
}