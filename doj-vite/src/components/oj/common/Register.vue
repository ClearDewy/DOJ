<template>
  <el-dialog v-if="web_state.modalStatus.mode==='register'" v-model="web_state.modalStatus.visible" center width="370px">
    <template #header>
      <span class="dialog-header">{{`${$t('Dialog_Register')} - ${web_state.websiteConfig.shortName}`}}</span>
    </template>
    <el-form
        :model="formRegister"
        :rules="rules"
        ref="ruleFormRef"
        label-width="0"
    >
      <el-form-item prop="username">
        <el-input
            v-model="formRegister.username"
            :prefix-icon="User"
            :placeholder="$t('Register_Username')"
            width="100%"
            size="large"
            @keyup.enter.native="handleRegister(ruleFormRef)"
        ></el-input>
      </el-form-item>
      <el-form-item prop="password">
        <el-input
            v-model="formRegister.password"
            :prefix-icon="Lock"
            size="large"
            :placeholder="$t('Register_Password')"
            type="password"
            @keyup.enter.native="handleRegister(ruleFormRef)"
        ></el-input>
      </el-form-item>
      <el-form-item prop="passwordAgain">
        <el-input
            v-model="formRegister.passwordAgain"
            :prefix-icon="Lock"
            size="large"
            :placeholder="$t('Register_Password_Again')"
            type="password"
            @keyup.enter.native="handleRegister(ruleFormRef)"
        ></el-input>
      </el-form-item>
      <el-form-item prop="email">
        <el-input
            v-model="formRegister.email"
            :prefix-icon="Message"
            size="large"
            :placeholder="$t('Register_Email')"
            type="password"
            @keyup.enter.native="handleRegister(ruleFormRef)"
        ><template #append>
          <el-button :icon="Message" @click=""/>
        </template></el-input>
      </el-form-item>
      <el-form-item prop="code">
        <el-input
            v-model="formRegister.code"
            :prefix-icon="Stamp"
            size="large"
            :placeholder="$t('Register_Email_Captcha')"
            type="password"
            @keyup.enter.native="handleRegister(ruleFormRef)"
        ></el-input>
      </el-form-item>
    </el-form>
    <el-button
        type="primary"
        @click="handleRegister(ruleFormRef)"
        :loading="btnRegisterLoading"
        class="login-btn"
    >{{ $t('Register_Btn') }}</el-button
    >
    <el-link
        v-if="web_state.websiteConfig.allowRegister"
        type="primary"
        @click="web_state.modalStatus.mode='login'"
    >{{ $t('Register_Already_Registed') }}</el-link
    >

  </el-dialog>
</template>

<script setup lang="ts">
import {web_state} from "../../../storage";
import {User,Lock,Message,Stamp} from "@element-plus/icons-vue";
import {reactive, ref} from "vue";
import {FormInstance, FormRules} from "element-plus";
import {useI18n} from "vue-i18n";

const {t} = useI18n();

const formRegister=reactive({
    username: '',
    password: '',
    passwordAgain:'',
    email:'',
    code:''
  }
)
const CheckUsernameNotExist=(rule: any, value: any, callback: any)=>{

}
const CheckEmailNotExist=(rule: any, value: any, callback: any)=>{

}
const CheckAgainPassword=(rule: any, value: any, callback: any)=>{
  if (value !== formRegister.password) {
    callback(new Error(t('Password_does_not_match')));
  }
  callback();
}

const rules=reactive<FormRules<typeof formRegister>>({
  username: [
    {
      required: true,
      message: t('Username_Check_Required'),
      trigger: 'blur',
    },
    {
      validator: CheckUsernameNotExist,
      trigger: 'blur',
      message: t('The_username_already_exists'),
    },
    {
      max: 20,
      message: t('Username_Check_Max'),
      trigger: 'blur',
    },
  ],

  email: [
    {
      required: true,
      message: t('Email_Check_Required'),
      trigger: 'blur',
    },
    {
      type: 'email',
      message: t('Email_Check_Format'),
      trigger: 'blur',
    },
    {
      validator: CheckEmailNotExist,
      message: t('The_email_already_exists'),
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
    }
  ],
  passwordAgain: [
    {
      required: true,
      message: t('Password_Again_Check_Required'),
      trigger: 'blur',
    },
    { validator: CheckAgainPassword, trigger: 'change' },
  ],
  code: [
    {
      required: true,
      message: t('Code_Check_Required'),
      trigger: 'blur',
    },
    {
      min: 6,
      max: 6,
      message: t('Code_Check_Length'),
      trigger: 'blur',
    },
  ],
})
const ruleFormRef = ref<FormInstance>()
const btnRegisterLoading=ref(false)

const handleRegister=async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  btnRegisterLoading.value=true
  await formEl.validate((valid, fields) => {
    if (valid) {
      btnRegisterLoading.value=false
    }else {
      btnRegisterLoading.value=false
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
