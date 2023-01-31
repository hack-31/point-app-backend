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

CREATE TABLE `transactions` (
    `id`                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '取引の識別子',
    `sending_user_id`    BIGINT UNSIGNED NOT NULL COMMENT '送信ユーザのID',
    `receiving_user_id`  BIGINT UNSIGNED NOT NULL COMMENT '受信ユーザのID',
    `transaction_point`  INT NOT NULL COMMENT '取引ポイント',
    `transaction_at`     DATETIME(6) NOT NULL COMMENT '取引日時',
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_sending_user_id`
        FOREIGN KEY (`sending_user_id`) REFERENCES `users` (`id`)
            ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT `fk_receiving_user_id`
        FOREIGN KEY (`receiving_user_id`) REFERENCES `users` (`id`)
            ON DELETE RESTRICT ON UPDATE RESTRICT
) Engine=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='取引';
