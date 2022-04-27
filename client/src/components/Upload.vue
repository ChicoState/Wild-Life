<script setup lang="ts">
import {inject, reactive} from 'vue';
import {Upload, UploadState} from '../upload';
import Results from "../views/Results.vue";
import Loading from "./Loading.vue";
import {Detection} from "../types";


interface UploadProps {
  upload: Upload,
  waiting: boolean
  context: boolean,
  error: any,
  response: {
    name: string,
    size: number,
    thumbnail: string,
    type: string,
    threshold: string,
    highlight: string,
    results: string,
    confidence: string,
    detections: Detection[],
    progress: UploadState[],
    token: string,
  },
}

let state = reactive<UploadProps>({
  upload: {} as Upload,
  waiting: false,
  error: "",
  response: {
    name: "",
    size: 0,
    type: "",
    thumbnail: "",
    threshold: "",
    highlight: "",
    results: "",
    confidence: "",
    detections: [] as Detection[],
    progress: [] as UploadState[],
    token: "",
  },
  context: false,
});


let cache: any = inject('cache')

function updateError(res: any) {
  state.error = res
}

function updateStatus(up: UploadState) {
  state.response.progress.push(up)
  switch (up.state) {
    case "processing":
      state.response.thumbnail = up.data
      state.context = true
      cache.history.push({
        data: up.data,
        name: state.response.name,
        type: state.response.type,
      })
      state.waiting = false
      break
    case "results":
      state.response.results = up.data
      let proto = JSON.parse(up.message)
      if (!proto) return
      state.response.detections = proto as Detection[]
      state.context = true
      break
  }
}

function reset() {
  state.response = {
    name: "",
    size: 0,
    thumbnail: "",
    type: "",
    threshold: "",
    highlight: "",
    results: "",
    confidence: "",
    progress: [] as UploadState[],
    detections: [] as Detection[],
    token: "",
  }
}

function uploadFile(event: any) {
  state.waiting = true
  reset()

  let file: File = event.target.files[0]

  if (file.type != 'image/png' && file.type != 'image/jpeg') {
    alert('Only PNG and JPG files are allowed')
    return
  }

  state.upload = new Upload()

  state.upload.addFile(file)
  state.upload.submit()


  state.response.name = file.name
  state.response.type = file.type
  state.response.size = file.size
  state.upload.update = updateStatus
  state.upload.error = updateError
}


</script>

<template>

  <div id="content-mobile">
    <span v-if="state.error">{{ state.error }}</span>
    <div v-if="state.context" class="">
      <a class="text-accent" href="#" @click="state.context = false"><i class="fas fa-chevron-left label-o4"
                                                                        style="text-decoration: none;">&nbsp;</i>Done
      </a>
      <h2 class="my-1">Upload Status</h2>
      <Results :response="state.response"></Results>
    </div>
    <div v-else class="d-flex flex-column">
      <h2>Upload An Image</h2>
      <div class="d-flex justify-content-between gap-0">
        <div class="flex-shrink-0">
          <div class = "d-block d-sm-block d-md-block d-lg-none">
            <label class="custom-file-upload button">
              <input id="camera" accept="image/png,image/jpeg" capture="environment
              " class="button" type="file"
                    @change="uploadFile">
              <i class="fa-solid fa-camera" style="text-align: center;"></i>
            </label>
          </div>
        </div>
        <div class="flex-grow-1">
          <label class="custom-file-upload button">
            <input id="upload" accept="image/png,image/jpeg" class="button" type="file" @change="uploadFile">
            <span v-if="state.waiting" class="d-flex gap">
              <Loading size="sm"></Loading>
              <span>&nbsp;&nbsp;Uploading</span>
            </span>
            <span v-else>
              <i class="fa fa-cloud-upload"></i>
              <span>&nbsp;&nbsp;Select File</span>
            </span>
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
