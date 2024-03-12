<template>
  <span class="problem-title">{{problem.title}}</span>
  <div>
    <el-tag v-if="!problem.tags||problem.tags.length===0" class="tag" style="color: #409EFF" color="#FFFFFF">{{$t("No_tag")}}</el-tag>
    <el-tag v-else v-for="tag in problem.tags" :color="tag.color" class="tag">{{tag.name}}</el-tag>
  </div>
  <div class="question-intr">

    <!--    <span>{{ $t('Time_Limit') }}：C/C++{{ problem.time_limit }}MS，-->
    <!--      {{$t('Other') }}{{ problem.time_limit * 2 }}MS</span><br />-->
    <!--      <span>{{ $t('Memory_Limit') }}：-->
    <!--      C/C++{{ problem.memory_limit }}MB，{{$t('Other') }}-->
    <!--      {{ problem.memory_limit * 2 }}MB</span><br />-->

    <span style="display: flex;align-items: center">{{ $t('Level') }}：
    <el-tag style="color: white" :color="PROBLEM_LEVEL[problem.difficulty in PROBLEM_LEVEL?problem.difficulty:1].color">{{PROBLEM_LEVEL[problem.difficulty in PROBLEM_LEVEL?problem.difficulty:1].name[locale]}}</el-tag>
    </span>
    <span>{{ $t('Created') }}：<el-link
        type="info"
        class="author-name"
        @click="goUserHome(problem.author)"
    >{{ problem.author }}</el-link></span><br />
  </div>
  <el-table :data="problem.languages" stripe style="width: 100%">
    <el-table-column :label="$t('Language')" prop="name"/>
    <el-table-column :label="$t('Time')">
      <template #default="{row}">
        {{row.time_limit}} ms
      </template>
    </el-table-column>
    <el-table-column :label="$t('Memory')">
      <template #default="{row}">
        {{row.memory_limit}} mb
      </template>
    </el-table-column>


  </el-table>
</template>

<script setup lang="ts">

import {PROBLEM_LEVEL} from "../../../common/constants";
import router from "../../../router";
import {useI18n} from "vue-i18n";
defineProps<{
  problem:any
}>();
const locale=useI18n().locale
const goUserHome=(username:string)=>{
  router.push({
    path:"/user-home",
    params:{
      username
    }
  })
}
</script>

<style scoped>

.tag{
  margin-right: 7px;
  color: white;
}
.question-intr {
  margin-top: 6px;
  border-radius: 4px;
  border: 1px solid #ddd;
  border-left: 2px solid #3498db;
  background: #fafafa;
  padding: 10px;
  line-height: 1.8;
  margin-bottom: 10px;
  font-size: 14px;
}

.problem-title {
  font-size: 32px;
  font-weight: 500;
  padding-top: 10px;
  padding-bottom: 20px;
  line-height: 30px;
}
</style>