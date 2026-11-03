create table user_info(
    id bigint primary key auto_increment,
    code varchar(255) comment 'user_code',
    birthday date comment '出生年月',
    gender tinyint unsigned DEFAULT 0 COMMENT '性别:0-未知,1-男,2-女',
    nickname varchar(255) DEFAULT '' COMMENT '昵称',
    username varchar(255) DEFAULT null COMMENT 'username',
    country_code int DEFAULT 84 COMMENT '国家码',
    is_logoff tinyint unsigned DEFAULT 0 COMMENT '是否注销:0-未注销,1-已注销',
    is_auth tinyint unsigned DEFAULT 0 COMMENT '是否认证:0-未认证,1-已认证',
    logoff_time bigint null default 0 COMMENT '注销时间',
    black_time bigint null default 0 COMMENT '平台拉黑/封禁时间',
    wealth_level int default 0 comment '财富等级',
    charm_level int default 0 comment '魅力等级',
    created_at timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP on update current_timestamp COMMENT '更新时间',
    unique key `udx_code` (`code`),
    unique key `udx_username` (`username`)
)default char set = `utf8mb4` comment '用户详情表';

create table register_info(
    id bigint primary key auto_increment comment '自增id',
    user_id bigint not null default 0 comment '用户id',
    phone varchar(255) default null comment '手机号',
    email varchar(255) default null comment '邮箱',
    status tinyint unsigned default 0 comment '注册状态 0正常,1封禁,2删除',
    created_at timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP on update current_timestamp COMMENT '更新时间',
    unique key `udx_phone` (phone),
    unique key `udx_email` (email),
    UNIQUE KEY `udx_user_id` (`user_id`)
)default char set = `utf8mb4` comment '注册信息表';

create table auth_code(
    id bigint primary key auto_increment comment '自增id',
    number varchar(255) NOT NULL comment '手机号/邮箱',
    code varchar(5) NOT NULL default '' comment '验证码',
    code_type tinyint unsigned not null default 0 comment '验证码发送场景:1-登录,2-注册,3,4,5,6,7',
    send_time BIGINT not null default 0 comment '发送时间',
    expires_time BIGINT NOT NULL DEFAULT 0 COMMENT '过期时间',
    deleted tinyint unsigned not null default 0 comment '是否删除:0-未删除,1-已删除',
    created_at timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP on update current_timestamp COMMENT '更新时间',
    UNIQUE KEY `udx_number_type` (`number`, `code_type`),
    key `idx_number` (`number`,deleted,expires_time)
)default char set = `utf8mb4` comment '验证码';

CREATE TABLE `auth_login` (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '自增ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    ip VARCHAR(45) NOT NULL DEFAULT '' COMMENT '登录IP',
    device_id VARCHAR(128) NOT NULL DEFAULT '' COMMENT '设备唯一标识',
    device_model VARCHAR(128) NOT NULL DEFAULT '' COMMENT '手机型号',
    system_version VARCHAR(64) NOT NULL DEFAULT '' COMMENT '手机系统版本',
    app_version VARCHAR(64) NOT NULL DEFAULT '' COMMENT '应用版本',
    platform TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '设备平台:0-未知,1-Android,2-iOS,3-Web',
    language VARCHAR(16) NOT NULL DEFAULT '' COMMENT '语言代码:如zh-CN,en-US,vi-VN',
    login_type TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '登录方式:0-未知,1-密码,2-手机验证码,3-邮箱验证码',
    action_type TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '操作类型:1-登录,2-登出',
    active_time BIGINT NOT NULL DEFAULT 0 COMMENT '最后活跃时间戳',
    country_code INT NOT NULL DEFAULT 0 COMMENT '国家区号',
    deleted TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否删除:0-未删除,1-已删除',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    KEY `idx_user_id_deleted_created` (`user_id`, `deleted`, `created_at`)
) DEFAULT CHARSET = utf8mb4 COMMENT '用户登录登出记录表';