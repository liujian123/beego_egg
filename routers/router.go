// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"beego_egg/controllers"
	_ "beego_egg/models"
	"beego_egg/utils"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"strings"
)

func init() {
	beego.Router("/login", &controllers.UserController{}, "post:Login")
	beego.Router("/user/info", &controllers.UserController{}, "get:Info")
	beego.Router("/user/collect", &controllers.UserController{}, "post:Collect")
	beego.Router("/user/withdraw", &controllers.UserController{}, "post:Withdraw")
	beego.Router("/user/balance", &controllers.UserController{}, "get:Balance")
	beego.Router("/connect", &controllers.WsSocketController{}, "*:HandleConnect")
	beego.InsertFilter("/user/*", beego.BeforeRouter, FilterAuthorization)
}

var FilterAuthorization = func(ctx *context.Context) {
	if ctx.Input.Header("Authorization") == "" && ctx.Request.RequestURI != "/login" {
		ctx.ResponseWriter.WriteHeader(401)
		ctx.ResponseWriter.Write([]byte("no permission"))
		return
	}
	if ctx.Request.RequestURI != "/login" && ctx.Input.Header("Authorization") != "" {
		et := utils.EasyToken{}
		authtoken := strings.TrimSpace(ctx.Input.Header("Authorization"))
		valido, uid, err := et.ValidateToken(authtoken)
		if !valido || uid == 0 || err != nil {
			ctx.ResponseWriter.WriteHeader(401)
			ctx.ResponseWriter.Write([]byte(fmt.Sprintf("%s", err)))
			return
		}
		ctx.Input.SetData("uid", uid)
	}
	return
}
