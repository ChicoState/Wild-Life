import {Upload} from "../upload";

test('Check for main routes', () => {
    let upload = new Upload()
    upload.error = (error) => {
        fail(error)
    }
    upload.update = () => {
        return
    }
});

