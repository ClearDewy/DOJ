<template>
  <el-row>
    <el-col :span="12">
      <el-tabs type="border-card" class="full-height">
        <el-tab-pane class="full-height-pane" v-loading="loading" >
          <template #label>
            <el-icon><Notification /></el-icon>
            {{$t("Problem_Information")}}
          </template>
          <ProblemInformation :problem="problem"/>
        </el-tab-pane>
        <el-tab-pane class="full-height-pane" v-loading="loading" >
          <template #label>
            <el-icon><DataBoard /></el-icon>
            {{$t("Problem_Description")}}
          </template>
          <ProblemStatement :problem="problem"/>
        </el-tab-pane>
        <el-tab-pane class="full-height-pane">
          <template #label>
            <el-icon><Clock /></el-icon>
            {{$t("My_Submission")}}
          </template>
          <ProblemMySubmission />
        </el-tab-pane>
      </el-tabs>

    </el-col>
    <el-col :span="12" v-loading="loading">
      <CodeEditor class="full-height" v-model="code" :problem_id="problem_id" :languages="problem.languages" :examples="problem.examples"/>
    </el-col>

  </el-row>

</template>

<script setup lang="ts">
import {DataBoard,Clock,Notification} from "@element-plus/icons-vue"

import CodeEditor from "../../../components/oj/problem/CodeEditor.vue";
import ProblemStatement from "../../../components/oj/problem/ProblemStatement.vue";
import ProblemMySubmission from "../../../components/oj/problem/ProblemMySubmission.vue";
import {onMounted, reactive, ref} from "vue";
import router from "../../../router";
import api from "../../../common/api";
import {ElMessage} from "element-plus";
import ProblemInformation from "../../../components/oj/problem/ProblemInformation.vue";

const loading=ref(true)

const problem=reactive<{
  languages:any[]
  examples:{input:string,output:string}[]
}>({
  languages:[],
  examples:[]
})

const code=ref("")
const problem_id=ref("")

onMounted(()=>{
  problem_id.value=router.currentRoute.value.params.problemID as string
  api.getProblemDetail(problem_id.value).then(res=>{
    console.log(res.data)
    Object.assign(problem,res.data)
    // console.log(JSON.parse(res.data.examples))
    problem.examples=JSON.parse(res.data.examples)
  }).catch(e=>{
    ElMessage({
      type:'error',
      message:"获取题目详情失败"
    })
  }).finally(()=>{
    loading.value=false
  })
})

</script>

<style scoped>
.full-height-pane{
  height: calc(100vh - 80px);
}
.full-height{
  height: 100vh;
}

</style>
