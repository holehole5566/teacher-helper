<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import StudentPicker from '../components/StudentPicker.svelte';
  import HistoryNav from '../components/HistoryNav.svelte';
  import { GetDailyOverview, GetDailyDates, SaveDailyAssignments } from '../../../wailsjs/go/main/App';
  import { RotateCcw, PartyPopper, BookOpen, Brush, UtensilsCrossed, FileText } from 'lucide-svelte';

  export let presentation = false;
  export let selectedDate = '';
  export let onDateChange: (value: string) => void = () => {};
  let dates: string[] = [];
  let date = '';
  let today = '';
  let error = '';
  let reviewFromYesterday = '';
  let busy = false;
  let touchX = 0;
  let touchY = 0;
  let clock = '';
  $: historical = date !== today;
  const periodLabels = ['1', '2', '3', '4', '午休', '5', '6', '7'];

  let displayDate = '';
  let isWorkday = false;
  let dutyStudents: any[] = [];
  let lunchAssignments: any[] = [];
  let loading = true;

  let dutyModified = false;
  let lunchModified = false;

  let allStudents: any[] = [];
  let showPicker = false;
  let pickerCallback: ((s: any) => void) | null = null;

  let todayClasses: { period: string; subject: string; idx: number }[] = [];
  let currentPeriod = -1;
  let periodTimes: string[] = [];
  let clockInterval: ReturnType<typeof setInterval>;
  let missingHomework: { subject: string; students: number[]; note: string }[] = [];

  export async function refresh() {
    loading = true;
    try {
      error = '';
      const record = await GetDailyOverview(date);
      const result = record.duty;
      displayDate = result.displayDate || date;
      isWorkday = result.isWorkday;
      dutyStudents = result.dutyStudents || [];
      lunchAssignments = result.lunchAssignments || [];
      dutyModified = record.manual_duty;
      lunchModified = record.manual_lunch;
      allStudents = record.students || [];
      periodTimes = record.period_times || [];
      missingHomework = record.missing_homework || [];
      reviewFromYesterday = record.review_from_yesterday || '';
      const tt = record.timetable;
      const dow = new Date(date + "T12:00:00").getDay();
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

    } catch (e: any) {
      error = String(e?.message || e);
      displayDate = date;
      isWorkday = false;
      reviewFromYesterday = ''; dutyModified = false; lunchModified = false;
    }
    loading = false;
  }

  function updateCurrentPeriod() {
    const now = new Date();
    const nowMin = now.getHours() * 60 + now.getMinutes();
    clock = now.toLocaleTimeString('zh-TW', { hour12: false });
    currentPeriod = -1;
    if (date !== today) return;
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
  async function persist(duty: any[], lunch: any[], reset = false) {
    busy = true;
    try {
      await SaveDailyAssignments(date, duty, lunch, reset);
      await refresh();
      closePicker();
    } catch (e) { error = `儲存失敗：${e}`; }
    finally { busy = false; }
  }
  function replaceDuty(index: number) {
    if (historical || presentation || busy) return;
    openPicker((student) => { const next = [...dutyStudents]; next[index] = student; persist(next, lunchAssignments); });
  }
  function replaceLunch(index: number) {
    if (historical || presentation || busy) return;
    openPicker((student) => { const next = [...lunchAssignments]; next[index] = {...next[index], student}; persist(dutyStudents, next); });
  }
  async function changeDate(next: string) { date = next; selectedDate = next; onDateChange(next); closePicker(); await refresh(); }
  function swipe(e: TouchEvent) {
    const dx = e.changedTouches[0].clientX - touchX;
    const dy = e.changedTouches[0].clientY - touchY;
    if (busy || loading || Math.abs(dx) < 80 || Math.abs(dx) < Math.abs(dy)*1.5) return;
    const next = dates.indexOf(date) + (dx > 0 ? 1 : -1);
    if (dates[next]) changeDate(dates[next]);
  }
  let checking = false;
  async function tick() {
    updateCurrentPeriod();
    if (checking) return;
    checking = true;
    try {
      const next = await GetDailyDates();
      if (next[0] !== today) {
        const wasToday = date === today;
        dates = next; today = next[0];
        if (wasToday || !dates.includes(date)) await changeDate(today);
      }
    } finally { checking = false; }
  }
  onMount(() => {
    (async () => { try { dates = await GetDailyDates(); today = dates[0]; await changeDate(dates.includes(selectedDate) ? selectedDate : today); } catch(e) { error = String(e); loading = false; } })();
    clockInterval = setInterval(() => { tick().catch(() => {}); }, 1000);
  });

  onDestroy(() => clearInterval(clockInterval));
</script>

{#if showPicker}
  <StudentPicker students={allStudents} onPick={(s) => pickerCallback && pickerCallback(s)} onClose={closePicker} />
{/if}

<div class="page" class:presentation on:touchstart={(e) => { touchX = e.touches[0].clientX; touchY = e.touches[0].clientY; }} on:touchend={swipe}>
  <HistoryNav {dates} {date} busy={busy || loading} onChange={changeDate} />
  <div class="date-header">
    <h2 class="page-title">{historical ? '歷史總覽' : '今日總覽'}</h2>
    <span class="date-display">{displayDate}</span>
    {#if presentation}<span class="clock">{clock}</span>{/if}
    {#if historical}<span class="badge badge-muted">歷史快照・唯讀</span>{/if}
    {#if !historical && !presentation && (dutyModified || lunchModified)}
      <button class="btn-outline btn-sm" disabled={busy} on:click={() => persist([], [], true)}><RotateCcw size={14} /> 重設排程</button>
    {/if}
  </div>

  {#if error}<p class="error-message" role="alert">{error}</p>{/if}
  {#if !loading && !error && reviewFromYesterday}
    <section class="card scheduled-review" aria-labelledby="scheduled-review-title">
      <h3 id="scheduled-review-title">今日複習考</h3>
      <p class="review-source">前一天聯絡簿的複習考提醒</p>
      <p class="review-content">{reviewFromYesterday}</p>
    </section>
  {/if}
  {#if loading}
    <div class="empty-state">載入中...</div>
  {:else if error && !isWorkday}
    <div class="empty-state">此日期沒有可顯示的資料</div>
  {:else if !isWorkday}
    <div class="card">
      <div class="rest-message">
        <span class="rest-icon"><PartyPopper size={48} /></span>
        <h3>今日休息</h3>
        <p>今天是週末或假期，無需值日</p>
      </div>
    </div>
  {:else}
    <div class="home-grid">
      <!-- Left: timetable -->
      <div class="card timetable-section">
        <div class="section-header">
          <span class="section-icon"><BookOpen size={18} /></span>
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
            <span class="section-icon"><Brush size={18} /></span>
            <h3 class="section-title">值日生</h3>
            {#if dutyModified}<span class="badge badge-muted">已調整</span>{/if}
          </div>
          {#if dutyStudents && dutyStudents.length > 0}
            <div class="chip-list">
              {#each dutyStudents as student, i}
                <button class="student-chip" disabled={historical || presentation || busy} title="點擊替換" on:click={() => replaceDuty(i)}>
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
            <span class="section-icon"><UtensilsCrossed size={18} /></span>
            <h3 class="section-title">抬餐</h3>
            {#if lunchModified}<span class="badge badge-muted">已調整</span>{/if}
          </div>
          {#if lunchAssignments && lunchAssignments.length > 0}
            <div class="lunch-list">
              {#each lunchAssignments as a, i}
                <button class="lunch-row" disabled={historical || presentation || busy} title="點擊替換" on:click={() => replaceLunch(i)}>
                  <span class="student-name">{a.student.seat_number}號 {a.student.name}</span>
                  <span class="bucket-tag">{a.bucket}</span>
                </button>
              {/each}
            </div>
          {:else}
            <p class="no-data">今日無抬餐</p>
          {/if}
        </div>

        {#if missingHomework.length > 0}
          <div class="card compact-card homework-accent">
            <div class="section-header">
              <span class="section-icon"><FileText size={18} /></span>
              <h3 class="section-title">作業未交</h3>
            </div>
            <div class="homework-list">
              {#each missingHomework as hw}
                <div class="homework-item">
                  <span class="hw-subject">{hw.subject}</span>
                  <div class="hw-chips">
                    {#each hw.students as seat}
                      <span class="hw-chip">{seat}號 {allStudents.find(s => s.seat_number === seat)?.name || ''}</span>
                    {/each}
                  </div>
                  {#if hw.note}
                    <span class="hw-note">{hw.note}</span>
                  {/if}
                </div>
              {/each}
            </div>
          </div>
        {/if}
      </div>
    </div>
  {/if}
</div>

<style>
  .scheduled-review { margin-bottom:16px; }
  .scheduled-review h3 { font-size: calc(var(--ui-font, 14px) * 18 / 14); margin-bottom:4px; }
  .review-source { font-size: calc(var(--ui-font, 14px) * 13 / 14); color:var(--text-secondary); margin-bottom:12px; }
  .review-content { font-size: calc(var(--ui-font, 14px) * 20 / 14); font-weight:600; white-space:pre-wrap; overflow-wrap:anywhere; }
  .presentation .review-content { font-size: calc(var(--ui-font, 14px) * 28 / 14); }
  .error-message { color:#b91c1c; padding:12px 0; }
  .clock { margin-left:auto; font-variant-numeric:tabular-nums; }
  .presentation { font-size: calc(var(--ui-font, 14px) * 20 / 14); }
  .presentation .home-grid { grid-template-columns:minmax(0, 1fr) minmax(280px, 1fr); }
  .presentation .subject, .presentation .name, .presentation .student-name { font-size: calc(var(--ui-font, 14px) * 22 / 14); }
  .presentation .hw-chip, .presentation .seat { font-size: calc(var(--ui-font, 14px) * 16 / 14); }
  button:disabled { cursor:default; }
  @media(max-width:700px) { .presentation .home-grid { grid-template-columns:1fr; } }
  .date-header {
    display: flex;
    align-items: baseline;
    flex-wrap:wrap;
    gap: 16px;
    margin-bottom: 20px;
  }
  .date-display {
    color: var(--text-secondary);
    font-size: calc(var(--ui-font, 14px) * 14 / 14);
  }
  .rest-message {
    text-align: center;
    padding: 40px 20px;
  }
  .rest-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 12px;
    color: #f59e0b;
  }
  .rest-message h3 {
    font-size: calc(var(--ui-font, 14px) * 18 / 14);
    margin-bottom: 8px;
  }
  .rest-message p {
    color: var(--text-secondary);
  }

  /* Grid layout: timetable left, duty/lunch right */
  .home-grid {
    display: grid;
    grid-template-columns: minmax(0, 1.7fr) minmax(240px, 1fr);
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
    display: inline-flex;
    align-items: center;
    color: var(--text-secondary);
  }
  .section-title {
    font-size: calc(var(--ui-font, 14px) * 15 / 14);
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
    font-size: calc(var(--ui-font, 14px) * 13 / 14);
    font-weight: 700;
    flex-shrink: 0;
  }
  .class-row.active .period-num {
    background: #22c55e;
    color: white;
  }
  .subject {
    font-size: calc(var(--ui-font, 14px) * 15 / 14);
    font-weight: 600;
    flex: 1;
  }
  .now-tag {
    font-size: calc(var(--ui-font, 14px) * 10 / 14);
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
    font-size: calc(var(--ui-font, 14px) * 11 / 14);
    color: var(--accent);
    font-weight: 600;
  }
  .name {
    font-size: calc(var(--ui-font, 14px) * 14 / 14);
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
    font-size: calc(var(--ui-font, 14px) * 13 / 14);
    font-weight: 600;
    color: var(--text-primary);
  }
  .bucket-tag {
    background: #fef3c7;
    color: #92400e;
    padding: 2px 8px;
    border-radius: 20px;
    font-size: calc(var(--ui-font, 14px) * 11 / 14);
    font-weight: 600;
  }

  .no-data {
    color: var(--text-secondary);
    font-size: calc(var(--ui-font, 14px) * 13 / 14);
  }

  .homework-accent {
    border-left: 4px solid #ef4444;
  }
  .homework-list {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }
  .homework-item {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }
  .hw-subject {
    font-size: calc(var(--ui-font, 14px) * 13 / 14);
    font-weight: 700;
    color: #991b1b;
  }
  .hw-chips {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
  }
  .hw-chip {
    background: #fee2e2;
    color: #991b1b;
    padding: 2px 8px;
    border-radius: 4px;
    font-size: calc(var(--ui-font, 14px) * 12 / 14);
    font-weight: 600;
  }
  .hw-note {
    font-size: calc(var(--ui-font, 14px) * 11 / 14);
    color: var(--text-secondary);
  }
  .page:not(.presentation) { --ui-font: clamp(14px, calc(2px + 1vw), 26px); font-size: var(--ui-font); }
  .page:not(.presentation) .right-panel { min-width: 0; }
  .subject, .student-name, .name { overflow-wrap: anywhere; min-width: 0; }
  .lunch-row { gap: 8px; text-align: left; }
  .bucket-tag { flex-shrink: 0; }
  .page:not(.presentation) .class-row { padding: 0.65em 0.8em; }
  .page:not(.presentation) .compact-card { padding: 1em !important; }
  .page:not(.presentation) .period-num { width: 2.2em; height: 2.2em; }
  @container workspace (max-width: 620px) {
    .home-grid { grid-template-columns: minmax(0, 1fr); }
    .date-header { gap: 8px 12px; }
  }
</style>
