CREATE TABLE IF NOT EXISTS cdp.user_register (
	userId String,
	name String,
	email String,
	phone String,
	gender String,
	birthday Date,
	region String,
	city String,
	ip   IPv4,
	createTime DateTime,
	souceId Int16
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(createTime)
ORDER BY userId,createTime;

