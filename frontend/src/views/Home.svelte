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
    { id: 'ask',       label: 'Ask a Question',  desc: 'Get an AI answer to any question', icon: '❓' },
    { id: 'chat',      label: 'Chat with AI',    desc: 'Multi-turn conversation with history', icon: '💬' },
    { id: 'translate', label: 'Translate Text',  desc: 'Translate between 11 languages', icon: '🌐' },
    { id: 'meal',      label: 'Meal Planner',    desc: 'AI-powered recipes and grocery lists', icon: '🍽️' },
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

  <section class="hotkeys">
    <h2>Global Keyboard Shortcuts</h2>
    <p class="hint">These hotkeys work system-wide while the app is running.</p>
    <table>
      <tr><td>Revise highlighted text</td><td><kbd>{shortcuts.ctrlRevise.join(' + ').toUpperCase()}</kbd></td></tr>
      <tr><td>Ask AI with highlighted text</td><td><kbd>{shortcuts.askAI.join(' + ').toUpperCase()}</kbd></td></tr>
      <tr><td>Translate highlighted text</td><td><kbd>{shortcuts.translate.join(' + ').toUpperCase()}</kbd></td></tr>
      <tr><td>Cycle through prompt options</td><td><kbd>Alt + P</kbd></td></tr>
      <tr><td>Read highlighted text aloud</td><td><kbd>Alt + R</kbd></td></tr>
    </table>
    <p class="hint" style="margin-top:8px">The window can be minimised; hotkeys keep working in the background.</p>
  </section>
</div>

<style>
  .home {
    padding: 24px 32px;
    overflow-y: auto;
    height: 100%;
  }

  header { margin-bottom: 24px; }

  h1 { font-size: 22px; color: #7cb9e8; margin-bottom: 4px; }
  .subtitle { color: #aaa; font-size: 13px; margin-bottom: 12px; }

  .status {
    display: inline-block;
    padding: 4px 10px;
    border-radius: 4px;
    font-size: 12px;
    background: #333;
    color: #aaa;
  }
  .status.ok  { background: #1a3a1a; color: #6dbf6d; }
  .status.err { background: #3a1a1a; color: #e07070; }

  .features {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
    margin-bottom: 24px;
  }

  .feature-card {
    display: flex;
    align-items: center;
    gap: 14px;
    padding: 16px;
    background: #252526;
    border: 1px solid #333;
    border-radius: 8px;
    cursor: pointer;
    text-align: left;
    color: #e0e0e0;
    transition: border-color 0.15s, background 0.15s;
  }

  .feature-card:hover {
    border-color: #7cb9e8;
    background: #2a2a2a;
  }

  .fi { font-size: 28px; }
  .feature-card strong { display: block; font-size: 14px; margin-bottom: 2px; }
  .feature-card p { font-size: 12px; color: #aaa; }

  .hotkeys { background: #252526; border: 1px solid #333; border-radius: 8px; padding: 16px; }
  h2 { font-size: 14px; color: #aaa; margin-bottom: 8px; }
  .hint { font-size: 12px; color: #888; margin-bottom: 8px; }

  table { width: 100%; border-collapse: collapse; font-size: 13px; }
  tr td { padding: 5px 8px; border-bottom: 1px solid #2d2d2d; }
  tr:last-child td { border-bottom: none; }
  td:first-child { color: #ccc; }
  td:last-child { text-align: right; }

  kbd {
    background: #333;
    border: 1px solid #555;
    border-radius: 3px;
    padding: 2px 6px;
    font-size: 11px;
    font-family: monospace;
    color: #ddd;
  }
</style>
