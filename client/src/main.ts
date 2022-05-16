import {createApp} from 'vue'
import App from './App.vue'
import router from "./router";
import './assets/app.scss'
import './assets/bootstrap-grid.min.css'

import OpenLayersMap from 'vue3-openlayers'
import 'vue3-openlayers/dist/vue3-openlayers.css'

let app = createApp(App)

app.use(router)
app.use(OpenLayersMap)
app.mount('#app')

