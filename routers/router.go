package routers

import (
	"distributionGame/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "GET:Login")
	beego.Router("/index", &controllers.MainController{}, "GET:Index")
	beego.Router("/wsOrderDistribution", &controllers.MainController{}, "GET:ServerWSOrderDistribution")
	beego.Router("/orderDistribute", &controllers.MainController{}, "GET:OrderDistributeIndex")
	beego.Router("/distribution", &controllers.MainController{}, "GET:DistributionIndex")
	beego.Router("/userListIndex", &controllers.MainController{}, "GET:UserListIndex")
	beego.Router("/viewer", &controllers.MainController{}, "GET:ViewerIndex")
	beego.Router("/wsViewer", &controllers.MainController{}, "GET:ServerWSViewer")
	beego.Router("/addressEdit", &controllers.MainController{}, "GET:AddressEditIndex")

	beego.Router("/distributors", &controllers.MainController{}, "GET:Distributors")
	beego.Router("/orders", &controllers.MainController{}, "GET:Orders")
	beego.Router("/uploadMapData", &controllers.MainController{}, "POST:UploadMapData")
	beego.Router("/mapData", &controllers.MainController{}, "GET:MapData")
}
