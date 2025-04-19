-- 创建主表 video_contents
CREATE TABLE video_contents (
                                id INT AUTO_INCREMENT PRIMARY KEY,
                                user_id VARCHAR(50) NOT NULL,
                                type ENUM('VIDEO', 'IMAGE', 'TEXT') NOT NULL,
                                url VARCHAR(255) NOT NULL,
                                title VARCHAR(255) NOT NULL,
                                description TEXT,
                                thumbnail_url VARCHAR(255),
                                is_favorite BOOLEAN DEFAULT FALSE,
                                is_private BOOLEAN DEFAULT FALSE,
                                note TEXT,
                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 创建元数据表 content_metadata
CREATE TABLE content_metadata (
                                  content_id INT PRIMARY KEY,
                                  title VARCHAR(255),
                                  description TEXT,
                                  thumbnail_url VARCHAR(255),
                                  author VARCHAR(100),
                                  site_name VARCHAR(100),
                                  favicon_url VARCHAR(255),
                                  content_type VARCHAR(100),
                                  language VARCHAR(50),
                                  is_article BOOLEAN DEFAULT FALSE,
                                  FOREIGN KEY (content_id) REFERENCES video_contents(id) ON DELETE CASCADE
);

-- 创建关键字表 content_keywords
CREATE TABLE content_keywords (
                                  id INT AUTO_INCREMENT PRIMARY KEY,
                                  content_id INT,
                                  keyword VARCHAR(100),
                                  FOREIGN KEY (content_id) REFERENCES video_contents(id) ON DELETE CASCADE
);

-- 创建标签表 content_tags
CREATE TABLE content_tags (
                              id INT AUTO_INCREMENT PRIMARY KEY,
                              content_id INT,
                              tag VARCHAR(100),
                              FOREIGN KEY (content_id) REFERENCES video_contents(id) ON DELETE CASCADE
);

-- 创建收藏集关联表 content_collections
CREATE TABLE content_collections (
                                     id INT AUTO_INCREMENT PRIMARY KEY,
                                     content_id INT,
                                     collection_id VARCHAR(50),
                                     FOREIGN KEY (content_id) REFERENCES video_contents(id) ON DELETE CASCADE
);

-- 插入主表数据
INSERT INTO video_contents (
    user_id, type, url, title, description, thumbnail_url,
    is_favorite, is_private, note
) VALUES (
             '74', 'VIDEO', 'https://livid-brace.net/', '游戏厅与此同时许多开做绿',
             '些认接由。此场王不他用更后受广。几听报量格准于马。数第阶。近速引系将格。',
             'https://hopeful-deck.net/', TRUE, TRUE, 'aute nisi anim'
         );

-- 获取最后插入的ID
SET @last_content_id = LAST_INSERT_ID();

-- 插入元数据
INSERT INTO content_metadata (
    content_id, title, description, thumbnail_url, author,
    site_name, favicon_url, content_type, language, is_article
) VALUES (
             @last_content_id, '僧袍百般在绿油油僧袍鳄鱼别靠近',
             '力更具活指次科科六。题量少支老道。多难一关很候质场。',
             'https://giving-legislature.org/', 'laboris eu dolor in',
             '恭燕', 'https://avatars.githubusercontent.com/u/70806301',
             'commodo laboris proident', 'cillum minim Excepteur', FALSE
         );

-- 插入关键字
INSERT INTO content_keywords (content_id, keyword) VALUES
                                                       (@last_content_id, 'tempor dolor esse'),
                                                       (@last_content_id, 'ullamco labore');

-- 插入标签
INSERT INTO content_tags (content_id, tag) VALUES
                                               (@last_content_id, 'quis sed nisi magna'),
                                               (@last_content_id, 'eiusmod magna');

-- 插入收藏集关联
INSERT INTO content_collections (content_id, collection_id) VALUES
                                                                (@last_content_id, '71'),
                                                                (@last_content_id, '70'),
                                                                (@last_content_id, '44');