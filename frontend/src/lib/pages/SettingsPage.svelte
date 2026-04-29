<script lang="ts">
  import { onMount } from 'svelte';
  import { GetSettings, SaveSettings } from '../../../wailsjs/go/main/App';

  let semesterStart = '';
  let dutyGroupSize = 2;
  let dutyStartNumber = 1;
  let lunchGroupSize = 5;
  let lunchStartNumber = 1;
  let mealBucketsStr = '';
  let countdownTimes: string[] = [];
  let periodTimes: string[] = ['','','','','','','',''];
  const periodTimeLabels = ['第1節', '第2節', '第3節', '第4節', '午休', '第5節', '第6節', '第7節'];
  let newTime = '';
  let saved = false;

  async function loadSettings() {
    const s = await GetSettings();
    semesterStart = s.semester_start_date || '';
    dutyGroupSize = s.duty_group_size || 2;
    dutyStartNumber = s.duty_start_number || 1;
    lunchGroupSize = s.lunch_group_size || 5;
    lunchStartNumber = s.lunch_start_number || 1;
    mealBucketsStr = (s.meal_buckets || []).join('，');
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
    const buckets = mealBucketsStr
      .split(/[,，]/)
      .map((b: string) => b.trim())
      .filter((b: string) => b);

    if (buckets.length > 0 && buckets.length !== lunchGroupSize) {
      if (!confirm(`餐桶數量（${buckets.length}）與抬餐每組人數（${lunchGroupSize}）不一致，確定要儲存嗎？`)) {
        return;
      }
    }

    await SaveSettings({
      semester_start_date: semesterStart,
      duty_group_size: dutyGroupSize,
      duty_start_number: dutyStartNumber,
      lunch_group_size: lunchGroupSize,
      lunch_start_number: lunchStartNumber,
      meal_buckets: buckets,
      auto_start: false,
      countdown_times: countdownTimes,
      period_times: periodTimes,
    });
    saved = true;
    setTimeout(() => { saved = false; }, 2000);
  }

  onMount(loadSettings);
</script>

<div class="page">
  <h2 class="page-title">設定</h2>

  <div class="card settings-form">
    <div class="form-group">
      <label>學期開始日期</label>
      <input type="date" bind:value={semesterStart} />
    </div>

    <div class="form-row">
      <div class="form-group">
        <label>值日每組人數</label>
        <input type="number" min="1" bind:value={dutyGroupSize} />
      </div>
      <div class="form-group">
        <label>值日起始座號</label>
        <input type="number" min="1" bind:value={dutyStartNumber} />
      </div>
    </div>

    <div class="form-row">
      <div class="form-group">
        <label>抬餐每組人數</label>
        <input type="number" min="1" bind:value={lunchGroupSize} />
      </div>
      <div class="form-group">
        <label>抬餐起始座號</label>
        <input type="number" min="1" bind:value={lunchStartNumber} />
      </div>
    </div>

    <div class="form-group">
      <label>餐桶名稱（以逗號分隔）</label>
      <input type="text" bind:value={mealBucketsStr} placeholder="飯，菜1，菜2，湯，餐具" />
    </div>

    <div class="form-group">
      <label>每節上課時間（用於展示模式判斷當前節次）</label>
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
</div>

<style>
  .settings-form {
    max-width: 560px;
  }
  .form-row {
    display: flex;
    gap: 16px;
  }
  .form-row .form-group {
    flex: 1;
  }
  .form-actions {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-top: 16px;
  }
  .save-btn {
    padding: 12px 32px;
    font-size: 16px;
    font-weight: 600;
  }
  .save-ok {
    color: var(--success);
    font-size: 13px;
    font-weight: 500;
  }
  .period-grid {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    margin-top: 6px;
  }
  .period-input {
    display: flex;
    align-items: center;
    gap: 4px;
    background: var(--bg-primary);
    border: 1px solid var(--border);
    border-radius: 6px;
    padding: 4px 8px;
  }
  .period-label {
    font-size: 12px;
    font-weight: 600;
    color: var(--text-secondary);
    white-space: nowrap;
  }
  .period-input input {
    width: 110px;
    padding: 4px 6px;
    font-size: 12px;
    border: none;
    background: transparent;
  }
  .time-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-top: 8px;
  }
  .time-tag {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    background: var(--bg-primary);
    border: 1px solid var(--border);
    border-radius: 20px;
    padding: 4px 10px;
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
</style>
