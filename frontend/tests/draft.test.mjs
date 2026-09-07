import { test, after } from 'node:test';
import assert from 'node:assert/strict';
import { mkdtempSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { pathToFileURL } from 'node:url';
import { buildSync } from 'esbuild';

const directory = mkdtempSync(join(tmpdir(), 'teacher-draft-test-'));
const outfile = join(directory, 'draft.mjs');
buildSync({ entryPoints: ['src/lib/stores/draft.ts'], bundle: true, platform: 'node', format: 'esm', outfile });
const draft = await import(pathToFileURL(outfile).href);
after(() => rmSync(directory, { recursive: true, force: true }));
const valueOf = store => { let value; const unsubscribe = store.subscribe(v => value = v); unsubscribe(); return value; };

function setup(save = async () => {}) {
  draft.cancelNavigation();
  draft.saving.set(false);
  const node = { inert: false };
  const guard = draft.draftGuard(node, { value: 'initial', ready: false, save });
  const update = (value, ready = true) => guard.update({ value, ready, save });
  update('loaded');
  return { node, guard, update };
}

test('loading sets baseline; edit and revert track actual changes', () => {
  const { guard, update, node } = setup();
  assert.equal(node.inert, false);
  assert.equal(valueOf(draft.dirty), false);
  update('edited');
  assert.equal(valueOf(draft.dirty), true);
  update('loaded');
  assert.equal(valueOf(draft.dirty), false);
  guard.destroy();
});

test('clean navigation runs immediately', () => {
  const { guard } = setup();
  let navigated = false;
  draft.requestNavigation(() => navigated = true);
  assert.equal(navigated, true);
  assert.equal(valueOf(draft.leavePrompt), false);
  guard.destroy();
});

test('cancel preserves draft and discard navigates', async () => {
  const { guard, update } = setup();
  update('edited');
  let navigations = 0;
  draft.requestNavigation(() => navigations++);
  assert.equal(valueOf(draft.leavePrompt), true);
  draft.cancelNavigation();
  assert.equal(navigations, 0);
  assert.equal(valueOf(draft.dirty), true);
  draft.requestNavigation(() => navigations++);
  await draft.resolveNavigation(false);
  assert.equal(navigations, 1);
  guard.destroy();
});

test('successful save updates baseline before leaving', async () => {
  const { guard, update } = setup(async () => draft.markDraftSaved());
  update('edited');
  let navigated = false;
  draft.requestNavigation(() => navigated = true);
  await draft.resolveNavigation(true);
  assert.equal(navigated, true);
  assert.equal(valueOf(draft.dirty), false);
  update('edited');
  assert.equal(valueOf(draft.dirty), false);
  update('new edit');
  assert.equal(valueOf(draft.dirty), true);
  guard.destroy();
});

for (const [name, save] of [
  ['cancelled authentication or validation', async () => {}],
  ['failed save', async () => { throw new Error('disk failure'); }],
]) {
  test(`${name} keeps draft and does not navigate`, async () => {
    const { guard, update } = setup(save);
    update('edited');
    let navigated = false;
    draft.requestNavigation(() => navigated = true);
    await draft.resolveNavigation(true);
    assert.equal(navigated, false);
    assert.equal(valueOf(draft.dirty), true);
    assert.equal(valueOf(draft.leavePrompt), true);
    guard.destroy();
  });
}

test('busy state prevents navigation and duplicate save', async () => {
  let saves = 0;
  const { guard, update } = setup(async () => saves++);
  update('edited');
  let navigated = false;
  draft.saving.set(true);
  draft.requestNavigation(() => navigated = true);
  await draft.resolveNavigation(true);
  assert.equal(navigated, false);
  assert.equal(saves, 0);
  draft.saving.set(false);
  guard.destroy();
  assert.equal(valueOf(draft.dirty), false);
});
