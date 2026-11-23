-- A3: 在空数据库按下列顺序建表；已有表的本机库只执行尚未创建的语句。

CREATE TABLE match_results (
    game_id VARCHAR(64) NOT NULL,
    won BOOLEAN NOT NULL,
    canceled BOOLEAN NOT NULL,
    PRIMARY KEY (game_id)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_bin;

CREATE TABLE match_players (
    game_id VARCHAR(64) NOT NULL,
    player_id TINYINT UNSIGNED NOT NULL,
    name VARCHAR(15) NOT NULL,
    reveals SMALLINT UNSIGNED NOT NULL,
    flags SMALLINT UNSIGNED NOT NULL,
    left_early BOOLEAN NOT NULL,
    PRIMARY KEY (game_id, player_id),
    FOREIGN KEY (game_id)
        REFERENCES match_results (game_id)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_bin;
