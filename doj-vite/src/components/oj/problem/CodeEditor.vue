<template>
<el-card class="full-height" shadow="never">
  <el-row justify="space-between">
    <el-col class="item-center" :span="18">
      <span id="lang">{{$t("Lang")}}:</span>
      <el-select v-model="currentLang.name" @change="selectLang">
        <el-option v-for="lang in languages" :key="lang.lid" :label="lang.name" :value="lang.name"/>
      </el-select>
      <el-button-group>
        <el-tooltip placement="top" :content="$t('Reset_Code')">
          <el-button :icon="Refresh"/>
        </el-tooltip>
        <el-tooltip placement="top" :content="$t('Get_Recently_Passed_Code')">
          <el-button :icon="Download"/>
        </el-tooltip>
      </el-button-group>

    </el-col>
    <el-col class="item-center" :span="6">
      <div class="flex-grow"/>
      <el-tooltip placement="top" :content="$t('Upload_file')">
        <el-button :icon="Upload"/>
      </el-tooltip>
      <el-popover
          placement="bottom"
          :width="300"
          trigger="click"
      >
        <template #reference>
          <div>
            <el-tooltip placement="top" :content="$t('Code_Editor_Setting')">
              <el-button :icon="Setting"/>
            </el-tooltip>
          </div>
        </template>
        <div class="setting-title">{{$t("Setting")}}</div>
        <div class="setting-item">
          <span class="setting-item-name">
              {{ $t('Theme') }}
            </span>
          <el-select v-model="themeName"
                     class="setting-item-value"
                     :teleported="false"
          >
            <el-option v-for="(value,key) in Themes" :key="key" :value="key" :label="key"/>
          </el-select>
        </div>
<!--        <div class="setting-item">-->
<!--          <span class="setting-item-name">-->
<!--              {{ $t('FontSize') }}-->
<!--            </span>-->
<!--        </div>-->
        <div class="setting-item">
          <span class="setting-item-name">
              {{ $t('TabSize') }}
            </span>
          <el-select
              :teleported="false"
              v-model="codeMirrorConfig.tab_size"
              class="setting-item-value"
          >
            <el-option
                :label="$t('Two_Spaces') "
                :value="2"
            >
              {{ $t('Two_Spaces') }}
            </el-option>
            <el-option
                :label="$t('Four_Spaces') "
                :value="4"
            >
              {{ $t('Four_Spaces') }}
            </el-option>
            <el-option
                :label="$t('Eight_Spaces') "
                :value="8"
            >
              {{ $t('Eight_Spaces') }}
            </el-option>
          </el-select>
        </div>
      </el-popover>
    </el-col>
  </el-row>
  <div style="position: relative">
    <Codemirror :modelValue="modelValue" @update:modelValue="handleModelValueUpdate"
                :extensions="codeMirrorConfig.extensions"
                :tab-size="codeMirrorConfig.tab_size"
                :autofocus="true"
                :indent-with-tab="true"
    />
    <el-drawer v-model="openTestCaseDrawer"
               direction="btt"
               :with-header="false"
               :modal="true"
               modal-class="opacity"
               :size="330"
    >
      <div style="position: relative">
        <el-tabs type="border-card" class="init">
          <el-tab-pane :label="$t('Test_Case')" style="height: 250px">
            <div style="margin-bottom: 15px">
              <el-check-tag v-for="(example,index) in examples"
                            :checked="index===runningExample.index" @click="inputExample(example,index)">{{$t("Fill_Case")}}{{index+1}} </el-check-tag>
            </div>
            <el-input type="textarea"
                      v-model="runningExample.input"
                      @change="runningExample.index=-1"
                      :rows="9" maxlength="1000"
                      show-word-limit/>
          </el-tab-pane>
          <el-tab-pane :label="$t('Test_Result')" style="height: 250px">

          </el-tab-pane>
        </el-tabs>
        <el-button :icon="VideoPlay" style="color: #67C23A;border-color: #67C23A;position: absolute;top: 4px;right: 20px"
                   @click.stop="submitTestJudge"
        >{{$t("Running_Test")}}</el-button>
      </div>
    </el-drawer>
  </div>

  <el-row>
    <el-col :span="8">

    </el-col>
    <el-col :span="16">
      <div  style="float: right">
        <el-button
            :type="openTestCaseDrawer?'success':''"
            :style="{color:openTestCaseDrawer?'#ffffff':'#67c23a'}"
            style="border-color: #67c23a;"
            @click="openTestCaseDrawer=true"
            effect="plain"
        >
          <svg
              t="1653665263421"
              class="icon"
              viewBox="0 0 1024 1024"
              version="1.1"
              xmlns="http://www.w3.org/2000/svg"
              p-id="1656"
              width="12"
              height="12"
              style="vertical-align: middle;"
          >
            <path
                d="M1022.06544 583.40119c0 11.0558-4.034896 20.61962-12.111852 28.696576-8.077979 8.077979-17.639752 12.117992-28.690436 12.117992L838.446445 624.215758c0 72.690556-14.235213 134.320195-42.718941 184.89915l132.615367 133.26312c8.076956 8.065699 12.117992 17.634636 12.117992 28.690436 0 11.050684-4.034896 20.614503-12.117992 28.691459-7.653307 8.065699-17.209964 12.106736-28.690436 12.106736-11.475356 0-21.040199-4.041036-28.690436-12.106736L744.717737 874.15318c-2.124384 2.118244-5.308913 4.88424-9.558703 8.283664-4.259 3.3984-13.180184 9.463536-26.78504 18.171871-13.598716 8.715499-27.415396 16.473183-41.439808 23.276123-14.029528 6.797823-31.462572 12.966313-52.289923 18.49319-20.827351 5.517667-41.446971 8.28571-61.842487 8.28571L552.801776 379.38668l-81.611739 0 0 571.277058c-21.668509 0-43.250036-2.874467-64.707744-8.615215-21.473057-5.734608-39.960107-12.749372-55.476499-21.039175-15.518438-8.289804-29.541827-16.572444-42.077328-24.867364-12.541641-8.290827-21.781072-15.193027-27.739784-20.714787l-9.558703-8.93244L154.95056 998.479767c-8.500605 8.921183-18.699897 13.386892-30.606065 13.386892-10.201339 0-19.335371-3.40454-27.409257-10.202363-8.079002-7.652284-12.437264-17.10968-13.080923-28.372188-0.633427-11.263531 2.659573-21.143553 9.893324-29.647227l128.787178-144.727219c-24.650423-48.464805-36.980239-106.699114-36.980239-174.710091L42.738895 624.207571c-11.057847 0-20.61655-4.041036-28.690436-12.111852-8.079002-8.082072-12.120039-17.640776-12.120039-28.696576 0-11.050684 4.041036-20.61962 12.120039-28.689413 8.073886-8.072863 17.632589-12.107759 28.690436-12.107759l142.81466 0L185.553555 355.156836l-110.302175-110.302175c-8.074909-8.077979-12.113899-17.640776-12.113899-28.691459 0-11.04966 4.044106-20.61962 12.113899-28.690436 8.071839-8.076956 17.638729-12.123109 28.691459-12.123109 11.056823 0 20.612457 4.052293 28.692482 12.123109l110.302175 110.302175 538.128077 0 110.303198-110.302175c8.070816-8.076956 17.632589-12.123109 28.690436-12.123109 11.050684 0 20.617573 4.052293 28.689413 12.123109 8.077979 8.070816 12.119015 17.640776 12.119015 28.690436 0 11.050684-4.041036 20.614503-12.119015 28.691459l-110.302175 110.302175 0 187.448206 142.815683 0c11.0558 0 20.618597 4.034896 28.690436 12.113899 8.076956 8.069793 12.117992 17.638729 12.117992 28.683273l0 0L1022.06544 583.40119 1022.06544 583.40119zM716.021162 216.158085 307.968605 216.158085c0-56.526411 19.871583-104.667851 59.616796-144.414087 39.733956-39.746236 87.88256-59.611679 144.411017-59.611679 56.529481 0 104.678084 19.865443 144.413064 59.611679C696.156742 111.48921 716.021162 159.631674 716.021162 216.158085L716.021162 216.158085 716.021162 216.158085 716.021162 216.158085z"
                p-id="1657"
                :fill="openTestCaseDrawer?'#ffffff':'#67c23a'"
            >
            </path>
          </svg>
          <span style="vertical-align: middle;">{{ $t('Online_Test') }}</span>
        </el-button>
        <el-button
            type="primary"
            :icon="Edit"
            :loading="submitting"
            @click="submitCode"
            :disabled="submitted"
        >
          <span v-if="submitting">{{ $t('Submitting') }}</span>
          <span v-else>{{ $t('Submit') }}</span>
        </el-button>
      </div>
    </el-col>

  </el-row>
</el-card>
</template>

<script setup lang="ts">
import {Codemirror} from 'vue-codemirror';
import {computed, reactive, ref} from "vue";
import {Refresh,Download,Upload,Setting} from "@element-plus/icons-vue"
import {oneDark} from "@codemirror/theme-one-dark"
import {StreamLanguage} from '@codemirror/language'
import {Edit,VideoPlay} from "@element-plus/icons-vue"

// import languages
import {c} from "@codemirror/legacy-modes/mode/clike"
import {cpp} from "@codemirror/lang-cpp"
import {java} from "@codemirror/lang-java"
import {python} from "@codemirror/lang-python"
import {javascript} from "@codemirror/lang-javascript"
import {go} from "@codemirror/legacy-modes/mode/go"
//@ts-ignore
import {csharp} from "@replit/codemirror-lang-csharp"
import {php} from "@codemirror/lang-php"
import {pascal} from "@codemirror/legacy-modes/mode/pascal"
import {d} from "@codemirror/legacy-modes/mode/d"
import {haskell} from "@codemirror/legacy-modes/mode/haskell"
import {perl} from "@codemirror/legacy-modes/mode/perl"
import {ruby} from "@codemirror/legacy-modes/mode/ruby"
import {rust} from "@codemirror/legacy-modes/mode/rust"
import {fortran} from "@codemirror/legacy-modes/mode/fortran"
import {ElMessage, ElNotification} from "element-plus";
import {useI18n} from "vue-i18n";
import api from "../../../common/api";
import {user_state} from "../../../storage/user";
import {web_state} from "../../../storage";

const {t}=useI18n()

type LanguageType={lid:number,name:string,content_type:string}
type ExampleType={input:string,output:string}

const props = defineProps<{
  modelValue:string
  languages:LanguageType[],
  examples?:ExampleType[],
  problem_id:string
}>();

const emits=defineEmits(['update:modelValue'])
const handleModelValueUpdate = (newValue:any) => {
  emits('update:modelValue',newValue)
}

const currentLang=reactive<LanguageType>({
  lid: -1,
  name: '',
  content_type: ''
})

const selectLang=(name:string)=>{
  Object.assign(currentLang,props.languages.find(lang=>lang.name===name))
}

const Language={
  c:StreamLanguage.define(c),
  cpp:cpp(),
  java:java(),
  python:python(),
  javascript:javascript(),
  go:StreamLanguage.define(go),
  csharp:csharp(),
  php:php(),
  pascal:StreamLanguage.define(pascal),
  d:StreamLanguage.define(d),
  haskell:StreamLanguage.define(haskell),
  perl:StreamLanguage.define(perl),
  ruby:StreamLanguage.define(ruby),
  rust:StreamLanguage.define(rust),
  fortran:StreamLanguage.define(fortran)
}
const Themes={
  "Default":undefined,
  "One Dark":oneDark,
}
const themeName=ref<keyof typeof Themes>("Default")

const codeMirrorConfig=reactive({
  tab_size:4,
  extensions : computed(()=>{
    const exts: any[]=[]
    if (currentLang.content_type in Language){
      exts.push(Language[currentLang.content_type as keyof typeof Language])
    }
    if(Themes[themeName.value]){
      exts.push(Themes[themeName.value])
    }
    return exts
  }),
  // fontSize:32
})

const openTestCaseDrawer=ref(false)

const runningExample=reactive({
  index:0,
  input:"",
  output:""
})


const inputExample=(example:ExampleType,index:number)=>{
  runningExample.index=index
  runningExample.input=example.input
  runningExample.output=example.output
}

const submitting=ref(false)
const submitted=ref(false)
const submitCode=()=>{
  submitting.value=true
  if (currentLang.lid<=0||props.modelValue===""||props.modelValue.length>65535){
    ElNotification({
      type:"error",
      title:t("Submit_code_fail"),
      message:t(currentLang.lid<=0?"Please_choose_code_language":props.modelValue===""?"Code_can_not_be_empty":"Code_Length_can_not_exceed_65535")
    })
    submitting.value=false
    return
  }
  api.submitCode({
    lid:currentLang.lid,
    code:props.modelValue,
    problem_id:props.problem_id
  }).then(res=>{
    submitted.value=true
    ElMessage({
      type:"success",
      message:t("Submit_code_successfully")
    })
  }).catch(e=>{
    ElNotification({
      type:"error",
      title:t("Submit_code_fail"),
      message:e.msg||t("Submit_code_fail")
    })
  }).finally(()=>{
    submitting.value=false
  })
}

const submitTestJudge=()=>{
  if (!user_state.token){
    web_state.modalStatus={ mode: 'login', visible: true }
    return
  }



}
</script>

<style scoped>
.full-height{
  height: 100%;
}
.flex-grow {
  flex-grow: 1;
}

.item-center{
  display: flex;
  align-items: center;
}
#lang{
  font-size: 14px;
  margin-right: 14px;
}
.setting-title {
  border-bottom: 1px solid #f3f3f6;
  color: #000;
  font-weight: 700;
  padding: 10px 0;
}
.setting-item {
  display: flex;
  padding: 15px 0 0;
}
.setting-item-name {
  flex: 2;
  color: #333;
  font-weight: 700;
  font-size: 13px;
  margin-top: 7px;
}
.setting-item-value {
  width: 140px;
  margin-left: 15px;
  flex: 5;
}
</style>
<style>
.cm-editor{
  height: calc(100vh - 64px - 40px)!important;
}
.el-overlay{
  position: absolute!important;
}

.opacity{
  color: transparent!important;
  background: transparent!important;
}
.init{
  color: initial!important;
  background: initial!important;
}

.el-drawer__body{
  padding: 0!important;
}
</style>