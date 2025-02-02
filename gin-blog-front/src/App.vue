<template>
  <!-- 顶部中间的消息提示 -->
  <UToast ref="messageRef" position="top" align="center" :timeout="3000" closeable />
  <!-- 右上方的消息通知 -->
  <UToast ref="notifyRef" position="top" align="right" :timeout="3000" closeable />

  <div class="h-full w-full flex flex-col">
    <!-- 顶部导航栏 -->
    <AppHeader class="z-10"/>

    <!-- 中间内容(包含底部信息) -->
    <article class="flex flex-1 flex-col">
      <!-- { Component, route } 是解构语法，表示从插槽中提取 Component 和 route 两个属性：
            Component：当前路由匹配的组件。
            route：当前路由的信息，包括 path、name 等属性。 
        -->
      <RouterView v-slot="{ Component, route }">  
        <!-- :key 的主要目的是确保组件在变化时能够正确地重新渲染，而不会出现状态或渲染的错误。 -->
        <component :is="Component" :key="route.path" />
      </RouterView>
    </article>

  </div>

</template>

<script setup lang="ts">
import UToast from './components/ui/UToast.vue';
import AppHeader from './components/layout/AppHeader.vue';
import { useAppStore, useUserStore } from '@/store';
import { onMounted, ref } from 'vue'

const appStore = useAppStore()
const userStore = useUserStore()

const messageRef = ref(null)
const notifyRef = ref(null)

onMounted(() => {
  // appStore.getPageList()
  // appStore.getBlogInfo()
  // userStore.getUserInfo()

  // 挂载全局提示
  window.$message = messageRef.value
  window.$notify = notifyRef.value
})
</script>

<style scoped lang="scss">
div {
  height: 300px;
}
</style>