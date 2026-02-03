package middleWare

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg":  "token不能为空",
				"code": 500,
			})
			c.Abort()
			return
		}
		getToken, err := GetToken(token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg":  "jwt解密失败",
				"code": 500,
			})
			c.Abort()
			return
		}
		//res := config.Rdb.SIsMember(context.Background(), "list:k1", getToken["userId"]).Val()
		//if res {
		//	c.JSON(http.StatusBadRequest, gin.H{
		//		"msg":  "用户黑名单",
		//		"code": 500,
		//	})
		//	c.Abort()
		//	return
		//}
		c.Set("userId", getToken["userId"].(string))
		c.Next()
	}
}

//handler, err := middleWare.TokenHandler(strconv.FormatInt(login.UserId, 10), time.Now().Add(time.Hour*time.Duration(1)).Unix())
/*
var good model.Goods
	err := good.FindGoodsByTitle(config.DB, in.Title)
	if err == nil {
		return nil, errors.New("数据不能重入库")
	}
	goods := model.Goods{
		Title:   in.Title,
		Url:     in.Title,
		Content: in.Content,
		Stock:   int(in.Stock),
		Price:   float64(in.Price),
	}
	begin := config.DB.Begin()
	err = goods.GoodsAdd(config.DB)
	if err != nil {
		begin.Rollback()
		return nil, errors.New("商品上架失败")
	}

	key := fmt.Sprintf("goods:%s", in.Title)
	marshaler, err := json.Marshal(goods)
	if err != nil {
		begin.Rollback()
		return nil, errors.New("商品序列化失败")
	}
	nx := config.Rdb.SetNX(config.Ctx, key, marshaler, time.Minute*30).Val()
	if !nx {
		begin.Rollback()
		return nil, errors.New("商品已存在,请勿重复上架")
	}
	begin.Commit()
	goodsMap := map[string]interface{}{
		"id":      goods.ID,
		"title":   in.Title,
		"url":     in.Title,
		"content": in.Content,
		"stock":   int(in.Stock),
		"price":   float64(in.Price),
	}
	_, err = config.Es.Index().Index("goods_map").BodyJson(goodsMap).Do(context.Background())
	if err != nil {
		return nil, errors.New("es同步失败")
	}
*/
/*
var detail model.Goods
	key := fmt.Sprintf("goods:%s", in.Title)
	err := config.Rdb.Exists(config.Ctx, key).Val()
	if err == 0 {
		return nil, errors.New("商品不存在")
	} else {
		result, _ := config.Rdb.Get(config.Ctx, key).Result()
		json.Unmarshal([]byte(result), &detail)
	}
*/
/*
 proto, _ := cmd.Flags().GetString("proto")
 exec.Command("protoc",proto,"--go_out=.","--go-grpc_out=.",).Run()
 pbCmd.Flags().StringVarP(&proto, "proto", "p", "", "生成proto文件")
*/
