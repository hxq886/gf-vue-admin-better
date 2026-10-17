/**
 * @description 导出默认网路配置 (已适配 GoFrame 后端统一响应结构 {code,msg,data})
 **/
const network = {
  // 接口地址: 开发环境走 rspack devServer 代理到 GoFrame 后端(127.0.0.1:8000);
  // 生产环境为相对路径, 适配前端 dist 与后端单二进制同源部署
  baseURL: process.env.NODE_ENV === "production" ? "./api" : "/api",
  //配后端数据的接收方式application/json;charset=UTF-8或者application/x-www-form-urlencoded;charset=UTF-8
  contentType: "application/json;charset=UTF-8",
  //消息框消失时间
  messageDuration: 3000,
  //最长请求时间
  requestTimeout: 15000,
  //操作正常code，支持String、Array、int多种类型 (后端统一响应 code=200 为成功)
  successCode: [200],
  //登录失效code (后端: 401 Token失效/未登录, 自动登出并跳转登录页)
  invalidCode: 401,
  //无权限code (后端: 403 权限不足)
  noPermissionCode: 403,
};
module.exports = network;
