<template>
  <el-scrollbar>
    <el-row justify="space-between">
      <el-col>

      </el-col>
      <div class="flex-grow" />
      <el-col>

      </el-col>
    </el-row>

    <div class="problem-content">
      <div v-if="problem.description">
        <p class="title">{{$t("Description")}}</p>
        <Editor editor-id="Description"
                v-model="problem.description"
                :preview-only="true"/>
      </div>
      <div v-if="problem.input">
        <p class="title">{{$t("Input")}}</p>
        <Editor editor-id="Input"
                v-model="problem.input"
                :preview-only="true"/>
      </div>
      <div v-if="problem.output">
        <p class="title">{{$t("Output")}}</p>
        <Editor editor-id="Output"
                v-model="problem.output"
                :preview-only="true"/>
      </div>
      <div v-if="problem.examples">
        <el-row v-for="(example,index) in problem.examples" justify="space-between">
          <el-col class="example-row" :span="11">
              <p class="title">{{$t("Sample_Input")}} {{index+1}}
              <el-icon class="copy" @click="copy(example.input)"><DocumentCopy /></el-icon></p>
            <pre>{{example.input}}</pre>
          </el-col>
          <el-col class="example-row"  :span="11">
              <p class="title">{{$t("Sample_Output")}} {{index+1}}
              <el-icon class="copy" @click="copy(example.output)"><DocumentCopy /></el-icon></p>
            <pre>{{example.output}}</pre>
          </el-col>
        </el-row>
      </div>
      <div v-if="problem.hint">
        <p class="title">{{$t("Hint")}}</p>
        <Editor editor-id="Hint"
                v-model="problem.hint"
                :preview-only="true"/>
      </div>
      <div v-if="problem.source">
        <p class="title">{{$t("Source")}}</p>
        <Editor editor-id="Source"
                v-model="problem.source"
                :preview-only="true"/>
      </div>
    </div>
  </el-scrollbar>
</template>

<script setup lang="ts">

import {PROBLEM_LEVEL} from "../../../common/constants";
import {useI18n} from "vue-i18n";
import router from "../../../router";
import Editor from "../common/MdEditor.vue";
import useClipboard from 'vue-clipboard3'
import {ElMessage} from "element-plus";
import {DocumentCopy} from "@element-plus/icons-vue"

defineProps<{
  problem:any
}>();

const { toClipboard } = useClipboard()
const {t}=useI18n()
const copy = async (value:string) => {
  try {
    await toClipboard(value)
    ElMessage({
      type:"success",
      message:t("Copied_successfully")
    })
  } catch (e) {
    console.error(e)
    ElMessage({
      type:"error",
      message:t("Copied_failed")
    })
  }
}

</script>

<style scoped>
.flex-grow {
  flex-grow: 1;
}
.problem-title {
  font-size: 32px;
  font-weight: 500;
  padding-top: 10px;
  padding-bottom: 20px;
  line-height: 30px;
}
.problem-content {

}
.problem-content .title {
  font-size: 16px;
  font-weight: 600;
  margin: 8px 0 20px 0;
  color: #3091f2;
}
.problem-content .copy {
  margin-left: 8px;
  cursor: pointer;
}
.problem-content .copy:hover {
  color: white;
  background: #3091f2;
}

.example-row{

}
.example-row {
  display: flex;
  flex-direction: column;
}

.example-row pre {
  flex-grow: 1;
  align-self: stretch;
  padding: 5px 10px;
  white-space: pre;
  margin-top: 10px;
  margin-bottom: 10px;
  background: #f1f1f1;
  border: 1px dashed #e9eaec;
  overflow: auto;
  font-size: 1.1em;
  font-family: Consolas,Menlo,Courier,monospace;
  word-break: break-all;
  display: block;
}
</style>
