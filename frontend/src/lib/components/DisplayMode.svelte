<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { GetTodayDuty, GetTimetable, GetSettings } from '../../../wailsjs/go/main/App';

  const dayLabels = ['日', '一', '二', '三', '四', '五', '六'];
  const periodLabels = ['1', '2', '3', '4', '午休', '5', '6', '7'];

  let now = new Date();
  let clockStr = '';
  let dateStr = '';
  let dutyStudents: any[] = [];
  let lunchAssignments: any[] = [];
  let todayClasses: { period: string; subject: string; idx: number }[] = [];
  let currentPeriod = -1;
  let periodTimes: string[] = [];
  let clockInterval: ReturnType<typeof setInterval>;

  function updateClock() {
    now = new Date();
    const hh = String(now.getHours()).padStart(2, '0');
    const mm = String(now.getMinutes()).padStart(2, '0');
    const ss = String(now.getSeconds()).padStart(2, '0');
    clockStr = `${hh}:${mm}:${ss}`;

    // Determine current period from period_times
    // periodTimes[i] is the start time of period i+1
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

  async function load() {
    const result = await GetTodayDuty();
    const d = new Date();
    const mm = String(d.getMonth() + 1).padStart(2, '0');
    const dd = String(d.getDate()).padStart(2, '0');
    dateStr = `${mm}/${dd}（${dayLabels[d.getDay()]}）`;
    dutyStudents = result.dutyStudents || [];
    lunchAssignments = result.lunchAssignments || [];

    const settings = await GetSettings();
    periodTimes = settings.period_times || [];

    const tt = await GetTimetable();
    const dow = d.getDay();
    if (dow >= 1 && dow <= 5 && tt) {
      const dayClasses = tt[dow - 1] || [];
      const classes: typeof todayClasses = [];
      for (let i = 0; i < dayClasses.length; i++) {
        if (i === 4) {
          classes.push({ period: periodLabels[i], subject: dayClasses[i] || '午休', idx: i });
        } else if (dayClasses[i]) {
          classes.push({ period: periodLabels[i], subject: dayClasses[i], idx: i });
        }
      }
      todayClasses = classes;
    }

    updateClock();
  }

  onMount(() => {
    load();
    clockInterval = setInterval(updateClock, 1000);
  });

  onDestroy(() => clearInterval(clockInterval));
</script>

<div class="display">
  <!-- Header: date + clock -->
  <div class="display-header">
    <span class="date">{dateStr}</span>
    <span class="clock">{clockStr}</span>
  </div>

  <div class="display-body">
    <!-- Left: timetable (main focus) -->
    <div class="panel-left">
      <div class="section-card timetable-card">
        <h2>📚 今日課表</h2>
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
    </div>

    <!-- Right: duty + lunch -->
    <div class="panel-right">
      <div class="section-card duty-card">
        <h2>🧹 今日值日生</h2>
        <div class="duty-list">
          {#each dutyStudents as s}
            <div class="duty-chip">
              <span class="seat">{s.seat_number}號</span>
              <span class="name">{s.name}</span>
            </div>
          {/each}
        </div>
      </div>

      <div class="section-card lunch-card">
        <h2>🍱 抬餐同學</h2>
        <div class="lunch-list">
          {#each lunchAssignments as a}
            <div class="lunch-row">
              <span class="student">{a.student.seat_number}號 {a.student.name}</span>
              <span class="bucket">{a.bucket}</span>
            </div>
          {/each}
        </div>
      </div>
    </div>
  </div>

  <div class="display-footer">按 ESC 退出展示模式</div>
</div>

<style>
  .display {
    position: fixed;
    inset: 0;
    z-index: 9000;
    background: #eef0f2;
    color: #1e293b;
    display: flex;
    flex-direction: column;
    padding: 32px 40px;
    font-family: var(--font-family);
    user-select: none;
  }

  /* Header */
  .display-header {
    display: flex;
    justify-content: space-between;
    align-items: baseline;
    margin-bottom: 28px;
    padding-bottom: 20px;
    border-bottom: 1px solid #e2e8f0;
  }
  .date {
    font-size: 36px;
    font-weight: 700;
  }
  .clock {
    font-size: 48px;
    font-weight: 700;
    font-variant-numeric: tabular-nums;
    color: #3b82f6;
  }

  /* Body */
  .display-body {
    flex: 1;
    display: grid;
    grid-template-columns: 1fr 340px;
    gap: 24px;
    min-height: 0;
  }
  .panel-left {
    display: flex;
    flex-direction: column;
  }
  .panel-right {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  /* Cards */
  .section-card {
    background: #ffffff;
    border-radius: 16px;
    padding: 24px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.08);
  }
  .section-card h2 {
    font-size: 18px;
    font-weight: 600;
    margin-bottom: 16px;
    opacity: 0.8;
  }
  .timetable-card {
    flex: 1;
    border-left: 4px solid #22c55e;
  }
  .duty-card {
    border-left: 4px solid #3b82f6;
  }
  .lunch-card {
    flex: 1;
    border-left: 4px solid #f59e0b;
  }

  /* Timetable */
  .class-list {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }
  .class-row {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 14px 20px;
    border-radius: 10px;
    transition: background 0.2s;
  }
  .class-row.active {
    background: rgba(34, 197, 94, 0.15);
  }
  .period-num {
    width: 36px;
    height: 36px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    background: #f1f5f9;
    font-size: 16px;
    font-weight: 700;
    flex-shrink: 0;
  }
  .class-row.active .period-num {
    background: #22c55e;
    color: #0f172a;
  }
  .subject {
    font-size: 22px;
    font-weight: 600;
    flex: 1;
  }
  .now-tag {
    font-size: 12px;
    font-weight: 700;
    background: #22c55e;
    color: #0f172a;
    padding: 3px 10px;
    border-radius: 20px;
    letter-spacing: 1px;
  }

  /* Duty */
  .duty-list {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
  }
  .duty-chip {
    display: flex;
    align-items: center;
    gap: 8px;
    background: #eff6ff;
    padding: 12px 18px;
    border-radius: 10px;
  }
  .duty-chip .seat {
    color: #3b82f6;
    font-size: 14px;
    font-weight: 600;
  }
  .duty-chip .name {
    font-size: 20px;
    font-weight: 700;
  }

  /* Lunch */
  .lunch-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  .lunch-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 14px;
    background: #fffbeb;
    border-radius: 8px;
  }
  .student {
    font-size: 15px;
    font-weight: 600;
  }
  .bucket {
    background: #fef3c7;
    color: #92400e;
    padding: 3px 10px;
    border-radius: 20px;
    font-size: 12px;
    font-weight: 600;
  }

  .no-data {
    opacity: 0.4;
    font-size: 16px;
  }

  /* Footer */
  .display-footer {
    text-align: center;
    padding-top: 16px;
    font-size: 12px;
    opacity: 0.4;
  }
</style>
