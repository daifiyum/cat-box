<template>
  <v-row class="mt-3" ref="el">
    <v-col v-for="item in subData" :key="item.id" cols="12" md="6" lg="6">
      <v-card
        elevation="2"
        @click="swOne(item.id)"
        :title="item.name"
        :subtitle="item.link"
      >
        <template v-slot:append v-if="item.active">
          <v-icon color="success" icon="mdi-check-circle-outline"></v-icon>
        </template>
        <v-divider></v-divider>
        <v-card-actions>
          <div class="d-flex align-center text-medium-emphasis">
            <v-icon icon="mdi-dots-grid" class="handle cursor-move me-2"></v-icon>
            <v-icon icon="mdi-update"></v-icon>
            <span>{{ formatTime(item.updated_at) }}更新</span>
          </div>
          <v-spacer></v-spacer>
          <EditSub :item="item" />
        </v-card-actions>
      </v-card>
    </v-col>
  </v-row>
</template>

<script setup>
import EditSub from "./EditSub.vue";
import { ref, onBeforeMount, onUnmounted, computed } from "vue";
import { get_sub, sw_sub, order_sub } from "@/api/home";
import { useDraggable } from "vue-draggable-plus";
import emitter from "@/utils/emitter";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import "dayjs/locale/zh-cn";
import { useSubStore } from "@/stores/sub";
import { storeToRefs } from "pinia";

const subStore = useSubStore();
const { subNum } = storeToRefs(subStore);

dayjs.extend(relativeTime);
dayjs.locale("zh-cn");

const el = ref()
const subData = ref([]);

function formatTime(isoString) {
  return dayjs(isoString).fromNow().replace(/\s+/g, "");
}

async function fetchData() {
  let res = await get_sub();
  subData.value = res.data;
  subNum.value = subData.value.length
}
onBeforeMount(async () => {
  fetchData();  
});

emitter.on("reloadData", () => {
  fetchData();
});

onUnmounted(() => {
  emitter.off("reloadData");
});

async function swOne(id) {
  await sw_sub(id);
  fetchData();
}

useDraggable(el, subData, {
  animation: 150,
  handle: '.handle',
  async onUpdate() {
    await order_sub(subData.value)
  }
})
</script>

<style lang="css" scoped></style>
