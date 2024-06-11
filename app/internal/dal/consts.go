package dal

import (
	"fmt"
	"reverse-turing/common/errno"
	"strings"
)

type TalkType string

const (
	// TalkTypeOpening 开场白
	TalkTypeOpening TalkType = "opening"
	// TalkTypeAsk 要求询问
	TalkTypeAsk TalkType = "ask"
	// TalkTypeAnswer 回答
	TalkTypeAnswer TalkType = "answer"
	// TalkTypeVote 投票
	TalkTypeVote TalkType = "vote"
)

var TalkQueryTplMap = map[TalkType]string{
	TalkTypeOpening: "请你说一段开场白",
	TalkTypeAsk:     "请你对自己怀疑的人提出一个问题",
	TalkTypeAnswer:  "%s向你提问：%s",
	TalkTypeVote:    "根据以下这些对话，请你选出最有可能是人类的对象，并给出理由。\n%s",
}

func GetQuery(talkType TalkType, args ...string) (string, error) {
	tpl, exists := TalkQueryTplMap[talkType]
	if !exists {
		return "", errno.UnsupportedTalkTypeErr
	}

	if len(args) > 0 && len(args) != strings.Count(tpl, "%s") {
		return "", errno.InvalidTalkTypeArgsErr
	}

	if len(args) > 0 {
		// 初始化一个空的接口切片，长度与 strings 切片相同
		var anys []any = make([]any, len(args))
		// 遍历 strings 切片并将每个元素添加到 anys 切片中
		for i, v := range args {
			anys[i] = v
		}
		return fmt.Sprintf(tpl, anys...), nil
	} else {
		return tpl, nil
	}
}
