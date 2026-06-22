<script lang="ts">
  import { onMount } from 'svelte';
  import {
    showPrompt,
    promptTitle,
    isSetupMode,
    errorMsg,
    hasPassword,
    submitPassword,
    cancelPrompt
  } from '../stores/auth';

  let passwordInput = '';
  let showPassword = false;
  let rememberMe = false;
  let inputEl: HTMLInputElement;

  function handleSubmit(e: Event) {
    e.preventDefault();
    submitPassword(passwordInput, rememberMe);
    passwordInput = '';
  }

  function handleCancel() {
    cancelPrompt();
    passwordInput = '';
  }

  onMount(() => {
    if (inputEl) {
      inputEl.focus();
    }
  });
</script>

{#if $showPrompt}
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div class="backdrop" on:click|self={$hasPassword ? handleCancel : null}>
    <div class="modal-card">
      <div class="modal-header">
        <h3>{$promptTitle}</h3>
        {#if $hasPassword}
          <button class="btn-icon close-btn" on:click={handleCancel}>✕</button>
        {/if}
      </div>

      <form on:submit={handleSubmit} class="modal-body">
        {#if $isSetupMode}
          <p class="modal-desc">
            首次使用或未設定密碼，請先設定管理員密碼。此密碼將用於保護後續所有寫入與修改操作。
          </p>
        {:else}
          <p class="modal-desc">
            此為敏感寫入操作，請輸入管理員密碼以確認儲存。
          </p>
        {/if}

        <div class="input-wrapper">
          <input
            bind:this={inputEl}
            type={showPassword ? 'text' : 'password'}
            placeholder={$isSetupMode ? '請輸入新密碼' : '請輸入密碼'}
            value={passwordInput}
            on:input={e => passwordInput = e.currentTarget.value}
            class="password-input"
            autocomplete="current-password"
          />
          <button
            type="button"
            class="toggle-visible"
            on:click={() => showPassword = !showPassword}
            tabindex="-1"
          >
            {showPassword ? '🙈' : '👁️'}
          </button>
        </div>

        {#if $errorMsg}
          <div class="error-banner">
            <span class="error-icon">⚠️</span>
            <span class="error-text">{$errorMsg}</span>
          </div>
        {/if}

        {#if !$isSetupMode}
          <label class="remember-label">
            <input type="checkbox" bind:checked={rememberMe} />
            <span>在接下來 5 分鐘內記住此密碼</span>
          </label>
        {/if}

        <div class="modal-footer">
          {#if $hasPassword}
            <button type="button" class="btn-outline" on:click={handleCancel}>取消</button>
          {/if}
          <button type="submit" class="btn-primary">
            {$isSetupMode ? '設定並啟用' : '確認送出'}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}

<style>
  .backdrop {
    position: fixed;
    inset: 0;
    z-index: 9999;
    background: rgba(15, 23, 42, 0.4);
    backdrop-filter: blur(12px);
    display: flex;
    align-items: center;
    justify-content: center;
    animation: fadeIn 0.2s ease-out;
  }

  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  .modal-card {
    background: rgba(255, 255, 255, 0.95);
    border: 1px solid rgba(226, 232, 240, 0.8);
    border-radius: 16px;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
    width: 380px;
    max-width: 90%;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    animation: slideUp 0.25s cubic-bezier(0.16, 1, 0.3, 1);
  }

  @keyframes slideUp {
    from { transform: translateY(16px) scale(0.98); opacity: 0; }
    to { transform: translateY(0) scale(1); opacity: 1; }
  }

  .modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 18px 20px;
    border-bottom: 1px solid var(--border);
  }

  .modal-header h3 {
    font-size: 16px;
    font-weight: 700;
    color: var(--text-primary);
  }

  .close-btn {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    transition: background 0.15s;
  }

  .close-btn:hover {
    background: #f1f5f9;
  }

  .modal-body {
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .modal-desc {
    font-size: 13px;
    color: var(--text-secondary);
    line-height: 1.6;
  }

  .input-wrapper {
    position: relative;
    display: flex;
    align-items: center;
  }

  .password-input {
    width: 100%;
    height: 42px;
    padding: 10px 44px 10px 14px;
    font-size: 14px;
    border: 1.5px solid var(--border);
    border-radius: 10px;
    transition: border-color 0.2s, box-shadow 0.2s;
    background: #f8fafc;
  }

  .password-input:focus {
    border-color: var(--accent);
    background: #fff;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.12);
  }

  .toggle-visible {
    position: absolute;
    right: 8px;
    background: transparent;
    border: none;
    padding: 8px;
    font-size: 16px;
    cursor: pointer;
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background 0.15s;
  }

  .toggle-visible:hover {
    background: #e2e8f0;
  }

  .error-banner {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 12px;
    background: #fef2f2;
    border: 1px solid #fee2e2;
    border-radius: 8px;
    color: var(--danger);
    font-size: 13px;
    font-weight: 500;
  }

  .error-icon {
    font-size: 14px;
  }

  .remember-label {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 13px;
    color: var(--text-secondary);
    cursor: pointer;
    user-select: none;
    padding: 2px 0;
  }

  .remember-label input {
    cursor: pointer;
    width: 16px;
    height: 16px;
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    margin-top: 8px;
  }

  .modal-footer button {
    height: 38px;
    font-size: 13px;
    font-weight: 600;
    padding: 0 18px;
  }
</style>
