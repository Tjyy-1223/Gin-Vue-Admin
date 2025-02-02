<template>
    <div
        class="group h-[430px] w-full flex flex-col animate-zoom-in animate-duration-700 items-center rounded-xl bg-white shadow-md transition-600 md:h-[280px] md:flex-row hover:shadow-2xl">

        <!-- 封面图 -->
        <div :class="isRightClass" class="h-[230px] w-full overflow-hidden md:h-full md:w-45/100">
            <RouterLink :to="`/article/${article.id}`">
                <img class="h-full w-full transition-600 hover:scale-110" :src="convertImgUrl(article.img)">
            </RouterLink>
        </div>


        <!-- 文章信息 -->
        <div class="my-4 w-9/10 md:w-55/100 space-y-4 md:px-10">
            <RouterLink :to="`/article/${article.id}`">
                <span class="text-2xl font-bold transition-300 group-hover:text-violet">
                    {{ article.title }}
                </span>
            </RouterLink>
            <div class="flex flex-wrap text-sm color-[#858585]">
                <!-- 置顶 -->
                <span v-if="article.is_top === 1" class="flex items-center color-[#ff7242]">
                    <Icon icon="carbon:align-vertical-top" style="font-size: 20px;" /> 置顶
                </span>
                <span v-if="article.is_top === 1" class="mx-1.5">|</span>
                <!-- 日期 -->
                <span class="flex items-center">
                    <Icon icon="mdi-calendar-month-outline" style="font-size: 20px;" />{{
                        dayjs(article.created_at).format('YYYY-MM-DD')
                    }}
                </span>
                <span class="mx-1.5">|</span>
                <!-- 分类 -->
                <RouterLink :to="`/categories/${article.category_id}?name=${article.category?.name}`"
                    class="flex items-center">
                    <Icon icon="mdi-inbox-full" style="font-size: 20px;" /> {{ article.category?.name }}
                </RouterLink>
                <span class="mx-1.5">|</span>
                <!-- 标签 -->
                <div class="flex gap-1">
                    <RouterLink v-for="tag in article.tags" :key="tag.id" :to="`/tags/${tag.id}?name=${tag.name}`"
                        class="flex items-center">
                        <Icon icon="mdi-tag-multiple" style="font-size: 20px;" /> {{ tag.name }}
                    </RouterLink>
                </div>
            </div>
            <div class="ell-4 text-sm leading-6">
                {{ article.content }}
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { convertImgUrl } from '@/utils'
import dayjs from 'dayjs'
import { Icon } from '@iconify/vue';
const props = defineProps({
    idx: Number,
    article: {},
})

// 判断图片放置位置 (左 or 右)
const isRightClass = computed(() => props.idx % 2 === 0
    ? 'rounded-t-xl md:order-0 md:rounded-l-xl md:rounded-tr-0'
    : 'rounded-t-xl md:order-1 md:rounded-r-xl md:rounded-tl-0')
</script>

<style scoped lang="scss">
.ell-4 {
    display: -webkit-box;
    overflow: hidden;
    text-overflow: ellipsis;
    -webkit-line-clamp: 4;
    -webkit-box-orient: vertical;
}
</style>