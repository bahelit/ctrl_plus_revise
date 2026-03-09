<script>
  export let value;
  export let options = []; // [{ value, label }]
  export let placeholder = '—';

  let open = false;
  let container;

  $: selectedLabel = options.find(o => o.value === value)?.label ?? placeholder;

  function pick(opt) {
    value = opt.value;
    open = false;
  }

  function onWindowClick() { open = false; }

  function onKeydown(e) {
    if (!open) return;
    if (e.key === 'Escape') { open = false; return; }
    const idx = options.findIndex(o => o.value === value);
    if (e.key === 'ArrowDown') { e.preventDefault(); pick(options[Math.min(idx + 1, options.length - 1)]); }
    if (e.key === 'ArrowUp')   { e.preventDefault(); pick(options[Math.max(idx - 1, 0)]); }
    if (e.key === 'Enter')     { open = false; }
  }
</script>

<svelte:window on:click={onWindowClick} on:keydown={onKeydown} />

<div class="select-wrap" bind:this={container} on:click|stopPropagation>
  <button
    class="trigger"
    class:open
    on:click={() => (open = !open)}
    aria-haspopup="listbox"
    aria-expanded={open}
  >
    <span class="val">{selectedLabel}</span>
    <svg class="chevron" class:open viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
      <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
    </svg>
  </button>

  {#if open}
    <div class="dropdown" role="listbox">
      {#each options as opt}
        <div
          class="option"
          class:selected={opt.value === value}
          role="option"
          aria-selected={opt.value === value}
          on:click={() => pick(opt)}
        >
          {opt.label}
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .select-wrap {
    position: relative;
    width: 100%;
  }

  .trigger {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
    padding: 8px 12px;
    background: var(--bg-input);
    border: 1px solid var(--border);
    border-radius: 6px;
    color: var(--text);
    font-size: 13px;
    font-family: inherit;
    cursor: pointer;
    transition: border-color 0.15s, box-shadow 0.15s;
    text-align: left;
  }

  .trigger:hover {
    border-color: var(--text-dim);
  }

  .trigger.open,
  .trigger:focus-visible {
    outline: none;
    border-color: var(--border-focus);
    box-shadow: 0 0 0 3px var(--accent-glow);
  }

  .val {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .chevron {
    width: 16px;
    height: 16px;
    flex-shrink: 0;
    color: var(--text-dim);
    transition: transform 0.2s ease, color 0.15s;
  }

  .trigger:hover .chevron,
  .trigger.open .chevron {
    color: var(--accent);
  }

  .chevron.open {
    transform: rotate(180deg);
  }

  .dropdown {
    position: absolute;
    top: calc(100% + 6px);
    left: 0;
    right: 0;
    background: var(--bg-card);
    border: 1px solid var(--border-focus);
    border-radius: 8px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5), 0 0 0 1px var(--accent-glow);
    z-index: 200;
    overflow: hidden;
    max-height: 260px;
    overflow-y: auto;
    animation: dropIn 0.15s ease;
  }

  .dropdown::-webkit-scrollbar { width: 6px; }
  .dropdown::-webkit-scrollbar-track { background: transparent; }
  .dropdown::-webkit-scrollbar-thumb {
    background: var(--border);
    border-radius: 3px;
  }

  @keyframes dropIn {
    from { opacity: 0; transform: translateY(-8px) scale(0.98); }
    to   { opacity: 1; transform: translateY(0)   scale(1);    }
  }

  .option {
    padding: 9px 14px;
    color: var(--text-muted);
    font-size: 13px;
    cursor: pointer;
    transition: background 0.1s, color 0.1s;
    user-select: none;
  }

  .option:hover {
    background: var(--bg-hover);
    color: var(--text);
  }

  .option.selected {
    color: var(--accent);
    background: var(--bg-active);
    font-weight: 500;
  }

  .option.selected::before {
    content: '✓ ';
    opacity: 0.8;
  }
</style>
