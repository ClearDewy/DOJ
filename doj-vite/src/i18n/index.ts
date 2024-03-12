import {createI18n, useI18n} from "vue-i18n";
import zhCN from "./zh-CN";
import enUS from "./en-US";

const i18n=createI18n<typeof zhCN, 'en-US' | 'zh-CN'>({
    locale:"zh-CN",
    fallbackLocale:"en-US",
    messages:{
        "en-US":enUS,
        "zh-CN":zhCN
    },
    legacy:false,
    globalInjection:true
})

export default i18n

export const changeWebLanguage=()=>{
    const locale=useI18n().locale
    locale.value = locale.value === 'zh-CN' ? 'en-US' : 'zh-CN'
}