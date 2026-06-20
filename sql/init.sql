-- 创建数据库（需修改库名和密码）
-- CREATE DATABASE IF NOT EXISTS bobby_test DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 用户表
CREATE TABLE IF NOT EXISTS `user` (
    `id`          INT           NOT NULL AUTO_INCREMENT,
    `name`        VARCHAR(64)   NOT NULL COMMENT '用户名',
    `password`    VARCHAR(256)  NOT NULL DEFAULT '' COMMENT 'bcrypt 哈希密码',
    `email`       VARCHAR(128)  NOT NULL DEFAULT '' COMMENT '邮箱',
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

-- 好友申请表
CREATE TABLE IF NOT EXISTS `friend_requests` (
    `id`          INT         NOT NULL AUTO_INCREMENT,
    `from_uid`    INT         NOT NULL COMMENT '发起者 UID',
    `to_uid`      INT         NOT NULL COMMENT '接收者 UID',
    `status`      VARCHAR(16) NOT NULL DEFAULT 'pending' COMMENT '状态: pending/accepted/rejected',
    `created_at`  DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_to_uid_status` (`to_uid`, `status`),
    KEY `idx_from_uid` (`from_uid`),
    UNIQUE KEY `uk_from_to_pending` (`from_uid`, `to_uid`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='好友申请表';

-- 黑名单表
CREATE TABLE IF NOT EXISTS `blacklist` (
    `uid`         INT      NOT NULL COMMENT '拉黑者 UID',
    `blocked_uid` INT      NOT NULL COMMENT '被拉黑者 UID',
    `created_at`  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '拉黑时间',
    PRIMARY KEY (`uid`, `blocked_uid`),
    KEY `idx_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户黑名单表';

-- 密码重置令牌表
CREATE TABLE IF NOT EXISTS `password_reset_tokens` (
    `id`         INT          NOT NULL AUTO_INCREMENT,
    `uid`        INT          NOT NULL COMMENT '用户 ID',
    `token`      VARCHAR(128) NOT NULL COMMENT '重置令牌',
    `expires_at` DATETIME     NOT NULL COMMENT '过期时间',
    `used`       TINYINT      NOT NULL DEFAULT 0 COMMENT '是否已使用',
    `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_token` (`token`),
    KEY `idx_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='密码重置令牌表';
