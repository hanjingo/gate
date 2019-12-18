namespace GateCliv1 {
    export class GateCli {
        private _conn: WebSocket; //ws连接

        //回调表
        private _callMap: Map<number, Function>;

        constructor(){
            this._callMap = new Map<number, Function>();
        }

        //注册回调
        public regHandler(id:number, f: Function):ErrorCode {
            if(PermMask.MASK & id < OpCode.PROTECT) {
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
            return ErrorCode.Success;
        }

        //路由
        public route(msg:any): ErrorCode {
            return ErrorCode.Success;
        }

        //ping
        public ping(): ErrorCode {
            return ErrorCode.Success;
        }

        //请求
        public req(): ErrorCode {
            return ErrorCode.Success;
        }

        //订阅
        public sub(): ErrorCode {
            return ErrorCode.Success;
        }

        //取消订阅
        public unSub(): ErrorCode {
            return ErrorCode.Success;
        }

        //注册服务
        public regService(): ErrorCode {
            return ErrorCode.Success;
        }

        //取消服务
        public unregService(): ErrorCode {
            return ErrorCode.Success;
        }

        //组播
        public multCast(): ErrorCode {
            return ErrorCode.Success;
        }

        //广播
        public BroadCast(): ErrorCode {
            return ErrorCode.Success;
        }

        //发布
        public Pub(): ErrorCode {
            return ErrorCode.Success;
        }

        //控制客户端
        public Control(): ErrorCode {
            return ErrorCode.Success;
        }
    }
}
