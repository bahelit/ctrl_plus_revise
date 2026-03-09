<script>
  import { onMount } from 'svelte';
  import { GetChats, SendChatMessage, DeleteChat } from '../../wailsjs/go/main/App.js';
  import { marked } from 'marked';

  let chats = [];
  let activeChatId = null;
  let activeChat = null;
  let message = '';
  let loading = false;
  let currentContext = [];

  onMount(loadChats);

  async function loadChats() {
    try {
      chats = await GetChats() || [];
    } catch {}
  }

  function selectChat(chat) {
    activeChatId = chat.id;
    activeChat = { ...chat };
    currentContext = chat.context || [];
  }

  function newChat() {
    activeChatId = null;
    activeChat = { id: 0, title: 'New Chat', questions: [], responses: [] };
    currentContext = [];
  }

  async function send() {
    if (!message.trim() || loading) return;
    loading = true;
    const isNew = !activeChatId || activeChatId === 0;
    const q = message;
    message = '';

    try {
      const resp = await SendChatMessage(activeChatId || 0, isNew, q, currentContext);
      currentContext = resp.context || [];

      if (!activeChat) activeChat = { id: 0, title: q.slice(0, 14), questions: [], responses: [] };
      activeChat.questions = [...activeChat.questions, q];
      activeChat.responses = [...activeChat.responses, resp.response];

      if (isNew && resp.chatId > 0) {
        activeChatId = resp.chatId;
        activeChat.id = resp.chatId;
        await loadChats();
      }
    } catch (e) {
      if (activeChat) {
        activeChat.responses = [...activeChat.responses, 'Error: ' + e.toString()];
      }
    } finally {
      loading = false;
    }
  }

  async function deleteChat(id) {
    try {
      await DeleteChat(id);
      if (activeChatId === id) { activeChatId = null; activeChat = null; }
      await loadChats();
    } catch {}
  }

  function handleKey(e) {
    if (e.key === 'Enter' && (e.ctrlKey || e.shiftKey)) send();
  }
</script>

<div class="chat-layout">
  <aside class="chat-list">
    <div class="chat-list-header">
      <span>Chats</span>
      <button class="btn-new" on:click={newChat}>+ New</button>
    </div>
    {#if chats.length === 0}
      <p class="empty">No saved chats</p>
    {/if}
    {#each chats as chat}
      <div class="chat-item" class:active={activeChatId === chat.id} on:click={() => selectChat(chat)}>
        <span class="chat-title">{chat.title || 'Untitled'}</span>
        <button class="del-btn" on:click|stopPropagation={() => deleteChat(chat.id)}>✕</button>
      </div>
    {/each}
  </aside>

  <div class="chat-main">
    {#if !activeChat}
      <div class="placeholder">
        <p>Select a chat or start a new one</p>
        <button class="btn-primary" on:click={newChat}>New Chat</button>
      </div>
    {:else}
      <div class="messages">
        {#each activeChat.questions as q, i}
          <div class="msg user">
            <div class="bubble">{q}</div>
          </div>
          {#if activeChat.responses[i] !== undefined}
            <div class="msg ai">
              <div class="bubble ai-bubble">{@html marked(activeChat.responses[i])}</div>
            </div>
          {/if}
        {/each}
        {#if loading}
          <div class="msg ai"><div class="bubble ai-bubble loading">Thinking…</div></div>
        {/if}
      </div>

      <div class="input-row">
        <textarea
          bind:value={message}
          placeholder="Type a message… (Ctrl+Enter to send)"
          rows="3"
          on:keydown={handleKey}
          disabled={loading}
        ></textarea>
        <button class="btn-primary" on:click={send} disabled={loading || !message.trim()}>Send</button>
      </div>
    {/if}
  </div>
</div>

<style>
  .chat-layout { display: flex; height: 100%; }

  .chat-list {
    width: 200px;
    min-width: 200px;
    background: #252526;
    border-right: 1px solid #333;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .chat-list-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 12px;
    border-bottom: 1px solid #333;
    font-size: 12px;
    color: #aaa;
  }

  .btn-new {
    background: #0078d7;
    color: #fff;
    border: none;
    border-radius: 4px;
    padding: 3px 8px;
    font-size: 11px;
    cursor: pointer;
  }
  .btn-new:hover { background: #106ebe; }

  .empty { padding: 12px; color: #666; font-size: 12px; }

  .chat-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 8px 12px;
    cursor: pointer;
    font-size: 12px;
    color: #ccc;
    border-bottom: 1px solid #2d2d2d;
  }
  .chat-item:hover { background: #2d2d2e; }
  .chat-item.active { background: #37373d; color: #7cb9e8; }

  .chat-title { flex: 1; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

  .del-btn {
    background: none;
    border: none;
    color: #666;
    cursor: pointer;
    font-size: 11px;
    padding: 0 2px;
    margin-left: 4px;
  }
  .del-btn:hover { color: #e07070; }

  .chat-main {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .placeholder {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 12px;
    color: #888;
    font-size: 14px;
  }

  .messages {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .msg { display: flex; }
  .msg.user { justify-content: flex-end; }
  .msg.ai   { justify-content: flex-start; }

  .bubble {
    max-width: 70%;
    padding: 10px 14px;
    border-radius: 8px;
    font-size: 13px;
    line-height: 1.5;
  }

  .msg.user .bubble { background: #0078d7; color: #fff; }
  .ai-bubble { background: #252526; border: 1px solid #333; color: #e0e0e0; }

  .loading { color: #888; font-style: italic; }

  .input-row {
    display: flex;
    gap: 8px;
    padding: 12px 16px;
    border-top: 1px solid #333;
    background: #1e1e1e;
  }

  textarea {
    flex: 1;
    background: #252526;
    border: 1px solid #444;
    border-radius: 6px;
    color: #e0e0e0;
    padding: 8px 10px;
    font-size: 13px;
    font-family: inherit;
    resize: none;
  }
  textarea:focus { outline: none; border-color: #7cb9e8; }

  .btn-primary {
    padding: 8px 16px;
    background: #0078d7;
    color: #fff;
    border: none;
    border-radius: 5px;
    font-size: 13px;
    cursor: pointer;
    align-self: flex-end;
  }
  .btn-primary:hover:not(:disabled) { background: #106ebe; }
  .btn-primary:disabled { opacity: 0.5; cursor: not-allowed; }

  :global(.ai-bubble p) { margin-bottom: 6px; }
  :global(.ai-bubble ul, .ai-bubble ol) { padding-left: 18px; margin-bottom: 6px; }
  :global(.ai-bubble code) { background: #1a1a1a; padding: 1px 4px; border-radius: 3px; font-size: 12px; }
  :global(.ai-bubble pre) { background: #1a1a1a; padding: 8px; border-radius: 4px; overflow-x: auto; }
  :global(.ai-bubble h1, .ai-bubble h2, .ai-bubble h3) { color: #7cb9e8; margin: 6px 0 4px; }
</style>
