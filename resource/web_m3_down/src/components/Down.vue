<template>
  <van-collapse v-model="activeName" accordion>
    <van-collapse-item title="新增m3u8" name="add_m3u8_url">
      <van-cell-group>
        <van-field v-model="name" label="名称" placeholder="请输入名称"/>
        <van-field v-model="url" label="url" placeholder="请输入url"/>
        <van-space direction="vertical" fill>
          <van-button type="primary" @click="add_m3u8" block>提交</van-button>
        </van-space>
      </van-cell-group>
    </van-collapse-item>
    <van-collapse-item title="新增m3u8-Html" name="add_m3u8_HTML">
      <van-cell-group>
        <van-field v-model="html_url" label="html_url" placeholder="请输入html"/>
        <van-space direction="vertical" fill>
          <van-button type="primary" @click="add_m3u8" block>提交</van-button>
        </van-space>
      </van-cell-group>
    </van-collapse-item>
  </van-collapse>

  <van-divider>下载列表</van-divider>

  <van-cell-group inset  v-for="i in dataList" >
    <van-cell :title="i.name" :value="i.progress + '%'" />
  </van-cell-group>


</template>

<script >
import {ref} from 'vue';
import api  from '@/utils/request'
import { Toast } from 'vant';
import 'vant/es/toast/style';
export default {
  name: "Down",
  setup() {
    const activeName = ref('add_m3u8_url');
    return {activeName};
  },
  data() {
    return {
      name:"",
      url:"",
      html_url:"",
      dataList: [],
      message: 'test!!',
    }
  },
  created() {
    this.initList()
  },
  methods: {
    add_m3u8() {
      api({
        url:'/start_down',
        method: 'get',
        params:({
          "name": this.name,
          "url": this.url
        })
      }).then(data =>{
        console.log(data.data)
       if (200 == data.status){
         Toast({
           message: data.data.msg,
           position: 'top',
         });
       }
      })
    },
    initList() {
      api({
        url:'/all',
        method: 'get',
      }).then(data =>{
        console.log(data.data)
        Toast({
          message: data.data.msg,
          position: 'top',
        });
        if (200 == data.status && 0 ==data.data.code){
         this.dataList = data.data.data
        }

      })
    }
  },
}
</script>

<style scoped>
.right {
  float: right;

}
</style>
