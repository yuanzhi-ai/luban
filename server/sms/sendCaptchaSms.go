// 短信服务的包
package sms

import (
	"fmt"
	"os"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/yuanzhi-ai/luban/server/log"
)

// 参考一下这个api https://next.api.aliyun.com/api-tools/demo/Dysmsapi/db7e1211-14e0-4b7b-9011-037dfb85d42e

func createClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

// 发送短信验证码
func SendMsg(code string, phoneNumber string) error {
	log.Debugf("into send msg")
	client, err := createClient(tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")), tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")))
	log.Debugf("key:%v secret:%v", os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"), os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"))
	if err != nil {
		return err
	}
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String("万象绘"),
		TemplateCode:  tea.String("SMS_462225349"),
		PhoneNumbers:  tea.String(phoneNumber),
		TemplateParam: tea.String(fmt.Sprintf("{\"code\":\"%v\"}", code)),
	}
	log.Debugf("ok request")
	runtime := &util.RuntimeOptions{}

	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		_response, _err := client.SendSmsWithOptions(sendSmsRequest, runtime)
		log.Debugf("send sms response:%v", _response)
		if _err != nil {
			return _err
		}
		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 如有需要，请打印 error
		_, _err := util.AssertAsString(error.Message)
		if _err != nil {
			return _err
		}
	}
	return nil
}
