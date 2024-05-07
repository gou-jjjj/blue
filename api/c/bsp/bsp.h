#ifndef BSP_H
#define BSP_H

#include <stdint.h>

// HeaderInter 接口定义了头部信息处理的方法集合。
typedef struct {
    // Type 方法返回头部的类型。
    uint8_t (*Type)();

    // Handle 方法返回头部的句柄。
    uint8_t (*Handle)();

    // HandleInfo 方法返回与头部关联的命令信息。
    void* (*HandleInfo)();

    // Bytes 方法返回头部的字节表示。
    uint8_t* (*Bytes)();
} HeaderInter;

// 定义头部类型的常量。
enum {
    // TypeMask 用于掩码头部类型。
    TypeMask = 0b11100000,

    // TypeSystem 表示系统类型的头部。
    TypeSystem = 0,
    // TypeDB 表示数据库类型的头部。
    TypeDB,
    // TypeNumber 表示数字类型的头部。
    TypeNumber,
    // TypeString 表示字符串类型的头部。
    TypeString,
    // TypeList 表示列表类型的头部。
    TypeList,
    // TypeSet 表示集合类型的头部。
    TypeSet,
    // TypeJson 表示JSON类型的头部。
    TypeJson
};

// Header 定义了头部信息的类型，是一个uint8的别名。
typedef uint8_t Header;

// HandleErr 定义了一个错误句柄值。
#define HandleErr 255

// NewHeader 创建一个新的头部实例。
Header NewHeader(Header handle);

#endif /* BSP_H */
