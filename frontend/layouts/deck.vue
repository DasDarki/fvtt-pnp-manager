<script setup lang="ts">
const ui = useUiStore()
const campaign = useCampaignStore()

onMounted(() => {
  campaign.ensure().catch(() => {})
})
</script>

<template>
  <div class="deck">
    <AtmosBackdrop />

    <button class="menu" :aria-label="'Menu'" @click="ui.toggleSidebar()">
      <Icon name="lucide:menu" />
    </button>

    <Transition name="veil">
      <div v-if="ui.sidebarOpen" class="drawer-veil" @click="ui.closeSidebar()" />
    </Transition>

    <AppSidebar class="side-slot" :class="{ open: ui.sidebarOpen }" />

    <main class="main">
      <TopCommandBar />
      <div class="content">
        <slot />
      </div>
    </main>

    <CommandPalette />
  </div>
</template>

<style lang="scss" scoped>
.deck {
  position: relative;
  z-index: 2;
  display: grid;
  grid-template-columns: 284px 1fr;
  height: 100vh;
}

.main {
  overflow-y: auto;
  height: 100vh;
}

.content {
  padding: 6px 34px 60px;
}

.menu {
  position: fixed;
  z-index: 25;
  top: 14px;
  left: 16px;
  display: none;
  width: 42px;
  height: 42px;
  border-radius: 12px;
  border: 1px solid var(--line-strong);
  background: var(--surface-2);
  color: var(--ink);
  place-items: center;
  cursor: pointer;

  :deep(svg) { width: 20px; height: 20px; }
}

.drawer-veil {
  position: fixed;
  inset: 0;
  z-index: 45;
  background: rgba(5, 4, 10, 0.5);
  backdrop-filter: blur(4px);
}
.veil-enter-active,
.veil-leave-active { transition: opacity 0.25s; }
.veil-enter-from,
.veil-leave-to { opacity: 0; }

@media (max-width: 820px) {
  .deck { grid-template-columns: 1fr; }
  .side-slot {
    position: fixed;
    z-index: 50;
    width: 284px;
    height: 100vh;
    transform: translateX(-100%);
    transition: transform 0.3s;
  }
  .side-slot.open { transform: none; }
  .menu { display: grid; }
  .content { padding: 6px 16px 50px; }
}
</style>
