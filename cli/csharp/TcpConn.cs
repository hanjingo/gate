using System;
using System.Collections;
using System.Collections.Generic;
using System.Net.Sockets;
using UnityEngine;
using System.Threading;

public class TcpConn
{
    /// <summary>
    /// 构造函数
    /// </summary>
    public TcpConn() { }

    /// <summary>
    /// 套接字
    /// </summary>
    Socket socket;

    /// <summary>
    /// 接收缓冲区
    /// </summary>
    ByteArray readBuf;

    /// <summary>
    /// 写入队列
    /// </summary>
    Queue<ByteArray> writeQueue;

    /// <summary>
    /// 连接状态
    /// </summary>
    public ConnStatus Status;

    /// <summary>
    /// 消息缓存表
    /// </summary>
    List<Msg> msgList = new List<Msg>();

    /// <summary>
    /// 一次update最多可处理的消息数量
    /// </summary>
    readonly static int OnceHandleMsgNum = 10;

    /// <summary>
    /// 是否启用心跳
    /// </summary>
    public bool DoHeartBit = false;

    /// <summary>
    /// 心跳周期(秒)
    /// </summary>
    public float PingInterval = 3;

    /// <summary>
    /// 上次发ping时间
    /// </summary>
    float lastPingTime = 0;

    /// <summary>
    /// 上次收pong时间
    /// </summary>
    float lastPongTime = 0;

    /// <summary>
    /// 消息委托
    /// </summary>
    /// <param name="msg"></param>
    public delegate void MsgListener(Msg msg);

    /// <summary>
    /// 消息监听表
    /// </summary>
    Dictionary<uint, MsgListener> msgListeners = new Dictionary<uint, MsgListener>();

    /// <summary>
    /// 添加监听
    /// </summary>
    /// <param name="id"></param>
    /// <param name="listener"></param>
    public void AddMsgListener(uint id, MsgListener listener) { 
        //添加
        if(msgListeners.ContainsKey(id)) { 
            msgListeners[id] += listener;
        }
        else { 
            msgListeners[id] = listener;
        }
    }

    /// <summary>
    /// 删除监听
    /// </summary>
    /// <param name="id"></param>
    /// <param name="listener"></param>
    public void RemoveMsgListener(uint id, MsgListener listener) { 
        if(msgListeners.ContainsKey(id)) { 
            msgListeners[id] -= listener;
        }
    }

    /// <summary>
    /// 处理消息
    /// </summary>
    /// <param name="msg"></param>
    public void HandleMsg(Msg msg) 
    {
        if (msgListeners.ContainsKey(msg.Id)) { 
            msgListeners[msg.Id](msg);
        }
    }

    /// <summary>
    /// 消息集合
    /// </summary>
    public Dictionary<UInt32, Type> MsgMap = new Dictionary<uint, Type>();

    /// <summary>
    /// 添加消息
    /// </summary>
    /// <param name="id"></param>
    /// <param name="TypeStr"></param>
    public void RegMsg(UInt32 id, System.Object obj) { MsgMap[id] = obj.GetType(); }

    /// <summary>
    /// 删消息
    /// </summary>
    /// <param name="id"></param>
    public void DelMsg(UInt32 id) { if (MsgMap.ContainsKey(id)) { MsgMap.Remove(id); } }

    /// <summary>
    /// 拿消息
    /// </summary>
    /// <param name="id"></param>
    /// <returns></returns>
    public Type GetMsg(UInt32 id) 
    {
        if (MsgMap.ContainsKey(id)) 
        {
            return MsgMap[id];
        }
        return null;
    }


    /// <summary>
    /// 连接
    /// </summary>
    /// <param name="ip"></param>
    /// <param name="port"></param>
    public void Connect(string ip, int port)
    { 
        if(socket != null && Status == ConnStatus.Connected) { 
            Debug.Log("已经连接,无需再次连接");
            return;    
        }
        if(Status == ConnStatus.Connecting) { 
            Debug.Log("正在连接,请稍后再试");
            return;
        }
        //初始化状态
        InitStatus();
        //关掉negal算法
        socket.NoDelay = true;
        Status = ConnStatus.Connecting;
        socket.BeginConnect(ip, port, ConnCallBack, socket);
    }

    /// <summary>
    /// 初始化状态
    /// </summary>
    private void InitStatus() 
    { 
        socket = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
        readBuf = new ByteArray();
        writeQueue = new Queue<ByteArray>();
        Status = ConnStatus.UnConnected;
        msgList = new List<Msg>();
        lastPingTime = Time.time;
        lastPongTime = Time.time;

        RegMsg(UnKnownMsgId, new UnKnownMsg());
        RegMsg(NewConnReqId, new NewConnReq());
        RegMsg(NewConnRspId, new NewConnRsp());
        RegMsg(PingId, new Ping());
        RegMsg(PongId, new Pong());

        AddMsgListener(PongId, doPong);
    }


    /// <summary>
    /// 连接回调
    /// </summary>
    /// <param name="ar"></param>
    public void ConnCallBack(IAsyncResult ar) 
    {
        try {
            Socket socket = (Socket)ar.AsyncState;
            socket.EndConnect(ar);

            //开始接受
            socket.BeginReceive(readBuf.Bytes, readBuf.WriteIdx, readBuf.Remain, 0, ReceiveCallBack, socket);
            Status = ConnStatus.Connected;
            Debug.Log("tcp连接成功");
        } catch(SocketException e) { 
            Debug.Log("tcp连接失败,错误:" + e.ToString());
            Close();
        }
    }

    /// <summary>
    /// 关闭连接
    /// </summary>
    public void Close() 
    { 
        if(socket == null || !socket.Connected) { return; }
        if(Status != ConnStatus.Connected) { return; }
        if(writeQueue.Count > 0) 
        { 
            Status = ConnStatus.Closing; 
            return; 
        }
        else
        {
            try { 
                socket.Shutdown(SocketShutdown.Both);
                socket.Close();
                Status = ConnStatus.Closed;
                Debug.Log("tcp关闭");
            }
            catch(Exception e) {
                Debug.Log("关闭连接异常:" + e.ToString());
            }
        }
    }

    /// <summary>
    /// 发送数据
    /// </summary>
    /// <param name="msg"></param>
    public void Send(Msg msg)
    { 
        if(socket == null || !socket.Connected) { 
            return; 
        }
        if(Status != ConnStatus.Connected) { 
            return; 
        }
        byte[] data = msg.Encode();
        ByteArray ba = new ByteArray(data);
        int count = 0;
        lock(writeQueue) { 
            writeQueue.Enqueue(ba);
            count = writeQueue.Count;
        }
        if(count == 1) { 
            writeQueue.Dequeue();
            socket.BeginSend(data, 0, data.Length, 0, SendCallBack, socket);
        }
    }

    /// <summary>
    /// send回调
    /// </summary>
    /// <param name="ar"></param>
    public void SendCallBack(IAsyncResult ar) 
    { 
        Socket socket = (Socket) ar.AsyncState;
        if(socket == null || !socket.Connected) { 
            return; 
        }
        int count = socket.EndSend(ar);
        ByteArray ba;
        lock (writeQueue) { 
            //新增的
            if(writeQueue.Count == 0) { 
                return;    
            }
            ba = writeQueue.Dequeue();
        }
        if(ba != null && ba.Length > Msg.HeadLen) { 
            socket.BeginSend(ba.Bytes, ba.ReadIdx, ba.Length, 0, SendCallBack, socket);
        }
        if(Status == ConnStatus.Closing && writeQueue.Count == 0) { 
            socket.Close();
        }
    }

    /// <summary>
    /// receive回调
    /// </summary>
    /// <param name="ar"></param>
    public void ReceiveCallBack(IAsyncResult ar) { 
        if(Status != ConnStatus.Connected) { 
            return;    
        }
        try { 
            Socket socket = (Socket) ar.AsyncState;
            int count = socket.EndReceive(ar);
            readBuf.WriteIdx += count;
            //缓存消息
            doReceive();
            //如果缓冲区不够
            if(readBuf.Remain < Msg.HeadLen) { 
                readBuf.MoveBytes(); //丢掉已经用过的
                readBuf.ReSize(readBuf.Length * 2); //翻倍
            }
            socket.BeginReceive(readBuf.Bytes, readBuf.WriteIdx, readBuf.Remain, 0, ReceiveCallBack, socket);
        } catch(SocketException e) { 
            Debug.Log("tcp接收回调失败,错误:" + e.ToString());
        }
    }

    /// <summary>
    /// 处理消息
    /// </summary>
    private void doReceive() {
        if (readBuf.Length < Msg.HeadLen)
            return;
        byte[] data1 = new byte[Msg.HeadLen];
        readBuf.Read(data1, 0, data1.Length);
        UInt32 id = Msg.GetId(data1);
        UInt16 len = Msg.GetLen(data1);
        if (!MsgMap.ContainsKey(id)) { 
            Debug.Log("新到未知消息id:" + id + ",不支持此消息!!!");
            return;    
        }

        Msg msg = new Msg();
        msg.Id = id;
        Type type = MsgMap[id];
        msg.Content = type.Assembly.CreateInstance(type.ToString());
        byte[] data2 = new byte[len - Msg.HeadLen];
        readBuf.Read(data2, 0, data2.Length);

        //合并数组
        byte[] data = new byte[len];
        Buffer.BlockCopy(data1, 0, data, 0, data1.Length);
        Buffer.BlockCopy(data2, 0, data, data1.Length, data2.Length);

        msg.Decode(data);
        lock (msgList) {
            msgList.Add(msg); 
        }
        if(readBuf.Length > Msg.HeadLen) { 
            doReceive(); 
        }
    }

    //刷新
    public void Update() { 
        MsgUpdate();
        PingUpdate();
    }

    //处理消息
    public void MsgUpdate() { 
        if(msgList.Count == 0)    
            return;
        for(int i = 0; i < OnceHandleMsgNum; i++) { 
            Msg msg = null;
            lock (msgList) { 
                if(msgList.Count > 0) { 
                    msg = msgList[0];
                    msgList.RemoveAt(0);
                }
            }
            if(msg != null)
                HandleMsg(msg);
        }
    }

    /// <summary>
    /// ping
    /// </summary>
    public void PingUpdate() { 
        if(!DoHeartBit)
            return;
        if(Time.time - lastPingTime > PingInterval) { 
            Msg msg = new Msg();
            msg.Id = PingId;
            msg.Content = new Ping();
            Send(msg);
            lastPingTime = Time.time;
        }
        if(Time.time - lastPongTime > PingInterval * 3)
            Close();
    }

    /// <summary>
    /// pong响应
    /// </summary>
    public void doPong(Msg msg) { 
        Pong rsp = (Pong)msg.Content;
        lastPongTime = Time.time;
    }

    /// <summary>
    /// 连接状态枚举
    /// </summary>
    public enum ConnStatus
    {
        UnConnected,
        Connecting,
        Connected,
        Closing,
        Closed,
    }
}
