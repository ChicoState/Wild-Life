<script setup>
import Upload from './components/Upload.vue'
import Logo from './components/Logo.vue'
</script>

<template>
  <div class="media">
    <div class="container">
      <div>
        <Logo></Logo>
      </div>
      <div>
        <Upload></Upload>
      </div>
    </div>
  </div>
</template>

<style>
.media {
  width: 100%;

  display: flex;
  justify-content: center;
  flex-direction: row;

}
</style>
