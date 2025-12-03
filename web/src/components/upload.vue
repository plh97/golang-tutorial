<script setup lang="ts">
import type { FileItem } from '@fesjs/fes-design/es/upload/interface'
import { FImage, FMessage, FSpace, FUpload } from '@fesjs/fes-design'
import { CloseCircleFilled } from '@fesjs/fes-design/icon'
import { ref } from 'vue'
import { request } from '@/api'

const props = defineProps({
  disabled: Boolean,
  modelValue: Array,
  count: Number,
  tip: Boolean,
})

const emit = defineEmits(['update:modelValue'])

const count = props.count ?? 6

const inputRef = ref()

const maxSize = ref(5) // MB

const fileList = ref<FileItem[]>(props?.modelValue?.map(url => ({ url })) as FileItem[])

async function handleUpload({ file, onProgress, onSuccess }: { file: File, onProgress: (progress: { percent: number }) => void, onSuccess: (url: string) => void }) {
  const fileExt = file.name.split('.').pop() || 'png'
  // const formData = new FormData()
  // formData.append('file_ext', fileExt)
  // formData.append('upload_scene', '1')
  // formData.append('file', file)
  const url = await request('/common/upload', {
    file_ext: fileExt,
    upload_scene: 1,
  }, { method: 'post' })
  await fetch(url.pre_signed_url, {
    method: 'put',
    headers: {
      'Content-Type': file.type,
    },
    // data: file,
    body: file,
  })
  onProgress({
    percent: 100,
  })
  onSuccess(url.endpoint_url)
}

function handleSuccess(res: { fileList: FileItem[] }) {
  fileList.value = res.fileList.map((file) => {
    if (file?.response) {
      return {
        uid: file.uid || Date.now().toString(),
        name: file.name || 'image',
        url: file.response,
      }
    }
    return file
  })
  emit('update:modelValue', fileList.value.map((file) => {
    return file.url
  }))
}
function handleDelete(url: string) {
  const newUrl = fileList.value.filter((file) => {
    return file.url !== url
  })
  fileList.value = newUrl
  emit('update:modelValue', newUrl.map((file) => {
    return file.url
  }))
}
function beforeUpload(file: File) {
  if (file.size > maxSize.value * 1024 * 1024) {
    FMessage.warn(`图片不能超过${maxSize.value}MB`)
    return false
  }
  return true
}
</script>

<template>
  <FUpload
    v-model:file-list="fileList"
    :before-upload="beforeUpload"
    :http-request="handleUpload" :accept="['image/*']" :disabled="props.disabled"
    @success="handleSuccess"
  >
    <button ref="inputRef" class="hide">
      +
    </button>
    <template #fileList="{ uploadFiles }">
      <div disabled :class="props.disabled ? 'disabled' : ''">
        <FSpace justify="start" vertical size="xsmall">
          <FSpace justify="start" align="center" wrap>
            <div v-for="(value, i) in uploadFiles" :key="i" class="list-container">
              <FImage class="upload-img" :src="value.url" :preview="!props.disabled" hide-on-click-modal fit="contain" />
              <CloseCircleFilled v-if="!props.disabled" class="icon" @click="handleDelete(value.url)" />
            </div>
            <div v-if="uploadFiles.length < count" class="upload-trigger" @click="!props.disabled && inputRef.click()">
              <span class="icon">+</span>
              <span>点击上传图片</span>
            </div>
          </FSpace>
          <div v-if="props.tip" class="f-upload__tip">
            大小不超过<span class="red">{{ maxSize }}MB</span>，尺寸为<span class="red">686X384</span>
          </div>
        </FSpace>
      </div>
    </template>
  </FUpload>
</template>

<style scoped lang="less">
:global(.fes-upload) {
  display: none;
}
.hide {
  display: none;
  height: 0px;
}

.upload-img {
  width: 100px;
  overflow: hidden;
  height: 100px;
  border: 1px solid #CED4DB;
  border-radius: 5px;
  width: 128px;
  height: 128px;
  cursor: pointer;
}

.upload-trigger {
  display: flex;
  border: 1px solid #000;
  border-radius: 5px;
  width: 128px;
  height: 128px;
  cursor: pointer;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  user-select: none;
  background-color: #F5F7FA;
  border: 1px dashed #CED4DB;

  .icon {
    font-size: 50px;
    margin-bottom: 20px;
  }
}

.list-container {
  position: relative;
  user-select: none;

  .icon {
    position: absolute;
    top: 0;
    right: 0;
    cursor: pointer;
    font-size: 20px;
    color: #333;
    transform: translate(50%, -50%);
  }
}
.f-upload__tip {
  color: #8D9094;
  font-size: 12px;
  font-weight: 400;
  line-height: 20px;
}
.red {
  color: #FA0658;
}
</style>
