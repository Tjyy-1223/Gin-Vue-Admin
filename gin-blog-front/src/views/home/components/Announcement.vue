<template>
    <div class="card-view hidden animate-zoom-in animate-duration-600 lg:block space-y-2">
        <p class="flex items-center text-lg">
            <Icon icon="fluent-emoji-flat:bell" style="font-size: 20px;" />
            <span> 公告 </span>
        </p>
        <div class="leading-6">
            {{ blogConfig.website_intro }}
        </div>
    </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useAppStore } from '@/store'
import { Icon } from '@iconify/vue';

const { blogConfig } = storeToRefs(useAppStore())
</script>

<style scoped></style>