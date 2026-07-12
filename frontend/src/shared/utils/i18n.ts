/**
 * 国际化消息定义
 */

import { reactive } from 'vue'

export type MessageKey = keyof typeof builtinMessages.zh

/**
 * Built-in translations. Custom languages are registered at runtime via
 * a user-provided i18n.json file (see registerCustomMessages).
 */
export const builtinMessages = {
  zh: {
    language: '语言',
    search: '搜索',
    noResults: '未找到匹配的接口',
    loading: '加载中...',
    loadError: '加载 openapi.json 失败',
    requestParams: '请求参数',
    response: '返回响应',
    debug: '调试',
    paramName: '参数名',
    paramType: '类型',
    paramRequired: '必填',
    paramDesc: '说明',
    yes: '是',
    no: '否',
    bodyParam: 'Body 参数',
    example: '示例',
    required: '必填',
    noResponseBody: '无响应体',
    prevPage: '上一页',
    nextPage: '下一页',
    clone: '克隆',
    export: '导出',
    onlineDebug: '在线运行',
    send: '发送',
    sending: '···',
    query: 'Query',
    queryParams: 'Query 参数',
    result: '返回结果',
    sendHint: '点击"发送"按钮获取返回结果',
    example1: '示例 1',
    addRow: '+ 添加',
    value: '值',
    pathParam: 'path 参数',
    queryParam: 'query 参数',
    headerParam: 'header 参数',
    cookieParam: 'cookie 参数',
    selectApi: '从左侧选择一个接口查看详情',
    pathParams: 'Path 参数',
    addParam: '添加参数',
    paramValue: '参数值',
    headers: 'Headers',
    cookie: 'Cookie',
    body: 'Body',
    bodyNone: '该请求没有请求体',
    chooseFile: '选择文件',
    header: 'Header',
    clickToCopy: '点击复制',
    copied: '已复制!',
    globalSettings: '全局设置',
    settings: '全局鉴权',
    settingsPageDesc: '配置后将自动注入到所有调试请求中，无需每次手动填写。',
    settingsType: '鉴权类型',
    settingsNone: '无',
    settingsNoneHint: '未配置全局鉴权，每次请求将不携带认证信息。',
    settingsBearerHint: '将自动在请求头中添加 Authorization: Bearer <token>',
    settingsKeyName: 'Header / Query 名称',
    settingsKeyValue: 'API Key 值',
    settingsKeyIn: '传递方式',
    settingsUsername: '用户名',
    settingsPassword: '密码',
    settingsClear: '清除鉴权信息',
    history: '历史',
    historyClear: '清空',
    historyEmpty: '暂无请求历史',
    selectFile: '选择文件',
    formDataParam: 'form-data 参数',
  },
  en: {
    language: 'Language',
    search: 'Search',
    noResults: 'No matching APIs',
    loading: 'Loading...',
    loadError: 'Failed to load openapi.json',
    requestParams: 'Request Parameters',
    response: 'Response',
    debug: 'Debug',
    paramName: 'Name',
    paramType: 'Type',
    paramRequired: 'Required',
    paramDesc: 'Description',
    yes: 'Yes',
    no: 'No',
    bodyParam: 'Body',
    example: 'Example',
    required: 'Required',
    noResponseBody: 'No response body',
    prevPage: 'Previous',
    nextPage: 'Next',
    clone: 'Clone',
    export: 'Export',
    onlineDebug: 'Debug',
    send: 'Send',
    sending: '···',
    query: 'Query',
    queryParams: 'Query Params',
    result: 'Response',
    sendHint: "Click 'Send' to see the response",
    example1: 'Example 1',
    addRow: '+ Add',
    value: 'Value',
    pathParam: 'path params',
    queryParam: 'query params',
    headerParam: 'header params',
    cookieParam: 'cookie params',
    selectApi: 'Select an API from the sidebar',
    pathParams: 'Path Params',
    addParam: 'Add Param',
    chooseFile: 'Choose file',
    paramValue: 'Value',
    headers: 'Headers',
    cookie: 'Cookie',
    body: 'Body',
    bodyNone: 'No request body',
    header: 'Header',
    clickToCopy: 'Click to copy',
    copied: 'Copied!',
    globalSettings: 'Settings',
    settings: 'Global Auth',
    settingsPageDesc: 'Credentials configured here will be automatically injected into every debug request.',
    settingsType: 'Auth Type',
    settingsNone: 'None',
    settingsNoneHint: 'No auth configured. Requests will be sent without credentials.',
    settingsBearerHint: 'Adds Authorization: Bearer <token> to every request.',
    settingsKeyName: 'Header / Query Name',
    settingsKeyValue: 'API Key Value',
    settingsKeyIn: 'Add to',
    settingsUsername: 'Username',
    settingsPassword: 'Password',
    settingsClear: 'Clear Auth',
    history: 'History',
    historyClear: 'Clear',
    historyEmpty: 'No request history yet',
    selectFile: 'Select file',
    formDataParam: 'form-data params',
  },
}

/**
 * The display name shown in the language switcher for the custom language.
 * Populated from the "name" field of the user-provided i18n.json.
 */
export const customLanguageName = reactive({ value: '' })

/**
 * Runtime translation registry: built-in languages plus an optional
 * custom language registered from i18n.json.
 */
export const messages: Record<string, Record<string, string>> = reactive({
  zh: builtinMessages.zh,
  en: builtinMessages.en,
})

/**
 * Register a custom language from a parsed i18n.json payload.
 * Shape: { name: string, messages: Record<string, string> }.
 * Returns true when a usable message table was registered.
 */
export function registerCustomMessages(payload: unknown): boolean {
  if (!payload || typeof payload !== 'object') return false

  const { name, messages: table } = payload as { name?: unknown; messages?: unknown }
  if (!table || typeof table !== 'object') return false

  messages.custom = table as Record<string, string>
  customLanguageName.value = typeof name === 'string' && name.trim() ? name.trim() : 'Custom'
  return true
}
