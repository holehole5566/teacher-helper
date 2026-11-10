<script lang="ts">
  import { onMount } from 'svelte';
  import { Monitor, Save } from 'lucide-svelte';
  import { GetDailyDates, GetMidtermExam, GetSettings, SaveMidtermExam } from '../../../wailsjs/go/main/App';
  import { models } from '../../../wailsjs/go/models';
  import { verifyAndRun } from '../stores/auth';
  import { draftGuard, markDraftSaved, requestNavigation, dirty, saving } from '../stores/draft';

  export let presentation = false;
  export let initialDate = '';
  export let onDisplay: (date: string) => void = () => {};
  const labels = ['第1節', '第2節', '第3節', '第4節', '午休', '第5節', '第6節', '第7節'];
  let date = '';
  let subjects: string[] = Array(8).fill('');
  let expected: number | undefined;
  let actual: number | undefined;
  let times: string[] = [];
  let ready = false;
  let error = '';
  $: draft = JSON.stringify({ date, subjects, expected, actual });
  $: absent = expected != null && actual != null ? expected - actual : null;

  async function load(next: string) {
    ready = false; error = ''; date = next;
    try {
      const [record, settings] = await Promise.all([GetMidtermExam(next), GetSettings()]);
      subjects = labels.map((_, i) => record.subjects?.[i] || '');
      expected = record.expected ?? undefined;
      actual = record.actual ?? undefined;
      times = settings.period_times || [];
      ready = true;
    } catch (e) { error = '載入失敗，請重新載入：' + String(e); }
  }
  function changeDate(next: string) {
    if (!next || next === date) return;
    requestNavigation(() => { load(next); });
  }
  async function save() {
    error = '';
    try {
      for (const count of [expected, actual]) {
        if (count != null && (!Number.isInteger(count) || count < 0)) throw new Error('請填入 0 或正整數人數');
      }
      if (absent != null && absent < 0) throw new Error('實到不可大於應到');
      const record = new models.MidtermExam({ date, subjects, expected: expected ?? null, actual: actual ?? null });
      await verifyAndRun(() => SaveMidtermExam(record), '儲存期中考');
      markDraftSaved();
    } catch (e) { error = String(e); throw e; }
  }
  onMount(() => {
    (async () => {
      try { await load(initialDate || (await GetDailyDates())[0]); }
      catch (e) { error = '載入失敗：' + String(e); }
    })();
  });
</script>

<div class="page adaptive-page exam-page" class:presentation>
  <div class="heading">
    <div><h2 class="page-title">期中考</h2><p class="description">{presentation ? date : '按日期安排考試或自習，時間沿用設定中的每節上課時間。'}</p></div>
    {#if !presentation}<button class="btn-primary" disabled={!ready || $saving} on:click={() => requestNavigation(() => onDisplay(date))}><Monitor size={16} /> 獨立展示</button>{/if}
  </div>
  {#if !presentation}
    <div class="date-picker"><label for="exam-date">考試日期</label><input id="exam-date" type="date" value={date} disabled={$saving} on:change={(event) => { const next = event.currentTarget.value; event.currentTarget.value = date; changeDate(next); }} /></div>
  {/if}
  {#if error}<p class="error" role="alert">{error}</p>{/if}
  {#if !ready}
    <div class="empty-state">{error ? '無法顯示此日期的期中考資料。' : '載入中…'}{#if error}<button class="btn-outline" on:click={() => load(date)}>重新載入</button>{/if}</div>
  {:else if presentation}
    <div class="attendance reading" aria-label="全日出席人數">
      <p>應到 <strong>{expected ?? '—'}</strong> 人</p><p>實到 <strong>{actual ?? '—'}</strong> 人</p><p>未到 <strong>{absent ?? '—'}</strong> 人</p>
    </div>
    <table><caption class="sr-only">{date} 期中考時間表</caption><thead><tr><th scope="col">節次</th><th scope="col">開始時間</th><th scope="col">考試科目／自習</th></tr></thead><tbody>{#each labels as label, i}<tr class:rest={i === 4}><th scope="row">{label}</th><td>{times[i] || '未設定'}</td><td>{subjects[i] || (i === 4 ? '午休' : '未安排')}</td></tr>{/each}</tbody></table>
  {:else}
    {#key date}
      <form use:draftGuard={{ value: draft, ready, save }} on:submit|preventDefault={() => { save().catch(() => {}); }}>
        <fieldset class="attendance"><legend>全日出席人數</legend><label for="exam-expected">應到（人）<input id="exam-expected" type="number" min="0" step="1" bind:value={expected} placeholder="尚未填寫" /></label><label for="exam-actual">實到（人）<input id="exam-actual" type="number" min="0" step="1" bind:value={actual} placeholder="尚未填寫" /></label><p>未到：{absent ?? '—'} 人</p></fieldset>
        <table><caption class="sr-only">期中考時間表設定</caption><thead><tr><th scope="col">節次</th><th scope="col">開始時間</th><th scope="col">考試科目／自習</th></tr></thead><tbody>
          {#each labels as label, i}<tr class:rest={i === 4}><th scope="row"><label for={'exam-subject-' + i}>{label}</label></th><td>{times[i] || '未設定'}</td><td><input id={'exam-subject-' + i} type="text" list="exam-subjects" bind:value={subjects[i]} placeholder={i === 4 ? '午休' : '輸入科目或選擇自習'} /></td></tr>{/each}
        </tbody></table>
        <datalist id="exam-subjects"><option value="自習" /><option value="國語" /><option value="數學" /><option value="英語" /><option value="自然" /><option value="社會" /></datalist>
        <div class="save-bar"><span class="description">{$dirty ? '有尚未儲存的修改' : '內容已儲存'}</span><button class="btn-primary" type="submit" disabled={$saving || !$dirty}><Save size={16} /> 儲存期中考</button></div>
      </form>
    {/key}
  {/if}
</div>

<style>
  .exam-page { overflow-y: auto; }
  .heading, .save-bar, .date-picker, .attendance { display: flex; align-items: center; gap: 16px; flex-wrap: wrap; }
  .heading, .save-bar { justify-content: space-between; }
  .heading { margin-bottom: 20px; }
  .heading .page-title { margin-bottom: 4px; }
  .description { color: var(--text-secondary); }
  .date-picker { margin-bottom: var(--page-gap); }
  button { display: inline-flex; align-items: center; gap: 8px; min-height: 36px; }
  button:disabled { opacity: .45; cursor: default; }
  input { font: inherit; color: var(--text-primary); background: var(--bg-secondary); border: 1px solid var(--border); border-radius: 8px; padding: 10px 12px; min-width: 0; }
  input::placeholder { color: var(--text-secondary); }
  input:focus-visible { outline: 2px solid var(--accent); outline-offset: 2px; }
  .attendance { border: 0; border-top: 1px solid var(--border); padding: var(--page-gap) 0; margin: 0; }
  legend { font-weight: 600; padding-right: 12px; }
  .attendance label { display: flex; align-items: center; gap: 12px; }
  .attendance input { width: 140px; }
  table { width: 100%; table-layout:fixed; border-collapse: collapse; font-variant-numeric: tabular-nums; }
  th, td { overflow-wrap:anywhere; text-align: left; padding: 12px 16px; border-bottom: 1px solid var(--border); }
  thead { background: var(--bg-secondary); color: var(--text-secondary); }
  td input { width: 100%; }
  .rest { background: var(--bg-secondary); }
  .save-bar { padding: 20px 0; }
  .error { color: var(--danger); margin-bottom: 16px; }
  .sr-only { position: absolute; width: 1px; height: 1px; padding: 0; overflow: hidden; clip-path: inset(50%); white-space: nowrap; }
  .presentation { font-size: clamp(16px, min(2.4vw, 3vh), 36px); }
  .presentation .attendance { justify-content: space-around; padding: 12px 0 20px; }
  .presentation strong { font-size: 1.5em; }
  .presentation th, .presentation td { padding: clamp(4px, .75vh, 12px) clamp(6px, 1vw, 18px); }
  .presentation td:last-child { font-weight: 600; overflow-wrap: anywhere; }
  @container workspace (max-width: 620px) {
    th, td { padding: 10px 6px; }
    .attendance label { width: 100%; justify-content: space-between; }
    .attendance input { flex: 1; max-width: 200px; }
  }
  th:first-child { width:20%; }
  th:nth-child(2) { width:25%; }
  .presentation .page-title { font-size:clamp(24px,min(3vw,4vh),44px); }
  .presentation .attendance { gap:8px; padding:8px 0 clamp(8px,1.5vh,20px); }
  @media(max-height:700px) {
    th, td { padding:6px 10px; }
    input { padding:6px 8px; }
    .attendance { gap:10px; }
  }
  @media(max-width:520px) {
    th, td { padding:8px 4px; }
    .attendance label { width:100%; justify-content:space-between; }
    .attendance input { flex:1; max-width:160px; }
  }
</style>
