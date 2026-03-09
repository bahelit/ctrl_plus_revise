<script>
  import { onMount } from 'svelte';
  import { GetShortcuts, SaveShortcuts, GetSettings } from '../../wailsjs/go/main/App.js';

  let shortcuts = { askAI: [], ctrlRevise: [], translate: [] };
  let speakEnabled = false;
  let saved = false;

  onMount(async () => {
    [shortcuts] = await Promise.all([GetShortcuts()]);
    try {
      const s = await GetSettings();
      speakEnabled = s.speakAIResponse;
    } catch {}
  });

  function keysDisplay(keys) {
    return keys.map(k => k.toUpperCase()).join(' + ');
  }

  function parseKeys(str) {
    return str.split('+').map(k => k.trim().toLowerCase()).filter(Boolean);
  }

  let askStr = '';
  let reviseStr = '';
  let translateStr = '';

  $: askStr = keysDisplay(shortcuts.askAI || []);
  $: reviseStr = keysDisplay(shortcuts.ctrlRevise || []);
  $: translateStr = keysDisplay(shortcuts.translate || []);

  async function save() {
    try {
      const s = {
        askAI: parseKeys(askStr),
        ctrlRevise: parseKeys(reviseStr),
        translate: parseKeys(translateStr),
      };
      await SaveShortcuts(s);
      saved = true;
      setTimeout(() => saved = false, 2000);
    } catch (e) {
      alert('Failed to save: ' + e);
    }
  }
</script>

<div class="view">
  <h1>Keyboard Shortcuts</h1>
  <p class="hint">Enter keys separated by <code>+</code> (e.g. <code>alt + c</code>). Changes require an app restart.</p>

  <div class="warn">⚠️ Restart the application for shortcut changes to take effect.</div>

  <section>
    <h2>Configurable Shortcuts</h2>

    <div class="shortcut-row">
      <label>Ask AI with highlighted text</label>
      <input type="text" bind:value={askStr} placeholder="alt + a" />
    </div>

    <div class="shortcut-row">
      <label>Ctrl+Revise highlighted text</label>
      <input type="text" bind:value={reviseStr} placeholder="alt + c" />
    </div>

    <div class="shortcut-row">
      <label>Translate highlighted text</label>
      <input type="text" bind:value={translateStr} placeholder="alt + t" />
    </div>
  </section>

  <section>
    <h2>Fixed Shortcuts</h2>
    <div class="fixed-row">
      <span>Cycle through AI actions</span>
      <kbd>Alt + P</kbd>
    </div>
    {#if speakEnabled}
      <div class="fixed-row">
        <span>Read highlighted text aloud</span>
        <kbd>Alt + R</kbd>
      </div>
    {/if}
  </section>

  <div class="save-row">
    <button class="btn-primary" on:click={save}>Save Shortcuts</button>
    {#if saved}<span class="saved-msg">✅ Saved — restart to apply</span>{/if}
  </div>
</div>

<style>
  .view { padding: 24px 32px; overflow-y: auto; height: 100%; display: flex; flex-direction: column; gap: 16px; }
  h1 { font-size: 20px; color: #7cb9e8; }
  .hint { font-size: 12px; color: #888; }
  code { background: #333; padding: 1px 4px; border-radius: 3px; font-size: 11px; }

  .warn {
    background: #3a2e1a;
    border: 1px solid #7a5c1a;
    border-radius: 5px;
    padding: 8px 14px;
    font-size: 13px;
    color: #d4a040;
  }

  section {
    background: #252526;
    border: 1px solid #333;
    border-radius: 6px;
    padding: 14px 16px;
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  h2 { font-size: 12px; color: #888; text-transform: uppercase; letter-spacing: 0.05em; margin-bottom: 4px; }

  .shortcut-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
  }

  .shortcut-row label {
    flex: 1;
    font-size: 13px;
    color: #ccc;
  }

  .shortcut-row input {
    width: 160px;
    background: #1e1e1e;
    border: 1px solid #444;
    border-radius: 4px;
    color: #e0e0e0;
    padding: 5px 10px;
    font-size: 12px;
    font-family: monospace;
    text-transform: uppercase;
  }
  .shortcut-row input:focus { outline: none; border-color: #7cb9e8; }

  .fixed-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 13px;
    color: #ccc;
  }

  kbd {
    background: #333;
    border: 1px solid #555;
    border-radius: 3px;
    padding: 3px 8px;
    font-size: 11px;
    font-family: monospace;
    color: #ddd;
  }

  .save-row { display: flex; align-items: center; gap: 12px; }

  .btn-primary {
    padding: 8px 24px;
    background: #0078d7;
    color: #fff;
    border: none;
    border-radius: 5px;
    font-size: 14px;
    cursor: pointer;
  }
  .btn-primary:hover { background: #106ebe; }

  .saved-msg { font-size: 13px; color: #6dbf6d; }
</style>
