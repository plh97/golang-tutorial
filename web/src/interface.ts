import type { YES_NO } from './enums'

// 1: 国家 2: 省份 3: 城市
export enum LOCATION_TYPE {
  国家 = 1,
  省份 = 2,
  城市 = 3,
}
export interface ILocationInfo {
  id: number
  name: string
  type: LOCATION_TYPE
}
export interface ITag {
  id: number
  title: string
}
export interface IUserProfile {
  userId: number
  image: string[]
  id: number
  user_id: number
  nickname: string
  height: number
  body_type: number
  personality_type: number
  location_city_info: ILocationInfo | null
  gender: number
  find_gender: number
  school: string
  job_id: number
  job_name: string
  education_level: number
  birth: number
  income_type: number
  income_display: number
  email: string
  images: string[]
  intro: string
  life_intro: string
  hobby_intro: string
  life_style_intro: string
  find_location_city_ids: ILocationInfo[] | null
  find_location_city_info: ILocationInfo[] | null
  find_age_min: number
  find_age_max: number
  find_height_min: number
  find_height_max: number
  find_body_types: number[] | null
  find_education_levels: number[] | null
  find_other_intro: string
  experience_tags: ITag[] | null
  interest_tags: ITag[] | null
  activity_tags: ITag[] | null
  love_expectation_tags: ITag[] | null
  life_style_tags: ITag[] | null
  required_fields: string[]
  is_stable: YES_NO
  operation_remark: string
  tags: { id: number, title: string }[]

  home_location_country_info: ILocationInfo | null
  home_location_state_info: ILocationInfo | null
  home_location_city_info: ILocationInfo | null

  home_location_city_id: number
  home_location_country_id: number
  home_location_state_id: number
  location_city_id: number
}
