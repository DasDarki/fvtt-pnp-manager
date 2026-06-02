<template>
  <div class="aw-shell">
    <AtmosBackdrop />
    <slot />
  </div>
</template>

<style lang="scss" scoped>
.aw-shell {
  position: relative;
  min-height: 100vh;
  overflow-x: hidden;
}
</style>
