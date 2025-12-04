<script setup lang="ts">
import type { IUserProfile } from '@/interface'
import { FButton, FDatePicker, FForm, FFormItem, FGrid, FGridItem, FInput, FMessage, FOption, FRadio, FRadioGroup, FSelect, FSpace } from '@fesjs/fes-design'
import { LoadingOutlined } from '@fesjs/fes-design/icon'
import { reactive, ref, watch } from 'vue'
import { useRequest } from 'vue-hooks-plus'
import { request } from '@/api'
import { useValidator } from '@/common/hooks'
import { BODY_SHAPE, EDUCATION_LEVEL, INCOME_TYPE, MBTI } from '@/enums'
import Upload from './upload.vue'

const props = defineProps<{
  loading?: boolean
  data?: Record<string, any>
  readonly?: boolean
  onSubmit?: (data: any) => void
  onCancel?: () => void
}>()

const data = reactive<Partial<IUserProfile>>({
  id: 0,
  user_id: 0,
  name: '',
  birth: 0,
  education_level: 0,
  gender: 0,
  find_gender: 0,
  school: '',
  job_id: 0,
  job_name: '',
  income_type: 0,
  email: '',
  intro: '',
  operation_remark: '',
  images: [],
  tags: [],
  height: 0,
  body_type: 0,
  personality_type: 0,
  location_city_info: null,
  income_display: 0,
  life_intro: '',
  hobby_intro: '',
  life_style_intro: '',
  find_location_city_ids: null,
  find_location_city_info: null,
  find_age_min: 0,
  find_age_max: 0,
  find_height_min: 0,
  find_height_max: 0,
  find_body_types: null,
  find_education_levels: null,
  find_other_intro: '',
  experience_tags: null,
  interest_tags: null,
  love_expectation_tags: null,
  life_style_tags: null,
  activity_tags: [],
  required_fields: [],
  is_stable: undefined,
  home_location_country_info: null,
  home_location_state_info: null,
  home_location_city_info: null,
  home_location_city_id: 0,
  home_location_country_id: 0,
  home_location_state_id: 0,
  location_city_id: 0,
})

// 监听 props 的 obj 变化，同步到本地状态
watch(
  () => props.data,
  (newValue) => {
    Object.keys(newValue ?? {}).forEach((k) => {
      if (k === 'activity_tags') {
        data.activity_tags = newValue?.activity_tags?.map((tag: { title: string }) => tag.title) || []
      }
      else {
        // eslint-disable-next-line ts/ban-ts-comment
        // @ts-ignore
        data[k] = newValue?.[k]
      }
    })
  },
  { immediate: true, deep: true }, // immediate 确保初始同步，deep 监听对象深层变化
)

function handleSelectChange() {
  if (data.activity_tags?.find(t => t?.title?.length > 5)) {
    FMessage.error({
      content: '标签长度不能超过5个字符',
    })
    data.activity_tags = data.activity_tags.filter(t => t?.title?.length <= 5)
  }
}

const format = ref()

interface JobConfig {
  id: number
  name: string
  title: string
}

const { data: jobConfig } = useRequest<JobConfig[]>(() => request('/user/job_config'))
const { data: tagList } = useRequest<JobConfig[]>(async () => {
  let res = await request('/tag/list')
  const map = {}
  res = res?.map(e => ({ label: e.title, value: e.title }))
  res = res.filter((item) => {
    if (!map[item.label]) {
      map[item.label] = item.label
      return true
    }
    return false
  })
  return res
})

function handleJobChange(id: number) {
  const jobItem = jobConfig?.value?.find((job) => {
    return job.id === id
  })
  if (!jobItem) {
    return
  }
  data.job_name = id === 1 ? '' : jobItem?.name
}

const validator = useValidator(data)
const formRef = ref<typeof FForm>()
async function handleSubmit() {
  await formRef.value?.validate()
  props.onSubmit?.(data)
}

function handleBirthChange(date: number) {
  data.birth = Math.floor(date / 1000)
}
function formatUserProfileTag(userProfile: IUserProfile) {
  // personality_type
  let res: string[] = []
  if (userProfile.personality_type) {
    res.push(MBTI[userProfile.personality_type])
  }
  // experience_tags
  // interest_tags
  // love_expectation_tags
  // life_style_tags
  if (userProfile.experience_tags?.length) {
    res = [...res, ...userProfile.experience_tags.map(t => t.title)]
  }
  if (userProfile.interest_tags?.length) {
    res = [...res, ...userProfile.interest_tags.map(t => t.title)]
  }
  if (userProfile.love_expectation_tags?.length) {
    res = [...res, ...userProfile.love_expectation_tags.map(t => t.title)]
  }
  if (userProfile.life_style_tags?.length) {
    res = [...res, ...userProfile.life_style_tags.map(t => t.title)]
  }
  return res.join(' ') || '-'
}
function locationDisplay(user: IUserProfile) {
  let homeLocation = ''
  if (user?.home_location_country_info?.id !== 4164) {
    homeLocation += user?.home_location_country_info?.name ?? ''
  }
  homeLocation += user?.home_location_state_info?.name ?? ''
  homeLocation += user?.home_location_city_info?.name ?? ''
  return homeLocation || '-'
}
</script>

<template>
  <div>
    <LoadingOutlined v-if="props.loading" class="loading" />
    <FForm
      v-show="!props.loading" ref="formRef" :model="data" label-position="right" :label-width="100" :span="12"
      align="flex-start"
    >
      <h2 class="sub-title">
        个人信息
      </h2>
      <FFormItem label="姓名" prop="name" :rules="[{ validator, required: true, type: 'string', message: '请输入姓名' }]">
        <FInput v-model="data.name" :maxlength="10" :disabled="props.readonly" placeholder="请输入" />
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
        <FGridItem :span="8">
          <FFormItem
            label="出生日期" prop="birth"
            :rules="[{ required: true, type: 'number', message: '请输入出生日期', validator }]"
          >
            <FDatePicker
              :model-value="data.birth! * 1000" :format="format" :disabled="props.readonly" type="date"
              @change="handleBirthChange"
            />
          </FFormItem>
        </FGridItem>
      </FGrid>

      <!-- <FFormItem label="出生日期" prop="birth" :rules="[{ required: true, type: 'number', message: '请输入出生日期', validator }]">
        <FDatePicker
          :model-value="data.birth! * 1000" :format="format" :disabled="props.readonly" type="date"
          @change="handleBirthChange"
        />
      </FFormItem> -->
      <FFormItem prop="images" :rules="[{ required: true, type: 'string', message: '请上传照片', validator }]" label="照片">
        <Upload :key="data.images?.length" v-model="data.images" :disabled="props.readonly" />
      </FFormItem>

      <FGrid :gutter="20">
        <FGridItem :span="8">
          <FFormItem
            prop="school" label="毕业学校"
            :rules="[{ required: true, type: 'string', message: '请输入毕业学校', validator }]"
          >
            <FInput
              v-model="data.school" type="school" :maxlength="100" :disabled="props.readonly"
              placeholder="请输入毕业学校"
            />
          </FFormItem>
        </FGridItem>
        <FGridItem :span="8">
          <FFormItem
            prop="school" :label-width="124" label="学历"
            :rules="[{ required: true, type: 'string', message: '请输入学历' }]"
          >
            <FSelect v-model="data.education_level" class="select">
              <FOption
                v-for="(id) in Object.keys(EDUCATION_LEVEL).filter((k) => isNaN(+(EDUCATION_LEVEL[k as any])))"
                :key="id" :value="+id"
              >
                {{ EDUCATION_LEVEL[id as any] }}
              </FOption>
            </FSelect>
          </FFormItem>
        </FGridItem>
        <!-- <FGridItem :span="8">
          <FFormItem label="身高(cm)" :rules="[{ required: true, type: 'string', message: '请输入身高' }]">
            <FInputNumber v-model="data.height" placeholder="请输入" />
          </FFormItem>
        </FGridItem> -->
      </FGrid>

      <FGrid :gutter="20">
        <FGridItem :span="8">
          <FFormItem label="城市" prop="id" :rules="[{ required: true, type: 'string' }]">
            <FInput :model-value="data.location_city_info?.name" disabled placeholder="请输入城市" />
          </FFormItem>
        </FGridItem>
        <FGridItem :span="8">
          <FFormItem label="家乡" prop="id" :rules="[{ required: true, type: 'string' }]">
            <FInput :model-value="locationDisplay(data as IUserProfile)" disabled placeholder="请输入家乡" />
          </FFormItem>
        </FGridItem>
        <FGridItem :span="8">
          <FFormItem label="是否有永驻" prop="id" :rules="[{ required: true, type: 'string' }]">
            <FRadioGroup v-model="data.is_stable" disabled>
              <FRadio :value="1">
                是
              </FRadio>
              <FRadio :value="2">
                否
              </FRadio>
            </FRadioGroup>
          </FFormItem>
        </FGridItem>
      </FGrid>

      <FGrid :gutter="20">
        <FGridItem :span="8">
          <FFormItem
            label="职业" :prop="JSON.stringify(['job_id', 'job_name'])"
            :rules="[{ required: true, type: 'array', message: '请输入职业', validator }]"
          >
            <FSelect v-model="data.job_id" class="select" @change="handleJobChange">
              <FOption v-for="({ id, name }) in jobConfig" :key="id" :value="id">
                {{ name }}
              </FOption>
            </FSelect>
            <FInput v-if="data.job_id === 1" v-model="data.job_name" style="margin-left: 5px;" />
          </FFormItem>
        </FGridItem>
        <FGridItem :span="8">
          <FFormItem
            label-width="124" label="年收(万日元)" prop="income_type"
            :rules="[{ required: true, type: 'string', message: '请输入身高', validator }]"
          >
            <FSelect v-model="data.income_type" class="select">
              <FOption
                v-for="(id) in Object.keys(INCOME_TYPE).filter((k: any) => isNaN(+(INCOME_TYPE[k])))" :key="id"
                :value="+id"
              >
                {{ INCOME_TYPE[id as any] }}
              </FOption>
            </FSelect>
          </FFormItem>
        </FGridItem>

        <FGridItem :span="8">
          <FFormItem label="是否对外展示年收" prop="id" :rules="[{ required: true, type: 'number' }]">
            <FRadioGroup v-model="data.income_display" disabled>
              <FRadio :value="1">
                是
              </FRadio>
              <FRadio :value="2">
                否
              </FRadio>
            </FRadioGroup>
          </FFormItem>
        </FGridItem>
      </FGrid>

      <FGrid :gutter="20">
        <FGridItem :span="8">
          <FFormItem label="身高" prop="id" :rules="[{ required: true, type: 'number' }]">
            <FInput :model-value="`${data.height}cm`" disabled placeholder="请输入身高" />
          </FFormItem>
        </FGridItem>
        <FGridItem :span="8">
          <FFormItem label="体型" prop="id" :rules="[{ required: true, type: 'number' }]">
            <FInput :model-value="BODY_SHAPE[data.body_type!]" disabled placeholder="请输入体型" />
          </FFormItem>
        </FGridItem>
      </FGrid>

      <FFormItem
        label="邮箱" prop="email" :rules="[
          { required: true, type: 'string', message: '请输入邮箱', validator },
          { required: true, type: 'email', message: '请输入正确的邮箱格式', validator },
        ]"
      >
        <FInput v-model="data.email" type="email" disabled placeholder="请输入" />
      </FFormItem>

      <FFormItem label="自我介绍">
        <FInput v-model="data.intro" show-word-limit :maxlength="500" type="textarea" disabled placeholder="请输入" />
      </FFormItem>

      <FFormItem label="上传标签">
        <FInput
          :model-value="formatUserProfileTag(data as IUserProfile)" show-word-limit :maxlength="500"
          type="textarea" disabled placeholder="请输入"
        />
      </FFormItem>

      <FFormItem label="生活近况与家庭背景">
        <FInput v-model="data.life_intro" show-word-limit :maxlength="500" type="textarea" disabled placeholder="请输入" />
      </FFormItem>

      <FFormItem label="性格与兴趣爱好">
        <FInput
          v-model="data.hobby_intro" show-word-limit :maxlength="500" type="textarea" disabled
          placeholder="请输入"
        />
      </FFormItem>

      <FFormItem label="生活方式与习惯">
        <FInput
          v-model="data.life_style_intro" show-word-limit :maxlength="500" type="textarea" disabled
          placeholder="请输入"
        />
      </FFormItem>

      <h2 class="sub-title">
        期望对方信息
      </h2>

      <FFormItem
        label="匹配性别" prop="find_gender"
        :rules="[{ validator, required: true, type: 'string', message: '请输入姓名' }]"
      >
        <FRadioGroup v-model="data.find_gender" disabled :cancelable="false">
          <FRadio :value="1">
            男
          </FRadio>
          <FRadio :value="2">
            女
          </FRadio>
        </FRadioGroup>
      </FFormItem>

      <FGrid :gutter="20">
        <FGridItem :span="6">
          <FFormItem label="对方年龄">
            <FInput disabled :model-value="`${data.find_age_min}-${data.find_age_max}岁`" />
          </FFormItem>
        </FGridItem>
        <FGridItem :span="6">
          <FFormItem label="对方身高">
            <FInput disabled :model-value="`${data.find_height_min}-${data.find_height_max}cm`" />
          </FFormItem>
        </FGridItem>
        <FGridItem :span="6">
          <FFormItem label="对方体型">
            <FInput :model-value="data.find_body_types?.map(t => BODY_SHAPE[t]).join(' ')" disabled />
          </FFormItem>
        </FGridItem>
        <FGridItem :span="6">
          <FFormItem label="对方学历">
            <FInput :model-value="data.find_education_levels?.map(t => EDUCATION_LEVEL[t]).join(' ')" disabled />
          </FFormItem>
        </FGridItem>
      </FGrid>

      <FFormItem label="对方城市">
        <FInput :model-value="data.find_location_city_ids?.map(t => t.name).join(' ')" disabled />
      </FFormItem>

      <FFormItem label="其他期望">
        <FInput
          v-model="data.find_other_intro" show-word-limit :maxlength="500" type="textarea" disabled
          placeholder="请输入"
        />
      </FFormItem>

      <h2 class="sub-title">
        运营编辑
      </h2>
      <FFormItem label="活动展示标签">
        <FSelect
          v-model="data.activity_tags!" tag filterable filter-text-highlight multiple :multiple-limit="3" class="select"
          :options="tagList"
          @change="handleSelectChange"
        />
      </FFormItem>
      <FFormItem label="运营备注">
        <FInput
          v-model="data.operation_remark" show-word-limit :maxlength="100" type="textarea"
          :disabled="props.readonly" placeholder="请输入"
        />
      </FFormItem>
      <FFormItem label=" ">
        <FSpace>
          <FButton v-if="props.onSubmit" type="primary" @click="handleSubmit">
            提交
          </FButton>
          <FButton v-if="props.onCancel" type="default" @click="props.onCancel">
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
