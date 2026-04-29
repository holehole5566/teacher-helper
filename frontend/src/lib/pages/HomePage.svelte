<script lang="ts">
  import { onMount } from 'svelte';
  import StudentPicker from '../components/StudentPicker.svelte';
  import { GetTodayDuty, GetStudents, GetTimetable, GetSettings } from '../../../wailsjs/go/main/App';

  const periodLabels = ['1', '2', '3', '4', '午休', '5', '6', '7'];

  let displayDate = '';
  let isWorkday = false;
  let dutyStudents: any[] = [];
  let lunchAssignments: any[] = [];
  let loading = true;

  let origDuty: any[] = [];
  let origLunch: any[] = [];
  let dutyModified = false;
  let lunchModified = false;

  let allStudents: any[] = [];
  let showPicker = false;
  let pickerCallback: ((s: any) => void) | null = null;

  let todayClasses: { period: string; subject: string; idx: number }[] = [];
  let currentPeriod = -1;
  let periodTimes: string[] = [];

  export async function refresh() {
    loading = true;
    try {
      const result = await GetTodayDuty();
      displayDate = result.displayDate || '';
      isWorkday = result.isWorkday || false;
      origDuty = result.dutyStudents || [];
      origLunch = result.lunchAssignments || [];
      dutyStudents = [...origDuty];
      lunchAssignments = origLunch.map((a: any) => ({ ...a }));
      dutyModified = false;
      lunchModified = false;
      allStudents = await GetStudents();

      const settings = await GetSettings();
      periodTimes = settings.period_times || [];

      const tt = await GetTimetable();
      const dow = new Date().getDay();
      todayClasses = [];
      if (dow >= 1 && dow <= 5 && tt) {
        const dayClasses = tt[dow - 1] || [];
        for (let i = 0; i < dayClasses.length; i++) {
          if (i === 4) {
            todayClasses.push({ period: '午休', subject: dayClasses[i] || '午休', idx: i });
          } else if (dayClasses[i]) {
            todayClasses.push({ period: periodLabels[i], subject: dayClasses[i], idx: i });
          }
        }
      }
      updateCurrentPeriod();
    } catch (e) {
      console.error('Failed to load:', e);
    }
    loading = false;
  }

  function updateCurrentPeriod() {
    const now = new Date();
    const nowMin = now.getHours() * 60 + now.getMinutes();
    currentPeriod = -1;
    for (let i = periodTimes.length - 1; i >= 0; i--) {
      if (!periodTimes[i]) continue;
      const [h, m] = periodTimes[i].split(':').map(Number);
      if (nowMin >= h * 60 + m) {
        currentPeriod = i;
        break;
      }
    }
  }

  function openPicker(cb: (s: any) => void) {
    pickerCallback = cb;
    showPicker = true;
  }
  function closePicker() {
    showPicker = false;
    pickerCallback = null;
  }
  function replaceDuty(index: number) {
    openPicker((s) => {
      dutyStudents[index] = s;
      dutyStudents = dutyStudents;
      dutyModified = true;
      closePicker();
    });
  }
  function replaceLunch(index: number) {
    openPicker((s) => {
      lunchAssignments[index] = { ...lunchAssignments[index], student: s };
      lunchAssignments = lunchAssignments;
      lunchModified = true;
      closePicker();
    });
  }
  function resetAll() {
    dutyStudents = [...origDuty];
    lunchAssignments = origLunch.map((a: any) => ({ ...a }));
    dutyModified = false;
    lunchModified = false;
  }

  onMount(refresh);
</script>

{#if showPicker}
  <StudentPicker students={allStudents} onPick={(s) => pickerCallback && pickerCallback(s)} onClose={closePicker} />
{/if}

<div class="page">
  <div class="date-header">
    <h2 class="page-title">今日總覽</h2>
    <span class="date-display">{displayDate}</span>
    {#if dutyModified || lunchModified}
      <button class="btn-outline btn-sm" on:click={resetAll}>↺ 重設排程</button>
    {/if}
  </div>

  {#if loading}
    <div class="empty-state">載入中...</div>
  {:else if !isWorkday}
    <div class="card">
      <div class="rest-message">
        <span class="rest-icon">🎉</span>
        <h3>今日休息</h3>
        <p>今天是週末或假期，無需值日</p>
      </div>
    </div>
  {:else}
    <div class="home-grid">
      <!-- Left: timetable -->
      <div class="card timetable-section">
        <div class="section-header">
          <span class="section-icon">📚</span>
          <h3 class="section-title">今日課表</h3>
        </div>
        {#if todayClasses.length > 0}
          <div class="class-list">
            {#each todayClasses as cls}
              <div class="class-row" class:active={cls.idx === currentPeriod}>
                <span class="period-num">{cls.period}</span>
                <span class="subject">{cls.subject}</span>
                {#if cls.idx === currentPeriod}
                  <span class="now-tag">NOW</span>
                {/if}
              </div>
            {/each}
          </div>
        {:else}
          <p class="no-data">今日無課程</p>
        {/if}
      </div>

      <!-- Right: duty + lunch compact -->
      <div class="right-panel">
        <div class="card compact-card duty-accent">
          <div class="section-header">
            <span class="section-icon">🧹</span>
            <h3 class="section-title">值日生</h3>
            {#if dutyModified}<span class="badge badge-muted">已調整</span>{/if}
          </div>
          {#if dutyStudents && dutyStudents.length > 0}
            <div class="chip-list">
              {#each dutyStudents as student, i}
                <button class="student-chip" title="點擊替換" on:click={() => replaceDuty(i)}>
                  <span class="seat">{student.seat_number}號</span>
                  <span class="name">{student.name}</span>
                </button>
              {/each}
            </div>
          {:else}
            <p class="no-data">今日無值日</p>
          {/if}
        </div>

        <div class="card compact-card lunch-accent">
          <div class="section-header">
            <span class="section-icon">🍱</span>
            <h3 class="section-title">抬餐</h3>
            {#if lunchModified}<span class="badge badge-muted">已調整</span>{/if}
          </div>
          {#if lunchAssignments && lunchAssignments.length > 0}
            <div class="lunch-list">
              {#each lunchAssignments as a, i}
                <button class="lunch-row" title="點擊替換" on:click={() => replaceLunch(i)}>
                  <span class="student-name">{a.student.seat_number}號 {a.student.name}</span>
                  <span class="bucket-tag">{a.bucket}</span>
                </button>
              {/each}
            </div>
          {:else}
            <p class="no-data">今日無抬餐</p>
          {/if}
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .date-header {
    display: flex;
    align-items: baseline;
    gap: 16px;
    margin-bottom: 20px;
  }
  .date-display {
    color: var(--text-secondary);
    font-size: 14px;
  }
  .rest-message {
    text-align: center;
    padding: 40px 20px;
  }
  .rest-icon {
    font-size: 48px;
    display: block;
    margin-bottom: 12px;
  }
  .rest-message h3 {
    font-size: 18px;
    margin-bottom: 8px;
  }
  .rest-message p {
    color: var(--text-secondary);
  }

  /* Grid layout: timetable left, duty/lunch right */
  .home-grid {
    display: grid;
    grid-template-columns: 1fr 260px;
    gap: 16px;
    flex: 1;
    min-height: 0;
  }

  /* Timetable section */
  .timetable-section {
    border-left: 4px solid #22c55e;
  }
  .section-header {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 12px;
  }
  .section-icon {
    font-size: 18px;
  }
  .section-title {
    font-size: 15px;
    font-weight: 600;
  }
  .class-list {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }
  .class-row {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 8px 12px;
    border-radius: 8px;
    transition: background 0.2s;
  }
  .class-row.active {
    background: rgba(34, 197, 94, 0.12);
  }
  .period-num {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    background: #f1f5f9;
    font-size: 13px;
    font-weight: 700;
    flex-shrink: 0;
  }
  .class-row.active .period-num {
    background: #22c55e;
    color: white;
  }
  .subject {
    font-size: 15px;
    font-weight: 600;
    flex: 1;
  }
  .now-tag {
    font-size: 10px;
    font-weight: 700;
    background: #22c55e;
    color: white;
    padding: 2px 8px;
    border-radius: 20px;
    letter-spacing: 1px;
  }

  /* Right panel */
  .right-panel {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
  .compact-card {
    padding: 14px !important;
  }
  .duty-accent {
    border-left: 4px solid var(--accent);
  }
  .lunch-accent {
    border-left: 4px solid var(--warning);
  }

  /* Duty chips */
  .chip-list {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
  }
  .student-chip {
    display: flex;
    align-items: center;
    gap: 6px;
    background: #eff6ff;
    border-radius: 8px;
    padding: 8px 12px;
    cursor: pointer;
    transition: background 0.15s;
    border: none;
    font-family: inherit;
  }
  .student-chip:hover {
    background: #dbeafe;
  }
  .seat {
    font-size: 11px;
    color: var(--accent);
    font-weight: 600;
  }
  .name {
    font-size: 14px;
    font-weight: 700;
    color: var(--text-primary);
  }

  /* Lunch list */
  .lunch-list {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }
  .lunch-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 6px 10px;
    background: #fffbeb;
    border-radius: 6px;
    border: none;
    width: 100%;
    cursor: pointer;
    font-family: inherit;
    transition: background 0.15s;
  }
  .lunch-row:hover {
    background: #fef3c7;
  }
  .student-name {
    font-size: 13px;
    font-weight: 600;
    color: var(--text-primary);
  }
  .bucket-tag {
    background: #fef3c7;
    color: #92400e;
    padding: 2px 8px;
    border-radius: 20px;
    font-size: 11px;
    font-weight: 600;
  }

  .no-data {
    color: var(--text-secondary);
    font-size: 13px;
  }
</style>
