-- 创建数据库（需修改库名和密码）
-- CREATE DATABASE IF NOT EXISTS bobby_test DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 用户表
CREATE TABLE IF NOT EXISTS `user` (
    `id`          INT           NOT NULL AUTO_INCREMENT,
    `name`        VARCHAR(64)   NOT NULL COMMENT '用户名',
    `password`    VARCHAR(256)  NOT NULL DEFAULT '' COMMENT 'bcrypt 哈希密码',
    `dob`         VARCHAR(16)   NOT NULL DEFAULT '' COMMENT '出生日期',
    `address`     VARCHAR(256)  NOT NULL DEFAULT '' COMMENT '地址',
    `description` VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '个人描述',
    `create_at`   DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `latitude`    DOUBLE        NOT NULL DEFAULT 0 COMMENT '纬度',
    `longitude`   DOUBLE        NOT NULL DEFAULT 0 COMMENT '经度',
    `loc_geohash` VARCHAR(16)   NOT NULL DEFAULT '' COMMENT 'geohash 位置编码',
    `deleted_at`  DATETIME      DEFAULT NULL COMMENT '软删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_name` (`name`),
    KEY `idx_geohash` (`loc_geohash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 好友关系表
CREATE TABLE IF NOT EXISTS `friends` (
    `uid`         INT      NOT NULL COMMENT '用户 ID',
    `fri`         INT      NOT NULL COMMENT '好友 ID',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '建立时间',
    PRIMARY KEY (`uid`, `fri`),
    KEY `idx_uid` (`uid`),
    KEY `idx_fri` (`fri`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='好友关系表';
