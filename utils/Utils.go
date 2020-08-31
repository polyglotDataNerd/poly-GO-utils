package utils

import (
	"os"
	"reflect"
)

type Props interface {
	GetEnv() string
	SetEnv()
	GetUser() string
	SetUser()
	GetPW() string
	SetPW()
}

type Mutator struct {
	SetterKeyEnv    string
	SetterValueEnv  string
	SetterKeyUser   string
	SetterValueUser string
	SetterKeyPW     string
	SetterValuePW   string
}

func (m *Mutator) SetEnv() {
	os.Setenv(m.SetterKeyEnv, m.SetterValueEnv)
}

func (m *Mutator) GetEnv() string {
	return os.Getenv(m.SetterKeyEnv)
}

func (m *Mutator) SetUser() {
	os.Setenv(m.SetterKeyUser, m.SetterValueUser)
}

func (m *Mutator) GetUser() string {
	return os.Getenv(m.SetterKeyUser)
}

func (m *Mutator) SetPW() {
	os.Setenv(m.SetterKeyPW, m.SetterValuePW)
}

func (m *Mutator) GetPW() string {
	return os.Getenv(m.SetterKeyPW)
}

func isNilMap(inMap map[string]interface{}) (response string) {
	val := reflect.ValueOf(inMap)

	if val.Kind() == reflect.Map {
		for _, k := range val.MapKeys() {
			v := val.MapIndex(k)
			if reflect.ValueOf(v).IsNil() {
				Error.Println("interface{} is null")
				response = "no value"
			}
		}
	}
	return response
}
