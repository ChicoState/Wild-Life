<script lang="ts" setup>
import Loading from "@/components/Loading.vue"
import {UploadState} from "../upload";


interface Result {
  upload: {
    name: string
    type: string
    size: number
  }
  thumbnails: {
    original: string
    processed: string
  }
}


interface Results {
  response: {
    name: string,
    size: number,
    type: string,
    thumbnail: string,
    threshold: string,
    highlight: string,
    results: string,
    confidence: string,
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

</script>

<template>
  <div class="d-flex flex-column gap-3">
    <div class="d-flex flex-row gap-3 mt-2">
      <div class="card image d-flex flex-column">
        <div class="d-flex">
          <div :style="`background-image: url('data:image/jpg;base64,${props.response.thumbnail}');`"
               class="image-preview"></div>
          <div class="d-flex flex-column px-3">
            <h3>
              <i class="fa-solid fa-circle-check text-accent label-c5"></i>&nbsp;&nbsp;
              {{ props.response.progress[props.response.progress.length - 1].state }}
            </h3>
            <div class="label-o3">
              Our scans concluded that there are no obvious signs of poison oak, poison ivy, or similar irritants.
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="d-flex flex-row gap-3">
      <div :style="`background-image: url('data:image/jpg;base64,${props.response.results}');`"
           class="preview-upload gap-1">
        <div v-if="props.response.results === ''"
             class="d-flex justify-content-center align-items-center align-content-center flex-column gap-1"
             style="height: 100%;">Computing
          <Loading></Loading>
        </div>
      </div>

      <div class="d-flex flex-lg-row flex-column flex-wrap gap-4">

        <div class="sidebar">
          <div class="subtitle">Processing</div>
          <div v-for="(state) in props.response.progress" class="pair">
            <div class="key-pair">{{ state.state }}</div>
            <div class="value-pair">done</div>
          </div>
          <div class="subtitle mt-4">Confidence</div>
          <div class="pair">
            <div class="key-pair">Irritants Found</div>
            <div class="value-pair">{{ props.response.confidence }}%</div>
          </div>
          <div class="pair">
            <div class="key-pair">name</div>
            <div class="value-pair">{{ props.response.name }}</div>
          </div>
          <div class="pair">
            <div class="key-pair">type</div>
            <div class="value-pair">{{ props.response.type }}</div>
          </div>
          <div class="pair">
            <div class="key-pair">type</div>
            <div class="value-pair">{{ Math.round(props.response.size / 1000 / 1000 * 100) / 100 }}MB</div>
          </div>
        </div>
      </div>
    </div>
  </div>

</template>
<style lang="scss" scoped>

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

  aspect-ratio: 4/3;
  width: 100%;
  background-size: contain;
  background-position: center;
  background-repeat: no-repeat;
  background-color: transparent;
  border: 1px solid rgba(255, 255, 255, 0.2);
  padding: 1rem;
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
  max-width: 100%;
  font-weight: 400;
  font-size: 0.9rem;
  font-family: "Roboto", serif;
  overflow: clip !important;
  text-overflow: ellipsis !important;
  color: rgba(255, 255, 255, 0.4);
}
</style>