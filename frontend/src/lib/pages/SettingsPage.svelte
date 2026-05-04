<script lang="ts">
  import { onMount } from 'svelte';
  import { GetSettings, SaveSettings, SelectCountdownMusics, GetCountdownMusicData, ValidateRandomPool } from '../../../wailsjs/go/main/App';

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

  type MusicTrack = { path: string; in_random: boolean };
  type CountdownTimeMusic = { time: string; mode: string; index: number };

  let countdownMusics: MusicTrack[] = [];
  let countdownTimeMusicMap: CountdownTimeMusic[] = [];
  let countdownVolume = 0.5;
  let testAudios: (HTMLAudioElement | null)[] = [];

  async function loadSettings() {
    const s = await GetSettings();
    semesterStart = s.semester_start_date || '';
    dutyGroupSize = s.duty_group_size || 2;
    dutyStartNumber = s.duty_start_number || 1;
    lunchGroupSize = s.lunch_group_size || 5;
    lunchStartNumber = s.lunch_start_number || 1;
    mealBucketsStr = (s.meal_buckets || []).join('，');
    countdownTimes = (s.countdown_times || []).slice().sort();
    countdownMusics = (s.countdown_musics || []).map((t: any) => ({ path: t.path, in_random: t.in_random }));
    countdownTimeMusicMap = (s.countdown_time_music_map || []).map((m: any) => ({ time: m.time, mode: m.mode || 'random', index: m.index || 0 }));
    countdownVolume = s.countdown_volume > 0 ? s.countdown_volume : 0.5;
    testAudios = countdownMusics.map(() => null);
    syncTimeMusicMap();
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
    syncTimeMusicMap();
  }

  async function addMusics() {
    const paths = await SelectCountdownMusics();
    if (!paths || paths.length === 0) return;
    const existing = new Set(countdownMusics.map(t => t.path));
    const newTracks: MusicTrack[] = paths
      .filter(p => !existing.has(p))
      .map(p => ({ path: p, in_random: true }));
    countdownMusics = [...countdownMusics, ...newTracks];
    testAudios = countdownMusics.map(() => null);
  }

  function removeTrack(i: number) {
    stopTrackAudio(i);
    countdownMusics = countdownMusics.filter((_, idx) => idx !== i);
    testAudios = countdownMusics.map(() => null);
    countdownTimeMusicMap = countdownTimeMusicMap.map(m => ({
      ...m,
      index: m.index > i ? m.index - 1 : m.index === i ? 0 : m.index
    }));
  }

  function stopTrackAudio(i: number) {
    if (testAudios[i]) { testAudios[i]!.pause(); testAudios[i] = null; testAudios = [...testAudios]; }
  }

  async function previewTrack(i: number) {
    if (testAudios[i]) { stopTrackAudio(i); return; }
    try {
      const url = await GetCountdownMusicData(i);
      if (!url) return;
      const a = new Audio(url);
      a.volume = countdownVolume;
      a.onended = () => { testAudios[i] = null; testAudios = [...testAudios]; };
      testAudios[i] = a;
      testAudios = [...testAudios];
      a.play().catch(() => {});
    } catch (e: any) {
      alert('無法播放：' + (e?.message || '檔案可能已被移動或刪除'));
    }
  }

  function getTimeMusicEntry(t: string): CountdownTimeMusic {
    return countdownTimeMusicMap.find(m => m.time === t) || { time: t, mode: 'random', index: 0 };
  }

  function setTimeMusicMode(t: string, mode: string) {
    const existing = countdownTimeMusicMap.find(m => m.time === t);
    if (existing) {
      existing.mode = mode;
      countdownTimeMusicMap = [...countdownTimeMusicMap];
    } else {
      countdownTimeMusicMap = [...countdownTimeMusicMap, { time: t, mode, index: 0 }];
    }
  }

  function setTimeMusicIndex(t: string, index: number) {
    const existing = countdownTimeMusicMap.find(m => m.time === t);
    if (existing) {
      existing.index = index;
      countdownTimeMusicMap = [...countdownTimeMusicMap];
    } else {
      countdownTimeMusicMap = [...countdownTimeMusicMap, { time: t, mode: 'index', index }];
    }
  }

  function syncTimeMusicMap() {
    countdownTimeMusicMap = countdownTimeMusicMap.filter(m => countdownTimes.includes(m.time));
  }

  function onTimeMusicModeChange(t: string, e: Event) {
    setTimeMusicMode(t, (e.target as HTMLSelectElement).value);
  }

  function onTimeMusicIndexChange(t: string, e: Event) {
    setTimeMusicIndex(t, Number((e.target as HTMLSelectElement).value));
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

    const hasRandomMode = countdownMusics.length > 0 && (
      countdownTimes.some(t => {
        const entry = countdownTimeMusicMap.find(m => m.time === t);
        return !entry || entry.mode === 'random';
      })
    );
    if (hasRandomMode) {
      const valid = await ValidateRandomPool();
      if (!valid) {
        alert('請至少將一首音樂加入隨機清單');
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
      countdown_volume: countdownVolume,
      countdown_musics: countdownMusics,
      countdown_time_music_map: countdownTimeMusicMap,
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
            <div class="time-tag-row">
              <span class="time-tag">
                {t}
                <button class="tag-remove" on:click={() => removeTime(t)}>×</button>
              </span>
              {#if countdownMusics.length > 0}
                <select
                  class="mode-select"
                  value={getTimeMusicEntry(t).mode}
                  on:change={e => onTimeMusicModeChange(t, e)}
                >
                  <option value="random">隨機</option>
                  <option value="index">指定</option>
                </select>
                {#if getTimeMusicEntry(t).mode === 'index'}
                  <select
                    class="track-select"
                    value={getTimeMusicEntry(t).index}
                    on:change={e => onTimeMusicIndexChange(t, e)}
                  >
                    {#each countdownMusics as track, i}
                      <option value={i}>{track.path.split(/[/\\]/).pop()}</option>
                    {/each}
                  </select>
                {/if}
              {/if}
            </div>
          {/each}
        </div>
      {/if}
    </div>

    <div class="form-group">
      <label>倒數音樂清單</label>
      <div class="inline-form" style="margin-bottom:8px">
        <button class="btn-primary btn-sm" on:click={addMusics}>＋ 新增音樂</button>
      </div>
      {#if countdownMusics.length > 0}
        <div class="music-list">
          {#each countdownMusics as track, i}
            <div class="music-row">
              <span class="music-name">{track.path.split(/[/\\]/).pop()}</span>
              <label class="random-toggle">
                <input type="checkbox" bind:checked={countdownMusics[i].in_random} />
                加入隨機
              </label>
              <button class="btn-sm btn-test" on:click={() => previewTrack(i)}>
                {testAudios[i] ? '⏹' : '▶'}
              </button>
              <button class="btn-sm btn-danger" on:click={() => removeTrack(i)}>✕</button>
            </div>
          {/each}
        </div>
        <div class="volume-row">
          <span class="volume-label">🔈</span>
          <input type="range" min="0" max="1" step="0.05" bind:value={countdownVolume} />
          <span class="volume-label">🔊</span>
          <span class="volume-value">{Math.round(countdownVolume * 100)}%</span>
        </div>
      {:else}
        <p class="empty-hint">尚未加入任何音樂</p>
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
  .music-list {
    display: flex;
    flex-direction: column;
    gap: 6px;
    margin-bottom: 10px;
  }
  .music-row {
    display: flex;
    align-items: center;
    gap: 8px;
    background: var(--bg-primary);
    border: 1px solid var(--border);
    border-radius: 6px;
    padding: 6px 10px;
  }
  .music-name {
    flex: 1;
    font-size: 13px;
    color: var(--text-secondary);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .random-toggle {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 12px;
    white-space: nowrap;
    cursor: pointer;
  }
  .empty-hint {
    font-size: 13px;
    color: var(--text-secondary);
    margin: 4px 0;
  }
  .time-tag-row {
    display: flex;
    align-items: center;
    gap: 6px;
    flex-wrap: wrap;
  }
  .mode-select, .track-select {
    font-size: 12px;
    padding: 2px 6px;
    border: 1px solid var(--border);
    border-radius: 6px;
    background: var(--bg-primary);
    cursor: pointer;
  }
  .track-select {
    max-width: 160px;
  }
  .btn-danger {
    background: var(--danger);
    color: white;
    border: none;
    border-radius: 6px;
    cursor: pointer;
  }
  .btn-danger:hover {
    background: var(--danger-hover);
  }
  .btn-test {
    background: var(--bg-primary);
    border: 1px solid var(--border);
    border-radius: 6px;
    cursor: pointer;
    font-size: 12px;
  }
  .btn-test:hover {
    background: var(--border);
  }
  .volume-row {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-top: 8px;
  }
  .volume-row input[type="range"] {
    flex: 1;
    max-width: 200px;
  }
  .volume-label {
    font-size: 14px;
  }
  .volume-value {
    font-size: 12px;
    color: var(--text-secondary);
    min-width: 36px;
  }
</style>
