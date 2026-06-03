<script setup lang="ts">
const { t } = useI18n()
const ui = useUiStore()
const auth = useAuthStore()
const campaignStore = useCampaignStore()
const providerStore = useProviderStore()

const userLabel = computed(() => auth.user?.name || auth.user?.email || 'Gast')
const userInitial = computed(() => (auth.user?.name || auth.user?.email || '?').charAt(0).toUpperCase())

async function logout() {
  await auth.logout()
  campaignStore.reset()
  providerStore.reset()
  await navigateTo('/login')
}

const primary = [
  { key: 'dashboard', icon: 'lucide:layout-dashboard', to: '/dashboard' },
  { key: 'characters', icon: 'lucide:users', to: '/characters', count: 28 },
]
const tree = [
  { key: 'allies', icon: 'lucide:folder', to: '#', count: 11 },
  { key: 'foes', icon: 'lucide:folder', to: '#', count: 17 },
]
const rest = [
  { key: 'scenes', icon: 'lucide:castle', to: '/scenes', count: 9 },
  { key: 'items', icon: 'lucide:gem', to: '/items', count: 54 },
  { key: 'images', icon: 'lucide:image', to: '/images', count: 12 },
  { key: 'memories', icon: 'lucide:sparkles', to: '/memories', count: 37 },
]
const system = [
  { key: 'dalle', icon: 'lucide:wand-sparkles', to: '/dalle' },
  { key: 'providers', icon: 'lucide:key-round', to: '/settings/providers' },
  { key: 'apikeys', icon: 'lucide:key-square', to: '/settings/api-keys' },
  { key: 'foundry', icon: 'lucide:plug-zap', to: '/settings/foundry' },
]
</script>

<template>
  <aside class="side">
    <div class="brand">
      <span class="sig" aria-hidden="true">
        <svg viewBox="0 0 24 24" fill="none" stroke="#fff" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round">
          <path d="M12 2l2.4 6.6L21 9l-5 4.2L17.5 21 12 17l-5.5 4L8 13.2 3 9l6.6-.4z" />
        </svg>
      </span>
      <span>Aetherwright<small>ARCANE DECK</small></span>
    </div>

    <CampaignSwitcher />

    <nav class="nav">
      <div class="grp">{{ t('nav.group.campaign') }}</div>
      <NuxtLink
        v-for="item in primary"
        :key="item.key"
        :to="item.to"
        class="nl"
        active-class="on"
        @click="ui.closeSidebar()"
      >
        <Icon :name="item.icon" />
        <span>{{ t(`nav.${item.key}`) }}</span>
        <span v-if="item.count" class="count">{{ item.count }}</span>
      </NuxtLink>

      <div class="tree">
        <NuxtLink v-for="item in tree" :key="item.key" :to="item.to" class="nl sm" @click="ui.closeSidebar()">
          <Icon :name="item.icon" />
          <span>{{ t(`nav.${item.key}`) }}</span>
          <span class="count">{{ item.count }}</span>
        </NuxtLink>
      </div>

      <NuxtLink v-for="item in rest" :key="item.key" :to="item.to" class="nl" active-class="on" @click="ui.closeSidebar()">
        <Icon :name="item.icon" />
        <span>{{ t(`nav.${item.key}`) }}</span>
        <span class="count">{{ item.count }}</span>
      </NuxtLink>

      <div class="grp">{{ t('nav.group.system') }}</div>
      <NuxtLink v-for="item in system" :key="item.key" :to="item.to" class="nl" @click="ui.closeSidebar()">
        <Icon :name="item.icon" />
        <span>{{ t(`nav.${item.key}`) }}</span>
      </NuxtLink>
    </nav>

    <div class="sidefoot">
      <div class="sync">
        <span class="ld" />
        <span>
          <b>{{ t('sidebar.connected') }}</b>
          <small>the-forge · v13</small>
        </span>
      </div>
      <div class="user">
        <span class="av" aria-hidden="true"><span>{{ userInitial }}</span></span>
        <span class="ui">
          <b>{{ userLabel }}</b>
          <small>{{ t('sidebar.role') }}</small>
        </span>
        <button class="gear" :title="t('sidebar.logout')" @click="logout">
          <Icon name="lucide:log-out" />
        </button>
      </div>
    </div>
  </aside>
</template>

<style lang="scss" scoped>
.side {
  display: flex;
  flex-direction: column;
  height: 100vh;
  border-right: 1px solid var(--line);
  background: linear-gradient(180deg, var(--surface), transparent);
  backdrop-filter: blur(14px);
  overflow: hidden;
}

.brand {
  display: flex;
  align-items: center;
  gap: 11px;
  padding: 22px 20px 18px;
  font-family: var(--font-deco);
  font-weight: 900;
  letter-spacing: 0.07em;
  font-size: 0.94rem;

  .sig {
    width: 32px;
    height: 32px;
    display: grid;
    place-items: center;
    border-radius: 9px;
    background: var(--grad-arcane);
    box-shadow: var(--glow-primary);
    position: relative;
    flex: none;

    &::before {
      content: '';
      position: absolute;
      inset: 1.5px;
      border-radius: 7px;
      background: var(--void-2);
    }
    svg {
      position: relative;
      z-index: 1;
      width: 17px;
      height: 17px;
    }
  }
  small {
    font-family: var(--font-mono);
    font-weight: 400;
    font-size: 0.56rem;
    letter-spacing: 0.28em;
    color: var(--primary);
    display: block;
    margin-top: 2px;
  }
}

.campaign {
  margin: 0 16px 6px;
  padding: 13px 14px;
  border: 1px solid var(--line);
  border-radius: 14px;
  background: var(--surface-2);
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  text-align: left;
  transition: 0.25s;
  color: var(--ink);

  &:hover { border-color: var(--line-strong); }
  .cv {
    width: 36px;
    height: 36px;
    border-radius: 10px;
    flex: none;
    background: conic-gradient(from 120deg, var(--secondary), var(--magenta), var(--gold), var(--secondary));
    box-shadow: var(--glow-secondary);
  }
  .ct { flex: 1; min-width: 0; }
  b {
    font-family: var(--font-display);
    font-weight: 600;
    font-size: 0.92rem;
    display: block;
    line-height: 1.1;
  }
  small { font-size: 0.68rem; color: var(--ink-faint); }
  .chev { color: var(--ink-faint); width: 16px; height: 16px; }
}

.nav {
  flex: 1;
  overflow-y: auto;
  padding: 14px 14px 8px;

  .grp {
    font-family: var(--font-mono);
    font-size: 0.58rem;
    letter-spacing: 0.2em;
    color: var(--ink-faint);
    text-transform: uppercase;
    padding: 16px 12px 7px;
  }
  .nl {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px 12px;
    border-radius: 11px;
    color: var(--ink-dim);
    text-decoration: none;
    font-size: 0.88rem;
    font-weight: 500;
    transition: 0.2s;
    position: relative;

    :deep(svg) { width: 17px; height: 17px; flex: none; }
    .count {
      margin-left: auto;
      font-family: var(--font-mono);
      font-size: 0.64rem;
      color: var(--ink-faint);
    }
    &:hover {
      background: var(--surface-2);
      color: var(--ink);
    }
    &.on {
      background: rgba(70, 232, 255, 0.1);
      color: var(--primary);

      &::before {
        content: '';
        position: absolute;
        left: -1px;
        top: 8px;
        bottom: 8px;
        width: 3px;
        border-radius: 3px;
        background: var(--grad-arcane);
        box-shadow: var(--glow-primary);
      }
    }
    &.sm { font-size: 0.82rem; padding: 8px 12px; }
  }
  .tree {
    margin-left: 18px;
    padding-left: 13px;
    border-left: 1px solid var(--line);
  }
}

.sidefoot {
  padding: 14px 16px;
  border-top: 1px solid var(--line);

  .sync {
    display: flex;
    align-items: center;
    gap: 11px;
    padding: 11px 13px;
    border: 1px solid rgba(55, 232, 164, 0.3);
    border-radius: 12px;
    background: rgba(55, 232, 164, 0.07);
    margin-bottom: 12px;

    .ld {
      width: 9px;
      height: 9px;
      border-radius: 50%;
      flex: none;
      background: var(--emerald);
      box-shadow: 0 0 10px var(--emerald);
      animation: blink 2s ease-in-out infinite;
    }
    b { font-size: 0.78rem; color: var(--emerald); display: block; line-height: 1.2; }
    small { font-size: 0.66rem; color: var(--ink-faint); font-family: var(--font-mono); }
  }
  .user {
    display: flex;
    align-items: center;
    gap: 11px;

    .av {
      width: 34px;
      height: 34px;
      border-radius: 10px;
      flex: none;
      background: conic-gradient(from 90deg, var(--primary), var(--secondary), var(--primary));
      position: relative;
      display: grid;
      place-items: center;
      font-family: var(--font-display);
      font-weight: 600;
      font-size: 0.9rem;

      &::before {
        content: '';
        position: absolute;
        inset: 2px;
        border-radius: 8px;
        background: var(--surface-solid);
      }
      span { position: relative; z-index: 1; }
    }
    b { font-size: 0.82rem; display: block; line-height: 1.1; }
    small { font-size: 0.66rem; color: var(--ink-faint); }
    .ui { min-width: 0; }
    .ui b { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; max-width: 150px; }
    .gear {
      margin-left: auto;
      flex: none;
      display: grid;
      place-items: center;
      width: 30px;
      height: 30px;
      border-radius: 9px;
      border: 1px solid var(--line);
      background: var(--surface);
      color: var(--ink-faint);
      cursor: pointer;
      transition: 0.2s;
    }
    .gear :deep(svg) { width: 16px; height: 16px; }
    .gear:hover { color: var(--ember); border-color: var(--ember); }
  }
}

@keyframes blink {
  50% { opacity: 0.4; }
}
</style>
