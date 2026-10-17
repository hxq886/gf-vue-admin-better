/**
 * @description 系统管理 API 封装 (用户/角色/菜单/字典), 对接 GoFrame 后端 /api 接口
 */
import request from "@/utils/request";

// ==================== 用户管理 ====================

// 用户列表 (分页)
export function getUserList(params) {
  return request({
    url: "/user/list",
    method: "get",
    params,
  });
}

// 新增用户
export function createUser(data) {
  return request({
    url: "/user/create",
    method: "post",
    data,
  });
}

// 修改用户 (账号不可改, 密码走重置接口)
export function updateUser(data) {
  return request({
    url: "/user/update",
    method: "put",
    data,
  });
}

// 重置用户密码
export function resetUserPassword(data) {
  return request({
    url: "/user/password",
    method: "put",
    data,
  });
}

// 批量删除用户
export function deleteUser(data) {
  return request({
    url: "/user/delete",
    method: "delete",
    data,
  });
}

// ==================== 角色管理 ====================

// 角色列表 (含已授权菜单ID)
export function getRoleList() {
  return request({
    url: "/role/list",
    method: "get",
  });
}

// 新增角色
export function createRole(data) {
  return request({
    url: "/role/create",
    method: "post",
    data,
  });
}

// 修改角色
export function updateRole(data) {
  return request({
    url: "/role/update",
    method: "put",
    data,
  });
}

// 批量删除角色
export function deleteRole(data) {
  return request({
    url: "/role/delete",
    method: "delete",
    data,
  });
}

// ==================== 菜单管理 ====================

// 菜单列表 (全部, 前端按 parentId 构建树)
export function getMenuList() {
  return request({
    url: "/menu/list",
    method: "get",
  });
}

// 新增菜单
export function createMenu(data) {
  return request({
    url: "/menu/create",
    method: "post",
    data,
  });
}

// 修改菜单
export function updateMenu(data) {
  return request({
    url: "/menu/update",
    method: "put",
    data,
  });
}

// 批量删除菜单
export function deleteMenu(data) {
  return request({
    url: "/menu/delete",
    method: "delete",
    data,
  });
}

// ==================== 字典管理 ====================

// 字典类型列表
export function getDictTypeList() {
  return request({
    url: "/dict/type/list",
    method: "get",
  });
}

// 新增字典类型
export function createDictType(data) {
  return request({
    url: "/dict/type/create",
    method: "post",
    data,
  });
}

// 修改字典类型 (类型键变更时同步字典数据)
export function updateDictType(data) {
  return request({
    url: "/dict/type/update",
    method: "put",
    data,
  });
}

// 批量删除字典类型 (级联删除其下字典数据)
export function deleteDictType(data) {
  return request({
    url: "/dict/type/delete",
    method: "delete",
    data,
  });
}

// 按类型键查询字典数据
export function getDictDataList(params) {
  return request({
    url: "/dict/data/list",
    method: "get",
    params,
  });
}

// 新增字典数据
export function createDictData(data) {
  return request({
    url: "/dict/data/create",
    method: "post",
    data,
  });
}

// 修改字典数据
export function updateDictData(data) {
  return request({
    url: "/dict/data/update",
    method: "put",
    data,
  });
}

// 批量删除字典数据
export function deleteDictData(data) {
  return request({
    url: "/dict/data/delete",
    method: "delete",
    data,
  });
}
