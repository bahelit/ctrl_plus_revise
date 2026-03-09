<script>
  import { GenerateMeal, GenerateMealPrep, GenerateGroceryList, GenerateBudgetGroceryList } from '../../wailsjs/go/main/App.js';
  import { marked } from 'marked';

  let info = {
    consumers: '2',
    meal: '',
    theme: '',
    cookware: '',
    dairy: '',
    veggies: '',
    protein: '',
    freezer: '',
    herbs: '',
    pantry: '',
    allergies: '',
  };

  let result = '';
  let loading = false;
  let activeAction = null;

  function parseList(s) {
    return s.split(',').map(x => x.trim()).filter(Boolean);
  }

  function buildInfo() {
    return {
      consumers: info.consumers,
      meal: info.meal,
      theme: info.theme,
      cookware: parseList(info.cookware),
      dairy: parseList(info.dairy),
      veggies: parseList(info.veggies),
      protein: parseList(info.protein),
      freezer: parseList(info.freezer),
      herbs: parseList(info.herbs),
      pantry: parseList(info.pantry),
      allergies: parseList(info.allergies),
    };
  }

  async function run(action, fn) {
    loading = true;
    activeAction = action;
    result = '';
    try {
      result = await fn(buildInfo());
    } catch (e) {
      result = 'Error: ' + e.toString();
    } finally {
      loading = false;
    }
  }
</script>

<div class="view">
  <h1>Meal Planner</h1>

  <div class="form-grid">
    <div class="field">
      <label>Number of people</label>
      <input type="text" bind:value={info.consumers} placeholder="2" />
    </div>
    <div class="field">
      <label>Meal / Dish</label>
      <input type="text" bind:value={info.meal} placeholder="pasta, stir fry, soup…" />
    </div>
    <div class="field">
      <label>Cuisine / Theme</label>
      <input type="text" bind:value={info.theme} placeholder="Italian, Asian, Mexican…" />
    </div>
    <div class="field">
      <label>Cookware (comma separated)</label>
      <input type="text" bind:value={info.cookware} placeholder="oven, pan, instant pot…" />
    </div>
    <div class="field">
      <label>Dairy available</label>
      <input type="text" bind:value={info.dairy} placeholder="milk, cheese, butter…" />
    </div>
    <div class="field">
      <label>Vegetables</label>
      <input type="text" bind:value={info.veggies} placeholder="carrots, spinach, onions…" />
    </div>
    <div class="field">
      <label>Protein</label>
      <input type="text" bind:value={info.protein} placeholder="chicken, tofu, eggs…" />
    </div>
    <div class="field">
      <label>Freezer items</label>
      <input type="text" bind:value={info.freezer} placeholder="frozen peas, shrimp…" />
    </div>
    <div class="field">
      <label>Herbs &amp; spices</label>
      <input type="text" bind:value={info.herbs} placeholder="basil, cumin, oregano…" />
    </div>
    <div class="field">
      <label>Pantry items</label>
      <input type="text" bind:value={info.pantry} placeholder="rice, pasta, canned beans…" />
    </div>
    <div class="field full">
      <label>Allergies</label>
      <input type="text" bind:value={info.allergies} placeholder="nuts, gluten, dairy…" />
    </div>
  </div>

  <div class="actions">
    <button class="btn" on:click={() => run('meal', GenerateMeal)} disabled={loading}>
      {loading && activeAction === 'meal' ? '…' : '🍽 Get Recipe'}
    </button>
    <button class="btn" on:click={() => run('prep', GenerateMealPrep)} disabled={loading}>
      {loading && activeAction === 'prep' ? '…' : '📅 Meal Prep Plan'}
    </button>
    <button class="btn" on:click={() => run('grocery', GenerateGroceryList)} disabled={loading}>
      {loading && activeAction === 'grocery' ? '…' : '🛒 Grocery List'}
    </button>
    <button class="btn" on:click={() => run('budget', GenerateBudgetGroceryList)} disabled={loading}>
      {loading && activeAction === 'budget' ? '…' : '💰 Budget Grocery'}
    </button>
  </div>

  {#if result}
    <div class="result-box">
      <div class="result-header">
        <span>Result</span>
        <button class="btn-sm" on:click={() => navigator.clipboard.writeText(result)}>Copy</button>
      </div>
      <div class="result-content">{@html marked(result)}</div>
    </div>
  {/if}
</div>

<style>
  .view { padding: 20px 28px; overflow-y: auto; height: 100%; display: flex; flex-direction: column; gap: 16px; }
  h1 { font-size: 20px; color: #7cb9e8; }

  .form-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 10px;
  }

  .field { display: flex; flex-direction: column; gap: 4px; }
  .field.full { grid-column: 1 / -1; }

  label { font-size: 11px; color: #aaa; }

  input {
    background: #252526;
    border: 1px solid #444;
    border-radius: 5px;
    color: #e0e0e0;
    padding: 6px 10px;
    font-size: 13px;
    font-family: inherit;
  }
  input:focus { outline: none; border-color: #7cb9e8; }

  .actions { display: flex; gap: 10px; flex-wrap: wrap; }

  .btn {
    padding: 8px 16px;
    background: #2d2d2e;
    border: 1px solid #444;
    border-radius: 5px;
    color: #e0e0e0;
    font-size: 13px;
    cursor: pointer;
  }
  .btn:hover:not(:disabled) { background: #37373d; border-color: #7cb9e8; }
  .btn:disabled { opacity: 0.5; cursor: not-allowed; }

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

  .result-box {
    background: #252526;
    border: 1px solid #444;
    border-radius: 6px;
    overflow: hidden;
    flex-shrink: 0;
  }

  .result-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 12px;
    background: #2d2d2e;
    border-bottom: 1px solid #333;
    font-size: 12px;
    color: #aaa;
  }

  .result-content {
    padding: 14px;
    font-size: 13px;
    line-height: 1.6;
    max-height: 380px;
    overflow-y: auto;
  }

  :global(.result-content h1, .result-content h2, .result-content h3) { color: #7cb9e8; margin: 8px 0 4px; }
  :global(.result-content p) { margin-bottom: 6px; }
  :global(.result-content ul, .result-content ol) { padding-left: 18px; margin-bottom: 6px; }
  :global(.result-content code) { background: #333; padding: 1px 4px; border-radius: 3px; }
  :global(.result-content strong) { color: #ddd; }
</style>
