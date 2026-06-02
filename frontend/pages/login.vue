<script setup lang="ts">
const { t } = useI18n()
const auth = useAuthStore()
const campaign = useCampaignStore()

const mode = ref<'login' | 'register'>('login')
const email = ref('')
const password = ref('')
const name = ref('')
const loading = ref(false)
const error = ref('')

useSeoMeta({ title: () => t('login.title') })

async function submit() {
  error.value = ''
  loading.value = true
  try {
    if (mode.value === 'register') {
      await auth.register(email.value, password.value, name.value)
    } else {
      await auth.login(email.value, password.value)
    }
    await campaign.ensure()
    await navigateTo('/dashboard')
  } catch (e: any) {
    error.value = e?.data?.error || t('login.error')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <main class="auth">
    <div class="card">
      <NuxtLink to="/" class="mark">
        <span class="sig" aria-hidden="true">
          <svg viewBox="0 0 24 24" fill="none" stroke="#fff" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 2l2.4 6.6L21 9l-5 4.2L17.5 21 12 17l-5.5 4L8 13.2 3 9l6.6-.4z" />
          </svg>
        </span>
        Aetherwright
      </NuxtLink>

      <h1>{{ mode === 'login' ? t('login.welcome') : t('login.create') }}</h1>
      <p class="sub">{{ t('login.subtitle') }}</p>

      <form @submit.prevent="submit">
        <label v-if="mode === 'register'" class="field">
          <span>{{ t('login.name') }}</span>
          <input v-model="name" type="text" autocomplete="name" />
        </label>
        <label class="field">
          <span>{{ t('login.email') }}</span>
          <input v-model="email" type="email" autocomplete="email" required />
        </label>
        <label class="field">
          <span>{{ t('login.password') }}</span>
          <input v-model="password" type="password" autocomplete="current-password" required minlength="8" />
        </label>

        <p v-if="error" class="err">{{ error }}</p>

        <button class="btn" type="submit" :disabled="loading">
          <Icon v-if="loading" name="lucide:loader-circle" class="spin" />
          {{ loading ? t('login.working') : mode === 'login' ? t('login.submitLogin') : t('login.submitRegister') }}
        </button>
      </form>

      <button class="toggle" @click="mode = mode === 'login' ? 'register' : 'login'">
        {{ mode === 'login' ? t('login.toRegister') : t('login.toLogin') }}
      </button>
    </div>
  </main>
</template>

<style lang="scss" scoped>
.auth {
  position: relative;
  z-index: 2;
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 40px 20px;
}

.card {
  width: min(420px, 100%);
  background: var(--surface-2);
  border: 1px solid var(--line-strong);
  border-radius: 24px;
  backdrop-filter: blur(18px);
  box-shadow: var(--shadow-panel);
  padding: 38px 34px;
  position: relative;
  overflow: hidden;
  animation: rise 0.7s cubic-bezier(0.2, 0.8, 0.2, 1) both;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 1px;
    background: linear-gradient(90deg, transparent, var(--line-strong), transparent);
  }
}

.mark {
  display: inline-flex;
  align-items: center;
  gap: 11px;
  font-family: var(--font-deco);
  font-weight: 900;
  letter-spacing: 0.06em;
  font-size: 1.05rem;
  color: var(--ink);
  text-decoration: none;
  margin-bottom: 24px;

  .sig {
    width: 34px;
    height: 34px;
    display: grid;
    place-items: center;
    border-radius: 10px;
    background: var(--grad-arcane);
    box-shadow: var(--glow-primary);
    position: relative;

    &::before {
      content: '';
      position: absolute;
      inset: 1.5px;
      border-radius: 8px;
      background: var(--void-2);
    }
    svg { position: relative; z-index: 1; width: 18px; height: 18px; }
  }
}

h1 {
  font-family: var(--font-display);
  font-weight: 600;
  font-size: 1.7rem;
  letter-spacing: 0.02em;
  text-shadow: var(--glow-text);
}
.sub { color: var(--ink-faint); font-size: 0.86rem; margin-top: 6px; margin-bottom: 24px; }

.field {
  display: flex;
  flex-direction: column;
  gap: 7px;
  margin-bottom: 16px;

  span {
    font-family: var(--font-mono);
    font-size: 0.64rem;
    letter-spacing: 0.16em;
    text-transform: uppercase;
    color: var(--ink-faint);
  }
  input {
    font-family: var(--font-body);
    font-size: 0.95rem;
    color: var(--ink);
    background: var(--surface);
    border: 1px solid var(--line-strong);
    border-radius: 12px;
    padding: 12px 15px;
    transition: 0.25s;

    &:focus {
      outline: 0;
      border-color: var(--primary);
      box-shadow: 0 0 0 3px rgba(70, 232, 255, 0.14), var(--glow-primary);
    }
  }
}

.err {
  color: var(--ember);
  font-size: 0.8rem;
  margin: 4px 0 14px;
  font-family: var(--font-mono);
}

.btn {
  width: 100%;
  font-family: var(--font-body);
  font-weight: 600;
  font-size: 0.92rem;
  border: 0;
  cursor: pointer;
  padding: 13px;
  border-radius: 13px;
  margin-top: 8px;
  background: var(--grad-arcane);
  color: #06040c;
  box-shadow: var(--glow-primary);
  transition: transform 0.2s, box-shadow 0.25s;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 9px;

  &:hover:not(:disabled) { transform: translateY(-2px); box-shadow: 0 0 36px rgba(70, 232, 255, 0.6); }
  &:disabled { opacity: 0.6; cursor: progress; }
  :deep(svg) { width: 16px; height: 16px; }
}
:global(html[data-theme='light']) .btn { color: #fff; }

.spin { animation: spin 0.8s linear infinite; }

.toggle {
  display: block;
  margin: 20px auto 0;
  background: none;
  border: 0;
  color: var(--ink-dim);
  font-size: 0.82rem;
  cursor: pointer;
  transition: 0.2s;

  &:hover { color: var(--primary); }
}

@keyframes rise {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: none; }
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
