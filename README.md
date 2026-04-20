# Teacher Wails — 班級值日與午餐排程工具

使用 [Wails](https://wails.io/)（Go + Svelte）開發的桌面應用程式，幫助老師自動排定學生的值日生輪值與午餐打菜分配。

## 功能

- **今日值日** — 根據學期開始日期與輪值設定，自動計算當天的值日生與午餐打菜人員
- **學生管理** — 新增、刪除學生，可個別開關值日 / 午餐參與
- **設定** — 學期起始日、值日人數、午餐人數、打菜桶名稱、起始座號等
- **假日管理** — 手動新增假日，或一鍵從政府開放資料（data.gov.tw）同步當年度行事曆
- **匯出 CSV** — 將未來 90 天的排程表匯出為 CSV 檔案

## 專案結構

```
teacher-wails/
├── app.go                  # 主要後端邏輯（Wails 綁定）
├── main.go                 # 程式進入點
├── internal/
│   ├── models/             # 資料模型
│   ├── services/           # 業務邏輯（排程計算、資料存取、匯出）
│   └── utils/              # 日期工具
├── frontend/               # Svelte 前端
│   └── src/lib/
│       ├── pages/          # 頁面元件
│       └── components/     # UI 元件
├── data/                   # 設定檔（開發用）
│   └── config.json
└── build/
    └── bin/                # 編譯輸出
        ├── teacher-wails.exe
        └── data/           # 設定檔（正式用，需與 exe 同目錄）
            └── config.json
```

## 開發

需要先安裝 [Wails CLI](https://wails.io/docs/gettingstarted/installation)。

```bash
# 開發模式（熱重載）
wails dev

# 編譯
wails build
```

## 範例資料

`data/config.json` 內附的學生名單為動漫角色，僅供示範用途，非真實人物。

## data 目錄說明

應用程式啟動時會依以下順序尋找 `data/config.json`：

1. exe 所在目錄下的 `data/`
2. 若找不到，退回到當前工作目錄下的 `data/`

發佈時請確保 `data/` 資料夾與 `teacher-wails.exe` 放在同一個目錄下。
