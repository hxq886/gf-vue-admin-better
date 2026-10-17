import request from '@/utils/request'

// 获取动态路由: GET /api/menu/getRouters
// 返回树形路由 (path/component/name/meta...), 字段结构严格适配 convertRouter
export function getRouterList() {
  return request({
    url: '/menu/getRouters',
    method: 'get',
  })
}
