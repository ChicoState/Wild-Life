import axios from 'axios'

export class Upload {


    constructor() {
        this.form = new FormData();
        this.onsuccess = (data) => {
        }
        this.onfail = (error) => {
        }
    }

    addFile(file) {
        this.form.append("file", file)
    }

    submit() {
        const headers = {'Content-Type': 'multipart/form-data'};
        const request = axios.post('http://localhost:5069/upload', this.form, {headers})
        request.then(this.onsuccess).catch(this.onfail)
    }


}

export default {
    Upload
}