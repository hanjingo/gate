/*
* 主要用来做http配置
 */

package gate

import (
	"encoding/json"
	"net/http"
	"strings"

	network "github.com/hanjingo/golib/network"
)

//url
const (
	URL_OPEN_SERVER   string = "/open_server"
	URL_CLOSE_SERVER  string = "/close_server"
	URL_GET_GATE_INFO string = "/get_gate_info"
)

//arg
const (
	ARG_CONF      string = "conf"
	ARG_USERNAME  string = "user_name"
	ARG_PASSWORD  string = "password"
	ARG_SERV_NAME string = "server_name"
)

func (gate *Gate) reg() {
	http.HandleFunc(URL_OPEN_SERVER, gate.openServer)
	http.HandleFunc(URL_CLOSE_SERVER, gate.closeServer)
	http.HandleFunc(URL_GET_GATE_INFO, gate.getGateInfo)
}

func (gate *Gate) check(w http.ResponseWriter, r *http.Request) bool {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "http解析失败", 405)
		return false
	}
	userName := r.Form.Get(ARG_USERNAME)
	password := r.Form.Get(ARG_PASSWORD)
	return userName == gate.conf.UserName && password == gate.conf.PassWord
}

//开启服务器
func (gate *Gate) openServer(w http.ResponseWriter, r *http.Request) {
	log.Debug("openServer>>")
	if !gate.check(w, r) {
		log.Error("开启服务器时,凭证检查失败")
		return
	}
	confStr := r.Form.Get(ARG_CONF)
	conf := &network.ServerConfig{SessionConfig: &network.SessionConfig{}}
	if err := json.Unmarshal([]byte(confStr), conf); err != nil {
		log.Error("参数解析失败,错误:%v", err)
		http.Error(w, "http解析参数失败", 405)
		return
	}
	log.Debug("openServer conf:%v", conf)
	if !conf.Check() {
		log.Error("开启服务器:%v 时配置检查失败", conf.Name)
		return
	}
	switch strings.ToUpper(conf.Type) {
	case "WS":
		s, err := network.NewWsServer(conf, gate.onConnClose, gate.onNewConn, gate.handleMsg)
		if err != nil {
			log.Error("服务器启动失败,错误:%v", err)
			return
		}
		s.Run(gate.wg)
		gate.servers[s.Name] = s
		log.Info("服务器:%v 启动", s.Name)
	case "TCP":
		s, err := network.NewTcpServer(conf, gate.onConnClose, gate.onNewConn, gate.handleMsg)
		if err != nil {
			log.Error("服务器启动失败,错误:%v", err)
			return
		}
		s.Run(gate.wg)
		gate.servers[s.Name] = s
		log.Info("服务器:%v 启动", s.Name)
	default:
		log.Error("开启服务时,未知服务类型:%v", conf.Type)
	}
}

//关闭服务器
func (gate *Gate) closeServer(w http.ResponseWriter, r *http.Request) {
	log.Debug("closeServer>>")
	if !gate.check(w, r) {
		return
	}
	servName := r.Form.Get(ARG_SERV_NAME)
	if _, ok := gate.servers[servName]; !ok {
		http.Error(w, "服务器不存在", 405)
		return
	}
	delete(gate.servers, servName)
}

//获得网关信息
func (gate *Gate) getGateInfo(w http.ResponseWriter, r *http.Request) {
	log.Debug("getGateInfo>>")
	if !gate.check(w, r) {
		return
	}
	gate.info() //todo
}

//热更新插件
func (gate *Gate) hotPatch(w http.ResponseWriter, r *http.Request) {
	//todo
}
