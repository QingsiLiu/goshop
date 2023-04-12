-- 创建数据库
create schema goshop collate utf8_general_ci;

-- 创建用户元信息表
CREATE TABLE `user`
(
    `id`              bigint(20)                             NOT NULL AUTO_INCREMENT COMMENT '用户id',
    `nickname`        varchar(50) COLLATE utf8mb4_unicode_ci  DEFAULT NULL COMMENT '昵称',
    `avatar`          varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '头像',
    `phone`           varchar(13) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '手机号',
    `password`        varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码',
    `system_id`       int(11)                                NOT NULL,
    `client_id`       int(11)                                NOT NULL,
    `union_id`        varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'unionid',
    `last_login_time` datetime                               NOT NULL COMMENT '最后一次登录时间',
    `create_time`     datetime                               NOT NULL COMMENT '创建时间',
    `is_deleted`      tinyint(1)                             NOT NULL COMMENT '逻辑删除',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1216
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='用户';

