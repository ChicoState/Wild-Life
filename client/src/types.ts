import {Buffer} from "buffer";

export interface Rectangle {
    Min: {
        X: number,
        Y: number,
        Max: {
            X: number,
            Y: number
        }
    }
}

export interface Detection {
    bounds: Rectangle;
    confidence: number;
    classes: number[],
    boxes: Rectangle[];
    confidences: number[];
    type: string
}

export interface fileType {
    id: string;
    name: string;
    type: string;
    size: number;
    data: string;
    status?: number;
    plant?: string;
    confidence?: number;
}

export interface User {

}

export interface Upload {
    thumbnail: Buffer
    result: Buffer
    threshold: Buffer
    highlight: Buffer
    confidence: number
}


export interface Application {
    user: User
    uploads: Upload
}
