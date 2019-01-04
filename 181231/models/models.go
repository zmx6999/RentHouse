package models

import (
	_ "github.com/go-sql-driver/mysql"
	"time"
	"sss/181231/utils"
		"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"fmt"
)

type User struct {
	Id int `json:"user_id"`
	Name string `orm:"size(32)" json:"name"`
	Password string `orm:"size(128)" json:"password"`
	Mobile string `orm:"size(11);unique" json:"mobile"`
	Real_name string `orm:"size(32);" json:"real_name"`
	Id_card string `orm:"size(20)" json:"id_card"`
	Avatar_url string `orm:"size(256)" json:"avatar_url"`

	Houses []*House `orm:"reverse(many)" json:"houses"`
	Orders []*Order `orm:"reverse(many)" json:"orders"`
}

type House struct {
	Id int `json:"house_id"`
	Title string `orm:"size(64)" json:"title"`
	Price int `orm:"default(0)" json:"price"`
	Address string `orm:"size(512);default('')" json:"address"`
	Room_count int `orm:"default(1)" json:"room_count"`
	Acreage int `orm:"default(0)" json:"acreage"`
	Unit string `orm:"size(32);default('')" json:"unit"`
	Capacity int `orm:"default(1)" json:"capacity"`
	Beds string `orm:"size(64);default('')" json:"beds"`
	Deposit int `orm:"default(0)" json:"deposit"`
	Min_days int `orm:"default(1)" json:"min_days"`
	Max_days int `orm:"default(0)" json:"max_days"`
	Order_count int `orm:"default(0)" json:"order_count"`
	Index_image_url string `orm:"size(256);default('')" json:"index_image_url"`
	Create_time time.Time `orm:"type(datetime)" json:"create_time"`

	User *User `orm:"rel(fk)" json:"user_id"`
	Area *Area `orm:"rel(fk)" json:"area_id"`

	Facilities []*Facility `orm:"reverse(many)" json:"facilities"`
	HouseImages []*HouseImage `orm:"reverse(many)" json:"house_images"`
	Orders []*Order `orm:"reverse(many)" json:"orders"`
}

func (this *House) Info() map[string]interface{} {
	loc,_:=time.LoadLocation("Asia/Saigon")
	info:=map[string]interface{} {
		"house_id": this.Id,
		"title": this.Title,
		"price": this.Price,
		"area_name": this.Area.Name,
		"index_image_url": utils.AddDomain2Url(this.Index_image_url),
		"room_count": this.Room_count,
		"order_count": this.Order_count,
		"address": this.Address,
		"user_avatar": utils.AddDomain2Url(this.User.Avatar_url),
		"create_time": this.Create_time.In(loc).Format("2006-01-02 15:04:05"),
	}
	return info
}

func (this *House) Desc() map[string]interface{} {
	desc:=map[string]interface{} {
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
	desc["img_urls"]=imgUrls

	facilities:=[]int{}
	for _,facility:=range this.Facilities{
		facilities=append(facilities,facility.Id)
	}
	desc["facilities"]=facilities

	comments:=[]interface{}{}
	orders:=[]Order{}
	o:=orm.NewOrm()
	_,err:=o.QueryTable("Order").Filter("House__Id",this.Id).Filter("Status",ORDER_STATUS_COMPLETE).OrderBy("-Create_time").Limit(10).All(&orders)
	if err!=nil {
		beego.Error(err)
	}
	n:=len(orders)
	for i:=0; i<n; i++ {
		o.LoadRelated(&orders[i],"User")
		username:=orders[i].User.Name
		if username=="" {
			username="匿名用户"
		}
		comment:=map[string]interface{} {
			"username": username,
			"comment": orders[i].Comment,
			"create_time":orders[i].Create_time.Format("2006-01-02 15:04:05"),
		}
		comments=append(comments,comment)
	}
	desc["comments"]=comments
	return desc
}

type Area struct {
	Id int `json:"area_id"`
	Name string `orm:"size(32)" json:"name"`
	Houses []*House `orm:"reverse(many)" json:"houses"`
}

type Facility struct {
	Id int `json:"facility_id"`
	Name string `orm:"size(32)" json:"name"`
	Houses []*House `orm:"rel(m2m)"`
}

type HouseImage struct {
	Id int `json:"house_image_id"`
	Url string `orm:"size(256)" json:"url"`
	House *House `orm:"rel(fk)" json:"house_id"`
}

const (
	ORDER_STATUS_WAIT_ACCEPT= "WAIT_ACCEPT"
	ORDER_STATUS_WAIT_PAYMENT= "WAIT_PAYMENT"
	ORDER_STATUS_PAID= "PAID"
	ORDER_STATUS_WAIT_COMMENT= "WAIT_COMMENT"
	ORDER_STATUS_COMPLETE= "COMPLETE"
	ORDER_STATUS_CANCELED= "CONCELED"
	ORDER_STATUS_REJECTED= "REJECTED"
)

type Order struct {
	Id int `json:"order_id"`
	User *User `orm:"rel(fk)" json:"user_id"`
	House *House `orm:"rel(fk)" json:"house_id"`
	Begin_date time.Time `orm:"type(datetime)" json:"begin_date"`
	End_date time.Time `orm:"type(datetime)" json:"end_date"`
	Days int `json:"days"`
	House_price int `json:"house_price"`
	Amount int `json:"amount"`
	Status string `orm:"default(WAIT_ACCEPT)" json:"status"`
	Comment string `orm:"size(512)" json:"comment"`
	Create_time time.Time `orm:"type(datetime)" json:"create_time"`
	Credit bool `json:"credit"`
}

func (this *Order) Info() map[string]interface{} {
	info:=map[string]interface{} {
		"order_id": this.Id,
		"title": this.House.Title,
		"image_url": utils.AddDomain2Url(this.House.Index_image_url),
		"begin_date": this.Begin_date.Format("2006-01-02"),
		"end_date": this.End_date.Format("2006-01-02"),
		"create_time": this.Create_time.Format("2006-01-02 15:04:05"),
		"days": this.Days,
		"amount": this.Amount,
		"status": this.Status,
		"comment": this.Comment,
		"credit": this.Credit,
	}
	return info
}

func init()  {
	orm.RegisterDriver("mysql",orm.DRMySQL)
	fmt.Print("root:123456@tcp("+utils.MySQLHost+":"+utils.MySQLPort+")/"+utils.MySQLDBName+"?charset=utf8")
	orm.RegisterDataBase("default","mysql","root:123456@tcp("+utils.MySQLHost+":"+utils.MySQLPort+")/"+utils.MySQLDBName+"?charset=utf8",30)
	orm.RegisterModel(new(User),new(House),new(Area),new(Facility),new(HouseImage),new(Order))
	orm.RunSyncdb("default",false,true)
}
