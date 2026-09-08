<script lang="ts">
  import { draftGuard, markDraftSaved } from '../stores/draft';
  let ready = false;
  import { onMount } from 'svelte';
  import { GetMissingHomework, SaveMissingHomework, GetTimetable, GetStudents } from '../../../wailsjs/go/main/App';
  import { verifyAndRun } from '../stores/auth';
  import { Save, CircleCheck, Trash2 } from 'lucide-svelte';

  interface HomeworkRecord {
    subject: string;
    students: number[];
    note: string;
  }

  let records: HomeworkRecord[] = [];
  let allStudents: { seat_number: number; name: string }[] = [];
  let subjects: string[] = [];
  let showAddDropdown = false;
  let expandedPicker: number | null = null;

  let originalRecords = '';
  let hasChanges = false;

  async function load() {
    allStudents = await GetStudents();
    const tt = await GetTimetable();
    records = await GetMissingHomework();
    originalRecords = JSON.stringify(records);
    hasChanges = false;
    const subjectSet = new Set<string>();
    for (const day of tt) {
      for (const s of day) {
        if (s && s.trim()) subjectSet.add(s.trim());
      }
    }
    subjects = [...subjectSet].sort();
  }

  function markChanged() {
    hasChanges = JSON.stringify(records) !== originalRecords;
  }

  async function handleSave() {
    try {
      await verifyAndRun(async () => {
        await SaveMissingHomework(records);
      }, '儲存作業紀錄');
      originalRecords = JSON.stringify(records);
      hasChanges = false;
      markDraftSaved();
    } catch (e: any) {
      if (e.message !== '驗證失敗或已取消') {
        alert('儲存失敗: ' + e.message);
      }
    }
  }

  function addRecord(subject: string) {
    records = [...records, { subject, students: [], note: '' }];
    showAddDropdown = false;
    markChanged();
  }

  function deleteRecord(index: number) {
    const record = records[index];
    if (!confirm(`確定刪除「${record.subject}」作業紀錄（${record.students.length} 位未交）？\n儲存後將無法復原。`)) return;
    expandedPicker = null;
    records = records.filter((_, i) => i !== index);
    markChanged();
  }

  function toggleStudent(recordIndex: number, seatNumber: number) {
    const rec = records[recordIndex];
    if (rec.students.includes(seatNumber)) {
      rec.students = rec.students.filter(s => s !== seatNumber);
    } else {
      rec.students = [...rec.students, seatNumber].sort((a, b) => a - b);
    }
    records = records;
    markChanged();
  }

  function removeStudent(recordIndex: number, seatNumber: number) {
    records[recordIndex].students = records[recordIndex].students.filter(s => s !== seatNumber);
    records = records;
    markChanged();
  }

  function updateNote(recordIndex: number, value: string) {
    records[recordIndex].note = value;
    records = records;
    markChanged();
  }

  function getStudentName(seatNumber: number): string {
    const s = allStudents.find(st => st.seat_number === seatNumber);
    return s ? s.name : `${seatNumber}號`;
  }

  // svelte-ignore unused-export-let
  function togglePicker(index: number) {
    expandedPicker = expandedPicker === index ? null : index;
  }

  onMount(() => { load().then(() => { ready = true; }).catch((e) => alert("載入失敗，請重新開啟此頁：" + e)); });
</script>

<div class="page" use:draftGuard={{ value: JSON.stringify(records), ready, save: handleSave }}>
  <div class="page-header">
    <h2 class="page-title">作業未交管理</h2>
    <div class="header-actions">
      {#if hasChanges}
        <button class="btn-success animate-fade" on:click={handleSave}>
          <Save size={14} /> 儲存變更
        </button>
      {/if}
      <div class="add-wrapper">
        <button class="btn-primary" on:click={() => (showAddDropdown = !showAddDropdown)}>
          + 新增作業
        </button>
        {#if showAddDropdown}
          <div class="dropdown">
            {#each subjects as subj}
              <button class="dropdown-item" on:click={() => addRecord(subj)}>{subj}</button>
            {/each}
            {#if subjects.length === 0}
              <div class="dropdown-empty">請先設定課表</div>
            {/if}
          </div>
        {/if}
      </div>
    </div>
  </div>

  {#if records.length === 0}
    <div class="empty-state">
      <span class="empty-icon"><CircleCheck size={48} /></span>
      <p>目前沒有未交作業紀錄</p>
    </div>
  {:else}
    <div class="records-list">
      {#each records as record, i}
        <div class="record-card">
          <div class="record-header">
            <h3 class="record-subject">{record.subject}</h3>
            <button class="btn-delete" on:click={() => deleteRecord(i)} title="刪除此紀錄"><Trash2 size={16} /></button>
          </div>

          <div class="note-row">
            <input
              type="text"
              class="note-input"
              placeholder="備註（選填）"
              value={record.note}
              on:input={(e) => updateNote(i, e.currentTarget.value)}
            />
          </div>

          <div class="selected-students">
            {#if record.students.length > 0}
              {#each record.students as seat}
                <span class="student-chip">
                  {seat}號 {getStudentName(seat)}
                  <button class="chip-remove" on:click={() => removeStudent(i, seat)}>×</button>
                </span>
              {/each}
            {:else}
              <span class="no-students">尚未選擇學生</span>
            {/if}
          </div>

          <button class="btn-outline btn-sm" on:click={() => togglePicker(i)}>
            {expandedPicker === i ? '收起' : '選擇學生'}
          </button>

          {#if expandedPicker === i}
            <div class="student-grid">
              {#each allStudents as student}
                <label class="student-checkbox">
                  <input
                    type="checkbox"
                    checked={record.students.includes(student.seat_number)}
                    on:change={() => toggleStudent(i, student.seat_number)}
                  />
                  <span>{student.seat_number}號 {student.name}</span>
                </label>
              {/each}
            </div>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .page-header {
    position: sticky;
    top: 0;
    z-index: 10;
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 20px;
    padding: 8px 0 12px;
    background: var(--bg-primary);
  }
  .page-header .page-title {
    margin-bottom: 0;
  }
  .add-wrapper {
    position: relative;
  }
  .dropdown {
    position: absolute;
    top: 100%;
    left: 0;
    margin-top: 4px;
    background: white;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0,0,0,0.1);
    z-index: 100;
    min-width: 140px;
    max-height: 240px;
    overflow-y: auto;
  }
  .dropdown-item {
    display: block;
    width: 100%;
    padding: 10px 16px;
    text-align: left;
    background: none;
    border: none;
    font-size: 14px;
    cursor: pointer;
    font-family: inherit;
  }
  .dropdown-item:hover {
    background: #f1f5f9;
  }
  .dropdown-empty {
    padding: 10px 16px;
    color: var(--text-secondary);
    font-size: 13px;
  }

  .empty-state {
    text-align: center;
    padding: 60px 20px;
    color: var(--text-secondary);
  }
  .empty-icon {
    font-size: 48px;
    display: block;
    margin-bottom: 12px;
  }

  .records-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }
  .record-card {
    background: white;
    border: 1px solid #e2e8f0;
    border-left: 4px solid #ef4444;
    border-radius: 12px;
    padding: 20px;
  }
  .record-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 12px;
  }
  .record-subject {
    font-size: 16px;
    font-weight: 700;
    color: var(--text-primary);
  }
  .btn-delete {
    background: none;
    border: none;
    font-size: 18px;
    cursor: pointer;
    opacity: 0.6;
    transition: opacity 0.15s;
  }
  .btn-delete:hover {
    opacity: 1;
  }

  .note-row {
    margin-bottom: 12px;
  }
  .note-input {
    width: 100%;
    padding: 8px 12px;
    border: 1px solid #e2e8f0;
    border-radius: 6px;
    font-size: 13px;
    font-family: inherit;
  }
  .note-input:focus {
    outline: none;
    border-color: var(--accent);
  }

  .selected-students {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    margin-bottom: 12px;
    min-height: 32px;
  }
  .student-chip {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    background: #fee2e2;
    color: #991b1b;
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 13px;
    font-weight: 600;
  }
  .chip-remove {
    background: none;
    border: none;
    color: #991b1b;
    font-size: 16px;
    cursor: pointer;
    padding: 0 2px;
    line-height: 1;
  }
  .no-students {
    color: var(--text-secondary);
    font-size: 13px;
  }

  .student-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
    gap: 6px;
    margin-top: 10px;
    padding: 12px;
    background: #f8fafc;
    border-radius: 8px;
  }
  .student-checkbox {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 13px;
    cursor: pointer;
  }
  .student-checkbox input {
    cursor: pointer;
  }

  .header-actions {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-left: auto;
  }

  .btn-success {
    background: var(--success);
    color: white;
  }

  .btn-success:hover {
    background: #16a34a;
  }

  .animate-fade {
    animation: fadeIn 0.2s ease-out;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: scale(0.95); }
    to { opacity: 1; transform: scale(1); }
  }
</style>
