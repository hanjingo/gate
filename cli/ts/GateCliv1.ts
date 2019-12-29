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
    }

    export class GateCli {
        private _conn: WebSocket; //ws连接

        //回调表
        private _callMap: Map<number, Function>;

        constructor(){
            this._callMap = new Map<number, Function>();
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
            if(conn) {
                
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
        //消息总长
        public len:number;
        //操作码
        public opcode:number;
        //发送者
        public sender:number;
        //接收者长度
        public recvLen:number;
        //接收者名单
        public recvs:number[];
        //时间戳
        public timeStamp:number;
        //内容
        public content:string;

        constructor(opcode:number, content:any, recvs:number[], ...sender:number[]) {
            this.len = 0; 
            this.opcode = opcode; 
            this.sender = 0; 
            if(sender.length > 0) {
                this.sender = sender[0];
            }
            this.recvLen = 0; 
            this.recvs = recvs; 
            this.timeStamp = this.timeToTimeStamp(new Date().getTime()); 
            this.content = content; 
        }

        //返回指定时间到"2006-01-02 15:04:05"的秒数
        private timeToTimeStamp(ms:number):number{
            return ms - StartTime;
        }
    }

    //编解码类
    export class Codec {
        //编码
        public Format(msg:Message):[string, ErrorCode]{
            let back = "";
            let contentData = this.StringToBytes(JSON.stringify(msg.content));
            let opData = msg.opcode<<32;
            let recvLenData = msg.recvs.length * 8;
            return [back, ErrorCode.Success];
        }

        //解码
        public UnFormat(arg:string, content:any):[any, ErrorCode]{
            let back = content;
            return [back, ErrorCode.Success];
        }

        //string转[]byte
        public StringToBytes(str):any {
            let bytes = new Array();
            var len, c;
            len = str.length;
            for (var i = 0; i < len; i++) {
                c = str.charCodeAt(i);
                if (c >= 0x010000 && c <= 0x10FFFF) {
                    bytes.push(((c >> 18) & 0x07) | 0xF0);
                    bytes.push(((c >> 12) & 0x3F) | 0x80);
                    bytes.push(((c >> 6) & 0x3F) | 0x80);
                    bytes.push((c & 0x3F) | 0x80);
                } else if (c >= 0x000800 && c <= 0x00FFFF) {
                    bytes.push(((c >> 12) & 0x0F) | 0xE0);
                    bytes.push(((c >> 6) & 0x3F) | 0x80);
                    bytes.push((c & 0x3F) | 0x80);
                } else if (c >= 0x000080 && c <= 0x0007FF) {
                    bytes.push(((c >> 6) & 0x1F) | 0xC0);
                    bytes.push((c & 0x3F) | 0x80);
                } else {
                    bytes.push(c & 0xFF);
                }
            }
            return bytes;
        }

        //[]byte转string
        public ByteToString(byte:any): string {
            if (typeof byte === 'string') {
                return byte;
            }
            var str = '',
            _arr = byte;
            for (var i = 0; i < _arr.length; i++) {
                var one = _arr[i].toString(2),
                v = one.match(/^1+?(?=0)/);
                if (v && one.length == 8) {
                    var bytesLength = v[0].length;
                    var store = _arr[i].toString(2).slice(7 - bytesLength);
                    for (var st = 1; st < bytesLength; st++) {
                      store += _arr[st + i].toString(2).slice(2);
                    }
                    str += String.fromCharCode(parseInt(store, 2));
                    i += bytesLength - 1;
                } else {
                    str += String.fromCharCode(_arr[i]);
                }
            }
            return str;
        }
    }
}


