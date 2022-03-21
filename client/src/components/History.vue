<script lang="ts" setup>

import {inject, onMounted} from "vue";
import {clearFiles, getAllFiles} from '../indexedDB';

const cache:any = inject('cache')

function clear() {
  cache.history = []
  clearFiles()
}

onMounted(() => {
  // Load images from indexedDB into cache
  getAllFiles().then(function (result: any) {
    console.log(result)
    result.forEach(file => {
      cache.history.push({
          image: `<img style="max-width: 100%; max-height: 25rem; object-fit: contain;" class="frame" src="data:${file.type};base64,${file.data}" alt=${file.name} />`
      })
    });
  })
})

</script>

<template>
  <div class="d-flex flex-column">
    <div class="d-flex justify-content-between align-items-center">
      <h2>Previous Uploads</h2>
      <a class="clear_btn" href="#" @click="clear">clear</a>
    </div>
    <div class="image-grid">
      <div v-for="img in cache.history" class="image">
        <div v-html="img.image"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.image {
  width: 100%;
  height: 8rem;
  aspect-ratio: 3/4;
  border-radius: 0.5rem;
  padding: 0.5rem;
  background-color: rgba(44, 44, 46, 1);
  color: rgba(255, 255, 255, 0.4);
  box-shadow: 0 0 8px 1px rgba(0, 0, 0, 0.05);
}

.image-grid {
  display: grid;
  align-items: center;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 1rem;
}

.clear_btn {
  float: right;
  color: rgb(108, 194, 2);
  text-decoration: none;
}

</style>