<script setup lang="ts">
import {inject, reactive} from 'vue';
import moment from 'moment';
import {Upload, UploadState} from '../rest';
import {addFile, getAllFiles} from '../indexedDB'
import {getBuffer} from '../upload'
import type {fileType} from '../types'
import Loading from "./Loading.vue";

const cache: any = inject('cache')


const state = reactive<{
  uploaded: fileType,
  upload: Upload,
  response: any,
  context: boolean,
  socket: any,
  imgs: any[]
}>({
  uploaded: {
    id: '',
    name: '',
    type: '',
    size: 0,
    data: '',
  },
  upload: {} as Upload,
  response: {
    thumbnail: "",
    threshold: "",
    highlight: "",
    results: "",
    progress: [] as UploadState[],
    token: "",
    resolve: {}
  },
  socket: {
    history: [],
    state: "",
    error: "",
    progress: ""
  },
  context: false,
  imgs: []
});


function updateStatus(up: UploadState) {
  state.response.progress.push(up)
  switch (up.state) {
    case "thumbnail":
      state.response.thumbnail = up.data
      state.context = true
      addFile(up)
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
      state.context = true
      break
  }
}

function reset() {
  state.response = {
    thumbnail: "",
    threshold: "",
    highlight: "",
    results: "",
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
    getAllFiles().then(function (result: any) {
      result.forEach((file: fileType) => {
        if (temp_uploaded.id == file.id) {
          //Reload to stop image from showing up until page refresh
          return
        }
      });
    })


    state.upload = new Upload()

    state.upload.addFile(file)
    state.upload.submit()
    state.upload.update = updateStatus


    let buf = new Promise(getBuffer(file))

    buf.then(function (data: any) {
      temp_uploaded.data = data
      addFile(temp_uploaded)
      state.uploaded = temp_uploaded
    }).catch(function (error) {
      console.log("Error: ", error)
    })
  }
}


</script>

<template>

  <div v-if="state.context" class="results-context" @click="state.context = !state.context">

    <div class="container">
      <div class="title">Results</div>
      <div class="d-flex flex-lg-row flex-column flex-wrap gap-1">
        <div class="result-grid">
          <div :style="`background-image: url('data:image/jpg;base64,${state.response.thumbnail}');`"
               class="preview-upload">
            <div>Original</div>
            <div v-if="state.response.thumbnail === ''"
                 class="d-flex justify-content-center align-items-center align-content-center flex-column"
                 style="height: 100%;">Computing
              <Loading></Loading>
            </div>
          </div>
          <div :style="`background-image: url('data:image/jpg;base64,${state.response.threshold}');`"
               class="preview-upload">
            <div>Threshold + Contour</div>
            <div v-if="state.response.threshold === ''"
                 class="d-flex justify-content-center align-items-center align-content-center flex-column"
                 style="height: 100%;">Computing
              <Loading></Loading>
            </div>
          </div>
          <div :style="`background-image: url('data:image/jpg;base64,${state.response.highlight}');`"
               class="preview-upload">
            <div>Highlight</div>
            <div v-if="state.response.highlight === ''"
                 class="d-flex justify-content-center align-items-center align-content-center flex-column gap-1"
                 style="height: 100%;">Computing
              <Loading></Loading>
            </div>
          </div>
          <div :style="`background-image: url('data:image/jpg;base64,${state.response.results}');`"
               class="preview-upload">
            <div>Results</div>
            <div v-if="state.response.results === ''"
                 class="d-flex justify-content-center align-items-center align-content-center flex-column gap-1"
                 style="height: 100%;">Computing
              <Loading></Loading>
            </div>
          </div>
        </div>
        <div class="sidebar">
          <div class="subtitle">Processing</div>
          <div v-for="(state,i) in state.response.progress" class="pair">
            <div class="key-pair">{{ state.state }}</div>
            <div class="value-pair">{{ moment(state.time).format('h:mm:ss') }}</div>
          </div>
          <div class="subtitle mt-4">Upload Details</div>
          <div class="pair">
            <div class="key-pair">name</div>
            <div class="value-pair">{{ state.uploaded.name }}</div>
          </div>
          <div class="pair">
            <div class="key-pair">size</div>
            <div class="value-pair">{{ Math.round(state.uploaded.size / 1000 / 1000 * 100) / 100 }} MB</div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div id="content-mobile">
    <div class="d-flex flex-column">
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

.result-grid {
  display: grid;
  flex-grow: 1;
  aspect-ratio: 1/1 !important;
  grid-gap: 1rem;
  grid-template-rows: repeat(2, minmax(20rem, 1fr));
  grid-template-columns: repeat(2, minmax(20rem, 1fr));
}


.thumbnail {
  background-size: contain;
  background-position: center;
  outline: 1px solid rgba(255, 255, 255, 0.2);
  width: 20rem !important;
  height: 20rem;
}

.title {
  font-weight: 500;
  font-size: 2.5rem;
  padding: 0.5rem 0;
  margin-bottom: 0.5rem;
  color: rgba(255, 255, 255, 0.8);
}

.subtitle {
  font-weight: 500;
  font-size: 1rem;
  margin-bottom: 0.5rem;
  color: rgba(255, 255, 255, 0.8);
}

.pair {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid rgba(255, 255, 255, 0.25);
  padding: 0.5rem 0.5rem;
}

.pair:nth-last-of-type(1) {
  border-bottom: 1px solid rgba(255, 255, 255, 0);
}

.key-pair {
  font-weight: 400;
  font-size: 1rem;
  color: rgba(255, 255, 255, 0.8);
}

.value-pair {
  font-weight: 400;
  font-size: 0.9rem;
  font-family: "Roboto", serif;

  overflow: clip;
  text-overflow: ellipsis;
  color: rgba(255, 255, 255, 0.4);
}


.sidebar {
  padding: 1rem;
  min-width: 18rem;
  border-radius: 0.5rem;
  backdrop-filter: blur(48px);
  background-color: rgba(0, 0, 0, 0.012);
  border: 1px solid rgba(255, 255, 255, 0.1)
}

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

.preview-upload {

  border-radius: 0.5rem;

  background-size: contain;
  background-position: center;
  background-repeat: no-repeat;
  background-color: black;
  border: 1px solid rgba(255, 255, 255, 0.2);
  padding: 1rem;
}

.gap-1 {
  gap: 0.25rem;
}

.child {
  display: inline-block;

  padding: 2px;

}

</style>
