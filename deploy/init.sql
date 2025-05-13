-- 用户表
use  sys;
create database tag;
use tag;
CREATE TABLE `user` (
                        `id` VARCHAR(36) NOT NULL COMMENT 'UUID',
                        `email` VARCHAR(255) NOT NULL COMMENT '邮箱',
                        `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希',
                        `username` VARCHAR(50) NOT NULL COMMENT '用户名',
                        `avatar_url` VARCHAR(255) NULL COMMENT '头像URL',
                        theme VARCHAR(20) DEFAULT 'light' COMMENT '主题偏好',
                        `bio`  VARCHAR(255) COMMENT '用户简介',
                        register_date DATETIME NOT NULL COMMENT '注册时间',
                        `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                        `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户基础信息表';-- 用户表

-- 用户订阅表
CREATE TABLE `user_subscription` (
                                     `id` VARCHAR(36) NOT NULL COMMENT 'UUID',
                                     `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
                                     `plan_type` VARCHAR(20) NOT NULL COMMENT '订阅类型:free,pro,premium',
                                     `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态:active,canceled,expired',
                                     `start_date` DATETIME NOT NULL COMMENT '开始时间',
                                     `end_date` DATETIME NULL COMMENT '结束时间',
                                     `payment_method` VARCHAR(50) NULL COMMENT '支付方式',
                                     `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                     `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                     PRIMARY KEY (`id`),
                                     KEY `idx_user_id` (`user_id`),
                                     CONSTRAINT `fk_subscription_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户订阅信息表';

-- 组织表
CREATE TABLE `organization` (
                                `id` VARCHAR(36) NOT NULL COMMENT 'UUID',
                                `code` VARCHAR(64) NOT NULL COMMENT '节点编码',
                                `parent_code` VARCHAR(64) DEFAULT '0' COMMENT '节点上级编码',
                                `parent_codes` VARCHAR(1000) DEFAULT '' COMMENT '节点所有上级编码',
                                `tree_leaf` CHAR(1) DEFAULT '1' COMMENT '是否叶子节点(0:否 1:是)',
                                `tree_level` DECIMAL(4,0) DEFAULT 0 COMMENT '节点层次级别(从0开始)',
                                `tree_names` VARCHAR(1000) DEFAULT '' COMMENT '节点全名称(用/分隔)',
                                `name` VARCHAR(100) NOT NULL COMMENT '组织名称',
                                `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
                                `description` VARCHAR(500) NULL COMMENT '描述',
                                `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序顺序',
                                `items_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '组织项目数',
                                `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                PRIMARY KEY (`id`),
                                UNIQUE KEY `uid_code` (`code`,`user_id`),
                                KEY `idx_parent_code` (`parent_code`),
                                KEY `idx_user_code` (`user_id`),
                                CONSTRAINT `fk_organization_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组织表';

-- 收藏项目表

CREATE TABLE `item` (
                        `id` VARCHAR(36) NOT NULL COMMENT 'UUID',
                        `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
                        `item_type` VARCHAR(20) NOT NULL COMMENT '类型:text,image,link,bookmark',
                        `title` VARCHAR(500) NULL COMMENT '标题',
                        `note` TEXT NULL COMMENT '内容/文本',
                        `url` VARCHAR(2000) NULL COMMENT '链接地址',
                        `thumbnail_url` VARCHAR(2000) NULL COMMENT '缩略图地址',
                        `tag_names` varchar(500) COMMENT '所有标签名',
                        `organization_path` varchar(500) COMMENT '组织路径',
                        `organization_id`VARCHAR(36) NOT NULL COMMENT '组织ID',
                        `deleted_at` DATETIME NULL COMMENT '删除时间',
                        `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                        `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                        PRIMARY KEY (`id`),
                        KEY `idx_user_id` (`user_id`),
                        KEY `idx_type` (`item_type`),
                        KEY `idx_deleted_at` (`deleted_at`),
                        CONSTRAINT `fk_item_organization_id` FOREIGN KEY (`organization_id`) REFERENCES `organization` (`id`) ON DELETE CASCADE,
                        CONSTRAINT `fk_item_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='收藏项目表';


-- 标签表
CREATE TABLE `tag` (
                       `id` VARCHAR(36) NOT NULL COMMENT 'UUID',
                       `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
                       `name` VARCHAR(100) NOT NULL COMMENT '标签名称',
                       `color` VARCHAR(20) NULL COMMENT '标签颜色',
                        `icon` varchar(100) NULL COMMENT '标签颜色',
                       `use_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '使用次数',
                       `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                       `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                       PRIMARY KEY (`id`),
                       UNIQUE KEY `idx_user_name` (`user_id`, `name`),
                       CONSTRAINT `fk_tag_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='标签表';

-- 项目标签关联表
CREATE TABLE `item_tag` (
                            `id` VARCHAR(36) NOT NULL COMMENT 'UUID',
                            `item_id` VARCHAR(36) NOT NULL COMMENT '项目ID',
                            `tag_id` VARCHAR(36) NOT NULL COMMENT '标签ID',
                            `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                            PRIMARY KEY (`id`),
                            UNIQUE KEY `idx_item_tag` (`item_id`, `tag_id`),
                            KEY `idx_tag_id` (`tag_id`),
                            CONSTRAINT `fk_it_item_id` FOREIGN KEY (`item_id`) REFERENCES `item` (`id`) ON DELETE CASCADE,
                            CONSTRAINT `fk_it_tag_id` FOREIGN KEY (`tag_id`) REFERENCES `tag` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='项目与标签关联表';

-- 用户事件表
CREATE TABLE `user_event` (
                              `id` VARCHAR(36) NOT NULL COMMENT 'UUID',
                              `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
                              `event_type` VARCHAR(50) NOT NULL COMMENT '事件类型',
                              `event_name` VARCHAR(100) NOT NULL COMMENT '事件名称',
                              `device_type` VARCHAR(20) NULL COMMENT '设备类型:web,mobile,app',
                              `platform` VARCHAR(20) NULL COMMENT '平台:ios,android,windows,macos,linux',
                              `browser` VARCHAR(50) NULL COMMENT '浏览器',
                              `ip_address` VARCHAR(50) NULL COMMENT 'IP地址',
                              `user_agent` VARCHAR(500) NULL COMMENT '用户代理',
                              `properties` JSON NULL COMMENT '事件属性',
                              `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                              PRIMARY KEY (`id`),
                              KEY `idx_user_id` (`user_id`),
                              KEY `idx_event_type` (`event_type`),
                              KEY `idx_created_at` (`created_at`),
                              CONSTRAINT `fk_event_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户事件记录表';

-- 页面访问表
CREATE TABLE `page_visit` (
                              `id` VARCHAR(36) NOT NULL COMMENT 'UUID',
                              `user_id` VARCHAR(36) NULL COMMENT '用户ID,未登录为null',
                              `session_id` VARCHAR(100) NOT NULL COMMENT '会话ID',
                              `page_url` VARCHAR(500) NOT NULL COMMENT '页面URL',
                              `page_title` VARCHAR(255) NULL COMMENT '页面标题',
                              `referrer` VARCHAR(500) NULL COMMENT '来源URL',
                              `device_type` VARCHAR(20) NULL COMMENT '设备类型',
                              `platform` VARCHAR(20) NULL COMMENT '平台',
                              `browser` VARCHAR(50) NULL COMMENT '浏览器',
                              `ip_address` VARCHAR(50) NULL COMMENT 'IP地址',
                              `entry_time` DATETIME NOT NULL COMMENT '进入时间',
                              `exit_time` DATETIME NULL COMMENT '离开时间',
                              `duration` INT UNSIGNED NULL COMMENT '停留时长(秒)',
                              `scroll_depth` INT UNSIGNED NULL COMMENT '滚动深度(%)',
                              `is_bounce` TINYINT(1) NULL COMMENT '是否跳出',
                              `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                              PRIMARY KEY (`id`),
                              KEY `idx_user_id` (`user_id`),
                              KEY `idx_session_id` (`session_id`),
                              KEY `idx_entry_time` (`entry_time`),
                              KEY `idx_page_url` (`page_url`(255)),
                              CONSTRAINT `fk_visit_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='页面访问记录表';

-- 功能使用表
CREATE TABLE `feature_usage` (
                                 `id` VARCHAR(36) NOT NULL COMMENT 'UUID',
                                 `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
                                 `feature_name` VARCHAR(100) NOT NULL COMMENT '功能名称',
                                 `usage_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '使用次数',
                                 `last_used_at` DATETIME NULL COMMENT '最后使用时间',
                                 `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                 `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                 PRIMARY KEY (`id`),
                                 UNIQUE KEY `idx_user_feature` (`user_id`, `feature_name`),
                                 CONSTRAINT `fk_feature_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='功能使用统计表';
create table `chat` (
                        `id` VARCHAR(36) NOT NULL COMMENT 'UUID',
                        `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',

                        `content` TEXT NOT NULL COMMENT '聊天内容',
                     `sender_type` VARCHAR(100) NOT NULL COMMENT '发送者类型',
                        `send_time` TIMESTAMP NOT NULL ,
                        `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                        `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                        PRIMARY KEY (`id`),
                        CONSTRAINT `fk_chat_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='聊天表';

