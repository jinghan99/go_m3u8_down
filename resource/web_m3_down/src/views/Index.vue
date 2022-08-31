<script setup>
import FileView from '@/components/File.vue'
import DownView from '@/components/Down.vue'
</script>


<template>
  <van-tabs v-model:active="active" sticky animated swipeable>
    <van-tab  title="下载列表">
      <DownView />
    </van-tab>
    <van-tab title="文件">
      <FileView />
    </van-tab>
  </van-tabs>
</template>

<script>
import {ref} from 'vue';


export default {
  name: "down",
  setup() {
    const activeName = ref('1');
    return {activeName};
  },
};
</script>

<style scoped>
.DownView{
  height: 100%;
}
.FileView{
  height: 100%;
}
</style>
