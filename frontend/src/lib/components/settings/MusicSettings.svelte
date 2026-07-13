<script lang="ts">
  import { onMount } from 'svelte';
  import { GetSettings, SaveSettings, SelectCountdownMusics, GetCountdownMusicData, ValidateRandomPool } from '../../../../wailsjs/go/main/App';
  import { verifyAndRun } from '../../stores/auth';
  import { Square, Play, X, Volume1, Volume2 } from 'lucide-svelte';

  type MusicTrack = { path: string; in_random: boolean };
  type CountdownTimeMusic = { time: string; mode: string; index: number };

  let countdownTimes: string[] = [];
  let countdownMusics: MusicTrack[] = [];
  let countdownTimeMusicMap: CountdownTimeMusic[] = [];
  let countdownVolume = 0.5;
  let audioOutputDevice = '';
  let audioDevices: { deviceId: string; label: string }[] = [];
  let testAudios: (HTMLAudioElement | null)[] = [];
  let saved = false;

  let timeMusicSettings: Record<string, { mode: string; index: number }> = {};

  async function refreshAudioDevices() {
    try {
      try {
        const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
        stream.getTracks().forEach(t => t.stop());
      } catch {}
      const devices = await navigator.mediaDevices.enumerateDevices();
      audioDevices = devices
        .filter(d => d.kind === 'audiooutput' && !d.label.startsWith('通訊') && !d.label.startsWith('預設'))
        .map(d => ({ deviceId: d.deviceId, label: d.label || `未知裝置 (${d.deviceId.slice(0, 8)})` }));
    } catch {
      audioDevices = [];
    }
  }

  async function loadSettings() {
    const s = await GetSettings();
    countdownTimes = (s.countdown_times || []).slice().sort();
    countdownMusics = (s.countdown_musics || []).map((t: any) => ({ path: t.path, in_random: t.in_random }));
    countdownTimeMusicMap = (s.countdown_time_music_map || []).map((m: any) => ({ time: m.time, mode: m.mode || 'random', index: m.index || 0 }));
    countdownVolume = s.countdown_volume > 0 ? s.countdown_volume : 0.5;
    audioOutputDevice = s.audio_output_device || '';
    await refreshAudioDevices();
    testAudios = countdownMusics.map(() => null);
    rebuildTimeMusicSettings();
  }

  function rebuildTimeMusicSettings() {
    const map: Record<string, { mode: string; index: number }> = {};
    for (const t of countdownTimes) {
      const entry = countdownTimeMusicMap.find(m => m.time === t);
      map[t] = entry ? { mode: entry.mode, index: entry.index } : { mode: 'random', index: 0 };
    }
    timeMusicSettings = map;
  }

  function syncMapFromSettings() {
    countdownTimeMusicMap = Object.entries(timeMusicSettings).map(([time, v]) => ({
      time, mode: v.mode, index: v.index
    }));
  }

  function onTimeMusicModeChange(t: string, e: Event) {
    timeMusicSettings[t] = { ...timeMusicSettings[t], mode: (e.target as HTMLSelectElement).value };
    timeMusicSettings = { ...timeMusicSettings };
    syncMapFromSettings();
  }

  function onTimeMusicIndexChange(t: string, e: Event) {
    timeMusicSettings[t] = { ...timeMusicSettings[t], index: Number((e.target as HTMLSelectElement).value) };
    timeMusicSettings = { ...timeMusicSettings };
    syncMapFromSettings();
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
    rebuildTimeMusicSettings();
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
      const deviceId = audioOutputDevice || 'default';
      try { await (a as any).setSinkId(deviceId); } catch {}
      a.onended = () => { testAudios[i] = null; testAudios = [...testAudios]; };
      testAudios[i] = a;
      testAudios = [...testAudios];
      a.play().catch(() => {});
    } catch (e: any) {
      alert('無法播放：' + (e?.message || '檔案可能已被移動或刪除'));
    }
  }

  async function handleSave() {
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

    try {
      const current = await GetSettings();
      await verifyAndRun(async () => {
        await SaveSettings({
          ...current,
          countdown_volume: countdownVolume,
          countdown_musics: countdownMusics,
          countdown_time_music_map: countdownTimeMusicMap,
          audio_output_device: audioOutputDevice,
        });
      }, '儲存音樂設定');
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
  <h3 class="section-title">倒數音樂</h3>
  <p class="section-desc">管理上課前倒數播放的音樂</p>

  <div class="form-group">
    <label>音樂清單</label>
    <button class="btn-primary btn-sm" on:click={addMusics}>＋ 新增音樂</button>
    {#if countdownMusics.length > 0}
      <div class="music-list">
        {#each countdownMusics as track, i}
          <div class="music-row">
            <span class="music-name">{track.path.split(/[/\\]/).pop()}</span>
            <label class="random-toggle">
              <input type="checkbox" bind:checked={countdownMusics[i].in_random} />
              隨機
            </label>
            <button class="btn-sm btn-test" on:click={() => previewTrack(i)}>
              {#if testAudios[i]}
                <Square size={12} />
              {:else}
                <Play size={12} />
              {/if}
            </button>
            <button class="btn-sm btn-danger" on:click={() => removeTrack(i)}><X size={12} /></button>
          </div>
        {/each}
      </div>
      <div class="volume-row">
        <span class="vol-icon"><Volume1 size={16} /></span>
        <input type="range" min="0" max="1" step="0.05" bind:value={countdownVolume} />
        <span class="vol-icon"><Volume2 size={16} /></span>
        <span class="vol-value">{Math.round(countdownVolume * 100)}%</span>
      </div>
    {:else}
      <p class="empty-hint">尚未加入任何音樂</p>
    {/if}
  </div>

  {#if countdownTimes.length > 0 && countdownMusics.length > 0}
    <div class="form-group">
      <label>各時段音樂配對</label>
      <div class="time-music-list">
        {#each countdownTimes as t}
          {#if timeMusicSettings[t]}
            <div class="time-music-row">
              <span class="tm-time">{t}</span>
              <select
                class="tm-select"
                value={timeMusicSettings[t].mode}
                on:change={e => onTimeMusicModeChange(t, e)}
              >
                <option value="random">隨機</option>
                <option value="index">指定</option>
                <option value="none">不播放</option>
              </select>
              {#if timeMusicSettings[t].mode === 'index'}
                <select
                  class="tm-select tm-track"
                  value={timeMusicSettings[t].index}
                  on:change={e => onTimeMusicIndexChange(t, e)}
                >
                  {#each countdownMusics as track, i}
                    <option value={i}>{track.path.split(/[/\\]/).pop()}</option>
                  {/each}
                </select>
              {/if}
            </div>
          {/if}
        {/each}
      </div>
    </div>
  {/if}

  <div class="form-group">
    <label>音訊輸出裝置</label>
    <div class="inline-form">
      <select bind:value={audioOutputDevice} style="flex:1">
        <option value="">跟隨系統預設</option>
        {#each audioDevices as device}
          <option value={device.deviceId}>{device.label}</option>
        {/each}
      </select>
      <button class="btn-primary btn-sm" on:click={refreshAudioDevices}>重新偵測</button>
    </div>
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
  .music-list {
    display: flex;
    flex-direction: column;
    gap: 6px;
    margin-top: 10px;
  }
  .music-row {
    display: flex;
    align-items: center;
    gap: 8px;
    background: var(--bg-primary);
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 8px 12px;
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
    display: inline-flex;
    align-items: center;
    justify-content: center;
  }
  .btn-test:hover {
    background: var(--border);
  }
  .volume-row {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-top: 10px;
  }
  .volume-row input[type="range"] {
    flex: 1;
    max-width: 200px;
  }
  .vol-icon {
    display: inline-flex;
    align-items: center;
    color: var(--text-secondary);
  }
  .vol-value {
    font-size: 12px;
    color: var(--text-secondary);
    min-width: 36px;
  }
  .empty-hint {
    font-size: 13px;
    color: var(--text-secondary);
    margin: 8px 0;
  }
  .time-music-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-top: 8px;
  }
  .time-music-row {
    display: flex;
    align-items: center;
    gap: 10px;
    background: var(--bg-primary);
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 8px 12px;
  }
  .tm-time {
    font-size: 13px;
    font-weight: 600;
    min-width: 48px;
  }
  .tm-select {
    font-size: 12px;
    padding: 4px 8px;
    border: 1px solid var(--border);
    border-radius: 6px;
    background: white;
    cursor: pointer;
  }
  .tm-track {
    max-width: 180px;
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
