<template>
  <v-snackbar v-model="snackbarVisible" :timeout="2000" color="primary">
    {{ message }}
    <template v-slot:actions>
      <v-btn variant="text" @click="snackbarVisible = false" icon="mdi-close">
      </v-btn>
    </template>
  </v-snackbar>
</template>

<script setup>
import { ref } from "vue";
const snackbarVisible = ref(false);
const message = ref("");

function openSnackbar(msg) {
  message.value = msg;
  snackbarVisible.value = true;
}

defineExpose({
  openSnackbar,
});
</script>

<style lang="css" scoped></style>
