<script lang="ts">
  import './style.css';
  import { onMount, onDestroy } from 'svelte';
  import Sidebar from './lib/components/Sidebar.svelte';
  import CountdownOverlay from './lib/components/CountdownOverlay.svelte';
  import DisplayMode from './lib/components/DisplayMode.svelte';
  import PasswordModal from './lib/components/PasswordModal.svelte';
  import DraftFeedback from './lib/components/DraftFeedback.svelte';
  import { requestNavigation, saving, leavePrompt } from './lib/stores/draft';
  import HomePage from './lib/pages/HomePage.svelte';
  import StudentsPage from './lib/pages/StudentsPage.svelte';
  import SettingsPage from './lib/pages/SettingsPage.svelte';
  import HolidaysPage from './lib/pages/HolidaysPage.svelte';
  import TimetablePage from './lib/pages/TimetablePage.svelte';
  import MissingHomeworkPage from './lib/pages/MissingHomeworkPage.svelte';
  import { ExportSchedule, GetSettings, SetFullscreen, ReportError, DebugLog } from '../wailsjs/go/main/App';
  import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime';
  import { checkPasswordStatus } from './lib/stores/auth';

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
      requestNavigation(() => { enterDisplayMode(); });
      return;
    }
    if (page === currentPage) return;
    requestNavigation(() => { currentPage = page; });
  }

  onMount(() => {
    checkPasswordStatus();
    EventsOn('countdown-trigger', onCountdownTrigger);
    window.addEventListener('keydown', handleKeydown);
  });

  onDestroy(() => {
    EventsOff('countdown-trigger');
    window.removeEventListener('keydown', handleKeydown);
  });
</script>

<PasswordModal />
<DraftFeedback />

{#if showCountdown}
  <CountdownOverlay seconds={60} triggerTime={countdownTriggerTime} onFinished={onCountdownFinished} />
{/if}

{#if showDisplay}
  <DisplayMode />
{:else}
  <div class="navigation" inert={$saving || $leavePrompt}>
    <Sidebar {currentPage} onNavigate={handleNavigate} />
  </div>

  <main class="main-content" inert={$saving || $leavePrompt}>
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
  .navigation { height: 100%; }
  .main-content {
    flex: 1;
    display: flex;
    overflow: hidden;
  }
</style>
