<script lang="ts">
  import { ChevronLeft, ChevronRight } from 'lucide-svelte';
  export let dates: string[] = [];
  export let date = '';
  export let busy = false;
  export let onChange: (date: string) => void;
  $: index = dates.indexOf(date);
</script>
<nav class="history-nav" aria-label="瀏覽最近七天">
  <button class="btn-outline" disabled={busy || index < 0 || index >= dates.length - 1} on:click={() => onChange(dates[index + 1])} aria-label="前一天"><ChevronLeft size={16} /></button>
  <label>日期 <select value={date} disabled={busy} on:change={(e) => onChange(e.currentTarget.value)}>{#each dates as day}<option value={day}>{day}{day === dates[0] ? '（今天）' : ''}</option>{/each}</select></label>
  <button class="btn-outline" disabled={busy || index <= 0} on:click={() => onChange(dates[index - 1])} aria-label="後一天"><ChevronRight size={16} /></button>
  <button class="btn-outline" disabled={busy || date === dates[0]} on:click={() => onChange(dates[0])}>回到今天</button>
</nav>
<style>
  .history-nav { display:flex; align-items:center; flex-wrap:wrap; gap:8px; margin-bottom:16px; }
  button { display:inline-flex; align-items:center; min-height:36px; }
  label { display:flex; align-items:center; gap:8px; }
  button:disabled { opacity:.45; cursor:default; }
</style>
