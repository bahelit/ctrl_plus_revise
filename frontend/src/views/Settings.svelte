<script>
  import { onMount } from 'svelte';
  import { GetSettings, SaveSettings, GetModels, GetPrompts, CheckOllama, InstallOllama, PullCurrentModel, TestOllamaURL } from '../../wailsjs/go/main/App.js';

  let settings = null;
  let models = [];
  let prompts = [];
  let saved = false;
  let ollamaStatus = '';
  let pullStatus = '';
  let installStatus = '';
  let testUrlStatus = '';

  onMount(async () => {
    [settings, models, prompts] = await Promise.all([GetSettings(), GetModels(), GetPrompts()]);
  });

  async function save() {
    try {
      await SaveSettings(settings);
      saved = true;
      setTimeout(() => saved = false, 2000);
    } catch (e) {
      alert('Failed to save: ' + e);
    }
  }

  async function checkOllama() {
    ollamaStatus = 'checking…';
    const ok = await CheckOllama();
    ollamaStatus = ok ? '✅ Connected' : '❌ Not reachable';
  }

  async function pullModel() {
    pullStatus = 'Pulling…';
    pullStatus = await PullCurrentModel();
  }

  async function installOllama() {
    installStatus = 'Installing…';
    installStatus = await InstallOllama();
  }

  async function testURL() {
    if (!settings.ollamaURL) return;
    testUrlStatus = 'Testing…';
    const ok = await TestOllamaURL(settings.ollamaURL);
    testUrlStatus = ok ? '✅ Connected' : '❌ Failed';
  }
</script>

<div class="view">
  {#if !settings}
    <p style="color:#888;padding:24px">Loading settings…</p>
  {:else}
    <h1>Settings</h1>

    <section>
      <h2>AI Model</h2>
      <div class="field">
        <label>Active Model</label>
        <select bind:value={settings.currentModel}>
          {#each models as m}
            <option value={m.id}>{m.name} ({m.memMB} MB)</option>
          {/each}
        </select>
      </div>
      <div class="field">
        <label>Chat Model</label>
        <select bind:value={settings.currentChatModel}>
          {#each models as m}
            <option value={m.id}>{m.name} ({m.memMB} MB)</option>
          {/each}
        </select>
      </div>
      <div class="field">
        <label>Default AI Action</label>
        <select bind:value={settings.currentPrompt}>
          {#each prompts as p}
            <option value={p.id}>{p.name}</option>
          {/each}
        </select>
      </div>
      <div class="btn-row">
        <button class="btn" on:click={pullModel}>Download / Update Model</button>
        {#if pullStatus}<span class="status-text">{pullStatus}</span>{/if}
      </div>
    </section>

    <section>
      <h2>Behaviour</h2>
      <label class="checkbox"><input type="checkbox" bind:checked={settings.replaceHighlightedText} /> Paste AI response over highlighted text</label>
      <label class="checkbox"><input type="checkbox" bind:checked={settings.showPopUp} /> Show response notification</label>
      <label class="checkbox"><input type="checkbox" bind:checked={settings.speakAIResponse} /> Speak AI response aloud (Alt+R to read selection)</label>
      <label class="checkbox"><input type="checkbox" bind:checked={settings.stopOllamaOnShutDown} /> Stop Ollama when app closes</label>
      <label class="checkbox"><input type="checkbox" bind:checked={settings.showStartWindow} /> Show window on startup</label>
    </section>

    <section>
      <h2>Translation</h2>
      <div class="inline-fields">
        <div class="field">
          <label>From</label>
          <input type="text" bind:value={settings.currentFromLang} />
        </div>
        <div class="field">
          <label>To</label>
          <input type="text" bind:value={settings.currentToLang} />
        </div>
      </div>
    </section>

    <section>
      <h2>Ollama Connection</h2>
      <label class="checkbox"><input type="checkbox" bind:checked={settings.useRemoteOllama} /> Connect to remote Ollama server</label>
      {#if settings.useRemoteOllama}
        <div class="field" style="margin-top:8px">
          <label>Ollama URL</label>
          <div class="input-row">
            <input type="text" bind:value={settings.ollamaURL} placeholder="http://localhost:11434" />
            <button class="btn-sm" on:click={testURL}>Test</button>
          </div>
          {#if testUrlStatus}<span class="status-text">{testUrlStatus}</span>{/if}
        </div>
      {/if}
      <label class="checkbox" style="margin-top:8px"><input type="checkbox" bind:checked={settings.useDocker} /> Run Ollama in Docker</label>
      <div class="btn-row" style="margin-top:8px">
        <button class="btn" on:click={checkOllama}>Check Connection</button>
        {#if ollamaStatus}<span class="status-text">{ollamaStatus}</span>{/if}
        <button class="btn" on:click={installOllama}>Install / Update Ollama</button>
        {#if installStatus}<span class="status-text">{installStatus}</span>{/if}
      </div>
    </section>

    <div class="save-row">
      <button class="btn-primary" on:click={save}>Save Settings</button>
      {#if saved}<span class="saved-msg">✅ Saved</span>{/if}
    </div>
  {/if}
</div>

<style>
  .view { padding: 20px 32px; overflow-y: auto; height: 100%; display: flex; flex-direction: column; gap: 8px; }
  h1 { font-size: 20px; color: #7cb9e8; margin-bottom: 4px; }

  section {
    background: #252526;
    border: 1px solid #333;
    border-radius: 6px;
    padding: 14px 16px;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  h2 { font-size: 12px; color: #888; text-transform: uppercase; letter-spacing: 0.05em; margin-bottom: 4px; }

  .field { display: flex; flex-direction: column; gap: 4px; }
  label:not(.checkbox) { font-size: 11px; color: #aaa; }

  .inline-fields { display: flex; gap: 12px; }
  .inline-fields .field { flex: 1; }

  select, input[type="text"] {
    background: #1e1e1e;
    border: 1px solid #444;
    border-radius: 4px;
    color: #e0e0e0;
    padding: 6px 10px;
    font-size: 13px;
    font-family: inherit;
    width: 100%;
  }
  select:focus, input:focus { outline: none; border-color: #7cb9e8; }

  .checkbox {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 13px;
    color: #ccc;
    cursor: pointer;
  }
  .checkbox input { cursor: pointer; accent-color: #0078d7; }

  .btn-row { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }

  .input-row { display: flex; gap: 8px; align-items: center; }
  .input-row input { flex: 1; }

  .btn {
    padding: 6px 14px;
    background: #2d2d2e;
    border: 1px solid #444;
    border-radius: 4px;
    color: #ddd;
    font-size: 12px;
    cursor: pointer;
    white-space: nowrap;
  }
  .btn:hover { background: #37373d; }

  .btn-sm {
    padding: 4px 10px;
    background: #333;
    border: 1px solid #555;
    border-radius: 4px;
    color: #ccc;
    font-size: 12px;
    cursor: pointer;
    white-space: nowrap;
  }
  .btn-sm:hover { background: #444; }

  .status-text { font-size: 12px; color: #aaa; }

  .save-row {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 8px 0;
  }

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
