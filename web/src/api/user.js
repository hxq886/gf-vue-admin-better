import request from "@/utils/request";

// 登录: POST /api/login (统一响应 data.accessToken)
export async function login(data) {
  return request({
    url: "/login",
    method: "post",
    data,
  });
}

// 获取登录用户信息: GET /api/user/getInfo (令牌经请求头自动携带)
// 返回 { username, nickname, avatar, roles, permissions }
export function getUserInfo() {
  return request({
    url: "/user/getInfo",
    method: "get",
  });
}

// 登出: POST /api/logout (服务端将令牌加入黑名单)
export function logout() {
  return request({
    url: "/logout",
    method: "post",
  });
}

// 注册 (保留原生注册页能力, 需后端自行实现 /api/register 后生效)
export function register() {
  return request({
    url: "/register",
    method: "post",
  });
}
