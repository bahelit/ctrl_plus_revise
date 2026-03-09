<script>
  import { AskQuestion } from '../../wailsjs/go/main/App.js';
  import { marked } from 'marked';

  let question = '';
  let response = '';
  let loading = false;
  let error = '';

  async function ask() {
    if (!question.trim()) return;
    loading = true;
    error = '';
    response = '';
    try {
      response = await AskQuestion(question);
    } catch (e) {
      error = e.toString();
    } finally {
      loading = false;
    }
  }

  function handleKey(e) {
    if (e.key === 'Enter' && (e.ctrlKey || e.shiftKey)) {
      ask();
    }
  }
</script>

<div class="view">
  <h1>Ask a Question</h1>
  <p class="hint">Type a question and press <kbd>Ctrl+Enter</kbd> or click Ask.</p>

  <textarea
    bind:value={question}
    placeholder="Type your question here…"
    rows="5"
    on:keydown={handleKey}
    disabled={loading}
  ></textarea>

  <button class="btn-primary" on:click={ask} disabled={loading || !question.trim()}>
    {loading ? 'Asking…' : 'Ask AI'}
  </button>

  {#if error}
    <div class="error">{error}</div>
  {/if}

  {#if response}
    <div class="response-box">
      <div class="response-header">
        <span>Response</span>
        <button class="btn-sm" on:click={() => navigator.clipboard.writeText(response)}>Copy</button>
      </div>
      <div class="response-content">{@html marked(response)}</div>
    </div>
  {/if}
</div>

<style>
  .view { padding: 24px 32px; overflow-y: auto; height: 100%; }
  h1 { font-size: 20px; color: #7cb9e8; margin-bottom: 6px; }
  .hint { font-size: 12px; color: #888; margin-bottom: 16px; }

  textarea {
    width: 100%;
    background: #252526;
    border: 1px solid #444;
    border-radius: 6px;
    color: #e0e0e0;
    padding: 10px 12px;
    font-size: 14px;
    font-family: inherit;
    resize: vertical;
    margin-bottom: 12px;
  }
  textarea:focus { outline: none; border-color: #7cb9e8; }

  .btn-primary {
    padding: 8px 20px;
    background: #0078d7;
    color: #fff;
    border: none;
    border-radius: 5px;
    font-size: 14px;
    cursor: pointer;
    margin-bottom: 16px;
  }
  .btn-primary:hover:not(:disabled) { background: #106ebe; }
  .btn-primary:disabled { opacity: 0.5; cursor: not-allowed; }

  .btn-sm {
    padding: 3px 10px;
    background: #333;
    color: #ccc;
    border: 1px solid #555;
    border-radius: 4px;
    font-size: 11px;
    cursor: pointer;
  }
  .btn-sm:hover { background: #444; }

  .error { color: #e07070; font-size: 13px; margin-bottom: 12px; }

  .response-box {
    background: #252526;
    border: 1px solid #444;
    border-radius: 6px;
    overflow: hidden;
  }

  .response-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 12px;
    background: #2d2d2e;
    border-bottom: 1px solid #333;
    font-size: 12px;
    color: #aaa;
  }

  .response-content {
    padding: 14px;
    font-size: 14px;
    line-height: 1.6;
    max-height: 400px;
    overflow-y: auto;
  }

  :global(.response-content h1, .response-content h2, .response-content h3) { color: #7cb9e8; margin: 10px 0 6px; }
  :global(.response-content p) { margin-bottom: 8px; }
  :global(.response-content ul, .response-content ol) { padding-left: 20px; margin-bottom: 8px; }
  :global(.response-content code) { background: #333; padding: 1px 4px; border-radius: 3px; font-size: 12px; }
  :global(.response-content pre) { background: #1a1a1a; padding: 10px; border-radius: 5px; overflow-x: auto; margin-bottom: 8px; }

  kbd { background: #333; border: 1px solid #555; border-radius: 3px; padding: 1px 5px; font-size: 11px; font-family: monospace; }
</style>
