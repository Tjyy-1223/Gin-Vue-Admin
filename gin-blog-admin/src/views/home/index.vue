<template>
    <!-- 外层容器是一个 AppPage 组件，包裹整个页面 -->
    <AppPage>
        <div class="flex-1">

            <!-- 用户信息卡片 -->
            <NCard>
                <div class="flex items-center">
                    <!-- 用户头像，圆形，大小为 60px -->
                    <NAvatar round :size="60" :src="avatar" />
                    <div class="ml-5">
                        <!-- 显示用户昵称 -->
                        <p> Hello, {{ nickname }} </p>
                        <!-- 使用线性渐变的文本显示一句话 -->
                        <NGradientText class="mt-1 op-60"
                            gradient="linear-gradient(90deg, red 0%, green 50%, blue 100%)">
                            {{ sentence }}
                        </NGradientText>
                    </div>
                    <div class="ml-auto flex items-center">
                        <!-- 显示 GitHub 项目的 Stars 数量 -->
                        <NStatistic label="Stars" class="w-[80px]">
                            <a href="https://github.com/szluyu99/gin-vue-blog" target="_blank">
                                <img alt="stars" src="https://badgen.net/github/stars/szluyu99/gin-vue-blog">
                            </a>
                        </NStatistic>
                        <!-- 显示 GitHub 项目的 Forks 数量 -->
                        <NStatistic label="Forks" class="ml-10 w-[100px]">
                            <a href="https://github.com/szluyu99/gin-vue-blog" target="_blank">
                                <img alt="forks" src="https://badgen.net/github/forks/szluyu99/gin-vue-blog">
                            </a>
                        </NStatistic>
                    </div>
                </div>
            </NCard>

            <!-- 首页统计信息展示，使用 NGrid 布局 -->
            <NGrid class="mt-4" x-gap="12" :cols="4">
                <!-- v-for 遍历数组渲染每个统计项 -->
                <template v-for="item of [
                    { icon: 'i-fa6-solid:users', color: 'text-[#40C9C6]', label: '访问量', key: 'view_count' },
                    { icon: 'i-heroicons:users-solid', color: 'text-[#34BFA3]', label: '用户量', key: 'user_count' },
                    { icon: 'i-material-symbols:article', color: 'text-[#F4516C]', label: '文章量', key: 'article_count' },
                    { icon: 'i-bxs:comment-dots', color: 'text-[#36A3F7]', label: '留言量', key: 'message_count' },
                ]" :key="item.key">
                    <!-- 每个统计项的布局 -->
                    <NGi>
                        <NCard>
                            <!-- 图标部分，大小为 60px，并根据数据动态设置颜色和图标 -->
                            <span class="text-[60px]" :class="[item.icon, item.color]" />
                            <!-- 显示统计数据 -->
                            <NStatistic class="float-right" :label="item.label">
                                {{ homeInfo[item.key] ?? 'unknown' }} <!-- 如果数据不存在，显示 'unknown' -->
                            </NStatistic>
                        </NCard>
                    </NGi>
                </template>
            </NGrid>

            <!-- 项目展示卡片，待完善首页设计 -->
            <NCard title="项目" size="small" class="mt-4">
                <!-- 卡片的右上角有一个额外的按钮 -->
                <template #header-extra>
                    <NButton text type="primary">
                        更多
                    </NButton>
                </template>
                <!-- v-for 渲染 5 个项目卡片 -->
                <NCard v-for="i in 5" :key="i" class="my-2 w-[300px] flex-shrink-0 cursor-pointer hover:shadow-lg"
                    title="Gin Blog Admin" size="small">
                    <!-- 项目描述 -->
                    <p class="op-60">
                        这是个基于 gin 开发的博客管理后台
                    </p>
                </NCard>
            </NCard>
        </div>
    </AppPage>
</template>


<script setup>
import { onMounted, ref } from 'vue'
import { NAvatar, NButton, NCard, NGi, NGradientText, NGrid, NStatistic } from 'naive-ui'

import AppPage from '@/components/common/AddPage.vue'
import { useUserStore } from '@/store'
import api from '@/api'

const { nickname, avatar } = useUserStore()

const homeInfo = ref({
    view_count: 0,
    user_count: 0,
    article_count: 0,
    message_count: 0,
})

onMounted(async () => {
    getOneSentence()
    const res = await api.getHomeInfo()
    homeInfo.value = res.data
})

// 一言
const sentence = ref('')
async function getOneSentence() {
    fetch('https://v1.hitokoto.cn?c=i')
        .then(resp => resp.json())
        .then(data => sentence.value = data.hitokoto)
        .catch(() => sentence.value = '宠辱不惊，看庭前花开花落；去留无意，望天上云卷云舒。')
}
</script>

<style lang="scss" scoped></style>