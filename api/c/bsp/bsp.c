#include "bsp.h"

// NewHeader 创建一个新的头部实例。
Header NewHeader(Header handle) {
    return handle;
}

// Type 返回头部的类型，通过与TypeMask进行与操作来获取。
uint8_t Type(uint8_t h) {
    return h & TypeMask;
}

// Handle 返回头部的句柄。
uint8_t Handle(uint8_t h) {
    return h;
}

// Bytes 返回头部的字节序列表示。
uint8_t* Bytes(uint8_t h) {
    static uint8_t byte[1];
    byte[0] = h;
    return byte;
}

// HandleInfo 返回与头部关联的命令信息。
void* HandleInfo(uint8_t h) {
    // 你需要实现这个函数来返回与头部关联的命令信息。
    return NULL;
}
