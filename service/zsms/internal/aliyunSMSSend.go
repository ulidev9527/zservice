package internal

import (
	"zservice/zservice"
	"zservice/zservice/zglobal"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type aliyunSMSSendConfig struct {
	Phone        string
	VerifyCode   string
	Key          string
	Secret       string
	TemplateCode string
	SignName     string
}

func aliyunSMSSend(ctx *zservice.Context, c *aliyunSMSSendConfig) (code uint32) {
	// 验证
	if c.Phone == "" ||
		c.VerifyCode == "" ||
		c.Key == "" ||
		c.Secret == "" ||
		c.TemplateCode == "" ||
		c.SignName == "" {
		return zglobal.Code_Zsms_SendParamsErr
	}

	// 请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID 和 ALIBABA_CLOUD_ACCESS_KEY_SECRET��
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例使用环境变量获取 AccessKey 的方式进行调用，仅供参考，建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html
	client, e := func() (_result *dysmsapi20170525.Client, _err error) {
		config := &openapi.Config{
			// 必填，您的 AccessKey ID
			AccessKeyId: tea.String(c.Key),
			// 必填，您的 AccessKey Secret
			AccessKeySecret: tea.String(c.Secret),
			Endpoint:        tea.String("dysmsapi.aliyuncs.com"),
		}

		_result = &dysmsapi20170525.Client{}
		_result, e := dysmsapi20170525.NewClient(config)
		return _result, e
	}()

	if e != nil {
		ctx.LogError(e)
		return zglobal.Code_Zsms_ErrorBreakoff
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(c.Phone),
		TemplateCode:  tea.String(c.TemplateCode),
		TemplateParam: tea.String("{\"code\":\"" + c.VerifyCode + "\"}"),
		SignName:      tea.String(c.SignName),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		_res, _err := client.SendSmsWithOptions(sendSmsRequest, runtime)

		if _err != nil {
			ctx.LogError(_e)
			return _err
		} else {
			ctx.LogInfo(_res.String())
		}

		return nil
	}()

	if tryErr != nil {
		return zglobal.Code_Zsms_ErrorBreakoff
	}
	return zglobal.Code_SUCC
}
