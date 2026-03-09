<script>
  import { onMount } from 'svelte';
  import { GetLanguages, TranslateText, GetSettings } from '../../wailsjs/go/main/App.js';

  let languages = [];
  let fromLang = 'English';
  let toLang = 'Spanish';
  let sourceText = '';
  let translatedText = '';
  let loading = false;
  let debounceTimer = null;

  onMount(async () => {
    languages = await GetLanguages();
    try {
      const s = await GetSettings();
      fromLang = s.currentFromLang || 'English';
      toLang   = s.currentToLang   || 'Spanish';
    } catch {}
  });

  function onSourceChange() {
    clearTimeout(debounceTimer);
    if (!sourceText.trim()) { translatedText = ''; return; }
    debounceTimer = setTimeout(translate, 1000);
  }

  async function translate() {
    if (!sourceText.trim() || loading) return;
    loading = true;
    try {
      translatedText = await TranslateText(sourceText, fromLang, toLang);
    } catch (e) {
      translatedText = 'Error: ' + e.toString();
    } finally {
      loading = false;
    }
  }

  function swap() {
    [fromLang, toLang] = [toLang, fromLang];
    [sourceText, translatedText] = [translatedText, sourceText];
  }
</script>

<div class="view">
  <h1>Translate Text</h1>

  <div class="lang-row">
    <select bind:value={fromLang}>
      {#each languages as lang}
        <option value={lang}>{lang}</option>
      {/each}
    </select>

    <button class="swap-btn" on:click={swap} title="Swap languages">⇄</button>

    <select bind:value={toLang}>
      {#each languages as lang}
        <option value={lang}>{lang}</option>
      {/each}
    </select>

    <button class="btn-primary" on:click={translate} disabled={loading || !sourceText.trim()}>
      {loading ? 'Translating…' : 'Translate'}
    </button>
  </div>

  <div class="panels">
    <div class="panel">
      <div class="panel-header">
        <span>{fromLang}</span>
      </div>
      <textarea
        bind:value={sourceText}
        on:input={onSourceChange}
        placeholder="Enter text to translate…"
        disabled={loading}
      ></textarea>
    </div>

    <div class="panel">
      <div class="panel-header">
        <span>{toLang}</span>
        {#if translatedText}
          <button class="btn-sm" on:click={() => navigator.clipboard.writeText(translatedText)}>Copy</button>
        {/if}
      </div>
      <textarea
        value={translatedText}
        placeholder="Translation will appear here…"
        readonly
      ></textarea>
    </div>
  </div>
</div>

<style>
  .view { padding: 24px 32px; display: flex; flex-direction: column; height: 100%; gap: 16px; }
  h1 { font-size: 20px; color: #7cb9e8; }

  .lang-row {
    display: flex;
    align-items: center;
    gap: 10px;
  }

  select {
    background: #252526;
    border: 1px solid #444;
    border-radius: 5px;
    color: #e0e0e0;
    padding: 6px 10px;
    font-size: 13px;
  }

  .swap-btn {
    background: #333;
    border: 1px solid #555;
    border-radius: 5px;
    color: #ccc;
    padding: 5px 10px;
    cursor: pointer;
    font-size: 16px;
  }
  .swap-btn:hover { background: #444; }

  .btn-primary {
    padding: 6px 16px;
    background: #0078d7;
    color: #fff;
    border: none;
    border-radius: 5px;
    font-size: 13px;
    cursor: pointer;
  }
  .btn-primary:hover:not(:disabled) { background: #106ebe; }
  .btn-primary:disabled { opacity: 0.5; cursor: not-allowed; }

  .btn-sm {
    padding: 2px 8px;
    background: #333;
    color: #ccc;
    border: 1px solid #555;
    border-radius: 4px;
    font-size: 11px;
    cursor: pointer;
  }
  .btn-sm:hover { background: #444; }

  .panels {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
    flex: 1;
    min-height: 0;
  }

  .panel {
    display: flex;
    flex-direction: column;
    border: 1px solid #333;
    border-radius: 6px;
    overflow: hidden;
    background: #252526;
  }

  .panel-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 6px 12px;
    background: #2d2d2e;
    border-bottom: 1px solid #333;
    font-size: 12px;
    color: #aaa;
  }

  textarea {
    flex: 1;
    background: transparent;
    border: none;
    color: #e0e0e0;
    padding: 12px;
    font-size: 14px;
    font-family: inherit;
    resize: none;
    line-height: 1.5;
  }
  textarea:focus { outline: none; }
  textarea[readonly] { color: #bbb; }
</style>
