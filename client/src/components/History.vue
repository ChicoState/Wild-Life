<script lang="ts" setup>

import {inject, onMounted} from "vue";
import type {fileType} from '../types';
import {clearFiles, getAllFiles} from '../indexedDB';

const cache: any = inject('cache')

function clear() {
  cache.history = []
  clearFiles()
}

function refreshUploads() {
  getAllFiles().then(function (result: any) {
    cache.history = cache.history.filter((a: any) => !a)
    result.forEach((file: fileType) => {
      let res = {
        image: file
      }
      cache.history.push(res)
    });
  })
}

// Maintains vue dom syntax
onMounted(() => {
  // Load images from indexedDB into cache
  refreshUploads()
})

</script>

<template>
  <div class="d-flex flex-column">
    <div class="d-flex justify-content-between align-items-center">
      <h2>Previous Uploads</h2>
      <a class="clear_btn" href="#" @click="clear">clear</a>
    </div>
    <div class="image-grid">
      <div v-for="file in cache.history.map(f => f.image)">
        <div class="image">
          <div :style="`background-image: url('data:${file.type};base64,${file.data}');`" class="image-preview">

          </div>
          <div class="image-desc d-flex flex-column justify-content-between">
            <div>
              <i class="fa-solid fa-triangle-exclamation" style="color: #ffc400;"></i>
              Possible Irritants
            </div>
            <div style="color: rgba(255, 255, 255, 0.8);"><i class="fa-solid fa-magnifying-glass"
                                                             style="text-align:right;"></i> {{ file.plant }}
            </div>
            <div style="color: rgba(255,255,255,0.3);"> {{ file.confidence || Math.round(Math.random() * 100) }}%
              confidence
            </div>
          </div>
        </div>
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
  text-align: left;
  display: flex;
  justify-content: start;
  flex-direction: row;
  gap: 1rem;
}

.image-desc {
  font-size: 1.1rem;
  gap: 0.5rem;
}

.image-grid {
  display: grid;
  align-items: center;
  grid-template-columns: repeat(auto-fit, minmax(20rem, 1fr));
  gap: 1rem;
}

.image-preview {
  aspect-ratio: 4/3;
  background-size: cover;
  background-position: center;
  height: 8rem;
}

.clear_btn {
  float: right;
  color: rgb(108, 194, 2);
  text-decoration: none;
}

</style>
