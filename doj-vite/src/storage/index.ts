// 封装本地存储器
import {reactive} from "vue";
import {CK} from "../common/constants";

export const storage={
    get(key:string){
        return JSON.parse(window.localStorage.getItem(key) as string)||null
    },
    set(key:string,value:any){
        window.localStorage.setItem(key,JSON.stringify(value))
    },
    remove(key:string){
        window.localStorage.removeItem(key)
    },
    clear(){
        window.localStorage.clear()
    }
}

export default storage

export const web_state = reactive({
    modalStatus: {
        mode: 'login'||'register', // or 'register',
        visible: false
    },
    websiteConfig:{
        recordName:'© 2022-2024',
        projectName:'Dewy Online Judge',
        shortName:'DOJ',
        recordUrl:'#',
        projectUrl:'#',
        introduction:'',
        allowRegister:true,
    },
    registerTimeOut: 60,
    resetTimeOut: 90,
    language:storage.get(CK.WEB_LANGUAGE) || 'zh-CN',
})