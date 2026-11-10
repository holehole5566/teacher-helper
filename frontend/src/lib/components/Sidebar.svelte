<script lang="ts">
  import { ClipboardList, BookOpen, Users, Settings, Calendar, FileText, Monitor, ChartBar } from 'lucide-svelte';

  export let currentPage: string;
  export let onNavigate: (page: string) => void;

  const iconMap: Record<string, any> = {
    ClipboardList, BookOpen, Users, Settings, Calendar, FileText, Monitor, ChartBar
  };

  const navItems = [
    { id: 'home', label: '今日值日', icon: 'ClipboardList' },
    { id: 'contactBook', label: '聯絡簿', icon: 'BookOpen' },
    { id: 'midterm', label: '期中考', icon: 'FileText' },
    { id: 'timetable', label: '課表設定', icon: 'BookOpen' },
    { id: 'students', label: '學生管理', icon: 'Users' },
    { id: 'settings', label: '設定', icon: 'Settings' },
    { id: 'holidays', label: '假期管理', icon: 'Calendar' },
    { id: 'missingHomework', label: '作業未交', icon: 'FileText' },
  ];
</script>

<aside class="sidebar">
  <div class="sidebar-header">
    <h1 class="sidebar-title" aria-label="值日生助手"><span class="full-title">值日生助手</span><span class="short-title" aria-hidden="true">值日</span></h1>
    <p class="sidebar-subtitle">班級管理工具</p>
  </div>
  <nav class="sidebar-nav">
    {#each navItems as item}
      <button
        class="nav-item"
        aria-label={item.label}
        aria-current={currentPage === item.id ? 'page' : undefined}
        title={item.label}
        class:active={currentPage === item.id}
        on:click={() => onNavigate(item.id)}
      >
        <span class="nav-icon"><svelte:component this={iconMap[item.icon]} size={18} /></span>
        <span class="nav-label">{item.label}</span>
      </button>
    {/each}
  </nav>
  <div class="sidebar-footer">
    <button class="nav-item display-btn" aria-label="展示模式" title="展示模式" on:click={() => onNavigate('display')}>
      <span class="nav-icon"><Monitor size={18} /></span>
      <span class="nav-label">展示模式</span>
    </button>
    <button class="nav-item export-btn" aria-label="匯出排程" title="匯出排程" on:click={() => onNavigate('export')}>
      <span class="nav-icon"><ChartBar size={18} /></span>
      <span class="nav-label">匯出排程</span>
    </button>
  </div>
</aside>

<style>
  .sidebar {
    width: clamp(180px, 16vw, 260px);
    min-width: 0;
    min-height: 0;
    overflow: hidden;
    background: var(--bg-sidebar);
    display: flex;
    flex-direction: column;
    height: 100%;
    user-select: none;
  }
  .sidebar-header {
    flex-shrink: 0;
    padding: 20px 16px 16px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  }
  .sidebar-title {
    color: var(--text-sidebar-active);
    font-size: calc(var(--ui-font, 14px) * 16 / 14);
    font-weight: 700;
  }
  .short-title { display: none; }
  .sidebar-subtitle {
    color: var(--text-sidebar);
    font-size: calc(var(--ui-font, 14px) * 11 / 14);
    margin-top: 2px;
    opacity: 0.7;
  }
  .sidebar-nav {
    min-height: 0;
    overflow-y: auto;
    scrollbar-width: thin;
    flex: 1;
    padding: 8px;
    display: flex;
    flex-direction: column;
    gap: 2px;
  }
  .nav-item {
    display: flex;
    align-items: center;
    gap: 10px;
    width: 100%;
    padding: 10px 12px;
    border-radius: 8px;
    background: transparent;
    color: var(--text-sidebar);
    font-size: calc(var(--ui-font, 14px) * 13 / 14);
    text-align: left;
    transition: background 0.15s, color 0.15s;
  }
  .nav-item:hover {
    background: var(--bg-sidebar-hover);
    color: var(--text-sidebar-active);
  }
  .nav-item.active {
    background: var(--bg-sidebar-active);
    color: var(--text-sidebar-active);
  }
  .nav-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 1.5em;
    flex-shrink: 0;
  }
  .sidebar-footer {
    flex-shrink: 0;
    padding: 8px;
    border-top: 1px solid rgba(255, 255, 255, 0.08);
  }
  .export-btn {
    opacity: 0.8;
  }
  .export-btn:hover {
    opacity: 1;
  }
  .display-btn {
    opacity: 0.8;
  }
  .display-btn:hover {
    opacity: 1;
    background: var(--bg-sidebar-active);
    color: var(--text-sidebar-active);
  }
  @media (max-width: 850px) {
    .sidebar { width: 68px; }
    .sidebar-header { padding: 16px 4px; text-align: center; }
    .sidebar-title { font-size: 14px; }
    .full-title { display: none; }
    .short-title { display: inline; }
    .sidebar-subtitle, .nav-label { display: none; }
    .nav-item { justify-content: center; padding: 12px 8px; min-height: 44px; }
  }
  @media (max-height: 480px) {
    .sidebar-header { padding-top: 8px; padding-bottom: 8px; }
    .sidebar-subtitle { display: none; }
    .sidebar-footer { padding-top: 4px; padding-bottom: 4px; }
  }
</style>
