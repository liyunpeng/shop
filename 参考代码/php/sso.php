<!DOCTYPE html> 
<body>
<?php 
	header('content-type:text/html;charset=utf-8'); 
	$phone = '13972165789';  // 改手机号测试 
	echo '测试的手机号：'.$phone. '<br>';

	$result = ssoUserCreate($phone);
	debugLog('创建大会员用户', $result);

	$result = ssoUserLogin($phone);
	debugLog('大会员用户登录', $result);

	$updateData = $arrayName = array(
		'address' => '地址1', 
		'nickname' => '名字1',
		'city' => '上海',
		);
	$result = ssoUserUpdate($phone, $updateData);
	debugLog('更新大会员用户', $result);

	$address=array(
		'receiver_name' => '姓名1',
		'receiver_phone' => '13972166666', 
		'province' => '廣東',
		'city' => '深圳',
		'area'  => '寶安',
		'address' => '深圳市寶安區龍華街道辦油松第十工業區東環二路二號',
		'type' => 0,
		'zip_code' => 10001, 
		'is_default' => 0,
		'status' => 1,
	);
	$address = '['.json_encode($address).']';
	$result = ssoUserAddaddress($phone, $address);
	debugLog('添加大会员用户地址', $result);


  function ssoUserAddaddress($phone, $address)
  {
  	$SSO_SERVER='https://iot.flnet.com/api';
		$BROKER='1707888012';
		$SECRET='LhDr8osERsWpfdnzAYDvgwS67TqvjrpnTzgiZMcO';

		$user_uuid = ssoUserLogin($phone)['user']['uuid'];			
		$url = $SSO_SERVER.'/user/'.$user_uuid.'/addresses';
		$method = 'POST';
		$postData = array(
			'client_uuid' => $BROKER,
			'client_secret' => $SECRET,
			'addresses' => $address,	
		);
		$ret = curlRequest($url, $method, $postData);
		return json_decode($ret, true);
  }

	function ssoUserCreate($phone)
	{
		try
		{
			$smsConfirmCode = getSmsConfirmCode($phone);
			if( !empty($smsConfirmCode['confirm_code'])) {
				return registerSsoUser($smsConfirmCode['confirm_code'], $phone);
			}else{
				debugLog('短信验证码错误', $smsConfirmCode);
				throw new Exception('no sms confirm code');
			}		
		}catch(Exception $e){
			return '创建大会员用户异常';
		}
	}

	function ssoUserLogin($phone)
	{
		$SSO_SERVER='https://iot.flnet.com/api';
		$BROKER='1707888012';
		$SECRET='LhDr8osERsWpfdnzAYDvgwS67TqvjrpnTzgiZMcO';
		$url = $SSO_SERVER.'/auth/login';
		$method = 'POST';
		$postData = array(
			'client_uuid' => $BROKER,
			'client_secret' => $SECRET,
			'username' => $phone,	
			'password' => '123456',
		);
		$ret = curlRequest($url, $method, $postData);
		return json_decode($ret, true);
	}
	
	function ssoUserUpdate($phone, $updateData)
	{
		$SSO_SERVER='https://iot.flnet.com/api';
		$BROKER='1707888012';
		$SECRET='LhDr8osERsWpfdnzAYDvgwS67TqvjrpnTzgiZMcO';
		$url = $SSO_SERVER.'/users/'.$phone;
		$method = 'PUT';
		$postData = array(
			'client_uuid' => $BROKER,
			'client_secret' => $SECRET,
			'mobile_country_code' => '+86',	
			'mobile' => $phone,
		);
		$postData = array_merge($postData, $updateData);
		$ret = curlRequest($url, $method, $postData);
		return json_decode($ret, true);		
	}

	function  getSmsConfirmCode($phone)
	{
		$SSO_SERVER='https://iot.flnet.com/api';
		$BROKER='1707888012';
		$SECRET='LhDr8osERsWpfdnzAYDvgwS67TqvjrpnTzgiZMcO';
		$url = $SSO_SERVER.'/sms/codes/'.$phone;
		$method = 'POST';
		$postData = array(
			'client_uuid' => $BROKER,
			'client_secret' => $SECRET,
			'type' => '1',
			'mobile_country_code' => '+86',
		);
		$ret = curlRequest($url, $method, $postData);
		return json_decode($ret, true);
	}
	
	function  registerSsoUser($smsConfirmCode, $phone)
	{
		$SSO_SERVER='https://iot.flnet.com/api';
		$BROKER='1707888012';
		$SECRET='LhDr8osERsWpfdnzAYDvgwS67TqvjrpnTzgiZMcO';
		$url = $SSO_SERVER.'/auth/register';
		$method = 'POST';
		$postData = array(
			'client_uuid' => $BROKER,
			'client_secret' => $SECRET,
			'mobile_country_code' => '+86',	
			'mobile_confirm_code' => strval($smsConfirmCode),
			'username' => $phone,
			'password' => '123456',
		);
		$ret = curlRequest($url, $method, $postData);
		return json_decode($ret, true);
	}

	function curlRequest($url, $method, $postData)
	{
		$curl = curl_init();
		curl_setopt($curl, CURLOPT_URL, $url);
		curl_setopt($curl, CURLOPT_HEADER, 0);
		curl_setopt($curl, CURLOPT_SSL_VERIFYPEER, 0);
		curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);
		curl_setopt($curl, CURLOPT_TIMEOUT, 30);
		if($method == 'POST'){
			curl_setopt($curl, CURLOPT_POST, 1);
			curl_setopt($curl, CURLOPT_POSTFIELDS, $postData);
		}else if($method == 'PUT'){
			curl_setopt($curl, CURLOPT_CUSTOMREQUEST, "PUT");
			curl_setopt($curl, CURLOPT_HTTPHEADER, array('Content-Type: application/x-www-form-urlencoded'));
			curl_setopt($curl, CURLOPT_POSTFIELDS, http_build_query($postData));
		}	
		$curlRet = curl_exec($curl);
		if (curl_errno($curl)) {
			throw new  Exception("curl error: ", curl_error($curl));
		}
		curl_close($curl);
		return $curlRet;
	}

	function debugLog($tag, $s)
	{	
		echo '<br> '.$tag.' dbginfo:<br>';
		var_dump($s);
	}
?>
</body>  
</html>  