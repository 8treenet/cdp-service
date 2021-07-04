# ************************************************************
# Sequel Pro SQL dump
# Version 5446
#
# https://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 5.7.18)
# Database: cdp
# Generation Time: 2021-07-04 05:04:29 +0000
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
  `processed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '非0已处理',
  `souceId` smallint(6) NOT NULL COMMENT '来源',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `featureId` (`featureId`,`processed`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='行为记录表';



# Dump of table cdp_behaviour_feature
# ------------------------------------------------------------

DROP TABLE IF EXISTS `cdp_behaviour_feature`;

CREATE TABLE `cdp_behaviour_feature` (
  `id` smallint(5) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(32) NOT NULL DEFAULT '',
  `warehouse` varchar(50) NOT NULL DEFAULT '' COMMENT 'clickhouse的表名',
  `categoryType` tinyint(6) NOT NULL COMMENT '0自定义行为，1系统提供行为，2系统提供行为(不可追加)',
  `category` varchar(64) NOT NULL DEFAULT '' COMMENT '行业',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `warehouse` (`warehouse`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='行为定义表';

LOCK TABLES `cdp_behaviour_feature` WRITE;
/*!40000 ALTER TABLE `cdp_behaviour_feature` DISABLE KEYS */;

INSERT INTO `cdp_behaviour_feature` (`id`, `title`, `warehouse`, `categoryType`, `category`, `created`, `updated`)
VALUES
	(1,'用户注册','user_register',2,'通用','2021-07-02 15:54:04','2021-07-02 15:54:46');

/*!40000 ALTER TABLE `cdp_behaviour_feature` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table cdp_behaviour_feature_metadata
# ------------------------------------------------------------

DROP TABLE IF EXISTS `cdp_behaviour_feature_metadata`;

CREATE TABLE `cdp_behaviour_feature_metadata` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `featureId` int(11) NOT NULL COMMENT '行为的特性ID',
  `variable` varchar(32) NOT NULL DEFAULT '' COMMENT '类型的名称',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '中文的标题',
  `kind` enum('String','Float32','Float64','UInt8','UInt16','UInt32','UInt64','Int8','Int16','Int32','Int64','DateTime','Date','ArrayString','ArrayFloat32','ArrayFloat64','ArrayUInt8','ArrayUInt16','ArrayUInt32','ArrayUInt64','ArrayInt8','ArrayInt16','ArrayInt32','ArrayInt64','ArrayDateTime','ArrayDate') NOT NULL DEFAULT 'Int32' COMMENT '类型',
  `dict` varchar(128) NOT NULL DEFAULT '' COMMENT '关联字典的key',
  `orderByNumber` tinyint(4) NOT NULL COMMENT 'ck排序键，非0排序',
  `partition` tinyint(4) NOT NULL COMMENT '1周分区 2月分区',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `featureId` (`featureId`,`variable`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='行为元数据';

LOCK TABLES `cdp_behaviour_feature_metadata` WRITE;
/*!40000 ALTER TABLE `cdp_behaviour_feature_metadata` DISABLE KEYS */;

INSERT INTO `cdp_behaviour_feature_metadata` (`id`, `featureId`, `variable`, `title`, `kind`, `dict`, `orderByNumber`, `partition`, `created`, `updated`)
VALUES
	(5,1,'userId','userId','Int64','',1,0,'2021-06-30 17:32:47','2021-06-30 17:32:47'),
	(6,1,'name','name','String','',0,0,'2021-06-30 17:34:59','2021-06-30 17:35:10'),
	(7,1,'email','email','String','',0,0,'2021-06-30 17:34:59','2021-06-30 17:35:10'),
	(8,1,'phone','phone','String','',0,0,'2021-06-30 17:34:59','2021-06-30 17:35:10'),
	(9,1,'gender','gender','String','',0,0,'2021-06-30 17:34:59','2021-06-30 17:36:11'),
	(11,1,'birthday','birthday','Date','',0,0,'2021-06-30 17:34:59','2021-07-01 18:48:48');

/*!40000 ALTER TABLE `cdp_behaviour_feature_metadata` ENABLE KEYS */;
UNLOCK TABLES;


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
	(66,'score','积分','Int32','','',0,1000,'2021-06-30 17:39:35','2021-06-30 17:39:35'),
	(67,'star','关注','UInt32','','',0,1000,'2021-06-30 17:39:35','2021-06-30 17:39:35'),
	(68,'addr','地址','String','','',0,1000,'2021-06-30 17:39:35','2021-06-30 17:39:35'),
	(69,'level','级别','UInt16','','',1,1000,'2021-06-30 17:39:35','2021-06-30 17:39:35');

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
  UNIQUE KEY `userId` (`userId`),
  UNIQUE KEY `uuid` (`uuid`,`sourceId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='临时用户表';

LOCK TABLES `cdp_customer_temporary` WRITE;
/*!40000 ALTER TABLE `cdp_customer_temporary` DISABLE KEYS */;

INSERT INTO `cdp_customer_temporary` (`id`, `uuid`, `userId`, `sourceId`, `created`, `updated`)
VALUES
	(13,'6d01db8adaef11ebb0df804a1460b6f5','6d01db8adaef11ebb0df804a1460b6f5user',1,'2021-07-02 12:39:03','2021-07-02 12:39:03'),
	(14,'6d01de3cdaef11ebb0df804a1460b6f5','6d01de3cdaef11ebb0df804a1460b6f5user',1,'2021-07-02 12:39:03','2021-07-02 12:39:03');

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
	(1265,'111.117.222.115','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1266,'119.7.146.48','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1267,'119.7.146.53','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1268,'119.7.146.66','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1269,'113.46.163.50','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1270,'119.7.146.35','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1271,'113.46.163.104','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1272,'111.117.222.100','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1273,'111.117.222.143','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1274,'111.117.222.136','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1275,'111.117.222.141','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1276,'111.117.222.148','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1277,'111.117.222.107','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1278,'111.117.222.133','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1279,'113.46.163.74','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1280,'113.46.163.90','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1281,'119.7.146.67','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1282,'119.7.146.81','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1283,'111.117.222.119','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1284,'113.46.163.73','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1285,'111.117.222.154','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1286,'113.46.163.102','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1287,'111.117.222.116','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1288,'111.117.222.137','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1289,'113.46.163.69','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1290,'119.7.146.49','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1291,'111.117.222.118','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1292,'111.117.222.139','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1293,'113.46.163.103','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1294,'111.117.222.106','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1295,'111.117.222.156','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1296,'119.7.146.36','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1297,'119.7.146.54','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1298,'119.7.146.69','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1299,'119.7.146.79','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1300,'119.7.146.86','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1301,'111.117.222.124','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1302,'113.46.163.79','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1303,'111.117.222.150','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1304,'111.117.222.155','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1305,'113.46.163.87','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1306,'113.46.163.101','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1307,'119.7.146.41','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1308,'119.7.146.61','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1309,'111.117.222.125','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1310,'111.117.222.127','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1311,'111.117.222.146','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1312,'119.7.146.77','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1313,'113.46.163.55','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1314,'119.7.146.45','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1315,'119.7.146.72','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1316,'111.117.222.122','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1317,'111.117.222.135','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1318,'111.117.222.128','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1319,'113.46.163.54','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1320,'119.7.146.68','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1321,'119.7.146.75','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1322,'111.117.222.102','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1323,'111.117.222.120','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1324,'111.117.222.142','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1325,'113.46.163.68','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1326,'113.46.163.85','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1327,'113.46.163.93','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1328,'119.7.146.42','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1329,'119.7.146.84','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1330,'111.117.222.112','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1331,'111.117.222.123','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1332,'113.46.163.61','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1333,'113.46.163.62','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1334,'113.46.163.82','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1335,'113.46.163.108','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1336,'119.7.146.74','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1337,'111.117.222.113','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1338,'111.117.222.158','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1339,'111.117.222.126','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1340,'113.46.163.66','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1341,'113.46.163.100','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1342,'119.7.146.37','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1343,'119.7.146.51','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1344,'119.7.146.82','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1345,'111.117.222.101','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1346,'111.117.222.117','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1347,'119.7.146.63','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1348,'111.117.222.108','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1349,'111.117.222.110','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1350,'113.46.163.78','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1351,'119.7.146.47','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1352,'119.7.146.83','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1353,'111.117.222.114','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1354,'111.117.222.147','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1355,'113.46.163.75','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1356,'113.46.163.89','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1357,'119.7.146.38','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1358,'119.7.146.71','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1359,'119.7.146.91','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1360,'111.117.222.157','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1361,'113.46.163.52','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1362,'113.46.163.53','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1363,'113.46.163.76','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1364,'113.46.163.91','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1365,'113.46.163.94','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1366,'113.46.163.97','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1367,'113.46.163.98','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1368,'111.117.222.111','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1369,'111.117.222.131','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1370,'119.7.146.70','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1371,'113.46.163.99','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1372,'119.7.146.58','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1373,'113.46.163.81','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1374,'113.46.163.88','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1375,'113.46.163.96','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1376,'119.7.146.39','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1377,'119.7.146.73','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1378,'111.117.222.130','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1379,'111.117.222.132','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1380,'113.46.163.86','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1381,'113.46.163.95','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1382,'113.46.163.56','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1383,'113.46.163.84','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1384,'119.7.146.62','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1385,'119.7.146.65','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1386,'119.7.146.92','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1387,'111.117.222.129','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1388,'113.46.163.63','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1389,'113.46.163.64','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1390,'113.46.163.71','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1391,'119.7.146.55','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1392,'119.7.146.88','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1393,'119.7.146.89','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1394,'119.7.146.93','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1395,'111.117.222.144','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1396,'111.117.222.152','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1397,'113.46.163.109','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1398,'119.7.146.40','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1399,'113.46.163.72','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1400,'113.46.163.107','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1401,'111.117.222.140','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1402,'111.117.222.149','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1403,'111.117.222.151','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1404,'119.7.146.50','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1405,'119.7.146.90','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1406,'111.117.222.105','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1407,'111.117.222.134','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1408,'113.46.163.92','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1409,'119.7.146.56','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1410,'119.7.146.80','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1411,'119.7.146.85','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1412,'113.46.163.51','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1413,'113.46.163.77','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1414,'113.46.163.80','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1415,'119.7.146.78','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1416,'113.46.163.60','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1417,'113.46.163.65','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1418,'113.46.163.59','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1419,'113.46.163.70','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1420,'113.46.163.83','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1421,'119.7.146.76','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1422,'111.117.222.103','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1423,'113.46.163.57','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1424,'113.46.163.58','中国','北京市','北京','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1425,'119.7.146.52','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1426,'119.7.146.57','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1427,'119.7.146.87','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1428,'111.117.222.145','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1429,'111.117.222.159','中国','吉林','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1430,'119.7.146.94','中国','北京市','朝阳','2021-06-24 15:27:35','2021-06-24 15:27:35'),
	(1431,'223.20.180.200','中国','北京市','北京','2021-06-25 18:33:00','2021-06-25 18:33:00'),
	(1432,'125.38.82.23','中国','天津市','天津','2021-06-25 18:35:10','2021-06-25 18:35:10'),
	(1433,'1.192.119.149','中国','河南','郑州市','2021-06-25 18:35:10','2021-06-25 18:35:10');

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
