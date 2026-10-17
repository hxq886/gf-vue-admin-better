# 接口文档

所有接口统一前缀 `/api`，统一响应结构：

```json
{ "code": 200, "msg": "success", "data": {} }
```

| code | 含义 | 前端行为 |
| --- | --- | --- |
| 200 | 成功 | 正常返回 data |
| 401 | 未登录 / Token 失效 | 自动登出并跳转登录页 |
| 403 | 权限不足 | 跳转 403 页面 |
| 500 | 服务异常 | 提示错误信息 |

> HTTP 状态码恒为 200（含 404 的未知接口也返回 JSON 结构），业务状态由 `code` 表达。

## 鉴权

- 登录成功后签发 JWT（有效期默认 24h，见 `server/manifest/config/config.yaml` 的 `jwt.timeout`）
- 后续请求经请求头 `accessToken` 携带（裸 token，无需 `Bearer` 前缀，后端自动适配）
- 登出后令牌进入黑名单立即失效

```
accessToken: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 基础接口（公开，无需令牌）

### 1. 登录

`POST /api/login`

请求：

```json
{ "username": "admin", "password": "123456" }
```

响应：

```json
{
  "code": 200,
  "msg": "success",
  "data": { "accessToken": "eyJhbGciOiJIUzI1NiJ9..." }
}
```

错误：账号不存在 / 密码错误 / 账号停用时返回 `code:500` 与对应提示。

### 2. 登出

`POST /api/logout`（需令牌）

响应：`{ "code": 200, "msg": "success", "data": {} }`，旧令牌立即失效。

---

## 用户接口（需令牌）

### 3. 获取登录用户信息

`GET /api/user/getInfo`

响应（权限标识数组供前端按钮级权限精确匹配）：

```json
{
  "code": 200,
  "data": {
    "username": "admin",
    "nickname": "管理员",
    "avatar": "",
    "roles": ["admin"],
    "permissions": ["system:user:create", "system:user:update", "..."]
  }
}
```

> `admin` 角色返回全部按钮权限；普通用户按"用户-角色-菜单"关联返回。

### 4. 用户列表（分页）

`GET /api/user/list?page=1&pageSize=10&username=&status=`

响应：`data: { list: [...用户(含 roleIds)], total: 100 }`

### 5. 新增用户

`POST /api/user/create`

```json
{ "username": "test", "nickname": "测试", "password": "123456", "roleIds": [2], "email": "", "mobile": "", "status": 1, "remark": "" }
```

### 6. 修改用户

`PUT /api/user/update`（body 含 `id`，账号与密码不可改）

### 7. 重置用户密码

`PUT /api/user/password`

```json
{ "id": 2, "password": "new123456" }
```

### 8. 批量删除用户

`DELETE /api/user/delete`

```json
{ "ids": [2, 3] }
```

> 内置超管（id=1）受保护，接口拒绝删除。

---

## 角色接口（需令牌）

| 接口 | 方法 | 说明 |
| --- | --- | --- |
| `/api/role/list` | GET | 全量角色列表（含 `menuIds` 已授权菜单） |
| `/api/role/create` | POST | 新增角色 `{ name, roleKey, sort, status, menuIds, remark }` |
| `/api/role/update` | PUT | 修改角色（body 含 `id`） |
| `/api/role/delete` | DELETE | 批量删除 `{ ids: [...] }` |

> 内置超管角色（id=1）不可删除、不可改权限字符与状态。

---

## 菜单接口（需令牌）

### 9. 动态路由（核心）

`GET /api/menu/getRouters`

返回当前登录用户可见的树形路由，字段结构严格适配 vue3-admin-better：

```json
{
  "code": 200,
  "data": [
    {
      "path": "/",
      "component": "Layout",
      "redirect": "/index",
      "name": "Root",
      "alwaysShow": false,
      "hidden": false,
      "children": [
        {
          "path": "/index",
          "component": "vab/index/index",
          "name": "Index",
          "hidden": false,
          "meta": { "title": "首页", "icon": "home", "affix": true, "noKeepAlive": false }
        }
      ]
    }
  ]
}
```

- `component`：视图相对路径（自动补 `views/` 前缀与 `.vue` 后缀）；`Layout` / `EmptyLayout` 为布局组件
- `meta.noKeepAlive`：`true` 表示不缓存该页面
- 超管返回全部启用菜单；普通用户按角色关联过滤（按钮类型不参与路由）

### 10. 菜单管理 CRUD

| 接口 | 方法 | 说明 |
| --- | --- | --- |
| `/api/menu/list` | GET | 全量菜单平铺列表（前端按 `parentId` 建树） |
| `/api/menu/create` | POST | 新增（`menuType`: M目录/C菜单/F按钮） |
| `/api/menu/update` | PUT | 修改（上级不可选自身） |
| `/api/menu/delete` | DELETE | 批量删除（存在子菜单时拒绝） |

---

## 字典接口（需令牌）

| 接口 | 方法 | 说明 |
| --- | --- | --- |
| `/api/dict/type/list` | GET | 字典类型列表 |
| `/api/dict/type/create` | POST | 新增类型（`dictType` 类型键唯一） |
| `/api/dict/type/update` | PUT | 修改类型（类型键变更自动同步字典数据） |
| `/api/dict/type/delete` | DELETE | 批量删除类型（级联删除其下数据） |
| `/api/dict/data/list?dictType=键` | GET | 按类型键查字典数据 |
| `/api/dict/data/create` | POST | 新增数据 `{ dictType, label, value, tagClass, sort, status }` |
| `/api/dict/data/update` | PUT | 修改数据 |
| `/api/dict/data/delete` | DELETE | 批量删除数据 |

---

## 通用错误示例

```json
// 未登录 / Token 失效
{ "code": 401, "msg": "登录状态已过期, 请重新登录", "data": null }

// 权限不足（当前版本未启用接口级权限拦截, 预留）
{ "code": 403, "msg": "权限不足", "data": null }

// 业务错误（如账号密码错误）
{ "code": 500, "msg": "账号或密码错误", "data": {} }

// 未知接口（/api 前缀下未匹配路由）
{ "code": 404, "msg": "接口不存在", "data": null }
```
