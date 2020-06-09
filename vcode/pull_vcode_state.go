package main

import (
	"encoding/json"
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20190711"//引入sms
)

func main() {
	/* 必要步骤：
	 * 实例化一个认证对象，入参需要传入腾讯云账户密钥对 secretId 和 secretKey
	 * 本示例采用从环境变量读取的方式，需要预先在环境变量中设置这两个值
	 * 您也可以直接在代码中写入密钥对，但需谨防泄露，不要将代码复制、上传或者分享给他人
	 * CAM 密匙查询: https://console.cloud.tencent.com/cam/capi*/
	credential := common.NewCredential(
		// os.Getenv("TENCENTCLOUD_SECRET_ID"),
		// os.Getenv("TENCENTCLOUD_SECRET_KEY"),
		"xxx",
		"xxx",
	)
	/* 非必要步骤:
	 * 实例化一个客户端配置对象，可以指定超时时间等配置 */
	cpf := profile.NewClientProfile()

	/* SDK 默认使用 POST 方法
	 * 如需使用 GET 方法，可以在此处设置，但 GET 方法无法处理较大的请求 */
	cpf.HttpProfile.ReqMethod = "POST"

	/* SDK 有默认的超时时间，非必要请不要进行调整
	 * 如有需要请在代码中查阅以获取最新的默认值 */
	//cpf.HttpProfile.ReqTimeout = 5

	/* SDK 会自动指定域名，通常无需指定域名，但访问金融区的服务时必须手动指定域名
	 * 例如 SMS 的上海金融区域名为 sms.ap-shanghai-fsi.tencentcloudapi.com */
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"

	/* SDK 默认用 TC3-HMAC-SHA256 进行签名
	 * 非必要请不要修改该字段 */
	cpf.SignMethod = "HmacSHA1"

	/* 实例化 SMS 的 client 对象
	 * 第二个参数是地域信息，可以直接填写字符串 ap-guangzhou，或者引用预设的常量 */
	client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)

	/* 实例化一个请求对象，根据调用的接口和实际情况，可以进一步设置请求参数
	 * 您可以直接查询 SDK 源码确定接口有哪些属性可以设置
	 * 属性可能是基本类型，也可能引用了另一个数据结构
	 * 推荐使用 IDE 进行开发，可以方便地跳转查阅各个接口和数据结构的文档说明 */
	request := sms.NewPullSmsSendStatusRequest()

	/* 基本类型的设置:
	 * SDK 采用的是指针风格指定参数，即使对于基本类型也需要用指针来对参数赋值。
	 * SDK 提供对基本类型的指针引用封装函数
	 * 帮助链接：
	 * 短信控制台：https://console.cloud.tencent.com/smsv2
	 * sms helper：https://cloud.tencent.com/document/product/382/3773 */

	/* 短信应用 ID: 在 [短信控制台] 添加应用后生成的实际 SDKAppID，例如1400006666 */
	request.SmsSdkAppid = common.StringPtr("1400787878")
	/* 拉取最大条数，最多100条 */
	request.Limit = common.Uint64Ptr(10)

	// 通过 client 对象调用想要访问的接口，需要传入请求对象
	response, err := client.PullSmsSendStatus(request)
	// 处理异常
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	// 非 SDK 异常，直接失败。实际代码中可以加入其他的处理
	if err != nil {
		panic(err)
	}
	b, _ := json.Marshal(response.Response)
	// 打印返回的 JSON 字符串
	fmt.Printf("%s", b)
}