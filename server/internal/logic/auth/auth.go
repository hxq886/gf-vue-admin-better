// Package auth 认证业务实现: 登录校验 / 令牌生成 / 登出拉黑。
package auth

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/crypto/bcrypt"

	v1 "gf-vue-admin-better/api/auth/v1"
	"gf-vue-admin-better/internal/dao"
	"gf-vue-admin-better/internal/middleware"
	"gf-vue-admin-better/internal/service"
)

// sAuth 认证业务实现
type sAuth struct{}

func init() {
	service.RegisterAuth(&sAuth{})
}

// Login 账号密码登录: 校验通过后签发 JWT 令牌。
func (s *sAuth) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	// 1. 查询用户
	record, err := dao.User.Ctx(ctx).
		Where(dao.User.Columns().Username, req.Username).
		One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询用户失败")
	}
	if record.IsEmpty() {
		return nil, gerror.New("账号或密码不正确")
	}

	// 2. 状态校验
	if record[dao.User.Columns().Status].Int() != 1 {
		return nil, gerror.New("账号已被停用, 请联系管理员")
	}

	// 3. 密码校验 (bcrypt)
	if err = bcrypt.CompareHashAndPassword(
		[]byte(record[dao.User.Columns().Password].String()),
		[]byte(req.Password),
	); err != nil {
		return nil, gerror.New("账号或密码不正确")
	}

	// 4. 签发令牌 (载荷携带用户ID, 键名 uid 与中间件 IdentityKey 一致, 避免污染请求参数 id 字段)
	token, _, err := middleware.Auth().TokenGenerator(g.Map{
		"uid": record[dao.User.Columns().Id].Uint64(),
	})
	if err != nil {
		return nil, gerror.Wrap(err, "令牌签发失败")
	}

	return &v1.LoginRes{AccessToken: token}, nil
}

// Logout 登出: 将当前令牌加入黑名单, 实现服务端"销毁Token"。
func (s *sAuth) Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	middleware.Auth().LogoutHandler(ctx)
	return &v1.LogoutRes{}, nil
}
