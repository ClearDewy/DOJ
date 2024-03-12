<template>
  <el-dialog v-if="web_state.modalStatus.mode==='login'" v-model="web_state.modalStatus.visible" center width="370px">
    <template #header>
      <span class="dialog-header">{{`${$t('Dialog_Login')} - ${web_state.websiteConfig.shortName}`}}</span>
    </template>
    <el-form
        :model="formLogin"
        :rules="rules"
        ref="ruleFormRef"
        label-width="0"
    >
      <el-form-item prop="username">
        <el-input
            v-model="formLogin.username"
            :prefix-icon="User"
            :placeholder="$t('Login_Username')"
            width="100%"
            size="large"
            @keyup.enter.native="handleLogin(ruleFormRef)"
        ></el-input>
      </el-form-item>
      <el-form-item prop="password">
        <el-input
            v-model="formLogin.password"
            :prefix-icon="Lock"
            size="large"
            :placeholder="$t('Login_Password')"
            type="password"
            @keyup.enter.native="handleLogin(ruleFormRef)"
        ></el-input>
      </el-form-item>
    </el-form>
    <el-button
        type="primary"
        @click="handleLogin(ruleFormRef)"
        :loading="btnLoginLoading"
        class="login-btn"
    >{{ $t('Login_Btn') }}</el-button
    >
    <el-link
        v-if="web_state.websiteConfig.register"
        type="primary"
        @click="web_state.modalStatus.mode='register'"
    >{{ $t('Login_No_Account') }}</el-link
    >
    <el-link
        type="primary"
        style="float: right"
    >{{ $t('Login_Forget_Password') }}</el-link
    >

  </el-dialog>
</template>

<script setup lang="ts">
import {web_state} from "../../../storage";
import {User,Lock} from "@element-plus/icons-vue";
import {reactive, ref} from "vue";
import {ElMessage, FormInstance, FormRules} from "element-plus";
import {useI18n} from "vue-i18n";
import api from "../../../common/api";
import {user_state} from "../../../storage/user";
import router from "../../../router";

const {t} = useI18n();

const formLogin=reactive({
      username: '',
      password: '',
    }
)
const rules=reactive<FormRules<typeof formLogin>>({
  username: [
    {
      required: true,
      message: t('Username_Check_Required'),
      trigger: 'blur',
    },
    {
      max: 20,
      message: t('Username_Check_Max'),
      trigger: 'blur',
    },
  ],
  password: [
    {
      required: true,
      message: t('Password_Check_Required'),
      trigger: 'blur',
    },
    {
      min: 6,
      max: 20,
      message: t('Password_Check_Between'),
      trigger: 'blur',
    },
  ],
},)
const ruleFormRef = ref<FormInstance>()
const btnLoginLoading=ref(false)

const handleLogin=async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  btnLoginLoading.value=true
  await formEl.validate((valid, fields) => {
    if (valid) {
      api.login(formLogin).then(res=>{
        user_state.userInfo=res.data
        ElMessage({
          message:t("Welcome_Back"),
          type:'success'
        })
        router.push(router.currentRoute.value.path)
      }).finally(()=>{
        btnLoginLoading.value=false
        web_state.modalStatus.visible=false
      })
    }else {
      btnLoginLoading.value=false
    }
  })
}
</script>

<style scoped>
.dialog-header{
  font-size: 22px;
  font-weight: 600;
  margin: 0;
}

.login-btn{
  width: 100%;
  margin-bottom: 10px;
}

</style>
