<?php
class topshop_ctl_selector_sku extends topshop_controller {
	public $limit = 18;

	public function loadSelectSkuModal()
    {
        $isImageModal = true;
        $pagedata['imageModal'] = true;
        $pagedata['textcol'] = input::get('textcol');
        $pagedata['view'] = input::get('view');
        $pagedata['shopCatList'] = app::get('topshop')->rpcCall('shop.authorize.cat',array('shop_id'=>$this->shopId));
        return view::make('topshop/selector/sku/index.html', $pagedata);
    }

    public function formatSelectedSkusRow()
    {
        $skuIds = input::get('item_id');
        $textcol = input::get('textcol');
        $ac = input::get('ac');
        $extendView = input::get('view');
        $searchParams['fields'] = 'sku_id,item_id,title,image_default_id,price,brand_id,spec_info';
        $searchParams['sku_id'] = implode(',', $skuIds);
        $skusList = app::get('topshop')->rpcCall('sku.search', $searchParams);

        //特殊判断 有待优化
        if($skusList['total_found']>40)
        {
            $pages = ceil($skusList['total_found']/40);

            for($i=2;$i<=$pages;$i++)
            {
                $searchParams = array(
                    'page_no' => $i,
                    'item_id' => implode(',',$itemIds),
                    //'approve_status' => 'onsale',
                    'fields' => 'item_id,title,image_default_id,cat_id,brand_id,price,spec_info',
                );
                $skusListData = app::get('syspromotion')->rpcCall('item.search',$searchParams);
                $skusList['list'] = array_merge($skusList['list'],$skusListData['list']);
            }
        }
        foreach($skusList['list'] as $key=>$value)
        {
            if(input::get('pricemin'))
            {
                $skusList['list'][$key]['discount_min'] = input::get('pricemin');
            }
            if(input::get('pricemax'))
            {
                $skusList['list'][$key]['discount_max'] = input::get('pricemax');
            }
        }

        $extends = input::get('extends');
        $extendsData = input::get('extends_data');
        if( count($extends) > 0 )
        {
            $fmtItemExtendsData = [];
            foreach($extendsData as $item)
            {
                $skuId = $item['sku_id'];

                $fmtItemExtendsData[$skuId] = $item;
            }

            foreach($skusList['list'] as $key=>$value)
            {
                $skuId = $value['sku_id'];
                $skusList['list'][$key]['extendsData'] = $fmtItemExtendsData[$skuId];
            }

            $pagedata['_input']['extends'] = $extends;
        }

        $datavalues = input::get('values');
        if(count($datavalues) > 0)
        {
            $valuesData = [];
            foreach($datavalues as $item)
            {
                $skuId = $item['sku_id'];

                $valuesData[$skuId] = $item;
            }

            foreach($skusList['list'] as $key=>$value)
            {
                $skuId = $value['sku_id'];
                $skusList['list'][$key]['datavalue'] = $valuesData[$skuId];
            }
        }

        $pagedata['_input']['skusList'] = $skusList['list'];
        $pagedata['_input']['view'] = $extendView;
        if(!$textcol)
        {
            $pagedata['_input']['_textcol'] = 'title';
        }
        else
        {
            $pagedata['_input']['_textcol'] = explode(',',$textcol);
        }
        $pagedata['ac'] = $ac;
        $pagedata['image_default_id'] = kernel::single('image_data_image')->getImageSetting('item');
        return view::make('topshop/selector/sku/input-row.html', $pagedata);
    }

    public function phoneOrderformatSelectedSkusRow()
    {
        $skuIds = input::get('item_id');
        $textcol = input::get('textcol');
        $ac = input::get('ac');
        $extendView = input::get('view');
        $searchParams['fields'] = 'sku_id,item_id,title,image_default_id,price,brand_id,spec_info';
        $searchParams['sku_id'] = implode(',', $skuIds);
        $skusList = app::get('topshop')->rpcCall('sku.search', $searchParams);

        //特殊判断 有待优化
        if($skusList['total_found']>40)
        {
            $pages = ceil($skusList['total_found']/40);

            for($i=2;$i<=$pages;$i++)
            {
                $searchParams = array(
                    'page_no' => $i,
                    'item_id' => implode(',',$itemIds),
                    //'approve_status' => 'onsale',
                    'fields' => 'item_id,title,image_default_id,cat_id,brand_id,price,spec_info',
                );
                $skusListData = app::get('syspromotion')->rpcCall('item.search',$searchParams);
                $skusList['list'] = array_merge($skusList['list'],$skusListData['list']);
            }
        }
        foreach($skusList['list'] as $key=>$value)
        {
            if(input::get('pricemin'))
            {
                $skusList['list'][$key]['discount_min'] = input::get('pricemin');
            }
            if(input::get('pricemax'))
            {
                $skusList['list'][$key]['discount_max'] = input::get('pricemax');
            }
        }

        $extends = input::get('extends');
        $extendsData = input::get('extends_data');
        if( count($extends) > 0 )
        {
            $fmtItemExtendsData = [];
            foreach($extendsData as $item)
            {
                $skuId = $item['sku_id'];

                $fmtItemExtendsData[$skuId] = $item;
            }

            foreach($skusList['list'] as $key=>$value)
            {
                $skuId = $value['sku_id'];
                $skusList['list'][$key]['extendsData'] = $fmtItemExtendsData[$skuId];
            }

            $pagedata['_input']['extends'] = $extends;
        }

        $datavalues = input::get('values');
        if(count($datavalues) > 0)
        {
            $valuesData = [];
            foreach($datavalues as $item)
            {
                $skuId = $item['sku_id'];

                $valuesData[$skuId] = $item;
            }

            foreach($skusList['list'] as $key=>$value)
            {
                $skuId = $value['sku_id'];
                $skusList['list'][$key]['datavalue'] = $valuesData[$skuId];
            }
        }

        $pagedata['_input']['skusList'] = $skusList['list'];
        $pagedata['_input']['view'] = $extendView;
        if(!$textcol)
        {
            $pagedata['_input']['_textcol'] = 'title';
        }
        else
        {
            $pagedata['_input']['_textcol'] = explode(',',$textcol);
        }
        $pagedata['ac'] = $ac;
        $pagedata['image_default_id'] = kernel::single('image_data_image')->getImageSetting('item');
        return view::make('topshop/selector/sku/phoneorder-input-row.html', $pagedata);
    }

    //根据商家id和3级分类id获取商家所经营的所有品牌
    public function getBrandList()
    {
        $shopId = $this->shopId;
        $catId = input::get('catId');
        $params = array(
            'shop_id'=>$shopId,
            'cat_id'=>$catId,
            'fields'=>'brand_id,brand_name,brand_url'
        );
        $brands = app::get('topshop')->rpcCall('category.get.cat.rel.brand',$params);
        return response::json($brands);
    }

    //根据商家类目id的获取商家所经营类目下的所有商品
    public function searchSku($json=true)
    {
        $shopId = $this->shopId;
        $catId = input::get('catId');
        $brandId = input::get('brandId');
        $keywords = input::get('searchname');
        $brandName = input::get('searchbrand');
        $bn = input::get('searchbn');
        $pages = input::get('pages');

        $searchParams = array(
            'shop_id' => $shopId,
            'brand_id' => $brandId,
            'search_keywords' => $keywords,
            'bn' => $bn,
            'page_no' => intval($pages),
            'page_size' => intval($this->limit),
        );
        if($catId)
        {
            $searchParams['cat_id'] = app::get('topshop')->rpcCall('category.cat.get.leafCatId',array('cat_id'=>intval($catId)));
        }
        if(trim($brandName) && trim($brandName) != 'undefined')
        {
            $searchBrandParams = array('brand_name'=>trim($brandName),'fields'=>'brand_id');
            $brand = app::get('topshop')->rpcCall('category.brand.get.list', $searchBrandParams);
            if($brand)
            {
                $tmpBrandIds = array_column($brand, 'brand_id');
                $searchParams['brand_id'] = implode(',', $tmpBrandIds);
            }
            else
            {
                return view::make('topshop/selector/sku/list.html', $pagedata);
            }
        }

        $searchParams['fields'] = 'sku_id,item_id,title,image_default_id,price,brand_id,spec_info,status';
        $skusList = app::get('topshop')->rpcCall('sku.search', $searchParams);
        $pagedata['skusList'] = $skusList['list'];
        $pagedata['total'] = $skusList['total_found'];
        $totalPage = ceil($skusList['total_found']/$this->limit);
        $filter = input::get();
        $filter['pages'] = time();
        $pagers = array(
            'link' => url::action('topshop_ctl_selector_sku@searchSku', $filter),
            'current' => $pages,
            'use_app' => 'topshop',
            'total' => $totalPage,
            'token' => time(),
        );
        $pagedata['pagers'] = $pagers;

        $pagedata['image_default_id'] = kernel::single('image_data_image')->getImageSetting('item');

        return view::make('topshop/selector/sku/list.html', $pagedata);
    }
}
