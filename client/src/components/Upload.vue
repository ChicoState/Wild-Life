<script setup lang="ts">
import {inject, reactive} from 'vue';
import {Upload, UploadState} from '../rest';
import type {fileType} from '../types'
import Results from "../views/Results.vue";

let state = reactive<{
  upload: Upload,
  response: any,
  context: boolean,
}>({
  upload: {} as Upload,
  response: {
    name: "",
    size: "",
    thumbnail: "",
    threshold: "",
    highlight: "",
    results: "",
    confidence: "",
    progress: [] as UploadState[],
    token: "",
  },
  context: false,
});


let cache: any = inject('cache')

function updateStatus(up: UploadState) {
  state.response.progress.push(up)
  switch (up.state) {
    case "thumbnail":
      state.response.thumbnail = up.data
      state.context = true
      cache.history.push({
        data: up.data,
        name: state.response.name,
        type: state.response.type,
      })
      break
    case "threshold":
      state.response.threshold = up.data
      state.context = true
      break
    case "highlight":
      state.response.highlight = up.data
      state.context = true
      break
    case "results":
      state.response.results = up.data
      state.response.confidence = up.message
      state.context = true
      break
  }
}

function reset() {
  state.response = {
    name: "",
    size: "",
    thumbnail: "",
    threshold: "",
    highlight: "",
    results: "",
    confidence: "",
    progress: [] as UploadState[],
    token: "",
  }
}

function uploadFile(event: any) {
  reset()
  //erase placeholder
  let file: File = event.target.files[0]
  let temp_uploaded: fileType = {
    name: file.name,
    size: file.size,
    type: file.type,
    id: "",
    data: ""
  }

  if (temp_uploaded.type != 'image/png' && temp_uploaded.type != 'image/jpeg') {
    alert('Only PNG and JPG files are allowed')
    return
  } else {
    const today = new Date();
    temp_uploaded.id = today.getFullYear() + '-' + (today.getMonth() + 1) + '-' + today.getDate() + '+' + file.name;
    // Load images from indexedDB into cache to prevent double image insertion
    // getAllFiles().then(function (result: any) {
    //   result.forEach((file: fileType) => {
    //     if (temp_uploaded.id == file.id) {
    //       //Reload to stop image from showing up until page refresh
    //       return
    //     }
    //   });
  }

    state.upload = new Upload()

  state.upload.addFile(file)
  state.upload.submit()
  state.response.name = file.name
  state.response.type = file.type
  state.response.size = file.size
  state.upload.update = updateStatus

  // let buf = new Promise(getBuffer(file))

  // buf.then(function (data: any) {
  //   temp_uploaded.data = data
  //   addFile(temp_uploaded)
  // }).catch(function (error) {
  //   console.log("Error: ", error)
  // })
}


</script>

<template>
  <div v-if="state.context" class="results-context" @mousedown="state.context = false">
    <Results :response="state.response"></Results>
  </div>
  <div id="content-mobile">
    <div class="d-flex flex-column">
      <h2>Upload Images</h2>
      <div class="d-flex justify-content-between gap-0">
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



.results-context {
  position: absolute;
  top: 0;
  left: 0;
  height: 100vh;
  width: 100vw;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  background-color: rgba(0, 0, 0, 0.2);

  backdrop-filter: blur(10px);
}

input[type="file"] {
  display: none;
}


.gap-1 {
  gap: 0.25rem;
}

.child {
  display: inline-block;

  padding: 2px;

}

</style>
