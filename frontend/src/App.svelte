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
  import MissingHomeworkPage from './lib/pages/MissingHomeworkPage.svelte';
  import { ExportSchedule, GetSettings, SetFullscreen, ReportError, DebugLog } from '../wailsjs/go/main/App';
  import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime';

  let currentPage = 'home';
  let homeRef: HomePage;
  let showCountdown = false;
  let showDisplay = false;
  let countdownTriggerTime = '';

  function onCountdownTrigger(triggerTime: string) {
    if (showCountdown) return;
    DebugLog(`[App] Countdown TRIGGERED from backend: triggerTime=${triggerTime}`);
    startCountdown(triggerTime);
  }

  async function startCountdown(triggerTime: string) {
    DebugLog(`[App] startCountdown called, triggerTime=${triggerTime}, showDisplay=${showDisplay}`);
    countdownTriggerTime = triggerTime;
    showCountdown = true;
    if (!showDisplay) {
      await SetFullscreen(true);
    }
  }

  async function onCountdownFinished() {
    showCountdown = false;
    if (!showDisplay) {
      await SetFullscreen(false);
    }
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
      } catch (e: any) {
        ReportError(`排班匯出失敗：${e?.message || e}`);
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
  }

  onMount(() => {
    EventsOn('countdown-trigger', onCountdownTrigger);
    window.addEventListener('keydown', handleKeydown);
  });

  onDestroy(() => {
    EventsOff('countdown-trigger');
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
    {:else if currentPage === 'missingHomework'}
      <MissingHomeworkPage />
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
