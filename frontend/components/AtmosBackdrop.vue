<template>
  <div class="aw-atmos" aria-hidden="true">
    <span class="smoke a" />
    <span class="smoke b" />
    <span class="smoke c" />
  </div>
  <div class="aw-grain" aria-hidden="true" />
</template>

<style lang="scss" scoped>
.aw-atmos {
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  overflow: hidden;

  .smoke {
    position: absolute;
    border-radius: 50%;
    filter: blur(80px);
    opacity: 0.55;
    mix-blend-mode: screen;
  }
  .smoke.a {
    width: 52vw;
    height: 52vw;
    top: -16vw;
    left: 8vw;
    background: radial-gradient(circle, rgba(70, 232, 255, 0.5), transparent 65%);
    animation: drift1 30s ease-in-out infinite;
  }
  .smoke.b {
    width: 48vw;
    height: 48vw;
    top: -6vw;
    right: -10vw;
    background: radial-gradient(circle, rgba(183, 104, 255, 0.5), transparent 65%);
    animation: drift2 36s ease-in-out infinite;
  }
  .smoke.c {
    width: 44vw;
    height: 44vw;
    bottom: -22vw;
    left: 36vw;
    background: radial-gradient(circle, rgba(255, 79, 163, 0.32), transparent 65%);
    animation: drift3 42s ease-in-out infinite;
  }
}

:global(html[data-theme='light']) .aw-atmos .smoke {
  opacity: 0.45;
  mix-blend-mode: multiply;
}

.aw-grain {
  position: fixed;
  inset: 0;
  z-index: 60;
  pointer-events: none;
  opacity: 0.045;
  mix-blend-mode: overlay;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='160' height='160'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.85' numOctaves='2'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)'/%3E%3C/svg%3E");
}

@keyframes drift1 {
  50% { transform: translate(8vw, 6vw) scale(1.12); }
}
@keyframes drift2 {
  50% { transform: translate(-6vw, 8vw) scale(1.08); }
}
@keyframes drift3 {
  50% { transform: translate(6vw, -7vw) scale(1.14); }
}
</style>
