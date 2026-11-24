UPDATE user_info SET is_logoff = 0 WHERE is_logoff IS NULL;
UPDATE register_info SET status = 0 WHERE status IS NULL;

ALTER TABLE user_info
    MODIFY COLUMN is_logoff TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否注销:0-未注销,1-已注销',
    ADD INDEX idx_user_phone_logoff (phone, is_logoff),
    ADD INDEX idx_user_email_logoff (email, is_logoff);

ALTER TABLE register_info
    MODIFY COLUMN status TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '注册状态 0正常,1封禁,2删除',
    DROP INDEX udx_phone,
    DROP INDEX udx_email,
    ADD UNIQUE INDEX uk_register_active_phone ((CASE WHEN status <> 2 THEN phone ELSE NULL END)),
    ADD UNIQUE INDEX uk_register_active_email ((CASE WHEN status <> 2 THEN email ELSE NULL END));

UPDATE register_info AS registration
SET registration.status = 2
WHERE registration.status <> 2
  AND (
      (registration.phone IS NOT NULL AND registration.email IS NULL
       AND EXISTS (SELECT 1 FROM user_info AS old_user
                   WHERE old_user.phone = registration.phone AND old_user.is_logoff = 1)
       AND NOT EXISTS (SELECT 1 FROM user_info AS active_user
                       WHERE active_user.phone = registration.phone AND active_user.is_logoff = 0))
      OR
      (registration.email IS NOT NULL AND registration.phone IS NULL
       AND EXISTS (SELECT 1 FROM user_info AS old_user
                   WHERE old_user.email = registration.email AND old_user.is_logoff = 1)
       AND NOT EXISTS (SELECT 1 FROM user_info AS active_user
                       WHERE active_user.email = registration.email AND active_user.is_logoff = 0))
  );
