import { 
    StringToBytes, 
    IntToBytes, 
    Uint16ToBytes, 
    Uint64ToBytes, 
    BytesToUint16, 
    BytesToInt, 
    BytesToUint64, 
    BytesToString 
} from "../../../util/transform";

namespace GateCliv1 {
    //标准起始时间 与golang的时间一致
    const TimeFormat:string = "2006-01-02 15:04:05";
    const StartTime:number = Date.parse(TimeFormat);

    /**
     * 权限掩码
     */
    export enum PermMask{
        MASK = 0x1fffffff,      //权限掩码过滤器
        ZERO = 0x0,             //0权限掩码
        NORMAL = 0x20000000,    //普通权限掩码
        CLI = 0x40000000,       //客户端权限掩码
        SERVER = 0x60000000,    //服务端权限掩码
        CONTROL = 0x80000000,   //控制网关权限掩码
        SYS = 0xa0000000,       //系统权限掩码
    }

    /**
     * 默认操作码(系统用)
     */
    export enum OpCode{
        PROTECT = 0xff,  //保护掩码 0~256 保护码以内的不允许用户使用

        //基础功能
        ROUTE = PermMask.NORMAL|0x1,              //路由
        PING = PermMask.NORMAL|0x2,               //ping
        PONG = PermMask.NORMAL|0x3,               //pong
        NEW_AGENT = PermMask.NORMAL|0x4,          //建立端点
        NEW_AGENT_RSP = PermMask.NORMAL|0x5,      //建立端点返回

        //客户端
        SUB = PermMask.CLI|0x1,                   //订阅
        UNSUB = PermMask.CLI|0x2,                 //取消订阅

        //服务器
        REG_SERVER = PermMask.SERVER|0x1,         //注册服务
        UNREG_SERVER = PermMask.SERVER|0x2,       //取消注册服务
        MULTCAST = PermMask.SERVER|0x3,           //多播
        BROADCAST = PermMask.SERVER|0x4,          //广播
        PUB = PermMask.SERVER|0x5,                //发布

        //控制码
        CLOSE_AGENT = PermMask.CONTROL|0x1,              //关闭节点
        SET_AGENT_FILTED = PermMask.CONTROL|0x2,         //过滤
        SET_AGENT_UNFILTED = PermMask.CONTROL|0x3,       //解除过滤
        SET_AGENT_LIMIT_STREAM = PermMask.CONTROL|0x4,   //设置节点控流限制
        SET_AGENT_UNLIMIT_STREAM = PermMask.CONTROL|0x5, //设置节点不限流
        CHANGE_PERM = PermMask.CONTROL|0x6,              //变更权限

        //系统码
        SYS = PermMask.SYS|0x1,
    }

    /**
     * 错误码
     */
    export enum ErrorCode {
        Fail = 0,           //失败
        Success = 1,        //成功
        CannotUseSysId=2,   //不许使用系统保留操作码
        AlreadyReg = 3,     //已经注册无需再次注册
        MsgNotComplete = 4, //消息不完整
    }

    export class GateCli {
        private _id: number; //客户端id
        private _conn: WebSocket; //ws连接
        private _codec: Codec; //编解码器

        //回调表
        private _callMap: Map<number, Function>;

        constructor(){
            this._callMap = new Map<number, Function>();
            this._codec = new Codec();
        }

        //发消息
        private writeMsg(ws: WebSocket, codec: Codec, msg: Message):ErrorCode {
            return ErrorCode.Success;
        }
        //收消息
        private readMsg(ws: WebSocket, codec: Codec, content: any): [any, ErrorCode] {
            let back: [null, ErrorCode.Success];
            return back;
        }

        //注册回调
        public regHandler(id:number, f: Function):ErrorCode {
            if(PermMask.MASK && id < OpCode.PROTECT) {
                return ErrorCode.CannotUseSysId;
            }
            if(this._callMap.has(id)){
                return ErrorCode.AlreadyReg;
            }
            this._callMap.set(id, f);
            return ErrorCode.Success;
        }

        //拨号
        public dial(addr: string): ErrorCode {
            let conn = new WebSocket(addr);
            let codec = this._codec;
            let callMap = this._callMap;
            if(conn) {
                conn.binaryType = 'arraybuffer';// 指定WebSocket接受ArrayBuffer实例作为参数
                conn.onopen = function () {
                    console.log("Send Text WS was opened.");
                };
                conn.onmessage = function (event) {
                    let recv = event.data;
                    let temp = codec.UnFormat(recv);
                    let rsp = temp[0];
                    let err = temp[1]
                    if(err == ErrorCode.Success && rsp != null && callMap.has(rsp.opcode)){
                        let f = callMap.get(rsp.opcode);
                        f.call(f, rsp.content);
                    }
                    console.log("response text msg: " + rsp);
                };
                conn.onerror = function () {
                    console.log("Send Text fired an error");
                };
                conn.onclose = function () {
                    console.log("WebSocket instance closed.");
                };

                //验证客户端 todo
                this._conn = conn;
                return ErrorCode.Success;
            }
            return ErrorCode.Fail;
        }

        //路由
        public route(content:any, ...recvs:number[]): ErrorCode {
            return ErrorCode.Success;
        }

        //ping
        public ping(): ErrorCode {
            return ErrorCode.Success;
        }

        //请求
        public req(opCode:number, content:any, ...recvs:number[]): ErrorCode {
            return ErrorCode.Success;
        }

        //订阅
        public sub(...topics:string[]): ErrorCode {
            return ErrorCode.Success;
        }

        //取消订阅
        public unSub(...topics:string[]): ErrorCode {
            return ErrorCode.Success;
        }

        //注册服务
        public regService(...apis:number[]): ErrorCode {
            return ErrorCode.Success;
        }

        //取消服务
        public unregService(...apis:number[]): ErrorCode {
            return ErrorCode.Success;
        }

        //组播
        public multCast(content:any, ...recvs:number[]): ErrorCode {
            return ErrorCode.Success;
        }

        //广播
        public BroadCast(content:any): ErrorCode {
            return ErrorCode.Success;
        }

        //发布
        public Pub(): ErrorCode {
            return ErrorCode.Success;
        }

        //控制命令
        public Control(cmd:number, ...args:any): ErrorCode {
            return ErrorCode.Success;
        }
    }

    //消息类
    export class Message {
        //操作码
        public opcode:number;
        //发送者
        public sender:number;
        //接收者名单
        public recvs:Array<number>;
        //内容
        public content:string;

        constructor() {
            this.opcode = 0; 
            this.sender = 0; 
            this.recvs = new Array<number>(); 
            this.content = ""; 
        }

        //返回指定时间到"2006-01-02 15:04:05"的秒数
        private timeToTimeStamp(ms:number):number{
            return ms - StartTime;
        }
    }

    //编解码类
    export class Codec {
        //编码
        public Format(msg:Message):[any, ErrorCode]{
            let back = new Array();
            let opData = IntToBytes(msg.opcode);
            let recvs = new Array();
            msg.recvs.forEach((v, _)=>{
                recvs.push(Uint64ToBytes(v));
            })
            let recvsLen = Uint16ToBytes(recvs.length);
            let contentData = StringToBytes(JSON.stringify(msg.content));
            let sender = Uint64ToBytes(msg.sender);
            let totalLen = Uint16ToBytes(2+4+2+8+recvs.length+contentData.length)

            back.push(totalLen, opData, recvsLen, sender, recvs, contentData);
            return [back, ErrorCode.Success];
        }

        //解码
        public UnFormat(arg:string):[Message, ErrorCode]{
            let data = StringToBytes(arg);
            if(data.length < 2+4+2+8 || data.length != this.ParseTotalLen(data)) {
                return [null, ErrorCode.MsgNotComplete];
            }
            let msg = new Message();
            msg.opcode = this.ParseOpCode(data);
            msg.sender = this.ParseSender(data);
            msg.recvs = this.ParseRecvs(data);
            msg.content = this.ParseContent(data);
            return [msg, ErrorCode.Success];
        }

        //解析总长
        public ParseTotalLen(data:any):number {
            if(data.length < 2) {
                return 0;
            }
            let temp = <Uint8Array>data;
            return BytesToUint16(temp.slice(0, 2));
        }

        //解析操作码
        public ParseOpCode(data:any):number {
            let temp = <Uint8Array>data;
            return BytesToInt(temp.slice(2, 5));
        }

        //解析发送者
        public ParseSender(data:any):number {
            let temp = <Uint8Array>data;
            return BytesToUint64(temp.slice(7, 15));
        }

        //解析收信人长度
        public ParseRecvsLen(data:any):number {
            let temp = <Uint8Array>data;
            let len = BytesToUint16(temp.slice(5,7));
            return len;
        }

        //解析收信人
        public ParseRecvs(data:any):Array<number> {
            let back = new Array<number>();
            let temp = <Uint8Array>data;
            let len = this.ParseRecvsLen(data);
            let recvData = temp.slice(15, len);
            for(let i = 0; i < recvData.length; i+=8) {
                back.push(BytesToUint64(recvData.slice(i, i+8)));
            }
            return back;
        }

        //解析内容
        public ParseContent(data:any):any {
            let temp = <Uint8Array>data;
            let start = 2 + 4 + 2 + 8;
            let end = this.ParseRecvsLen(temp);
            let contentData = temp.slice(start, end);
            return JSON.parse(BytesToString(contentData));
        }
    }
}


