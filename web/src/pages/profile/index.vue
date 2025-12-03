<script setup lang="ts">
import type { IUserProfile } from '@/interface'
import { defineRouteMeta } from '@fesjs/fes'
import { FButton, FForm, FFormItem, FGrid, FGridItem, FInput, FMessage, FRadio, FRadioGroup, FSpace } from '@fesjs/fes-design'
import { LoadingOutlined } from '@fesjs/fes-design/icon'
import { reactive, ref } from 'vue'
import { useRequest } from 'vue-hooks-plus'
import { request } from '@/api'
import { useValidator } from '@/common/hooks'
import Upload from '@/components/upload.vue'

const props = defineProps<{
  loading?: boolean
  data?: Record<string, any>
  readonly?: boolean
  onSubmit?: (data: any) => void
  onCancel?: () => void
}>()

defineRouteMeta({
  name: 'profile',
  title: '个人中心',
  layout: {
    // navigation: null,
  },
})

const data = reactive<Partial<IUserProfile>>({
  userId: 0,
  nickname: '32',
  email: '',
  gender: 1,
  image: [],
})

const { loading } = useRequest(() => {
  return request('/profile')
}, {
  onSuccess(res) {
    Object.keys(res).forEach((k) => {
      if (k === 'image') {
        data[k] = res?.[k] ? [res?.[k]] : []
        return
      }
      Object.assign(data, res)
    })
  },
})

const validator = useValidator(data)
const formRef = ref<typeof FForm>()
async function handleSubmit() {
  await formRef.value?.validate()
  await request('/profile', {
    ...data,
    image: data.image?.toString() || '',
  }, { method: 'put' })
  FMessage.success('保存成功')
}
</script>

<template>
  <div>
    <LoadingOutlined v-if="loading" class="loading" />
    <FForm
      v-if="!loading" ref="formRef" :model="data" label-position="right" :label-width="100" :span="12"
      align="flex-start"
    >
      <h2 class="sub-title">
        个人信息
      </h2>

      <FFormItem prop="image" :rules="[{ required: true, type: 'string', message: '请上传头像', validator }]" label="头像">
        <Upload v-model="data.image" :count="1" />
      </FFormItem>

      <FFormItem label="姓名" prop="nickname" :rules="[{ validator, required: true, type: 'string', message: '请输入姓名' }]">
        <FInput v-model="data.nickname" :maxlength="10" placeholder="请输入" />
      </FFormItem>

      <FGrid :gutter="20">
        <FGridItem :span="8">
          <FFormItem
            label="性别" prop="gender"
            :rules="[{ required: true, type: 'string', message: '请输入性别', validator }]"
          >
            <FRadioGroup v-model="data.gender" :cancelable="false">
              <FRadio :value="1">
                男
              </FRadio>
              <FRadio :value="2">
                女
              </FRadio>
            </FRadioGroup>
          </FFormItem>
        </FGridItem>
      </FGrid>

      <FFormItem
        label="邮箱" prop="email" :rules="[
          { required: true, type: 'string', message: '请输入邮箱', validator },
          { required: true, type: 'email', message: '请输入正确的邮箱格式', validator },
        ]"
      >
        <FInput v-model="data.email" type="email" placeholder="请输入" />
      </FFormItem>

      <FFormItem label=" ">
        <FSpace>
          <FButton type="primary" @click="handleSubmit">
            保存
          </FButton>
          <FButton type="default" @click="props.onCancel">
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
  overflow: hidden;
  height: 100px;
  border: 1px solid #333;
  border-radius: 5px;
}

.sub-title {
  color: #000;
  font-size: 20px;
  font-weight: 600;
  line-height: 24px;
  margin-top: 0px;
}
</style>
