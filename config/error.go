// 全局错误码定义

package config

// 定义一些通用的错误码常量，并将其与相应的错误信息进行了映射。
// 设计错误码 一个良好结构的错误码有助于简化问题描述, 当前设计的错误码共有五位, 结构如下:
//  1 第一位是服务级别, 1 为系统错误, 2 为普通错误.
// 00 第二三位是模块, 模块不是指 Go 中的模块, 而是指代某个范围, 比如数据库错误, 认证错误.
// 01 第四五位是具体错误, 比如数据库错误中的插入错误, 找不到数据等.
// 定义错误码的时候不光有 Code 数字, 也会有对应的文本信息。
// 通常，文本分为两类：一类是给用户看的, 另一类是用于 debug 的。

// 错误码常量定义
const (
	ErrUnknown          = 10000 // 未知错误
	ErrBadRequest       = 10001 // 错误的请求
	ErrUnauthorized     = 10004 // 未授权访问
	ErrInternal         = 10007 // 内部服务器错误
	ErrNotFound         = 10010 // 请求的资源不存在
	ErrParamEmpty       = 10020 // 传入的参数为空或者不合法
	ErrFrequent         = 10030 // 请求过于频繁请稍后再试
	ErrDatabase         = 20000 // 数据库错误
	ErrNoRecords        = 20010 // 没有找到指定类型的记录
	ErrResExists        = 20020 // 资源已存在
	ErrEmptyIdle        = 20022 // 闲置空资源待认领
	ErrUserDisabled     = 20040 // 帐号已被禁用
	ErrNoPermission     = 20050 // 帐号无权限操作
	ErrNoActivate       = 20060 // 帐号尚未激活
	ErrValidation       = 20100 // 用户名或密码错误
	ErrHeader           = 20104 // 授权标头为空
	ErrToken            = 20200 // Token错误或过期
	ErrVerificationCode = 20301 // 验证码错误
	ErrExpireCode       = 20302 // 验证码已过期
	ErrFormat           = 20400 // 格式不正确
	ErrNotVip           = 20500 // 您不是VIP或已过期

)

// 错误码对应的错误信息
var ErrMsgs = map[int]string{
	ErrUnknown:          "未知错误",
	ErrBadRequest:       "错误的请求",
	ErrUnauthorized:     "未授权访问",
	ErrInternal:         "内部服务器错误",
	ErrNotFound:         "请求的资源不存在",
	ErrParamEmpty:       "传入的参数为空或者不合法",
	ErrFrequent:         "请求过于频繁请稍后再试",
	ErrDatabase:         "数据库错误",
	ErrNoRecords:        "没有找到指定类型的记录",
	ErrResExists:        "资源已存在",
	ErrEmptyIdle:        "闲置空资源待认领",
	ErrUserDisabled:     "帐号已被禁用",
	ErrNoPermission:     "帐号无权限操作",
	ErrNoActivate:       "帐号尚未激活",
	ErrValidation:       "用户名或密码错误",
	ErrHeader:           "授权标头为空",
	ErrToken:            "Token错误或过期",
	ErrVerificationCode: "验证码错误",
	ErrExpireCode:       "验证码已过期",
	ErrFormat:           "格式不正确",
	ErrNotVip:           "您不是VIP或已过期",
}
