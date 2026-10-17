-- ----------------------------------------------------------------------------
-- gf-vue-admin-better 数据库初始化脚本
-- 用法: mysql -h127.0.0.1 -P3306 -uroot -p < init.sql
-- 说明: 包含 RBAC 权限体系(用户/角色/菜单/关联表)与字典表结构及种子数据
-- ----------------------------------------------------------------------------

CREATE DATABASE IF NOT EXISTS `gf_vue_admin_better` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE `gf_vue_admin_better`;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------------------------------------------------------
-- 1. 系统用户表
-- ----------------------------------------------------------------------------
DROP TABLE IF EXISTS `sys_user`;
CREATE TABLE `sys_user` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username`   VARCHAR(64)  NOT NULL                COMMENT '登录账号',
  `nickname`   VARCHAR(64)  NOT NULL DEFAULT ''     COMMENT '用户昵称',
  `password`   VARCHAR(100) NOT NULL                COMMENT '密码(bcrypt加密存储)',
  `avatar`     VARCHAR(255) NOT NULL DEFAULT ''     COMMENT '头像地址',
  `email`      VARCHAR(128) NOT NULL DEFAULT ''     COMMENT '邮箱',
  `mobile`     VARCHAR(20)  NOT NULL DEFAULT ''     COMMENT '手机号',
  `status`     TINYINT      NOT NULL DEFAULT 1      COMMENT '状态(1启用 0停用)',
  `remark`     VARCHAR(255) NOT NULL DEFAULT ''     COMMENT '备注',
  `created_at` DATETIME     NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME     NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '系统用户表';

-- ----------------------------------------------------------------------------
-- 2. 系统角色表
-- ----------------------------------------------------------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name`       VARCHAR(64)  NOT NULL                COMMENT '角色名称',
  `role_key`   VARCHAR(64)  NOT NULL                COMMENT '角色权限字符(接口鉴权用)',
  `sort`       INT          NOT NULL DEFAULT 0      COMMENT '显示顺序',
  `status`     TINYINT      NOT NULL DEFAULT 1      COMMENT '状态(1启用 0停用)',
  `remark`     VARCHAR(255) NOT NULL DEFAULT ''     COMMENT '备注',
  `created_at` DATETIME     NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME     NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '系统角色表';

-- ----------------------------------------------------------------------------
-- 3. 系统菜单表 (目录M / 菜单C / 按钮F, 同时驱动动态路由与按钮权限)
-- ----------------------------------------------------------------------------
DROP TABLE IF EXISTS `sys_menu`;
CREATE TABLE `sys_menu` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  `parent_id`   BIGINT UNSIGNED NOT NULL DEFAULT 0     COMMENT '父菜单ID(顶级为0)',
  `menu_type`   CHAR(1)       NOT NULL DEFAULT 'C'     COMMENT '类型(M目录 C菜单 F按钮)',
  `title`       VARCHAR(64)   NOT NULL                COMMENT '菜单标题',
  `name`        VARCHAR(64)   NOT NULL DEFAULT ''     COMMENT '路由name(唯一)',
  `path`        VARCHAR(191)  NOT NULL DEFAULT ''     COMMENT '路由path',
  `component`   VARCHAR(191)  NOT NULL DEFAULT ''     COMMENT '组件(Layout/EmptyLayout/视图相对路径)',
  `redirect`    VARCHAR(191)  NOT NULL DEFAULT ''     COMMENT '重定向地址(noRedirect/子路由path)',
  `always_show` TINYINT(1)    NOT NULL DEFAULT 0      COMMENT '始终显示根路由',
  `hidden`      TINYINT(1)    NOT NULL DEFAULT 0      COMMENT '隐藏菜单(1隐藏)',
  `icon`        VARCHAR(64)   NOT NULL DEFAULT ''     COMMENT '菜单图标(FontAwesome风格名称)',
  `keep_alive`  TINYINT(1)    NOT NULL DEFAULT 0      COMMENT '是否缓存(1缓存)',
  `affix`       TINYINT(1)    NOT NULL DEFAULT 0      COMMENT '是否固定在标签栏',
  `permission`  VARCHAR(191)  NOT NULL DEFAULT ''     COMMENT '权限标识(按钮: system:user:create)',
  `sort`        INT           NOT NULL DEFAULT 0      COMMENT '显示顺序',
  `status`      TINYINT       NOT NULL DEFAULT 1      COMMENT '状态(1启用 0停用)',
  `remark`      VARCHAR(255)  NOT NULL DEFAULT ''     COMMENT '备注',
  `created_at`  DATETIME      NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`  DATETIME      NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '系统菜单表';

-- ----------------------------------------------------------------------------
-- 4. 用户角色关联表
-- ----------------------------------------------------------------------------
DROP TABLE IF EXISTS `sys_user_role`;
CREATE TABLE `sys_user_role` (
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `role_id` BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
  PRIMARY KEY (`user_id`, `role_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户角色关联表';

-- ----------------------------------------------------------------------------
-- 5. 角色菜单关联表
-- ----------------------------------------------------------------------------
DROP TABLE IF EXISTS `sys_role_menu`;
CREATE TABLE `sys_role_menu` (
  `role_id` BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
  `menu_id` BIGINT UNSIGNED NOT NULL COMMENT '菜单ID',
  PRIMARY KEY (`role_id`, `menu_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '角色菜单关联表';

-- ----------------------------------------------------------------------------
-- 6. 字典类型表
-- ----------------------------------------------------------------------------
DROP TABLE IF EXISTS `sys_dict_type`;
CREATE TABLE `sys_dict_type` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '字典类型ID',
  `name`       VARCHAR(64)   NOT NULL                COMMENT '字典名称',
  `dict_type`  VARCHAR(64)   NOT NULL                COMMENT '类型键(唯一)',
  `status`     TINYINT       NOT NULL DEFAULT 1      COMMENT '状态(1启用 0停用)',
  `remark`     VARCHAR(255)  NOT NULL DEFAULT ''     COMMENT '备注',
  `created_at` DATETIME      NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME      NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_dict_type` (`dict_type`)
) ENGINE = InnoDB AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '字典类型表';

-- ----------------------------------------------------------------------------
-- 7. 字典数据表
-- ----------------------------------------------------------------------------
DROP TABLE IF EXISTS `sys_dict_data`;
CREATE TABLE `sys_dict_data` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '字典数据ID',
  `dict_type`  VARCHAR(64)   NOT NULL                COMMENT '所属类型键',
  `label`      VARCHAR(64)   NOT NULL                COMMENT '标签(展示用)',
  `value`      VARCHAR(64)   NOT NULL                COMMENT '键值(存储用)',
  `tag_class`  VARCHAR(32)   NOT NULL DEFAULT ''     COMMENT '标签样式(primary/success/warning/danger)',
  `sort`       INT           NOT NULL DEFAULT 0      COMMENT '显示顺序',
  `status`     TINYINT       NOT NULL DEFAULT 1      COMMENT '状态(1启用 0停用)',
  `remark`     VARCHAR(255)  NOT NULL DEFAULT ''     COMMENT '备注',
  `created_at` DATETIME      NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME      NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_dict_type` (`dict_type`)
) ENGINE = InnoDB AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '字典数据表';

-- ----------------------------------------------------------------------------
-- 初始数据: 超级管理员 (账号 admin / 密码 123456, bcrypt 加密)
-- ----------------------------------------------------------------------------
INSERT INTO `sys_user` (`id`, `username`, `nickname`, `password`, `avatar`, `status`, `remark`) VALUES
(1, 'admin', '超级管理员', '$2a$10$18rdbdnvTkU5.a01eNtFbeT8O7V.4Y8NHPY9VMN5MdE/c8Ya0g.bO',
 'https://gcore.jsdelivr.net/gh/zxwk1998/image/avatar/avatar_1.png', 1, '内置超级管理员,拥有全部权限');

-- ----------------------------------------------------------------------------
-- 初始数据: 角色
-- ----------------------------------------------------------------------------
INSERT INTO `sys_role` (`id`, `name`, `role_key`, `sort`, `status`, `remark`) VALUES
(1, '超级管理员', 'admin', 1, 1, '拥有全部权限,不受角色菜单限制'),
(2, '普通角色',   'common', 2, 1, '示例角色,可按需分配菜单');

-- ----------------------------------------------------------------------------
-- 初始数据: 菜单 (与 vue3-admin-better 动态路由字段严格对齐)
-- component 说明: Layout=主框架 / EmptyLayout=空布局 / 其余为 src/views 下相对路径(不带 views 前缀)
-- ----------------------------------------------------------------------------
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_type`, `title`, `name`, `path`, `component`, `redirect`, `always_show`, `hidden`, `icon`, `keep_alive`, `affix`, `permission`, `sort`, `status`, `remark`) VALUES
-- 首页 (固定在工作台)
(1,  0, 'M', '首页',   '',          '/',      'Layout', 'index',     0, 0, 'home',    0, 0, '', 0, 1, '根路由'),
(2,  1, 'C', '首页',   'Index',     'index',  'index/index', '',      0, 0, 'home',    0, 1, '', 1, 1, '工作台'),
-- 系统管理目录
(3,  0, 'M', '系统管理', 'System',  '/system', 'Layout', 'noRedirect', 1, 0, 'cogs',    0, 0, '', 2, 1, '系统管理目录'),
(4,  3, 'C', '用户管理', 'SystemUser',  'user',  'vab/system/user/index', '', 0, 0, 'user',     0, 0, '', 1, 1, '用户管理页面'),
(5,  3, 'C', '角色管理', 'SystemRole',  'role',  'vab/system/role/index', '', 0, 0, 'user-cog', 0, 0, '', 2, 1, '角色管理页面'),
(6,  3, 'C', '菜单管理', 'SystemMenu',  'menu',  'vab/system/menu/index', '', 0, 0, 'list-ul',  0, 0, '', 3, 1, '菜单管理页面'),
(7,  3, 'C', '字典管理', 'SystemDict',  'dict',  'vab/system/dict/index', '', 0, 0, 'database', 0, 0, '', 4, 1, '字典管理页面'),
-- 用户管理按钮
(8,  4, 'F', '新增用户', '', '', '', '', 0, 0, '', 0, 0, 'system:user:create',   1, 1, ''),
(9,  4, 'F', '修改用户', '', '', '', '', 0, 0, '', 0, 0, 'system:user:update',   2, 1, ''),
(10, 4, 'F', '删除用户', '', '', '', '', 0, 0, '', 0, 0, 'system:user:delete',   3, 1, ''),
(11, 4, 'F', '重置密码', '', '', '', '', 0, 0, '', 0, 0, 'system:user:password', 4, 1, ''),
-- 角色管理按钮
(12, 5, 'F', '新增角色', '', '', '', '', 0, 0, '', 0, 0, 'system:role:create', 1, 1, ''),
(13, 5, 'F', '修改角色', '', '', '', '', 0, 0, '', 0, 0, 'system:role:update', 2, 1, ''),
(14, 5, 'F', '删除角色', '', '', '', '', 0, 0, '', 0, 0, 'system:role:delete', 3, 1, ''),
-- 菜单管理按钮
(15, 6, 'F', '新增菜单', '', '', '', '', 0, 0, '', 0, 0, 'system:menu:create', 1, 1, ''),
(16, 6, 'F', '修改菜单', '', '', '', '', 0, 0, '', 0, 0, 'system:menu:update', 2, 1, ''),
(17, 6, 'F', '删除菜单', '', '', '', '', 0, 0, '', 0, 0, 'system:menu:delete', 3, 1, ''),
-- 字典管理按钮
(18, 7, 'F', '新增字典', '', '', '', '', 0, 0, '', 0, 0, 'system:dict:create', 1, 1, ''),
(19, 7, 'F', '修改字典', '', '', '', '', 0, 0, '', 0, 0, 'system:dict:update', 2, 1, ''),
(20, 7, 'F', '删除字典', '', '', '', '', 0, 0, '', 0, 0, 'system:dict:delete', 3, 1, '');

-- ----------------------------------------------------------------------------
-- 初始数据: 管理员关联超管角色
-- ----------------------------------------------------------------------------
INSERT INTO `sys_user_role` (`user_id`, `role_id`) VALUES (1, 1);

-- 超管角色授权全部菜单 (逻辑层对 admin 角色放行全部权限,此处数据保证完整)
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`)
SELECT 1, `id` FROM `sys_menu`;

-- 普通角色仅授权工作台
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES (2, 1), (2, 2);

-- ----------------------------------------------------------------------------
-- 初始数据: 字典
-- ----------------------------------------------------------------------------
INSERT INTO `sys_dict_type` (`id`, `name`, `dict_type`, `status`, `remark`) VALUES
(1, '用户状态', 'sys_user_status', 1, '用户/角色/字典通用启停状态'),
(2, '菜单类型', 'sys_menu_type',   1, '菜单类型说明字典');

INSERT INTO `sys_dict_data` (`id`, `dict_type`, `label`, `value`, `tag_class`, `sort`, `status`, `remark`) VALUES
(1, 'sys_user_status', '启用', '1', 'success', 1, 1, '正常状态'),
(2, 'sys_user_status', '停用', '0', 'danger',  2, 1, '禁用状态'),
(3, 'sys_menu_type',   '目录', 'M', 'primary', 1, 1, ''),
(4, 'sys_menu_type',   '菜单', 'C', 'success', 2, 1, ''),
(5, 'sys_menu_type',   '按钮', 'F', 'warning', 3, 1, '');

SET FOREIGN_KEY_CHECKS = 1;
