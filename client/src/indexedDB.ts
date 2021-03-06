//store file info in db
export const addFile = (fileInfo: any) => {
    return new Promise((resolve, reject) => {
        getDB().then(function (db: IDBDatabase) {
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

// gets a file by id from indexedDB
export const getFile = (fileID: string) => {
    return new Promise((resolve: any) => {
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

// gets all files from indexedDB
export const getAllFiles = () => {
    return new Promise((resolve: any) => {
        getDB().then(function(db: IDBDatabase) {
            var request: any
            const tx = db.transaction("files", "readonly");
            const store = tx.objectStore("files");
            request = store.getAll();
            request.onerror = (event:any) => {
                console.log("Error: " + (event.target as any).errorCode)
            }
            request.onsuccess = () => {
                return resolve(request.result);
            }
            request.onupgradeneeded = () => {
                console.log("upgrade needed");
                var db = request.result;
                db.createObjectStore("files");
            }
        })
    })
}

//delete file info from db
export const rmFile = (fileID: string) => {
    return new Promise((resolve: any) => {
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

// delete all files from db
export const clearFiles = () => {
    return new Promise((resolve: any) => {
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

// get indexedDB db
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