using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class ByteArray
{
    const int DEFAULT_SZIE = 1024;
    //初始大小
    int initSize = 0;
    //缓冲区
    public byte[] Bytes;
    //读写位置
    public int ReadIdx = 0;
    public int WriteIdx = 0;
    //容量
    private int capa = 0;
    //剩余空间
    public int Remain { get { return capa - WriteIdx; } }
    //数据长度
    public int Length { get { return WriteIdx - ReadIdx; } }

    /// <summary>
    /// 构造函数1
    /// </summary>
    /// <param name="size"></param>
    public ByteArray(int size = DEFAULT_SZIE) { 
        Bytes = new byte[size];
        capa = size;
        initSize = size;
        ReadIdx = 0;
        WriteIdx = 0;
    }

    /// <summary>
    /// 构造函数2
    /// </summary>
    /// <param name="defaultBytes"></param>
    public ByteArray(byte[] defaultBytes) { 
        Bytes = defaultBytes;
        capa = defaultBytes.Length;
        initSize = defaultBytes.Length;
        ReadIdx = 0;
        WriteIdx = defaultBytes.Length;
    }

    /// <summary>
    /// 重设尺寸
    /// </summary>
    /// <param name="size"></param>
    public void ReSize(int size) { 
        if(size < Length) return;
        if(size < initSize) return;
        int n = 1;
        while(n < size) n *= 2;
        capa = n;
        byte[] newBytes = new byte[capa];
        Array.Copy(Bytes, ReadIdx, newBytes, 0, WriteIdx-ReadIdx);
        Bytes = newBytes;
        WriteIdx = Length;
        ReadIdx = 0;
    }

    /// <summary>
    /// 写入数据
    /// </summary>
    /// <param name="bs"></param>
    /// <param name="offset"></param>
    /// <param name="count"></param>
    /// <returns></returns>
    public int Write(byte[] bs, int offset, int count) { 
        if(Remain < count) { 
            ReSize(Length + count);
        }
        Array.Copy(bs, offset, Bytes, WriteIdx, count);
        WriteIdx += count;
        return count;
    }

    /// <summary>
    /// 读取数据
    /// </summary>
    /// <param name="bs"></param>
    /// <param name="offset"></param>
    /// <param name="count"></param>
    /// <returns></returns>
    public int Read(byte[] bs, int offset, int count) { 
        count = Math.Min(count, Length);
        Array.Copy(Bytes, ReadIdx, bs, offset, count);
        ReadIdx += count;
        CheckAndMoveBytes();
        return count;
    }

    /// <summary>
    /// 检查并移动数据
    /// </summary>
    public void CheckAndMoveBytes() { 
        if(Length < 8) { 
            MoveBytes();
        }
    }

    /// <summary>
    /// 移动数据
    /// </summary>
    public void MoveBytes() { 
        Array.Copy(Bytes, ReadIdx, Bytes, 0, Length);
        WriteIdx = Length;
        ReadIdx = 0;
    }
}
