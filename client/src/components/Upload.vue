<script setup lang="ts">
import {inject, reactive} from 'vue';
import {Upload} from '../rest';
import {addFile} from '../indexedDB'
import {getBuffer} from '../upload'
import type {fileType} from '../types'

const cache:any = inject('cache')

const data = reactive<{
  uploaded: fileType,
  response: any,
  imgs: any[]
}>({
  uploaded: {
    id: '',
    name: '',
    type: '',
    size: 0,
    data: '',
  },
  response: {},
  imgs: []
});

function uploadFile(event: any) {

  //erase placeholder
  let upload = new Upload()
  let file: File = event.target.files[0]
  let temp_uploaded: fileType = {
    name: "",
    size: 0,
    type: file.type,
    id: "",
    data: ""
  }

  upload.addFile(file)

  upload.submit()
  if (temp_uploaded.type != 'image/png' && temp_uploaded.type != 'image/jpeg') {
    alert('Only PNG and JPG files are allowed')
    return
  } else {
    const today = new Date();
    temp_uploaded.id = today.getFullYear() + '-' + (today.getMonth() + 1) + '-' + today.getDate() + '+' + file.name;
    temp_uploaded.name = file.name
    temp_uploaded.size = file.size
    var buf = new Promise(getBuffer(file))
    buf.then(function (data: any) {
      temp_uploaded.data = data
      console.log(temp_uploaded)
      addFile(temp_uploaded)
      cache.history.push({
        image: `<img style="max-width: 100%; max-height: 25rem; object-fit: contain;" class="frame" src="data:${temp_uploaded.type};base64,${temp_uploaded.data}" alt=${temp_uploaded.name} />`
      });
    }).catch(function (error) {
      console.log("Error: ", error)
    })
  }
}

</script>

<template>
  <div class="d-flex flex-column">
    <div class="d-flex flex-column">
      <h2>Upload Images</h2>
      <label class="custom-file-upload button">
        <input id="upload" accept="image/png,image/jpeg" class="button" type="file" @change="uploadFile">
        <i class="fa fa-cloud-upload"></i> Select File
      </label>
      {{ data }}
    </div>
  </div>
</template>

<style scoped>

input[type="file"] {
  display: none;
}

</style>
