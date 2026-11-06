<script lang="ts">
  import HomePage from '../pages/HomePage.svelte';
  import ContactBookPage from '../pages/ContactBookPage.svelte';
  export let contact = false;
  export let date = '';
  export let onDateChange: (date: string) => void = () => {};
  export let onExit: () => void;
</script>
<div class="display">
  <header><span>{contact ? '聯絡簿展示' : '班級總覽'}</span><button class="btn-outline" on:click={onExit}>結束展示（ESC）</button></header>
  {#if contact}<ContactBookPage presentation initialDate={date} />{:else}<HomePage presentation selectedDate={date} {onDateChange} />{/if}
</div>
<style>
  .display { position:fixed; inset:0; z-index:9000; background:var(--bg-primary); display:flex; flex-direction:column; padding:20px 32px; color:var(--text-primary); }
  header { display:flex; align-items:center; justify-content:space-between; gap:12px; padding-bottom:12px; }
  header span { font-weight:600; color:var(--text-secondary); }
  @media(max-width:700px) { .display { padding:12px; } }
</style>
