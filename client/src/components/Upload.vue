<script setup>
import {reactive} from 'vue';

const data = reactive({
  uploaded: {
    size: 0,
    filename: '',
    time: '',
  },
});

function uploadFile(event) {
    //Upload size in bytes
    data.uploaded.size = event.target.files[0].size;
    //Example file name to display and add to db
    data.uploaded.filename = event.target.files[0].name;

    //time in "dd/mm/yyyy, h:mm:ss PM" format
    data.uploaded.time = new Date().toLocaleString();
    console.log(event)

   // check for IndexedDB support
   //Supported well in chrome
    if (!window.indexedDB) {
        console.log(`Your browser doesn't support IndexedDB`);
        return;
    }

    // open the wildlife database with the version 1
    const request = indexedDB.open('wildlife', 1);

    // create the image object store and indexes
    request.onupgradeneeded = (event) => {
        let db = event.target.result;

        // create the Contacts object store 
        // with auto-increment id
        let store = db.createObjectStore('image', {
            autoIncrement: true
        });

        // create an index on the email property
        let index = store.createIndex('name', 'name', {
            unique: false
        });
    };

    // handle the error event
    request.onerror = (event) => {
        console.error(`Database error: ${event.target.errorCode}`);
    };

    // handle the success event
    request.onsuccess = (event) => {
        const db = event.target.result;

        //Insert the filename, size, and time into db
        insertName(db, {
            name: data.uploaded.filename,
            size: data.uploaded.size,
            time: data.uploaded.time,
        });
    };

    function insertName(db, image) {
        // create a new transaction
        const txn = db.transaction('image', 'readwrite');

        // get the image object store
        const store = txn.objectStore('image');
        //
        let query = store.put(image);

        // handle success case
        query.onsuccess = function (event) {
            console.log(event);
        };

        // handle the error case
        query.onerror = function (event) {
            console.log(event.target.errorCode);
        }

        // close the database once the 
        // transaction completes
        txn.oncomplete = function () {
            db.close();
        };
    }
}

</script>

<template>
  <div>
    <h1>Upload Images</h1>
    <div>
      data:{{ data.uploaded }}
      <br>
      <input id="upload" accept="image/png,image/jpg" class="file-upload" name="File" type="file" @change="uploadFile" >
    </div>
  </div>
</template>

<style scoped>

</style>

