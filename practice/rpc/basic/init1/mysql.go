package init1

import (
	"fmt"
	"zhongyao/aa/crawler/practice/rpc/basic/config"
	"zhongyao/aa/crawler/practice/rpc/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMysql() {
	var err error
	mysqlConfig := config.Config.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConfig.User,
		mysqlConfig.Password,
		mysqlConfig.Host,
		mysqlConfig.Port,
		mysqlConfig.Database,
	)
	config.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println("mysql连接成功")
	err = config.DB.AutoMigrate(&model.User{}, &model.Order{}, &model.Goods{}, &model.OrderItem{})
	if err != nil {
		panic(err)
	}
	fmt.Println("迁移成功")
}
