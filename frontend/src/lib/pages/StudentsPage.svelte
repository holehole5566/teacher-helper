<script lang="ts">
  import { onMount } from 'svelte';
  import { GetStudents, AddStudent, DeleteStudent, ToggleDuty, ToggleLunch, ReportError } from '../../../wailsjs/go/main/App';
  import { verifyAndRun } from '../stores/auth';

  let students: any[] = [];
  let newSeat = '';
  let newName = '';
  let errorMsg = '';

  async function loadStudents() {
    students = await GetStudents();
  }

  async function handleAdd() {
    errorMsg = '';
    const seat = parseInt(newSeat);
    const name = newName.trim();
    if (!seat || !name) {
      errorMsg = '請輸入座號和姓名';
      return;
    }
    try {
      await verifyAndRun(async () => {
        await AddStudent(seat, name);
      }, '新增學生');
      newSeat = '';
      newName = '';
      await loadStudents();
    } catch (e: any) {
      if (e.message !== '驗證失敗或已取消') {
        errorMsg = e?.message || String(e);
        ReportError(`新增學生失敗：${e?.message || e}`);
      }
    }
  }

  async function handleDelete(seatNumber: number) {
    try {
      await verifyAndRun(async () => {
        await DeleteStudent(seatNumber);
      }, '刪除學生');
      await loadStudents();
    } catch (e: any) {
      if (e.message !== '驗證失敗或已取消') {
        alert('刪除學生失敗: ' + e.message);
      }
    }
  }

  async function handleToggleDuty(seatNumber: number) {
    try {
      await verifyAndRun(async () => {
        await ToggleDuty(seatNumber);
      }, '修改值日狀態');
      await loadStudents();
    } catch (e: any) {
      if (e.message !== '驗證失敗或已取消') {
        alert('修改值日狀態失敗: ' + e.message);
      }
    }
  }

  async function handleToggleLunch(seatNumber: number) {
    try {
      await verifyAndRun(async () => {
        await ToggleLunch(seatNumber);
      }, '修改抬餐狀態');
      await loadStudents();
    } catch (e: any) {
      if (e.message !== '驗證失敗或已取消') {
        alert('修改抬餐狀態失敗: ' + e.message);
      }
    }
  }

  onMount(loadStudents);
</script>

<div class="page">
  <h2 class="page-title">學生管理</h2>

  <div class="card add-section">
    <div class="inline-form">
      <input type="number" placeholder="座號" bind:value={newSeat} min="1" style="width: 80px; flex: none;" />
      <input type="text" placeholder="姓名" bind:value={newName} />
      <button class="btn-primary" on:click={handleAdd}>新增</button>
    </div>
    {#if errorMsg}
      <p class="error-msg">{errorMsg}</p>
    {/if}
  </div>

  <div class="card" style="margin-top: 16px;">
    {#if students.length === 0}
      <div class="empty-state">尚無學生，請先新增</div>
    {:else}
      <table>
        <thead>
          <tr>
            <th style="width: 60px;">座號</th>
            <th>姓名</th>
            <th style="width: 80px;">值日</th>
            <th style="width: 80px;">抬餐</th>
            <th style="width: 160px;">操作</th>
          </tr>
        </thead>
        <tbody>
          {#each students as student}
            <tr>
              <td>{student.seat_number}</td>
              <td>{student.name}</td>
              <td>
                <button
                  class="badge"
                  class:badge-success={student.duty_enabled}
                  class:badge-muted={!student.duty_enabled}
                  on:click={() => handleToggleDuty(student.seat_number)}
                >
                  {student.duty_enabled ? '啟用' : '停用'}
                </button>
              </td>
              <td>
                <button
                  class="badge"
                  class:badge-success={student.lunch_enabled}
                  class:badge-muted={!student.lunch_enabled}
                  on:click={() => handleToggleLunch(student.seat_number)}
                >
                  {student.lunch_enabled ? '啟用' : '停用'}
                </button>
              </td>
              <td>
                <div class="actions">
                  <button class="btn-danger btn-sm" on:click={() => handleDelete(student.seat_number)}>刪除</button>
                </div>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    {/if}
  </div>
</div>

<style>
  .add-section {
    margin-bottom: 0;
  }
  .error-msg {
    color: var(--danger);
    font-size: 12px;
    margin-top: 8px;
  }
  .badge {
    cursor: pointer;
    border: none;
  }
</style>
