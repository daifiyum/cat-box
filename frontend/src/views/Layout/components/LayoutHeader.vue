<template>
  <v-app-bar class="px-2">
    <template v-slot:prepend>
      <v-img
        src="@/assets/logo.svg"
        width="35"
        transition="scale-transition"
      ></v-img>
    </template>
    <v-app-bar-title>订阅</v-app-bar-title>
    <template v-slot:append>
      <v-btn
        variant="outlined"
        rounded
        color="primary"
        prepend-icon="mdi-vpn"
        class="mr-2"
        link
        :href="toProxyManager"
      >
        代理面板
      </v-btn>
      <v-btn icon="mdi-cog-outline" @click="drawer = !drawer"></v-btn>
    </template>
  </v-app-bar>
</template>

<script setup>
import { ref } from "vue";
import { storeToRefs } from "pinia";
import { useDrawerStore } from "@/stores/drawer";
const drawerStore = useDrawerStore();
const { drawer } = storeToRefs(drawerStore);

// 代理面板地址
const toProxyManager = ref(`http://${window.location.hostname}:9090/ui/`);
</script>

<style lang="css" scoped></style>
