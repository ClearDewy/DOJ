<template>



</template>

<script setup lang="ts">

import api from "../../../common/api";
import {initUserState, user_state} from "../../../storage/user";
import {onMounted} from "vue";
import router from "../../../router";

onMounted(()=>{
  api.logout().finally(()=>{
    Object.assign(user_state.userInfo,initUserState)
    router.push("/home")
  })
})

</script>

<style scoped>

</style>
