<script lang="ts">
  export let currentPage: string;
  export let onNavigate: (page: string) => void;

  const navItems = [
    { id: 'home', label: '今日值日', icon: '📋' },
    { id: 'timetable', label: '課表設定', icon: '📚' },
    { id: 'students', label: '學生管理', icon: '👥' },
    { id: 'settings', label: '設定', icon: '⚙️' },
    { id: 'holidays', label: '假期管理', icon: '📅' },
    { id: 'missingHomework', label: '作業未交', icon: '📝' },
  ];
</script>

<aside class="sidebar">
  <div class="sidebar-header">
    <h1 class="sidebar-title">值日生助手</h1>
    <p class="sidebar-subtitle">班級管理工具</p>
  </div>
  <nav class="sidebar-nav">
    {#each navItems as item}
      <button
        class="nav-item"
        class:active={currentPage === item.id}
        on:click={() => onNavigate(item.id)}
      >
        <span class="nav-icon">{item.icon}</span>
        <span class="nav-label">{item.label}</span>
      </button>
    {/each}
  </nav>
  <div class="sidebar-footer">
    <button class="nav-item display-btn" on:click={() => onNavigate('display')}>
      <span class="nav-icon">🖥️</span>
      <span class="nav-label">展示模式</span>
    </button>
    <button class="nav-item export-btn" on:click={() => onNavigate('export')}>
      <span class="nav-icon">📊</span>
      <span class="nav-label">匯出排程</span>
    </button>
  </div>
</aside>

<style>
  .sidebar {
    width: 200px;
    min-width: 200px;
    background: var(--bg-sidebar);
    display: flex;
    flex-direction: column;
    height: 100%;
    user-select: none;
  }
  .sidebar-header {
    padding: 20px 16px 16px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  }
  .sidebar-title {
    color: var(--text-sidebar-active);
    font-size: 16px;
    font-weight: 700;
  }
  .sidebar-subtitle {
    color: var(--text-sidebar);
    font-size: 11px;
    margin-top: 2px;
    opacity: 0.7;
  }
  .sidebar-nav {
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
    font-size: 13px;
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
    font-size: 16px;
    width: 20px;
    text-align: center;
  }
  .sidebar-footer {
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
</style>
