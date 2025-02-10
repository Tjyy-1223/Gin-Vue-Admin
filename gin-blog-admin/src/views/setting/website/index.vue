<template>
    <CommonPage :show-header="false" show-footer>
        <NTabs type="line" animated>
            <NTabPane name="website" tab="网站信息">
                <NForm ref="formRef" label-placement="left" label-align="left" :label-width="120" :model="form"
                    class="mt-4 w-[500px]">
                    <NFormItem label="网站头像" path="website_avatar">
                        <UploadOne v-model:preview="form.website_avatar" :width="120" />
                    </NFormItem>
                    <NFormItem label="网站名称" path="website_name">
                        <NInput v-model:value="form.website_name" placeholder="请输入网站名称" />
                    </NFormItem>
                    <NFormItem label="网站作者" path="website_author">
                        <NInput v-model:value="form.website_author" placeholder="请输入网站作者" />
                    </NFormItem>
                    <NFormItem label="网站简介" path="website_intro">
                        <NInput v-model:value="form.website_intro" placeholder="请输入网站简介" />
                    </NFormItem>
                    <NFormItem label="网站创建日期" path="website_createtime">
                        <NDatePicker v-model:formatted-value="form.website_createtime"
                            value-format="yyyy-MM-dd HH:mm:ss" type="datetime" />
                    </NFormItem>
                    <NFormItem label="网站公告" path="website_notice">
                        <NInput v-model:value="form.website_notice" type="textarea" placeholder="请输入网站公告"
                            :autosize="{ minRows: 4, maxRows: 6 }" />
                    </NFormItem>
                    <NFormItem label="网站备案号" path="website_record">
                        <NInput v-model:value="form.website_record" placeholder="请输入网站备案号" />
                    </NFormItem>
                    <!-- TODO: 第三方登录 -->
                    <!-- <n-form-item label="第三方登录" path="social_login_list">
              <n-checkbox-group v-model:value="cities">
                <n-space item-style="display: flex;">
                  <n-checkbox value="QQ" label="QQ" />
                  <n-checkbox value="WeiBo" label="微博" />
                  <n-checkbox value="WeChat" label="微信" />
                </n-space>
              </n-checkbox-group>
            </n-form-item> -->
                    <NButton type="primary" @click="handleSave">
                        确认
                    </NButton>
                </NForm>
            </NTabPane>
            <NTabPane name="contact" tab="社交信息">
                <NForm ref="formRef" label-placement="left" label-align="left" :label-width="120" :model="form"
                    class="mt-4 w-[500px]">
                    <NFormItem label="QQ" path="qq">
                        <NInput v-model:value="form.qq" placeholder="请输入 QQ" />
                    </NFormItem>
                    <NFormItem label="Github" path="github">
                        <NInput v-model:value="form.github" placeholder="请输入 Github" />
                    </NFormItem>
                    <NFormItem label="Gitee" path="gitee">
                        <NInput v-model:value="form.gitee" placeholder="请输入 Gitee" />
                    </NFormItem>
                    <NButton type="primary" @click="handleSave">
                        确认
                    </NButton>
                </NForm>
            </NTabPane>
            <NTabPane name="other" tab="其他设置">
                <NForm ref="formRef" label-placement="left" label-align="left" :label-width="120" :model="form"
                    class="mt-4">
                    <NForm ref="formRef" label-align="left" :label-width="120" :model="form" inline>
                        <NFormItem label="用户头像" path="user_avatar">
                            <UploadOne v-model:preview="form.user_avatar" :width="120" />
                        </NFormItem>
                        <NFormItem label="游客头像" path="tourist_avatar">
                            <UploadOne v-model:preview="form.tourist_avatar" :width="120" />
                        </NFormItem>
                        <!-- <n-form-item label="微信收款码" path="tourist_avatar">
                <n-image border-dashed border-1 text-gray width="120" :src="form.tourist_avatar" />
              </n-form-item>
              <n-form-item label="支付宝收款码" path="tourist_avatar">
                <n-image border-dashed border-1 text-gray width="120" :src="form.tourist_avatar" />
              </n-form-item> -->
                    </NForm>
                    <NFormItem label-placement="top" label="文章默认封面" path="article_cover">
                        <UploadOne v-model:preview="form.article_cover" :width="300" />
                    </NFormItem>
                    <NFormItem label="评论默认审核" path="is_comment_review">
                        <NRadioGroup v-model:value="form.is_comment_review" name="is_comment_review">
                            <NRadio value="true">
                                关闭
                            </NRadio>
                            <NRadio value="false">
                                开启
                            </NRadio>
                        </NRadioGroup>
                    </NFormItem>
                    <NFormItem label="留言默认审核" path="is_message_review">
                        <NRadioGroup v-model:value="form.is_message_review" name="is_message_review">
                            <NRadio value="true">
                                关闭
                            </NRadio>
                            <NRadio value="false">
                                开启
                            </NRadio>
                        </NRadioGroup>
                    </NFormItem>
                    <!-- <NFormItem label="邮箱通知" path="is_email_notice">
              <NRadioGroup v-model:value="form.is_email_notice" name="is_email_notice">
                <NRadio :value="0">
                  关闭
                </NRadio>
                <NRadio :value="1">
                  开启
                </NRadio>
              </NRadioGroup>
            </NFormItem> -->
                    <NButton type="primary" @click="handleSave">
                        确认
                    </NButton>
                </NForm>
            </NTabPane>
        </NTabs>
    </CommonPage>
</template>


<script setup>
import { onMounted, ref } from 'vue'  // 导入 Vue 相关的 API，创建响应式数据
import { NButton, NDatePicker, NForm, NFormItem, NInput, NRadio, NRadioGroup, NTabPane, NTabs } from 'naive-ui'  // 导入 Naive UI 组件库的相关组件

import CommonPage from '@/components/common/CommonPage.vue'  // 导入一个通用页面布局组件
import UploadOne from '@/components//UploadOne.vue'  // 导入一个文件上传组件（用于上传单个文件）

import api from '@/api'  // 导入 API 模块，用于与后台进行交互

defineOptions({ name: '网站管理' })  // 设置组件的名称为 "网站管理"

// 定义响应式数据
const formRef = ref(null)  // 用于引用表单组件，进行表单验证
const form = ref({  // 存储网站配置的表单数据
    website_avatar: '',  // 网站头像
    website_name: 'Tjyy的个人博客',  // 网站名称
    website_author: 'Tjyy',  // 网站作者
    website_intro: 'coding is coding',  // 网站简介
    website_notice: '博客后端基于 gin、gorm 开发\n博客前端基于 Vue3、TS、NaiveUI 开发\n努力学习中...冲冲冲！加油！',  // 网站公告
    website_createtime: '2023-12-27 22:40:22',  // 网站创建时间
    website_record: '鲁ICP备2022040119号',  // 网站备案号
    qq: '123456789',  // QQ号
    github: 'https://github.com/Tjyy-1223',  // GitHub 链接
    gitee: 'https://github.com/Tjyy-1223',  // Gitee 链接
    tourist_avatar: 'https://cdn.hahacode.cn/16815451239215dc82548dcadcd578a5bbc8d5deaa.jpg',  // 游客头像
    user_avatar: 'https://cdn.hahacode.cn/2299fc4d14c94e6183b082973b35855d.png',  // 用户头像
    article_cover: 'https://cdn.hahacode.cn/1679461519cc592408198d67faf1290ff8969dc614.png',  // 文章封面图
    is_comment_review: 1,  // 是否启用评论审核（1 表示启用）
    is_message_review: 1,  // 是否启用留言审核（1 表示启用）
    // is_email_notice: 0,  // 是否启用邮件通知（注释掉，暂时没有使用）
    // social_login_list: [],  // 社交登录列表（注释掉，暂时没有使用）
    // social_url_list: [],  // 社交 URL 列表（注释掉，暂时没有使用）
    // is_reward: 0,  // 是否启用打赏（注释掉，暂时没有使用）
    // wechat_qrcode: 'http://dummyimage.com/100x100',  // 微信二维码（注释掉，暂时没有使用）
    // alipay_ode: 'http://dummyimage.com/100x100',  // 支付宝二维码（注释掉，暂时没有使用）
})

onMounted(async () => {  // 在组件挂载完成后执行
    fetchData()  // 调用 fetchData 函数获取网站配置信息
})

async function fetchData() {  // 获取网站配置信息
    const resp = await api.getConfig()  // 调用 API 获取配置数据
    form.value = resp.data  // 将返回的数据赋值给 form
}

function handleSave() {  // 保存网站配置
    formRef.value?.validate(async (err) => {  // 对表单进行验证
        if (!err) {  // 如果表单验证没有错误
            try {
                $loadingBar?.start()  // 显示加载条
                await api.updateConfig(form.value)  // 调用 API 更新网站配置
                $loadingBar?.finish()  // 完成加载条
                $message.success('博客信息更新成功')  // 显示成功消息
                // fetchData()  // 更新数据，可以选择刷新数据（此行被注释掉）
            }
            catch (err) {  // 如果保存过程中发生错误
                $loadingBar?.error()  // 显示加载条错误状态
            }
        }
    })
}
</script>

<style lang="scss" scoped></style>