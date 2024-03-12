import storage from "./index";
import {reactive, watch} from "vue";
import {CK} from "../common/constants";

export type userInfoType={
    username:string,
    school:string,
    major:string,
    number:string,
    name:string,
    gender:number|undefined,
    cf_username:string,
    email:string,
    avatar:string,
    signature:string,
    title_name:string,
    title_color:string,

    system_auth:number,
    user_auth:number,
    problem_auth:number,
    context_auth:number,
    train_auth:number,
    problem_status_auth:number,

    submit_auth:number,
    context_attend_auth:number,
    train_attend_auth:number,
}

export type userStateType={
    userInfo:userInfoType,
    token:string,
    loginFailNum:number,
}

export const initUserState:userStateType={
    userInfo:{
        username:"",
        school:"",
        major:"",
        number:"",
        name:"",
        gender:undefined,
        cf_username:"",
        email:"",
        avatar:"",
        signature:"",
        title_name:"",
        title_color:"",

        system_auth:0,
        user_auth:0,
        problem_auth:0,
        context_auth:0,
        train_auth:0,
        problem_status_auth:0,

        submit_auth:1,
        context_attend_auth:1,
        train_attend_auth:1,
    },
    token: "",
    loginFailNum:0,
}

if (!storage.get(CK.USER_STATE)){
    storage.set(CK.USER_STATE,initUserState)
}

// storage.set(CK.USER_STATE,initUserState)

export const user_state = reactive<userStateType>(storage.get(CK.USER_STATE))

watch(user_state,(new_user_state)=>{
    storage.set(CK.USER_STATE,new_user_state)
})
