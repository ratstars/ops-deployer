// 这个包是脚本包, 包括将脚本解析成命令, 以及脚本解析工具
// 脚本格式如下：
// define executor1 type {ip:"ip", username:"username", password:"password"}
// define executor2 type {ip:"ip", username:"username", password:"password"}
// define ...
// 
// executor1 command
// in 60 expect regular_expression
// executor1 command
// executor1 command
// # comments
// 
// define开头的行定义了执行器, 以以下行为例
// define executor1 type {ip:"ip", username:"username", password:"password"}
// executor1 为自定义的执行器名字, type表示执行器类型, {......}是执行器的参数,
// 这个值会传给执行器的factory
//
// executor1 command
// executor1表示执行器的名字, command表示执行的名字
// in 60 expect regular_expression
// in 60表示执行会在60s内结束, 如果超过这个时间应会照时
// expect regular_expression 希望出现的正则表达式
// 也可以写成unexpect regular_expression 表示不希望出现的正则表达式
// #开头的行为注释行，在系统执行的过程中表示要暂停提示
// 注意，分词只支持空格不支持制表符等不可见字符
package script

import (

)

