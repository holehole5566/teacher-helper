<script lang="ts">
  import { onMount } from 'svelte';
  import { GetSettings, SaveSettings } from '../../../../wailsjs/go/main/App';
  import { verifyAndRun } from '../../stores/auth';

  let countdownTimes: string[] = [];
  let periodTimes: string[] = ['','','','','','','',''];
  const periodTimeLabels = ['第1節', '第2節', '第3節', '第4節', '午休', '第5節', '第6節', '第7節'];
  let newTime = '';
  let saved = false;

  async function loadSettings() {
    const s = await GetSettings();
    countdownTimes = (s.countdown_times || []).slice().sort();
    const pt = s.period_times || [];
    for (let i = 0; i < 8; i++) periodTimes[i] = pt[i] || '';
    periodTimes = periodTimes;
  }

  function addTime() {
    const t = newTime.trim();
    if (t && /^\d{2}:\d{2}$/.test(t) && !countdownTimes.includes(t)) {
      countdownTimes = [...countdownTimes, t].sort();
      newTime = '';
    }
  }

  function removeTime(t: string) {
    countdownTimes = countdownTimes.filter(x => x !== t);
  }

  async function handleSave() {
    try {
      const current = await GetSettings();
      await verifyAndRun(async () => {
        await SaveSettings({
          ...current,
          countdown_times: countdownTimes,
          period_times: periodTimes,
        });
      }, '儲存課程時間設定');
      saved = true;
      setTimeout(() => { saved = false; }, 2000);
    } catch (e: any) {
      if (e.message !== '驗證失敗或已取消') {
        alert('儲存失敗: ' + e.message);
      }
    }
  }

  onMount(loadSettings);
</script>

<div class="settings-section">
  <h3 class="section-title">課程時間</h3>
  <p class="section-desc">設定每節課開始時間與倒數觸發點</p>

  <div class="form-group">
    <label>每節上課時間</label>
    <div class="period-grid">
      {#each periodTimes as _, i}
        <div class="period-input">
          <span class="period-label">{periodTimeLabels[i]}</span>
          <input type="time" bind:value={periodTimes[i]} />
        </div>
      {/each}
    </div>
  </div>

  <div class="form-group">
    <label>上課倒數時間點（到達該時間前 1 分鐘開始倒數）</label>
    <div class="inline-form">
      <input type="time" bind:value={newTime} />
      <button class="btn-primary btn-sm" on:click={addTime}>新增</button>
    </div>
    {#if countdownTimes.length > 0}
      <div class="time-tags">
        {#each countdownTimes as t}
          <span class="time-tag">
            {t}
            <button class="tag-remove" on:click={() => removeTime(t)}>×</button>
          </span>
        {/each}
      </div>
    {/if}
  </div>

  <div class="form-actions">
    <button class="btn-primary save-btn" on:click={handleSave}>儲存設定</button>
    {#if saved}
      <span class="save-ok">已儲存</span>
    {/if}
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
  .period-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 8px;
    margin-top: 6px;
  }
  .period-input {
    display: flex;
    align-items: center;
    gap: 8px;
    background: var(--bg-primary);
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 8px 12px;
  }
  .period-label {
    font-size: 12px;
    font-weight: 600;
    color: var(--text-secondary);
    white-space: nowrap;
    min-width: 36px;
  }
  .period-input input {
    flex: 1;
    padding: 4px 6px;
    font-size: 13px;
    border: none;
    background: transparent;
  }
  .time-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-top: 10px;
  }
  .time-tag {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    background: var(--bg-primary);
    border: 1px solid var(--border);
    border-radius: 20px;
    padding: 6px 12px;
    font-size: 13px;
    font-weight: 500;
  }
  .tag-remove {
    background: none;
    border: none;
    padding: 0 2px;
    font-size: 16px;
    color: var(--danger);
    cursor: pointer;
    line-height: 1;
  }
  .tag-remove:hover {
    color: var(--danger-hover);
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
</style>
