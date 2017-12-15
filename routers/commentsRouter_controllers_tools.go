package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["vil_tools/controllers/tools:ToolControllers"] = append(beego.GlobalControllerRouter["vil_tools/controllers/tools:ToolControllers"],
		beego.ControllerComments{
			Method: "GetCourrntTime",
			Router: `/getCourrntTime`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
