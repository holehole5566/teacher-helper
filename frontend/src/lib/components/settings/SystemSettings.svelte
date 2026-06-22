<script lang="ts">
  import { onMount } from 'svelte';
  import { GetSettings, SaveSettings, GetCurrentVersion, CheckForUpdate, DoUpdate } from '../../../../wailsjs/go/main/App';
  import { verifyAndRun } from '../../stores/auth';

  let discordWebhook = '';
  let saved = false;

  let currentVersion = '';
  let updateStatus: 'idle' | 'checking' | 'available' | 'updating' | 'upToDate' | 'error' = 'idle';
  let latestVersion = '';
  let updateError = '';

  async function loadSettings() {
    const s = await GetSettings();
    discordWebhook = s.discord_webhook || '';
    currentVersion = await GetCurrentVersion();
  }

  async function handleSave() {
    try {
      const current = await GetSettings();
      await verifyAndRun(async () => {
        await SaveSettings({
          ...current,
          discord_webhook: discordWebhook,
        });
      }, '儲存系統設定');
      saved = true;
      setTimeout(() => { saved = false; }, 2000);
    } catch (e: any) {
      if (e.message !== '驗證失敗或已取消') {
        alert('儲存失敗: ' + e.message);
      }
    }
  }

  async function checkUpdate() {
    updateStatus = 'checking';
    updateError = '';
    try {
      const result = await CheckForUpdate();
      if (result.has_update) {
        updateStatus = 'available';
        latestVersion = result.latest_version;
      } else {
        updateStatus = 'upToDate';
      }
    } catch (e: any) {
      updateStatus = 'error';
      updateError = e?.message || '檢查更新失敗';
    }
  }

  async function doUpdate() {
    updateStatus = 'updating';
    try {
      await DoUpdate();
    } catch (e: any) {
      updateStatus = 'error';
      updateError = e?.message || '更新失敗';
    }
  }

  onMount(loadSettings);
</script>

<div class="settings-section">
  <h3 class="section-title">系統設定</h3>
  <p class="section-desc">通知與軟體更新</p>

  <div class="form-group">
    <label>Discord Webhook（錯誤通知用）</label>
    <input type="text" bind:value={discordWebhook} placeholder="https://discord.com/api/webhooks/..." />
  </div>

  <div class="form-actions">
    <button class="btn-primary save-btn" on:click={handleSave}>儲存設定</button>
    {#if saved}
      <span class="save-ok">已儲存</span>
    {/if}
  </div>

  <div class="divider"></div>

  <div class="update-section">
    <h4 class="subsection-title">軟體更新</h4>
    <div class="update-card">
      <div class="update-info">
        <span class="version-label">目前版本</span>
        <span class="version-value">{currentVersion || '...'}</span>
      </div>
      <div class="update-actions">
        {#if updateStatus === 'idle'}
          <button class="btn-primary btn-sm" on:click={checkUpdate}>檢查更新</button>
        {:else if updateStatus === 'checking'}
          <span class="update-msg">檢查中...</span>
        {:else if updateStatus === 'upToDate'}
          <span class="update-msg update-ok">已是最新版本</span>
          <button class="btn-outline btn-sm" on:click={checkUpdate}>重新檢查</button>
        {:else if updateStatus === 'available'}
          <span class="update-msg update-new">新版本 {latestVersion} 可用</span>
          <button class="btn-primary btn-sm" on:click={doUpdate}>立即更新</button>
        {:else if updateStatus === 'updating'}
          <span class="update-msg">更新中，請稍候...</span>
        {:else if updateStatus === 'error'}
          <span class="update-msg update-err">{updateError}</span>
          <button class="btn-outline btn-sm" on:click={checkUpdate}>重試</button>
        {/if}
      </div>
    </div>
  </div>
</div>

<style>
  .settings-section {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }
  .section-title {
    font-size: 16px;
    font-weight: 600;
    margin: 0;
  }
  .section-desc {
    font-size: 12px;
    color: var(--text-secondary);
    margin-bottom: 16px;
  }
  .form-actions {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-top: 8px;
  }
  .save-btn {
    padding: 10px 28px;
    font-size: 14px;
    font-weight: 600;
  }
  .save-ok {
    color: var(--success);
    font-size: 13px;
    font-weight: 500;
  }
  .divider {
    height: 1px;
    background: var(--border);
    margin: 20px 0;
  }
  .update-section {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }
  .subsection-title {
    font-size: 14px;
    font-weight: 600;
    margin: 0;
  }
  .update-card {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: var(--bg-primary);
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 12px 16px;
    gap: 16px;
    flex-wrap: wrap;
  }
  .update-info {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .version-label {
    font-size: 12px;
    color: var(--text-secondary);
  }
  .version-value {
    font-size: 13px;
    font-weight: 600;
    font-family: monospace;
  }
  .update-actions {
    display: flex;
    align-items: center;
    gap: 10px;
  }
  .update-msg {
    font-size: 13px;
  }
  .update-ok {
    color: var(--success);
  }
  .update-new {
    color: var(--accent);
    font-weight: 500;
  }
  .update-err {
    color: var(--danger);
  }
</style>
