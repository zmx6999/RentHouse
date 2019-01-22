package models

import (
	"time"
	"190120/utils"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id int `json:"user_id"`
	Name string `json:"name" orm:"size(32)"`
	Password string `json:"password" orm:"size(128)"`
	Mobile string `json:"mobile" orm:"size(11);unique"`
	Real_name string `json:"real_name" orm:"size(32);default()"`
	Id_card string `json:"id_card" orm:"size(18);default()"`
	Avatar_url string `json:"avatar_url" orm:"size(256);default()"`

	Houses []*House `orm:"reverse(many)"`
	Orders []*Order `orm:"reverse(many)"`
}

type House struct {
	Id int `json:"house_id"`
	Title string `json:"title" orm:"size(64)"`
	Price int `json:"price" orm:"default(0)"`
	Address string `json:"address" orm:"size(512);default()"`
	Room_count int `json:"room_count" orm:"default(1)"`
	Acreage int `json:"acreage" orm:"default(0)"`
	Unit string `json:"unit" orm:"size(32);default()"`
	Capacity int `json:"capacity" orm:"default(1)"`
	Beds string `json:"beds" orm:"size(64);default()"`
	Deposit int `json:"deposit" orm:"default(0)"`
	Min_days int `json:"min_days" orm:"default(1)"`
	Max_days int `json:"max_days" orm:"default(0)"`
	Order_count int `json:"order_count" orm:"default(0)"`
	Index_image_url string `json:"index_image_url" orm:"size(256);default()"`
	Create_time time.Time `json:"create_time" orm:"type(datetime)"`

	Area *Area `orm:"rel(fk)"`
	User *User `orm:"rel(fk)"`

	Facilities []*Facility `orm:"reverse(many)"`
	HouseImages []*HouseImage `orm:"reverse(many)"`
	Orders []*Order `orm:"reverse(many)"`
}

func (this *House) Info() map[string]interface{} {
	loc,_:=time.LoadLocation(utils.TimeZone)
	data:=map[string]interface{}{
		"house_id": this.Id,
		"title": this.Title,
		"price": this.Price,
		"area": this.Area.Name,
		"index_image_url": utils.AddDomain2Url(this.Index_image_url),
		"room_count": this.Room_count,
		"order_count": this.Order_count,
		"address": this.Address,
		"user_avatar": utils.AddDomain2Url(this.User.Avatar_url),
		"create_time": this.Create_time.In(loc).Format("2006-01-02 15:04:05"),
	}
	return data
}

func (this *House) Desc() map[string]interface{} {
	data:=map[string]interface{}{
		"house_id": this.Id,
		"user_id": this.User.Id,
		"username": this.User.Name,
		"user_avatar": utils.AddDomain2Url(this.User.Avatar_url),
		"title": this.Title,
		"price": this.Price,
		"address": this.Address,
		"room_count": this.Room_count,
		"acreage": this.Acreage,
		"uint": this.Unit,
		"capacity": this.Capacity,
		"beds": this.Beds,
		"deposit": this.Deposit,
		"min_days": this.Min_days,
		"max_days": this.Max_days,
	}

	imgUrls:=[]string{}
	for _,img:=range this.HouseImages{
		imgUrls=append(imgUrls,utils.AddDomain2Url(img.Url))
	}
	data["img_urls"]=imgUrls

	facilities:=[]string{}
	for _,facility:=range this.Facilities{
		facilities=append(facilities,facility.Name)
	}
	data["facilities"]=facilities

	comments:=[]map[string]interface{}{}
	loc,_:=time.LoadLocation(utils.TimeZone)
	for _,order:=range this.Orders{
		comment:=map[string]interface{}{
			"username": order.User.Name,
			"comment": order.Comment,
			"create_time": order.Create_time.In(loc).Format("2006-01-02 15:04:05"),
		}
		comments=append(comments,comment)
	}
	data["comments"]=comments
	return data
}

type Area struct {
	Id int `json:"area_id"`
	Name string `json:"name" orm:"size(32)"`

	Houses []*House `orm:"reverse(many)"`
}

type Facility struct {
	Id int `json:"facility_id"`
	Name string `json:"name" orm:"size(32)"`

	Houses []*House `orm:"rel(m2m)"`
}

type HouseImage struct {
	Id int `json:"house_image_id"`
	Url string `json:"url" orm:"size(256)"`

	House *House `orm:"rel(fk)"`
}

type Order struct {
	Id int `json:"order_id"`
	Begin_date time.Time `json:"begin_date" orm:"type(datetime)"`
	End_date time.Time `json:"end_date" orm:"type(datetime)"`
	Days int `json:"days"`
	House_price int `json:"house_price"`
	Amount int `json:"amount"`
	Status string `json:"status" orm:"default(WAIT_ACCEPT)"`
	Comment string `json:"comment" orm:"size(512);default()"`
	Create_time time.Time `json:"create_time" orm:"type(datetime)"`
	Credit bool `json:"credit"`

	User *User `orm:"rel(fk)"`
	House *House `orm:"rel(fk)"`
}

func (this *Order) Info() map[string]interface{} {
	loc,_:=time.LoadLocation(utils.TimeZone)
	data:=map[string]interface{}{
		"order_id": this.Id,
		"title": this.House.Title,
		"image_url": utils.AddDomain2Url(this.House.Index_image_url),
		"begin_date": this.Begin_date.In(loc).Format("2006-01-02"),
		"end_date": this.End_date.In(loc).Format("2006-01-02"),
		"create_time": this.Create_time.In(loc).Format("2006-01-02 15:04:05"),
		"days": this.Days,
		"amount": this.Amount,
		"status": this.Status,
		"comment": this.Comment,
		"credit": this.Credit,
	}
	return data
}

func init()  {
	orm.RegisterDataBase("default","mysql",utils.MySQLUserName+":"+utils.MySQLPassword+"@tcp("+utils.MySQLHost+":"+utils.MySQLPort+")/"+utils.MySQLDBName+"?charset=utf8")
	orm.RegisterModel(new(User),new(House),new(Area),new(Facility),new(HouseImage),new(Order))
	orm.RunSyncdb("default",false,true)
}
