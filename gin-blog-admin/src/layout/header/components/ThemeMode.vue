<template>
    <NIcon class="cursor-pointer" size="18" @click="themeStore.switchDarkMode">
        <p v-if="themeStore.darkMode" class="i-mdi-moon-waning-crescent" />
        <p v-else class="i-mdi-white-balance-sunny" />
    </NIcon>
</template>

<script setup>
import { NIcon } from 'naive-ui'
import { useThemeStore } from '@/store'

const themeStore = useThemeStore()
</script>

<style lang="scss" scoped></style>