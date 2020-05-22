OAuth:一个关于授权（authorization）的开放网络标准，目前版本是2.0版。


### 为何要使用OAuth协议呢？OAuth协议的应用场景。
第三方服务方提供服务，某些服务需要用户的同意才能够做到，好比客厅要装修，需要得到主人的同意，拿到钥匙，才能装修，提供服务。
传统做法：
把所有钥匙（账号密码）给工人。但这样，工人可能用这个钥匙开卧室的门。甚至打一个新的钥匙。
缺点：（不安全）
* 服务提供方只是提供一个服务，为了保证服务能提供，就会保存账号密码以供下次提供，这显然不安全。
* 服务提供方有了账号密码，就拥有了用户所有的权利，用户没办法限制服务提供方获得权限的范围和有效期。
* 用户只有修改密码，才能收回赋给服务方的权力。这样做，会使得其他所有获得用户授权的第三方全部失效。
* 只要有一个第三方应用程序被破解，就会导致用户密码泄漏，以及所有被密码保护的数据泄漏


OAuth可以解决这些问题。

### oauth2.0 概念名词：
* Third-party application：第三方应用程序（client）。
* HTTP service：HTTP服务提供商。
* Resource Owner：资源所有者-"用户"（user）。
* User Agent：用户代理-浏览器。
* Authorization server：认证服务器，即服务提供商专门用来处理认证的服务器。
* Resource server：资源服务器，即服务提供商存放用户生成的资源的服务器。它与认证服务器，可以是同一台服务器，也可以是不同的服务器。
OAuth的作用就是让"客户端"安全可控地获取"用户"的授权，与"服务商提供商"进行互动。

### Oauth设计理念
* OAuth在"客户端"与"服务提供商"之间，设置一个授权层（authorization layer）。
* "客户端"不能直接登录"服务提供商"，只能登录授权层，以此将用户与客户端区分开来。
* "客户端"登录授权层所用的令牌（token），与用户的密码不同。用户可在登录时，指定授权层令牌的权限范围和有效期。
* "客户端"登录授权层以后，"服务提供商"根据令牌的权限范围和有效期，向"客户端"开放用户储存的资料。

### OAuth 2.0的运行流程
* 1.用户打开客户端，客户端要求用户给予授权。
* 2.用户同意给予客户端授权。
* 3.客户端使用上一步获得的授权（一般是Code），向认证服务器申请令牌TOKEN。
* 4.认证服务器对客户端进行认证以后，确认无误，同意发放令牌。
* 5.客户端使用令牌，向资源服务器申请获取资源（用户信息等）。
* 6.资源服务器确认令牌无误，同意向客户端开放资源。
重点是如何获取Token

### 客户端获取授权的五种模式：
客户端必须得到用户的授权（authorization grant），才能获得令牌（access token）。

OAuth 2.0定义了五种授权方式：
* 授权码模式（authorization code）
* 简化模式（implicit）
* 密码模式（resource owner password credentials）
* 客户端模式（client credentials）
* 扩展模式（Extension）

#### 授权码模式（authorization code）：
授权码模式是功能最完整、流程最严密的授权模式。
特点是通过客户端的后台服务器，与"服务提供商"的认证服务器进行互动。

以微信公众平台公众号网页应用开发流程为例。步骤如下：
* （A）用户访问客户端，客户端将用户导向认证服务器。
* （B）用户选择是否给予客户端授权。
* （C）若用户给予授权，认证服务器将用户导向客户端指定的"重定向URI"（redirection URI），同时附上授权码code。
* （D）客户端收到授权码code，附上早先的"重定向URI"，向认证服务器申请token。这一步是在客户端的后台的服务器上完成的，对用户不可见。
* （E）认证服务器核对了授权码和重定向URI，确认无误后，向客户端发送访问令牌（access token）和更新令牌（refresh token）

一些重要参数：
* response_type：表示授权类型，必选项，此处的值固定为"code"
* appid：表示客户端的ID，必选项
* redirect_uri：表示重定向URI，可选项
* scope：表示申请的权限范围，可选项
* state：表示客户端的当前状态，可以指定任意值，认证服务器会原封不动地返回这个值。用于防止恶意攻击
 

1.引导用户跳转到授权页面：
```
https://open.weixin.qq.com/connect/oauth2/authorize?appid=APPID&redirect_uri=REDIRECT_URI&response_type=code&scope=SCOPE&state=STATE#wechat_redirect
```
参数：
* appid      公众号的唯一标识
* redirect_uri    授权后重定向的回调链接地址， 请使用 urlEncode 对链接进行处理
* response_type      返回类型，请填写code
* scope     应用授权作用域，有snsapi_base 、snsapi_userinfo 两种
* state       重定向后会带上state参数，开发者可以填写a-zA-Z0-9的参数值，最多128字节

2.通过授权码code获取Token
```
https://api.weixin.qq.com/sns/oauth2/access_token?appid=APPID&secret=SECRET&code=CODE&grant_type=authorization_code
```
参数：
* appid      公众号的唯一标识
* secret     公众号的appsecret
* code       填写获取的code参数（存在有效期，通常设为10分钟，客户端只能使用该码一次，否则会被授权服务器拒绝。该码与客户端ID和重定向URI，是一一对应关系）
* grant_type     填写为authorization_code

返回结果：
```
{
  "access_token":"ACCESS_TOKEN", //网页授权接口调用凭证,注意：此access_token与基础支持的access_token不同
  "expires_in":7200,  // access_token接口调用凭证超时时间，单位（秒）
  "refresh_token":"REFRESH_TOKEN", //用户刷新access_token
  "openid":"OPENID",  //用户唯一标识
  "scope":"SCOPE"  //用户授权的作用域，使用逗号（,）分隔
} 
```
返回结果解释：
* access_token：表示访问令牌，必选项。
* token_type：表示令牌类型，该值大小写不敏感，必选项，可以是bearer类型或mac类型。
* expires_in：表示过期时间，单位为秒。如果省略该参数，必须其他方式设置过期时间。
* refresh_token：表示更新令牌，用来获取下一次的访问令牌，可选项。
* scope：表示权限范围，如果与客户端申请的范围一致，此项可省略。

####  其他四种模式，用的很少，略过
