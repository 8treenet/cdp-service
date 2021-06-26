# ************************************************************
# Sequel Pro SQL dump
# Version 5446
#
# https://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 5.7.18)
# Database: cdp
# Generation Time: 2021-06-26 02:45:12 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table cdp_behaviour
# ------------------------------------------------------------

DROP TABLE IF EXISTS `cdp_behaviour`;

CREATE TABLE `cdp_behaviour` (
  `id` bigint(11) unsigned NOT NULL AUTO_INCREMENT,
  `wechatUnionId` varchar(29) NOT NULL DEFAULT '' COMMENT '微信唯一id',
  `userKey` varchar(65) NOT NULL DEFAULT '' COMMENT '用户自定义key',
  `userPhone` varchar(15) NOT NULL DEFAULT '' COMMENT '用户手机号',
  `tempUserId` varchar(80) NOT NULL DEFAULT '' COMMENT '临时用户唯一id',
  `userIpAddr` varchar(16) NOT NULL DEFAULT '' COMMENT '用户ip地址',
  `featureId` smallint(128) NOT NULL COMMENT '行为的类型',
  `createTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '行为的时间',
  `data` json NOT NULL COMMENT '数据',
  `souceId` smallint(6) NOT NULL COMMENT '来源',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='行为记录表';



# Dump of table cdp_behaviour_feature
# ------------------------------------------------------------

DROP TABLE IF EXISTS `cdp_behaviour_feature`;

CREATE TABLE `cdp_behaviour_feature` (
  `id` smallint(5) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(32) NOT NULL DEFAULT '',
  `warehouse` varchar(50) NOT NULL DEFAULT '' COMMENT 'clickhouse的表名',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='行为定义表';



# Dump of table cdp_behaviour_feature_metadata
# ------------------------------------------------------------

DROP TABLE IF EXISTS `cdp_behaviour_feature_metadata`;

CREATE TABLE `cdp_behaviour_feature_metadata` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `featureId` int(11) NOT NULL COMMENT '行为的特性ID',
  `variable` varchar(32) NOT NULL DEFAULT '' COMMENT '类型的名称',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '中文的标题',
  `kind` enum('String','Float32','Float64','UInt8','UInt16','UInt32','UInt64','Int8','Int16','Int32','Int64','DateTime','Date') NOT NULL DEFAULT 'Int32' COMMENT '类型',
  `dict` varchar(128) NOT NULL DEFAULT '' COMMENT '关联字典的key',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `featureId` (`featureId`,`variable`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='行为元数据';



# Dump of table cdp_customer
# ------------------------------------------------------------

DROP TABLE IF EXISTS `cdp_customer`;

CREATE TABLE `cdp_customer` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `userId` varchar(40) NOT NULL DEFAULT '',
  `name` varchar(32) NOT NULL DEFAULT '' COMMENT '名称',
  `email` varchar(32) NOT NULL DEFAULT '' COMMENT '邮箱',
  `phone` varchar(15) NOT NULL DEFAULT '' COMMENT '电话',
  `gender` varchar(8) NOT NULL DEFAULT '' COMMENT '性别',
  `birthday` date DEFAULT NULL COMMENT '生日',
  `city` varchar(64) NOT NULL DEFAULT '' COMMENT '市',
  `region` varchar(64) NOT NULL DEFAULT '' COMMENT '省',
  `wechatUnionId` varchar(29) NOT NULL DEFAULT '' COMMENT '微信唯一id',
  `userKey` varchar(65) NOT NULL DEFAULT '' COMMENT '自定义识别key',
  `sourceId` smallint(6) NOT NULL DEFAULT '0' COMMENT '来源id',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `userId` (`userId`),
  KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='客户表';

LOCK TABLES `cdp_customer` WRITE;
/*!40000 ALTER TABLE `cdp_customer` DISABLE KEYS */;

INSERT INTO `cdp_customer` (`id`, `userId`, `name`, `email`, `phone`, `gender`, `birthday`, `city`, `region`, `wechatUnionId`, `userKey`, `sourceId`, `created`, `updated`)
VALUES
	(1,'88c60a38d5a111eb8636804a1460b6f5','yangshu3333','4932004@qq.com','135135179333','男','1989-05-13','北京','北京市','10012133333','yangshu611113513517944333',2,'2021-06-25 18:38:53','2021-06-25 18:38:53'),
	(2,'953d0154d5a111eb8636804a1460b6f5','yangshuList-0','4932004@qq.com','135135179390','男','1989-05-13','天津','天津市','100121333330','yangshu6111135135179443330',0,'2021-06-25 18:39:14','2021-06-25 18:39:14'),
	(3,'953deb28d5a111eb8636804a1460b6f5','yangshuList-1','4932004@qq.com','135135179391','男','1989-05-13','郑州市','河南','100121333331','yangshu6111135135179443331',0,'2021-06-25 18:39:14','2021-06-25 18:39:14'),
	(4,'953e42bcd5a111eb8636804a1460b6f5','yangshuList-2','4932004@qq.com','135135179392','男','1989-05-13','太原','山西','100121333332','yangshu6111135135179443332',0,'2021-06-25 18:39:14','2021-06-25 18:39:14'),
	(5,'953e941ad5a111eb8636804a1460b6f5','yangshuList-3','4932004@qq.com','135135179393','男','1989-05-13','太原','山西','100121333333','yangshu6111135135179443333',0,'2021-06-25 18:39:14','2021-06-25 18:39:14'),
	(6,'953ed916d5a111eb8636804a1460b6f5','yangshuList-4','4932004@qq.com','135135179394','男','1989-05-13','太原','山西','100121333334','yangshu6111135135179443334',0,'2021-06-25 18:39:14','2021-06-25 18:39:14');

/*!40000 ALTER TABLE `cdp_customer` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table cdp_customer_extension
# ------------------------------------------------------------

DROP TABLE IF EXISTS `cdp_customer_extension`;

CREATE TABLE `cdp_customer_extension` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `userId` varchar(40) NOT NULL DEFAULT '',
  `data` json NOT NULL COMMENT '扩展数据',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`userId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='客户扩展数据';

LOCK TABLES `cdp_customer_extension` WRITE;
/*!40000 ALTER TABLE `cdp_customer_extension` DISABLE KEYS */;

INSERT INTO `cdp_customer_extension` (`id`, `userId`, `data`, `created`, `updated`)
VALUES
	(1,'88c60a38d5a111eb8636804a1460b6f5','{\"addr\": \"11111\", \"star\": 50, \"level\": 20, \"score\": 100}','2021-06-25 18:38:53','2021-06-25 18:38:53'),
	(2,'953d0154d5a111eb8636804a1460b6f5','{\"addr\": \"刺恒小区sss-0\", \"star\": 50, \"level\": 20, \"score\": 100}','2021-06-25 18:39:14','2021-06-25 18:39:14'),
	(3,'953deb28d5a111eb8636804a1460b6f5','{\"addr\": \"刺恒小区sss-1\", \"star\": 51, \"level\": 21, \"score\": 101}','2021-06-25 18:39:14','2021-06-25 18:39:14'),
	(4,'953e42bcd5a111eb8636804a1460b6f5','{\"addr\": \"刺恒小区sss-2\", \"star\": 52, \"level\": 22, \"score\": 102}','2021-06-25 18:39:14','2021-06-25 18:39:14'),
	(5,'953e941ad5a111eb8636804a1460b6f5','{\"addr\": \"刺恒小区sss-3\", \"star\": 53, \"level\": 23, \"score\": 103}','2021-06-25 18:39:14','2021-06-25 18:39:14'),
	(6,'953ed916d5a111eb8636804a1460b6f5','{\"addr\": \"刺恒小区sss-4\", \"star\": 54, \"level\": 24, \"score\": 104}','2021-06-25 18:39:14','2021-06-25 18:39:14');

/*!40000 ALTER TABLE `cdp_customer_extension` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table cdp_customer_extension_metadata
# ------------------------------------------------------------

DROP TABLE IF EXISTS `cdp_customer_extension_metadata`;

CREATE TABLE `cdp_customer_extension_metadata` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `variable` varchar(32) NOT NULL DEFAULT '' COMMENT '类型的名称',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '中文的标题',
  `kind` enum('String','Float32','Float64','UInt8','UInt16','UInt32','UInt64','Int8','Int16','Int32','Int64','DateTime','Date') NOT NULL DEFAULT 'Int32' COMMENT '类型',
  `dict` varchar(128) NOT NULL DEFAULT '' COMMENT '关联字典的key',
  `reg` varchar(128) NOT NULL DEFAULT '' COMMENT '正则',
  `required` tinyint(4) NOT NULL DEFAULT '0' COMMENT '1 必填',
  `sort` int(11) NOT NULL DEFAULT '1000' COMMENT '排序',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `variable` (`variable`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='客户元数据';

LOCK TABLES `cdp_customer_extension_metadata` WRITE;
/*!40000 ALTER TABLE `cdp_customer_extension_metadata` DISABLE KEYS */;

INSERT INTO `cdp_customer_extension_metadata` (`id`, `variable`, `title`, `kind`, `dict`, `reg`, `required`, `sort`, `created`, `updated`)
VALUES
	(58,'score','积分','Int32','','',0,1000,'2021-06-24 18:48:46','2021-06-24 18:48:46'),
	(59,'star','关注','UInt32','','',0,1000,'2021-06-24 18:48:47','2021-06-24 18:48:47'),
	(60,'addr','地址','String','','',0,1002,'2021-06-24 18:48:47','2021-06-24 18:49:14'),
	(61,'level','级别','UInt16','','',1,1000,'2021-06-24 18:48:47','2021-06-24 18:48:47');

/*!40000 ALTER TABLE `cdp_customer_extension_metadata` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table cdp_customer_key
# ------------------------------------------------------------

DROP TABLE IF EXISTS `cdp_customer_key`;

CREATE TABLE `cdp_customer_key` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `userId` varchar(40) NOT NULL DEFAULT '0' COMMENT '用户userid',
  `userKey` varchar(65) NOT NULL DEFAULT '' COMMENT '唯一Key',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `userid` (`userId`),
  UNIQUE KEY `key` (`userKey`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='自定义key用户客户表';

LOCK TABLES `cdp_customer_key` WRITE;
/*!40000 ALTER TABLE `cdp_customer_key` DISABLE KEYS */;

INSERT INTO `cdp_customer_key` (`id`, `userId`, `userKey`, `created`, `updated`)
VALUES
	(1,'88c60a38d5a111eb8636804a1460b6f5','yangshu611113513517944333','2021-06-25 18:38:53','2021-06-25 18:38:53'),
	(2,'953d0154d5a111eb8636804a1460b6f5','yangshu6111135135179443330','2021-06-25 18:39:14','2021-06-25 18:39:14'),
	(3,'953deb28d5a111eb8636804a1460b6f5','yangshu6111135135179443331','2021-06-25 18:39:14','2021-06-25 18:39:14'),
	(4,'953e42bcd5a111eb8636804a1460b6f5','yangshu6111135135179443332','2021-06-25 18:39:14','2021-06-25 18:39:14'),
	(5,'953e941ad5a111eb8636804a1460b6f5','yangshu6111135135179443333','2021-06-25 18:39:14','2021-06-25 18:39:14'),
	(6,'953ed916d5a111eb8636804a1460b6f5','yangshu6111135135179443334','2021-06-25 18:39:14','2021-06-25 18:39:14');

/*!40000 ALTER TABLE `cdp_customer_key` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table cdp_customer_phone
# ------------------------------------------------------------

DROP TABLE IF EXISTS `cdp_customer_phone`;

CREATE TABLE `cdp_customer_phone` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `userId` varchar(40) NOT NULL DEFAULT '0' COMMENT '用户userid',
  `phone` varchar(15) NOT NULL DEFAULT '' COMMENT '手机号',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `phone` (`phone`),
  UNIQUE KEY `userId` (`userId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='电话客户表';

LOCK TABLES `cdp_customer_phone` WRITE;
/*!40000 ALTER TABLE `cdp_customer_phone` DISABLE KEYS */;

INSERT INTO `cdp_customer_phone` (`id`, `userId`, `phone`, `created`, `updated`)
VALUES
	(1,'88c60a38d5a111eb8636804a1460b6f5','135135179333','2021-06-25 18:38:53','2021-06-25 18:38:53'),
	(2,'953d0154d5a111eb8636804a1460b6f5','135135179390','2021-06-25 18:39:14','2021-06-25 18:39:14'),
	(3,'953deb28d5a111eb8636804a1460b6f5','135135179391','2021-06-25 18:39:14','2021-06-25 18:39:14'),
	(4,'953e42bcd5a111eb8636804a1460b6f5','135135179392','2021-06-25 18:39:14','2021-06-25 18:39:14'),
	(5,'953e941ad5a111eb8636804a1460b6f5','135135179393','2021-06-25 18:39:14','2021-06-25 18:39:14'),
	(6,'953ed916d5a111eb8636804a1460b6f5','135135179394','2021-06-25 18:39:14','2021-06-25 18:39:14');

/*!40000 ALTER TABLE `cdp_customer_phone` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table cdp_customer_temporary
# ------------------------------------------------------------

DROP TABLE IF EXISTS `cdp_customer_temporary`;

CREATE TABLE `cdp_customer_temporary` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(80) NOT NULL DEFAULT '' COMMENT '临时用户的唯一id',
  `userId` varchar(40) NOT NULL DEFAULT '',
  `sourceId` smallint(20) NOT NULL COMMENT '来源id',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  UNIQUE KEY `userId` (`userId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='临时用户表';

LOCK TABLES `cdp_customer_temporary` WRITE;
/*!40000 ALTER TABLE `cdp_customer_temporary` DISABLE KEYS */;

INSERT INTO `cdp_customer_temporary` (`id`, `uuid`, `userId`, `sourceId`, `created`, `updated`)
VALUES
	(4,'1624608330','1c0537eed58c11ebb60a804a1460b6f5',0,'2021-06-25 16:05:31','2021-06-25 16:05:31');

/*!40000 ALTER TABLE `cdp_customer_temporary` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table cdp_customer_wechat
# ------------------------------------------------------------

DROP TABLE IF EXISTS `cdp_customer_wechat`;

CREATE TABLE `cdp_customer_wechat` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `userId` varchar(40) NOT NULL DEFAULT '0' COMMENT '用户userid',
  `unionId` char(29) NOT NULL DEFAULT '' COMMENT '唯一id',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unionid` (`unionId`),
  UNIQUE KEY `userid` (`userId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='微信客户表';

LOCK TABLES `cdp_customer_wechat` WRITE;
/*!40000 ALTER TABLE `cdp_customer_wechat` DISABLE KEYS */;

INSERT INTO `cdp_customer_wechat` (`id`, `userId`, `unionId`, `created`, `updated`)
VALUES
	(1,'88c60a38d5a111eb8636804a1460b6f5','10012133333','2021-06-25 18:38:53','2021-06-25 18:38:53'),
	(2,'953d0154d5a111eb8636804a1460b6f5','100121333330','2021-06-25 18:39:14','2021-06-25 18:39:14'),
	(3,'953deb28d5a111eb8636804a1460b6f5','100121333331','2021-06-25 18:39:14','2021-06-25 18:39:14'),
	(4,'953e42bcd5a111eb8636804a1460b6f5','100121333332','2021-06-25 18:39:14','2021-06-25 18:39:14'),
	(5,'953e941ad5a111eb8636804a1460b6f5','100121333333','2021-06-25 18:39:14','2021-06-25 18:39:14'),
	(6,'953ed916d5a111eb8636804a1460b6f5','100121333334','2021-06-25 18:39:14','2021-06-25 18:39:14');

/*!40000 ALTER TABLE `cdp_customer_wechat` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table cdp_ip_addr
# ------------------------------------------------------------

DROP TABLE IF EXISTS `cdp_ip_addr`;

CREATE TABLE `cdp_ip_addr` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `addr` varchar(16) NOT NULL DEFAULT '' COMMENT 'ip地址',
  `country` varchar(60) NOT NULL DEFAULT '' COMMENT '国家',
  `region` varchar(60) NOT NULL DEFAULT '' COMMENT '省',
  `city` varchar(60) NOT NULL DEFAULT '' COMMENT '市',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `addr` (`addr`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `cdp_ip_addr` WRITE;
/*!40000 ALTER TABLE `cdp_ip_addr` DISABLE KEYS */;

INSERT INTO `cdp_ip_addr` (`id`, `addr`, `country`, `region`, `city`, `created`, `updated`)
VALUES
	(1251,'113.46.163.105','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1252,'119.7.146.59','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1253,'119.7.146.60','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1254,'111.117.222.104','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1255,'113.46.163.67','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1256,'119.7.146.44','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1257,'119.7.146.64','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1258,'111.117.222.153','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1259,'113.46.163.106','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1260,'111.117.222.121','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1261,'111.117.222.138','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1262,'119.7.146.43','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1263,'119.7.146.46','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1264,'111.117.222.109','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1425,'119.7.146.52','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1426,'119.7.146.57','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1427,'119.7.146.87','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1428,'111.117.222.145','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1429,'111.117.222.159','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1430,'119.7.146.94','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1431,'223.20.180.200','中国','北京市','北京','2021-06-25 18:33:00','2021-06-25 18:33:00'),
	(1432,'125.38.82.23','中国','天津市','天津','2021-06-25 18:35:10','2021-06-25 18:35:10'),
	(1433,'1.192.119.149','中国','河南','郑州市','2021-06-25 18:35:10','2021-06-25 18:35:10'),
	(1434,'','','','','2021-06-25 18:35:10','2021-06-25 18:35:10');

/*!40000 ALTER TABLE `cdp_ip_addr` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table cdp_source
# ------------------------------------------------------------

DROP TABLE IF EXISTS `cdp_source`;

CREATE TABLE `cdp_source` (
  `id` smallint(6) unsigned NOT NULL AUTO_INCREMENT,
  `source` varchar(50) NOT NULL DEFAULT '',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `source` (`source`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='来源定义';

LOCK TABLES `cdp_source` WRITE;
/*!40000 ALTER TABLE `cdp_source` DISABLE KEYS */;

INSERT INTO `cdp_source` (`id`, `source`, `created`, `updated`)
VALUES
	(2,'ali','2021-06-25 16:11:51','2021-06-25 16:11:51');

/*!40000 ALTER TABLE `cdp_source` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table cdp_system_config
# ------------------------------------------------------------

DROP TABLE IF EXISTS `cdp_system_config`;

CREATE TABLE `cdp_system_config` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL DEFAULT '' COMMENT '配置名称',
  `value` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '配置数据',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='配置表';




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
