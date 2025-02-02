<template>
    <!-- 自定义动画类（需要在 CSS 或 Tailwind 配置中定义），用于给横幅添加渐变或淡入效果。 -->
    <div class="banner-fade-down absolute top-20 bottom-0 left-0 right-0 h-screen text-center text-white z-0" :style="coverStyle">
        <!-- 这个容器被定位为绝对定位，inset-x-0 表示左右两侧的内边距为 0，mt-[43vh] 为元素顶部的外边距设置为 43vh（视口高度的 43%），目的是将内容从屏幕上方略微向下偏移。 -->
        <div class="absolute inset-x-0 mt-[43vh] text-center space-y-3">
            <h1 class="animate-zoom-in text-4xl font-bold lg:text-5xl">
                {{ blogConfig.website_name }}
            </h1>
            <div class="text-lg lg:text-xl">
                {{ typer.output }}
                <span class="animate-ping"> | </span>
            </div>
        </div>


        <!-- 社交信息（移动端专用） -->
        <div class="text-2xl space-x-5">
            <a :href="`http://wpa.qq.com/msgrd?v=3&uin=${blogConfig.qq}&site=qq&menu=yes`" target="_blank">
                <span class="i-ant-design:qq-circle-filled inline-block">
                    <Icon icon="ant-design:qq-circle-filled" style="font-size: 20px;" />
                </span>
            </a>
            <a :href="blogConfig.github" target="_blank">
                <span class="i-mdi:github inline-block">
                    <Icon icon="mdi:github" style="font-size: 20px;" />
                </span>
            </a>
            <a :href="blogConfig.gitee" target="_blank">
                <span class="i-simple-icons:gitee inline-block">
                    <Icon icon="simple-icons:gitee" style="font-size: 20px;" />
                </span>
            </a>
        </div>


        <!-- 向下滚动 -->
        <div class="absolute bottom-0 w-full cursor-pointer" @click="scrollDown">
            <span class="inline-block animate-bounce text-2xl text-white" >
                <Icon icon="ep:arrow-down-bold" style="font-size: 30px;" />
            </span>
        </div>
    </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { useAppStore } from '@/store'
import { computed, onMounted, reactive } from 'vue'
import { storeToRefs } from 'pinia'
import EasyTyper from 'easy-typer-js'
// 从 Pinia store 中获取并解构出两个响应式数据：pageList 和 blogConfig
const { pageList, blogConfig } = storeToRefs(useAppStore())

// 打字机特效配置
// reactive() 是一个用于创建响应式数据的函数，它使得对象或数组成为响应式的，也就是说，当数据发生变化时，视图会自动更新。
const typer = reactive({
    output: '',
    isEnd: false, // 全局控制是否终止
    speed: 300, // 打字速度
    singleBack: false, // 单次的回滚
    sleep: 0, // 完整输出一句话后, 睡眠一定时候后触发回滚事件
    type: 'normal', // rollback, normal
    backSpeed: 80, // 回滚速度
    sentencePause: false, // 运行完毕后, 句子是否暂停显示
})

onMounted(() => {
    getOneSentence()
})

// 随机获取一句名言
function getOneSentence() {
    // 一言 + 打字机特效
    fetch('https://v1.hitokoto.cn/?c=i')
        .then(res => res.json())
        .then(data => new EasyTyper(typer, data.hitokoto, () => { }, () => { }))
        .catch(() => new EasyTyper(typer, '宠辱不惊，看庭前花开花落；去留无意，望天上云卷云舒。', () => { }, () => { }))
}

// 这段代码定义了一个名为 scrollDown 的函数，目的是平滑地滚动页面到一个特定的位置。
// 意味着页面将被滚动到当前视口的高度。也就是说，页面会滚动到下方一个完整的视口高度（例如，屏幕的下半部分）。
function scrollDown() {
    window.scrollTo({
        behavior: 'smooth',
        top: document.documentElement.clientHeight,
    })
}

// 根据后端配置动态获取封面
const coverStyle = computed(() => {
    const page = pageList.value.find(e => e.label === 'home')
    return page
        ? `background: url('${page.cover}') center center / cover no-repeat;`
        : 'background: grey center center / cover no-repeat;'
})

</script>

<style scoped></style>