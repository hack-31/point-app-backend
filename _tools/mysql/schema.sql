CREATE TABLE `users` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ユーザーの識別子',
    `family_name`      VARCHAR(256) NOT NULL COMMENT '苗字',
    `family_name_kana`  VARCHAR(256) NOT NULL COMMENT '苗字カナ',
    `first_name`       VARCHAR(256) NOT NULL COMMENT '名前',
    `first_name_kana`   VARCHAR(256) NOT NULL COMMENT '名前カナ',
    `email`           VARCHAR(256) NOT NULL COMMENT 'メールアドレス',
    `password`        VARCHAR(256) NOT NULL COMMENT 'パスワードハッシュ',
    `sending_point`   INT NOT NULL COMMENT '送信可能ポイント',
    `created_at`      DATETIME(6) NOT NULL COMMENT 'レコード作成日時',
    `update_at`       DATETIME(6) NOT NULL COMMENT 'レコード修正日時',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uix_email` (`email`) USING BTREE
) Engine=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='ユーザー';

