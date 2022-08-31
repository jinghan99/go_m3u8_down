<style scoped>

</style>

<template>
  <RouterView/>
</template>


<script>
import {RouterLink, RouterView,useRouter} from 'vue-router'
import router from "@/router";
import {ref} from 'vue';
export default {
  setup() {
    const active = ref(0);
    return {active};
  },
};

</script>
