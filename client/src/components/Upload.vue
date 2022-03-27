<script setup lang="ts">
import {inject, reactive} from 'vue';
import {Upload} from '../rest';
import {addFile, getAllFiles} from '../indexedDB'
import {getBuffer} from '../upload'
import type {fileType} from '../types'

const cache: any = inject('cache')

const state = reactive<{
  uploaded: fileType,
  response: any,
  context: boolean,
  imgs: any[]
}>({
  uploaded: {
    id: '',
    name: '',
    type: '',
    size: 0,
    data: '',
  },
  response: {
    image: {
      data: "",
      type: ""
    },
    token: "",
    resolve: {}
  },
  context: false,
  imgs: []
});

function beginProcessing(res: any) {
  state.response.image.data = res.data.data.token
  state.response.image.type = "image/png"
}

function uploadFailure(data: any) {
  console.log("Rock", data)
}

function uploadFile(event: any) {

  //erase placeholder

  let file: File = event.target.files[0]
  let temp_uploaded: fileType = {
    name: "",
    size: 0,
    type: file.type,
    id: "",
    data: ""
  }
  let upload = new Upload()
  upload.addFile(file)
  const uploaded = upload.submit()
  uploaded.then(beginProcessing).catch(beginProcessing)
  if (temp_uploaded.type != 'image/png' && temp_uploaded.type != 'image/jpeg') {
    alert('Only PNG and JPG files are allowed')
    return
  } else {
    const today = new Date();
    temp_uploaded.id = today.getFullYear() + '-' + (today.getMonth() + 1) + '-' + today.getDate() + '+' + file.name;
    // Load images from indexedDB into cache to prevent double image insertion
    getAllFiles().then(function (result: any) {
      result.forEach((file: fileType) => {
        if (temp_uploaded.id == file.id) {
          //Reload to stop image from showing up until page refresh
          window.location.reload()
          alert("You cannot re-upload the same file!")
          return
        }
      });
    })

    temp_uploaded.name = file.name
    temp_uploaded.size = file.size

    let buf = new Promise(getBuffer(file))

    buf.then(function (data: any) {
      temp_uploaded.data = data
      addFile(temp_uploaded)
      state.context = true

    }).catch(function (error) {
      console.log("Error: ", error)
    })
  }
}


</script>

<template>
  <div v-if="state.context" class="results-context" @click="state.context = !state.context">
    <div><h1>Results</h1></div>
    <div class="d-flex gap-1">
      <div :style="`background-image:  url('data:${state.response.image.type};base64,${state.response.image.data}');`"
           class="preview-upload"></div>
      <div class="sidebar">
        Sidebar
      </div>
    </div>
  </div>
  <div id="content-mobile">
    <div class="d-flex flex-column ">
      <h2>Upload Images</h2>
      <div class="d-flex justify-content-between gap-1">
        <div class="flex-shrink-0">
          <label class="custom-file-upload button">
            <input id="camera" accept="image/png,image/jpeg" capture="user" class="button" type="file"
                   @change="uploadFile">
            <i class="fa-solid fa-camera" style="text-align: center;"></i>
          </label>
        </div>
        <div class="flex-grow-1">
          <label class="custom-file-upload button">
            <input id="upload" accept="image/png,image/jpeg" class="button" type="file" @change="uploadFile">
            <i class="fa fa-cloud-upload"></i> Select File
          </label>
        </div>

      </div>

    </div>
  </div>

</template>

<style scoped>
.sidebar {
  padding: 2rem;
  border-radius: 0.5rem;
  background-color: rgba(64, 64, 64, 0.2);
}

.results-context {
  position: absolute;
  top: 0;
  left: 0;
  height: 100vh;
  width: 100vw;
  padding: 2rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  background-color: rgba(64, 64, 64, 0.2);
  backdrop-filter: blur(12px);
}

input[type="file"] {
  display: none;
}

.preview-upload {
  width: 70%;
  aspect-ratio: 6/4;
  border-radius: 0.5rem;

  background-size: contain;
  background-position: center;
  background-repeat: no-repeat;
  background-color: rgba(255, 255, 255, 0.05);
}

.gap-1 {
  gap: 0.25rem;
}

.child {
  display: inline-block;

  padding: 2px;

}

</style>
