<template>
  <!-- 顶部中间的消息提示 -->
  <UToast ref="messageRef" position="top" align="center" :timeout="3000" closeable />
  <!-- 右上方的消息通知 -->
  <UToast ref="notifyRef" position="top" align="right" :timeout="3000" closeable />

  <div class="h-full w-full flex flex-col">
    <!-- 顶部导航栏 -->
    <AppHeader></AppHeader>
  </div>

</template>

<script setup lang="ts">
import UToast from './components/ui/UToast.vue';
import AppHeader from './components/layout/AppHeader.vue';
</script>

<style scoped lang="scss">
div {
  height: 300px;
}
</style>