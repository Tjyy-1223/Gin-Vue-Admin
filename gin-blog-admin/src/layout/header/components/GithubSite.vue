<template>
    <NIcon class="cursor-pointer" size="18" @click="handleLinkClick">
        <p class="i-mdi:github" />
    </NIcon>
</template>


<script setup>
import { NIcon } from 'naive-ui'

function handleLinkClick() {
    window.open('https://github.com/Tjyy-1223')
}
</script>

<style lang="scss" scoped></style>