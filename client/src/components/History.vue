<script lang="ts" setup>

import {inject, onMounted} from "vue";
import type {fileType} from '../types';
import {clearFiles, getAllFiles} from '../indexedDB';

const cache:any = inject('cache')

function clear() {
  cache.history = []
  clearFiles()
}

// Maintains vue dom syntax
onMounted(() => {
  // Load images from indexedDB into cache
  getAllFiles().then(function (result: any) {
    console.log(result)
    result.forEach((file:fileType) => {
      cache.history.push({
          image: `<div class = "row">
          <div class = "col"><img style="max-width: 100%; max-height: 100px; object-fit: contain; border-radius: 0.5rem;" class="frame" src="data:${file.type};base64,${file.data}" alt=${file.name} /></div>
          <div class = "col" style = "line-height: 1px; font-size:90%; text-align:center">
            <h2 style = "color: rgba(255, 255, 255, 0.8);"><i style = "color: yellow;" class="fa-solid fa-triangle-exclamation"></i> Possible Irritants</h2>
            <h4 style = "color: rgba(255, 255, 255, 0.8);"><i style = "text-align:right; " class="fa-solid fa-magnifying-glass"></i>    ${file.plant}</h4>
            <h4 style = "color: rgba(255,255,255,0.4);"> ${file.confidence} confidence </h4>
          </div> 
          </div>`
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

  border-radius: 0.5rem;
  padding: 0.5rem;
  background-color: rgba(44, 44, 46, 1);
  box-shadow: 0 0 8px 1px rgba(0, 0, 0, 0.05);
 text-align:left;
}

.image-grid {
  display: grid;
  align-items: center;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 1rem;
}

.clear_btn {
  float: right;
  color: rgb(108, 194, 2);
  text-decoration: none;
}

</style>
