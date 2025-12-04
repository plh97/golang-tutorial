<script setup>
import { defineRouteMeta, useRoute, useRouter } from '@fesjs/fes'
import { FButton, FDivider, FEllipsis, FForm, FFormItem, FImage, FInput, FInputNumber, FMessage, FModal, FOption, FPagination, FSelect, FSpace, FTable, FTableColumn } from '@fesjs/fes-design'
import { LoadingOutlined } from '@fesjs/fes-design/icon'
import { computed, reactive, ref, watch } from 'vue'
import { useRequest } from 'vue-hooks-plus'
import { request } from '@/api'
import { useValidator } from '@/common/hooks'
import { formatArray, formatTimestamp, fromBirthToAge } from '@/common/utils'
import Upload from '@/components/upload.vue'
import { EDUCATION_LEVEL, GENDER, INCOME_TYPE, PAY_METHOD, SIGN_UP_PAY_STATUS, SIGN_UP_STATUS, YES_NO } from '@/enums'

defineRouteMeta({
  name: 'activity/signup',
  title: '报名列表',
})

const { params } = useRoute()
const changeSignUpStatusModalForm = ref()
const state = reactive({
  activeItemLocal: {},
  activeItem: {},
  remarkModalVisible: false,
  paymentModalVisible: false,
  statusModalVisible: false,
  selectIDs: [],
})

const pageState = reactive({
  current_page: 1,
  page_size: 10,
  total: 0,
  status: 0,
  is_pay: 0,
  is_return: 0,
})

const { loading, data, run } = useRequest(() => request('/sign_up/list', {
  activity_id: +params.id,
  status: pageState.status ? pageState.status : undefined,
  is_pay: pageState.is_pay ? pageState.is_pay : undefined,
  is_return: pageState.is_return ? pageState.is_return : undefined,
  page: {
    current_page: pageState.current_page,
    page_count: 0,
    page_size: pageState.page_size, // 获取全部活动
    total: 0,
  },
}, {
  method: 'POST',
  transformData: (data) => {
    // 处理响应内容异常
    if (data?.code !== 0) {
      // Reject the promise with an error object containing code and message
      FMessage.error({
        content: data?.message,
      })
    }
    return {
      list: data?.data?.list?.map((item) => {
        return {
          ...item,
          remark: item?.order_sn_json?.remark,
          pay_method: item?.order_sn_json?.pay_type,
          pay_price: item?.order_sn_json?.pay_amount,
          pay_proof: item?.order_sn_json?.pay_proof,
          refund_method: item?.order_sn_json?.refund_type,
          refund_price: item?.order_sn_json?.refund_amount,
          refund_proof: item?.order_sn_json?.refund_proof,
        }
      }),
      page: data?.data?.page,
    }
  },
}))

const loadingOnce = ref(loading.value)

watch(() => loading.value, (val) => {
  loadingOnce.value = val
}, {
  once: true,
})

function reset() {
  pageState.status = 0
  pageState.is_pay = 0
  pageState.is_return = 0
  run()
}

const router = useRouter()
const action = [
  {
    label: '备注',
    func: (row) => {
      state.activeItem = row
      state.activeItemLocal = { ...row }
      state.remarkModalVisible = true
    },
    onConfirm: () => {
      request('/sign_up/update', {
        id: state.activeItemLocal?.id,
        remark: state.activeItemLocal?.remark,
      }, {
        method: 'POST',
      }).then(() => {
        FMessage.success({
          content: '备注成功',
        })
        run()
      }).finally(() => {
        state.activeItem = {}
        state.activeItemLocal = {}
        state.remarkModalVisible = false
      })
    },
  },
  {
    label: '编辑付款/退款',
    func: (row) => {
      state.activeItem = row
      state.activeItemLocal = { ...row }
      state.activeItemLocal.pay_proof_list = row.pay_proof ? row.pay_proof.split(',') : []
      state.activeItemLocal.refund_proof_list = row.refund_proof ? row.refund_proof.split(',') : []
      state.paymentModalVisible = true
    },
    onConfirm: () => {
      request('/sign_up/update', {
        id: state.activeItemLocal?.id,
        pay_type: state.activeItemLocal?.pay_method,
        pay_amount: state.activeItemLocal?.pay_price,
        pay_proof: state.activeItemLocal.pay_proof_list.join(','),
        refund_type: state.activeItemLocal?.refund_method,
        refund_amount: state.activeItemLocal?.refund_price,
        refund_proof: state.activeItemLocal?.refund_proof_list.join(','),
      }, {
        method: 'POST',
      }).then(() => {
        FMessage.success({
          content: '编辑付款/退款成功',
        })
        run()
      }).finally(() => {
        state.activeItem = {}
        state.activeItemLocal = {}
        state.paymentModalVisible = false
      })
    },
  },
]
function handleBack() {
  router.push(`/activity`)
}
const calculateTableHeight = computed(() => {
  // Get viewport height
  const viewportHeight = window.innerHeight
  return viewportHeight - 390
})
function handleChange(page, pageSize) {
  pageState.current_page = page
  pageState.page_size = pageSize
  run()
}
async function handleChangeStatus() {
  await changeSignUpStatusModalForm.value?.validate()
  const list = data.value?.list ?? []
  // 应该过滤掉没有变化的
  const ids = state.selectIDs
    .map(id => list.find(i => i.id === id))
    .filter(signUp => signUp.status !== state.status)
    .map(i => (i.id))
  const genders = state.selectIDs
    .map(id => list.find(i => i.id === id))
    .filter(signUp => signUp.status !== state.status)
    .map(i => (i?.user_profile?.gender))
  request('/sign_up/batch_update', {
    activity_id: +params.id,
    sign_up_ids: ids,
    genders,
    status: state.status,
  }, {
    method: 'post',
  }).then(() => {
    FMessage.success({
      content: '修改成功',
    })
    state.statusModalVisible = false
    state.selectIDs = []
    run()
  })
}
const validator = useValidator(state)
</script>

<template>
  <!-- remark modal -->
  <FModal v-model:show="state.remarkModalVisible" title="添加备注" display-directive="show" @ok="action[0].onConfirm">
    <FForm label-position="right" :label-width="100" :span="12" align="flex-start">
      <FFormItem label="备注">
        <FInput v-model="state.activeItemLocal.remark" :maxlength="100" show-word-limit type="textarea" />
      </FFormItem>
    </FForm>
  </FModal>
  <!-- remark modal -->
  <!-- payment modal -->
  <FModal v-model:show="state.paymentModalVisible" title="编辑付款/退款" display-directive="show" @ok="action[1].onConfirm">
    <FForm label-position="right" :label-width="100" :span="12" align="flex-start">
      <FFormItem label="付款方式">
        <FSelect v-model="state.activeItemLocal.pay_method">
          <FOption :value="0">
            -
          </FOption>
          <FOption
            v-for="(id) in Object.keys(PAY_METHOD).filter(key => isNaN(+(PAY_METHOD[key])))" :key="id"
            :value="+id"
          >
            {{ PAY_METHOD[id] }}
          </FOption>
        </FSelect>
      </FFormItem>
      <FFormItem label="付款价格">
        <FInputNumber v-model="state.activeItemLocal.pay_price" :min="0" />
      </FFormItem>
      <FFormItem label="上传付款凭证">
        <Upload :key="state.activeItemLocal.pay_proof" v-model="state.activeItemLocal.pay_proof_list" />
      </FFormItem>

      <FFormItem label="退款方式">
        <FSelect v-model="state.activeItemLocal.refund_method">
          <FOption :value="0">
            -
          </FOption>
          <FOption
            v-for="(id) in Object.keys(PAY_METHOD).filter(key => isNaN(+(PAY_METHOD[key])))" :key="id"
            :value="+id"
          >
            {{ PAY_METHOD[id] }}
          </FOption>
        </FSelect>
      </FFormItem>
      <FFormItem label="退款价格">
        <FInputNumber v-model="state.activeItemLocal.refund_price" :min="0" />
      </FFormItem>
      <FFormItem label="上传退款凭证">
        <Upload :key="state.activeItemLocal.refund_proof" v-model="state.activeItemLocal.refund_proof_list" />
      </FFormItem>
    </FForm>
  </FModal>
  <!-- payment modal -->
  <!-- mul-change sign up status modal -->
  <FModal v-model:show="state.statusModalVisible" title="更改报名状态" display-directive="show" @ok="handleChangeStatus">
    <FForm
      ref="changeSignUpStatusModalForm" :model="state" label-position="right" :label-width="100" :span="12"
      align="flex-start"
    >
      <FFormItem label="状态" prop="status" :rules="[{ required: true, type: 'string', message: '请选择状态', validator }]">
        <FSelect v-model="state.status">
          <FOption
            v-for="(id) in Object.keys(SIGN_UP_STATUS).filter(key => isNaN(+(SIGN_UP_STATUS[key])))" :key="id"
            :value="+id"
          >
            {{ SIGN_UP_STATUS[id] }}
          </FOption>
        </FSelect>
      </FFormItem>
    </FForm>
  </FModal>
  <!-- mul-change sign up status modal -->
  <nav>
    <FSpace justify="space-between" align="center">
      <h1>
        报名列表
      </h1>
      <FButton type="primary" @click="handleBack">
        返回
      </FButton>
    </FSpace>
    <FForm>
      <FSpace justify="space-between" align="center" :wrap="false">
        <FSpace justify="space-between" align="center" :wrap="false">
          <FFormItem label="报名状态：" style="width: 200px">
            <FSelect v-model="pageState.status">
              <FOption :value="0">
                全部
              </FOption>
              <FOption
                v-for="(id) in Object.keys(SIGN_UP_STATUS).filter(key => isNaN(+(SIGN_UP_STATUS[key])))"
                :key="id" :value="+id"
              >
                {{ SIGN_UP_STATUS[id] }}
              </FOption>
            </FSelect>
          </FFormItem>
          <FFormItem label="付款状态：" style="width: 200px">
            <FSelect v-model="pageState.is_pay">
              <FOption :value="0">
                全部
              </FOption>
              <FOption v-for="(id) in Object.keys(YES_NO).filter(key => isNaN(+(YES_NO[key])))" :key="id" :value="+id">
                {{ YES_NO[id] }}
              </FOption>
            </FSelect>
          </FFormItem>
          <FFormItem label="是否退款：" style="width: 200px">
            <FSelect v-model="pageState.is_return">
              <FOption :value="0">
                全部
              </FOption>
              <FOption v-for="(id) in Object.keys(YES_NO).filter(key => isNaN(+(YES_NO[key])))" :key="id" :value="+id">
                {{ YES_NO[id] }}
              </FOption>
            </FSelect>
          </FFormItem>
        </FSpace>
        <FSpace justify="space-between" style="width: 170px;" align="center" :wrap="false">
          <FFormItem label=" ">
            <FButton @click="pageState.current_page = 0; run()">
              查询
            </FButton>
          </FFormItem>
          <FFormItem label=" ">
            <FButton @click="reset">
              重置
            </FButton>
          </FFormItem>
        </FSpace>
      </FSpace>
    </FForm>
    <FButton :disabled="state.selectIDs.length === 0" @click="state.statusModalVisible = true">
      更改报名状态
    </FButton>
  </nav>
  <FDivider />
  <div v-if="loading" class="loading">
    <LoadingOutlined class="icon" />
  </div>
  <FTable
    v-show="!loading" v-model:checked-keys="state.selectIDs" always-scrollbar :height="calculateTableHeight"
    class="table" size="small" row-key="id" :data="data?.list ?? []"
  >
    <FTableColumn type="selection" :width="30" fixed="left" />
    <FTableColumn prop="id" label="ID" :min-width="50" fixed="left" />
    <FTableColumn label="姓名" fixed="left">
      <template #default="{ row }">
        <FEllipsis :content="row.user_profile?.name" style="max-width: 100%" />
      </template>
    </FTableColumn>
    <FTableColumn label="性别">
      <template #default="{ row }">
        {{ GENDER[row.user_profile?.gender] }}
      </template>
    </FTableColumn>
    <FTableColumn label="匹配性别">
      <template #default="{ row }">
        {{ GENDER[row.user_profile?.find_gender] }}
      </template>
    </FTableColumn>
    <FTableColumn prop="age" label="年龄">
      <template #default="{ row }">
        {{ fromBirthToAge(row.user_profile?.birth) }}
      </template>
    </FTableColumn>
    <FTableColumn prop="images" label="照片">
      <template #default="{ row }">
        <FImage
          v-for="(img, i) of formatArray(row.user_profile?.images)" :key="i" class="preview-image" preview
          fit="contain" hide-on-click-modal :src="img"
        />
      </template>
    </FTableColumn>
    <FTableColumn label="学历">
      <template #default="{ row }">
        {{ EDUCATION_LEVEL[row.user_profile?.education_level] }}
      </template>
    </FTableColumn>
    <FTableColumn label="毕业学校">
      <template #default="{ row }">
        {{ row.user_profile?.school }}
      </template>
    </FTableColumn>
    <FTableColumn label="邮箱">
      <template #default="{ row }">
        {{ row.user_profile?.email }}
      </template>
    </FTableColumn>
    <!-- <FTableColumn label="身高">
        <template #default="{ row }">
          {{ row.user_profile.height }}
        </template>
      </FTableColumn> -->
    <FTableColumn label="职业">
      <template #default="{ row }">
        {{ row.user_profile?.job_name }}
      </template>
    </FTableColumn>
    <FTableColumn label="年收">
      <template #default="{ row }">
        {{ INCOME_TYPE[row.user_profile?.income_type] }}
      </template>
    </FTableColumn>
    <FTableColumn label="报名状态">
      <template #default="{ row }">
        {{ SIGN_UP_STATUS[row.status] }}
      </template>
    </FTableColumn>
    <FTableColumn label="付款状态">
      <template #default="{ row }">
        {{ SIGN_UP_PAY_STATUS[row.pay_price > 0 ? 2 : 1] }}
      </template>
    </FTableColumn>
    <FTableColumn label="付款方式">
      <template #default="{ row }">
        {{ PAY_METHOD[row.pay_method] }}
      </template>
    </FTableColumn>
    <FTableColumn label="付款金额">
      <template #default="{ row }">
        {{ row.pay_price }}
      </template>
    </FTableColumn>
    <FTableColumn prop="pay_proof" label="付款凭证">
      <template #default="{ row }">
        <FSpace v-if="row.pay_proof">
          <FImage
            v-for="url in row.pay_proof.split(',')" :key="url" class="proof-image" preview hide-on-click-modal
            :src="url"
          />
        </FSpace>
      </template>
    </FTableColumn>
    <FTableColumn label="是否退款">
      <template #default="{ row }">
        {{ row.refund_price > 0 ? '是' : '否' }}
      </template>
    </FTableColumn>
    <FTableColumn prop="return_amount" label="退款金额">
      <template #default="{ row }">
        {{ row.refund_price }}
      </template>
    </FTableColumn>
    <FTableColumn label="退款方式">
      <template #default="{ row }">
        {{ PAY_METHOD[row.refund_method] }}
      </template>
    </FTableColumn>
    <FTableColumn label="退款凭证">
      <template #default="{ row }">
        <FSpace v-if="row.refund_proof">
          <FImage
            v-for="url in row.refund_proof.split(',')" :key="url" class="proof-image" preview hide-on-click-modal
            :src="url"
          />
        </FSpace>
      </template>
    </FTableColumn>
    <FTableColumn :min-width="163" label="报名时间">
      <template #default="{ row }">
        {{ formatTimestamp(row.create_time * 1000) }}
      </template>
    </FTableColumn>
    <FTableColumn prop="remark" label="备注" />
    <FTableColumn :min-width="163" prop="update_time" label="更新时间">
      <template #default="{ row }">
        {{ formatTimestamp(row.update_time * 1000) }}
      </template>
    </FTableColumn>
    <FTableColumn :min-width="163" prop="update_time" label="取消报名时间">
      <template #default="{ row }">
        {{ formatTimestamp(row.status_time * 1000) }}
      </template>
    </FTableColumn>
    <FTableColumn label="操作" align="center" fixed="right" :action="action" :width="210" />
  </FTable>
  <FPagination
    v-if="!loadingOnce"
    show-total
    class="pagination"
    :page-size="pageState.page_size"
    :total-count="data?.page?.total"
    show-size-changer
    show-quick-jumper
    @change="handleChange"
  />
</template>

<style scoped>
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

.proof-image {
  width: 30px;
  height: 30px;
}

.pagination {
  margin-top: 10px;
  align-self: center;
}

.table {
  flex: 1;
}

.query-form {
  width: 500px;
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
</style>
