<?php

/**
 * BBC licence
 *
 * @copyright  Copyright (c) 2005-2010 BBC Technologies Inc. (http://www.BBC.cn)
 * @license    http://ecos.BBC.cn/ BBC License
 */
class topshop_ctl_passport extends topshop_controller
{

    //简单的登录页面
    public function simpleSignin()
    {
        return view::make('topshop/passport/simpleSignin.html');
    }

    /**
     * @brief 显示登录页面
     */
    public function signin()
    {
        $this->contentHeaderTitle = app::get('topshop')->_('企业账号登录');
        $this->set_tmpl('passport');
        if (pamAccount::isEnableVcode('sysshop')) {
            $pagedata['isShowVcode'] = 'true';
        }

        $pagedata['backgroundImgUrl'] = base_storager::modifier(app::get('sysconf')->getConf('sysconf_setting.shop_login_background'));
        $pagedata['backgroundImgUrl'] = $pagedata['backgroundImgUrl'] ?: app::get('topshop')->res_url . '/images/bj_01.jpg';
        $pagedata['openid'] = input::get('openid');
        $pagedata['oauthcode'] = input::get('oauthcode');

        return $this->page('topshop/passport/signin.html', $pagedata);
    }

    /**
     * @brief 会员登录处理
     * @return
     */
    public function login()
    {
        if (pamAccount::isEnableVcode('sysshop')) {
            // 验证图片验证码
            if (!base_vcode::verify(input::get('imagevcodekey'), input::get('imgcode'))) {
                $msg = app::get('topshop')->_('图片验证码错误');
                $url = url::action('topshop_ctl_passport@signin');
                return $this->splash('error', $url, $msg, true);
            }
        }
        try {
            $shop = app::get('sysshop')->model('shop')->getRow('seller_id', array('email' => input::get('login_account')));
            if ($shop) {
                $account = app::get('sysshop')->model('account')->getRow('login_account', array('seller_id' => $shop['seller_id']));
            } else {
                $url = url::action('topshop_ctl_passport@signin');
                return $this->splash('error', $url, app::get('topshop')->_('供应商账号不存在'), true);
            }
            shopAuth::login($account['login_account'], input::get('login_password'));
        } catch (Exception $e) {
            $url = url::action('topshop_ctl_passport@signin');
            $msg = $e->getMessage();
        }
        if (pamAccount::check()) {
            if (input::get('remember_me')) {
                setcookie('LOGINNAME', trim(input::get('login_account')), time() + 31536000, kernel::base_url() . '/');
            }

            $url = url::action('topshop_ctl_index@index');
            $msg = app::get('topshop')->_('登录成功');
            $this->sellerlog('账号登录。账号名是' . input::get('login_account'));

            if (request::ajax())
                return $this->splash('success', $url, $msg, true);
            else
                return redirect::to($url);
        } else {
            return $this->splash('error', $url, $msg, true);
        }

    }

    /**
     * @brief 显示登录注册
     */
    public function signup()
    {
        $this->contentHeaderTitle = app::get('topshop')->_('企业账号注册');
        $this->set_tmpl('pwdfind');
        $pagedata['license'] = app::get('sysshop')->getConf('sysshop.register.setting_sysshop_license');
        return $this->page('topshop/passport/signup.html', $pagedata);
    }

    public function isExists()
    {
        switch (input::get('type')) {
            case 'account':
                $str = input::get('login_account');
                break;
            case 'mobile':
                $str = input::get('mobile');
                break;
            case 'email':
                $str = input::get('email');
                break;
        }
        $flag = shopAuth::isExists($str, input::get('type'));
        $status = $flag ? 'false' : 'true';
        return $this->isValidMsg($status);
    }

    /**
     * @brief 创建商家会员
     * @return json
     */
    public function create()
    {
        if (input::get('license') != 'on') {
            $msg = $this->app->_('同意注册条款后才能注册');
            throw new \LogicException($msg);
        }

        try {
            $request = input::get();
            if (!$_SESSION['register']['mobile'])
                throw new LogicException('手机号码错误');
            $request['mobile'] = $_SESSION['register']['mobile'];
            $request['auth_type'] = 'AUTH_MOBILE';
            shopAuth::signupSeller($request);
        } catch (Exception $e) {
            $msg = $e->getMessage();
        }

        if (pamAccount::check()) {
            $url = url::action('topshop_ctl_index@index');
            $msg = app::get('topshop')->_('注册成功');
            return $this->splash('success', $url, $msg, true);
        } else {
            return $this->splash('error', null, $msg, true);
        }
    }

    public function logout()
    {
        pamAccount::logout();
        return redirect::action('topshop_ctl_passport@signin');
    }

    /**
     * @brief 会员密码修改
     */
    public function update()
    {
        return view::make('topshop/passport/update.html');
    }

    public function updatepwd()
    {
        try {
            shopAuth::modifyPwd(input::get());
        } catch (Exception $e) {
            $msg = $e->getMessage();
            return $this->splash('error', null, $msg, true);
        }

        $this->sellerlog('修改当前账号密码。');
        $url = url::action('topshop_ctl_passport@signin');
        $msg = app::get('topshop')->_('修改成功,请重新登陆');
        pamAccount::logout();

        return $this->splash('success', $url, $msg, true);
    }

    public function initialize()
    {
        if (empty($_GET) || $_GET['sign'] != 'activate') {
            $pagedata['status'] = 'error';
            $pagedata['msg'] = app::get('topc')->_('未定义的操作');
        }

        $shopId = substr($_GET['s'], 4, strlen($_GET['s']) - 4);
        $shop = app::get('sysshop')->model('shop')->getRow('seller_id', array('shop_id' => $shopId));
        $res = app::get('sysshop')->model('account')->getRow('*', array('seller_id' => $shop['seller_id']));
        $pagedata['account'] = $res;

        if (!$res || $res['active']) {
            $pagedata['status'] = 'error';
            $pagedata['msg'] = app::get('topc')->_('该链接已失效');
        }

        $seller = app::get('sysshop')->model('seller')->getRow('email', array('seller_id' => $shop['seller_id']));
        if (md5($seller['email'] . '1234') != $_GET['e']) {
            $pagedata['status'] = 'error';
            $pagedata['msg'] = app::get('topc')->_('您所用的账户不存在');
        } else {
            $pagedata['account']['email'] = $seller['email'];
        }

        $pagedata['backgroundImgUrl'] = base_storager::modifier(app::get('sysconf')->getConf('sysconf_setting.shop_login_background'));
        $pagedata['backgroundImgUrl'] = $pagedata['backgroundImgUrl'] ?: app::get('topshop')->res_url . '/images/bj_01.jpg';
        $this->set_tmpl('passport');
        return $this->page('topshop/passport/initialize.html', $pagedata);
    }

    public function replace()
    {
        $post = utils::_filter_input(input::get());
        try {
            if (!$post['seller_id']) {
                throw new \LogicException(app::get('topc')->_('供应商账号不存在'));
            }

            if (!app::get('sysshop')->model('seller')->getRow('*', array('seller_id' => $post['seller_id']))) {
                throw new \LogicException(app::get('topc')->_('供应商账号不存在'));
            }

            $sellerId = $post['seller_id'];
            $data = array();
            $data['login_password'] = $post['login_password'];
            $data['psw_confirm'] = $post['psw_confirm'];
            shopAuth::resetPwd($sellerId, $data);

            $pamShopData['seller_id'] = $post['seller_id'];
            $pamShopData['active'] = '1';
            $pamShopData['modified_time'] = time();
            app::get('sysshop')->model('account')->save($pamShopData);

            $seller = [
                'seller_od' => $post['seller_id'],
                'auth_type' => 'AUTH_EMAIL',
            ];
            app::get('sysshop')->model('seller')->save($seller);
        } catch (Exception $e) {
            $msg = $e->getMessage();
            return $this->splash('error', null, $msg);
        }

        $url = url::action('topshop_ctl_passport@signin');
        return $this->splash('success', $url, app::get('sysuser')->_('该帐号已被激活，谢谢！'));
    }
}

