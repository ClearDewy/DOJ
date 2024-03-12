import {TRAINING_TYPE} from "../common/constants";
import {reactive} from "vue";

export const train_state = reactive({
    intoAccess: true, // 训练进入权限
    training: {
        auth: TRAINING_TYPE.Public.name,
        rankShowName:'username',
        gid:null
    },
    trainingProblemList: [],
    itemVisible: {
        table: true,
        chart: true,
    },
    groupTrainingAuth: 0,
})