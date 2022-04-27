<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { LongLat } from '../map';
let center = ref([-121.837479, 39.728493])
const projection = ref('EPSG:4326')
const zoom = ref(8)
const rotation = ref(0)
onMounted(() => {
    let options = {
        enableHighAccuracy: true,
        timeout: 5000,
        maximumAge: 0
    };
    const success = (pos: any) => {
        console.log(pos.coords)
        center = ref([pos.coords.lng, pos.coords.lat])
    }
    const error = (err: any) => {
        console.warn(`ERROR(${err.code}): ${err.message}`);
    }
    if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(success,error,options);
    }
    else{
        console.log("Geolocation is not supported by this browser.");
    }
})
</script>
<template>
    <div class="d-flex flex-column">
        <div class="d-flex justify-content-between align-items-center">
            <h2>Location Data {{ center[0] }} {{ center[1] }}</h2>
        </div>
        <ol-map :loadTilesWhileAnimating="true" :loadTilesWhileInteracting="true" style="height:70vh" class="d-flex justify-content-between align-items-center">
        <ol-view ref="view" :center="center" :rotation="rotation" :zoom="zoom" :projection="projection" />

        <ol-tile-layer>
            <ol-source-osm />
        </ol-tile-layer>
        
        </ol-map>
    </div>
</template>