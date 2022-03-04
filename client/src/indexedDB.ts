//store file info in db
export const addFile = (fileInfo) => {
    return new Promise(() => {
        getDB().then(function(db: IDBDatabase) {
            const request = db.transaction(["files"], "readwrite")
                .objectStore("files")
                .add(fileInfo, fileInfo.id);
            request.onsuccess = () => {
                console.log("file added");
            }
            request.onerror = (event) => {
                console.log("Error adding file: ", event);
            }
        })
    })
}

export const getFile = (fileID) => {
    return new Promise((resolve) => {
        getDB().then(function(db: IDBDatabase) {
            const tx = db.transaction(["files"], "readonly");
            const store = tx.objectStore("files");
            const request = store.get(fileID);
            request.onerror = (event) => {
                console.log(event);
            }
            request.onsuccess = (event) => {
                return resolve(request.result);
            }
        });
    })
}

export const getAllFiles = () => {
    return new Promise((resolve) => {
        getDB().then(function(db: IDBDatabase) {
            var request: any
            const tx = db.transaction("files", "readonly");
            const store = tx.objectStore("files");
            request = store.getAll();
            request.onerror = (event) => {
                console.log("Error: "+(event.target as any).errorCode)
            }
            request.onsuccess = () => {
                return resolve(request.result);
            }
            request.onupgradeneeded = () => {
                console.log("upgrade needed");
                var db = request.result;
                db.createObjectStore("files");
            }
        });
    })
}

//delete file info from db
export const rmFile = (fileID) => {
    return new Promise((resolve) => {
        getDB().then(function(db: IDBDatabase) {
            const tx = db.transaction("files", "readwrite");
            const store = tx.objectStore("files");
            const request = store.delete(fileID);
            request.onerror = (event) => {
                console.log("Error: ",event);
            }
            request.onsuccess = (event) => {
                return resolve(request.result);
            }
        })
    })
}

export const clearFiles = () => {
    return new Promise((resolve) => {
        getDB().then(function(db: IDBDatabase) {
            const tx = db.transaction("files", "readwrite");
            const store = tx.objectStore("files");
            const request = store.clear();
            request.onerror = (event) => {
                console.log("Error: ",event);
            }
            request.onsuccess = (event) => {
                return resolve(request.result);
            }
        })
    })
}

export const getDB = () => {return new Promise((resolve: (value: IDBDatabase) => void) => {
        var request = indexedDB.open("file-upload", 1);
        request.onsuccess = (e) => {
            //check if files exists
            return resolve(request.result);
        }
        request.onupgradeneeded = () => {
            console.log("upgrade needed");
            var db = request.result;
            db.createObjectStore("files");
        }
    })
}