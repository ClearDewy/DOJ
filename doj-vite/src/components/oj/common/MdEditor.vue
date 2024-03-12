<template>
	<MdPreview
		v-if="props.previewOnly"
		:modelValue="modelValue"
		previewTheme="vuepress"
		codeTheme="github"
		:showCodeRowNumber="true"
	/>
	<MdEditor
		v-else ref="editorRef"
		noIconfont
		:tabWidth="4"
		:editorId="editorId"
		:modelValue="modelValue"
		previewTheme="vuepress"
    codeTheme="github"
		@update:modelValue="handleModelValueUpdate"
		:showCodeRowNumber="true"
		:toolbarsExclude="toolbarsExclude"
	/>
</template>

<script setup lang="ts">
//@ts-ignore
import {ExposeParam, MdEditor, MdPreview} from "md-editor-v3";
import 'md-editor-v3/lib/style.css';
import {ref} from "vue";

const props=defineProps({
	editorId:String,
	modelValue:String,
	previewOnly:{
		type:Boolean,
		default:false
	}
})
const emits=defineEmits(['update:modelValue'])
const handleModelValueUpdate = (newValue:any) => {
	emits('update:modelValue',newValue)
}
const editorRef = ref<ExposeParam>();

const toolbarsExclude=['save','github']

</script>


<style scoped>
.md-editor-previewOnly:hover{
  box-shadow: var(--el-box-shadow-light);  /* 这是一个浅色的阴影效果 */
	transition: box-shadow 0.3s ease-in-out;  /* 平滑的阴影过渡 */
}
</style>
<style>

p {
  margin:0 0 10px 0 !important;
}
</style>
