<template>
    <!-- <Teleport> 的 to 属性可以动态改变，允许你在不同目标之间切换 -->
    <!-- to = "body" 代表其中的元素将被渲染到 body 元素中 -->
    <Teleport to="body">
        <div class="fixed inset-0 overflow-y-auto transition-all ease-in" :class="[
            isOpen ? 'visible' : 'invisible duration-100 ease-in',
        ]" :style="{ 'z-index': zIndex }">

            <!-- overlay -->
            <!-- @click.self="close": 这表示当用户点击这个 div 元素本身时触发 close 方法 -->
            <div class="fixed inset-0 bg-black transition-opacity" :class="[
                isOpen ? 'opacity-50 duration-75 ease-out' : 'opacity-0 duration-75 ease-in',
            ]" @click.self="close"></div>

            <!-- dialog -->
            <div class="min-h-full flex items-center justify-center p-3">
                <div v-bind="$attrs"
                    class="relative inline-block w-full rounded-lg bg-white shadow-xl transition-all dark:bg-gray-900"
                    :class="[
                        padded ? 'p-4 lg:py-5 lg:px-7' : 'p-1',
                        isOpen
                            ? 'translate-y-0 opacity-100 duration-300 sm:scale-100'
                            : 'translate-y-4 opacity-0 duration-300 sm:translate-y-0 sm:scale-95',
                    ]" :style="{
                    width: `${width}px`,
                }">

                    <button v-if="dismissButton"
                        class="absolute right-5 top-5 h-6 w-6 rounded-full bg-gray-100 p-1 text-gray-700 hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-500"
                        @click="close">
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                    <slot />

                </div>
            </div>
        </div>
    </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'

// props 是用来接收父组件传递给子组件的数据的。你可以把 props 想象成组件的输入，它定义了子组件可以接受的参数。
// 在这段代码中，props 定义了一些属性（即子组件可以接受的外部输入），它们的作用如下：
// modelValue: 这是一个布尔类型的属性，表示对话框是否打开。它通常与父组件的 v-model 进行双向绑定（v-model 语法糖会将 modelValue 与父组件的变量绑定）。
// dismissible: 一个布尔类型的属性，表示对话框是否可以通过点击外部区域来关闭，默认值为 true。
// dismissButton: 一个布尔类型的属性，表示是否显示关闭按钮，默认值为 true。
// padded: 一个布尔类型的属性，表示对话框是否有内边距，默认值为 true。
// width: 一个数字类型的属性，表示对话框的宽度，默认值为 500（即 500px）。
// zIndex: 一个数字类型的属性，表示对话框的 z-index，用于控制其层级，默认值为 40
const props = defineProps({
    modelValue: {
        type: Boolean,
        required: true,
    },
    dismissible: {
        type: Boolean,
        default: true,
    },
    dismissButton: {
        type: Boolean,
        default: true,
    },
    padded: {
        type: Boolean,
        default: true,
    },
    width: {
        type: Number,
        default: 500,
    },
    zIndex: {
        type: Number,
        default: 40,
    },
})

// emit 用于触发事件，通知父组件发生了某些事情。在 Vue 3 中，子组件通过 emit 触发的事件可以让父组件做出相应的处理。
// update:modelValue: 当对话框的打开状态（isOpen）变化时，子组件通过 emit('update:modelValue', val) 通知父组件更新 v-model 绑定的值。这是实现双向绑定的关键。
// close: 当对话框关闭时，子组件触发 close 事件，父组件可以通过监听该事件来执行一些自定义的行为（比如关闭其他相关的UI，或执行清理操作）。
const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void
    (e: 'close'): void
}>()

// computed 是 Vue 3 中的一个响应式 API，用于创建计算属性。计算属性是基于它们的依赖进行缓存的，并且只有在依赖项发生变化时才会重新计算。计算属性可以用来封装逻辑，避免重复的计算。
// get 是计算属性的获取器，它定义了计算属性的值是如何被获取的。在这里，get 返回了父组件传递的 props.modelValue。
// set 是计算属性的设置器，它定义了计算属性的值是如何被设置的。当你修改计算属性（如 isOpen）的值时，实际上会触发 set 方法，从而执行一些副作用。
// 在这里，set 用来将 isOpen 的变化通知给父组件。它会触发 emit('update:modelValue', val)，通知父组件更新 modelValue，实现双向绑定。
const isOpen = computed({ // isOpen 代表  props.modelValue 的值
    get: () => props.modelValue,
    set: val => emit('update:modelValue', val), // isOpen 变化时调用 emit('update:modelValue', val)
})


function close() {
    if (props.dismissible)
        isOpen.value = false
    emit('close')
}

</script>

<style scoped></style>