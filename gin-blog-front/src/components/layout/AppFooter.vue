<template>
    <footer class="f-c-c">
        <div class="footer-wrap w-full px-5 py-10 text-center text-white leading-8">
            <p> &copy;{{ 2022 }} - {{ new Date().getFullYear() }} By Tjyy</p>
            <a href="https://beian.miit.gov.cn/" target="_blank">
                {{ blogConfig.website_record }}
            </a>
        </div>
    </footer>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useAppStore } from '@/store'

const { blogConfig } = storeToRefs(useAppStore())
</script>

<style scoped>
.footer-wrap {
    background: linear-gradient(-45deg, #ee7752, #ce3e75, #23a6d5, #23d5ab);
    background-size: 400% 400%;
    animation: Gradient 8s ease infinite;
}

@keyframes Gradient {
    0% {
        background-position: 0 50%;
    }

    50% {
        background-position: 100% 50%;
    }

    100% {
        background-position: 0 50%;
    }
}
</style>