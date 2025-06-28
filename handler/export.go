// handlers/export.go
package handler

import (
	"makehugocontent/utils"
	"net/http"
)

// ExportCSVHandler streams content list as CSV for download.
func ExportCSVHandler(w http.ResponseWriter, r *http.Request) {
	if !checkLogin(w, r) { // 复用你的会话校验
		return
	}

	rows, err := utils.ScanDir(CONTENT_PATH)
	if err != nil {
		http.Error(w, "导出失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=\"content_export.csv\"")

	if err := utils.WriteCSV(rows, w); err != nil {
		http.Error(w, "CSV 写入失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
