import { get, writable } from 'svelte/store';

export const dirty = writable(false);
export const saving = writable(false);
export const notice = writable('');
export const leavePrompt = writable(false);
let saveDraft: (() => Promise<void>) | null = null;
let acceptBaseline: (() => void) | null = null;
let pendingNavigation: (() => void) | null = null;

// Compare serialized form data, including changes made by add/remove buttons.
export function draftGuard(node: HTMLElement, initial: { value: string; ready: boolean; save: () => Promise<void> }) {
  let baseline: string | null = null;
  let latest = initial.value;
  const accept = () => { baseline = latest; dirty.set(false); };
  function update(options: typeof initial) {
    node.inert = !options.ready;
    latest = options.value;
    saveDraft = options.save;
    acceptBaseline = accept;
    if (!options.ready) return;
    if (baseline === null) baseline = latest;
    dirty.set(latest !== baseline);
  }
  update(initial);
  return {
    update,
    destroy() {
      if (acceptBaseline === accept) {
        saveDraft = null;
        acceptBaseline = null;
        dirty.set(false);
      }
    }
  };
}

export function markDraftSaved() {
  acceptBaseline?.();
}

export function requestNavigation(navigate: () => void) {
  if (get(saving)) return;
  if (!get(dirty)) { navigate(); return; }
  pendingNavigation = navigate;
  leavePrompt.set(true);
}

export function cancelNavigation() {
  pendingNavigation = null;
  leavePrompt.set(false);
}

export async function resolveNavigation(save: boolean) {
  if (get(saving)) return;
  if (save) {
    try {
      await saveDraft?.();
    } catch {
      return;
    }
    // A cancelled password prompt, validation error or failed save must not leave.
    if (get(dirty)) return;
  }
  const navigate = pendingNavigation;
  cancelNavigation();
  navigate?.();
}
