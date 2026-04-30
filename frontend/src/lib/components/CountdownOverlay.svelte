<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { GetCountdownMusicData, GetSettings } from '../../../wailsjs/go/main/App';

  export let seconds: number = 60;
  export let onFinished: () => void = () => {};

  let remaining = seconds;
  let audio: HTMLAudioElement | null = null;

  const interval = setInterval(() => {
    remaining--;
    if (remaining <= 0) {
      clearInterval(interval);
      stopMusic();
      onFinished();
    }
  }, 1000);

  function stopMusic() {
    if (audio) {
      audio.pause();
      audio = null;
    }
  }

  onMount(async () => {
    try {
      const [dataUrl, settings] = await Promise.all([GetCountdownMusicData(), GetSettings()]);
      if (dataUrl) {
        audio = new Audio(dataUrl);
        audio.volume = settings.countdown_volume > 0 ? settings.countdown_volume : 0.5;
        audio.play().catch(() => {});
      }
    } catch {}
  });

  onDestroy(() => {
    clearInterval(interval);
    stopMusic();
  });

  $: minutes = Math.floor(remaining / 60);
  $: secs = remaining % 60;
  $: display = minutes > 0
    ? `${minutes}:${String(secs).padStart(2, '0')}`
    : `${secs}`;
  $: progress = remaining / seconds;
</script>

<div class="overlay">
  <div class="countdown-container">
    <div class="ring-wrapper">
      <svg viewBox="0 0 200 200" class="ring">
        <circle cx="100" cy="100" r="90" class="ring-bg" />
        <circle
          cx="100" cy="100" r="90"
          class="ring-progress"
          style="stroke-dashoffset: {565.48 * (1 - progress)}"
        />
      </svg>
      <span class="time">{display}</span>
    </div>
    <p class="label">上課倒數</p>
    <p class="sublabel">請回到座位準備上課</p>
  </div>
</div>

<style>
  .overlay {
    position: fixed;
    inset: 0;
    z-index: 99999;
    background: rgba(15, 23, 42, 0.97);
    display: flex;
    align-items: center;
    justify-content: center;
    user-select: none;
    cursor: default;
  }
  .countdown-container {
    text-align: center;
    color: white;
  }
  .ring-wrapper {
    position: relative;
    width: 280px;
    height: 280px;
    margin: 0 auto 32px;
  }
  .ring {
    width: 100%;
    height: 100%;
    transform: rotate(-90deg);
  }
  .ring-bg {
    fill: none;
    stroke: rgba(255, 255, 255, 0.1);
    stroke-width: 8;
  }
  .ring-progress {
    fill: none;
    stroke: #3b82f6;
    stroke-width: 8;
    stroke-linecap: round;
    stroke-dasharray: 565.48;
    transition: stroke-dashoffset 1s linear;
  }
  .time {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 80px;
    font-weight: 700;
    font-variant-numeric: tabular-nums;
    letter-spacing: -2px;
  }
  .label {
    font-size: 28px;
    font-weight: 600;
    margin-bottom: 8px;
  }
  .sublabel {
    font-size: 16px;
    opacity: 0.6;
  }
</style>
