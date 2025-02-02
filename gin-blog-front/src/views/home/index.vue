<template>
    <!-- 首页封面图 -->
    <HomeBanner />
    <!-- 内容 -->
    <div class="mx-auto mb-8 max-w-[1230px] flex flex-col justify-center px-3" style="margin-top: calc(100vh + 30px)">
        <!-- 这表示容器会被划分为 12 列，每列占据相同的空间 -->
        <div class="grid grid-cols-12 gap-4">
            <!-- 左半部分 -->
            <div class="col-span-0 lg:col-span-9 space-y-5">
                <!-- 说说轮播 -->
                <TalkingCarousel />
                <!-- 文章列表 -->
                <div class="space-y-5">
                    <ArticleCard v-for="(item, idx) in articleList" :key="item.id" :article="item" :idx="idx" />
                </div>

                <!-- 无限加载 -->
                <div class="f-c-c">
                    <InfiniteLoading class="mt-2 lg:mt-5" @infinite="getArticlesInfinite">
                        <!-- TODO: 优化界面 -->
                        <template #spinner>
                            <span class="animate-pulse text-xl">
                                loading...
                            </span>
                        </template>
                        <template #complete>
                            <span class="flex justify-center items-center gap-2 text-gray">
                                没有更多文章啦!
                                <button class="flex items-center" @click="backTop">
                                    点击回到顶部
                                    <Icon icon="mdi:arrow-up-bold-box" style="font-size: 20px;" />
                                </button>
                            </span>
                        </template>
                    </InfiniteLoading>
                </div>
            </div>


            <!-- 右半部分 -->
            <div class="col-span-0 lg:col-span-3">
                <!-- sticky 实现悬浮固定效果 -->
                <div class="sticky top-5 space-y-5">
                    <!-- 博主信息 -->
                    <AuthorInfo />
                    <!-- 公告 -->
                    <Announcement />
                    <!-- 网站资讯 -->
                    <WebsiteInfo />
                </div>
            </div>
        </div>
    </div>
    <!-- 底部 -->
    <AppFooter />
</template>

<script setup lang="ts">
import HomeBanner from './components/HomeBanner.vue'
import TalkingCarousel from './components/TalkingCarousel.vue'
import ArticleCard from "./components/ArticleCard.vue"
import AuthorInfo from './components/AuthorInfo.vue'
import WebsiteInfo from './components/WebsiteInfo.vue'
import Announcement from './components/Announcement.vue'
import AppFooter from '@/components/layout/AppFooter.vue'
import { onMounted, reactive, ref } from 'vue'
import InfiniteLoading from 'v3-infinite-loading'
import { marked } from 'marked'
import api from '@/api'
import { Icon } from '@iconify/vue';


const articleList = ref([])
const loading = ref(false)

// 无限加载文章
const params = reactive({ page_size: 5, page_num: 1 }) // 列表加载参数
async function getArticlesInfinite($state: any) {
    if (!loading.value) {
        try {
            const resp = await api.getArticles(params)
            // 加载完成
            if (!resp.data.length) {
                $state.complete()
                return
            }
            // 非首次加载, 都是往列表中添加数据
            articleList.value.push(...resp.data)
            // 过滤 Markdown 符号
            articleList.value.forEach(e => e.content = filterMdSymbol(e.content))
            params.page_num++
            $state.loaded()
        }
        catch (error) {
            $state.error()
        }
    }
}

onMounted(async () => {
    loading.value = true
    // 首次加载
    const resp = await api.getArticles(params)
    articleList.value = resp.data
    // 过滤 Markdown 符号
    articleList.value.forEach(e => e.content = filterMdSymbol(e.content))
    params.page_num++
    loading.value = false
    //   console.log(articleList.value)
})

// 过滤 Markdown 符号: 先转 Html 再去除 Html 标签
// 这段代码的作用是 过滤 Markdown 格式中的符号，将 Markdown 转换为纯文本，并去除其中的 HTML 标签以及一些其他不需要的字符。
function filterMdSymbol(md:any) {
    return marked(md) // 转 HTML
        .replace(/<\/?[^>]*>/g, '') // 正则去除 Html 标签
        .replace(/[|]*\n/, '')
        .replace(/&npsp;/gi, '')
}

function backTop() {
    window.scrollTo({ behavior: 'smooth', top: 0 })
}



</script>

<style scoped></style>