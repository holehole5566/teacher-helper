<script lang="ts">
  import { X } from 'lucide-svelte';
  export let students: Array<{ seat_number: number; name: string }> = [];
  export let onPick: (student: { seat_number: number; name: string }) => void = () => {};
  export let onClose: () => void = () => {};
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div class="backdrop" on:click|self={onClose}>
  <div class="picker">
    <div class="picker-header">
      <h3>選擇學生</h3>
      <button class="btn-icon" on:click={onClose}><X size={18} /></button>
    </div>
    <div class="picker-list">
      {#each students as s}
        <button class="pick-item" on:click={() => onPick(s)}>
          <span class="seat">{s.seat_number}號</span>
          <span class="name">{s.name}</span>
        </button>
      {/each}
    </div>
  </div>
</div>

<style>
  .backdrop {
    position: fixed;
    inset: 0;
    z-index: 1000;
    background: rgba(0,0,0,0.3);
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .picker {
    background: var(--bg-secondary);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-lg);
    width: 320px;
    max-height: 400px;
    display: flex;
    flex-direction: column;
  }
  .picker-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 14px 16px;
    border-bottom: 1px solid var(--border);
  }
  .picker-header h3 {
    font-size: 15px;
    font-weight: 600;
  }
  .picker-list {
    overflow-y: auto;
    padding: 8px;
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
  }
  .pick-item {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 8px 12px;
    border-radius: var(--radius);
    background: var(--bg-primary);
    font-size: 13px;
    transition: background 0.15s;
  }
  .pick-item:hover {
    background: #dbeafe;
  }
  .seat {
    color: var(--accent);
    font-weight: 600;
    font-size: 12px;
  }
  .name {
    font-weight: 500;
  }
</style>
