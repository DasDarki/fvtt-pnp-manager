<script setup lang="ts">
const props = defineProps<{
  title: string
  body: string
  subjectLabel: string
  moreCount?: number
  onAck?: () => void | Promise<void>
}>()
const { t } = useI18n()
const dismissed = ref(false)

const moreLabel = computed(() =>
  props.moreCount ? ` · ${t('critical.more', { n: props.moreCount })}` : '',
)

async function ack() {
  dismissed.value = true
  await props.onAck?.()
}
</script>

<template>
  <Transition name="fade">
    <div v-if="!dismissed" class="alert">
      <span class="ai"><Icon name="lucide:triangle-alert" /></span>
      <div class="at">
        <div class="att">
          <span class="badge">{{ t('critical.badge') }}</span>
          <span>{{ title }}</span>
        </div>
        <p>{{ body }}</p>
        <span class="lk">↳ {{ subjectLabel }}{{ moreLabel }}</span>
      </div>
      <button class="ack" @click="ack">
        <Icon name="lucide:check" />
        <span class="hide-s">{{ t('critical.ack') }}</span>
      </button>
    </div>
  </Transition>
</template>

<style lang="scss" scoped>
.alert {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 20px;
  border-radius: 22px;
  margin-bottom: 24px;
  border: 1px solid rgba(255, 106, 85, 0.42);
  background: linear-gradient(100deg, rgba(255, 106, 85, 0.16), rgba(255, 79, 163, 0.05));
  box-shadow: 0 0 34px -8px rgba(255, 106, 85, 0.5), var(--shadow-panel);
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: 4px;
    background: var(--grad-ember);
    box-shadow: 0 0 16px var(--ember);
  }

  .ai {
    width: 46px;
    height: 46px;
    border-radius: 13px;
    display: grid;
    place-items: center;
    flex: none;
    background: rgba(255, 106, 85, 0.18);
    color: var(--ember);
    box-shadow: 0 0 20px -4px var(--ember);
    animation: apulse 2.4s ease-in-out infinite;

    :deep(svg) { width: 23px; height: 23px; }
  }
  .at { flex: 1; min-width: 0; }
  .att {
    font-family: var(--font-display);
    font-weight: 600;
    font-size: 1.02rem;
    display: flex;
    align-items: center;
    gap: 11px;
    flex-wrap: wrap;
  }
  .badge {
    font-family: var(--font-mono);
    font-size: 0.56rem;
    letter-spacing: 0.14em;
    color: var(--ember);
    border: 1px solid var(--ember);
    border-radius: 999px;
    padding: 3px 10px;
    text-transform: uppercase;
  }
  p { font-size: 0.84rem; color: var(--ink-dim); margin-top: 4px; }
  .lk {
    font-family: var(--font-mono);
    font-size: 0.64rem;
    color: var(--ember);
    margin-top: 7px;
    display: inline-block;
  }
  .ack {
    flex: none;
    font-family: var(--font-body);
    font-weight: 600;
    font-size: 0.78rem;
    cursor: pointer;
    padding: 10px 15px;
    border-radius: 10px;
    border: 1px solid rgba(255, 106, 85, 0.4);
    background: rgba(255, 106, 85, 0.1);
    color: var(--ember);
    transition: 0.25s;
    display: inline-flex;
    align-items: center;
    gap: 7px;

    :deep(svg) { width: 14px; height: 14px; }
    &:hover { background: var(--ember); color: #1a0a06; box-shadow: 0 0 22px rgba(255, 106, 85, 0.5); }
  }
}

@keyframes apulse {
  50% { box-shadow: 0 0 32px 0 rgba(255, 106, 85, 0.6); }
}

.fade-leave-active { transition: opacity 0.4s, transform 0.4s; }
.fade-leave-to { opacity: 0; transform: translateY(-8px); }

@media (max-width: 560px) {
  .alert .ack .hide-s { display: none; }
}
</style>
