<script lang="ts">
  import HomePage from '../pages/HomePage.svelte';
  import ContactBookPage from '../pages/ContactBookPage.svelte';
  import MidtermPage from '../pages/MidtermPage.svelte';
  export let exam = false;
  export let contact = false;
  export let date = '';
  export let onDateChange: (date: string) => void = () => {};
  export let onExit: () => void;
</script>
<div class="display">
  <header><span>{exam ? '期中考展示' : contact ? '聯絡簿展示' : '班級總覽'}</span><button class="btn-outline" on:click={onExit}>結束展示（ESC）</button></header>
  {#if exam}<MidtermPage presentation initialDate={date} />{:else if contact}<ContactBookPage presentation initialDate={date} />{:else}<HomePage presentation selectedDate={date} {onDateChange} />{/if}
</div>
<style>
  .display { position:fixed; inset:0; z-index:9000; background:var(--bg-primary); display:flex; flex-direction:column; padding:clamp(8px, 1.8vh, 20px) clamp(12px, 2vw, 32px); container:workspace / inline-size; min-height:0; overflow:hidden; color:var(--text-primary); }
  header { flex-shrink:0; flex-wrap:wrap; display:flex; align-items:center; justify-content:space-between; gap:12px; padding-bottom:12px; }
  header span { font-weight:600; color:var(--text-secondary); }
  @media(max-width:700px) { .display { padding:12px; } }
</style>
