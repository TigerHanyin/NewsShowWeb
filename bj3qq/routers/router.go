package routers

import (
	"bj3qq/controllers"
	"github.com/astaxie/beego"
    "github.com/astaxie/beego/context"
)

func init() {
    beego.InsertFilter("/article/*",beego.BeforeExec,filterFunc)
    beego.Router("/", &controllers.MainController{})
    beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandleRegister")
    //登录业务处理
    beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
    //首页展示
    beego.Router("/article/index",&controllers.ArticleController{},"get:ShowIndex")
    //添加文章业务
    beego.Router("/article/addArticle",&controllers.ArticleController{},
    "get:ShowAddArticle;post:HandleAddArticle")
    //查看文章详情
    beego.Router("/article/content",&controllers.ArticleController{},"get:ShowContent")
    //编辑文章
    beego.Router("/article/update",&controllers.ArticleController{},"get:ShowUpdate;post:HandleUpdate")
    //删除文章
    beego.Router("/article/delete",&controllers.ArticleController{},"get:HandleDelete")
    //展示添加分类页面
    beego.Router("/article/addType",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
    beego.Router("/article/addType",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
    beego.Router("/article/logout",&controllers.UserController{},"get:Handlelogout")
}
func filterFunc(ctx *context.Context) {
    userName := ctx.Input.Session("userName")
    if userName ==nil{
        ctx.Redirect(302,"/login")
        return
    }
}