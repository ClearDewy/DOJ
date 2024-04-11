<template>
  <el-menu
      :default-active="activeIndex"
      class="el-menu-demo"
      mode="horizontal"
      :ellipsis="false"
      router
  >
    <div class="logo">
      <el-tooltip
          :content="$t('Click_To_Change_Web_Language')"
          placement="bottom"
          effect="dark"
      >
        <el-image
            style="width: 139px; height: 50px"
            :src="logo"
            fit="scale-down"
            @click="changeWebLanguage"
        ></el-image>
      </el-tooltip>
    </div>

    <el-menu-item index="/home"
    >
      <el-icon>
        <HomeFilled/>
      </el-icon>
      {{ $t('NavBar_Home') }}</el-menu-item
    >
    <el-menu-item index="/problem">
      <el-icon>
        <Grid/>
    </el-icon>{{ $t('NavBar_Problem') }}</el-menu-item>
    <el-menu-item index="/training"
    ><el-icon>
      <Checked/>
    </el-icon>{{ $t('NavBar_Training') }}</el-menu-item>
    <el-menu-item index="/contest">
      <el-icon>
        <Trophy/>
      </el-icon>
      {{ $t('NavBar_Contest') }}</el-menu-item>
    <el-menu-item index="/status">
      <el-icon>
        <TrendCharts/>
      </el-icon>
      {{ $t('NavBar_Status') }}</el-menu-item>
    <el-sub-menu index="rank">
      <template #title>
        <el-icon>
          <Histogram/>
        </el-icon>
        {{ $t('NavBar_Rank') }}</template>
      <el-menu-item index="/acm-rank">{{
          $t('NavBar_ACM_Rank')
        }}</el-menu-item>
      <el-menu-item index="/oi-rank">{{
          $t('NavBar_OI_Rank')
        }}</el-menu-item>
    </el-sub-menu>

    <el-menu-item index="/introduction">
      <el-icon>
        <InfoFilled/>
      </el-icon>
      {{ $t('NavBar_Introduction') }}</el-menu-item>

    <div class="flex-grow" />
    <template v-if="!user_state.token">
      <div class="btn-menu">
        <el-button round size="large" @click="web_state.modalStatus={mode: 'login',visible: true}"
        >{{ $t('NavBar_Login') }}
        </el-button>
        <el-button
            v-if="web_state.websiteConfig.allowRegister"
            round
            @click="web_state.modalStatus={mode: 'register',visible: true}"
            style="margin-left: 5px"
            size="large"
        >{{ $t('NavBar_Register') }}
        </el-button>
      </div>
    </template>
    <template v-else>

      <el-dropdown
          class="drop-menu"
          @command="handleRoute"
          placement="bottom"
          trigger="hover"
      >
            <span class="el-dropdown-link">
              <el-avatar
                  :username="user_state.userInfo.username"
                  :inline="true"
                  :size="32"
                  :src="user_state.userInfo.avatar||defaultAvatar"
                  class="avatar"
              ></el-avatar>
              {{ user_state.userInfo.username }}
              <el-icon class="el-icon--right" f>
                <ArrowDown />
              </el-icon>
            </span>

        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="/user-home">{{
                $t('NavBar_UserHome')
              }}</el-dropdown-item>
            <el-dropdown-item command="/status?onlyMine=true">{{
                $t('NavBar_Submissions')
              }}</el-dropdown-item>
            <el-dropdown-item command="/setting">{{
                $t('NavBar_Setting')
              }}</el-dropdown-item>
            <el-dropdown-item command="/admin">{{
                $t('NavBar_Management')
              }}</el-dropdown-item>
            <el-dropdown-item divided command="/logout">{{
                $t('NavBar_Logout')
              }}</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </template>
  </el-menu>

</template>

<script setup lang="ts">
import logo from "../../../assets/logo.png"
import {ref} from "vue";
import {changeWebLanguage} from "../../../i18n";
import {user_state} from "../../../storage/user";
import {web_state} from "../../../storage";
import router from "../../../router";
import {HomeFilled,Grid,Checked,Trophy,TrendCharts,Histogram,InfoFilled,ArrowDown} from "@element-plus/icons-vue"
import defaultAvatar from "../../../assets/default.jpg"

const activeIndex=ref("1")
const handleRoute=(route:string)=> {
  //电脑端导航栏路由跳转事件
  if (route && route.split('/')[1] != 'admin') {
    router.push(route)
  } else {
    window.open('/admin/');
  }
}


</script>

<style scoped>
.flex-grow {
  flex-grow: 1;
}

.logo {
  cursor: pointer;
  margin-left: 2%;
  margin-right: 2%;
  float: left;
  width: 139px;
  height: 42px;
  margin-top: 5px;
}

.el-menu--popup .el-menu-item{
  justify-content: center;
}

.btn-menu {
  font-size: 16px;
  float: right;
  margin-right: 10px;
  margin-top: 10px;
}
.drop-menu {

  margin-top: 9px;
  margin-bottom: 9px;
  margin-right: 30px;
  font-weight: 500;
  font-size: 18px;

}

.el-dropdown-link {
  cursor: pointer;
  color: #409eff !important;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
}
.avatar{
  margin-right: 5px;
}
</style>

