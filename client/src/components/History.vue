<script lang="ts" setup>

import {inject} from "vue";
import type {Cache, UploadResult} from "../upload";

let cache: Cache = inject('cache') || {} as Cache

function clear() {
  cache.history = cache.history.filter((g: UploadResult) => !g) as UploadResult[]
}

</script>

<template>
  <div class="d-flex flex-column">
    <div class="d-flex justify-content-between align-items-center">
      <h2>Previous Uploads</h2>
      <a class="clear_btn" href="#" @click="clear">clear</a>
    </div>
    <div class="image-grid">
      <div v-for="file in cache.history">
        <div v-if="file" class="image">
          <div :style="`background-image: url('data:image/jpg;base64,${file.data}');`" class="image-preview">

          </div>
          <div class="image-desc d-flex flex-column justify-content-between">
            <div class="label-c4 label-o4 overflow-ellipse">{{ file.name }}</div>
            <div v-if="file.confidence > 0.7">
              <div class="label-c4 label-o4">
                <i class="fa-solid fa-triangle-exclamation" style="color: #ffc400;"></i>
                Irritants detected
              </div>
            </div>
            <div v-else>
              <div class="label-c4 label-o4">
                <i class="fa-solid fa-check" style="color: #00c853;"></i>
                No irritants
              </div>
            </div>
            <div class="label-c4 label-o4"><i class="fa-solid fa-magnifying-glass" style="text-align:right;"></i>
              {{ file.result }}
            </div>
            <div class="label-c4 label-o3"> {{ ((file.confidence * 100) || 0).toPrecision(4) }}%
              confidence
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
.image {
  border-radius: 0.5rem;
  padding: 0.5rem;
  background-color: rgba(44, 44, 46, 1);
  box-shadow: 0 0 8px 1px rgba(0, 0, 0, 0.05);
  text-align: left;
  display: flex;
  justify-content: start;
  flex-direction: row;
  gap: 0.5rem;
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
  border-radius: 0.25rem;
  height: 7.5rem;
}

.clear_btn {
  float: right;
  color: rgb(108, 194, 2);
  text-decoration: none;
}

</style>
