<template>
  <el-row justify="space-between" >
    <el-col :span="18">

      <el-card class="box-card">
        <template #header>
          <div class="card-header">
            <el-row justify="space-between">
              <el-col :span="6">
                <span class="problem-list-title">{{ $t('Problem_List') }}</span>
              </el-col>
              <el-col :span="5">
                <el-input v-model="query.keyword" @keyup.enter="refreshProblemList" :placeholder="$t('Enter_keyword')">
                  <template #append>
                    <el-button :icon="Search" @click="refreshProblemList"/>
                  </template>
                </el-input>
              </el-col>
              <el-col :span="2">
                <el-checkbox v-model="showTags" :label="$t('Show_Tags')"/>
              </el-col>
              <el-col :span="2">
                <el-button type="primary" :icon="Refresh">{{$t('Refresh')}}</el-button>
              </el-col>
            </el-row>
          </div>
        </template>
        <section>
          <el-row justify="space-between">
            <el-col :span="2" style="align-items: center">
              <b class="problem-filter">{{ $t('Problem_Bank') }} :</b>
            </el-col>
            <el-col :span="22" >
              <el-check-tag
                  class="filter-item"
                  :checked="query.oj === 'All'"
                  @click="filterByOJ('All')"
              >{{ $t('All') }}</el-check-tag
              >
              <el-check-tag
                  class="filter-item"
                  :checked="!query.oj || query.oj==='Mine'"
                  @click="filterByOJ('Mine')"
              >{{ $t('My_OJ') }}</el-check-tag
              >
              <el-check-tag
                  class="filter-item"
                  v-for="(remoteOj, index) in REMOTE_OJ"
                  :checked="query.oj === remoteOj.key"
                  :key="index"
                  @click="filterByOJ(remoteOj.key)"
              >{{ remoteOj.name }}</el-check-tag
              >
            </el-col>
          </el-row>
        </section>

        <section style="margin-top: 20px">
          <el-row justify="space-between">
            <el-col :span="2" style="align-items: center">
              <b class="problem-filter">{{ $t('Level') }} :</b>
            </el-col>
            <el-col :span="22" >
              <el-check-tag
                  class="filter-item"
                  :checked="query.difficulty === 'All'||!query.difficulty"
                  @click="filterByDifficulty('All')"
              >{{ $t('All') }}</el-check-tag
              >

              <el-check-tag
                  class="filter-item"
                  v-for="(value, key, index) in PROBLEM_LEVEL"
                  :checked="query.difficulty === key"
                  :key="index"
                  @click="filterByDifficulty(key)"
              >{{ PROBLEM_LEVEL[key].name[locate] }}</el-check-tag>
            </el-col>
          </el-row>
        </section>
        <section style="margin-top: 20px" v-if="filterTagList.length">
          <el-row justify="space-between">
            <el-col :span="2" style="align-items: center">
              <b class="problem-filter">{{ $t('Tags') }} :</b>
            </el-col>
            <el-col :span="22" >
              <el-tag v-for="(tag,index) in filterTagList" closable @close="removeTag(index)" :color="tag.color" class="tag">{{tag.name}}</el-tag>
            </el-col>
          </el-row>
        </section>
        <el-table :data="problemList" stripe @row-click="openProblem">
          <el-table-column width="30">
            <template #default="{row}">
              <el-tooltip v-if="row.myStatus in JUDGE_STATUS" :content="JUDGE_STATUS[row.myStatus].name">
                <el-icon v-if="row.myStatus===1" color="#67C23A"><Select /></el-icon>
                <el-icon v-if="row.myStatus in PROBLEM_LEVEL" :color="PROBLEM_LEVEL[row.myStatus].color"><SemiSelect /></el-icon>
              </el-tooltip>
            </template>
          </el-table-column>
          <el-table-column min-width="50" :label="$t('Problem_ID')" show-overflow-tooltip prop="problem_id"/>
          <el-table-column min-width="150" :label="$t('Problem')" show-overflow-tooltip>
            <template #default="{row}">
              <span class="title-a">{{row.title}}</span>
            </template>

          </el-table-column>
          <el-table-column min-width="100" :label="$t('Level')">
            <template #default="{row}">
              <el-tag style="color: white" :color="PROBLEM_LEVEL[row.difficulty in PROBLEM_LEVEL?row.difficulty:1].color">{{PROBLEM_LEVEL[row.difficulty in PROBLEM_LEVEL?row.difficulty:1].name[locate]}}</el-tag>
            </template>
          </el-table-column>
          <el-table-column v-if="showTags" min-width="230" :label="$t('Tags')">
            <template #default="{row}">
              <el-tag v-for="tag in row.tags" :color="tag.color" class="tag  clickable" @click.stop="addTag(tag)">{{tag.name}}</el-tag>
            </template>
          </el-table-column>
          <el-table-column min-width="80" :label="$t('Total')" prop="total"/>
          <el-table-column min-width="150" :label="$t('AC_Rate')">
            <template #default="{row}">
              <el-tooltip :content="`${row.ac}/${row.total}`">
                <el-progress :text-inside="true"
                             stroke-width="20"
                             :percentage="row.total===0?0:row.ac/row.total"
                              :color="customColors"/>
              </el-tooltip>

            </template>

          </el-table-column>

        </el-table>
        <template #footer>
          <el-pagination
              style="float: right;margin-bottom: 10px"
              v-if="problemList.length"
              v-model:current-page="query.currentPage"
              v-model:page-size="query.limit"
              :page-sizes="pageSizes"
              layout="sizes, prev, pager, next"
              :total="problemList.length"
              @change="refreshProblemList"
          />
        </template>
      </el-card>


    </el-col>

    <el-col :span="5">


    </el-col>


  </el-row>


</template>

<script setup lang="ts">
import {Search, Refresh, Select,SemiSelect} from "@element-plus/icons-vue";
import router from "../../../router";
import {REMOTE_OJ,PROBLEM_LEVEL,JUDGE_STATUS} from "../../../common/constants";
import {computed, onMounted, reactive, ref} from "vue";
import {useI18n} from "vue-i18n";
import {ProblemListType, TagType} from "../../../common/type";
import api from "../../../common/api";
import {ElMessage} from "element-plus";

const showTags=ref(false)
const locate=useI18n().locale
const filterTagList=ref<TagType[]>([])

const removeTag=(index:number)=>{
  filterTagList.value.splice(index,1)
  refreshProblemList()
}

const addTag=(tag:TagType)=>{
  filterTagList.value.push(tag)
  refreshProblemList()
}

const query = reactive({
  ...{
    oj:"Mine",
    difficulty:"",
    limit:30,
    currentPage:1,
    keyword:"",
    tags:computed(()=>filterTagList.value.map(tag=>tag.id))
  },
  ...router.currentRoute.value.query});

const refreshProblemList=()=>{
  router.push({
    name:router.currentRoute.value.name!,
    query:{
      ...router.currentRoute.value.query,
      ...query
    }
  })
}

const filterByOJ=(oj:string)=>{
  query.oj=oj
  refreshProblemList()
}

const filterByDifficulty=(difficulty:string)=>{
  query.difficulty=difficulty
  refreshProblemList()
}


const problemList=ref<ProblemListType[]>([])

const pageSizes=ref([10,30,50,100])
const customColors= [
  { color: '#909399', percentage: 20 },
  { color: '#f56c6c', percentage: 40 },
  { color: '#e6a23c', percentage: 60 },
  { color: '#1989fa', percentage: 80 },
  { color: '#67c23a', percentage: 100 },
]

const getProblemList=()=>{
  api.getProblemList({
    ...router.currentRoute.value.query,
    ...query
  }).then(res=>{
    console.log("题目列表")
    console.log(res)
    console.log(res.data)
    problemList.value=res.data||[]
  }).catch(e=>{
    ElMessage({
      type:"error",
      message:"获取题目列表失败"
    })
  })
}

const openProblem=(row:any)=>{
  console.log(row)
  router.push("/problem/"+row.problem_id)
}

onMounted(()=>{
  getProblemList()
})

</script>

<style scoped>
.problem-list-title {
  font-size: 2em;
  font-weight: 500;
  line-height: 30px;
}
.problem-filter {
  font-weight: bolder;
  white-space: nowrap;
  font-size: 16px;
}
.filter-item {
  margin-right: 1em;
  font-size: 13px;
}

.tag{
  margin-right: 7px;
  color: white;
}
.clickable{
  cursor: pointer;
}
.title-a {
  color: #495060;
  font-family: inherit;
  font-size: 14px;
  font-weight: 500;
}
</style>
