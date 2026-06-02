<script setup lang="ts">
const { t, locale, setLocale } = useI18n()
const { theme, toggle } = useTheme()
const appConfig = useAppConfig()

function switchLocale() {
  setLocale(locale.value === 'de' ? 'en' : 'de')
}
</script>

<template>
  <main class="hero">
    <div class="toolbar">
      <button class="ctl" :aria-label="t('a11y.toggleLang')" @click="switchLocale">
        {{ locale === 'de' ? 'EN' : 'DE' }}
      </button>
      <button class="ctl" :aria-label="t('a11y.toggleTheme')" @click="toggle">
        {{ theme === 'dark' ? '☾' : '☀' }}
      </button>
    </div>

    <span class="eyebrow">{{ t('home.eyebrow') }}</span>

    <h1 class="wordmark">
      <span class="sig" aria-hidden="true">
        <svg viewBox="0 0 24 24" fill="none" stroke="#fff" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round">
          <path d="M12 2l2.4 6.6L21 9l-5 4.2L17.5 21 12 17l-5.5 4L8 13.2 3 9l6.6-.4z" />
        </svg>
      </span>
      {{ appConfig.site.name }}
    </h1>

    <p class="sub">{{ t('brand.deck') }}</p>
    <p class="lead">{{ t('home.lead') }}</p>

    <div class="actions">
      <NuxtLink class="btn btn-primary" to="/dashboard">{{ t('home.enter') }}</NuxtLink>
      <a class="btn btn-ghost" href="https://nuxt.com" target="_blank" rel="noopener">{{ t('home.docs') }}</a>
    </div>

    <div class="chips">
      <span class="chip">Nuxt 3</span>
      <span class="chip">Tailwind</span>
      <span class="chip">vue-i18n</span>
      <span class="chip">FoundryVTT v13</span>
    </div>
  </main>
</template>

<style lang="scss" scoped>
.hero {
  position: relative;
  z-index: 2;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 80px 24px;
}

.toolbar {
  position: fixed;
  top: 22px;
  right: 24px;
  display: flex;
  gap: 10px;
  z-index: 40;

  .ctl {
    width: 42px;
    height: 42px;
    display: grid;
    place-items: center;
    border-radius: 12px;
    border: 1px solid var(--line-strong);
    background: var(--surface);
    color: var(--ink);
    font-family: var(--font-mono);
    font-size: 0.82rem;
    cursor: pointer;
    backdrop-filter: blur(12px);
    transition: 0.25s;

    &:hover {
      color: var(--primary);
      border-color: var(--primary);
      box-shadow: var(--glow-primary);
      transform: translateY(-1px);
    }
  }
}

.eyebrow {
  font-family: var(--font-mono);
  font-size: 0.7rem;
  letter-spacing: 0.42em;
  text-transform: uppercase;
  color: var(--primary);
  display: inline-flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 26px;
  animation: rise 0.8s cubic-bezier(0.2, 0.8, 0.2, 1) both;

  &::before,
  &::after {
    content: '';
    width: 44px;
    height: 1px;
    background: linear-gradient(90deg, transparent, var(--primary));
  }
  &::after {
    transform: scaleX(-1);
  }
}

.wordmark {
  display: flex;
  align-items: center;
  gap: 22px;
  font-family: var(--font-deco);
  font-weight: 900;
  letter-spacing: 0.02em;
  line-height: 0.96;
  font-size: clamp(2.6rem, 9vw, 6.6rem);
  background: var(--grad-arcane);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  filter: drop-shadow(0 0 40px rgba(70, 232, 255, 0.22));
  animation: rise 0.8s cubic-bezier(0.2, 0.8, 0.2, 1) 0.08s both;

  .sig {
    width: clamp(48px, 9vw, 82px);
    height: clamp(48px, 9vw, 82px);
    display: grid;
    place-items: center;
    border-radius: 18px;
    background: var(--grad-arcane);
    box-shadow: var(--glow-primary);
    position: relative;
    flex: none;

    &::before {
      content: '';
      position: absolute;
      inset: 3px;
      border-radius: 15px;
      background: var(--void-2);
    }
    svg {
      position: relative;
      z-index: 1;
      width: 50%;
      height: 50%;
    }
  }
}

.sub {
  font-family: var(--font-display);
  font-weight: 500;
  letter-spacing: 0.42em;
  text-transform: uppercase;
  font-size: clamp(0.78rem, 2.2vw, 1.1rem);
  color: var(--ink-dim);
  margin-top: 18px;
  animation: rise 0.8s cubic-bezier(0.2, 0.8, 0.2, 1) 0.16s both;
}

.lead {
  max-width: 560px;
  margin-top: 28px;
  color: var(--ink-dim);
  font-size: 1.02rem;
  font-weight: 300;
  animation: rise 0.8s cubic-bezier(0.2, 0.8, 0.2, 1) 0.24s both;
}

.actions {
  display: flex;
  gap: 14px;
  margin-top: 36px;
  flex-wrap: wrap;
  justify-content: center;
  animation: rise 0.8s cubic-bezier(0.2, 0.8, 0.2, 1) 0.32s both;
}

.btn {
  font-family: var(--font-body);
  font-weight: 600;
  font-size: 0.9rem;
  border: 0;
  cursor: pointer;
  padding: 13px 26px;
  border-radius: 13px;
  text-decoration: none;
  transition: transform 0.2s, box-shadow 0.25s;
}
.btn-primary {
  background: var(--grad-arcane);
  color: #06040c;
  box-shadow: var(--glow-primary);
}
.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 0 36px rgba(70, 232, 255, 0.6), 0 10px 30px -8px rgba(183, 104, 255, 0.5);
}
.btn-ghost {
  background: var(--surface);
  color: var(--ink);
  border: 1px solid var(--line-strong);
}
.btn-ghost:hover {
  transform: translateY(-2px);
  color: var(--primary);
  border-color: var(--primary);
  box-shadow: var(--glow-primary);
}

.chips {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  justify-content: center;
  margin-top: 40px;
  animation: rise 0.8s cubic-bezier(0.2, 0.8, 0.2, 1) 0.4s both;

  .chip {
    font-family: var(--font-mono);
    font-size: 0.66rem;
    letter-spacing: 0.08em;
    color: var(--ink-dim);
    padding: 6px 13px;
    border-radius: 8px;
    border: 1px solid var(--line);
    background: var(--surface);
  }
}

@keyframes rise {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: none;
  }
}
</style>
