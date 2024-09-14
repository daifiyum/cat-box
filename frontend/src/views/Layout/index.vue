<template>
  <LayoutDrawer />
  <LayoutHeader />

  <v-main>
    <v-container>
      <router-view />
    </v-container>
  </v-main>
</template>

<script setup>
import LayoutDrawer from "./components/LayoutDrawer.vue";
import LayoutHeader from "./components/LayoutHeader.vue";
</script>

<style lang="scss" scoped></style>
