#ifndef BSP_H
#define BSP_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// Header 结构体的定义，你需要提供 Header 结构体的实现。
typedef struct Header {
    // 根据你的需求定义 Header 结构体的成员变量。
} Header;

// BspProto 结构体的定义
typedef struct BspProto {
    Header header;
    char *key;
    char **value;
    char *buf;
} BspProto;

// 函数声明
BspProto* NewBspProto();
void PutBspProto(BspProto *b);
void SetHeader(BspProto *b, Header h);
char* Key(BspProto *b);
void SetKey(BspProto *b, char *key);
char* ValueBytes(BspProto *b);
char* ValueStr(BspProto *b);
void SetValue(BspProto *b, char *value);
void SetValues(BspProto *b, char **values);
char* Buf(BspProto *b);
void SetBuf(BspProto *b, char *buf);

#endif /* BSP_H */
