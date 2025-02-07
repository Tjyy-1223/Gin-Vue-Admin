<template>
    <NIcon size="18" class="cursor-pointer" @click="toggle">
        <p v-if="isFullscreen" class="i-ant-design:fullscreen-exit-outlined" />
        <p v-else class="i-ant-design:fullscreen-outlined" />
    </NIcon>
</template>


<script setup>
import { NIcon } from 'naive-ui'
import { useFullscreen } from '@vueuse/core'

const { isFullscreen, toggle } = useFullscreen()
</script>

<style lang="scss" scoped></style>