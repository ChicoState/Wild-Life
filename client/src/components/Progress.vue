<script lang="ts" setup>

import moment from "moment";

interface TaskState {
  state: string
  time?: string
}

interface Progress {
  progress: TaskState[]
  final: TaskState[]
}

let props = defineProps<Progress>()

function getDelta(state: string) {
  let pot = props.progress.find((p: TaskState) => p.state === state)

  let before = props.progress.find((p: TaskState) => p.state === state)
  if (!pot) return ""
  return moment(pot.time).subtract(props.progress[0].time)
}

</script>

<template>
  <div class="progress-stepped">
    <div v-for="step in props.final" :key="step.state"
         :class="`${props.final.filter((p: TaskState) => props.progress.filter(q => q.state === p.state))?'completed':'active'}`"
         class="progress-stepped-item">
      <div class="step-counter"><i v-if="props.progress.find((p: TaskState) => p.state === step.state)"
                                   class="fa-solid fa-check"></i></div>
      <div class="step-name">{{ step.state }}</div>

    </div>
  </div>

</template>

<style lang="scss" scoped>
.progress-stepped {
  margin-top: auto;
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

.progress-stepped-item {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;

  @media (max-width: 768px) {
    font-size: 12px;
  }

}

.progress-stepped-item::before {
  position: absolute;
  content: "";
  border-bottom: 2px solid #ccc;
  width: 100%;
  top: 20px;
  left: -50%;
  z-index: 2;
}

.progress-stepped-item::after {
  position: absolute;
  content: "";
  border-bottom: 2px solid #ccc;
  width: 100%;
  top: 20px;
  left: 50%;
  z-index: 2;
}

.progress-stepped-item .step-counter {
  position: relative;
  z-index: 5;
  display: flex;
  justify-content: center;
  align-items: center;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #ccc;
  margin-bottom: 6px;
}

.progress-stepped-item.active {
  font-weight: bold;
}

.progress-stepped-item.completed .step-counter {
  background-color: #4bb543;
}

.progress-stepped-item.completed::after {
  position: absolute;
  content: "";
  border-bottom: 2px solid #4bb543;
  width: 100%;
  top: 20px;
  left: 50%;
  z-index: 3;
}

.progress-stepped-item:first-child::before {
  content: none;
}

.progress-stepped-item:last-child::after {
  content: none;
}
</style>