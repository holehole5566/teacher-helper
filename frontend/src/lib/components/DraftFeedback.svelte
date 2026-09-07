<script lang="ts">
  import { dirty, saving, notice, leavePrompt, cancelNavigation, resolveNavigation } from '../stores/draft';

  function focusDialog(node: HTMLElement) {
    const previous = document.activeElement as HTMLElement | null;
    const cancel = node.querySelector<HTMLButtonElement>('[data-cancel]');
    cancel?.focus();
    function keydown(event: KeyboardEvent) {
      if (event.key === 'Escape' && !$saving) {
        event.preventDefault();
        event.stopPropagation();
        cancelNavigation();
      }
      if (event.key === 'Tab') {
        const buttons = Array.from(node.querySelectorAll<HTMLButtonElement>('button:not(:disabled)'));
        const first = buttons[0];
        const last = buttons[buttons.length - 1];
        if (event.shiftKey && document.activeElement === first) { event.preventDefault(); last?.focus(); }
        if (!event.shiftKey && document.activeElement === last) { event.preventDefault(); first?.focus(); }
      }
    }
    node.addEventListener('keydown', keydown);
    return { destroy() { node.removeEventListener('keydown', keydown); if (previous?.isConnected) previous.focus(); } };
  }
</script>

{#if $dirty || $saving || $notice}
  <div class="feedback" role="status" aria-live="polite">
    {#if $saving}正在驗證／儲存，請稍候…{:else if $dirty}尚未儲存變更{:else}{$notice}{/if}
  </div>
{/if}

{#if $leavePrompt}
  <div class="backdrop">
    <section class="dialog" role="dialog" aria-modal="true" aria-labelledby="draft-title" use:focusDialog>
      <h2 id="draft-title">尚有未儲存的變更</h2>
      <p>離開此頁會遺失修改，是否先儲存？</p>
      <div class="buttons">
        <button class="btn-outline" data-cancel disabled={$saving} on:click={cancelNavigation}>取消</button>
        <button class="btn-danger" disabled={$saving} on:click={() => resolveNavigation(false)}>捨棄變更</button>
        <button class="btn-primary" disabled={$saving} on:click={() => resolveNavigation(true)}>{$saving ? '處理中…' : '儲存並離開'}</button>
      </div>
    </section>
  </div>
{/if}

<style>
  .feedback { position: fixed; bottom: 16px; left: 50%; transform: translateX(-50%); z-index: 9500; background: #1e293b; color: white; padding: 10px 20px; border-radius: 10px; box-shadow: var(--shadow-lg); pointer-events: none; }
  .backdrop { position: fixed; inset: 0; z-index: 9800; background: rgba(15,23,42,.45); display: flex; align-items: center; justify-content: center; }
  .dialog { width: min(460px, 92vw); background: white; padding: 24px; border-radius: 12px; box-shadow: var(--shadow-lg); }
  h2 { font-size: 20px; margin-bottom: 12px; }
  p { color: var(--text-secondary); margin-bottom: 24px; }
  .buttons { display: flex; justify-content: flex-end; flex-wrap: wrap; gap: 8px; }
</style>
