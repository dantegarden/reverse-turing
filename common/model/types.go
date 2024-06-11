package model

import (
	"database/sql/driver"
	"encoding/json"
	"k8s.io/apimachinery/pkg/runtime"
)

type ListString []string

func (list ListString) Value() (driver.Value, error) {
	return json.Marshal(list)
}

func (list *ListString) Scan(input interface{}) error {
	val := make([]string, 0)
	if err := json.Unmarshal(input.([]byte), &val); err != nil {
		return err
	}
	*list = val
	return nil
}

type MapStringToAny map[string]interface{}

func (mp MapStringToAny) Value() (driver.Value, error) {
	return json.Marshal(mp)
}

func (mp *MapStringToAny) Scan(input interface{}) error {
	val := make(map[string]interface{})
	if err := json.Unmarshal(input.([]byte), &val); err != nil {
		return err
	}
	*mp = val
	return nil
}

func (mp MapStringToAny) DeepCopy() MapStringToAny {
	if mp == nil {
		return nil
	}
	return runtime.DeepCopyJSON(mp)
}

type MapStringToString map[string]string

func (mp MapStringToString) Value() (driver.Value, error) {
	return json.Marshal(mp)
}

func (mp *MapStringToString) Scan(input interface{}) error {
	val := make(map[string]string)
	if err := json.Unmarshal(input.([]byte), &val); err != nil {
		return err
	}
	*mp = val
	return nil
}

func (mp MapStringToString) DeepCopy() MapStringToString {
	if mp == nil {
		return nil
	}
	out := MapStringToString{}
	for k, v := range mp {
		out[k] = v
	}
	return out
}

// Pair 定义一个通用的结构体来存储两个值
type Pair[T1 any, T2 any] struct {
	Left  T1
	Right T2
}

// Turple 定义一个通用的结构体来存储三个值
type Turple[T1 any, T2 any, T3 any] struct {
	Left   T1
	Right  T2
	Middle T3
}
