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
	// beego.Router("/wsViewer", &controllers.MainController{}, "GET:ServerWSViewer")
	beego.Router("/addressEdit", &controllers.MainController{}, "GET:AddressEditIndex")
	beego.Router("/rankIndex", &controllers.MainController{}, "GET:RankIndex")
	beego.Router("/newGameIndex", &controllers.MainController{}, "GET:NewGameIndex")
	beego.Router("/gameListIndex", &controllers.MainController{}, "GET:GameListIndex")

	beego.Router("/distributors", &controllers.MainController{}, "GET:Distributors")
	beego.Router("/gameList", &controllers.MainController{}, "GET:GameList")
	beego.Router("/mapNameList", &controllers.MainController{}, "GET:MapNameList")
	// beego.Router("/orders", &controllers.MainController{}, "GET:Orders")
	beego.Router("/uploadMapData", &controllers.MainController{}, "POST:UploadMapData")
	beego.Router("/mapData", &controllers.MainController{}, "GET:MapData")
	beego.Router("/restartGame", &controllers.MainController{}, "GET:RestartGame")
	beego.Router("/newGame", &controllers.MainController{}, "post:NewGame")
}
