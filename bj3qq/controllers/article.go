package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"bj3qq/models"
	"math"
	"strconv"
)

type ArticleController struct {
	beego.Controller
}

//展示首页
func(this*ArticleController)ShowIndex(){
	//获取所有文章数据，展示到页面
	o := orm.NewOrm()
	qs := o.QueryTable("Article")
	var articles []models.Article
	//qs.All(&articles)
	typeName:=this.GetString("select")
	var count int64
	if typeName==""{
		count,_=qs.RelatedSel("ArticleType").Count()
	}else{
		count,_=qs.RelatedSel("ArticleType").Filter("ArticleType__typeName").Count()
	}

	//获取总页数
	pageIndex := 2


	pageCount := math.Ceil(float64(count) / float64(pageIndex))
	//获取首页和末页数据
	//获取页码
	pageNum ,err := this.GetInt("pageNum")
	if err != nil {
		pageNum = 1
	}
	beego.Info("数据总页数未:",pageNum)
	if typeName ==""{
		qs.Limit(pageIndex,pageIndex*(pageNum-1)).RelatedSel("ArticleType").All(&articles)
	}else {qs.Limit(pageIndex,(pageIndex*pageNum-1)).RelatedSel("ArticleType").Filter("ArticleType__typeName").All(articles)}

	//获取对应页的数据   获取几条数据     起始位置
	qs.Limit(pageIndex,pageIndex * (pageNum - 1)).All(&articles)

	this.Data["articles"] = articles
	this.Data["count"] = count
	this.Data["pageCount"] = pageCount
	this.Data["pageNum"] = pageNum
	this.Layout="layout.html"
	this.TplName = "index.html"
}

//展示添加文章页面
func(this*ArticleController)ShowAddArticle(){
	//获取所有类型并绑定下拉框
	o := orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)

	this.Data["articleTypes"] = articleTypes
	this.Layout="layout.html"
	this.TplName = "add.html"
}

//处理添加文章业务
func(this*ArticleController)HandleAddArticle(){
	//获取数据
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	typeName := this.GetString("select")

	//校验数据
	if articleName == "" || content == "" || typeName == ""{
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "获取数据错误"
		this.Layout="layout.html"
		this.TplName = "add.html"
		return
	}

	//获取图片
	//返回值 文件二进制流  文件头    错误信息
	file,head,err := this.GetFile("uploadname")
	if err != nil {
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "图片上传失败"
		this.Layout="layout.html"
		this.TplName = "add.html"
		return
	}
	defer file.Close()
	//校验文件大小
	if head.Size >5000000{
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "图片数据过大"
		this.Layout="layout.html"
		this.TplName = "add.html"
		return
	}

	//校验格式 获取文件后缀
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "上传文件格式错误"
		this.Layout="layout.html"
		this.TplName = "add.html"
		return
	}

	//防止重名
	fileName := time.Now().Format("200601021504052222")


	//jianhuangcaozuo

	//把上传的文件存储到项目文件夹
	this.SaveToFile("uploadname","./static/img/"+fileName+ext)

	//处理数据
	//把数据存储到数据库
	//获取orm对象
	o := orm.NewOrm()
	//获取插入独享
	var article models.Article
	//给插入对象赋值
	article.Title = articleName
	article.Content = content
	article.Img = "/static/img/"+fileName+ext

	//获取一个类型对象，并插入到文章中
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Read(&articleType,"TypeName")

	article.ArticleType = &articleType
	//插入数据
	_,err = o.Insert(&article)
	if err != nil {
		beego.Error("获取数据错误",err)
		this.Data["errmsg"] = "数据插入失败"
		this.Layout="layout.html"
		this.TplName = "add.html"
		return
	}

	//返回数据  跳转页面
	this.Redirect("/article/index",302)
}

//查看文章详情
func(this*ArticleController)ShowContent(){
	//获取数据
	id,err := this.GetInt("id")
	//校验数据
	if err != nil {
		beego.Error("获取文章id错误")
		this.Redirect("/article/index",302)  //渲染  如果页面本身有数据加载，不能直接渲染
		return
	}
	//处理数据
	//查询文章数据
	o := orm.NewOrm()
	//获取查询对象
	var article models.Article
	//给查询条件赋值
	article.Id = id
	//查询
	o.Read(&article)

	//给更新条件赋值
	article.ReadCount += 1
	o.Update(&article)

	//返回数据
	this.Data["article"] = article
	userName:=this.GetSession("userNmae")
	var user models.User
	user.Name=userName.(string)
	o.Read(&user,"Name")
	m2m:=o.QueryM2M(&article,"Users")
	m2m.Add(user)

	this.Layout="layout.html"
	this.TplName = "content.html"
}

//展示文章编辑页面
func(this*ArticleController)ShowUpdate(){
	//获取数据
	id,err := this.GetInt("id")
	//校验数据
	if err != nil {
		beego.Error("获取文章ID错误")
		this.Redirect("/article/index",302)
		return
	}
	//处理数据
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Read(&article)


	//返回数据
	this.Data["article"] = article
	this.Layout="layout.html"
	this.TplName = "update.html"
}

//封装上传文件处理函数
func UploadFile(this *ArticleController,filePath string,errHtml string)string{
	//获取图片
	//返回值 文件二进制流  文件头    错误信息
	file,head,err := this.GetFile(filePath)
	if err != nil {
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "图片上传失败"
		this.TplName = errHtml
		return ""
	}
	defer file.Close()
	//校验文件大小
	if head.Size >5000000{
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "图片数据过大"
		this.TplName = errHtml
		return ""
	}

	//校验格式 获取文件后缀
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "上传文件格式错误"
		this.TplName = errHtml
		return ""
	}

	//防止重名
	fileName := time.Now().Format("200601021504052222")


	//jianhuangcaozuo

	//把上传的文件存储到项目文件夹
	this.SaveToFile(filePath,"./static/img/"+fileName+ext)
	return "/static/img/"+fileName+ext

}

//处理文章编辑
func(this*ArticleController)HandleUpdate(){
	//获取数据
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	savePath := UploadFile(this,"uploadname","update.html")
	id,_ := this.GetInt("id")  //隐藏域传值
	//校验数据
	if articleName == "" || content == "" ||savePath == "" {
		beego.Error("获取数据失败")
		this.Redirect("/article/update?id="+strconv.Itoa(id),302)
		return
	}
	//处理数据
	//更新操作
	o := orm.NewOrm()
	var article models.Article
	//先查询要更新的文章是否存在
	article.Id = id
	//必须查询
	o.Read(&article)
	//更新   需要先赋新值   beego中的ORM如果需要更新，更新的对象Id必须有值
	article.Title = articleName
	article.Content = content
	article.Img = savePath
	o.Update(&article)


	//返回数据
	this.Redirect("/article/index",302)
}


//删除文章
func(this*ArticleController)HandleDelete(){
	//获取数据
	id,err := this.GetInt("id")
	//校验数据
	if err != nil {
		beego.Error("获取Id错误")
		this.Redirect("/article/index",302)
		return
	}
	//处理数据
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Delete(&article,"Id")

	//返回数据
	this.Redirect("/article/index",302)
}

//展示添加分类页面
func(this*ArticleController)ShowAddType(){
	//获取所有类型，并展示到页面上
	//获取所有用all
	o := orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)

	//返回数据
	this.Data["articleTypes"] = articleTypes
	this.Layout="layout.html"

	this.TplName = "addType.html"
}

//处理添加类型请求
func(this*ArticleController)HandleAddType(){
	//获取数据
	typeName := this.GetString("typeName")
	//校验数据
	if typeName == ""{
		beego.Error("类型名称传输失败")
		this.Redirect("/article/addType",302)
		return
	}
	//处理数据
	//插入操作
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Insert(&articleType)

	//返回数据
	this.Redirect("/article/addType",302)
}