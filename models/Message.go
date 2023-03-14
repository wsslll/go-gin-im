package models

import (
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Message 消息models
type Message struct {
	gorm.Model
	FormId         uint   //发送者
	TargetId       uint   //接收者
	MessageType    string //消息类型 群聊 广播 私聊
	MessageMedia   int    //消息类型 文字 图片 音频
	MessageContent string //消息内容
	MessagePic     string //图片
	MessageUrl     string //url地址
	MessageDesc    string //描述
	MessageAmount  int    // 其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系
var clientMap = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(w http.ResponseWriter, r *http.Request) {
	//获取参数
	query := r.URL.Query()
	//token := query.Get("token") //TODO 检验Token
	userId, _ := strconv.Atoi(query.Get("userId"))
	//messageType := query.Get("messageType")
	//targetId := query.Get("targetId")
	//content := query.Get("content")
	isvalida := true
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln(err)
	}

	//获取连接
	node := &Node{
		conn:      conn,
		GroupSets: set.New(set.ThreadSafe),
		DataQueue: make(chan []byte),
	}

	//用户关系

	//userId和node绑定
	rwLocker.Lock()
	clientMap[int64(userId)] = node
	// 完成发送逻辑
	go sendProc(node)
	//完成接受逻辑
	go recvProc(node)

}
func sendProc(n *Node) {

}

func recvProc(n *Node) {

}
