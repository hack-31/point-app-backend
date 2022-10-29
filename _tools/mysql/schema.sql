CREATE TABLE `roles`
(
    `id`       BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ロールの識別子',
    `name`     VARCHAR(20)  NOT NULL COMMENT 'ロール名',
    PRIMARY KEY (`id`)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='ロール';


CREATE TABLE `users`
(
    `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ユーザーの識別子',
    `name`        VARCHAR(20) NOT NULL COMMENT 'ユーザー名',
    `email`       VARCHAR(20) NOT NULL COMMENT 'メールアドレス',
    `password`    VARCHAR(80) NOT NULL COMMENT 'パスワードハッシュ',
    `role_id`     BIGINT UNSIGNED NOT NULL COMMENT 'ロール',
    `created_at`  DATETIME(6) NOT NULL COMMENT 'レコード作成日時',
    `modified_at` DATETIME(6) NOT NULL COMMENT 'レコード修正日時',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uix_email` (`email`) USING BTREE,
    CONSTRAINT `fk_role`
        FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`)
            ON DELETE RESTRICT ON UPDATE RESTRICT
) Engine=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='ユーザー';

