package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func matchDetailHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := r.PathValue("gameID")

		if gameID == "" || len(gameID) > 64 {
			http.Error(w, "对局编号无效", http.StatusBadRequest)
			return
		}

		result, err := loadMatchResult(r.Context(), db, gameID)

		if err == sql.ErrNoRows {
			http.Error(w, "战绩不存在", http.StatusNotFound)
			return
		}
		if err != nil {
			// TODO(A3)：在服务端日志里记录查询失败的对局编号和原始错误。
			http.Error(w, "查询战绩失败", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Printf("写入战绩详情响应失败：%v", err)
		}
	}
}

func matchListHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从 URL 的查询参数取出 limit 和 offset，用 strconv.Atoi 转成整数。
		text := r.URL.Query().Get("limit") // 读到字符串 "2"
		count, err := strconv.Atoi(text)   // 转成整数 2，err 表示转换是否失败
		// 缺少参数、转换失败，或 limit 不在 1 到 maxMatchPageSize、offset 小于 0 时，返回 400。
		if err != nil || count < 1 || count > maxMatchPageSize {
			http.Error(w, "分页参数不合法", http.StatusBadRequest)
			return
		}
		otext := r.URL.Query().Get("offset")
		offset, offsetErr := strconv.Atoi(otext)

		if offsetErr != nil || offset < 0 {
			http.Error(w, "跳过参数不合法", http.StatusBadRequest)
			return
		}
		// 用请求的上下文调用 listMatchResults。查询失败时记录错误并返回 500。
		page, perr := listMatchResults(r.Context(), db, count, offset)
		if perr != nil {
			http.Error(w, "查询出错", http.StatusInternalServerError)
			return
		}

		// 复用详情的 JSON 响应写法，返回这一页摘要。
		http.Error(w, "战绩分页接口尚未实现", http.StatusNotImplemented)
	}
}
