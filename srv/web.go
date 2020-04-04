package srv

import (
	"context"
	"fmt"
	"shop/models"
	"strconv"
)

type Hrefs struct{}

func (s *Hrefs) GetOrderByUser(ctx context.Context, req int, rsp []*models.Order) error {
	fmt.Println("微服务请求参数 req=", req)
	rsp = models.OrderFindByUser(strconv.Itoa(req))
	fmt.Println("微服务响应： rsp=", rsp)
	return nil
}

func (s *Hrefs) GetOrderById(ctx context.Context, req int, rsp *models.Order) error {
	fmt.Println("微服务请求参数 req=", req)
	*rsp = *models.OrderFindById(req)
	fmt.Println("微服务响应： rsp.name=", rsp.Username)
	return nil
}

func (s *Hrefs) IndexLinks(ctx context.Context, req int, rsp *models.User) error {
	//items, err := domain.IndexLinks()
	//rsp.TableName()
	//if err != nil {
	//	return err
	//}

	rsp.Name = "aa"
	fmt.Println("1111111   micro为服务响应, 接收到的请求参数为：",req)

	return nil
}

//func (s *Hrefs) TopArticles(ctx context.Context, req bool, rsp *model.ArticleItems) error {
//	items, err := domain.TopArticles()
//	if err != nil {
//		return err
//	}
//
//	rsp.Items = items
//
//	return nil
//}
//
//func (s *Hrefs) TopCusLinks(ctx context.Context, req bool, rsp *model.CusLinkItems) error {
//	items, err := domain.TopCusLinks()
//	if err != nil {
//		return err
//	}
//
//	rsp.Items = items
//
//	return nil
//}
//
//func (s *Hrefs) ListLinks(ctx context.Context, req string, rsp *model.LinkItems) error {
//	items, err := domain.ListLinks(req)
//	if err != nil {
//		return err
//	}
//
//	rsp.Items = items
//
//	return nil
//}
//
//func (s *Hrefs) ListCusLinksByCatId(ctx context.Context, req string, rsp *model.CusLinkItems) error {
//	items, err := domain.ListCusLinksByCatId(req)
//	if err != nil {
//		return err
//	}
//
//	rsp.Items = items
//
//	return nil
//}
//
//func (s *Hrefs) ListCusLinks(ctx context.Context, req bool, rsp *model.CusLinkItems) error {
//	items, err := domain.ListCusLinks()
//	if err != nil {
//		return err
//	}
//
//	rsp.Items = items
//
//	return nil
//}
//
//func (s *Hrefs) ListArticles(ctx context.Context, req bool, rsp *model.ArticleItems) error {
//	items, err := domain.ListArticles()
//	if err != nil {
//		return err
//	}
//
//	rsp.Items = items
//
//	return nil
//}
//
//func (s *Hrefs) GetArticle(ctx context.Context, req string, rsp *model.Article) error {
//	result, err := domain.GetArticle(req)
//	*rsp = *result
//	return err
//}
//
//func (s *Hrefs) UpdateArticleVisited(ctx context.Context, req string, rsp *bool) error {
//	err := domain.UpdateArticleVisited(req)
//	return err
//}
//
//func (s *Hrefs) LinkVisitedCount(ctx context.Context, req bool, rsp *int64) error {
//	result, err := domain.LinkVisitedCount()
//	*rsp = result
//	return err
//}
//
//func (s *Hrefs) GetLinkUrl(ctx context.Context, req string, rsp *string) error {
//	url, err := domain.GetLinkUrl(req)
//	*rsp = url
//	return err
//}
//
//func (s *Hrefs) UpdateLinkVisited(ctx context.Context, req string, rsp *bool) error {
//	err := domain.UpdateLinkVisited(req)
//	return err
//}
//
//func (s *Hrefs) GetCusLinkUrl(ctx context.Context, req int, rsp *string) error {
//	url, err := domain.GetCusLinkUrl(req)
//	*rsp = url
//	return err
//}
//
//func (s *Hrefs) UpdateCusLinkVisited(ctx context.Context, req int, rsp *bool) error {
//	err := domain.UpdateCusLinkVisited(req)
//	return err
//}
