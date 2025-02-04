<template>
    <Transition name="slide-fade" appear>
        <div class="card-view space-y-2">
            <div class="flex items-center">
                <span class="i-fa-solid:list-ul" />
                <span class="ml-2">目录</span>
            </div>
            <ul>
                <li v-for="anchor of anchors" :key="anchor.id">
                    <div class="cursor-pointer border-l-4 border-transparent rounded py-1 text-sm color-#666261 hover:bg-#00c4b6 hover:bg-opacity-30"
                        :class="anchor.id === selectAnchor && 'bg-#00c4b6 text-white border-l-#009d92'"
                        :style="{ paddingLeft: `${5 + anchor.indent * 15}px` }" @click="handleClickAnchor(anchor.id)">
                        {{ anchor.name }}
                    </div>
                </li>
            </ul>
        </div>
    </Transition>
</template>


<script setup>
import { onMounted, ref } from 'vue' // 引入 Vue 的生命周期钩子和响应式引用
import { useWindowScroll, watchThrottled } from '@vueuse/core' // 引入窗口滚动钩子和节流函数

// 接收父组件传递的 previewRef 属性，这个属性是一个 DOM 元素引用
const { previewRef } = defineProps({
    previewRef: { type: Object, required: true, },
})

onMounted(() => {
    buildAnchors() // 在组件挂载完成后调用 buildAnchors 方法来生成目录锚点
})

// 定义响应式变量
const selectAnchor = ref('') // 当前选中的目录项的 id
const anchors = ref([]) // 存储所有的锚点信息
const headings = Array.from(previewRef.querySelectorAll('h1,h2,h3,h4,h5,h6')) // 获取所有标题元素（h1-h6）

// 生成目录的锚点
function buildAnchors() {
    // 筛选出文本非空的标题
    const titleList = Array.from(headings).filter(t => !!t.innerText.trim())
    // 获取所有标题的标签名（h1、h2、h3 等），并去重排序
    const hTags = Array.from(new Set(titleList.map(t => t.tagName))).sort()

    let count = 0 // 用于生成唯一的 ID，避免同名标题冲突
    for (let i = 0; i < headings.length; i++) {
        const anchor = headings[i].textContent.trim() // 获取标题文本内容
        // 给每个标题生成一个唯一的 ID，防止重名，在 ID 后加上序号
        headings[i].id = `${anchor}-${count++}`
        // 将每个标题的锚点信息（ID、名称、层级）保存到 anchors 数组中
        anchors.value.push({
            id: headings[i].id,
            name: headings[i].innerText,
            indent: hTags.indexOf(headings[i].tagName), // 根据标签类型确定标题的层级
        })
    }
}

// 处理点击目录项时的平滑滚动
function handleClickAnchor(id) {
    const anchorElement = document.getElementById(id) // 获取目标标题的 DOM 元素
    window.scrollTo({
        behavior: 'smooth', // 启用平滑滚动
        top: anchorElement.offsetTop - 40, // 滚动到目标标题的位置，并留出 40px 的偏移量
    })
    // 延迟设置选中的锚点 ID，以便高亮显示
    setTimeout(() => selectAnchor.value = id, 600)
}

// 实现目录高亮当前滚动位置对应的标题
// 思路是通过滚动条的位置，循环检测每个标题距离顶部的距离，判断当前应该高亮哪个标题
const { y } = useWindowScroll() // 获取当前的滚动位置
watchThrottled(y, () => {
    anchors.value.forEach((e) => {
        const value = headings.find(ee => ee.id === e.id) // 查找当前锚点对应的标题
        if (value && y.value >= value.offsetTop - 50) { // 如果滚动位置大于标题的顶部偏移量，并且小于标题的底部
            selectAnchor.value = value.id // 设置当前选中的锚点 ID
        }
    })
}, { throttle: 200 }) // 设定滚动事件的节流时间，避免高频触发
</script>


<style lang="scss" scoped></style>