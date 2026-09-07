<script lang="ts">
  import { draftGuard, markDraftSaved } from '../../stores/draft';
  let ready = false;
  import { onMount } from 'svelte';
  import { GetSettings, SaveSettings } from '../../../../wailsjs/go/main/App';
  import { verifyAndRun } from '../../stores/auth';

  let semesterStart = '';
  let dutyGroupSize = 2;
  let dutyStartNumber = 1;
  let lunchGroupSize = 5;
  let lunchStartNumber = 1;
  let mealBucketsStr = '';
  let saved = false;

  async function loadSettings() {
    const s = await GetSettings();
    semesterStart = s.semester_start_date || '';
    dutyGroupSize = s.duty_group_size || 2;
    dutyStartNumber = s.duty_start_number || 1;
    lunchGroupSize = s.lunch_group_size || 5;
    lunchStartNumber = s.lunch_start_number || 1;
    mealBucketsStr = (s.meal_buckets || []).join('，');
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

    try {
      await verifyAndRun(async () => {
        const current = await GetSettings();
        await SaveSettings(Object.assign(current, {
          semester_start_date: semesterStart,
          duty_group_size: dutyGroupSize,
          duty_start_number: dutyStartNumber,
          lunch_group_size: lunchGroupSize,
          lunch_start_number: lunchStartNumber,
          meal_buckets: buckets,
        }));
      }, '儲存值日設定');
      markDraftSaved();
      saved = true;
      setTimeout(() => { saved = false; }, 2000);
    } catch (e: any) {
      if (e.message !== '驗證失敗或已取消') {
        alert('儲存失敗: ' + e.message);
      }
    }
  }

  onMount(() => { loadSettings().then(() => { ready = true; }).catch((e) => alert("載入失敗，請重新開啟此頁：" + e)); });
</script>

<div class="settings-section" use:draftGuard={{ value: JSON.stringify([semesterStart, dutyGroupSize, dutyStartNumber, lunchGroupSize, lunchStartNumber, mealBucketsStr]), ready, save: handleSave }}>
  <h3 class="section-title">值日與抬餐</h3>
  <p class="section-desc">設定學期日期與輪值規則</p>

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
