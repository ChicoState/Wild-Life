import {createRouter, createWebHistory} from "vue-router";
// Views
import Home from "./views/Home.vue";
import Results from "./views/Results.vue";
// Authorization
import Login from "./views/Login.vue";
import Register from "./views/Register.vue";

const routes = [
    {
        path: '/',
        name: 'Home',
        component: Home,
        icon: 'fa-house',
    },
    {
        path: '/results/{token}',
        name: 'Results',
        component: Results
    },
    {
        path: '/login',
        name: 'Login',
        component: Login
    },
    {
        path: '/register',
        name: 'Register',
        component: Register
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes: routes,
})

export default router
