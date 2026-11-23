package main

import (
	"context"
	"database/sql"
	"errors"

	"github.com/siyad01/multiplayer-terminal-games/internal/minesweeper"
)

// saveMatchResult 在同一个事务中保存一局对局及其玩家记录。
func saveMatchResult(ctx context.Context, db *sql.DB, result minesweeper.MatchResult) error {
	// 使用 ctx 开启事务。
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// 提前返回时回滚事务。
	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO match_results (game_id, won, canceled) VALUES (?, ?, ?)",
		result.GameID, result.Won, result.Cancel,
	)
	if err != nil {
		return err
	}

	// 遍历 result.PlayerResult，为每位玩家向 match_players 写入一行。
	for _, p := range result.PlayerResult {
		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO match_players(game_id,player_id,name,reveals,flags,left_early) VALUES(?,?,?,?,?,?)",
			result.GameID, p.ID, p.Name, p.Reveals, p.Flags, p.Left,
		)
		// 沿用同一个事务和对局编号；写入失败就返回错误。
		if err != nil {
			return err
		}
	}

	// 等所有写入都成功后再提交事务，并返回提交结果。
	return tx.Commit()

}

// loadMatchResult 按对局编号读取一局战绩及其玩家记录。
func loadMatchResult(ctx context.Context, db *sql.DB, gameID string) (minesweeper.MatchResult, error) {
	var result minesweeper.MatchResult
	// 根据 gameID 查 match_results 的一行，用 Scan 填入 result 的对局字段。
	const matchQuery = `SELECT game_id, won, canceled FROM match_results WHERE game_id = ?`
	err := db.QueryRowContext(ctx, matchQuery, gameID).Scan(
		&result.GameID,
		&result.Won,
		&result.Cancel,
	)
	// 查不到对局时返回 sql.ErrNoRows；其他查询错误也直接返回。
	if err != nil {
		return result, err
	}

	// 根据同一个 gameID 查 match_players，按 player_id 升序返回。
	const playersQuery = `
		SELECT player_id, name, reveals, flags, left_early
		FROM match_players
		WHERE game_id = ?
		ORDER BY player_id ASC`
	rows, err := db.QueryContext(ctx, playersQuery, gameID)
	if err != nil {
		return result, err
	}
	// 取得 rows 后安排关闭；把每一行读入一份 MinePlayer，再追加到 result.PlayerResult。
	defer rows.Close()

	for rows.Next() {
		var p minesweeper.MinePlayer
		if err := rows.Scan(&p.ID, &p.Name, &p.Reveals, &p.Flags, &p.Left); err != nil {
			return result, err
		}
		result.PlayerResult = append(result.PlayerResult, p)
	}

	// 检查遍历过程中是否出错，确认无误后返回完整结果。
	if err := rows.Err(); err != nil {
		return result, err
	}
	return result, nil
}

// matchSummary 是列表中一局对局的摘要，不包含玩家详情。
type matchSummary struct {
	GameID string
	Won    bool
	Cancel bool
}

const maxMatchPageSize = 50

// listMatchResults 读取一页对局摘要。limit 是条数，offset 是从开头跳过的条数。
func listMatchResults(ctx context.Context, db *sql.DB, limit, offset int) ([]matchSummary, error) {
	// limit 不在 1 到 maxMatchPageSize，或 offset 小于 0 时，返回错误。
	if limit < 1 || limit > maxMatchPageSize || offset < 0 {
		return nil, errors.New("分页参数不合法")
	}
	// 用两个 SQL 参数限制条数、跳过前面的行，不拼接参数。
	const listQuery = `
    SELECT game_id,won,canceled
    FROM match_results
    ORDER BY game_id ASC
    LIMIT ? OFFSET ?`
	rows, err := db.QueryContext(ctx, listQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	// ：只查 match_results 的编号和两个状态字段，按唯一主键升序排列。
	matches := make([]matchSummary, 0)
	// 参照 loadMatchResult 的玩家查询，关闭 rows、逐行 Scan 并检查 rows.Err()。
	for rows.Next() {
		var match matchSummary
		if err := rows.Scan(&match.GameID, &match.Won, &match.Cancel); err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return matches, nil

	// 合法但超出最后一页时，返回空切片和 nil，不返回 sql.ErrNoRows。
}
