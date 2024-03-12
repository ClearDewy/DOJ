export const JUDGE_STATUS = {
  '-4': {
    name: 'Judging',
    color: 'blue',
    type: '',
    rgb:'#2d8cf0'
  },
  '-3':{
    name: 'Compiling',
    short: 'CP',
    color: 'green',
    type: 'info',
    rgb:'#25bb9b'
  },
  '-2': {
    name: 'Pending',
    color: 'yellow',
    type: 'warning',
    rgb:'#f90'
  },
  '-1': {
    name: 'Submitting',
    color: 'yellow',
    type: 'warning',
    rgb:'#f90'
  },
  '1': {
    name: 'Accepted',
    short: 'AC',
    color: 'green',
    type: 'success',
    rgb:'#19be6b'
  },
  '2': {
    name: 'Wrong Answer',
    short: 'WA',
    color: 'red',
    type: 'error',
    rgb:'#ed3f14'
  },

  '3': {
    name: 'Time Limit Exceeded',
    short: 'TLE',
    color: 'red',
    type: 'error',
    rgb:'#ed3f14'
  },
  '4': {
    name: 'Memory Limit Exceeded',
    short: 'MLE',
    color: 'red',
    type: 'error',
    rgb:'#ed3f14'
  },
  '5': {
    name: 'Compile Error',
    short: 'CE',
    color: 'yellow',
    type: 'warning',
    rgb:'#f90'
  },
  '6': {
    name: 'Runtime Error',
    short: 'RE',
    color: 'red',
    type: 'error',
    rgb:'#ed3f14'
  },
  '7': {
    name: 'Partial Accepted',
    short: 'PAC',
    color: 'blue',
    type: '',
    rgb:'#2d8cf0'
  },
  '8': {
    name: 'System Error',
    short: 'SE',
    color: 'gray',
    type: 'info',
    rgb:'#909399'
  },
  '9': {
    name: 'Presentation Error',
    short: 'PE',
    color: 'yellow',
    type: 'warning',
    rgb:'#f90'
  },
  '10':{
    name:"Submitted Failed",
    color:'gray',
    short:'SF',
    type: 'info',
    rgb:'#909399',
  },
  '11': {
    name: 'Not Submitted',
    short: 'NS',
    color: 'gray',
    type: 'info',
    rgb:'#909399'
  },
  '12': {
    name: 'Cancelled',
    short: 'CA',
    color: 'purple',
    type: 'info',
    rgb:'#676fc1'
  },
  '13': {
    name: 'Submitted Unknown Result',
    short: 'SNR',
    color: 'gray',
    type: 'info',
    rgb:'#909399'
  },
}

export const JUDGE_STATUS_RESERVE={
  "Judging":-4,
  "Compiling":-3,
  "Pending":-2,
  "Submitting":-1,
  "Accepted":1,
  "Wrong Answer":2,
  "Time Limit Exceeded":3,
  "Memory Limit Exceeded":4,
  "Compile Error":5,
  "Runtime Error":6,
  "Partial Accepted":7,
  "System Error":8,
  "Presentation Error":9,
  "Submitted Failed":10,
  "Not Submitted":11,
  "Cancelled":12,
  "Submitted Unknown Result":13
}

export const PROBLEM_LEVEL={
  '1':{
    name:{
      'zh-CN':'未知',
      'en-US':'Unkonwn',
    },
    color:'#5f5f5f'
  },
  '2':{
    name:{
      'zh-CN':'入门',
      'en-US':'Entry',
    },
    color:'#19be6b'
  },
  '3':{
    name:{
      'zh-CN':'普及',
      'en-US':'Universal',
    },
    color:'#ffc116'
  },
  '4':{
    name:{
      'zh-CN':'提高',
      'en-US':'Develop',
    },
    color:'#3498d8'
  },
  '5':{
    name:{
      'zh-CN':'铜牌',
      'en-US':'BronzeMedal',
    },
    color:'#fe4c61'
  },
  '6':{
    name:{
      'zh-CN':'银牌',
      'en-US':'SilverMedal',
    },
    color:'#9d3dcf'
  },
  '7':{
    name:{
      'zh-CN':'金牌',
      'en-US':'GoldMedal',
    },
    color:'#0e1d69'
  },
}


export const REMOTE_OJ = [
  {
    name:'洛谷',
    key:"LuoGu"
  },
  {
    name:'HDU',
    key:"HDU"
  },
  {
    name:"Codeforces",
    key:"CF"
  },
  {
    name:"POJ",
    key:"POJ"
  },
  {
    name:"GYM",
    key:"GYM"
  },
  {
    name:"AtCoder",
    key:"AC"
  },
  {
    name:"SPOJ",
    key:"SPOJ"
  }
]

export const CONTEST_STATUS = {
  'SCHEDULED': -1,
  'RUNNING': 0,
  'ENDED': 1
}

export const CONTEST_STATUS_REVERSE = {
  '-1': {
    name: 'Scheduled',
    color: '#f90'
  },
  '0': {
    name: 'Running',
    color: '#19be6b'
  },
  '1': {
    name: 'Ended',
    color: '#ed3f14'
  }
}

export const TRAINING_TYPE = {
  'Public':{
    color:'success',
    name:'Public'
  },
  'Private':{
    color:'danger',
    name:'Private'
  }
}

export const GROUP_TYPE = {
  PUBLIC: 1,
  PROTECTED: 2,
  PRIVATE: 3
}

export const GROUP_TYPE_REVERSE = {
  '1':{
    name: 'Public',
    color: 'success',
    tips: 'Group_Public_Tips',
  },
  '2':{
    name: 'Protected',
    color: 'warning',
    tips: 'Group_Protected_Tips',
  },
  '3':{
    name: 'Private',
    color: 'danger',
    tips: 'Group_Private_Tips',
  }
}

export const RULE_TYPE = {
  ACM: 0,
  OI: 1
}

export const CONTEST_TYPE_REVERSE = {
  '0': {
    name:'Public',
    color:'success',
    tips:'Public_Tips',
    submit:true,              // 公开赛可看可提交
    look:true,
  },
  '1':{
    name:'Private',
    color:'danger',
    tips:'Private_Tips',
    submit:false,         // 私有赛 必须要密码才能看和提交
    look:false,
  },
  '2':{
    name:'Protected',
    color:'warning',
    tips:'Protected_Tips',
    submit:false,       //保护赛，可以看但是不能提交，提交需要附带比赛密码
    look:true,
  }
}

export const CONTEST_TYPE = {
  PUBLIC: 0,
  PRIVATE: 1,
  PROTECTED: 2
}

export const USER_TYPE = {
  REGULAR_USER: 'user',
  ADMIN: 'admin',
  PROBLEM_ADMIN:'problem_admin',
  SUPER_ADMIN: 'root'
}

export const JUDGE_CASE_MODE = {
  DEFAULT: 'default',
  SUBTASK_LOWEST: 'subtask_lowest',
  SUBTASK_AVERAGE: 'subtask_average',
  ERGODIC_WITHOUT_ERROR: 'ergodic_without_error'
}

export const FOCUS_MODE_ROUTE_NAME = {
  'TrainingFullProblemDetails': 'TrainingProblemDetails',
  'ContestFullProblemDetails': 'ContestProblemDetails',
  'GroupFullProblemDetails':'GroupProblemDetails',
  'GroupTrainingFullProblemDetails': 'GroupTrainingProblemDetails'
}

// constant_key
export const CK = {
  AUTHED: 'authed',
  PROBLEM_CODE_AND_SETTING: 'hojProblemCodeAndSetting',
  languages: 'languages',
  CONTEST_ANNOUNCE:'hojContestAnnounce',
  individualLanguageAndSetting:'hojIndividualLanguageAndSetting',
  CONTEST_RANK_CONCERNED:'hojContestRankConcerned',
  AUTHORIZATION:"authorization",
  WEB_LANGUAGE:"WebLanguage",
  USER_INFO:"userInfo",
  TOKEN:"token",
  USER_STATE:"user_state"
}

export function buildIndividualLanguageAndSettingKey () {
  return `${CK.individualLanguageAndSetting}`
}

export function buildProblemCodeAndSettingKey (problemID:number, contestID = null) {
  if (contestID) {
    return `${CK.PROBLEM_CODE_AND_SETTING}_${contestID}_${problemID}`
  }
  return `${CK.PROBLEM_CODE_AND_SETTING}_NoContest_${problemID}`
}

export function buildContestAnnounceKey (uid:number, contestID:number) {
  return `${CK.CONTEST_ANNOUNCE}_${uid}_${contestID}`
}

export function buildContestRankConcernedKey(contestID:number) {
  return `${CK.CONTEST_RANK_CONCERNED}_${contestID}`
}

