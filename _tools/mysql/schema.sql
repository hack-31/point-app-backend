CREATE TABLE `users` (
    `id`                      BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ユーザーの識別子',
    `family_name`             VARCHAR(256) NOT NULL COMMENT '苗字',
    `family_name_kana`        VARCHAR(256) NOT NULL COMMENT '苗字カナ',
    `first_name`              VARCHAR(256) NOT NULL COMMENT '名前',
    `first_name_kana`         VARCHAR(256) NOT NULL COMMENT '名前カナ',
    `email`                   VARCHAR(256) NOT NULL COMMENT 'メールアドレス',
    `password`                VARCHAR(256) NOT NULL COMMENT 'パスワードハッシュ',
    `sending_point`           INT NOT NULL COMMENT '送信可能ポイント',
    `created_at`              DATETIME(6) NOT NULL COMMENT 'レコード作成日時',
    `update_at`               DATETIME(6) NOT NULL COMMENT 'レコード修正日時',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uix_email` (`email`) USING BTREE
) Engine=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='ユーザー';

CREATE TABLE `delete_users` (
    `id`                      BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ユーザーの識別子',
    `notification_latest_id`  BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '最新のお知らせID',
    `family_name`             VARCHAR(256) NOT NULL COMMENT '苗字',
    `family_name_kana`        VARCHAR(256) NOT NULL COMMENT '苗字カナ',
    `first_name`              VARCHAR(256) NOT NULL COMMENT '名前',
    `first_name_kana`         VARCHAR(256) NOT NULL COMMENT '名前カナ',
    `email`                   VARCHAR(256) NOT NULL COMMENT 'メールアドレス',
    `password`                VARCHAR(256) NOT NULL COMMENT 'パスワードハッシュ',
    `sending_point`           INT NOT NULL COMMENT '送信可能ポイント',
    `created_at`              DATETIME(6) NOT NULL COMMENT 'レコード作成日時',
    `update_at`               DATETIME(6) NOT NULL COMMENT 'レコード修正日時',
    `delete_at`               DATETIME(6) NOT NULL COMMENT 'レコード削除日時',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uix_email` (`email`) USING BTREE
) Engine=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='削除ユーザー';


CREATE TRIGGER `after_delete_users`
AFTER DELETE ON `users`
FOR EACH ROW
BEGIN
    INSERT INTO `delete_users` (`first_name`, `first_name_kana`, `family_name`, `family_name_kana`, `email`, `password`, `sending_point`, `created_at`, `update_at`, `delete_at`)
    VALUES (OLD.`first_name`, OLD.`first_name_kana`, OLD.`family_name`, OLD.`family_name_kana`, OLD.`email`, OLD.`password`, OLD.`sending_point`, OLD.`created_at`, OLD.`update_at`,  NOW());
END; 

CREATE TABLE `transactions` (
    `id`                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '取引の識別子',
    `sending_user_id`    BIGINT UNSIGNED NOT NULL COMMENT '送信ユーザのID',
    `receiving_user_id`  BIGINT UNSIGNED NOT NULL COMMENT '受信ユーザのID',
    `transaction_point`  INT NOT NULL COMMENT '取引ポイント',
    `transaction_at`     DATETIME(6) NOT NULL COMMENT '取引日時',
    PRIMARY KEY (`id`)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='取引';

CREATE TABLE `notification_types` (
    `id`                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'お知らせ種別ID',
    `title`             VARCHAR(256) NOT NULL COMMENT 'タイトル',
    PRIMARY KEY (`id`)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='お知らせ種別';

CREATE TABLE `notifications` (
    `id`                    BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'お知らせID',
    `to_user_id`            BIGINT UNSIGNED NOT NULL COMMENT 'お知らせ先ユーザID',
    `from_user_id`          BIGINT UNSIGNED NOT NULL COMMENT 'お知らせ元ユーザID',
    `is_checked`            BOOLEAN NOT NULL DEFAULT 0 COMMENT 'チェックフラグ',
    `notification_type_id`  BIGINT UNSIGNED NOT NULL COMMENT 'お知らせ種別ID',
    `description`           VARCHAR(256) NOT NULL COMMENT '説明',
    `created_at`            DATETIME(6) NOT NULL COMMENT '作成日時',
    PRIMARY KEY (`id`),
    INDEX `idx_to_user_id` (`to_user_id`, `id`),
    CONSTRAINT `fk_notification_type_id`
        FOREIGN KEY (`notification_type_id`) REFERENCES `notification_types` (`id`) ON UPDATE CASCADE
) Engine=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='お知らせ';
