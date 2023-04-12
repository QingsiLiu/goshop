package common

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SwapToStruct 通过json tag 进行结构体赋值
func SwapToStruct(req, target interface{}) (err error) {
	dataByte, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = json.Unmarshal(dataByte, target)
	return
}

type H struct {
	Code                   string
	Message                string
	TraceId                string
	Data                   interface{}
	Rows                   interface{}
	Total                  interface{}
	SkyWalkingFynamicField string
}

func Resp(w http.ResponseWriter, code, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	h := H{
		Code:    code,
		Data:    data,
		Message: message,
	}

	ret, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
	}

	w.Write(ret)
}

func RespList(w http.ResponseWriter, code, message string, data interface{}, rows, total interface{}, skyWalkingFynamicField string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	h := H{
		Code:                   code,
		Data:                   data,
		Message:                message,
		Rows:                   rows,
		Total:                  total,
		SkyWalkingFynamicField: skyWalkingFynamicField,
	}

	ret, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
	}

	w.Write(ret)
}

/**
200 OKLoginSuccessVO
201 Created
401 Unauthorized
403 Forbidden
404 Not Found
**/

func RespOK(w http.ResponseWriter, message string, data interface{}) {
	Resp(w, "SUCCESS.", message, data)
}

func RespFail(w http.ResponseWriter, message string, data interface{}) {
	Resp(w, "TOKEN_FAIL", message, data)
}

func RespListOK(w http.ResponseWriter, message string, data, rows, total interface{}, skyWalkingFynamicField string) {
	RespList(w, "SUCCESS.", message, data, rows, total, skyWalkingFynamicField)
}

func RespListFail(w http.ResponseWriter, message string, data, rows, total interface{}, skyWalkingFynamicField string) {
	RespList(w, "TOKEN_FILE", message, data, rows, total, skyWalkingFynamicField)
}
