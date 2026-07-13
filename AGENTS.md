# GEMINI.md

這個檔案提供了此專案（Teacher Wails / 班級值日與午餐排程工具）的上下文與開發指南，供 Gemini 讀取以進行高效的程式碼理解與修改。

## 專案概覽 (Project Overview)

- **專案名稱**：Teacher Helper (班級值日生助手)
- **技術棧**：
  - **後端 (Backend)**：Go + Wails v2
  - **前端 (Frontend)**：Svelte 4 + TypeScript + Vite (無外部 UI 庫，純 CSS 樣式)
- **資料儲存**：無資料庫，全部資料（學生、設定、假日、課表、缺交作業）以單一 JSON 檔案 `data/config.json` 進行持久化儲存。
- **核心功能**：
  1. **今日值日**：依學期起始日、扣除假日後的工作天數，自動計算值日生與午餐打菜人員。
  2. **手動更換**：臨時替換今日值日生/打菜生，不變動核心排程。
  3. **課表設定**：設定 5×8 的課表內容。
  4. **全螢幕展示**：投影模式，自動高亮當前課堂節次。
  5. **上課倒數**：上課前 1 分鐘自動開啟 60 秒倒數，可撥放背景 MP3 音樂。
  6. **國定假日同步**：串接政府開放資料 (data.gov.tw) 行事曆 API 自動匯入非週末假日。
  7. **匯出 CSV**：匯出未來 90 天的值日與午餐預排表。

---

## 常用開發指令 (Commands)

### Wails 開發與編譯
```bash
# 啟動熱重載開發模式 (同時啟動 Go 後端與 Svelte 前端)
wails dev

# 編譯成發佈版本 (輸出至 build/bin/)
wails build
```

### 前端獨立測試 (在 frontend/ 目錄下)
```bash
npm install        # 安裝套件
npm run dev        # 單獨啟動 Vite 開發伺服器
npm run build      # 打包前端資源
```

### 後端測試 (Go)
```bash
# 執行所有服務層測試
go test ./internal/services/...

# 執行特定測試
go test ./internal/services/ -run <TestName>
```

---

## 專案架構 (Project Architecture)

### 後端架構 (`internal/`)
- `models/models.go` — 定義所有 JSON 與 API 回傳的資料結構（如 `Student`, `Settings`, `Config`, `TodayDutyResult`）。
- `services/data_handler.go` — 負責載入與寫入 `data/config.json`，在啟動時會優先搜尋執行檔路徑下的 `data/` 資料夾，若不存在則退回當前工作目錄。
- `services/duty_calculator.go` — 排程核心邏輯：
  - **值日生**：以工作天數為基礎進行日輪替 (`calculateGroupRotation`，模式為 `"daily"`)。
  - **午餐打菜**：以 ISO 週數為基礎進行週輪替 (模式為 `"weekly"`)，打菜桶 (Meal Buckets) 的分配在當天以 `daySeed` (工作天數) 進行隨機打散。
- `services/student_manager.go` — 處理學生的增刪改查、切換值日/午餐參與狀態。
- `services/schedule_exporter.go` — 負責將未來的排程資料轉換成 CSV 字串。
- `utils/date_utils.go` — 包含日期的解析、計算工作天、判斷是否為工作天等工具函式。

### Wails 綁定層 (`app.go`, `main.go`)
- `main.go` — 應用程式進入點，配置視窗寬高與綁定 `App` 實體。
- `app.go` — 前後端溝通的主要橋樑。定義了所有匯出給前端呼叫的 Go 方法，並啟動了 `countdownWatcher` 背景執行緒監聽上課時間。

### 前端架構 (`frontend/src/`)
- `App.svelte` — 根元件與頁面路由控制（`currentPage` 狀態）。
- `lib/pages/` — 各個主要頁面（例如首頁、學生管理、設定頁、展示頁、課表頁等）。
- `lib/components/` — 抽離出來的 UI 元件。
- `style.css` — 全域與元件使用的純 CSS 樣式表。

---

## 開發與修改規範 (Development Rules)

1. **前後端呼叫**：
   - 前端不能使用 `fetch` 或 `axios` 等 HTTP 請求，必須使用 Wails 自動產生的 `frontend/wailsjs/go/main/App.js` 模組直接呼叫 Go 函式。
   - 所有後端新增的方法，在執行 `wails dev` 或 `wails build` 時會自動在 `frontend/wailsjs/` 產生 TypeScript 的型別定義與 API 綁定檔。
2. **資料持久化**：
   - 每次修改學生、設定、課表等，都必須透過 `data_handler` 儲存回 `config.json`，以便下次開啟時能回復狀態。
3. **時間處理**：
   - 一律使用 `utils/date_utils.go` 的工具來進行日期轉換，並遵循 `YYYY-MM-DD` 格式以避免時區與格式解析錯誤。
