-- 用户表
CREATE TABLE `user` (
                        `id` VARCHAR(36) NOT NULL COMMENT 'UUID',
                        `email` VARCHAR(255) NOT NULL COMMENT '邮箱',
                        `phone` VARCHAR(20) NULL COMMENT '电话',
                        `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希',
                        `username` VARCHAR(50) NOT NULL COMMENT '用户名',
                        `avatar_url` VARCHAR(255) NULL COMMENT '头像URL',
                        `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-正常，0-禁用',
                        `last_login_at` DATETIME NULL COMMENT '最后登录时间',
                        `login_count` INT UNSIGNED DEFAULT 0 COMMENT '登录次数',
                        `registration_source` VARCHAR(20) DEFAULT 'email' COMMENT '注册来源: email,phone,wechat,github等',
                        `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                        `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `idx_email` (`email`),
                        UNIQUE KEY `idx_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户基础信息表';

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

-- 收藏夹表
CREATE TABLE `collection` (
                              `id` VARCHAR(36) NOT NULL COMMENT 'UUID',
                              `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
                              `parent_id` VARCHAR(36) NULL COMMENT '父收藏夹ID',
                              `name` VARCHAR(100) NOT NULL COMMENT '收藏夹名称',
                              `description` VARCHAR(500) NULL COMMENT '描述',
                              `is_default` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否默认收藏夹',
                              `is_shared` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否共享',
                              `share_code` VARCHAR(32) NULL COMMENT '分享码',
                              `share_expire_at` DATETIME NULL COMMENT '分享过期时间',
                              `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序顺序',
                              `items_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '收藏项目数',
                              `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                              `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                              PRIMARY KEY (`id`),
                              KEY `idx_user_parent` (`user_id`, `parent_id`),
                              KEY `idx_share_code` (`share_code`),
                              CONSTRAINT `fk_collection_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE,
                              CONSTRAINT `fk_collection_parent_id` FOREIGN KEY (`parent_id`) REFERENCES `collection` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户收藏夹表';

-- 收藏项目表
CREATE TABLE `item` (
                        `id` VARCHAR(36) NOT NULL COMMENT 'UUID',
                        `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
                        `item_type` VARCHAR(20) NOT NULL COMMENT '类型:text,image,link,bookmark',
                        `title` VARCHAR(500) NULL COMMENT '标题',
                        `content` TEXT NULL COMMENT '内容/文本',
                        `url` VARCHAR(2000) NULL COMMENT '链接地址',
                        `image_url` VARCHAR(2000) NULL COMMENT '图片地址',
                        `thumbnail_url` VARCHAR(2000) NULL COMMENT '缩略图地址',
                        `source_domain` VARCHAR(255) NULL COMMENT '来源网站域名',
                        `source_page_title` VARCHAR(500) NULL COMMENT '来源页面标题',
                        `is_favorited` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否收藏',
                        `favorited_at` DATETIME NULL COMMENT '收藏时间',
                        `is_deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否删除',
                        `deleted_at` DATETIME NULL COMMENT '删除时间',
                        `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                        `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                        PRIMARY KEY (`id`),
                        KEY `idx_user_id` (`user_id`),
                        KEY `idx_type` (`item_type`),
                        KEY `idx_user_deleted` (`user_id`, `is_deleted`),
                        KEY `idx_deleted_at` (`deleted_at`),
                        KEY `idx_domain` (`source_domain`),
                        CONSTRAINT `fk_item_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='收藏项目表';

-- 收藏夹项目关联表
CREATE TABLE `collection_item` (
                                   `id` VARCHAR(36) NOT NULL COMMENT 'UUID',
                                   `collection_id` VARCHAR(36) NOT NULL COMMENT '收藏夹ID',
                                   `item_id` VARCHAR(36) NOT NULL COMMENT '项目ID',
                                   `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序顺序',
                                   `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                   PRIMARY KEY (`id`),
                                   UNIQUE KEY `idx_collection_item` (`collection_id`, `item_id`),
                                   KEY `idx_item_id` (`item_id`),
                                   CONSTRAINT `fk_ci_collection_id` FOREIGN KEY (`collection_id`) REFERENCES `collection` (`id`) ON DELETE CASCADE,
                                   CONSTRAINT `fk_ci_item_id` FOREIGN KEY (`item_id`) REFERENCES `item` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='收藏夹与项目关联表';

-- 标签表
CREATE TABLE `tag` (
                       `id` VARCHAR(36) NOT NULL COMMENT 'UUID',
                       `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
                       `name` VARCHAR(100) NOT NULL COMMENT '标签名称',
                       `is_system` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否系统标签',
                       `color` VARCHAR(20) NULL COMMENT '标签颜色',
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
