import { createRouter, createWebHashHistory } from "vue-router";
const login = () => import ("../components/Login.vue")
const main = () => import ("../components/MainView.vue")
const routes = [
    {
        path: "/",
        name: "Home",
        component: login,
    },
    {
        path: "/main",
        name: "MainView",
        component: main,
    },
];

const router = createRouter({
    history: createWebHashHistory(),
    routes,
});

export default router;
