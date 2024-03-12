import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import {storage} from "../storage";
import adminRoutes from './adminRoutes'
import ojRoutes from './ojRoutes'
import NProgress from 'nprogress' // nprogress插件
import 'nprogress/nprogress.css' // nprogress样式

// 配置NProgress进度条选项  —— 动画效果
NProgress.configure({ easing: 'ease', speed: 1000,showSpinner: false })

const routes=[...adminRoutes,...ojRoutes] as Array<RouteRecordRaw>

const router=createRouter({
    history: createWebHistory(),
    routes,
    scrollBehavior:(to, from, savedPosition) => {
        if (savedPosition) {
            return savedPosition;
        } else {
            return { left: 0, top: 0 };
        }
    }
})

router.beforeEach((to,from,next)=>{
    NProgress.start()
    if (to.meta.requiresAuth){
        const token=storage.get("authorization")
        if (token){
            to.name&&(document.title =to.name.toString()+" | Aoki")
            next()
        }else{
            storage.set("redirectPath",to.fullPath)
            to.name&&(document.title ="登录 | Aoki")
            next({
                // path:routerPath.Login,
            })
        }
    }else{
        to.name&&(document.title =to.name.toString()+" | Aoki")
        next()
    }
})

router.afterEach(()=>{
    NProgress.done()
})


export default router
