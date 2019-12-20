namespace GateCliv1 {
    /**
     * 权限掩码
     */
    enum PermMask{
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
    enum OpCode{
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
    enum ErrorCode {
        Fail = 0,           //失败
        Success = 1,        //成功
        CannotUseSysId=2,   //不许使用系统保留操作码
        AlreadyReg = 3,     //已经注册无需再次注册
    }

    public class GateCli {}
}