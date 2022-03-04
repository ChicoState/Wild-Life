<script setup>
import {reactive} from 'vue';

import {Upload} from '../rest.js';

const data = reactive({
  uploaded: {
    status: 0,
    size: 0,
    name: 0,
  },
  response: {}
});

function uploadSucceeded(res) {
  data.response = res
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
      {{ data.response }}
    </div>
  </div>
</template>

<style scoped>
input[type="file"] {
  display: none;
}


</style>

