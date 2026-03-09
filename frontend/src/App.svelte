<script>
  import { onMount } from 'svelte';
  import Home from './views/Home.svelte';
  import Chat from './views/Chat.svelte';
  import AskQuestion from './views/AskQuestion.svelte';
  import Translate from './views/Translate.svelte';
  import MealPlanner from './views/MealPlanner.svelte';
  import Settings from './views/Settings.svelte';
  import Shortcuts from './views/Shortcuts.svelte';
  import { EventsOn } from '../wailsjs/runtime/runtime.js';
  import { currentTheme, themeDefinitions } from './lib/theme.js';

  let currentView = 'home';
  let notification = null;

  const views = [
    { id: 'home',       label: 'Home',           icon: '🏠' },
    { id: 'ask',        label: 'Ask Question',   icon: '❓' },
    { id: 'chat',       label: 'Chat with AI',   icon: '💬' },
    { id: 'translate',  label: 'Translate',      icon: '🌐' },
    { id: 'meal',       label: 'Meal Planner',   icon: '🍽️' },
    { id: 'settings',   label: 'Settings',       icon: '⚙️' },
    { id: 'shortcuts',  label: 'Shortcuts',      icon: '⌨️' },
  ];

  onMount(() => {
    EventsOn('notification', (data) => {
      notification = data;
      setTimeout(() => { notification = null; }, 4000);
    });
  });

  $: theme = themeDefinitions[$currentTheme];
  $: cssVars = Object.entries(theme.vars).map(([k, v]) => `${k}:${v}`).join(';');
</script>

<div class="app" style={cssVars}>
  <nav class="sidebar">
    <div class="logo">
      <span class="logo-text">Ctrl+Revise</span>
      <div class="logo-bar"></div>
    </div>

    <div class="nav-list">
      {#each views as view}
        <button
          class="nav-btn"
          class:active={currentView === view.id}
          on:click={() => currentView = view.id}
        >
          <span class="icon">{view.icon}</span>
          <span class="label">{view.label}</span>
          {#if currentView === view.id}
            <span class="active-dot"></span>
          {/if}
        </button>
      {/each}
    </div>

  </nav>

  <main class="content">
    {#if currentView === 'home'}
      <Home on:navigate={e => currentView = e.detail} />
    {:else if currentView === 'ask'}
      <AskQuestion />
    {:else if currentView === 'chat'}
      <Chat />
    {:else if currentView === 'translate'}
      <Translate />
    {:else if currentView === 'meal'}
      <MealPlanner />
    {:else if currentView === 'settings'}
      <Settings />
    {:else if currentView === 'shortcuts'}
      <Shortcuts />
    {/if}
  </main>

  {#if notification}
    <div class="toast">
      <strong>{notification.title}</strong>
      <p>{notification.content}</p>
    </div>
  {/if}
</div>

<style>
  :global(body) {
    background: #000;
  }

  :global(*::-webkit-scrollbar) { width: 6px; height: 6px; }
  :global(*::-webkit-scrollbar-track) { background: transparent; }
  :global(*::-webkit-scrollbar-thumb) {
    background: var(--border, #333);
    border-radius: 3px;
  }
  :global(*::-webkit-scrollbar-thumb:hover) {
    background: var(--text-dim, #555);
  }

  .app {
    display: flex;
    height: 100vh;
    background: var(--bg-app);
    color: var(--text);
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, Ubuntu, Cantarell, sans-serif;
  }

  /* ── Sidebar ── */
  .sidebar {
    width: 188px;
    min-width: 188px;
    background: var(--sidebar-bg);
    display: flex;
    flex-direction: column;
    border-right: 1px solid var(--border);
    position: relative;
  }

  /* subtle inner glow on right edge */
  .sidebar::after {
    content: '';
    position: absolute;
    inset: 0 0 0 auto;
    width: 1px;
    background: linear-gradient(180deg, transparent, var(--accent-glow) 40%, transparent);
    pointer-events: none;
  }

  /* ── Logo ── */
  .logo {
    padding: 18px 14px 14px;
    user-select: none;
  }

  .logo-text {
    display: block;
    font-size: 15px;
    font-weight: 800;
    letter-spacing: -0.02em;
    background: var(--logo-gradient);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .logo-bar {
    margin-top: 10px;
    height: 1px;
    background: linear-gradient(90deg, var(--accent-glow), transparent);
  }

  /* ── Nav ── */
  .nav-list {
    flex: 1;
    display: flex;
    flex-direction: column;
    padding: 6px 8px;
    gap: 2px;
  }

  .nav-btn {
    position: relative;
    display: flex;
    align-items: center;
    gap: 9px;
    padding: 9px 10px;
    background: none;
    border: none;
    border-radius: 7px;
    color: var(--text-muted);
    cursor: pointer;
    text-align: left;
    font-size: 13px;
    font-family: inherit;
    transition: background 0.15s, color 0.15s;
    width: 100%;
  }

  .nav-btn:hover {
    background: var(--bg-hover);
    color: var(--text);
  }

  .nav-btn.active {
    background: var(--bg-active);
    color: var(--accent);
    font-weight: 600;
  }

  .active-dot {
    position: absolute;
    left: 0;
    top: 50%;
    transform: translateY(-50%);
    width: 3px;
    height: 60%;
    background: var(--accent);
    border-radius: 0 3px 3px 0;
    box-shadow: 0 0 8px var(--accent-glow);
  }

  .icon {
    font-size: 15px;
    width: 20px;
    text-align: center;
    flex-shrink: 0;
  }

  /* ── Main content ── */
  .content {
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    background: var(--bg-app);
  }

  /* ── Toast ── */
  .toast {
    position: fixed;
    bottom: 20px;
    right: 20px;
    background: var(--bg-card);
    border: 1px solid var(--border-focus);
    border-radius: 10px;
    padding: 14px 18px;
    max-width: 320px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5), 0 0 0 1px var(--accent-glow);
    z-index: 1000;
    animation: slideIn 0.2s ease;
  }

  .toast strong {
    display: block;
    color: var(--accent);
    margin-bottom: 4px;
    font-size: 13px;
  }

  .toast p {
    font-size: 12px;
    color: var(--text-muted);
    white-space: pre-wrap;
  }

  @keyframes slideIn {
    from { transform: translateX(24px); opacity: 0; }
    to   { transform: translateX(0);    opacity: 1; }
  }
</style>
