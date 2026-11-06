<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { Plus, Trash2, ArrowUp, ArrowDown, Save, Monitor } from 'lucide-svelte';
  import { GetDailyDates, GetDailyOverview, SaveContactBook } from '../../../wailsjs/go/main/App';
  import { verifyAndRun } from '../stores/auth';
  import { draftGuard, markDraftSaved, requestNavigation, dirty, saving } from '../stores/draft';
  import HistoryNav from '../components/HistoryNav.svelte';
  export let presentation = false;
  export let initialDate = '';
  export let onDisplay: (date: string) => void = () => {};
  let dates: string[] = [];
  let date = '';
  let items: string[] = [];
  let legacyMorningReview = '';
  let nextDayReview = '';
  let ready = false;
  let loaded = false;
  let rolloverDates: string[] = [];
  let error = '';
  let clockInterval: ReturnType<typeof setInterval>;
  $: historical = date !== dates[0];
  $: draft = JSON.stringify({date, items, nextDayReview});
  async function load(next: string) {
    ready = false; loaded = false; error = ''; rolloverDates = []; date = next; items = []; legacyMorningReview = ''; nextDayReview = '';
    try {
      const record = await GetDailyOverview(date);
      items = record.contact_items || []; legacyMorningReview = record.morning_review || ''; nextDayReview = record.next_day_review || ''; loaded = true;
    } catch (e) { error = String(e); }
    finally { ready = true; }
  }
  async function save() {
    error = '';
    try {
      await verifyAndRun(() => SaveContactBook(date, items.filter(item => item.trim()), legacyMorningReview, nextDayReview.trim()), '儲存聯絡簿');
      markDraftSaved();
    } catch (e) { error = String(e); throw e; }
  }
  function move(index: number, direction: number) {
    const next = [...items];
    [next[index], next[index + direction]] = [next[index + direction], next[index]];
    items = next;
  }
  function changeDate(next: string) { requestNavigation(() => { load(next); }); }
  let checking = false;
  async function checkDate() {
    if (checking) return;
    checking = true;
    try {
      const next = await GetDailyDates();
      if (dates[0] !== next[0]) {
        // Keep unsaved previous-day text visible; saving refuses to write it to a new date.
        if ($dirty) { rolloverDates = next; error = '日期已變更。請先複製需要的文字，再切換到今天；昨天的草稿無法儲存。'; return; }
        const wasToday = date === dates[0]; dates = next;
        if (wasToday || !dates.includes(date)) await load(dates[0]);
      }
    } finally { checking = false; }
  }
  onMount(() => {
    (async () => {
      try { dates = await GetDailyDates(); await load(dates.includes(initialDate) ? initialDate : dates[0]); }
      catch(e) { error = String(e); }
    })();
    clockInterval = setInterval(() => { checkDate().catch(() => {}); }, 10000);
  });
  onDestroy(() => clearInterval(clockInterval));
</script>

<div class="page contact-page" class:presentation>
  <div class="heading">
    <div><h2 class="page-title">聯絡簿</h2><p class="description">{presentation ? date : '每日聯絡事項與明日複習考，保留最近七個日曆天。'}</p></div>
    {#if !presentation}<button class="btn-primary" disabled={!ready || !loaded || $saving} on:click={() => requestNavigation(() => onDisplay(date))}><Monitor size={16} /> 獨立展示</button>{/if}
  </div>
  <HistoryNav {dates} {date} busy={!ready || $saving} onChange={changeDate} />
  {#if error}<p class="error" role="alert">{error}</p>{/if}
  {#if rolloverDates.length}<button class="btn-outline" on:click={() => requestNavigation(() => { dates = rolloverDates; load(dates[0]); })}>切換到今天</button>{/if}
  {#if !ready}<div class="empty-state">載入中…</div>
  {:else if !loaded}<div class="empty-state">此日期沒有可顯示的聯絡簿。<button class="btn-outline" on:click={() => load(date)}>重新載入</button></div>
  {:else if presentation || historical}
    {#if historical}<p class="description">歷史快照・唯讀</p>{/if}
    <section class="reading">
      {#if nextDayReview}<div class="review"><h3>明日複習考</h3><p>{nextDayReview}</p></div>{/if}
      <h3>聯絡事項</h3>
      {#if items.length}<ol>{#each items as item}<li>{item}</li>{/each}</ol>
      {:else}<p class="description">此日尚無聯絡事項</p>{/if}
    </section>
  {:else}
    {#key date}
      <form use:draftGuard={{value: draft, ready, save}} on:submit|preventDefault={() => { save().catch(() => {}); }}>
        <section class="editor-section">
          <label for="next-day-review">明日複習考</label>
          <p class="description field-help" id="next-day-review-help">今天抄在聯絡簿，明天會顯示在首頁。沒有複習考可留白。</p>
          <textarea id="next-day-review" rows="2" bind:value={nextDayReview} aria-describedby="next-day-review-help" placeholder="例如：數學第四單元複習考"></textarea>
        </section>
        <section class="editor-section">
          <div class="section-heading"><h3>聯絡事項</h3><button type="button" class="btn-outline" on:click={() => items = [...items, '']}><Plus size={16} /> 新增事項</button></div>
          {#if !items.length}<p class="description empty">今天尚未新增事項，按「新增事項」開始編寫。</p>{/if}
          {#each items as item, i}
            <div class="item-editor">
              <label for={'contact-item-' + i}>{i + 1}.</label>
              <textarea id={'contact-item-' + i} aria-label={'聯絡事項 ' + (i + 1)} rows="2" bind:value={items[i]} placeholder="輸入作業、攜帶物品或提醒事項"></textarea>
              <div class="item-actions">
                <button type="button" class="btn-outline" aria-label="上移" disabled={i === 0} on:click={() => move(i,-1)}><ArrowUp size={16} /></button>
                <button type="button" class="btn-outline" aria-label="下移" disabled={i === items.length - 1} on:click={() => move(i,1)}><ArrowDown size={16} /></button>
                <button type="button" class="btn-outline" aria-label={'刪除事項 ' + (i + 1)} on:click={() => items = items.filter((_, index) => index !== i)}><Trash2 size={16} /></button>
              </div>
            </div>
          {/each}
        </section>
        <div class="save-bar"><span class="description">{$dirty ? '有尚未儲存的修改' : '內容已儲存'} </span><button class="btn-primary" type="submit" disabled={$saving || !$dirty || rolloverDates.length > 0}><Save size={16} /> 儲存聯絡簿</button></div>
      </form>
    {/key}
  {/if}
</div>
<style>
  .contact-page { overflow-y:auto; }
  .heading,.section-heading,.save-bar { display:flex; align-items:center; justify-content:space-between; gap:16px; flex-wrap:wrap; }
  .heading { margin-bottom:20px; }
  .heading .page-title { margin-bottom:4px; }
  .description { color:var(--text-secondary); font-size: calc(var(--ui-font, 14px) * 14 / 14); }
  button { display:inline-flex; align-items:center; gap:8px; min-height:36px; }
  button:disabled { opacity:.45; cursor:default; }
  form,.reading { max-width:960px; }
  .editor-section { padding:20px 0; border-top:1px solid var(--border); }
  label,h3 { font-size: calc(var(--ui-font, 14px) * 16 / 14); font-weight:600; }
  textarea { display:block; width:100%; resize:vertical; border:1px solid var(--border); border-radius:8px; padding:12px; font:inherit; color:var(--text-primary); background:var(--bg-secondary); line-height:1.7; }
  textarea:focus-visible { outline:2px solid var(--accent); outline-offset:2px; }
  textarea::placeholder { color:#64748b; }
  .field-help { margin-top:8px; }
  .editor-section > textarea { margin-top:12px; }
  .item-editor { display:grid; grid-template-columns:24px minmax(0,1fr) auto; gap:12px; align-items:start; margin-top:16px; }
  .item-editor label { padding-top:12px; }
  .item-actions { display:flex; gap:4px; }
  .item-actions button { padding:8px; }
  .empty { padding:24px 0; }
  .save-bar { border-top:1px solid var(--border); padding:20px 0; }
  .error { color:#b91c1c; padding:12px 0; }
  .reading { background:var(--bg-secondary); border-radius:12px; padding:28px; margin-top:16px; overflow-wrap:anywhere; }
  .reading p, .reading li { white-space:pre-wrap; }
  .reading ol { padding-left:28px; margin-top:16px; }
  .reading li { padding:10px 0; border-bottom:1px solid var(--border); }
  .review { background:#eff6ff; color:#1e3a8a; padding:20px; border-radius:12px; margin-bottom:28px; }
  .review h3 { margin-bottom:8px; }
  .presentation .reading { max-width:none; font-size:clamp(24px,3vw,44px); }
  .presentation h3 { font-size: calc(var(--ui-font, 14px) * 24 / 14); }
  .contact-page:not(.presentation) form, .contact-page:not(.presentation) .reading { max-width: none; }
  @container workspace (max-width: 620px) {
    .item-editor { grid-template-columns: 20px minmax(0, 1fr); }
    .item-actions { grid-column: 2; }
  }
  @container workspace (min-width: 1100px) {
    form { display: grid; grid-template-columns: minmax(240px, 1fr) minmax(0, 2fr); column-gap: 24px; align-items: start; }
    .save-bar { grid-column: 1 / -1; }
  }
</style>
