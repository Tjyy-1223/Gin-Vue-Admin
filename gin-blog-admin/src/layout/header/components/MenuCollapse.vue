<template>
    <NIcon size="18" class="cursor-pointer" @click="themeStore.switchCollapsed">
        <p v-if="themeStore.collapsed" class="i-mdi-format-indent-increase" />
        <p v-else class="i-mdi-format-indent-decrease" />
    </NIcon>
</template>


<script setup>
import { NIcon } from 'naive-ui'
import { useThemeStore } from '@/store'

const themeStore = useThemeStore()
</script>

<style lang="scss" scoped></style>