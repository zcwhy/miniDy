# 抖音极简版

- [项目实现](#项目实现)
  - [技术选型](#技术选型)
  - [架构设计](#架构设计)
  - [数据库结构设计](#数据库结构设计)
    - [数据库建立说明](#数据库建立说明)
  - [项目代码介绍](#项目代码介绍)
- [项目总结与反思](#项目总结与反思)
  - [目前仍存在的问题](#目前仍存在的问题)
- [项目运行](#项目运行)
  - [运行所需环境](#运行所需环境)
  - [运行需要更改配置](#运行需要更改配置)
  - [运行所需命令](#运行所需命令)

## 项目实现
### 技术选型
项目选用gin + gorm作为主体框架
同时采用了一些开源组件库来简化部分功能的实现，包括
  1. 采用ffmpeg来对视频进行截取封面
  2. 采用golang-jwt来完成jwt中间件编写
  3. 采用toml来完成配置文件的读取
  4. 采用imaging来完成视频封面的格式等
### 架构设计
整体架构设计
![image](https://user-images.githubusercontent.com/87178677/220897373-0ad5473b-980b-4678-8c7e-3089375e06c7.png)
项目整体采用三层架构的设计方式，其中
handler层主要完成：
1. 解析得到request中的参数
2. 将参数传递给下层的service层，开始进行service层逻辑
service层主要完成：
1. 校验参数
2. 执行相应的业务逻辑
3. 将从model层得到的数据封装返回给handler
model层主要完成：
对数据库的CRUD操作

以用户登录为例：
1. 首先通过中间件，对密码进行加密
2. userLoginHandler获取userName和加密过后的password
3. userLoginService对userName和password进行校验， 并且执行userLoginDao查询对应用户的密码并进行比对
4. 将比对结果返回给handler层，handler层根据结果构造相应的response
### 数据库结构设计
根据文档提供的需求分析，构建以下几个表：
- user_logins：存下用户的用户名和密码
- user_infos：存下用户的基本信息
- videos：存下视频的基本信息
- comments：存下每个评论的基本信息
- messages：存下用户发送的消息
- user_favor_videos：存下用户喜欢的作品id
- user_relations：存下用户关注的对方用户id
所有的表都有自己的id主键为唯一的标识。         
表之间的关系如下图所示：        
![image](https://user-images.githubusercontent.com/87178677/220898753-909fad4f-86e4-4afd-84bb-692e65f1ae45.png)
![image](https://user-images.githubusercontent.com/87178677/220899050-5e5f3afe-546a-4ac5-a06a-7aa69b6c690b.png)
#### 数据库建立说明
数据库各表的建立无需自己实现额外的建表操作，一切都由gorm框架自动建表，具体逻辑在models层的代码中。

### 项目代码介绍
![image](https://user-images.githubusercontent.com/87178677/220897644-2533e93f-aee4-42e4-8375-b772a9742c79.png)

![image](https://user-images.githubusercontent.com/87178677/220897669-6fcbd2b3-a218-4650-882d-c621186bafb3.png)

项目按照handler-service-model层进行分包设计，其中handler和service层按照不同的模块进行进一步的划分

router文件夹包含所有的路由信息，util文件夹下包含项目会使用到的工具，static文件夹下用来存储上传的文件以及截取的封面，constant文件夹则包含一些项目使用到的常量

main.go文件是项目的启动点

## 项目总结与反思
### 目前仍存在的问题
1. 项目仅实现基本的接口，未进行并发优化，redis优化等
2. 对于已实现的功能，未进行详细的性能测试，压力测试等
3. 项目未完成日志功能，bug排查起来较为繁琐

## 项目运行
> 本项目运行不需要手动建表，项目启动后会自动建表。
### 运行所需环境
>mysql 5.7及以上                    
>ffmpeg（需要自行安装并配置环境变量）       

### 运行需要更改配置
>进入config目录和util目录更改对应的mysql、server信息。      
- config/config.toml：首先要先创建项目数据库，再修改mysql相关的配置信息      
- config/config,go:需要修改为自己的config.toml的路径   
- util/videos.go：需要将IP改为自己主机IP。

### 运行所需命令
首先切换到项目根目录，然后：
```
go run main.go
```
