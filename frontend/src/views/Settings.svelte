<script>
  import { onMount } from 'svelte';
  import { GetSettings, SaveSettings, GetModels, GetPrompts, CheckOllama, TestOllamaURL } from '../../wailsjs/go/main/App.js';
  import Select from '../lib/Select.svelte';
  import { currentTheme, themeDefinitions } from '../lib/theme.js';

  let settings = null;
  let models = [];
  let prompts = [];
  let saved = false;
  let ollamaStatus = '';
  let testUrlStatus = '';

  onMount(async () => {
    [settings, models, prompts] = await Promise.all([GetSettings(), GetModels(), GetPrompts()]);
  });

  $: modelOptions = models.map(m => ({
    value: m.name,
    label: m.name + (m.isRecommended ? ' ⭐' : '') + (!m.isInstalled ? ' (not installed)' : '')
  }));
  $: promptOptions = prompts.map(p => ({ value: p.id, label: p.name }));

  async function save() {
    try {
      await SaveSettings(settings);
      saved = true;
      setTimeout(() => (saved = false), 2500);
    } catch (e) {
      alert('Failed to save: ' + e);
    }
  }

  async function checkOllama() {
    ollamaStatus = 'checking…';
    const ok = await CheckOllama();
    ollamaStatus = ok ? '✅ Connected' : '❌ Not reachable';
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
    <div class="loading">Loading settings…</div>
  {:else}
    <div class="page-header">
      <h1>Settings</h1>
      <p class="subtitle">Configure models, behaviour, and appearance</p>
    </div>

    <!-- AI Model -->
    <section class="card">
      <h2>AI Model</h2>

      <div class="field">
        <label>Active Model</label>
        <Select bind:value={settings.currentModel} options={modelOptions} />
      </div>

      <div class="field">
        <label>Chat Model</label>
        <Select bind:value={settings.currentChatModel} options={modelOptions} />
      </div>

      <div class="field">
        <label>Default AI Action</label>
        <Select bind:value={settings.currentPrompt} options={promptOptions} />
      </div>

    </section>

    <!-- Behaviour -->
    <section class="card">
      <h2>Behaviour</h2>
      <label class="checkbox"><input type="checkbox" bind:checked={settings.replaceHighlightedText} /> Paste AI response over highlighted text</label>
      <label class="checkbox"><input type="checkbox" bind:checked={settings.showPopUp} /> Show response notification</label>
      <label class="checkbox"><input type="checkbox" bind:checked={settings.speakAIResponse} /> Speak AI response aloud <span class="kbd">Alt+R</span></label>
      <label class="checkbox"><input type="checkbox" bind:checked={settings.stopOllamaOnShutDown} /> Stop Ollama when app closes</label>
      <label class="checkbox"><input type="checkbox" bind:checked={settings.showStartWindow} /> Show window on startup</label>
    </section>

    <!-- Translation -->
    <section class="card">
      <h2>Translation</h2>
      <div class="inline-fields">
        <div class="field">
          <label>From language</label>
          <input type="text" bind:value={settings.currentFromLang} />
        </div>
        <div class="field">
          <label>To language</label>
          <input type="text" bind:value={settings.currentToLang} />
        </div>
      </div>
    </section>

    <!-- Ollama Connection -->
    <section class="card">
      <h2>Ollama Connection</h2>
      <label class="checkbox"><input type="checkbox" bind:checked={settings.useRemoteOllama} /> Connect to remote Ollama server</label>

      {#if settings.useRemoteOllama}
        <div class="field">
          <label>Ollama URL</label>
          <div class="input-row">
            <input type="text" bind:value={settings.ollamaURL} placeholder="http://localhost:11434" />
            <button class="btn-sm" on:click={testURL}>Test</button>
          </div>
          {#if testUrlStatus}<span class="status-text">{testUrlStatus}</span>{/if}
        </div>
      {/if}

      <div class="btn-row">
        <button class="btn" on:click={checkOllama}>🔌 Check Connection</button>
        {#if ollamaStatus}<span class="status-text">{ollamaStatus}</span>{/if}
      </div>
    </section>

    <!-- Theme -->
    <section class="card">
      <h2>Appearance</h2>
      <p class="section-hint">Choose a colour theme. Your choice is saved automatically.</p>
      <div class="theme-grid">
        {#each Object.entries(themeDefinitions) as [id, def]}
          <button
            class="theme-card"
            class:active={$currentTheme === id}
            on:click={() => currentTheme.set(id)}
            style="--tc-accent: {def.accent}"
          >
            <span class="theme-dot" style="background: {def.accent}; box-shadow: 0 0 8px {def.accent}88"></span>
            <span class="theme-name">{def.name}</span>
            {#if $currentTheme === id}
              <span class="theme-check">✓</span>
            {/if}
          </button>
        {/each}
      </div>
    </section>

    <!-- Save -->
    <div class="save-row">
      <button class="btn-primary" on:click={save}>Save Settings</button>
      {#if saved}<span class="saved-msg">✅ Saved</span>{/if}
    </div>
  {/if}
</div>

<style>
  .view {
    padding: 20px 28px 32px;
    overflow-y: auto;
    height: 100%;
    display: flex;
    flex-direction: column;
    gap: 12px;
    scrollbar-gutter: stable;
  }

  .loading {
    padding: 32px;
    color: var(--text-muted);
    font-size: 14px;
  }

  /* ── Page header ── */
  .page-header { margin-bottom: 4px; }

  h1 {
    font-size: 22px;
    font-weight: 800;
    letter-spacing: -0.02em;
    background: var(--logo-gradient);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    margin-bottom: 2px;
  }

  .subtitle {
    font-size: 12px;
    color: var(--text-dim);
  }

  /* ── Card / section ── */
  .card {
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 16px 18px;
    display: flex;
    flex-direction: column;
    gap: 10px;
    box-shadow: var(--card-shadow);
    transition: border-color 0.2s;
  }

  .card:hover {
    border-color: var(--text-dim);
  }

  h2 {
    font-size: 10px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    color: var(--text-dim);
    margin-bottom: 2px;
  }

  .section-hint {
    font-size: 12px;
    color: var(--text-dim);
    margin-top: -4px;
  }

  /* ── Fields ── */
  .field {
    display: flex;
    flex-direction: column;
    gap: 5px;
  }

  label:not(.checkbox) {
    font-size: 11px;
    font-weight: 600;
    letter-spacing: 0.03em;
    color: var(--text-muted);
  }

  .inline-fields {
    display: flex;
    gap: 14px;
  }

  .inline-fields .field { flex: 1; }

  /* ── Native inputs (text) ── */
  input[type="text"] {
    background: var(--bg-input);
    border: 1px solid var(--border);
    border-radius: 6px;
    color: var(--text);
    padding: 8px 12px;
    font-size: 13px;
    font-family: inherit;
    width: 100%;
    transition: border-color 0.15s, box-shadow 0.15s;
  }

  input[type="text"]:focus {
    outline: none;
    border-color: var(--border-focus);
    box-shadow: 0 0 0 3px var(--accent-glow);
  }

  input[type="text"]::placeholder {
    color: var(--text-dim);
  }

  /* ── Checkbox ── */
  .checkbox {
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 13px;
    color: var(--text-muted);
    cursor: pointer;
    padding: 2px 0;
    transition: color 0.15s;
  }

  .checkbox:hover { color: var(--text); }

  .checkbox input {
    cursor: pointer;
    accent-color: var(--accent);
    width: 15px;
    height: 15px;
    flex-shrink: 0;
  }

  .kbd {
    font-size: 10px;
    background: var(--bg-input);
    border: 1px solid var(--border);
    border-radius: 4px;
    padding: 1px 5px;
    color: var(--text-dim);
    font-family: monospace;
    margin-left: 2px;
  }

  /* ── Button rows ── */
  .btn-row {
    display: flex;
    align-items: center;
    gap: 10px;
    flex-wrap: wrap;
    margin-top: 2px;
  }

  .input-row {
    display: flex;
    gap: 8px;
    align-items: center;
  }

  .input-row input { flex: 1; }

  /* ── Buttons ── */
  .btn {
    padding: 7px 14px;
    background: var(--btn-bg);
    border: 1px solid var(--btn-border);
    border-radius: 6px;
    color: var(--text-muted);
    font-size: 12px;
    font-family: inherit;
    cursor: pointer;
    white-space: nowrap;
    transition: background 0.15s, color 0.15s, border-color 0.15s;
  }

  .btn:hover {
    background: var(--btn-bg-hover);
    color: var(--text);
    border-color: var(--text-dim);
  }

  .btn-sm {
    padding: 6px 12px;
    background: var(--btn-bg);
    border: 1px solid var(--btn-border);
    border-radius: 6px;
    color: var(--text-muted);
    font-size: 12px;
    font-family: inherit;
    cursor: pointer;
    white-space: nowrap;
    transition: background 0.15s, color 0.15s;
  }

  .btn-sm:hover {
    background: var(--btn-bg-hover);
    color: var(--text);
  }

  .status-text {
    font-size: 12px;
    color: var(--text-muted);
  }

  /* ── Theme grid ── */
  .theme-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 8px;
  }

  .theme-card {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 14px;
    background: var(--bg-input);
    border: 1px solid var(--border);
    border-radius: 8px;
    cursor: pointer;
    font-family: inherit;
    font-size: 13px;
    color: var(--text-muted);
    transition: background 0.15s, border-color 0.15s, color 0.15s;
    position: relative;
    text-align: left;
  }

  .theme-card:hover {
    background: var(--bg-hover);
    border-color: var(--tc-accent);
    color: var(--text);
  }

  .theme-card.active {
    border-color: var(--tc-accent);
    background: var(--bg-active);
    color: var(--text);
    font-weight: 600;
    box-shadow: 0 0 14px var(--accent-glow);
  }

  .theme-dot {
    width: 14px;
    height: 14px;
    border-radius: 50%;
    flex-shrink: 0;
  }

  .theme-name { flex: 1; }

  .theme-check {
    color: var(--tc-accent);
    font-size: 14px;
    font-weight: 700;
  }

  /* ── Save row ── */
  .save-row {
    display: flex;
    align-items: center;
    gap: 14px;
    padding: 4px 0 4px;
  }

  .btn-primary {
    padding: 10px 28px;
    background: var(--btn-primary);
    color: #fff;
    border: none;
    border-radius: 7px;
    font-size: 14px;
    font-weight: 600;
    font-family: inherit;
    cursor: pointer;
    transition: background 0.15s, box-shadow 0.15s, transform 0.1s;
    box-shadow: 0 2px 12px var(--accent-glow);
  }

  .btn-primary:hover {
    background: var(--btn-primary-hover);
    box-shadow: 0 4px 20px var(--accent-glow);
    transform: translateY(-1px);
  }

  .btn-primary:active {
    transform: translateY(0);
  }

  .saved-msg {
    font-size: 13px;
    color: var(--success-text);
    animation: fadeIn 0.2s ease;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateX(-8px); }
    to   { opacity: 1; transform: translateX(0); }
  }
</style>
