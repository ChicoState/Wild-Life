<script setup lang="ts">
import Logo from "@/components/Logo.vue"
import {onMounted, provide, reactive, watch} from "vue";
import {UploadResult} from "./upload";

import router from "./router";


interface App {
  history: UploadResult[]
}

let state = reactive<App>({
  history: [] as UploadResult[]
})

onMounted(() => {
  let local = localStorage.getItem("cache")
  if (!local) return
  state.history = JSON.parse(local).history as UploadResult[]
})

watch(state, (recent: any, old: any) => {
  if (recent !== [])
    localStorage.setItem("cache", JSON.stringify(recent))
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
