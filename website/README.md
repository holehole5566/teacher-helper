# Teacher Helper 介紹網站

正式網址：<https://teacher-helper.pe4nut.com/>

使用 Cloudflare Workers Static Assets 託管純 HTML / CSS，不執行 Worker 後端程式、不使用資料庫，也不讀取桌面應用程式的班級資料。僅 `public/` 會被部署；不會上傳專案的 `data/`、備份或桌面程式設定。

## 本機預覽

```bash
cd website
npm ci
npm run dev
```

開啟 Wrangler 印出的本機網址。頁面不需要建置，也不需要前端 JavaScript。

## 部署

```bash
cd website
npx wrangler login      # 第一次使用時授權，已有登入可略過
npx wrangler whoami     # 確認是管理 pe4nut.com 的帳號
npm run check          # dry run，不修改線上資源
npm run deploy
```

`wrangler.jsonc` 設定 Custom Domain 為 `teacher-helper.pe4nut.com`。`pe4nut.com` 必須是部署帳號內的有效 Cloudflare zone；Cloudflare 會管理該子網域的 DNS 與 TLS 憑證。若已有同名 DNS／服務，先確認用途，不要直接覆寫或刪除。

此網站與桌面版本的 Release CI 相互獨立，不需要為文案修改新增桌面版本 tag。

## GitHub Actions CI/CD

設定檔：`.github/workflows/website.yml`，與桌面 Release workflow 分開。

- PR 到 `master`：網站或此 workflow 有修改時，執行 `npm ci` 和 `npm run check`，不部署、不使用 Cloudflare 憑證。
- Push 到 `master`：相同路徑有修改時，先驗證，通過後執行 `npm run deploy`。
- 手動部署：GitHub → Actions → Website CI/CD → Run workflow，選 `master`。其他分支只驗證、不部署。
- 部署工作序列化，避免同時上傳。只有 `public/` 會部署，無需建置桌面程式。

### 首次設定憑證

本機 `cf` / Wrangler 的 OAuth 登入不會傳給 GitHub runner，也不應複製短效 OAuth token 當 CI 憑證。

1. 在 Cloudflare 的 API Tokens 頁建立專用 token，可從 **Edit Cloudflare Workers** 模板開始，將帳號限制為部署帳號、zone 限制為 `pe4nut.com`。此網站使用 Workers Static Assets 與 Custom Domain，保留模板中部署 Workers 和管理網域所需的權限，不需另外開 R2、D1 或資料庫權限。
2. 在 GitHub repository → Settings → Secrets and variables → Actions 加入：
   - `CLOUDFLARE_API_TOKEN`：上述 token。
   - `CLOUDFLARE_ACCOUNT_ID`：Cloudflare 帳號 ID（不是 zone ID）。
3. 提交並 push workflow 和網站修改，或在 Actions 手動啟動。

也可使用已登入的 GitHub CLI，以下指令會互動詢問值，勿將 token 貼進命令參數或聊天：

```bash
gh secret set CLOUDFLARE_API_TOKEN --repo holehole5566/teacher-helper
gh secret set CLOUDFLARE_ACCOUNT_ID --repo holehole5566/teacher-helper
```

使用 GitHub Actions 後，請勿同時啟用相同 Worker 的 Cloudflare Workers Builds Git 自動部署，以免同一次 push 重複部署。手動 `npm run deploy` 仍可使用。請勿將 OAuth token、API token 或 `.env` 提交到 Git。

參考：<https://developers.cloudflare.com/workers/ci-cd/external-cicd/github-actions/>

## 維護

- `public/index.html`：繁體中文介紹、功能、FAQ、下載入口。
- `public/styles.css`：響應式版型與鍵盤焦點樣式。
- `public/_headers`：CSP 與安全標頭。
- `public/404.html`：不存在頁面的錯誤畫面。
- `public/robots.txt`、`public/sitemap.xml`：搜尋引擎入口。
- 下載按鈕連到 GitHub Releases 的 `latest/download`，不寫死版本。
- `public/screenshots/`：八個功能的實際 Svelte 前端截圖，使用虛構名單、課表與作業資料；不是原生桌面視窗截圖。
- 首屏以文字與日常使用情境介紹；八張功能卡預設不顯示截圖，點選「查看真實畫面」後以原生 details 展開大圖，可收合或另開原尺寸圖片，不需 JavaScript。

## 更新功能截圖

啟動 `frontend/` 的 `npm run dev`，安裝 Playwright 後執行 `node website/scripts/capture.mjs`（從專案根目錄）。可用 `PLAYWRIGHT_MODULE` 指定 Playwright 模組路徑、`CHROME_PATH` 指定 Chrome 執行檔、`APP_URL` 指定 Vite 網址。

腳本只在獨立瀏覽器注入虛構 Wails 回傳資料，使用原本的前端元件與樣式，不啟動 Go 後端、不讀寫 `data/config.json`。腳本與示範資料不在 `public/` 內，不會部署；更新圖片後請檢查文字、捲動區域與版面。

如果 Release 附件名稱改變，請同步更新 `index.html` 的下載連結。
