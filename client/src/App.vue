<script setup lang="ts">
import Logo from "@/components/Logo.vue"
import {provide, reactive} from "vue";

import router from "./router";

let state = reactive({
  history: []
})

provide('cache', state)

</script>

<template>
  <div class="container">
    <div class="d-flex flex-column justify-content-between">
      <div class="d-flex flex-row justify-content-between align-items-center">
        <Logo></Logo>
        <a v-if="router.currentRoute.value.fullPath !== '/login'" class="text-accent" href="/login">Login</a>
        <a v-else class="text-accent" href="/register">register</a>
      </div>
      <router-view/>
    </div>
  </div>

</template>

<style lang="scss">

</style>
