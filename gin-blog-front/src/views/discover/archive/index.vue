<!-- <template> 包裹中的内容整体会插入在 BannerPage 中的 <slot /> 中 -->
<template>
    <BannerPage title="归档" label="archive" :loading="loading" card>
        <p class="pb-5 text-lg lg:text-2xl">
            目前共计 {{ archiveList.length }} 篇文章，继续加油！
        </p>

        <!-- <template> 本身是不会渲染任何内容的，除非它嵌套在具有渲染功能的元素中，或者是作为条件/循环渲染的容器。  -->
        <!-- <template v-if="true">
            <div>222</div> 
        </template> -->

        <template v-for="(item, idx) of archiveList" :key="item.id">
            <div class="flex items-center gap-2">
                <div class="i-mdi:circle bg-blue text-sm" />
                <span class="text-sm color-#666 lg:text-base">
                    {{ dayjs(item.created_at).format('YYYY-MM-DD') }}
                </span>
                <a class="color-#666 lg:text-lg hover:text-orange" @click="router.push(`/article/${item.id}`)">
                    {{ item.title }}
                </a>
            </div>
            <hr v-if="idx !== archiveList.length - 1" class="my-4 border-1 border-color-#d2ebfd border-dashed">
        </template>

        <!-- TODO: 分页 -->
        <!-- <div class="my-15 mt-20 f-c-c">
            <NPagination v-model:page="current" :page-count="Math.ceil(total / 10)" />
        </div> -->
    </BannerPage>
</template>

<script setup>
import { onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import dayjs from 'dayjs'

import BannerPage from '@/components/BannerPage.vue'
import api from '@/api'


const router = useRouter()

const loading = ref(true)
const total = ref(0)
const archiveList = ref([])

// 监听当前页面变化
const current = ref(1) // 当前页数
watch(current, () => getArchives())

async function getArchives() {
    const resp = await api.getArchives({
        page_num: current.value,
        page_size: 50,
    })
    archiveList.value = resp.data.page_data
    total.value = resp.data.total
    loading.value = false
}

onMounted(() => {
    getArchives()
})

</script>

<style lang="scss" scoped></style>