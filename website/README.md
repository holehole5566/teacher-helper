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

## 可選：Cloudflare Git 自動部署

在 Cloudflare 的 Workers & Pages 中，將此 Worker 的 Builds 連接到 GitHub 儲存庫 `holehole5566/teacher-helper`：

- Production branch：`master`
- Root directory：`website`
- Build command：留空（純靜態檔）
- Deploy command：`npm run deploy`
- 如可設定 watch paths，限制為 `website/**`

連接儲存庫需要在 Cloudflare 控制台授權 GitHub；本機 CLI 部署不會自動建立這個連接。請勿將 OAuth token、API token 或 `.env` 提交到 Git。

## 維護

- `public/index.html`：繁體中文介紹、功能、FAQ、下載入口。
- `public/styles.css`：響應式版型與鍵盤焦點樣式。
- `public/_headers`：CSP 與安全標頭。
- `public/404.html`：不存在頁面的錯誤畫面。
- `public/robots.txt`、`public/sitemap.xml`：搜尋引擎入口。
- 下載按鈕連到 GitHub Releases 的 `latest/download`，不寫死版本。
- 畫面為手工製作的概念示意，不是桌面程式實際截圖；所有姓名皆為虛構。

如果 Release 附件名稱改變，請同步更新 `index.html` 的下載連結。
