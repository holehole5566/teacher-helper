<script lang="ts">
  import { UtensilsCrossed, Pencil } from 'lucide-svelte';
  export let assignments: Array<{
    student: { seat_number: number; name: string };
    bucket: string;
  }>;
  export let onReplace: (index: number) => void = () => {};
  export let modified = false;
</script>

<div class="card lunch-card">
  <div class="card-header">
    <span class="card-icon"><UtensilsCrossed size={20} /></span>
    <h3 class="card-title">抬餐負責人</h3>
    {#if modified}<span class="badge badge-muted">已手動調整</span>{/if}
  </div>
  <div class="card-body">
    {#if assignments && assignments.length > 0}
      <div class="assignment-list">
        {#each assignments as a, i}
          <button class="assignment-row" title="點擊替換" on:click={() => onReplace(i)}>
            <div class="student-info">
              <span class="seat">{a.student.seat_number}號</span>
              <span class="name">{a.student.name}</span>
              <span class="edit-hint"><Pencil size={14} /></span>
            </div>
            <span class="bucket-tag">{a.bucket}</span>
          </button>
        {/each}
      </div>
    {:else}
      <p class="no-data">今日無抬餐</p>
    {/if}
  </div>
</div>

<style>
  .lunch-card {
    border-left: 4px solid var(--warning);
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
  .assignment-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  .assignment-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 14px;
    background: #fffbeb;
    border-radius: 8px;
    border: none;
    width: 100%;
    cursor: pointer;
    font-family: inherit;
    transition: background 0.15s, box-shadow 0.15s;
  }
  .assignment-row:hover {
    background: #fef3c7;
    box-shadow: var(--shadow);
  }
  .student-info {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .seat {
    font-size: 13px;
    color: var(--warning);
    font-weight: 600;
  }
  .name {
    font-size: 15px;
    font-weight: 600;
    color: var(--text-primary);
  }
  .edit-hint {
    font-size: 14px;
    color: var(--text-secondary);
    opacity: 0;
    transition: opacity 0.15s;
  }
  .assignment-row:hover .edit-hint {
    opacity: 1;
  }
  .bucket-tag {
    background: #fef3c7;
    color: #92400e;
    padding: 4px 10px;
    border-radius: 20px;
    font-size: 12px;
    font-weight: 600;
  }
  .no-data {
    color: var(--text-secondary);
    font-size: 14px;
  }
</style>
