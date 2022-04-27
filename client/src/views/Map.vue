<script lang="ts" setup>
import { onMounted, reactive } from 'vue'
import { string } from 'yargs';
import { LongLat } from '../map';

let state = reactive({
    center: [-121,49],
    projection: 'EPSG:4326',
    zoom: 8,
    rotation: 0
})

function success(pos: GeolocationPosition){
    state.center[0] = pos.coords.longitude
    state.center[1] = pos.coords.latitude
    //update ol-view
}
const error = (err: any) => {
    console.warn(`ERROR(${err.code}): ${err.message}`);
}
//Gets random range of numbers for area, can be used later to get user data.
const getRandomInRange = (from, to, fixed) => {
    return (Math.random() * (to - from) + from).toFixed(fixed) * 1;
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
    <div class="d-flex  flex-column">
        <div class="d-flex justify-content-between align-items-center">
            <h2 class="my-1">Plant Map {{ state.center[0] }} {{ state.center[1] }}</h2>
        </div>
        <ol-map :loadTilesWhileAnimating="true" :loadTilesWhileInteracting="true" style="height:70vh" class="d-flex justify-content-between align-items-center fixed-top rounded card">
       <ol-view ref="view" :center="state.center" :rotation="state.rotation" :zoom="state.zoom" :projection="state.projection" />
        <ol-tile-layer>
            <ol-source-osm />
        </ol-tile-layer>
        <ol-vector-layer>
            <ol-source-cluster :distance="40">
                <ol-source-vector>
                    <ol-feature v-for="index in 10" :key="index">
                        <ol-geom-point :coordinates="[-121.837479, 39.728493]"></ol-geom-point>
                    </ol-feature>
                </ol-source-vector>
            </ol-source-cluster>
        <ol-style>
            <ol-style-fill color="rgba(255,255,255,1)"></ol-style-fill>
            <ol-style-circle :radius="10">
                <ol-style-fill color="rgba(51, 153, 204, 0.6)"></ol-style-fill>
                <ol-style-stroke color="rgba(255,255,255,1)" :width="1"></ol-style-stroke>
            </ol-style-circle>
        </ol-style>
        </ol-vector-layer>
        </ol-map>
    </div>
</template>
<style lang="scss" scoped>

</style>