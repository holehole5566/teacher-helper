<script lang="ts">
  import { onMount } from 'svelte';
  import { GetTimetable, SaveTimetable } from '../../../wailsjs/go/main/App';

  const days = ['一', '二', '三', '四', '五'];
  const periods = 8;
  const periodLabels = ['第1節', '第2節', '第3節', '第4節', '午休', '第5節', '第6節', '第7節'];

  // timetable[day][period]
  let timetable: string[][] = Array.from({ length: 5 }, () => Array(periods).fill(''));
  let saved = false;

  async function load() {
    const data = await GetTimetable();
    if (data) {
      for (let d = 0; d < 5; d++) {
        for (let p = 0; p < periods; p++) {
          timetable[d][p] = data[d]?.[p] || '';
        }
      }
      timetable = timetable;
    }
  }

  async function handleSave() {
    await SaveTimetable(timetable);
    saved = true;
    setTimeout(() => { saved = false; }, 2000);
  }

  onMount(load);
</script>

<div class="page">
  <h2 class="page-title">課表設定</h2>

  <div class="card">
    <div class="table-wrapper">
      <table class="timetable">
        <thead>
          <tr>
            <th class="period-col">節次</th>
            {#each days as day}
              <th>週{day}</th>
            {/each}
          </tr>
        </thead>
        <tbody>
          {#each Array(periods) as _, p}
            <tr>
              <td class="period-col">{periodLabels[p]}</td>
              {#each days as _, d}
                <td class="cell">
                  <input
                    type="text"
                    bind:value={timetable[d][p]}
                    placeholder="—"
                    class="cell-input"
                  />
                </td>
              {/each}
            </tr>
          {/each}
        </tbody>
      </table>
    </div>

    <div class="form-actions">
      <button class="btn-primary" on:click={handleSave}>儲存課表</button>
      {#if saved}
        <span class="save-ok">已儲存</span>
      {/if}
    </div>
  </div>
</div>

<style>
  .table-wrapper {
    overflow-x: auto;
    margin-bottom: 16px;
  }
  .timetable {
    width: 100%;
    border-collapse: collapse;
  }
  .timetable th {
    text-align: center;
    padding: 8px;
    font-size: 13px;
  }
  .period-col {
    width: 60px;
    text-align: center;
    font-weight: 600;
    color: var(--text-secondary);
    font-size: 13px;
  }
  .cell {
    padding: 3px;
  }
  .cell-input {
    width: 100%;
    text-align: center;
    padding: 8px 4px;
    border: 1px solid var(--border);
    border-radius: 6px;
    font-size: 13px;
    background: var(--bg-primary);
    transition: border-color 0.15s;
  }
  .cell-input:focus {
    border-color: var(--accent);
    background: white;
  }
  .form-actions {
    display: flex;
    align-items: center;
    gap: 12px;
  }
  .save-ok {
    color: var(--success);
    font-size: 13px;
    font-weight: 500;
  }
</style>
