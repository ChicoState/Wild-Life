<script setup>
import {reactive} from 'vue';

import {Upload} from '../rest.js';

const data = reactive({
  uploaded: {
    status: 0,
    size: 0,
    name: 0,
  },
  img: "",
  response: {}
});

function uploadSucceeded(res) {
  data.response = res
  data.img = res.data.data.token
}

function uploadFailed(error) {
  data.response = error
}

function uploadFile(event) {

  let upload = new Upload()

  upload.addFile(event.target.files[0])

  upload.onsuccess = uploadSucceeded
  upload.onfail = uploadFailed

  upload.submit()
}

</script>

<template>
  <div>
    <h2>Upload Images</h2>
    <div>
      <label class="custom-file-upload button">
        <input id="upload" accept="image/png,image/jpg" class="button" type="file" @change="uploadFile">
        <i class="fa fa-cloud-upload"></i> Select File
      </label>
      <div v-if="data.img !== ''">
        <img :src="`data:image/png;base64,${data.img}`" alt="img" style="width: 75%; padding: 1rem;"/>
      </div>

    </div>
  </div>
</template>

<style scoped>
input[type="file"] {
  display: none;
}


</style>
