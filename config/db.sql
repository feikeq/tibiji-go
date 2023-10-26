-- 用户API接口数据库设计
-- DB架构师：肥客泉 FK68.net
-- 创造时间：2012-04-17
-- 修改时间：2023-07-10


-- --------------------------------------------------------
-- (系统表)用户中心表 `sys_user`
-- --------------------------------------------------------
-- uid 用户ID16位(1552276542005575) = 毫秒时间戳13位(1552276542005) + 3位随机数(001)
-- ciphers 密码 = MD5(注册时间格式字符串+MD5(明文密码串))
-- intime 注册时间 为手动填入且注册时间不能变否则密码核对失败
-- object 预留字段 因是TEXT类型不能设置缺省值所以入库时尽量初始化为空字符串减少程序错误
CREATE TABLE IF NOT EXISTS `sys_user` (
    `uid` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID(自动)',
    `username` varchar(32) NOT NULL DEFAULT '' COMMENT '帐号(开放平台末绑定自动添加openid)',
    `ciphers` char(32) NOT NULL DEFAULT '' COMMENT '密码(为空则是开放台的没有绑定)',
    `email` varchar(128) NOT NULL DEFAULT '' COMMENT '邮箱(也可做登录名)',
    `cell` varchar(128) NOT NULL DEFAULT '' COMMENT '电话(也可做登录名)',
    `nickname` varchar(128) NOT NULL DEFAULT '' COMMENT '昵称',
    `headimg` varchar(255) NOT NULL DEFAULT '' COMMENT '头像',
    `sex` tinyint(1) NOT NULL DEFAULT '0' COMMENT '性别(0未知 1男 2女)',
    `birthday` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '出生日期',
    `company` varchar(100) NOT NULL DEFAULT '' COMMENT '公司',
    `address` varchar(100) NOT NULL DEFAULT '' COMMENT '地址',
    `city` varchar(50) NOT NULL DEFAULT '' COMMENT '城市',
    `province` varchar(50) NOT NULL DEFAULT '' COMMENT '省份',
    `country` varchar(50) NOT NULL DEFAULT '' COMMENT '国家',
    `regip` varchar(128) NOT NULL DEFAULT '0.0.0.0' COMMENT '注册IP地址',
    `referer` varchar(50) NOT NULL DEFAULT '' COMMENT '用户来源',
    `inviter` bigint unsigned NOT NULL DEFAULT '0' COMMENT '邀请者UID',
    `userclan` varchar(128) NOT NULL DEFAULT '' COMMENT '用户拓谱图(以逗号分隔)',
    `fname` varchar(50) NOT NULL DEFAULT '' COMMENT '真实姓名',
    `bankcard` varchar(50) NOT NULL DEFAULT '' COMMENT '银行卡',
    `identity_card` varchar(50) NOT NULL DEFAULT '' COMMENT '身份证(也可做登录名)',
    `grouptag` varchar(128) NOT NULL DEFAULT '' COMMENT '用户组',
    `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
    `object` text COMMENT '预留字段(不检索,可存json数组)',
    `state` tinyint NOT NULL DEFAULT '1' COMMENT '状态(0停用 1正常 2待激活)',
    `intime` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '注册时间',
    `uptime` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '更新时间',
    PRIMARY KEY (`uid`),
    UNIQUE KEY `username` (`username`),
    KEY `identity_card` (`identity_card`),
    KEY `email` (`email`),
    KEY `cell` (`cell`),
    KEY `grouptag` (`grouptag`),
    KEY `referer` (`referer`),
    KEY `inviter` (`inviter`),
    KEY `remark` (`remark`)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户中心表';

-- --------------------------------------------------------
-- (系统表)开放平台用户表 `sys_oauth`
-- --------------------------------------------------------
-- unionid 不一定有所以不能用作唯一索引，优先判断openid不存在再去判断unionid
-- privilege 用户特权信息 因是TEXT类型不能设置缺省值所以入库时尽量初始化为空字符串减少程序错误
-- tidings 用户动态 因是TEXT类型不能设置缺省值所以入库时尽量初始化为空字符串减少程序错误
-- object 预留字段 因是TEXT类型不能设置缺省值所以入库时尽量初始化为空字符串减少程序错误
CREATE TABLE IF NOT EXISTS `sys_oauth` (
  `oid` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '序号(自动)',
  `uid` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户ID(同一平台只能绑定一个用户ID)',
  `platfrom` varchar(50) NOT NULL DEFAULT '' COMMENT '外接平台名(微信、支付宝、AppID..)',
  `openid` varchar(128) BINARY NOT NULL DEFAULT '' COMMENT '外接平台身份ID(区分大小写)',
  `headimg` varchar(255) NOT NULL DEFAULT '' COMMENT '头像',
  `unionid` varchar(128) NOT NULL DEFAULT '' COMMENT '外接唯一标识[英文帐号]',
  `nickname` varchar(128) NOT NULL DEFAULT '' COMMENT '昵称',
  `sex` char(1) NOT NULL DEFAULT '0' COMMENT '性别(0未知 1男 2女)',
  `city` varchar(50) NOT NULL DEFAULT '' COMMENT '城市',
  `province` varchar(50) NOT NULL DEFAULT '' COMMENT '省份',
  `country` varchar(50) NOT NULL DEFAULT '' COMMENT '国家',
  `language` varchar(50) NOT NULL DEFAULT '' COMMENT '语言(zh_CN、EN)',
  `privilege` text COMMENT '用户特权信息(json数组)[关注数粉丝数等]',
  `token` varchar(128) DEFAULT '' COMMENT '平台接口授权凭证',
  `expires` varchar(128) NOT NULL DEFAULT '' COMMENT '平台授权失效时间',
  `refresh` varchar(128) NOT NULL DEFAULT '' COMMENT '平台刷新token',
  `scope` varchar(255) NOT NULL DEFAULT '' COMMENT '用户授权的作用域(使用逗号分隔)',
  `subscribe` int NOT NULL DEFAULT '0' COMMENT '是否关注[或是follow数]',
  `subscribetime` varchar(128) NOT NULL DEFAULT '' COMMENT '关注时间(时间戳)',
  `grouptag` varchar(128) NOT NULL DEFAULT '' COMMENT '分组',
  `tidings` text COMMENT '用户动态(json数组)[或最近一条社交信息]',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注(是否认证)',
  `object` text COMMENT '预留字段(不检索,可存json数组)',
  `intime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '入库时间(自动)',
  `uptime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间(自动)',
  PRIMARY KEY (`oid`),
  UNIQUE KEY `platfrom_openid` (`platfrom`,`openid`),
  KEY `uid` (`uid`),
  KEY `platfrom` (`platfrom`),
  KEY `openid` (`openid`),
  KEY `unionid` (`unionid`),
  KEY `subscribetime` (`subscribetime`),
  KEY `subscribe` (`subscribe`),
  KEY `remark` (`remark`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='开放平台用户表';

-- --------------------------------------------------------
-- (系统表)用户附属资料表 `sys_material`
-- --------------------------------------------------------
-- 存放经常变动的数据，与表用户中心是一对一关系，可以自由添加其他字段
-- object 预留字段 因是TEXT类型不能设置缺省值所以入库时尽量初始化为空字符串减少程序错误
CREATE TABLE IF NOT EXISTS `sys_material` (
  `uid` bigint unsigned NOT NULL COMMENT '用户UID',
  `cid` bigint unsigned NOT NULL DEFAULT '0' COMMENT '联系人ID',
  `balance` decimal(18,2) NOT NULL DEFAULT '0.00' COMMENT '余额（元）',
  `vip` char(1) NOT NULL DEFAULT '0' COMMENT '会员(0否 1是 2超级会员)',
  `exptime` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '会员到期时间',
  `manage` char(1) NOT NULL DEFAULT '0' COMMENT '管理权限(0普通用户 1管理员)',
  `tag` varchar(128) NOT NULL DEFAULT '' COMMENT '权限标识',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `object` text COMMENT '预留字段(不检索,可存json数组)',
  `intime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '入库时间(自动)',
  `uptime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间(自动)',
  PRIMARY KEY (`uid`),
  UNIQUE KEY `uid_cid` (`uid`,`cid`),
  KEY `balance` (`balance`),
  KEY `vip` (`vip`),
  KEY `exptime` (`exptime`),
  KEY `intime` (`intime`),
  KEY `manage` (`manage`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户附属资料表';

-- --------------------------------------------------------
-- (系统表)用户日志表 `sys_logs`
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS `sys_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `uid` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `action` varchar(50) NOT NULL DEFAULT '' COMMENT '操作标识(login,pay,sys)',
  `note` varchar(255) NOT NULL DEFAULT '' COMMENT '操作说明(微信登录、修改密码、充值等)',
  `actip` varchar(128) NOT NULL DEFAULT '' COMMENT '操作IP地址',
  `ua` varchar(255) NOT NULL DEFAULT '' COMMENT '设备信息',
  `intime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录时间',
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`),
  KEY `action` (`action`),
  KEY `intime` (`intime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户日志表';










-- --------------------------------------------------------
-- (业务表)记事本表 `tbj_notepad`
-- --------------------------------------------------------
-- content 内容 因是TEXT类型不能设置缺省值所以入库时尽量初始化为空字符串减少程序错误
CREATE TABLE IF NOT EXISTS `tbj_notepad` (
  `nid` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自动id',
  `uid` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `url` varchar(50) BINARY NOT NULL DEFAULT '' COMMENT '访问地址(区分大小写)',
  `share` varchar(50) BINARY NOT NULL DEFAULT '' COMMENT '共享地址(区分大小写)',
  `content` text COMMENT '内容',
  `pwd` varchar(30) DEFAULT '' COMMENT '密码',
  `caret` mediumint(8) DEFAULT '0' COMMENT '光标位置',
  `scroll` mediumint(8) DEFAULT '0' COMMENT '滚动位置',
  `ip` varchar(128) DEFAULT '0.0.0.0' COMMENT 'IP地址',
  `referer` varchar(255) DEFAULT '' COMMENT '纸张来源',
  `state` tinyint(2) NOT NULL DEFAULT '1' COMMENT '状态(0锁定 1正常)',
  `intime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间(自动)',
  `uptime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间(自动)',
  PRIMARY KEY (`nid`),
  UNIQUE KEY `url` (`url`),
  UNIQUE KEY `share` (`share`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='记事本表';

-- --------------------------------------------------------
-- (业务表)联系人表 `tbj_contact`
-- --------------------------------------------------------
-- phone 电话(手机,住宅,单位,自定义)格式 [TEL::xxx||TEL_CELL_WORK::xxx||TEL_FAX::xxx||TEL_FAX_HOME::xxx]
-- mail 邮箱(个人,单位,家用,自定义)格式 [EMAIL::xxx||EMAIL_HOME::xxx||EMAIL_WORK::xxx]
-- im 聊天(QQ,WeChat,Skype,Line,自定义)格式 [X-QQ::xxx||X-WECHAT::xxx]
-- http 网址(个人,单位,自定义)格式 [URL::xxx||URL_WORK::xxx]
-- 地址(家庭,单位,其他,自定义)格式 [ADR::xxx||ADR_HOME::xxx]
-- remind 提醒方式(email,phone)格式 {email::7,1,0||phone::7,1,0||..} 提前7天、1天、0当天 
CREATE TABLE IF NOT EXISTS `tbj_contact` (
  `cid` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自动id',
  `uid` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `fullname` varchar(50) NOT NULL DEFAULT '' COMMENT '姓名',
  `pinyin` varchar(50) NOT NULL DEFAULT '' COMMENT '拼音',
  `nickname` varchar(32) DEFAULT '' COMMENT '昵称绰号',
  `picture` varchar(255) DEFAULT '' COMMENT '相片照片',
  `phone` varchar(255) DEFAULT '' COMMENT '电话与传真(手机,住宅,单位,自定义)',
  `mail` varchar(255) DEFAULT '' COMMENT '邮箱(个人,单位,家用,自定义)',
  `im` varchar(255) DEFAULT '' COMMENT '聊天(QQ,WeChat,Skype,LINE,自定义)',
  `http` varchar(255) DEFAULT '' COMMENT '网址(个人,单位,自定义)',
  `company` varchar(100) DEFAULT '' COMMENT '公司部门',
  `position` varchar(100) DEFAULT '' COMMENT '职位头衔',
  `address` varchar(255) DEFAULT '' COMMENT '地址(家庭,单位,其他,自定义)',
  `gender` char(1) NOT NULL DEFAULT '0' COMMENT '性别(0未知 1男 2女)',
  `birthday` datetime DEFAULT '0000-00-00 00:00:00' COMMENT '出生时间',
  `lunar` tinyint(1) DEFAULT '0' COMMENT '是否为农历',
  `grouptag` varchar(32) DEFAULT '' COMMENT '分组',
  `remind` varchar(255) DEFAULT '' COMMENT '提醒方式(email,phone)',
  `relation` varchar(50) DEFAULT '' COMMENT '关系',
  `family` varchar(50) DEFAULT '' COMMENT '家庭户主',
  `note` varchar(2000) DEFAULT '' COMMENT '备注',
  `state` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态类别(0删除 1正常 2纪念日 3闹铃)',
  `intime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间(自动)',
  `uptime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间(自动)',
  PRIMARY KEY (`cid`),
  KEY `uid` (`uid`),
  KEY `fullname` (`fullname`),
  KEY `pinyin` (`pinyin`),
  KEY `birthday` (`birthday`),
  KEY `grouptag` (`grouptag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='联系人表';

-- --------------------------------------------------------
-- (业务表)记帐表 `tbj_account`
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tbj_account` (
  `aid` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自动id',
  `uid` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `item` varchar(10) NOT NULL DEFAULT '支出' COMMENT '操作项目(支出,收入,借款,还款)',  
  `class` varchar(10) NOT NULL DEFAULT '' COMMENT '主分类',
  `sort` varchar(10) NOT NULL DEFAULT '' COMMENT '子类别',
  `cid` bigint unsigned NOT NULL DEFAULT '0' COMMENT '收支对象ID',
  `object` varchar(32) NOT NULL DEFAULT '' COMMENT '收支对象',
  `accounts` varchar(10) NOT NULL DEFAULT '现金' COMMENT '操作账户(现金,银行卡,存折)',
  `money` decimal(18,2) NOT NULL DEFAULT '0.00' COMMENT '金额（元）',
  `note` varchar(88) DEFAULT '' COMMENT '备注说明',
  `btime` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '帐单时间',
  `intime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '入库时间(自动)',
  `uptime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间(自动)',
  PRIMARY KEY (`aid`),
  KEY `item` (`item`),
  KEY `class` (`class`),
  KEY `sort` (`sort`),
  KEY `cid` (`cid`),
  KEY `accounts` (`accounts`),
  KEY `note` (`note`),
  KEY `object` (`object`),
  KEY `btime` (`btime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='记帐表';

-- --------------------------------------------------------
-- (业务表)其它业务表
-- --------------------------------------------------------
