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
</script>

<div class="app">
  <nav class="sidebar">
    <div class="logo">Ctrl+Revise</div>
    {#each views as view}
      <button
        class="nav-btn"
        class:active={currentView === view.id}
        on:click={() => currentView = view.id}
      >
        <span class="icon">{view.icon}</span>
        <span class="label">{view.label}</span>
      </button>
    {/each}
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
  .app {
    display: flex;
    height: 100vh;
    background: #1e1e1e;
    color: #e0e0e0;
  }

  .sidebar {
    width: 180px;
    min-width: 180px;
    background: #252526;
    display: flex;
    flex-direction: column;
    padding: 0;
    border-right: 1px solid #333;
  }

  .logo {
    padding: 16px 12px;
    font-size: 14px;
    font-weight: 700;
    color: #7cb9e8;
    border-bottom: 1px solid #333;
    user-select: none;
  }

  .nav-btn {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 12px;
    background: none;
    border: none;
    color: #ccc;
    cursor: pointer;
    text-align: left;
    font-size: 13px;
    transition: background 0.15s;
    width: 100%;
  }

  .nav-btn:hover {
    background: #2d2d2e;
    color: #fff;
  }

  .nav-btn.active {
    background: #37373d;
    color: #7cb9e8;
    border-left: 3px solid #7cb9e8;
  }

  .icon {
    font-size: 16px;
    width: 20px;
    text-align: center;
  }

  .content {
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  .toast {
    position: fixed;
    bottom: 20px;
    right: 20px;
    background: #323233;
    border: 1px solid #555;
    border-radius: 6px;
    padding: 12px 16px;
    max-width: 320px;
    box-shadow: 0 4px 12px rgba(0,0,0,0.4);
    z-index: 1000;
    animation: slideIn 0.2s ease;
  }

  .toast strong {
    display: block;
    color: #7cb9e8;
    margin-bottom: 4px;
    font-size: 13px;
  }

  .toast p {
    font-size: 12px;
    color: #ccc;
    white-space: pre-wrap;
  }

  @keyframes slideIn {
    from { transform: translateX(20px); opacity: 0; }
    to   { transform: translateX(0);    opacity: 1; }
  }
</style>
