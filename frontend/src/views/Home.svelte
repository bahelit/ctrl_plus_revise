<script>
  import { createEventDispatcher } from 'svelte';
  import { CheckOllama, GetSettings } from '../../wailsjs/go/main/App.js';

  const dispatch = createEventDispatcher();

  let ollamaStatus = 'checking';
  let shortcuts = { ctrlRevise: [], askAI: [], translate: [] };

  async function checkOllama() {
    try {
      const ok = await CheckOllama();
      ollamaStatus = ok ? 'connected' : 'disconnected';
    } catch {
      ollamaStatus = 'disconnected';
    }
  }

  async function loadInfo() {
    try {
      const s = await GetSettings();
      shortcuts = {
        ctrlRevise: s.ctrlReviseShortcut || [],
        askAI: s.askAIShortcut || [],
        translate: s.translateShortcut || [],
      };
    } catch {}
  }

  checkOllama();
  loadInfo();

  const features = [
    { id: 'ask',       label: 'Ask a Question',  desc: 'Get an AI answer to any question',       icon: '❓' },
    { id: 'chat',      label: 'Chat with AI',    desc: 'Multi-turn conversation with history',    icon: '💬' },
    { id: 'translate', label: 'Translate Text',  desc: 'Translate between 11 languages',          icon: '🌐' },
    { id: 'meal',      label: 'Meal Planner',    desc: 'AI-powered recipes and grocery lists',    icon: '🍽️' },
  ];
</script>

<div class="home">
  <header>
    <h1>Welcome to Ctrl+Revise</h1>
    <p class="subtitle">AI-powered text tools, running locally with Ollama</p>
    <div class="status" class:ok={ollamaStatus === 'connected'} class:err={ollamaStatus === 'disconnected'}>
      {#if ollamaStatus === 'checking'}
        ⏳ Checking Ollama connection…
      {:else if ollamaStatus === 'connected'}
        ✅ Ollama connected
      {:else}
        ❌ Ollama not reachable — go to Settings to configure
      {/if}
    </div>
  </header>

  <section class="features">
    {#each features as f}
      <button class="feature-card" on:click={() => dispatch('navigate', f.id)}>
        <span class="fi">{f.icon}</span>
        <div>
          <strong>{f.label}</strong>
          <p>{f.desc}</p>
        </div>
      </button>
    {/each}
  </section>

  <section class="hotkeys-card">
    <h2>Global Keyboard Shortcuts</h2>
    <p class="hint">These hotkeys work system-wide while the app is running.</p>
    <table>
      <tr><td>Revise highlighted text</td><td><kbd>{shortcuts.ctrlRevise.join(' + ').toUpperCase()}</kbd></td></tr>
      <tr><td>Ask AI with highlighted text</td><td><kbd>{shortcuts.askAI.join(' + ').toUpperCase()}</kbd></td></tr>
      <tr><td>Translate highlighted text</td><td><kbd>{shortcuts.translate.join(' + ').toUpperCase()}</kbd></td></tr>
      <tr><td>Cycle through prompt options</td><td><kbd>Alt + P</kbd></td></tr>
      <tr><td>Read highlighted text aloud</td><td><kbd>Alt + R</kbd></td></tr>
    </table>
    <p class="hint" style="margin-top: 8px">The window can be minimised; hotkeys keep working in the background.</p>
  </section>
</div>

<style>
  .home {
    padding: 24px 28px 32px;
    overflow-y: auto;
    height: 100%;
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  header { display: flex; flex-direction: column; gap: 6px; }

  h1 {
    font-size: 22px;
    font-weight: 800;
    letter-spacing: -0.02em;
    background: var(--logo-gradient);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .subtitle {
    font-size: 13px;
    color: var(--text-dim);
  }

  .status {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 5px 12px;
    border-radius: 20px;
    font-size: 12px;
    background: var(--bg-card);
    border: 1px solid var(--border);
    color: var(--text-muted);
    width: fit-content;
  }

  .status.ok  { background: var(--success-bg); color: var(--success-text); border-color: var(--success-text); }
  .status.err { background: var(--error-bg);   color: var(--error-text);   border-color: var(--error-text);   }

  /* ── Feature cards ── */
  .features {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 10px;
  }

  .feature-card {
    display: flex;
    align-items: center;
    gap: 14px;
    padding: 16px;
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: 10px;
    cursor: pointer;
    text-align: left;
    color: var(--text);
    font-family: inherit;
    transition: border-color 0.15s, background 0.15s, transform 0.15s, box-shadow 0.15s;
    box-shadow: var(--card-shadow);
  }

  .feature-card:hover {
    border-color: var(--accent);
    background: var(--bg-hover);
    transform: translateY(-2px);
    box-shadow: var(--card-shadow), 0 0 16px var(--accent-glow);
  }

  .fi { font-size: 28px; line-height: 1; }

  .feature-card strong {
    display: block;
    font-size: 14px;
    font-weight: 600;
    margin-bottom: 3px;
    color: var(--text);
  }

  .feature-card p {
    font-size: 12px;
    color: var(--text-muted);
  }

  /* ── Hotkeys ── */
  .hotkeys-card {
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 16px 18px;
    box-shadow: var(--card-shadow);
  }

  h2 {
    font-size: 10px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    color: var(--text-dim);
    margin-bottom: 8px;
  }

  .hint {
    font-size: 12px;
    color: var(--text-dim);
    margin-bottom: 10px;
  }

  table { width: 100%; border-collapse: collapse; font-size: 13px; }

  tr td {
    padding: 6px 4px;
    border-bottom: 1px solid var(--border);
    color: var(--text-muted);
  }

  tr:last-child td { border-bottom: none; }
  td:first-child { color: var(--text-muted); }
  td:last-child  { text-align: right; }

  kbd {
    background: var(--bg-input);
    border: 1px solid var(--border);
    border-radius: 4px;
    padding: 2px 7px;
    font-size: 11px;
    font-family: monospace;
    color: var(--accent);
    letter-spacing: 0.02em;
  }
</style>
