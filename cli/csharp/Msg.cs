

namespace GateCliv1 
{
    /// <summary>
    /// 消息
    /// </summary>
    class Msg
    {
        //操作码
        public uint32 OpCode;

        //接受者
        public uint64[] Recvs;

        //发送者
        public uint64 Sender;

        //内容
        public System.Object Content;

        Msg(System.Object content) {
            this.Content = content;
        }
    }

    /// <summary>
    /// 编解码器
    /// 格式:消息总长(2)+操作码(4)+收信人名单长度(2)+发送者(8)+收信人名单(0~65535)+内容(0~65535)
    /// </summary>
    public static class Codec
    {
        //总长字段长度
        private static int lenLen = 2;
        //操作码字段长度
        private static int opcodeLen = 4;
        //发送者字段长度
        private static int senderLen = 8;
        //接收者长度字段长度
        private static int recvLenLen = 2;
        //id字段长度
        private static int idLen = 8;

        /// <summary>
        /// 编码 (统一用大端)
        /// </summary>
        /// <returns></returns>
        public static byte[] Format(Msg msg)
        {
            //opcode
            byte[] data2 = BitConverter.GetBytes(msg.OpCode);
            if (BitConverter.IsLittleEndian)
            {
                ByteArray.Reverse(data2);
            }

            //发送者
            byte[] data4 = BitConverter.GetBytes(msg.Sender);
            if (BitConverter.IsLittleEndian)
            {
                ByteArray.Reverse(data4);
            }

            //收信人
            byte[] data5;
            for (int i = 0; i < msg.Recvs.Length; i++)
            {
                byte[] temp = BitConverter.GetBytes(msg.Sender);
                if (BitConverter.IsLittleEndian)
                {
                    ByteArray.Reverse(temp);
                }
                Array.Copy(temp, 0, data5, i * 8, (i + 1) * 8);
            }

            //收信人长度
            byte[] data3 = BitConverter.GetBytes((UInt16)data5.Length);
            if (BitConverter.IsLittleEndian)
            {
                ByteArray.Reverse(data3);
            }

            //内容
            byte[] data6;
            using (MemoryStream ms = new MemoryStream())
            {
                Serializer.Serialize(ms, msg.Content);
                data6 = ms.ToArray();
            }

            //总长度
            int totalLen = 2 + data2.Length + data3.Length + data4.Length + data5.Length + data6.Length;
            byte[] data1 = BitConverter.GetBytes((UInt16)totalLen);
            if (BitConverter.IsLittleEndian)
            {
                ByteArray.Reverse(data1);
            }

            //合并
            byte[] back = new byte[totalLen];
            ByteArray array = new ByteArray();
            array.Write(data1, 0, data1.Length);
            array.Write(data2, 0, data2.Length);
            array.Write(data3, 0, data3.Length);
            array.Write(data4, 0, data4.Length);
            array.Write(data5, 0, data5.Length);
            array.Write(data6, 0, data6.Length);
            array.Read(back, 0, totalLen);
            return back;
        }

        /// <summary>
        /// 解码
        /// </summary>
        /// <param name="data"></param>
        /// <returns></returns>
        public static Msg UnFormat(byte[] data, System.Object content)
        {
            Msg back = new Msg(content);
            back.OpCode = ParseOpCode(data);
            back.Sender = ParseSender(data);
            back.Recvs = ParseRecvs(data);
            back.Content = ParseContent(data);
            return back;
        }

        /// <summary>
        /// 获得总长度
        /// </summary>
        /// <param name="data"></param>
        /// <returns></returns>
        public static int ParseTotalLen(byte[] data)
        {
            int start = 0;
            int end = start + lenLen;
            if (data == null || data.Length < end)
                return 0;
            
            byte[] temp = new byte[lenLen];
            Array.Copy(data, start, temp, 0, end);
            if (BitConverter.IsLittleEndian)
            {
                Array.Reverse(temp);
            }
            UInt16 back = BitConverter.ToUInt16(temp, 0);
            return back;
        }

        /// <summary>
        /// 解析opcode
        /// </summary>
        /// <param name="data"></param>
        /// <returns></returns>
        public static UInt32 ParseOpCode(byte[] data)
        {
            int start = lenLen;
            int end = start + opcodeLen;
            if (data == null || data.length < end) 
                return 0;

            byte[] temp = new byte[opcodeLen];
            Array.Copy(data, start, temp, 0, opcodeLen);
            if (BitConverter.IsLittleEndian)
            {
                Array.Reverse(temp);
            }
            UInt32 back = BitConverter.ToUInt32(temp, 0);
            return back;
        }

        /// <summary>
        /// 解析收信人名单长度
        /// </summary>
        /// <param name="data"></param>
        /// <returns></returns>
        public static UInt16 ParseRecvLen(byte[] data)
        {
            int start = lenLen + opcodeLen;
            int end = start + recvLenLen;
            if (data == null || data.length < end)
                return 0;

            byte[] temp = new byte[recvLenLen];
            Array.Copy(data, start, temp, 0, recvLenLen);
            if (BitConverter.IsLittleEndian)
            {
                Array.Reverse(temp);
            }
            UInt16 back = BitConverter.ToUInt16(temp, 0);
            return back;
        }

        /// <summary>
        /// 解析发送者
        /// </summary>
        /// <param name="data"></param>
        /// <returns></returns>
        public static UInt64 ParseSender(byte[] data)
        {
            int start = lenLen + opcodeLen + recvLenLen;
            int end = start + senderLen;
            if (data == null || data.length < end)
                return 0;

            byte[] temp = new byte[senderLen];
            Array.Copy(data, start, temp, 0, senderLen);
            if (BitConverter.IsLittleEndian)
            {
                Array.Reverse(temp);
            }
            UInt64 back = BitConverter.ToUInt64(temp, 0);
            return back;
        }

        /// <summary>
        /// 解析收信人名单
        /// </summary>
        /// <param name="data"></param>
        /// <returns></returns>
        public static UInt64[] ParseRecvs(byte[] data)
        {
            int recvLen = ParseRecvLen(data);
            int start = lenLen + opcodeLen + recvLenLen + senderLen;
            int end = start + recvLen;
            int idNum = recvLen / idLen;
            if (data == null || data.length < end)
                return 0;

            UInt64[] back = new UInt64[idNum];
            for (int i = start, j = 0; i <= end; i += idLen) {
                byte[] temp = new byte[idLen];
                Array.Copy(data, start, temp, 0, idLen);
                if (BitConverter.IsLittleEndian)
                    Array.Reverse(temp);
                UInt64 id = BitConverter.ToUInt64(temp, 0);
                back[j] = id;
                j++;
            }
            return back;
        }

        /// <summary>
        /// 解析内容
        /// </summary>
        /// <param name="data"></param>
        public static System.Object ParseContent(byte[] data)
        {
            int recvLen = ParseRecvLen(data);
            int start = lenLen + opcodeLen + recvLenLen + 
                senderLen + recvLen;
            int end = ParseTotalLen(data);
            if (data == null || data.length < end)
                return 0;
            byte[] temp = new byte[end - start];
            Array.Copy(data, start, temp, 0, temp.length);
            //todo 解析json
            return null;
        }
    }
}