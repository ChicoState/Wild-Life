<script lang="ts" setup>
import Loading from "@/components/Loading.vue"
import {reactive} from "vue";

let props = defineProps({
  response: {}
})

let state = reactive({
  response: props.response
})

</script>

<template>
  <div class="d-flex flex-column">
    <div class="title">Results</div>
    <div class="d-flex flex-lg-row flex-column flex-wrap gap-4">
      <div class="result-grid">
        <div :style="`background-image: url('data:image/jpg;base64,${props.response.thumbnail}');`"
             class="preview-upload">
          <div>Original</div>
          <div v-if="props.response.thumbnail === ''"
               class="d-flex justify-content-center align-items-center align-content-center flex-column"
               style="height: 100%;">Computing
            <Loading></Loading>
          </div>
        </div>
        <div :style="`background-image: url('data:image/jpg;base64,${props.response.threshold}');`"
             class="preview-upload">
          <div>Threshold + Contour</div>
          <div v-if="props.response.threshold === ''"
               class="d-flex justify-content-center align-items-center align-content-center flex-column"
               style="height: 100%;">Computing
            <Loading></Loading>
          </div>
        </div>
        <div :style="`background-image: url('data:image/jpg;base64,${props.response.highlight}');`"
             class="preview-upload">
          <div>Highlight</div>
          <div v-if="props.response.highlight === ''"
               class="d-flex justify-content-center align-items-center align-content-center flex-column gap-1"
               style="height: 100%;">Computing
            <Loading></Loading>
          </div>
        </div>
        <div :style="`background-image: url('data:image/jpg;base64,${props.response.results}');`"
             class="preview-upload">
          <div>Results</div>
          <div v-if="props.response.results === ''"
               class="d-flex justify-content-center align-items-center align-content-center flex-column gap-1"
               style="height: 100%;">Computing
            <Loading></Loading>
          </div>
        </div>
      </div>
      <div class="sidebar">
        <div class="subtitle">Processing</div>
        <div v-for="(state,i) in props.response.progress" class="pair">
          <div class="key-pair">{{ state.state }}</div>
          <div class="value-pair">done</div>
        </div>
        <div class="subtitle mt-4">Confidence</div>
        <div class="pair">
          <div class="key-pair">detections</div>
          <div class="value-pair">{{ props.response.confidence }}</div>
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


</template>
<style scoped>

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
  grid-template-rows: repeat(2, minmax(10rem, 1fr));
  grid-template-columns: repeat(2, minmax(10rem, 1fr));
}


.thumbnail {
  background-size: contain;
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

  background-size: contain;
  background-position: center;
  background-repeat: no-repeat;
  background-color: black;
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