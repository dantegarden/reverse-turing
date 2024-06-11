package errno

import "github.com/zeromicro/x/errors"

const (
	OkMsg = "Ok"
)

func NewErrNo(code int, msg string) error {
	return errors.New(code, msg)
}

const (
	OKCode                   = 0
	BindAndValidateErrorCode = 41400
	ForbiddenErrorCode       = 41403
	ServerErrorCode          = 51500
	UnknownErrorCode         = 51999
)

const (
	_ = iota + 41000

	CharacterNotFoundCode
	AgentNotFoundCode
	GameNotFoundCode
	GameCharacterNotFoundCode
	CharacterDisable
	AgentDisable
	DeviceIdNotFoundCode
	InvalidAiNumCode
	NoAvailableCharacterCode
	NoAvailableAgentCode
	UnsupportedTalkTypeCode
	InvalidTalkTypeArgsCode
)

var (
	Success            = NewErrNo(OKCode, OkMsg)
	ForbiddenError     = NewErrNo(ForbiddenErrorCode, "非法操作")
	BindAndValidateErr = NewErrNo(BindAndValidateErrorCode, "参数绑定/校验错误")
	ServerErr          = NewErrNo(ServerErrorCode, "服务端错误")

	CharacterNotFoundErr     = NewErrNo(CharacterDisable, "角色不存在")
	GameCharacterNotFoundErr = NewErrNo(GameCharacterNotFoundCode, "GameCharacter不存在")
	AgentNotFoundErr         = NewErrNo(AgentDisable, "Agent不存在")
	GameNotFoundErr          = NewErrNo(GameNotFoundCode, "Game不存在")
	DeviceIdNotFoundErr      = NewErrNo(DeviceIdNotFoundCode, "DeviceId不存在")
	InvalidAiNumErr          = NewErrNo(InvalidAiNumCode, "aiNum必须是大于2的偶数")
	NoAvailableCharacterErr  = NewErrNo(NoAvailableCharacterCode, "没有可用的角色")
	NoAvailableAgentErr      = NewErrNo(NoAvailableAgentCode, "没有可用的Agent")
	UnsupportedTalkTypeErr   = NewErrNo(UnsupportedTalkTypeCode, "不支持的对话类型")
	InvalidTalkTypeArgsErr   = NewErrNo(InvalidTalkTypeArgsCode, "对话模板参数值错误或数量不对")
)
