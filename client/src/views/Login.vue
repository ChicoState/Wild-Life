<script lang="ts" setup>

</script>

<template>
  <div class="d-flex justify-content-center">
    <div class="col-12 col-md-8 col-lg-5">
      <h2 class="header">Login</h2>
      <div class="card d-flex flex-column gap-4">
        <div class="d-flex flex-column">
          <label for="email">Email</label>
          <input id="email" class="text-input" type="text"/>
        </div>
        <div class="d-flex flex-column">
          <label for="email">Password</label>
          <input id="password" class="text-input" type="password"/>
        </div>
        <div class="d-flex justify-content-between mt-1 align-items-center">
          <a class="text-accent" href="/register">register</a>
          <div class="button" style="background-color: #40826D;">Login</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.header {
  margin-bottom: 0.5rem;
  margin-top: 4rem;
}
</style>