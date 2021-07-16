CREATE TABLE IF NOT EXISTS cdp.user_register (
	`name` String,
	`email` String,
	`phone` String,
	`gender` String,
	`birthday` Date,
	`userId` String,
	`region` String,
	`city` String,
	`ip`   IPv4,
	`createTime` DateTime,
	`sourceId` Int16
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(createTime)
ORDER BY (userId,createTime);

CREATE TABLE IF NOT EXISTS cdp.whole_flow (
	`feature` String,
	`userId` String,
	`region` String,
	`city` String,
	`ip`   IPv4,
	`createTime` DateTime,
	`sourceId` Int16
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(createTime)
ORDER BY (userId,createTime);