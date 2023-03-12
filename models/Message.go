package models

import "gorm.io/gorm"

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
