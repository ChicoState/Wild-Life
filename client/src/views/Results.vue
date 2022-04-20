<script lang="ts" setup>
import Loading from "@/components/Loading.vue"
import {UploadState} from "../upload";
import {reactive} from "vue";
import type {Detection} from "../types";

interface Result {
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
}


interface Results {
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

const target = [
  {
    state: "uploaded",
  },
  {
    state: "queued",
  },
  {
    state: "processing",
  },
  {
    state: "analyzing",
  },
  {
    state: "compiling",
  },
  {
    state: "results",
  }
]

let props = defineProps<Results>()

let state = reactive({})

</script>

<template>
  <div class="d-flex flex-column gap-3">
    <div class="d-flex flex-row gap-3 mt-2">
      <div class="card image d-flex flex-column">
        <div class="d-flex">
          <div :style="`background-image: url('data:image/jpg;base64,${props.response.thumbnail}');`"
               class="image-preview"></div>
          <div class="d-flex flex-column px-3">
            <h3 class="mb-0">
              <span v-if="props.response.detections.length === 0">
                <i class="fa-solid fa-circle-check text-accent label-c5 fa-fw"></i>&nbsp;&nbsp;Low irritation
              </span>
              <span v-else>
                <i class="fa-solid fa-triangle-exclamation text-orange label-c5"></i>&nbsp;&nbsp;{{ props.response.detections.length }} Irritants Detected
              </span>

            </h3>
            <div>
              <div v-if="props.response.detections.length === 0 " class="label-o3 label-c4">
                Our scans concluded that there are no obvious signs of poison oak, poison ivy, or similar irritants.
              </div>
              <div v-else class="label-o3 label-c4">
                We have identified multiple possible irritants.
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="d-flex flex-row gap-3">

      <div class="flex-grow-1">
        <img v-if="props.response.results !== ''" :src="`data:image/jpg;base64,${props.response.results}`" alt="results"
             class="preview-upload"/>
        <div v-else
             class="d-flex card justify-content-center align-items-center align-content-center flex-column gap-1"
             style="height: 100%;">Computing
          <Loading></Loading>
        </div>
      </div>

      <div class="d-flex flex-lg-row flex-column flex-wrap gap-4">

        <div class="sidebar">

          <div class="subtitle">Detections</div>
          <div class="shadow-box">
            <div v-for="d in props.response.detections.sort((a, b) => b.confidence - a.confidence).slice(0, 4)">
              <div class="card py-2 px-2 my-1 d-flex justify-content-between ">
                <div class="label-c4">{{ d.type }}</div>
                <div class="label-c4 label-o3">{{ Math.round(d.confidence * 1000) / 10 }}%</div>
              </div>
            </div>
            <div class="d-flex justify-content-between">
              <div></div>
              <div v-if="props.response.detections.length > 4" class="text-accent label-c3 px-2">+
                {{ props.response.detections.length - 4 }} more
              </div>
            </div>
          </div>
          <div class="pair">
            <div class="value-pair">{{ props.response.confidence }}</div>
          </div>
          <div class="subtitle mt-3">Upload</div>
          <div class="pair">
            <div class="key-pair">name</div>
            <div class="value-pair">{{ props.response.name }}</div>
          </div>
          <div class="pair">
            <div class="key-pair">type</div>
            <div class="value-pair">{{ props.response.type }}</div>
          </div>
          <div class="pair">
            <div class="key-pair">size</div>
            <div class="value-pair">{{ Math.round(props.response.size / 1000 / 1000 * 100) / 100 }}MB</div>
          </div>
        </div>
      </div>
    </div>
  </div>

</template>
<style lang="scss" scoped>
.shadow-box {

}

.sidebar {
  padding: 1rem;
  min-width: 18rem;
  border-radius: 0.5rem;
  backdrop-filter: blur(48px);
  background-color: rgba(0, 0, 0, 0.012);
  border: 1px solid rgba(255, 255, 255, 0.1)
}

.result-grid {
  display: grid;
  flex-grow: 1;
  aspect-ratio: 1/1 !important;
  grid-gap: 1rem;
  grid-template-rows: repeat(2, 1fr);
  grid-template-columns: repeat(2, 1fr);
}


.thumbnail {
  background-size: cover;
  background-position: center;
  outline: 1px solid rgba(255, 255, 255, 0.2);
  width: 10rem !important;
  height: 10rem;
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


.preview-upload {

  border-radius: 0.5rem;

  width: 100%;
  background-size: contain;
  background-position: center;
  background-repeat: no-repeat;
  background-color: transparent;
  border: 1px solid rgba(255, 255, 255, 0.1);
  animation: loadIn 200ms ease-out forwards;
}

@keyframes loadIn {
  0% {
    opacity: 0.7;
    transform: scale(0.98);
  }
  100% {
    opacity: 1;
    transform: scale(1);
  }
}

.pair {
  display: flex;
  justify-content: space-between;
  align-items: center;

  padding: 0.2rem 0.2rem;
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
  max-width: 100%;
  font-weight: 400;
  font-size: 0.9rem;

  overflow: clip !important;
  text-overflow: ellipsis !important;
  color: rgba(255, 255, 255, 0.4);
}
</style>