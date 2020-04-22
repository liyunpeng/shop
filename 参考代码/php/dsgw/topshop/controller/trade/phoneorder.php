0<?php
class topshop_ctl_trade_phoneorder extends topshop_controller { 

    public function index()
    {
        $this->contentHeaderTitle = app::get('topshop')->_('订单录入');
        return $this->page('topshop/trade/phoneorder/edit.html', $pagedata);
    }

    public function loadAddress()
    {
        $userAccount = input::get('accountPhone').'tel';
        try{
            if(kernel::single('sysuser_passport')->isExistsAccount($userAccount)){
                $userFilter['login_account'] = $userAccount;
                $userId = app::get('sysuser')->model('account')->getRow('user_id',$userFilter)['user_id'];
                $addressParam['user_id'] = intval($userId);
                $addressParam['def_addr'] = 1;
                $userAddrList = app::get('topshop')->rpcCall('user.address.list',$addressParam)['list'];

                $addressInfo = $userAddrList[0];
                list($regions,$region_id) = explode(':', $addressInfo['area']);
                $addressInfo['area'] = $regions;
                $addressInfo['region_id'] = str_replace('/', ',', $region_id);
                return response::json($addressInfo);    
            }else{
                return 'error';
            }
        }catch(Exception $e){
            $this->log("加载地址:".$e->getMessage());
            return 'error';
        }


    }
    
    public function loadConfirmModal()
    {
        return view::make('topshop/trade/phoneorder/confirm_dialog.html');
    }

    public function loadModalContent()
    {
        $pagedata = input::get();
        $total_price = 0;
        $total_delivery_fee = 0;
        $total_real_payment = 0;
        $skuId = implode(',', array_values($pagedata['sku_id']));
        $skuNum = $pagedata['sku_num'];

        try
        {
            $objLibMath = kernel::single('ectools_math');
            $skuInfo = app::get('topshop')->rpcCall('sku.list', ['sku_ids'=>$skuId, 'fields'=>'*']);

            foreach ($skuInfo as $key => $value) {
                $skuInfo[$key]['quantity']= intval($skuNum[$value['sku_id']]);

                $skuInfo[$key]['subtotal_price'] = $objLibMath->number_multiple(array($skuInfo[$key]['price'], $skuInfo[$key]['quantity']));

                $skuInfo[$key]['delivery_fee']  = $this->getDelivery(intval($skuInfo[$key]['item']['dlytmpl_id']), $skuInfo[$key]['quantity'], intval($key));

                $skuInfo[$key]['subtotal_payment'] = $objLibMath->number_plus(array($skuInfo[$key]['subtotal_price'], $skuInfo[$key]['delivery_fee']));

                $total_price = $objLibMath->number_plus(array($total_price, $skuInfo[$key]['subtotal_price']));

                $total_delivery_fee =  $objLibMath->number_plus(array($total_delivery_fee, $skuInfo[$key]['delivery_fee']));

                $total_real_payment = $objLibMath->number_plus(array($total_real_payment,$skuInfo[$key]['subtotal_payment']));
            }
            $pagedata['sku_info']= $skuInfo;
            $pagedata['total_price'] = $total_price;
            $pagedata['total_delivery_fee'] = $total_delivery_fee;
            $pagedata['total_real_payment'] = $total_real_payment;

            $area = app::get('topshop')->rpcCall('logistics.area',array('area'=>$pagedata['area'][0]));
            $pagedata['area_addr'] = str_replace('/','', $area).' '.$pagedata['addr'];
        }catch(Exception $e){
            $this->log('电话订单对话框异常：'.$e->getMessage());
        }

        return  view::make('topshop/trade/phoneorder/confirm_dialog_content.html', $pagedata);
    }

    public function getCounpon1($userId)
    {
        $params['user_id'] = $userId;

        $couponCode = kernel::single('systrade_cart_coupon_redis')->get($params['user_id'], $this->shop_id);

        $this->log('coupon: '.$couponCode);

        if ($couponCode) {
            $userCoupon = $userCouponModel->getRow('coupon_id', ['coupon_code' => $couponCode]);
            $coupon = $couponModel->getRow('coupon_name', ['coupon_id' => $userCoupon['coupon_id']]);
            $cartItem['coupon_name'] = $coupon['coupon_name'];
            $this->log('coupon: '.$coupon);
        }
    }
     

    public function loadCoupon()
    {

        $userAccount = input::get('accountPhone').'tel';
        try{
            if(kernel::single('sysuser_passport')->isExistsAccount($userAccount)){
                $userFilter['login_account'] = $userAccount;
                $userId = app::get('sysuser')->model('account')->getRow('user_id',$userFilter)['user_id'];

                $params = array(
                    'page_no' => 0,
                    'page_size' => 100,
                    'fields' => '*',
                    'user_id' => intval($userId),
                    // 'shop_id' => intval($shop_id),
                    // 'is_valid' => 1,
                    // 'platform' => 'app',
                    );

                $couponListData = app::get('topapi')->rpcCall('user.coupon.list', $params);

                $this->log($couponListData);
                
                $couponList = $couponListData['coupons'];
                // $this->log('coupon1: '.response::json($couponList));
                return response::json($couponList);
                // return $couponList;
            }else{
                return 'error';
            }
        }catch(Exception $e){
            $this->log("加载优惠券:".$e->getMessage());
            return 'error';
        }

    }

    public function createUser($userAccount)
    {
        if ( kernel::single('sysuser_passport')->isExistsAccount($userAccount) ) {
            $userFilter['login_account'] = $userAccount;
            $userId = app::get('sysuser')->model('account')->getRow('user_id',$userFilter)['user_id'];
        }else{
            $user['account'] = $userAccount;
            $user['password'] = 'admin123';
            $user['pwd_confirm'] = 'admin123';
            $user['reg_ip'] = request::getClientIp();
            $userId = app::get('topshop')->rpcCall('user.create',$user,'buyer');
        }

        $userId=intval($userId);        

        return $userId;
    }

    public function createAddress($userId, $param)
    {
        $param['user_id'] = $userId;
        $area = app::get('topshop')->rpcCall('logistics.area',array('area'=>$param['area'][0]));
        $validator = validator::make(
            [
             'area' => $area,
             'addr' => $param['addr'] ,
             'name' => $param['name'],
             'mobile' => $param['mobile'],
             'user_id' =>$param['user_id'],
             'zip' =>intval($param['zip']),
            ],
            [
            'area' => 'required|max:20',
            'addr' => 'required',
            'name' => 'required',
            'mobile' => 'required|mobile',
            'user_id' => 'required',
             'zip' =>'numeric|max:999999',
            ],
            [
             'area' => '地区不存在!',
             'addr' => '会员街道地址必填!',
             'name' => '收货人姓名未填写!',
             'mobile' => '手机号码必填!|手机号码格式不正确!',
             'user_id' => '缺少参数!',
             'zip' =>'邮编必须为6位数的整数|邮编最大为999999',
            ]
        );
        if ($validator->fails())
        {
            $messages = $validator->messagesInfo();

            foreach( $messages as $error )
            {
                throw new RuntimeException('地址验证错误'.$error[0]);
            }
        }
        $areaId =  str_replace(",","/", $param['area'][0]);
        $param['area'] = $area . ':' . $areaId;
        $param['def_addr'] = 1;     
        app::get('topshop')->rpcCall('user.address.add',$param,'buyer');

        $addressParam['user_id'] = $userId;
        $addressParam['def_addr'] = 1;
        $userDefAddr = app::get('topshop')->rpcCall('user.address.list',$addressParam)['list'][0];
        return intval($userDefAddr['addr_id']);
    }

    public function createCart($skuItems, $userId)
    {
        foreach( $skuItems  as $skuId=>$skuNum ){
            $addCartParam['user_id'] = $userId;
            $addCartParam['quantity'] = $skuNum;
            $addCartParam['sku_id'] = $skuId;
            $addCartParam['obj_type'] = 'item';
     
            $addCartParam['mode'] = 'cart';           
            $addCartRet = app::get('topshop')->rpcCall('trade.cart.add', $addCartParam, 'buyer');
                
            $updateCartParam['mode'] = 'cart';
            $updateCartParam['obj_type'] = 'item';
            $updateCartParam['cart_id'] = $addCartRet['cart_id'];
            $updateCartParam['totalQuantity'] = $addCartParam['quantity'];
            $updateCartParam['user_id'] = $userId;
            $updateCartParam['is_checked'] = '1';
            $updateCartRet = app::get('topc')->rpcCall('trade.cart.update',$updateCartParam);
        }
        return app::get('topshop')->rpcCall('trade.cart.getBasicCartInfo',
           array('user_id'=>$userId, 'platform'=>'pc', 'mode'=>'cart','checked'=>1), 'buyer');
    }

    public function createOrder($userId, $userAccount, $addressId, $cartInfo, $mark)
    {
        $orderParam['user_id']   = $userId;
        $orderParam['user_name'] =  $userAccount;
        $orderParam['addr_id']   = $addressId;
        $markArray[] = array(
            'shop_id'=>$this->shopId,
            'memo'=> $mark,
            );
        $orderParam['mark'] = json_encode($markArray);
        $orderParam['payment_type']   = 'offline';
        $orderParam['mode']      = 'cart';
        $orderParam['source_from'] = 'phone';
        $orderParam['invoice_type']  = 'notuse';
        $orderParam['shipping_type'][] = ['shop_id' =>  $this->shopId, 'type'=> 'express'];
        $orderParam['shipping_type'] = json_encode($orderParam['shipping_type']);
        $orderParam['md5_cart_info'] = md5(serialize(
            utils::array_ksort_recursive($cartInfo)));

        $this->getCounpon1($userId);
        
        return  app::get('topshop')->rpcCall('trade.create',$orderParam);
    }

    public function deleteCart($userId, $cartInfo)
    {
        $cartIds = array_column($cartInfo, 'cart_id');
        $cartIds = implode(',', $cartIds);     
        return  app::get('systrade')->rpcCall('trade.cart.delete', array('cart_id'=>$cartIds,'mode'=>'cart','user_id'=>$userId ));
    }

    public function save()
    {
        $inputParam = input::get();

        $userAccount = $inputParam['accountPhone'].'tel';

        try
        {
            $userId = $this->createUser($userAccount);
            
            $addressId = $this->createAddress($userId, $inputParam);

            $cartInfo = $this->createCart($inputParam['sku_num'], $userId);

            $this->createOrder($userId, $userAccount, $addressId, $cartInfo, $inputParam['mark']);

            $this->deleteCart($userId, $cartInfo);

        }catch(Exception $e){
            $this->log('电话订单创建异常:'.$e->getMessage());
            $msg = app::get('topshop')->_($e->getMessage());
            return $this->splash('error',$url,$msg,true);
        }

        $msg = app::get('topshop')->_('电话订单创建成功');
        $url = url::action('topshop_ctl_trade_list@index');
        return $this->splash('success',$url,$msg,true);
    }

    public function getDelivery($dlytmplId, $quantity, $key)
    {
        $objLibMath = kernel::single('ectools_math');
        $dlytmplInfo = app::get('sysitem')->rpcCall('logistics.dlytmpl.get', ['template_id'=>$dlytmplId, 'fields'=>'*']);
        $deliveryFee = 0;

        if( $key == 0){
            if( $quantity == 1){
                $deliveryFee =  $dlytmplInfo['fee_conf'][0]['start_fee'] ;
            }else if ($quantity > 1) {
                $addFee = $objLibMath->number_multiple(array($dlytmplInfo['fee_conf'][0]['add_fee'], $quantity-1));
                $deliveryFee = $objLibMath->number_plus(array( $dlytmplInfo['fee_conf'][0]['start_fee'], $addFee));
            }
        }else if ($key > 0){
            $deliveryFee = $objLibMath->number_multiple(array($dlytmplInfo['fee_conf'][0]['add_fee'], $quantity));
        }
        return $deliveryFee;
    }

    public function log($s)
    {
        logger::info('[phoneorderLogTag]');
        var_dump($s);
    }
}