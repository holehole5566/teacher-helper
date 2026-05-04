<script lang="ts">
  import './style.css';
  import { onMount, onDestroy } from 'svelte';
  import Sidebar from './lib/components/Sidebar.svelte';
  import CountdownOverlay from './lib/components/CountdownOverlay.svelte';
  import DisplayMode from './lib/components/DisplayMode.svelte';
  import HomePage from './lib/pages/HomePage.svelte';
  import StudentsPage from './lib/pages/StudentsPage.svelte';
  import SettingsPage from './lib/pages/SettingsPage.svelte';
  import HolidaysPage from './lib/pages/HolidaysPage.svelte';
  import TimetablePage from './lib/pages/TimetablePage.svelte';
  import { ExportSchedule, GetSettings, SetFullscreen } from '../wailsjs/go/main/App';

  let currentPage = 'home';
  let homeRef: HomePage;
  let showCountdown = false;
  let showDisplay = false;
  let countdownTriggerTime = '';
  let countdownTimes: string[] = [];
  let triggeredToday = new Set<string>();
  let checkInterval: ReturnType<typeof setInterval>;

  async function loadCountdownTimes() {
    const s = await GetSettings();
    countdownTimes = s.countdown_times || [];
  }

  function checkCountdown() {
    const now = new Date();
    const hh = String(now.getHours()).padStart(2, '0');
    const mm = String(now.getMinutes()).padStart(2, '0');
    const currentHHMM = `${hh}:${mm}`;
    const todayKey = now.toDateString();

    if (triggeredToday.size > 0 && !Array.from(triggeredToday).some(k => k.startsWith(todayKey))) {
      triggeredToday = new Set();
    }

    for (const t of countdownTimes) {
      const [th, tm] = t.split(':').map(Number);
      let triggerH = th;
      let triggerM = tm - 1;
      if (triggerM < 0) { triggerM = 59; triggerH--; }
      if (triggerH < 0) triggerH = 23;
      const triggerHHMM = `${String(triggerH).padStart(2, '0')}:${String(triggerM).padStart(2, '0')}`;

      const key = `${todayKey}-${t}`;
      if (currentHHMM === triggerHHMM && !triggeredToday.has(key) && !showCountdown) {
        triggeredToday.add(key);
        startCountdown(t);
        break;
      }
    }
  }

  async function startCountdown(triggerTime: string) {
    countdownTriggerTime = triggerTime;
    showCountdown = true;
    await SetFullscreen(true);
  }

  async function onCountdownFinished() {
    showCountdown = false;
    await SetFullscreen(false);
  }

  async function enterDisplayMode() {
    showDisplay = true;
    await SetFullscreen(true);
  }

  async function exitDisplayMode() {
    showDisplay = false;
    await SetFullscreen(false);
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape' && showDisplay) {
      exitDisplayMode();
    }
  }

  async function handleNavigate(page: string) {
    if (page === 'export') {
      try {
        const path = await ExportSchedule();
        if (path) alert(`已匯出至: ${path}`);
      } catch (e) {
        alert('匯出失敗: ' + e);
      }
      return;
    }
    if (page === 'display') {
      enterDisplayMode();
      return;
    }
    currentPage = page;
    if (page === 'home' && homeRef) homeRef.refresh();
    if (page !== 'settings') loadCountdownTimes();
  }

  onMount(() => {
    loadCountdownTimes();
    checkInterval = setInterval(checkCountdown, 1000);
    window.addEventListener('keydown', handleKeydown);
  });

  onDestroy(() => {
    clearInterval(checkInterval);
    window.removeEventListener('keydown', handleKeydown);
  });
</script>

{#if showCountdown}
  <CountdownOverlay seconds={60} triggerTime={countdownTriggerTime} onFinished={onCountdownFinished} />
{/if}

{#if showDisplay}
  <DisplayMode />
{:else}
  <Sidebar {currentPage} onNavigate={handleNavigate} />

  <main class="main-content">
    {#if currentPage === 'home'}
      <HomePage bind:this={homeRef} />
    {:else if currentPage === 'students'}
      <StudentsPage />
    {:else if currentPage === 'settings'}
      <SettingsPage />
    {:else if currentPage === 'holidays'}
      <HolidaysPage />
    {:else if currentPage === 'timetable'}
      <TimetablePage />
    {/if}
  </main>
{/if}

<style>
  .main-content {
    flex: 1;
    display: flex;
    overflow: hidden;
  }
</style>
