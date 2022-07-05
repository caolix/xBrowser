import { createRouter, createWebHashHistory } from "vue-router";
const login = () => import ("../components/Login.vue")
const main = () => import ("../components/MainView.vue")
const bucketList = () => import ("../components/BucketList.vue")
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
        children: [
            {
                path: '',
                components: {
                    default: bucketList()
                }
            }
        ]
    },
];

const router = createRouter({
    history: createWebHashHistory(),
    routes,
});

export default router;
