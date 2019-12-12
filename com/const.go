package com

//未知编码id
const UNKNOWN_CODECID uint8 = 0

//掩码
const (
	MASK_ACTION uint32 = 0xff000000 //动作掩码
)

//动作
const (
	//0~7 基础功能
	ACTION_ROUTE uint8 = 0 //路由
	ACTION_PING  uint8 = 1 //ping
	ACTION_PONG  uint8 = 2 //pong

	//8～15 客户端
	ACTION_REQ uint8 = 8 //请求
	ACTION_SUB uint8 = 9 //订阅

	//64～127 服务器
	ACTION_REG_SERVER uint8 = 64 //注册服务
	ACTION_RSP        uint8 = 65 //服务器返回消息
	ACTION_MULTCAST   uint8 = 66 //多播
	ACTION_BROADCAST  uint8 = 67 //广播
	ACTION_PUB        uint8 = 68 //发布

	//128~254 控制码
	ACTION_CLOSE_AGENT              uint8 = 128 //关闭节点
	ACTION_SET_AGENT_FILTED         uint8 = 129 //过滤
	ACTION_SET_AGENT_UNFILTED       uint8 = 130 //解除过滤
	ACTION_SET_AGENT_LIMIT_STREAM   uint8 = 131 //设置节点控流限制
	ACTION_SET_AGENT_UNLIMIT_STREAM uint8 = 132 //设置节点不限流
	ACTION_CHANGE_PERM              uint8 = 133 //变更权限

	//255 系统码
	ACTION_SYS uint8 = 255
)

//权限
const (
	PERM_NORMAL  uint8 = 7   //普通权限
	PERM_CLI     uint8 = 15  //客户端权限
	PERM_SERVER  uint8 = 127 //服务器权限
	PERM_CONTROL uint8 = 254 //控制权限
	PERM_SYS     uint8 = 255 //系统权限
)

//编解码器id
const (
	CODEC_ID1 uint8 = 1 //1号编解码器
	CODEC_ID2 uint8 = 2 //2号编解码器
	CODEC_ID3 uint8 = 3 //3号编解码器
	CODEC_ID4 uint8 = 4 //4号编解码器
	CODEC_ID5 uint8 = 5 //5号编解码器
	CODEC_ID6 uint8 = 6 //6号编解码器
)
