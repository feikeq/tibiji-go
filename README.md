# 提笔记 - 肥客接口架构  Ver 3.3.0
[![The Go Programming Language](https://img.shields.io/badge/Go-v1.22-green)](https://github.com/golang/go)
[![Iris](https://img.shields.io/badge/Iris-v12.2.4-green)](https://github.com/kataras/iris)
[![sqlx](https://img.shields.io/badge/sqlx-v1.4.0-green)](https://github.com/jmoiron/sqlx)
[![lunar](https://img.shields.io/badge/lunar-v1.3.13-green)](https://github.com/6tail/lunar-go)
[![go-pinyin](https://img.shields.io/badge/go%20pinyin-v0.20.0-green)](https://github.com/mozillazg/go-pinyin)
[![Tencent Cloud SDK for Go](https://img.shields.io/badge/tencentcloud%20sdk%20go-v1.0.990-green)](https://github.com/tencentcloud/tencentcloud-sdk-go/)
[![Generator of unique non-sequential short Ids](https://img.shields.io/badge/shortid-v0.0.0-green)](https://github.com/teris-io/shortid)
[![go-vcard](https://img.shields.io/badge/go%20vcard-v0.3.0-green)](https://github.com/emersion/go-vcard)



![TiBiJ image](https://www.tibiji.com/web/images/logo.png)


肥客联邦官网：
[FK68.net](http://fk68.net)

作者：肥客泉 - [https://github.com/feikeq](https://github.com/feikeq)




## 🏆 一个全新的go项目：
1. 采用iris-go的MVC模式来更好地组织路由和控制器和模型的代码
2. 使用sqlx连接MySQL数据库并使用config/mysql.go文件单独管理
3. 并且拥有user等多个模块同时使用了控制器加模型的MVC设计模式 
4. 集成微信SDK和腾讯云短信也可采用smtp发送邮件通知
5. 集合农历库能对阴历生日进行有效计算
6. 集合中文拼音功能让检索数据更加方便
7. 支持手机通讯录vCard文件导入至联系人


## ⚙️ 项目配置
配置文件所在位置 config/cfg.ini 修改后需要重启才能生效
```ini
[Other]
    SERV_ADDR = ":8888" # 端口
    SERV_NAME = "提笔记服务端" # 项目名
    SERV_EXPIRES_TIME = 172800  # 设置token的有效时间(秒) 2 天  
    SERV_KEY_SECRET = "123456789" # API高级密钥
    SERV_OPEN_CHECK = true # 是否开启验证(登录注册是否验证)
    SERV_LIST_SIZE = 20 # 默认单页条数
    SERV_SAFE_GTIME = 30 # 获取验证码和密保的时间间隔(秒)。 30秒半分钟
    SERV_SAFE_ETIME = 1800 # 验证码和密保的有效时间(秒)。 1800秒半小时
    SERV_NOTEPAD_MAX = 1 # 非VIP用户每人最大云纸张(记事本)数量 
```


## 💻 环境变量
TIBIJI_SERV_ENV 如果没有设置默认为正式生产环境，
你可以通过在终端或命令提示符中使用以下命令来设置环境变量：

Unix系统 
```bash
export TIBIJI_SERV_ENV=development
```
Windows系统 
```bat
set TIBIJI_SERV_ENV=development
```

## ✅ API基本数据结构
* code : 错误编码 （当它返回为0时则正常否则为异常）
* msg : 错误信息 （无错误时通常为空但也有正常情况下返回的）
* data : 数据对象 （字符、数字、数组、对像等）

```json
{
    "code":0,
    "msg":"",
    "data":null
}
```

## 🏗️初始化项目
安装依赖（ go mod init <项目名> 是初始化项目依赖生成go.mod文件的 ）
也可用于移除未使用的依赖库。如果您的项目中使用了很多依赖库，但实际上只使用了其中的一部分，您可以尝试移除未使用的依赖库。
这可以通过使用Go的 go mod tidy 命令来完成执行加载依赖。
```sh
go mod tidy
```

## 🚀 启动服务
```sh
go run main.go
```

## 🛰 编译运行
不使用 go build 而使用名为 build.sh 的Shell脚本来进行编译
```sh
chmod +x build.sh
./build.sh
```
执行上面代码将编译main.go并创建3个名为tibiji-go的不同平台的可执行文件。
将编译后的可执行文件/dist/目录上传到服务器，在服务器上你可以直接运行相应平台的这个可执行文件。

这里以Linux可执行文件示例：
```sh
chmod +x tibiji-go
./tibiji-go
```
确保服务器上安装了Go语言环境，如果你的应用程序需要特定版本的Go，可以在源代码中设置go.mod文件指定版本。其中go get 命令可以简单理解为 npm install
```sh
go mod init tibiji-go
go get
go build -o tibiji-go
```
如要后台启动并使其常驻内存
```sh
 Windows ("start /b tibiji-go.exe")
 Linux ("nohup ./tibiji-go &")
 macOS ("nohup ./tibiji-go.mac &")
```


#netstat-tulpn 显示tcp和udp的侦听端口并且显示相应的程序id和程序名
(例：查看80端口是否占用 netstat -tulpn |grep 80
MAC下查看80端口是否占用 lsof -iTCP:80 | grep LISTEN )
netstat -ntlp   //查看当前所有tcp端口
也使用以下命令来查看当前占用该端口的进程：
sudo lsof -i :8888
通过PID是3230。你可以使用以下命令杀掉该进程，这将强制终止该进程。
sudo kill -9 3230

强杀进程:
#kill -9 {进程ID/进程名} 或#pkill -9 {进程名}  
sudo kill tibiji-go
 


后台启动并使其常驻内存
Linux后台运行的方法1:
--------------------
运行之前: rm nohup.out (删除一下日志文件)
运行命令： nohup ./steamcmd.sh +runscript update_csgoserv.txt &
备注:程序输出写nohup.out是按一定量字节数写入的，就是说，程序先输出到缓存区，待缓存区满1K数据之后一次性写入nohup.out文件中
查看当前进度 : tail -f nohup.out
查看steam进程: ps -ef|grep steam
杀掉steam进程: kill 进程ID


查看当前进度 : tail -f nohup.out
查看steam进程: pgrep -l srcds
杀掉steam进程: killall 名字



## 生产部署
nginx的配置署到服务器(在 Go 语言的部署过程中，通常不需要使用 PM2。)
第一种配置：直接使用 proxy_pass
```ini
location / {
    proxy_pass http://127.0.0.1:8888; 
}
```
功能：这一段配置用于设置一个简单的反向代理。当客户端请求到达 Nginx 时，它会将这些请求转发到 http://127.0.0.1:8888 这个后端服务器。
使用场景：适用于简单的代理配置，通常用在一个服务的请求不需要负载均衡或其他复杂需求时。
请求处理：该配置直接将所有匹配 / 的请求转发到指定的后端服务器上。
宝塔面板可直接在网站设-反向代理-添加反向代理 中填写 目标URLhttp://127.0.0.1:8888即可

当然你也可以加添加其它配置
```ini
location / {
  index  index.html index.htm index.php;  
  proxy_set_header Host $host;
  proxy_set_header X-Real-IP $remote_addr;
  proxy_set_header REMOTE-HOST $remote_addr;
  proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
  proxy_pass http://127.0.0.1:3000; # 后端服务器GO程序访问地址和端口，具体配置upstream部分即可  
}
```

第二种配置：使用 upstream
使用 通过 upstream nodejs 可以配置多台 nodejs 节点，做负载均衡
```ini
upstream nodenuxt {
    server 127.0.0.1:8888; 
    keepalive 64;  # 指定与后端服务器之间的保持活动连接的数量提高性能
}

location / {
    proxy_pass http://nodenuxt; 
}
```
功能：这一段配置引入了 upstream 块，定义了一个名为 nodenuxt 的后端服务器组。在此配置中，你可以将多个后端服务器添加到 upstream 中供 Nginx 负载均衡使用。
keepalive：这里的 keepalive 64 指定了与后端服务器之间的保持活动连接的数量，提高了性能。
使用场景：适合需要负载均衡、请求分发、故障转移等更复杂需求的情况。如果你将来需要将更多服务器加入到 nodenuxt 中，你只需要在 upstream 块中添加更多的 server 行，而不需要对 location 配置进行改变。
请求处理：在请求到达时，Nginx 将按照一定的负载均衡策略将请求转发到 upstream 定义的后端服务器。这可以提供更好的扩展性和容错能力。
灵活性：upstream 配置提供了更多的灵活性，可以 facilement 扩展到多个后端服务器，而直接的 proxy_pass 适合简单场景。
性能优化：upstream 可以使用 keepalive 连接，提高性能。
负载均衡：upstream 可以支持多种负载均衡算法，而直接的 proxy_pass 只能转发到一个后端目标。
一般来说，如果你需要简单的代理，可以使用直接的 proxy_pass；但如果计划将来扩展或需要负载均衡功能，使用 upstream 更为合适。
```ini
upstream nodenuxt {
    server 127.0.0.1:3000; # nuxt 项目监听端口
    keepalive 64; # keepalive 设置存活时间。如果不设置可能会产生大量的 timewait
}
server {
  listen 80;
  server_name example.com www.example.com;

  location / {
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
    proxy_set_header Host $host;
    proxy_set_header X-Nginx-Proxy true;
    proxy_cache_bypass $http_upgrade;
    proxy_pass http://nodenuxt; # 反向代理 # proxy_pass 反向代理转发 http://nodejs
  }
}
```

第三种部署方式 - 宝塔面板（推荐）
1. 网站管理中添加Go项目
2. 选择执行文件填写名称和端口
3. 执行命令有参数就用空格连接
4. 运行用户选www网站用户
5. 勾选开机自动运行
6. 如果端口不是80端口则无需绑定域名也无需放行端口否则往下执行
7. 如果端口不是80则再添加一个反向代理项目
8. 填写域名和目标 http://127.0.0.1:8888 你的本地项目端口号
9. 发送域名默认 $http_host 即可再绑上ssl证书


## 👿故障排除
Go环境未正确安装或者损坏 或 你的Go项目路径下缺少必要的源文件 或 你的Go环境可能存在缓存问题，导致编译器无法正确识别源文件。

1. 确认 Go 环境已正确安装，并且版本符合你的项目要求。
2. 确保你的项目路径中包含了各种包的源代码。如果是标准库中的包，它通常已经包含在 Go 的安装中，不需要额外处理。
3. 清除 Go 的模块缓存，可以通过执行 go clean -modcache 来清除缓存。
4. 运行 go mod tidy 后再 go run main.go 启动。
5. 如还不行可尝试删除 go.mod 文件和 go.sum 文件，然后重新运行 go mod init <module-name> 初始化模块，再运行 go mod tidy。

如果以上步骤无法解决问题，可能根据提报错上下文信息来进一步诊断问题。


## 🐬 数据割接
当前基于 MySql 8.0.19
```sql

-- 导入用户表
INSERT INTO myapp.sys_user (`uid`,`username`,`ciphers`,`email`,`cell`,`intime`,`uptime`,`regip`,`referer`,`state`,`object`)
SELECT `uid`,`username`,`password`,`email`,`cell`,`intime`,`uptime`,`regip`,`txt`,`status`,`notepad` FROM tibiji.tbj_user;

-- 导入用户附属资料表
INSERT INTO myapp.sys_material (`uid`,`cid`,`balance`,`object`,`remark`,`intime`,`uptime`)
SELECT `uid`,COALESCE(`cid`,0),`money`,`lastlog`,`allmoney`,`intime`,`uptime` FROM tibiji.tbj_user;

-- 导入用户最后一次登录日志
INSERT INTO myapp.sys_logs (`uid`,`action`,`note`,`actip`,`intime`)
SELECT `uid`, "login","username", SUBSTRING_INDEX(SUBSTRING_INDEX(lastlog,"||", 1),"::",1) AS "ip",SUBSTRING_INDEX(SUBSTRING_INDEX(lastlog,"||", 1),"::",-1) AS "time"  
FROM tibiji.tbj_user;


-- 导入开放平台用户表 (使用分割字符串为多行多条记录的方式)
INSERT INTO myapp.sys_oauth (`uid`,`platfrom`,`nickname`,`openid`,`intime`,`uptime`,`privilege`,`object`,`tidings`)
SELECT a.uid,SUBSTRING_INDEX(SUBSTRING_INDEX(a.otherid,'||',1),'::',1) as "type" ,a.username, SUBSTRING_INDEX(SUBSTRING_INDEX(a.otherid,'||',1),'::',-1) AS 'otherid', a.intime , a.uptime ,"","","" FROM tibiji.tbj_user a  WHERE a.otherid !='' 
UNION (SELECT b.uid,SUBSTRING_INDEX(SUBSTRING_INDEX(SUBSTRING_INDEX(SUBSTRING_INDEX(b.otherid,'||',2),SUBSTRING_INDEX(b.otherid,'||',1),-1), '||', -1),'::',1) as "type" , b.username, SUBSTRING_INDEX(SUBSTRING_INDEX(SUBSTRING_INDEX(SUBSTRING_INDEX(b.otherid,'||',2),SUBSTRING_INDEX(b.otherid,'||',1),-1), '||', -1),'::',-1) AS 'otherid', b.intime ,b.uptime,"","","" FROM tibiji.tbj_user b  WHERE b.otherid !='' AND SUBSTRING_INDEX(SUBSTRING_INDEX(SUBSTRING_INDEX(SUBSTRING_INDEX(b.otherid,'||',2),SUBSTRING_INDEX(b.otherid,'||',1),-1), '||', -1),'::',1) !='' )
ORDER BY `otherid` DESC 

-- 导入用户联系人
INSERT INTO myapp.tbj_contact 
(`cid`,`uid`,`fullname`,`pinyin`,`nickname`,`picture`,`phone`,`mail`,`im`,`http`,`company`,`position`,`address`,`gender`,`birthday`,`lunar`,`grouptag`,`remind`,`relation`,`family`,`note`,`state`,`intime`,`uptime`)
SELECT `cid`,`uid`,`fname`,COALESCE(`pinyin`,''),COALESCE(`nickname`,''),
 REPLACE(COALESCE(`picture`,''), './', '/')  
,COALESCE(`phone`,''),COALESCE(`mail`,''),COALESCE(`im`,''),COALESCE(`http`,''),COALESCE(`company`,''),COALESCE(`position`,''),COALESCE(`address`,''),CASE `gender` WHEN '男' THEN 1 WHEN '女' THEN  2 ELSE 0 END  AS "tonum",COALESCE(`birthday`,'0000-00-00 00:00:00'),lunar,COALESCE(`group`,''),COALESCE(`remind`,''),COALESCE(`relation`,''),COALESCE(`family`,''),COALESCE(`note`,''),`status`,`intime`,`uptime`
FROM tibiji.tbj_contact;



-- 导入记帐表
INSERT INTO myapp.tbj_account (`uid`,`item`,`class`,`sort`,`cid`,`object`,`accounts`,`money`,`note`,`btime`,`intime`,`uptime`)
SELECT `uid`,`item`,`class`,`sort`,0,`object`,`accounts`,`money`,`note`,`intime`,`intime`,`intime` FROM tibiji.tbj_account;
-- 更新记帐表CID
UPDATE `myapp`.`tbj_account` a JOIN tibiji.tbj_contact b ON a.uid = b.uid AND a.object = b.fname   SET a.`cid` =  b.`cid` WHERE a.object !='';



-- 导入记事本表
INSERT INTO myapp.tbj_notepad (`nid`,`uid`,`url`,`share`,`content`,`pwd`,`state`,`caret`,`scroll`,`ip`,`referer`,`intime`,`uptime`)
SELECT `id`,0,`url`,`share`,`content`,`pwd`,`status`,`caret`,`scroll`,`ip`,`referer`,`intime`,`uptime` FROM tibiji.tbj_notepad;
-- 更新记事本CID
UPDATE myapp.tbj_notepad  SET `uid` = `referer` WHERE referer !='' AND referer  REGEXP '^[0-9]+$';


```


## 🧩 短信模版
验证码短信：每个变量取值最多支持6位纯数字
非验证码短信：每个变量取值最多支持6个字符。对于超出变量可支持长度的内容，建议可以固定部分内容，如订单号固定前几位放到模板正文中、时间可设定为{1}年{2}月{3}日。

2196589	验证码  验证码为：{1}，您正在登录，若非本人操作，请勿泄露。

2271267	生日提醒 [提笔记生日提醒]{1}是{2}的{3}岁生日

2271314	纪念提醒 [提笔记纪念提醒]{1}是{2}的{3}周年纪念

1815721	闹铃提醒 [提笔记闹铃提醒]{1}的{2}闹铃


## 📢 通知任务
使用宝塔面板的配置方式
1. 打开计划任务添加任务
2. 任务类型选择访问URL-GET
3. 执行周期选每天固定时间
4. 填写URL地址访问的网址 User-Agent 可填 “scheduled tasks”
5. 首次填写 任务队列 的URL
6. 之后的 发送提醒 任务全基于 任务队列 之后

没发成功的会自后往后排，确保所有通知发送完毕
```
任务队列 每天的08:30执行一次
发送提醒 每天的09:00执行一次
发送提醒 每天的09:10执行一次
发送提醒 每天的09:20执行一次
发送提醒 每天的09:30执行一次
```

## 📄 文档地址
 [提笔记API接口文档接口文档](https://doc.apipost.net/docs/322e38c4e464000)

 
## 🛡 安全漏洞

如果您发现在 TiBiJi 存在安全漏洞，请发送电子邮件至 [service@tibiji.com](mailto:service@tibiji.com)。所有安全漏洞将会得到及时解决。


## 📝 开源协议

就像其它开源项目的协议一样。
