<script lang="ts">
  import { Brush, Pencil } from 'lucide-svelte';
  export let students: Array<{ seat_number: number; name: string }>;
  export let onReplace: (index: number) => void = () => {};
  export let modified = false;
</script>

<div class="card duty-card">
  <div class="card-header">
    <span class="card-icon"><Brush size={20} /></span>
    <h3 class="card-title">今日值日生</h3>
    {#if modified}<span class="badge badge-muted">已手動調整</span>{/if}
  </div>
  <div class="card-body">
    {#if students && students.length > 0}
      <div class="student-list">
        {#each students as student, i}
          <button class="student-chip" title="點擊替換" on:click={() => onReplace(i)}>
            <span class="seat">{student.seat_number}號</span>
            <span class="name">{student.name}</span>
            <span class="edit-hint"><Pencil size={14} /></span>
          </button>
        {/each}
      </div>
    {:else}
      <p class="no-data">今日無值日</p>
    {/if}
  </div>
</div>

<style>
  .duty-card {
    border-left: 4px solid var(--accent);
  }
  .card-header {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 16px;
  }
  .card-icon {
    display: inline-flex;
    align-items: center;
    color: var(--text-secondary);
  }
  .card-title {
    font-size: 15px;
    font-weight: 600;
  }
  .student-list {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
  }
  .student-chip {
    display: flex;
    align-items: center;
    gap: 8px;
    background: #eff6ff;
    border-radius: 10px;
    padding: 12px 18px;
    cursor: pointer;
    transition: background 0.15s, box-shadow 0.15s;
    border: none;
    font-family: inherit;
  }
  .student-chip:hover {
    background: #dbeafe;
    box-shadow: var(--shadow);
  }
  .seat {
    font-size: 13px;
    color: var(--accent);
    font-weight: 600;
  }
  .name {
    font-size: 18px;
    font-weight: 700;
    color: var(--text-primary);
  }
  .edit-hint {
    font-size: 14px;
    color: var(--text-secondary);
    opacity: 0;
    transition: opacity 0.15s;
  }
  .student-chip:hover .edit-hint {
    opacity: 1;
  }
  .no-data {
    color: var(--text-secondary);
    font-size: 14px;
  }
</style>
