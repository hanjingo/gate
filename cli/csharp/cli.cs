using System;
using NetWork;
using Protocol;
using Protocols.V4;
using Util;

namespace GateCli
{
    /**
     * 权限掩码
     */
    public enum PermMask : UInt32
    {
        MASK = 0x1fffffff,      //权限掩码过滤器
        ZERO = 0x0,             //0权限掩码
        NORMAL = 0x10000000,    //普通权限掩码
        CLI = 0x20000000,       //客户端权限掩码
        SERVER = 0x40000000,    //服务端权限掩码
        CONTROL = 0x80000000,   //控制网关权限掩码
        SYS = 0xa0000000,       //系统权限掩码
    }

    /**
     * 默认操作码(系统用)
     */
    public enum OpCode : UInt32
    {
        PROTECT = 0xff,  //保护掩码 0~256 保护码以内的不允许用户使用

        //基础功能
        ROUTE = PermMask.NORMAL | 0x1,              //路由
        PING = PermMask.NORMAL | 0x2,               //ping
        PONG = PermMask.NORMAL | 0x3,               //pong
        NEW_AGENT_REQ = PermMask.NORMAL | 0x4,      //建立端点请求
        NEW_AGENT_RSP = PermMask.NORMAL | 0x5,      //建立端点返回
        AGENT_CONNECT_SUCCESS = PermMask.NORMAL | 0x6, //连接成功通知

        //客户端
        SUB = PermMask.CLI | 0x1,                   //订阅
        UNSUB = PermMask.CLI | 0x2,                 //取消订阅

        //服务器
        REG_SERVER = PermMask.SERVER | 0x1,         //注册服务
        UNREG_SERVER = PermMask.SERVER | 0x2,       //取消注册服务
        MULTCAST = PermMask.SERVER | 0x3,           //多播
        BROADCAST = PermMask.SERVER | 0x4,          //广播
        PUB = PermMask.SERVER | 0x5,                //发布

        //控制码
        CLOSE_AGENT = PermMask.CONTROL | 0x1,              //关闭节点

        SET_AGENT_FILTED = PermMask.CONTROL | 0x2,         //过滤
        SET_AGENT_UNFILTED = PermMask.CONTROL | 0x3,       //解除过滤
        ON_FILT = PermMask.CONTROL | 0x4,                  //过滤消息

        SET_AGENT_LIMIT_STREAM = PermMask.CONTROL | 0x5,   //设置节点限流
        SET_AGENT_UNLIMIT_STREAM = PermMask.CONTROL | 0x6, //设置节点不限流
        ON_STREAM = PermMask.CONTROL | 0x7,                //控流

        CHANGE_PERM = PermMask.CONTROL | 0x8, //变更权限

        //系统码
        SYS = PermMask.SYS | 0x1,
    }

    /**
     * 错误码
     */
    public enum ErrorCode : UInt32
    {
        Fail = 0,           //失败
        Success = 1,        //成功
        CannotUseSysId = 2,   //不许使用系统保留操作码
        AlreadyReg = 3,     //已经注册无需再次注册
    }

    /*
    * 网关客户端
    */
    public class GateCli
    {
        /// <summary>
        /// 客户端id
        /// </summary>
        public UInt64 Id;

        /// <summary>
        /// 回调器
        /// </summary>
        private HandlerI handler;

        /// <summary>
        /// 编解码器
        /// </summary>
        private CodecI codec;

        /// <summary>
        /// 连接
        /// </summary>
        private TcpConn conn;

        /// <summary>
        /// 回调器
        /// </summary>
        public class Handler
        {
            /// <summary>
            /// 处理
            /// </summary>
            /// <param name="arr"></param>
            /// <returns></returns>
            public byte[] Handle(ByteArray arr)
            {
                return null;
            }

            /// <summary>
            /// 处理ping
            /// </summary>
            /// <returns></returns>
            public byte[] Ping()
            {
                return null;
            }
        }

        public GateCli()
        {
            handler = (HandlerI)new Handler();
            codec = (CodecI)new Codec();
            conn = new TcpConn(handler);
        }

        /// <summary>
        /// 写消息
        /// </summary>
        /// <param name="opcode"></param>
        /// <param name="sender"></param>
        /// <param name="obj"></param>
        /// <param name="recvs"></param>
        /// <returns></returns>
        private ErrorCode write(OpCode opcode, UInt64 sender, System.Object obj, UInt64[] recvs)
        {
            Msg msg = new Msg((UInt32)opcode, sender, recvs, obj);
            conn.Send(codec.Format(msg));
            return ErrorCode.Success;
        }

        /// <summary>
        /// 拨号
        /// </summary>
        /// <param name="ip"></param>
        /// <param name="port"></param>
        /// <param name="token"></param>
        /// <returns></returns>
        public ErrorCode Dial(string ip, int port, string token)
        {
            conn.Connect(ip, port);
            return ErrorCode.Success;
        }

        /// <summary>
        /// 路由
        /// </summary>
        /// <param name="obj"></param>
        /// <param name="recvs"></param>
        /// <returns></returns>
        public ErrorCode Route(System.Object obj, UInt64[] recvs)
        {
            return write(OpCode.ROUTE, Id, obj, recvs);
        }

        /// <summary>
        /// ping
        /// </summary>
        /// <returns></returns>
        public ErrorCode Ping()
        {
            Ping obj = new Ping();
            return write(OpCode.PING, Id, obj, null);
        }

        /// <summary>
        /// 请求
        /// </summary>
        /// <param name="code"></param>
        /// <param name="obj"></param>
        /// <param name="recvs"></param>
        /// <returns></returns>
        public ErrorCode Req(UInt32 code, System.Object obj, UInt64[] recvs)
        {
            return write((OpCode)((UInt32)PermMask.CLI & code), Id, obj, recvs);
        }

        /// <summary>
        /// 响应
        /// </summary>
        /// <param name="code"></param>
        /// <param name="obj"></param>
        /// <param name="recvs"></param>
        /// <returns></returns>
        public ErrorCode Rsp(UInt32 code, System.Object obj, UInt64[] recvs)
        {
            return write((OpCode)((UInt32)PermMask.CLI & code), Id, obj, recvs);
        }

        /// <summary>
        /// 订阅
        /// </summary>
        /// <param name="topics"></param>
        /// <returns></returns>
        public ErrorCode Sub(string[] topics)
        {
            Sub msg = new Sub(topics);
            return write(OpCode.SUB, Id, msg, null);
        }

        /// <summary>
        /// 发布
        /// </summary>
        /// <returns></returns>
        public ErrorCode Pub(string topic, System.Object obj)
        {
            Pub msg = new Pub(topic, obj);
            return write(OpCode.PUB, Id, msg, null);
        }

        /// <summary>
        /// 取消订阅
        /// </summary>
        /// <returns></returns>
        public ErrorCode UnSub(string[] topics)
        {
            UnSub msg = new UnSub(topics);
            return write(OpCode.UNSUB, Id, msg, null);
        }

        /// <summary>
        /// 注册服务
        /// </summary>
        /// <param name="apis"></param>
        /// <returns></returns>
        public ErrorCode RegApi(UInt32[] apis)
        {
            RegApi msg = new RegApi(Id, apis);
            return write(OpCode.REG_SERVER, Id, msg, null);
        }

        /// <summary>
        /// 取消注册
        /// </summary>
        /// <param name="apis"></param>
        /// <returns></returns>
        public ErrorCode UnRegApi(UInt32[] apis)
        {
            UnRegApi msg = new UnRegApi(Id, apis);
            return write(OpCode.UNREG_SERVER, Id, msg, null);
        }

        /// <summary>
        /// 多播
        /// </summary>
        /// <param name="obj"></param>
        /// <param name="recvs"></param>
        /// <returns></returns>
        public ErrorCode MultCast(System.Object obj, UInt64[] recvs)
        {
            return write(OpCode.MULTCAST, Id, obj, recvs);
        }

        /// <summary>
        /// 广播
        /// </summary>
        /// <param name="obj"></param>
        /// <returns></returns>
        public ErrorCode BroadCast(System.Object obj)
        {
            return write(OpCode.BROADCAST, Id, obj, null);
        }

        /// <summary>
        /// 控制
        /// </summary>
        /// <param name="cmd"></param>
        /// <param name="msg"></param>
        /// <returns></returns>
        public ErrorCode Ctl(OpCode cmd, System.Object[] objs)
        {
            Control msg = new Control(cmd, objs);
            return write(cmd, Id, msg, null);
        }
    }

    /// <summary>
    /// 连接成功通知
    /// </summary>
    class AgentConnSucc
    {
        public UInt64 Id;
    }

    /// <summary>
    /// 连接验证请求
    /// </summary>
    class NewAgentReq
    {
        public string Token;
    }

    /// <summary>
    /// 连接验证返回
    /// </summary>
    class NewAgentRsp
    {
        ErrorCode Result;
        UInt64 id;
    }

    /// <summary>
    /// ping
    /// </summary>
    class Ping
    {

    }

    /// <summary>
    /// pong
    /// </summary>
    class Pong
    {

    }

    /// <summary>
    /// 订阅
    /// </summary>
    class Sub
    {
        string[] Topics;
        public Sub(string[] topics) { Topics = topics; }
    }

    /// <summary>
    /// 取消订阅
    /// </summary>
    class UnSub
    {
        string[] Topics;
        public UnSub(string[] topics) { Topics = topics; }
    }

    /// <summary>
    /// 发布
    /// </summary>
    class Pub
    {
        string Topic;
        System.Object Info;
        public Pub(string topic, System.Object info)
        {
            Topic = topic;
            Info = info;
        }
    }

    /// <summary>
    /// 注册api
    /// </summary>
    class RegApi
    {
        UInt32[] Apis;
        UInt64 Id;
        public RegApi(UInt64 id, UInt32[] apis)
        {
            Apis = apis;
            Id = id;
        }
    }

    /// <summary>
    /// 取消注册api
    /// </summary>
    class UnRegApi
    {
        UInt32[] Apis;
        UInt64 Id;
        public UnRegApi(UInt64 id, UInt32[] apis)
        {
            Apis = apis;
            Id = id;
        }
    }

    /// <summary>
    /// 控制
    /// </summary>
    class Control
    {
        UInt32 Cmd;
        System.Object[] Args;
        public Control(OpCode cmd, System.Object[] args)
        {
            Cmd = (UInt32)cmd;
            Args = args;
        }
    }
}