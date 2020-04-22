<?php

class topshop_ctl_promotion_coupon extends topshop_controller
{

    public function list_coupon()
    {
        $this->contentHeaderTitle = app::get('topshop')->_('优惠券管理');
        $filter = input::get();
        if (!$filter['pages']) {
            $filter['pages'] = 1;
        }
        $pageSize = 10;
        $params = array(
            'page_no' => $filter['pages'],
            'page_size' => intval($pageSize),
            'fields' => '*',
            'shop_id' => $this->shopId,
        );
        $couponListData = app::get('topshop')->rpcCall('promotion.coupon.list', $params, 'seller');
        $count = $couponListData['count'];
        $pagedata['couponList'] = $couponListData['coupons'];

        //处理翻页数据
        $current = $filter['pages'] ? $filter['pages'] : 1;
        $filter['pages'] = time();
        if ($count > 0) $total = ceil($count / $pageSize);
        $pagedata['pagers'] = array(
            'link' => url::action('topshop_ctl_promotion_coupon@list_coupon', $filter),
            'current' => $current,
            'total' => $total,
            'use_app' => 'topshop',
            'token' => $filter['pages'],
        );

        $pagedata['now'] = time();
        $pagedata['total'] = $count;
        $pagedata['examine_setting'] = app::get('sysconf')->getConf('shop.promotion.examine');

        return $this->page('topshop/promotion/coupon/index.html', $pagedata);
    }

    public function edit_coupon()
    {
        $this->contentHeaderTitle = app::get('topshop')->_('新添/编辑优惠券');
        $apiData['coupon_id'] = input::get('coupon_id');
        $apiData['coupon_itemList'] = true;
        $pagedata['valid_time'] = date('Y/m/d H:i', time() + 60) . '-' . date('Y/m/d H:i', time() + 120); //默认时间
        $pagedata['cansend_time'] = date('Y/m/d H:i', time() + 60) . '-' . date('Y/m/d H:i', time() + 120); //默认时间
        if ($apiData['coupon_id']) {
            $pagedata = app::get('topshop')->rpcCall('promotion.coupon.get', $apiData);
            $pagedata['valid_time'] = date('Y/m/d H:i', $pagedata['canuse_start_time']) . '-' . date('Y/m/d H:i', $pagedata['canuse_end_time']);
            $pagedata['cansend_time'] = date('Y/m/d H:i', $pagedata['cansend_start_time']) . '-' . date('Y/m/d H:i', $pagedata['cansend_end_time']);
            if ($pagedata['shop_id'] != $this->shopId) {
                return $this->splash('error', '', '您没有权限编辑此优惠券', true);
            }
            $notItems = array_column($pagedata['itemsList'], 'item_id');
            $pagedata['notEndItem'] = json_encode($notItems, true);
        }

        $valid_grade = explode(',', $pagedata['valid_grade']);
        $pagedata['gradeList'] = app::get('topshop')->rpcCall('user.grade.list');
        foreach ($pagedata['gradeList'] as &$v) {
            if (in_array($v['grade_id'], $valid_grade)) {
                $v['is_checked'] = true;
            }
        }
        foreach ($pagedata['itemsList'] as $itemRow) {
            $pagedata['item_sku'][$itemRow['item_id']] = $itemRow['sku_ids'];
        }
        return $this->page('topshop/promotion/coupon/edit.html', $pagedata);
    }

    //查看优惠券
    public function show_coupon()
    {
        $this->contentHeaderTitle = app::get('topshop')->_('查看优惠券');
        $apiData['coupon_id'] = input::get('coupon_id');
        $apiData['coupon_itemList'] = true;
        if ($apiData['coupon_id']) {
            $pagedata = app::get('topshop')->rpcCall('promotion.coupon.get', $apiData);
            $pagedata['valid_time'] = date('Y/m/d H:i', $pagedata['canuse_start_time']) . ' ~ ' . date('Y/m/d H:i', $pagedata['canuse_end_time']);
            $pagedata['send_time'] = date('Y/m/d H:i', $pagedata['cansend_start_time']) . ' ~ ' . date('Y/m/d H:i', $pagedata['cansend_end_time']);
            if ($pagedata['shop_id'] != $this->shopId) {
                return $this->splash('error', '您没有权限查看此优惠券', true);
            }
            $notItems = array_column($pagedata['itemsList'], 'item_id');
            $pagedata['notEndItem'] = json_encode($notItems, true);
        }

        $valid_grade = explode(',', $pagedata['valid_grade']);
        $pagedata['gradeList'] = app::get('topshop')->rpcCall('user.grade.list');
        $gradeIds = array_column($pagedata['gradeList'], 'grade_id');
        if (!array_diff($gradeIds, $valid_grade)) {
            $gradeStr = ' 所有会员';
        } else {
            foreach ($pagedata['gradeList'] as $member) {
                if (in_array($member['grade_id'], $valid_grade)) {
                    $gradeStr .= $member['grade_name'] . ',';
                }
            }
            $gradeStr = rtrim($gradeStr, ',');
        }
        $pagedata['grade_str'] = $gradeStr;
        $pagedata['ac'] = input::get('ac', '');
        return $this->page('topshop/promotion/coupon/show.html', $pagedata);
    }

    public function save_coupon()
    {
        $params = input::get();

        if (!$params['image_id'] || !$params['detail_image_id'] || !$params['detail_received_image_id'] || !$params['other_image_id'] || !$params['received_image_id']) {
            $url = url::action('topshop_ctl_promotion_coupon@edit_coupon', array('coupon_id' => $params['coupon_id']));
            return $this->splash('error', $url, '请上传优惠券图片', true);
        }

        $apiData = $params;
        $apiData['shop_id'] = $this->shopId;
        // 可使用的有效期
        $canuseTimeArray = explode('-', $params['valid_time']);
        $apiData['canuse_start_time'] = strtotime($canuseTimeArray[0]);
        $apiData['canuse_end_time'] = strtotime($canuseTimeArray[1]);
        // 可以领取的时间段
        $cansendTimeArray = explode('-', $params['cansend_time']);
        $apiData['cansend_start_time'] = strtotime($cansendTimeArray[0]);
        $apiData['cansend_end_time'] = strtotime($cansendTimeArray[1]);
        // 可以使用的会员等级
        $apiData['valid_grade'] = implode(',', $params['grade']);
        $couponRelItem = null;
        foreach ((array)$params['item_id'] as $key => $itemId) {
            $itemData['item_id'] = $itemId;
            $itemData['sku_id'] = $params['item_sku'][$key];
            $couponRelItem[] = $itemData;
        }
        $apiData['coupon_rel_item'] = $couponRelItem ? json_encode($couponRelItem) : null;

        try {
            if ($params['coupon_id']) {
                // 修改优惠券
                $result = app::get('topshop')->rpcCall('promotion.coupon.update', $apiData);
            } else {
                // 新添优惠券
                $result = app::get('topshop')->rpcCall('promotion.coupon.add', $apiData);
            }
        } catch (\LogicException $e) {
            $msg = $e->getMessage();
            $url = url::action('topshop_ctl_promotion_coupon@edit_coupon', array('coupon_id' => $params['coupon_id']));
            return $this->splash('error', $url, $msg, true);
        }
        $this->sellerlog('添加/修改优惠券。优惠券名称是 ' . $apiData['coupon_name']);
        $url = url::action('topshop_ctl_promotion_coupon@list_coupon');
        $msg = app::get('topshop')->_('保存优惠券成功');
        return $this->splash('success', $url, $msg, true);
    }

    public function submit_approve()
    {
        $apiData = input::get();
        try {
            $couponInfo = app::get('topshop')->rpcCall('promotion.coupon.get', $apiData);
            if ($couponInfo['cansend_end_time'] <= time()) {
                throw new \LogicException('您的活动已过期，无法提交审核!');
            }
            $result = app::get('topshop')->rpcCall('promotion.coupon.approve', $apiData);
        } catch (\LogicException $e) {
            $msg = $e->getMessage();
            return $this->splash('error', $url, $msg, true);
        }
        $this->sellerlog('更新优惠券。优惠券ID是 ' . $apiData['coupon_id']);
        $url = url::action('topshop_ctl_promotion_coupon@list_coupon');
        $msg = app::get('topshop')->_('提交审核成功');
        return $this->splash('success', $url, $msg, true);
    }

    public function delete_coupon()
    {
        $apiData['shop_id'] = $this->shopId;
        $apiData['coupon_id'] = input::get('coupon_id');
        $url = url::action('topshop_ctl_promotion_coupon@list_coupon');
        try {
            app::get('topshop')->rpcCall('promotion.coupon.delete', $apiData);
        } catch (\LogicException $e) {
            $msg = $e->getMessage();
            return $this->splash('error', $url, $msg, true);
        }
        $this->sellerlog('删除优惠券。优惠券ID是 ' . $apiData['coupon_id']);
        $msg = app::get('topshop')->_('删除优惠券成功');
        return $this->splash('success', $url, $msg, true);
    }

    public function cancel_coupon()
    {
        $apiData['shop_id'] = $this->shopId;
        $apiData['coupon_id'] = input::get('coupon_id');
        $url = url::action('topshop_ctl_promotion_coupon@list_coupon');
        try {
            app::get('topshop')->rpcCall('promotion.coupon.cancel', $apiData);
        } catch (\LogicException $e) {
            $msg = $e->getMessage();
            return $this->splash('error', $url, $msg, true);
        }
        $this->sellerlog('取消优惠券促销。优惠券促销ID是 ' . $apiData['coupon_id']);
        $msg = app::get('topshop')->_('取消优惠券促销成功');
        return $this->splash('success', $url, $msg, true);
    }

    //根据商家id和3级分类id获取商家所经营的所有品牌
    public function getBrandList()
    {
        $shopId = $this->shopId;
        $catId = input::get('catId');
        $params = array(
            'shop_id' => $shopId,
            'cat_id' => $catId,
            'fields' => 'brand_id,brand_name,brand_url',
        );
        $brands = app::get('topshop')->rpcCall('category.get.cat.rel.brand', $params);
        return response::json($brands);
    }

    public function ajaxCouponList()
    {
        // 店铺优惠券信息,
        $params = array(
            'page_no' => input::get('pages', 1),
            'page_size' => input::get('page_size', 10),
            'fields' => 'deduct_money,coupon_name,coupon_id,canuse_start_time,canuse_end_time,cansend_start_time,cansend_end_time',
            'shop_id' => $this->shopId,
            'platform' => input::get('platform', 'pc'),
            'is_cansend' => 1,
        );

        $filter = input::get();

        $couponListData = app::get('topc')->rpcCall('promotion.coupon.list', $params);
        $count = $couponListData['count'];
        $pagedata['list'] = $couponListData['coupons'];

        //处理翻页数据
        $current = $filter['pages'] ? $filter['pages'] : 1;
        $filter['pages'] = time();
        if ($count > 0) $total = ceil($count / $pageSize);
        $pagedata['pagers'] = array(
            'link' => url::action('topshop_ctl_promotion_coupon@ajaxCouponList', $filter),
            'current' => $current,
            'total' => $total,
            'use_app' => 'topshop',
            'token' => $filter['pages'],
        );

        $pagedata['now'] = time();
        $pagedata['total'] = $count;
        return view::make('topshop/promotion/coupon/ajaxList.html', $pagedata);
    }
}

