<template>
  <el-config-provider :locale="e_locale">
    <el-scrollbar>
      <transition name="el-zoom-in-bottom">
        <div id="app">
          <el-backtop :right="10"></el-backtop>
          <div v-if="!isAdminView">
            <el-container>
              <el-header><NavBar/></el-header>
              <el-main><router-view/></el-main>
            </el-container>
            <Login/>
            <Register/>
          </div>
          <div v-else>
            <el-container>
              <el-aside>

              </el-aside>
              <el-main>

              </el-main>
            </el-container>
          </div>
        </div>
      </transition>
    </el-scrollbar>
  </el-config-provider>
</template>

<script setup lang="ts">
import NavBar from "./components/oj/common/NavBar.vue"
//@ts-ignore
import zhCN from 'element-plus/dist/locale/zh-cn.mjs'
//@ts-ignore
import enUS from 'element-plus/dist/locale/en.mjs'
import {computed, onMounted, ref, watch} from "vue";
import {useI18n} from "vue-i18n";
import router from "./router";
import Login from "./components/oj/common/Login.vue";
import Register from "./components/oj/common/Register.vue";
import api from "./common/api";
import {web_state} from "./storage";
import {ElMessage} from "element-plus";
const locale=useI18n().locale
// 设置element-plus内置语言
const e_locale = computed(() => (locale.value === 'zh-CN' ? zhCN : enUS))

const isAdminView=ref(false)

watch(router.currentRoute,(new_route,old_route)=>{
  isAdminView.value=(new_route.path!==old_route.path&&old_route.path.startsWith("/admin"))
})


onMounted(()=>{
  api.getWebConfig().then(res=>{
    Object.assign(web_state.websiteConfig,res)
  }).catch(e=>{
    ElMessage({
      type:"error",
      message:e.msg||"获取服务器配置失败"
    })
  })
})

</script>

<style scoped>
.el-header{
  padding: 0;
  box-shadow: var(--el-box-shadow-lighter)
}
.el-main{
  padding: 20px 20px 0 20px
}
.el-container{
  padding: 0;
}
</style>
