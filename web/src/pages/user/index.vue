<script setup lang="ts">
import type { ILocationInfo, IUserProfile } from '@/interface'
import { defineRouteMeta, useRouter } from '@fesjs/fes'
import {
  FButton,
  FForm,
  FFormItem,
  FImage,
  FInput,
  FOption,
  FPagination,
  FSelect,
  FTable,
  FTableColumn,
  FTag,
} from '@fesjs/fes-design'
import { LoadingOutlined } from '@fesjs/fes-design/icon'
import { reactive, ref, watch } from 'vue'
import { useRequest } from 'vue-hooks-plus'
import { request } from '@/api'
import { formatTimestamp, fromBirthToAge } from '@/common/utils'
import {
  AGE_RANGE,
  BODY_SHAPE,
  EDUCATION_LEVEL,
  GENDER,
  HEIGHT_RANGE,
  INCOME_TYPE,
  MBTI,
  YES_NO,
} from '@/enums'

defineRouteMeta({
  name: 'user/list',
  title: '用户资料',
})

const pageState = reactive({
  current_page: 1,
  page_size: 10,
  total: 0,
})

const defaultValue = {
  id: '',
  name: '',
  gender: 0,
  birth_max: 0,
  birth_min: 0,
  height_max: 0,
  height_min: 0,
  is_stable: 0,
  education_level: 0,
  job_id: 0,
  income_type: 0,
  height: 1,
  body_type: 0,
  location_city_id: 0,
  home_location: '',
  user_upload_tag: '',
}
const searchState = reactive(defaultValue)
const router = useRouter()

function reset() {
  router.go(0)
}

function reqUserProfile() {
  const data = {
    id: +searchState.id,
    name: searchState.name,
    gender: searchState.gender,
    is_stable: searchState.is_stable,
    birth_min: searchState.birth_min,
    birth_max: searchState.birth_max,
    height_min: searchState.height_min,
    height_max: searchState.height_max,
    education_level: searchState.education_level,
    job_id: searchState.job_id,
    income_type: searchState.income_type,
    body_type: searchState.body_type,
    location_city_id: searchState.location_city_id,
    user_upload_tag: '',
    personality_type: 0,
    home_location: searchState.home_location,
    page: {
      current_page: pageState.current_page,
      page_count: 0,
      page_size: pageState.page_size, // 获取全部活动
      total: 0,
    },
  }
  if (MBTI[searchState.user_upload_tag as 'INTP']) {
    data.personality_type = MBTI[searchState.user_upload_tag as 'INTP']
  }
  else {
    data.user_upload_tag = searchState.user_upload_tag
  }
  return request('/user_profile/list', data, { method: 'post' })
}

const {
  loading,
  data,
  run: getUserProfileList,
} = useRequest<{ list: IUserProfile[], page: { total: number } }>(
  reqUserProfile,
)

const loadingOnce = ref(loading.value)
const formRef = ref<typeof FForm>()
watch(
  () => loading.value,
  val => (loadingOnce.value = val),
  { once: true },
)

const action = [
  {
    label: '编辑',
    func: (row: IUserProfile) => {
      router.push(`/user/${row.id}`)
    },
  },
]
function handleChange(page: number, pageSize: number) {
  pageState.current_page = page
  pageState.page_size = pageSize
  getUserProfileList()
}

interface JobConfig {
  id: number
  name: string
  title: string
}

const { data: jobConfig } = useRequest<JobConfig[]>(() =>
  request('/user/job_config'),
)

const { data: JapanCityList } = useRequest<{ id: number, name: string, type: number }[]>(() =>
  request('/config/location_info', { parent_id: 3025 }, { method: 'POST' }).then(res => res.list),
)
// const { data: tagList } = useRequest<JobConfig[]>(() => request('/tag/list'))

function formatUserProfileTag(userProfile: IUserProfile) {
  // personality_type
  let res: string[] = []
  if (userProfile.personality_type) {
    res.push(MBTI[userProfile.personality_type])
  }
  if (userProfile.experience_tags?.length) {
    res = [...res, ...userProfile.experience_tags.map(t => t.title)]
  }
  if (userProfile.interest_tags?.length) {
    res = [...res, ...userProfile.interest_tags.map(t => t.title)]
  }
  if (userProfile.love_expectation_tags?.length) {
    res = [
      ...res,
      ...userProfile.love_expectation_tags.map(t => t.title),
    ]
  }
  if (userProfile.life_style_tags?.length) {
    res = [...res, ...userProfile.life_style_tags.map(t => t.title)]
  }
  return res
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
function handleAgeRangeChange(val: 0 | AGE_RANGE) {
  if (val === 0) {
    searchState.birth_min = 0
    searchState.birth_max = 0
  }
  else if (val === AGE_RANGE['80前']) {
    searchState.birth_min = 0
    searchState.birth_max = Math.floor(new Date('1980-1-1').getTime() / 1000)
  }
  else if (val === AGE_RANGE['80-84']) {
    searchState.birth_min = Math.floor(new Date('1980-1-1').getTime() / 1000)
    searchState.birth_max = Math.floor(new Date('1985-1-1').getTime() / 1000) - 1
  }
  else if (val === AGE_RANGE['85-89']) {
    searchState.birth_min = Math.floor(new Date('1985-1-1').getTime() / 1000)
    searchState.birth_max = Math.floor(new Date('1990-1-1').getTime() / 1000) - 1
  }
  else if (val === AGE_RANGE['90-94']) {
    searchState.birth_min = Math.floor(new Date('1990-1-1').getTime() / 1000)
    searchState.birth_max = Math.floor(new Date('1995-1-1').getTime() / 1000) - 1
  }
  else if (val === AGE_RANGE['95-99']) {
    searchState.birth_min = Math.floor(new Date('1995-1-1').getTime() / 1000)
    searchState.birth_max = Math.floor(new Date('2000-1-1').getTime() / 1000) - 1
  }
  else if (val === AGE_RANGE['00-04']) {
    searchState.birth_min = Math.floor(new Date('2000-1-1').getTime() / 1000)
    searchState.birth_max = Math.floor(new Date('2005-1-1').getTime() / 1000) - 1
  }
  else if (val === AGE_RANGE['05后']) {
    searchState.birth_min = Math.floor(new Date('2005-1-1').getTime() / 1000)
    searchState.birth_max = 0
  }
}
function handleHeightRangeChange(val: 0 | HEIGHT_RANGE) {
  if (val === 0) {
    searchState.height_min = 0
    searchState.height_max = 0
  }
  else if (val === HEIGHT_RANGE['150及以下']) {
    searchState.height_min = 0
    searchState.height_max = 150
  }
  else if (val === HEIGHT_RANGE['151-160']) {
    searchState.height_min = 151
    searchState.height_max = 160
  }
  else if (val === HEIGHT_RANGE['161-170']) {
    searchState.height_min = 161
    searchState.height_max = 170
  }
  else if (val === HEIGHT_RANGE['171-175']) {
    searchState.height_min = 171
    searchState.height_max = 175
  }
  else if (val === HEIGHT_RANGE['176-180']) {
    searchState.height_min = 176
    searchState.height_max = 180
  }
  else if (val === HEIGHT_RANGE['181-185']) {
    searchState.height_min = 181
    searchState.height_max = 185
  }
  else if (val === HEIGHT_RANGE['185以上']) {
    searchState.height_min = 185
    searchState.height_max = 0
  }
}
function handleLocationCityChange(val: string) {
  const cityId = JapanCityList.value?.find(city => city.name === val)?.id
  if (val) {
    searchState.location_city_id = cityId || 9999
  }
  else {
    searchState.location_city_id = 0
  }
}
</script>

<template>
  <nav>
    <h1>用户资料</h1>
    <div>
      <FForm
        ref="formRef" :model="data" label-position="right" :span="12" align="flex-start"
        class="user-profile-search-form" @keydown.enter="getUserProfileList"
      >
        <FFormItem prop="id" label="ID:">
          <FInput v-model="searchState.id" placeholder="请输入ID" />
        </FFormItem>
        <FFormItem prop="name" label="姓名:">
          <FInput v-model="searchState.name" placeholder="请输入姓名" />
        </FFormItem>
        <FFormItem prop="gender" label="性别:">
          <FSelect v-model="searchState.gender">
            <FOption :value="0">
              全部
            </FOption>
            <FOption :value="1">
              男性
            </FOption>
            <FOption :value="2">
              女性
            </FOption>
          </FSelect>
        </FFormItem>
        <FFormItem prop="age" label="年龄段:">
          <FSelect :model-value="0" @change="handleAgeRangeChange">
            <FOption :value="0">
              全部
            </FOption>
            <FOption
              v-for="(id) in Object.keys(AGE_RANGE).filter((k: any) => isNaN(+(AGE_RANGE[k])))"
              :key="id"
              :value="+id"
            >
              {{ AGE_RANGE[+id] }}
            </FOption>
          </FSelect>
        </FFormItem>
        <FFormItem prop="education" label="学历:">
          <FSelect v-model="searchState.education_level">
            <FOption :value="0">
              全部
            </FOption>
            <FOption
              v-for="(id) in Object.keys(EDUCATION_LEVEL).filter((k) => isNaN(+(EDUCATION_LEVEL[k as any])))"
              :key="id" :value="+id"
            >
              {{ EDUCATION_LEVEL[+id] }}
            </FOption>
          </FSelect>
        </FFormItem>

        <FFormItem prop="location_city_id" label="城市:">
          <FInput placeholder="请输入城市" @input="handleLocationCityChange" />
        </FFormItem>
        <FFormItem prop="hometown" label="家乡:">
          <FInput v-model="searchState.home_location" placeholder="请输入家乡" />
        </FFormItem>

        <FFormItem prop="is_stable" label="是否永驻:">
          <FSelect v-model="searchState.is_stable">
            <FOption :value="0">
              全部
            </FOption>
            <FOption
              v-for="(id) in Object.keys(YES_NO).filter((k) => isNaN(+(YES_NO[+k])))" :key="id"
              :value="+id"
            >
              {{ YES_NO[+id] }}
            </FOption>
          </FSelect>
        </FFormItem>
        <FFormItem prop="job" label="职业:">
          <FSelect v-model="searchState.job_id" class="select">
            <FOption :value="0">
              全部
            </FOption>
            <FOption v-for="{ id, name } in jobConfig" :key="id" :value="id">
              {{ name }}
            </FOption>
          </FSelect>
        </FFormItem>
        <FFormItem prop="income" label="年收:">
          <FSelect v-model="searchState.income_type" class="select">
            <FOption :value="0">
              全部
            </FOption>
            <FOption
              v-for="(id) in Object.keys(INCOME_TYPE).filter((k: any) => isNaN(+(INCOME_TYPE[k])))" :key="id"
              :value="+id"
            >
              {{ INCOME_TYPE[+id] }}
            </FOption>
          </FSelect>
        </FFormItem>
        <FFormItem prop="height" label="身高:">
          <FSelect :model-value="0" class="select" @change="handleHeightRangeChange">
            <FOption :value="0">
              全部
            </FOption>
            <FOption
              v-for="(id) in Object.keys(HEIGHT_RANGE).filter((k: any) => isNaN(+(HEIGHT_RANGE[k])))" :key="id"
              :value="+id"
            >
              {{ HEIGHT_RANGE[+id] }}
            </FOption>
          </FSelect>
        </FFormItem>

        <FFormItem prop="body_type" label="体型:">
          <FSelect v-model="searchState.body_type" class="select">
            <FOption :value="0">
              全部
            </FOption>
            <FOption
              v-for="(id) in Object.keys(BODY_SHAPE).filter((k: any) => isNaN(+(BODY_SHAPE[k])))" :key="id"
              :value="+id"
            >
              {{ BODY_SHAPE[+id] }}
            </FOption>
          </FSelect>
        </FFormItem>

        <FFormItem prop="tags" label="上传标签:">
          <FInput v-model="searchState.user_upload_tag" placeholder="请输入用户上传的标签" />
        </FFormItem>

        <FFormItem style="float: right" label=" ">
          <FButton type="primary" @click="getUserProfileList">
            查询
          </FButton>
          &nbsp;&nbsp;&nbsp;
          <FButton @click="reset">
            重置
          </FButton>
        </FFormItem>
      </FForm>
    </div>
  </nav>
  <div v-if="loading" class="loading">
    <LoadingOutlined class="icon" />
  </div>
  <FTable
    v-show="!loading" always-scrollbar :height="0" class="table" size="small" row-key="id"
    :data="data?.list ?? []"
    row-class-name="32r32r"
  >
    <FTableColumn fixed="left" prop="id" label="ID" :min-width="50" />
    <FTableColumn fixed="left" prop="name" label="姓名" />
    <FTableColumn fixed="left" label="性别">
      <template #default="{ row }">
        {{ GENDER[row.gender] }}
      </template>
    </FTableColumn>
    <FTableColumn label="匹配性别">
      <template #default="{ row }">
        {{ GENDER[row.find_gender] }}
      </template>
    </FTableColumn>
    <FTableColumn prop="birth" label="年龄" sortable>
      <template #default="{ row }">
        {{ fromBirthToAge(row.birth) }}
      </template>
    </FTableColumn>
    <FTableColumn prop="school" label="毕业学校" />
    <FTableColumn prop="education_level" label="学历">
      <template #default="{ row }">
        {{ EDUCATION_LEVEL[row.education_level] }}
      </template>
    </FTableColumn>
    <FTableColumn prop="images" label="照片">
      <template #default="{ row }">
        <FImage
          v-for="(img, i) of row.images" :key="i" class="preview-image" preview fit="contain" hide-on-click-modal
          :src="img"
        />
      </template>
    </FTableColumn>
    <FTableColumn label="城市">
      <template #default="{ row }">
        {{ row.location_city_info?.name ?? "-" }}
      </template>
    </FTableColumn>
    <FTableColumn label="家乡">
      <template #default="{ row }">
        {{ locationDisplay(row) }}
      </template>
    </FTableColumn>
    <FTableColumn label="是否永驻">
      <template #default="{ row }">
        {{ YES_NO[row.is_stable] ?? "-" }}
      </template>
    </FTableColumn>
    <FTableColumn prop="job_name" label="职业" />
    <FTableColumn prop="income_type" label="年收" sortable>
      <template #default="{ row }">
        {{ INCOME_TYPE[row.income_type] }}
      </template>
    </FTableColumn>
    <FTableColumn prop="height" label="身高(cm)">
      <template #default="{ row }">
        {{ row.height || "-" }}
      </template>
    </FTableColumn>
    <FTableColumn prop="body_type" label="体型">
      <template #default="{ row }">
        {{ BODY_SHAPE[row.body_type] ?? "-" }}
      </template>
    </FTableColumn>
    <FTableColumn prop="email" label="邮箱" />
    <FTableColumn label="上传标签" sortable>
      <template #default="{ row }">
        <FTag v-for="tag in formatUserProfileTag(row)" :key="tag" type="info">
          {{ tag }}
        </FTag>
      </template>
    </FTableColumn>
    <FTableColumn prop="find_body_types" label="对方城市">
      <template #default="{ row }">
        {{
          row.find_location_city_ids
            ?.map((e: ILocationInfo) => e.name)
            .join(" ") || "-"
        }}
      </template>
    </FTableColumn>
    <FTableColumn label="对方年龄">
      <template #default="{ row }">
        {{ row.find_age_min || "" }}-{{ row.find_age_max || "" }}
      </template>
    </FTableColumn>
    <FTableColumn label="对方身高(cm)">
      <template #default="{ row }">
        {{ row.find_height_min || "" }}-{{ row.find_height_max || "" }}
      </template>
    </FTableColumn>
    <FTableColumn prop="find_body_types" label="对方体型">
      <template #default="{ row }">
        {{
          row.find_body_types
            ?.map((e: number) => BODY_SHAPE[e])
            ?.join(" ") || "-"
        }}
      </template>
    </FTableColumn>
    <FTableColumn prop="find_education_levels" label="对方学历">
      <template #default="{ row }">
        {{
          row.find_education_levels
            ?.map((e: number) => EDUCATION_LEVEL[e])
            ?.join(" ") || "-"
        }}
      </template>
    </FTableColumn>
    <!-- <FTableColumn prop="remark" label="补充信息" /> -->
    <FTableColumn prop="tags" label="活动标签">
      <template #default="{ row }">
        <FTag v-for="tag in row.activity_tags" :key="tag" type="info">
          {{ tag.title }}
        </FTag>
      </template>
    </FTableColumn>
    <FTableColumn prop="operation_remark" label="运营备注">
      <template #default="{ row }">
        {{ row.operation_remark || "-" }}
      </template>
    </FTableColumn>
    <FTableColumn :min-width="163" prop="create_time" label="创建时间">
      <template #default="{ row }">
        {{ formatTimestamp(row.create_time * 1000) }}
      </template>
    </FTableColumn>
    <FTableColumn :min-width="163" prop="update_time" label="最近更新时间">
      <template #default="{ row }">
        {{ formatTimestamp(row.update_time * 1000) }}
      </template>
    </FTableColumn>
    <FTableColumn label="操作" align="center" fixed="right" :action="action" />
  </FTable>
  <FPagination
    v-if="!loadingOnce" show-total class="pagination" :total-count="data?.page?.total"
    show-size-changer
    show-quick-jumper
    @change="handleChange"
  />
</template>

<style scoped lang="less">
.loading {
  width: 100%;
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;

  .icon {
    margin-top: 20px;
    font-size: 40px;
  }
}

.preview-image {
  width: 30px;
  height: 30px;
  margin-right: 3px;
  display: inline-block;
  overflow: hidden;
  border: 1px solid #333;
  border-radius: 3px;
}

.pagination {
  margin-top: 10px;
  align-self: center;
}

nav {
  margin-bottom: 20px;
}

.table {
  flex: 1;
  display: flex;
  flex-direction: column;
  :global(.table .fes-table-body-wrapper) {
    flex: 1;
  }
}

.user-profile-search-form {
  :global(& .fes-form-item) {
    // display: inline-block;
    display: inline-flex;
    flex-direction: row;

    &+& {
      margin-right: 24px;
    }
  }

  :global(& .fes-input) {
    width: 150px;
  }

  :global(& .fes-select) {
    width: 100px;
  }
}
</style>
