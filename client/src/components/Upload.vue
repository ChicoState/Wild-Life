<script setup>
import {reactive} from 'vue';

const data = reactive({
  uploaded: {
    size: 0
  },
});

function uploadFile(event) {
  data.uploaded.size = event.target.size
  console.log(event)
}

</script>

<template>
  <div>
    <h1>Upload Images</h1>
    <div>
      Size:{{ data.uploaded }}
      <br>
      <input id="upload" accept="image/png,image/jpg" class="file-upload" name="File" type="file" @change="uploadFile">
    </div>
  </div>
</template>

<style scoped>

</style>
