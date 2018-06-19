## fastcgi header
```
typedef struct
{
    unsigned char version;              //版本
    unsigned char type;                 //操作类型
    unsigned char requestIdB1;          //请求id
    unsigned char requestIdB0;
    unsigned char contentLengthB1;      //内容长度
    unsigned char contentLengthB0;
    unsigned char paddingLength;        //填充字节的长度
    unsigned char reserved;             //保留字节
}FCGI_Header;

.version:用来表示FCGI的版本信息，如果是web服务器给php-fpm发送的消息，请求头中只需要将其置0就可以
.type:此字段用来说明每次所发送消息的类型，其具体值可以为如下
    type值	具体含义
    1	在与php-fpm建立连接之后发送的第一个消息中的type值就得为1，用来表明此消息为请求开始的第一个消息
    2	异常断开与php-fpm的交互
    3	在与php-fpm交互中所发的最后一个消息中type值为此，以表明交互的正常结束
    4	在交互过程中给php-fpm传递环境参数时，将type设为此，以表明消息中包含的数据为某个name-value对
    5	web服务器将从浏览器接收到的POST请求数据(表单提交等)以消息的形式发给php-fpm,这种消息的type就得设为5
    6	php-fpm给web服务器回的正常响应消息的type就设为6
    7	php-fpm给web服务器回的错误响应设为7
.requestId:此字段占俩个字节，它表示这某个特有的交互，因为php-fpm(可以理解为服务器)可以同时处理多个交互
.contentLength:此字段也占2个字节，它用来表示此消息中的消息体中数据的长度(我们上面一直说的请求头也可以叫其消息头)，我们可以据此在读消息时，能够知道读多长能读出一条完整的消息
.paddingLength:填充长度的值，为了提高处理消息的能力，我们的每个消息大小都必须为8的倍数，此长度标示，我们在消息的尾部填充的长度
.reserved:保留字段
```

## fastcgi body
消息头用一个统一的结构体表示，但是消息体不同于消息头，对于请求开始(及一次交互的第一个)的消息，有其自己的消息体格式，对于请求结束(一次交互的最后一个)的消息，有其自己的消息体格式，对于传递PARAMS参数……消息头type字段标示的消息类型不同，对应的消息体的格式就可能不同
```
1.type为1
typedef struct
{
    unsigned char roleB1;       //web服务器所期望php-fpm扮演的角色，具体取值下面有
    unsigned char roleB0;
    unsigned char flags;        //确定php-fpm处理完一次请求之后是否关闭
    unsigned char reserved[5];  //保留字段
}FCGI_BeginRequestBody;
根据上述可知type值为1的消息(标识开始请求)的消息的消息体为固定大小8字节，其中各个字段的具体含义如下
.role:此字段占2个字节，用来说明我们对php-fpm发起请求时，我们想让php-fpm为我们扮演什么角色，其常见的3个取值如下:
    role值	具体含义
    1	最常用的值，php-fpm接受我们的http所关联的信息，并产生个响应
    2	php-fpm会对我们的请求进行认证，认证通过的其会返回响应，认证不通过则关闭请求
    3	过滤请求中的额外数据流，并产生过滤后的http响应
.flags:字段确定是否与php-fpm建立长连接，为1长连接，为0则在每次请求处理结束之后关闭连接
.reserved:保留字段

2.type值为3
type值为3表示结束消息，其消息体的c定义如下
typedef struct
{
    unsigned char appStatusB3;      //结束状态，0为正常
    unsigned char appStatusB2;
    unsigned char appStatusB1;
    unsigned char appStatusB0;
    unsigned char protocolStatus;   //协议状态
    unsigned char reserved[3];
}FCGI_EndRequestBody;
同样我们可以看出结束消息体也为固定8字节大小，其各字段的具体含义如下:
.appStatus:此字段共4个字节，用来表示结束状态，0为正常结束
.protocolStatus:为协议所处的状态，0为正常状态
.reserved:为保留字节

3.type为4
此值表示此消息体为传递PARAMS(环境参数)，环境参数其实就是name-value对
我们可以使用自己定义的name-value传给php-fpm或者传递php-fpm已有的name-value对，以下为我们后面实例将会使用到的php-fpm以有的name-value对如下
    name	value
    SCRIPT_FILENAME	value值为具体.php文件所处的位置php-fpm将根据value值找到所要处理的.php文件
    REQUEST_METHOD	value值一般为GET,表示http请求的方式
回到主体，当我们要传递消息体为环境参数时，我们的消息体的格式如下
typedef struct {
     unsigned char nameLengthB3; /* nameLengthB0 >> 7 == 0 */
     unsigned char nameLengthB2;
     unsigned char nameLengthB1;
     unsigned char nameLengthB0;
     unsigned char valueLengthB3; /* nameLengthB0 >> 7 == 0 */
     unsigned char valueLengthB2;
     unsigned char valueLengthB1;
     unsigned char valueLengthB0;
     unsigned char nameData[(B3 & 0x7f) << 24) + (B2 << 16) + (B1 << 8) + B0];
     unsigned char valueData[valueLength
     ((B3 & 0x7f) << 24) + (B2 << 16) + (B1 << 8) + B0];
} FCGI_NameValue;
可以看出消息体前8个字节为固定的，其字段具体含义为
.nameLength:此字段占用4字节，用来说明name的长度
.valueLength:此字段为4个字节，用来说明value的长度
前8个字节之后紧跟的为nameLength长度的name值，接着是valueLength长度的value值

4.type值为5,6,7
当消息为输入，输出，错误时，它的消息头之后便直接跟具体数据
```

## 完整消息record
fastcgi将一个完整的消息称为record，我们每次发送的单位就是record。通过上面的介绍，我们可以总结出常见的记录格式
```
type值	record
1	    header(消息头) + 开始请求体(8字节)
3	    header + 结束请求体(8字节)
4	    header + name-value长度(8字节) + 具体的name-value
5,6,7	header + 具体内容

```