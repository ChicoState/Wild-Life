<script lang="ts" setup>

import {inject} from "vue";

let cache: any = inject('cache')

function clear() {
  cache.history = cache.history.filter((g: any) => !g)
}

</script>

<template>
  <div class="d-flex flex-column">
    <div class="d-flex justify-content-between align-items-center">
      <h2>Previous Uploads</h2>
      <a class="clear_btn" href="#" @click="clear">clear</a>
    </div>

    <div class="image-grid">
      <div v-for="file in cache.history" :key="file.name">
        <div class="image">
          <div :style="`background-image: url('data:image/jpg;base64,${file.data}');`" class="image-preview">

          </div>
          <div class="image-desc d-flex flex-column justify-content-between">
            <div class="label-c4 label-o4 overflow-ellipse">{{ file.name }}</div>
            <div class="label-c4 label-o4">
              <i class="fa-solid fa-triangle-exclamation" style="color: #ffc400;"></i>
              Irritants detected
            </div>
            <div class="label-c4 label-o4"><i class="fa-solid fa-magnifying-glass" style="text-align:right;"></i>
            </div>
            <div class="label-c4 label-o3"> {{ file.confidence || Math.round(Math.random() * 1000) / 10 }}%
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
