

namespace GateCliv1 
{
    //格式:			消息总长(2)+操作码(4)+收信人名单长度(2)+发送者(8)+收信人名单(0~65535)+内容(0~65535)
    //消息编解码
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

        //编码 (统一用大端)
        public byte[] Format() 
        {
            //opcode
            byte[] data2 = BitConverter.GetBytes(this.OpCode);
            if(BitConverter.IsLittleEndian)
            {
                ByteArray.Reverse(data2);
            }

            //发送者
            byte[] data4 = BitConverter.GetBytes(this.Sender);
            if(BitConverter.IsLittleEndian)
            {
                ByteArray.Reverse(data4);
            }

            //收信人
            byte[] data5;
            for(int i = 0; i < this.Recvs.Length; i++) 
            {
                byte[] temp = BitConverter.GetBytes(this.Sender);
                if(BitConverter.IsLittleEndian)
                {
                    ByteArray.Reverse(temp);
                }
                Array.Copy(temp, 0, data5, i*8, (i+1)*8);
            }

            //收信人长度
            byte[] data3 = BitConverter.GetBytes((UInt16)data5.Length);
            if(BitConverter.IsLittleEndian)
            {
                ByteArray.Reverse(data3);
            }

            //内容
            byte[] data6;
            using(MemoryStream ms = new MemoryStream())
            {
                Serializer.Serialize(ms, this.Content);
                data6 = ms.ToArray();
            }

            //总长度
            int totalLen = 2 + data2.Length + data3.Length + data4.Length + data5.Length + data6.Length;
            byte[] data1 = BitConverter.GetBytes((UInt16)totalLen);
            if(BitConverter.IsLittleEndian)
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

        //解码
        public Msg UnFormat(byte[] data) 
        {
            
        }

        //获得总长度
        public static int GetTotalLen(byte[] data) 
        {
            if(data.Length < 2) 
                return 0;
            byte[] temp = new byte[2];
            Array.Copy(data, 0, temp, 0, 2);
            if(BitConverter.IsLittleEndian)
            {
                Array.Reverse(temp);
            }
        }
    }
}