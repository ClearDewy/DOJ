import api from '../common/api'
import { CONTEST_STATUS, CONTEST_TYPE } from '../common/constants'
import {reactive, watch} from "vue";
import storage from "./index";
export const contest_state = reactive({
  now: new Date(),
  intoAccess: false, // 比赛进入权限
  submitAccess:false, // 保护比赛的提交权限
  forceUpdate: false, // 强制实时榜单
  removeStar: false, // 榜单去除打星队伍
  concernedList:[], // 关注队伍
  isContainsAfterContestJudge: false, // 是否包含比赛结束后的提交
  contest: {
    auth: CONTEST_TYPE.PUBLIC,
    openPrint: false,
    rankShowName:'username',
    allowEndSubmit: false,
  },
  contestProblems: [],
  itemVisible: {
    table: true,
    chart: true,
  },
  disPlayIdMapColor:{}, // 展示id对应的气球颜色
  groupContestAuth: 0,
})