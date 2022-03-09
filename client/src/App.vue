<script setup>
import Upload from './components/Upload.vue'
import Logo from './components/Logo.vue'
import History from "./components/History.vue";
import {provide, reactive} from "vue";

const cache = reactive({
  history: []
})

provide('cache', cache)

</script>

<template>

  <div class="container">
    <div class="d-flex flex-column justify-content-between">
      <div class="d-flex flex-row justify-content-between align-items-center">
        <Logo></Logo>
        <div>Login</div>
      </div>
      <div class="flex-grow-1">
        <Upload></Upload>
      </div>
      <div>
        <History></History>
      </div>
    </div>
  </div>
</template>

<style>

</style>
