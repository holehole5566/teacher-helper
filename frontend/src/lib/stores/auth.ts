import { writable, get } from 'svelte/store';
import { saving, notice } from './draft';

let noticeTimer: ReturnType<typeof setTimeout>;

import { HasPassword, SetPassword, VerifyPassword } from '../../../wailsjs/go/main/App';

export const hasPassword = writable<boolean>(false);
export const showPrompt = writable<boolean>(false);
export const promptTitle = writable<string>('');
export const isSetupMode = writable<boolean>(false);
export const errorMsg = writable<string>('');

let resolvePromise: ((value: boolean) => void) | null = null;
let currentPasswordVal = '';
let rememberMeVal = false;
let sessionExpiryTime = 0; // timestamp

export async function checkPasswordStatus() {
  const hasPw = await HasPassword();
  hasPassword.set(hasPw);
  if (!hasPw) {
    isSetupMode.set(true);
    promptTitle.set('設定管理員密碼');
    showPrompt.set(true);
  }
}

export async function verifyAndRun<T>(action: () => Promise<T>, title: string): Promise<T> {
  if (get(saving)) throw new Error('操作進行中，請稍候');
  saving.set(true);
  notice.set('');
  try {
    const result = await runVerified(action, title);
    notice.set(`${title}完成`);
    clearTimeout(noticeTimer);
    noticeTimer = setTimeout(() => notice.set(''), 4000);
    return result;
  } finally {
    saving.set(false);
  }
}

function runVerified<T>(action: () => Promise<T>, title: string): Promise<T> {
  return new Promise<T>(async (resolve, reject) => {
    // If no password is set on the backend, check status first
    let hasPw: boolean;
    try {
      hasPw = await HasPassword();
    } catch (error) {
      reject(error);
      return;
    }
    hasPassword.set(hasPw);
    if (!hasPw) {
      isSetupMode.set(true);
      promptTitle.set('設定管理員密碼');
      showPrompt.set(true);
      // Wait for password setup
      const setupSuccess = await promptForPassword();
      if (!setupSuccess) {
        reject(new Error('未設定密碼，無法進行此操作'));
        return;
      }
    }

    // Check if session is unlocked
    if (Date.now() < sessionExpiryTime && currentPasswordVal) {
      try {
        // Silently verify password to authorize backend write
        const ok = await VerifyPassword(currentPasswordVal);
        if (ok) {
          try {
            const result = await action();
            resolve(result);
          } catch (actionErr) {
            reject(actionErr);
          }
          return;
        }
      } catch (err) {
        // Ignore and fall through to manual prompt
      }
    }

    // Otherwise prompt for password
    isSetupMode.set(false);
    promptTitle.set(title);
    errorMsg.set('');
    showPrompt.set(true);

    const success = await promptForPassword();
    if (success) {
      try {
        const result = await action();
        resolve(result);
      } catch (err) {
        reject(err);
      }
    } else {
      reject(new Error('驗證失敗或已取消'));
    }
  });
}

function promptForPassword(): Promise<boolean> {
  return new Promise<boolean>((resolve) => {
    resolvePromise = resolve;
  });
}

export async function submitPassword(password: string, remember: boolean) {
  errorMsg.set('');
  if (!password) {
    errorMsg.set('密碼不能為空');
    return;
  }

  const isSetup = get(isSetupMode);
  if (isSetup) {
    try {
      // Temporarily bypass checkWriteAuth in backend SetPassword since we verify isSetup
      // To satisfy Go SetPassword, verify with empty (initial setup allows setting if no password is set)
      await SetPassword(password);
      hasPassword.set(true);
      showPrompt.set(false);
      isSetupMode.set(false);
      if (resolvePromise) {
        resolvePromise(true);
        resolvePromise = null;
      }
    } catch (err: any) {
      errorMsg.set(err?.message || String(err));
    }
  } else {
    try {
      const ok = await VerifyPassword(password);
      if (ok) {
        currentPasswordVal = password;
        rememberMeVal = remember;
        if (remember) {
          sessionExpiryTime = Date.now() + 5 * 60 * 1000; // 5 minutes
        } else {
          sessionExpiryTime = 0;
        }
        showPrompt.set(false);
        if (resolvePromise) {
          resolvePromise(true);
          resolvePromise = null;
        }
      } else {
        errorMsg.set('密碼錯誤');
      }
    } catch (err: any) {
      errorMsg.set(err?.message || String(err));
    }
  }
}

export function cancelPrompt() {
  showPrompt.set(false);
  if (resolvePromise) {
    resolvePromise(false);
    resolvePromise = null;
  }
}
