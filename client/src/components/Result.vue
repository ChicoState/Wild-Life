<script lang="ts" setup>
import {inject, reactive} from 'vue';
import Loading from "./Loading.vue";
import {Buffer} from "buffer";

let cache: any = inject('cache')

const state = reactive({});

interface FileUpload {
  name: string
  date: Date
  type: string
  size: string
  token: string
}

interface Score {
  detections: number
  buckets: number
  confidence: number
  prognosis: string
}

interface Results {
  token: string
  file: FileUpload
  thumbnail: Buffer
  score: Score
}

let props = defineProps<Results>()


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
            <div class="value-pair">done</div>
          </div>
          <div class="subtitle mt-4">Confidence</div>
          <div class="pair">
            <div class="key-pair">detections</div>
            <div class="value-pair">{{ state.response.confidence }}</div>
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
  border-bottom: 1px solid rgba(255, 255, 255, 0.125);
  padding: 0.5rem 0.5rem;
}

.pair:nth-last-of-type(1) {
  /*border-bottom: 1px solid rgba(255, 255, 255, 0);*/
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
