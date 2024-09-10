<template>
  <v-form v-model="valid" class="d-flex align-center mt-3">
    <v-text-field
      variant="solo"
      v-model="subLink"
      :rules="inputRules"
      label="订阅链接"
      clearable
      hide-details
      required
    >
      <template v-slot:append-inner>
        <v-btn
          class="ml-3"
          color="primary"
          :disabled="!valid"
          :loading="loading"
          prepend-icon="mdi-file-download-outline"
          @click="addSubData"
        >
          导入
        </v-btn>
      </template>
    </v-text-field>
  </v-form>
</template>

<script setup>
import { ref, reactive } from "vue";
import { add_sub } from "@/api/home";
import emitter from "@/utils/emitter";

let valid = ref(true);
let subLink = ref("");
let inputRules = reactive([(value) => !!value]);
let loading = ref(false);

async function addSubData() {
  loading.value = true;
  valid.value = !valid.value;
  await add_sub({
    name: "订阅",
    link: subLink.value,
  });
  emitter.emit("reloadData");
  loading.value = false;
  valid.value = !valid.value;
}
</script>

<style lang="css" scoped></style>
