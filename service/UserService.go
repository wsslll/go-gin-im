package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUserList 获取用户列表
// @Summary 用户列表
// @Description 获取用户列表接口
// @Tags 用户模块
// @Success 200 string json{"code","data"}
// @Router /user/list [get]
func GetUserList(ctx *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}

// CreateUser 新增用户接口
// @Summary 新增用户
// @Description 新增用户接口
// @param username formData string false "用户名"
// @param password formData string false "用户密码"
// @param repassword formData string false "确认密码"
// @Param phone formData string false "电话号码"
// @Param email formData string false "邮箱"
// @Tags 用户模块
// @Success 200 string json{"code","message"}
// @Router /user/create [post]
func CreateUser(ctx *gin.Context) {
	username, _ := ctx.GetPostForm("username")
	password, _ := ctx.GetPostForm("password")
	repassword, _ := ctx.GetPostForm("repassword")
	user := models.UserBasic{}
	if username == "" || password == "" || repassword == "" {
		ctx.JSON(-1, gin.H{
			"code":    -1,
			"message": "参数不能为空",
		})
		return
	}
	if password != repassword {
		ctx.JSON(-1, gin.H{
			"code":    -1,
			"message": "两次密码不一致",
		})
		return
	}
	user.Name = username
	salt := fmt.Sprintf("%06d", rand.Int31())
	user.Salt = salt
	user.PassWord = utils.MakePassword(password, salt)
	user.Phone, _ = ctx.GetPostForm("phone")
	user.Email, _ = ctx.GetPostForm("email")
	name := models.GetUserByName(user.Name)
	if name.Name != "" {
		ctx.JSON(-1, gin.H{
			"code":    -1,
			"message": "该用户名已注册",
		})
		return
	}
	phone := models.GetUserByPhone(user.Phone)
	if phone.Phone != "" {
		ctx.JSON(-1, gin.H{
			"code":    -1,
			"message": "该手机号已注册",
		})
		return
	}
	email := models.GetUserByEmail(user.Email)
	if email.Email != "" {
		ctx.JSON(-1, gin.H{
			"code":    -1,
			"message": "该邮箱已注册",
		})
		return
	}
	_, err := govalidator.ValidateStruct(user)
	if err != nil {

		ctx.JSON(-1, gin.H{
			"code":    -1,
			"message": "数据格式不匹配",
		})
		return
	}
	models.CreateUser(user)
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "添加用户成功",
	})
}

// UpdateUser 修改用户接口
// @Summary 修改用户
// @Description 修改用户接口
// @Tags 用户模块
// @Param id formData uint true "用户ID"
// @Param username formData string false "用户名称"
// @Param password formData string false "用户密码"
// @Param phone formData string false "电话号码"
// @Param email formData string false "邮箱"
// @Success 200 string json{"code","message"}
// @Router /user/update [put]
func UpdateUser(ctx *gin.Context) {
	user := models.UserBasic{}
	id, _ := ctx.GetPostForm("id")
	intID, _ := strconv.Atoi(id)
	user.ID = uint(intID)
	user.Name, _ = ctx.GetPostForm("username")
	user.PassWord, _ = ctx.GetPostForm("password")
	user.Phone, _ = ctx.GetPostForm("phone")
	user.Email, _ = ctx.GetPostForm("email")
	_, err := govalidator.ValidateStruct(user)
	if err != nil {

		ctx.JSON(-1, gin.H{
			"code":    -1,
			"message": "数据格式不匹配",
		})
		return
	}
	result := models.UpdateUser(user)
	if result.Error != nil {
		ctx.JSON(-1, gin.H{
			"code":    -1,
			"message": "修改用户失败",
		})
		panic(result.Error)
		return
	}
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "修改用户成功",
	})
}

// DeleteUser 删除用户接口
// @Summary 删除用户
// @Description 删除用户接口
// @Tags 用户模块
// @Success 200 string json{"code","message"}
// @Router /user/delete/:id [delete]
func DeleteUser(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	user := models.UserBasic{}
	user.ID = uint(id)
	result := models.DeleteUser(user)
	if result.Error != nil {
		ctx.JSON(-1, gin.H{
			"code":    -1,
			"message": "删除用户失败",
		})
		panic(result.Error)
		return
	}
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "删除用户成功",
	})
}

// Login 用户登录接口
// @Summary 用户登录
// @Description 用户登录接口
// @param username formData string ture "用户名"
// @param password formData string ture "用户密码"
// @Tags 用户模块
// @Success 200 string json{"code","message","data"}
// @Router /user/login [post]
func Login(ctx *gin.Context) {
	username, _ := ctx.GetPostForm("username")
	password, _ := ctx.GetPostForm("password")
	user := models.GetUserByName(username)
	if user.ID == 0 {
		ctx.JSON(-1, gin.H{
			"code":    -1,
			"message": "用户名或密码错误",
		})
		return
	}
	checkPassword := utils.CheckPassword(password, user.PassWord, user.Salt)
	if !checkPassword {
		ctx.JSON(-1, gin.H{
			"code":    -1,
			"message": "用户名或密码错误",
		})
		return
	}
	nowPassword := utils.MakePassword(password, user.Salt)
	userNew := models.GetUserByNameAndPwd(username, nowPassword)
	if userNew.ID == 0 {
		ctx.JSON(-1, gin.H{
			"code":    -1,
			"message": "用户名或密码错误",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": userNew,
	})

}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(ctx *gin.Context) {
	ws, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		panic(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			panic(err)
			return
		}
	}(ws)
	MsgHandler(ws, ctx)
}

func MsgHandler(ws *websocket.Conn, ctx *gin.Context) {

	msg, err := utils.Subscribe(ctx, utils.PublishKey)
	if err != nil {
		panic(err)
		return
	}
	time := time.Now().Format("2006-01-02 15:04:05")
	m := fmt.Sprintf("[ws][%s]:%s", time, msg)
	err = ws.WriteMessage(1, []byte(m))
	if err != nil {
		panic(err)
		return
	}
}
