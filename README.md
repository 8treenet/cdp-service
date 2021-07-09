# cdp-service
cdp数据平台，帮助企业充分了解客户，实现千人千面的精准营销。

#### 开发中....

## 功能预览

提供API、MQ、SDK埋点方式，以及文件方式的导入数据。

提供实时计算的用户画像、经营报表、广告推荐以及通过CRM系统的用户触达。

![](https://tcs.teambition.net/storage/3127f48b96de8f842d52457fbf80112b1814?Signature=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBcHBJRCI6IjU5Mzc3MGZmODM5NjMyMDAyZTAzNThmMSIsIl9hcHBJZCI6IjU5Mzc3MGZmODM5NjMyMDAyZTAzNThmMSIsIl9vcmdhbml6YXRpb25JZCI6IiIsImV4cCI6MTYyNjQyMjA4NywiaWF0IjoxNjI1ODE3Mjg3LCJyZXNvdXJjZSI6Ii9zdG9yYWdlLzMxMjdmNDhiOTZkZThmODQyZDUyNDU3ZmJmODAxMTJiMTgxNCJ9.jYduJSu8BrxjLX5VjWA_R_TyuVYqblBXpSiJQ0z5ANA&download=preview.png "")



## 架构预览

本项目以2个服务组成，cdp-service是核心服务，backend-service是围绕管理后台的服务，如果有自己的管理后台服务，可以直接接入cdp-service。

![](https://tcs.teambition.net/storage/3127e711963dc6c3d194e488c0408f482aab?Signature=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBcHBJRCI6IjU5Mzc3MGZmODM5NjMyMDAyZTAzNThmMSIsIl9hcHBJZCI6IjU5Mzc3MGZmODM5NjMyMDAyZTAzNThmMSIsIl9vcmdhbml6YXRpb25JZCI6IiIsImV4cCI6MTYyNjQyMjA4NywiaWF0IjoxNjI1ODE3Mjg3LCJyZXNvdXJjZSI6Ii9zdG9yYWdlLzMxMjdlNzExOTYzZGM2YzNkMTk0ZTQ4OGMwNDA4ZjQ4MmFhYiJ9.YaK-VfihPJ_oDccyplShDhyDFcsmV30x0FL4zW0K1Ys&download=map.png "")

## 中间件模块介绍

### Kafka模块

用于行为埋点的消峰，在上万同时埋点的情况下使用mq方式。

### Prometheus模块

用于系统的监控，大数据环境的数据操作和读取较慢，需要普罗米修斯监控性能。

### JWT模块

用于后台的管理员账户鉴权。

### Echarts模块

可视化图表的HTML的处理，用于展示各种图形的汇总数据

### Flink模块

用于各种行为埋点的数据合并和数据清洗。

### Istio模块

用于系统的灰度发布，负载均衡的流量控制。



## backend-service模块介绍

### 菜单管理、权限管理、角色管理模块

动态路由，动态菜单，角色增删查改，用户增删查改，casbin鉴权。

### 文件导入、文件导出模块

转换文件形式的埋点数据，并导入cdp-service，以及导出报表数据。

### dict模块

字典常用功能，可以支持枚举等数据转换。

### 可视化报表模块

组织cdp-service服务计算出的汇总数据并且使用Echarts模块展示各种图形报表。

### 条件编辑模块

主要用来转移图形化的条件。为画像、触达、等事件图形化转移后输入到cdp-service。

### 活动管理模块

促销优惠、拉新、转化等活动创建和管理。

### 事件编辑模块

画像、触达、推荐等事件管理。

### 元数据编辑模块

创建和管理客户、行为等元数据。

### 配置管理模块

管理cdp-service的配置。



## cdp-service模块介绍

### 行为管理模块

- 创建行为

- 管理行为的元数据

- 行为缓冲区

- 批量行为入数仓

- 系统行为/通用行为/行业行为的抽象处理

### 客户管理模块

- 创建客户

- 管理客户的元数据

- 批量客户入数仓

### 事件管理和事件引擎模块

- 创建和修改事件

- 定时器执行事件

- 抽象设计事件引擎为活动、画像、推荐提供支持

- 智能触达

### 推送模块

- 抽象可扩展的第三方接入设计

- 提供默认推送功能短信、微信、scrm

### 画像分析模块

- 基于事件引擎的离线标签(可配置权重)。

- 查询引擎的实时计算(可配置权重)。

- 画像条件管理。

- 查询满足画像客户列表。

- 手动打标签。

### 相似推荐模块

- 基于相似画像客户的行为推荐(可配置权重)。

- 基于相似推荐条件的行为推荐(可配置权重)。

- 基于CatBoots机器学习的行为推荐(可配置权重)。

### 数据清洗、和数据合并模块

使用Flink做数据清洗和数据合并，可在backend-service服务配置条件。

### GEO模块

- 为行为的IP地址提供地理位置查询。

- 为行为的经纬度提供地理位置查询。

### AB/Test

- 增删查改管理。

- 为事件引擎提供分流。

- 提供实验结果报表

### 查询引擎

- 解析步骤、条件、关联成查询对象。

- 解析查询对象为SQL。

- 可扩展、抽象的链式对象设计。

- Clickhouse特性封装。

- DSL查询描述层。