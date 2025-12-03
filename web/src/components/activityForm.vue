<script setup lang="ts">
import type { IDomEditor, IEditorConfig } from '@wangeditor/editor'
import { useRoute } from '@fesjs/fes'
import { FButton, FDatePicker, FForm, FFormItem, FInput, FInputNumber, FOption, FSelect, FSpace } from '@fesjs/fes-design'
import { LoadingOutlined } from '@fesjs/fes-design/icon'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import { onBeforeUnmount, reactive, ref, shallowRef, watch } from 'vue'
import { request } from '@/api'
import { useValidator } from '@/common/hooks'
import { ACTIVITY_STATUS } from '@/enums'
import Upload from './upload.vue'
import '@wangeditor/editor/dist/css/style.css'

const props = defineProps({
  loading: Boolean,
  data: Object,
  readonly: Boolean,
  onSubmit: Function,
  onCancel: Function,
})

const route = useRoute()

const isEdit = route.fullPath.includes('edit')

const data = reactive({
  id: -1,
  title: '',
  desc: '',
  location: '',
  man_sign_up_count: 0,
  woman_sign_up_count: 0,
  man_sign_up_fake_count: 0,
  woman_sign_up_fake_count: 0,
  man_sign_up_limit: 0,
  woman_sign_up_limit: 0,
  banner_img: '',
  man_price: 0,
  man_price_jp: 0,
  woman_price: 0,
  woman_price_jp: 0,
  start_time: Math.floor(Date.now() / 1000),
  status: 0,
  end_time: Math.floor(Date.now() / 1000),
  create_time: 0,
  update_time: 0,
})

const localImageList = ref<string[]>([])
// 监听 props 的 obj 变化，同步到本地状态
watch(
  () => props.data,
  (newValue) => {
    Object.keys(newValue ?? {}).forEach((k) => {
      // eslint-disable-next-line ts/ban-ts-comment
      // @ts-expect-error
      data[k] = newValue?.[k]
    })
    if (data.banner_img) {
      localImageList.value = [data.banner_img]
    }
  },
  { immediate: true, deep: true }, // immediate 确保初始同步，deep 监听对象深层变化
)

const format = ref()
const formRef = ref<typeof FForm>()
const editorToolRef = shallowRef()
const toolbarConfig = {}
const editorConfig: Partial<IEditorConfig> = {
  placeholder: '请输入内容...',
  MENU_CONF: {
    uploadImage: {
      async customUpload(file: File, insertFn: any) {
        const formData = new FormData()
        formData.append('upload_scene', '1')
        formData.append('file', file)
        const url = await request('/common/upload', formData, {
          method: 'post',
        })
        insertFn(url, 'baidu logo', url)
      },
    },
  },
}

// 组件销毁时，也及时销毁编辑器
onBeforeUnmount(() => {
  const editor = editorToolRef.value
  if (editor == null) {
    return
  }
  editor.destroy()
})
function handleCreated(editor: IDomEditor) {
  editorToolRef.value = editor // 记录 editor 实例，重要！
  if (props.readonly) {
    editorToolRef.value.disable()
  }
}
function handleDateChange([start, end]: number[]) {
  data.start_time = Math.floor(start / 1000)
  data.end_time = Math.floor(end / 1000)
}
async function handleSubmit() {
  await formRef.value?.validate()
  if (!data.desc || data.desc === '<p><br></p>') {
    return
  }
  if (data.status === 3) {
    // @ts-ignore
    delete data.woman_sign_up_fake_count
    // @ts-ignore
    delete data.man_sign_up_fake_count
  }
  if (isEdit) {
    // @ts-ignore
    delete data.woman_sign_up_fake_count
    // @ts-ignore
    delete data.man_sign_up_fake_count
    // @ts-ignore
    delete data.woman_sign_up_limit
    // @ts-ignore
    delete data.man_sign_up_limit
  }
  props.onSubmit?.(data)
}
watch(
  () => localImageList.value,
  ([url]) => {
    data.banner_img = url
  },
)
const validator = useValidator(data)
</script>

<template>
  <div>
    <LoadingOutlined v-if="props.loading" class="loading" />
    <FForm
      v-show="!props.loading"
      ref="formRef"
      :disabled="props.readonly" :model="data" label-position="right" :label-width="100" :span="12"
      align="flex-start"
    >
      <FFormItem
        prop="banner_img" :rules="[{ required: true, type: 'number', message: '请上传图片', validator }]"
        label="封面"
      >
        <Upload :key="data.banner_img" v-model="localImageList" tip :disabled="props.readonly" :count="1" />
      </FFormItem>
      <FFormItem label="活动标题" prop="title" :rules="[{ required: true, type: 'string', message: '请输入活动标题', validator }]">
        <FInput v-model="data.title" placeholder="请输入" :maxlength="20" />
      </FFormItem>
      <FFormItem
        label="活动地点" prop="location"
        :rules="[{ required: true, type: 'string', message: '请输入活动地点', validator }]"
      >
        <FInput v-model="data.location" placeholder="请输入" :maxlength="100" />
      </FFormItem>
      <FFormItem
        label="活动时间" :prop="JSON.stringify(['start_time', 'end_time'])"
        :rules="[{ required: true, type: 'array', message: '请输入活动时间', validator }]"
      >
        <FDatePicker
          :min-date="new Date('2025-01-01')" :format="format" type="datetimerange"
          :model-value="[data.start_time * 1000, data.end_time * 1000]" clearable @change="handleDateChange"
        >
          <template #separator>
            至
          </template>
        </FDatePicker>
      </FFormItem>
      <FFormItem
        label="活动状态" prop="status"
        :rules="[{ required: true, type: 'string', message: '请输入活动状态', validator }]"
      >
        <FSelect v-model="data.status">
          <FOption
            v-for="(id) in Object.keys(ACTIVITY_STATUS).filter(key => isNaN(+(ACTIVITY_STATUS[+key])))" :key="id"
            :value="+id"
          >
            {{ ACTIVITY_STATUS[+id] }}
          </FOption>
        </FSelect>
      </FFormItem>
      <FFormItem
        label="活动价格"
        :prop="JSON.stringify(['man_price_jp', 'woman_price_jp', 'man_price', 'woman_price'])"
        :rules="[{ required: true, type: 'array', message: '请输入活动价格', validator }]"
      >
        <div class="container">
          <div class="grid">
            <span>日元</span>
            <span>男</span>
            <FInputNumber v-model="data.man_price_jp" :min="0" placeholder="请输入" />
            <span>女</span>
            <FInputNumber v-model="data.woman_price_jp" :min="0" placeholder="请输入" />
          </div>
          <div class="grid">
            <span>人民币</span>
            <span>男</span>
            <FInputNumber v-model="data.man_price" :min="0" placeholder="请输入" />
            <span>女</span>
            <FInputNumber v-model="data.woman_price" :min="0" placeholder="请输入" />
          </div>
        </div>
      </FFormItem>
      <FFormItem
        label="活动人数"
        :prop="JSON.stringify(['man_sign_up_limit', 'woman_sign_up_limit'])"
        :rules="[{ required: true, type: 'array', message: '请输入活动人数', validator }]"
      >
        <!-- <FInputNumber v-model="data.sign_up_limit" :min="0" placeholder="请输入" /> -->
        <div class="container">
          <div class="grid">
            <br>
            <span>男</span>
            <FInputNumber v-model="data.man_sign_up_limit" :disabled="isEdit" :min="0" placeholder="请输入" />
            <span>女</span>
            <FInputNumber v-model="data.woman_sign_up_limit" :disabled="isEdit" :min="0" placeholder="请输入" />
          </div>
        </div>
      </FFormItem>
      <FFormItem label="活动介绍" prop="desc" :rules="[{ required: true, type: 'string', message: '请输入活动介绍', validator }]">
        <div class="editor-container">
          <Toolbar
            style="border-bottom: 1px solid #ccc" :editor="editorToolRef" :default-config="toolbarConfig"
            mode="simple"
          />
          <Editor
            v-model="data.desc"
            style="height: 200px; overflow-y: hidden;"
            :default-config="editorConfig"
            mode="simple" @on-created="handleCreated"
          />
        </div>
        <div v-if="!data.desc || data.desc === '<p><br></p>'" class="fes-form-item-error">
          请输入活动介绍
        </div>
      </FFormItem>
      <FFormItem label="虚拟人数" prop="sign_up_limit">
        <div class="container">
          <div class="grid">
            <br>
            <span>男</span>
            <FInputNumber v-model="data.man_sign_up_fake_count" :disabled="isEdit" :min="0" placeholder="请输入" />
            <span>女</span>
            <FInputNumber v-model="data.woman_sign_up_fake_count" :disabled="isEdit" :min="0" placeholder="请输入" />
          </div>
        </div>
      </FFormItem>
      <FFormItem label=" ">
        <FSpace>
          <FButton v-if="props.onSubmit" type="primary" @click="handleSubmit()">
            提交
          </FButton>
          <FButton v-if="props.onCancel" type="default" @click="props.onCancel()">
            取消
          </FButton>
        </FSpace>
      </FFormItem>
    </FForm>
  </div>
</template>

<style scoped lang="less">
.loading {
  width: 100%;
  margin-top: 20px;
  font-size: 40px;
}

.container {
  display: flex;
  flex-direction: column;
  gap: 14px;

  .grid {
    align-items: center;
    display: grid;
    grid-template-columns: 60px 40px 184px 40px 144px;
  }
}

.upload-img {
  width: 100px;
  height: 100px;
  border: 1px solid #333;
  overflow: hidden;
  border-radius: 5px;
}
:global(input:disabled, select:disabled) {
  color: var(--f-text-color-disabled) !important;
}
.editor-container {
  border: 1px solid #ccc;
  border-radius: 5px;
  overflow: hidden;
}
</style>
