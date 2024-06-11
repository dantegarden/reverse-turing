package dal

import (
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"reverse-turing/app/internal/types"
	"reverse-turing/common/errno"
	"reverse-turing/common/model"
	"time"
)

type Agent struct {
	model.Model
	ModelName    string                  `json:"modelName" gorm:"column:model_name;size:32;"`
	Endpoint     string                  `json:"endpoint" gorm:"column:endpoint;size:300;"`
	ApiKey       string                  `json:"apiKey" gorm:"column:api_key;size:32;"`
	PromptParams model.MapStringToString `json:"promptParams" gorm:"column:prompt_params;"`
	Enable       bool                    `json:"enable" gorm:"column:enable;"`
}

func (c *Agent) Convert() interface{} {
	var pbResp types.Agent
	_ = copier.CopyWithOption(&pbResp, c, copier.Option{Converters: []copier.TypeConverter{
		{
			SrcType: time.Time{},
			DstType: copier.String,
			Fn: func(src interface{}) (interface{}, error) {
				s, ok := src.(time.Time)

				if !ok {
					logx.Errorf("src type not matching")
					return nil, errno.ServerErr
				}

				return s.Format(time.DateTime), nil
			},
		}},
		DeepCopy: true,
	})

	return &pbResp
}
