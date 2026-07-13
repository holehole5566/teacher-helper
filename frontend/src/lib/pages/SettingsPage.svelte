<script lang="ts">
  import DutySettings from '../components/settings/DutySettings.svelte';
  import ScheduleSettings from '../components/settings/ScheduleSettings.svelte';
  import MusicSettings from '../components/settings/MusicSettings.svelte';
  import SystemSettings from '../components/settings/SystemSettings.svelte';
  import { ClipboardList, Clock, Music, Wrench } from 'lucide-svelte';

  let activeTab = 'duty';

  const iconMap: Record<string, any> = {
    ClipboardList, Clock, Music, Wrench
  };

  const tabs = [
    { id: 'duty', label: '值日/抬餐', icon: 'ClipboardList' },
    { id: 'schedule', label: '課程時間', icon: 'Clock' },
    { id: 'music', label: '倒數音樂', icon: 'Music' },
    { id: 'system', label: '系統', icon: 'Wrench' },
  ];
</script>

<div class="page">
  <h2 class="page-title">設定</h2>

  <div class="settings-layout">
    <div class="tab-bar">
      {#each tabs as tab}
        <button
          class="tab-btn"
          class:active={activeTab === tab.id}
          on:click={() => activeTab = tab.id}
        >
          <span class="tab-icon"><svelte:component this={iconMap[tab.icon]} size={16} /></span>
          <span class="tab-label">{tab.label}</span>
        </button>
      {/each}
    </div>

    <div class="tab-content">
      {#if activeTab === 'duty'}
        <DutySettings />
      {:else if activeTab === 'schedule'}
        <ScheduleSettings />
      {:else if activeTab === 'music'}
        <MusicSettings />
      {:else if activeTab === 'system'}
        <SystemSettings />
      {/if}
    </div>
  </div>
</div>

<style>
  .settings-layout {
    display: flex;
    flex-direction: column;
    gap: 0;
    max-width: 640px;
  }
  .tab-bar {
    display: flex;
    gap: 4px;
    padding: 4px;
    background: var(--bg-secondary);
    border-radius: var(--radius-lg) var(--radius-lg) 0 0;
    box-shadow: var(--shadow);
    border-bottom: 1px solid var(--border);
  }
  .tab-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 10px 16px;
    border-radius: 8px;
    background: transparent;
    color: var(--text-secondary);
    font-size: 13px;
    font-weight: 500;
    transition: all 0.15s;
  }
  .tab-btn:hover {
    background: var(--bg-primary);
    color: var(--text-primary);
  }
  .tab-btn.active {
    background: var(--accent);
    color: white;
  }
  .tab-icon {
    display: inline-flex;
    align-items: center;
  }
  .tab-content {
    background: var(--bg-secondary);
    border-radius: 0 0 var(--radius-lg) var(--radius-lg);
    box-shadow: var(--shadow);
    padding: 24px;
  }
</style>
