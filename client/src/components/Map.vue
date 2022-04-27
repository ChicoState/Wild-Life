<script lang="ts" setup>
import { onMounted, reactive } from 'vue'
import { string } from 'yargs';
import { LongLat } from '../map';

let state = reactive({
    center: [40,40],
    projection: 'EPSG:4326',
    zoom: 8,
    rotation: 0
})

function success(pos: any){
    console.log(pos.coords)
    state.center[0] = pos.coords.lng
    state.center[1] = pos.coords.lat
    //update ol-view
}
const error = (err: any) => {
    console.warn(`ERROR(${err.code}): ${err.message}`);
}

onMounted(() => {
    let options = {
        enableHighAccuracy: true,
        timeout: 5000,
        maximumAge: 0
    };

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
            <h2>Location Data {{ state.center[0] }} {{ state.center[1] }}</h2>
        </div>
        <ol-map :loadTilesWhileAnimating="true" :loadTilesWhileInteracting="true" style="height:70vh" class="d-flex justify-content-between align-items-center">
        <ol-view ref="view" :center="state.center" :rotation="state.rotation" :zoom="state.zoom" :projection="state.projection" />

        <ol-tile-layer>
            <ol-source-osm />
        </ol-tile-layer>
        
        </ol-map>
    </div>
</template>