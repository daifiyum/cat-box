<template>
  <v-dialog v-model="dialog" max-width="450">
    <template v-slot:activator="{ props: activatorProps }">
      <v-btn
        variant="outlined"
        rounded
        color="primary"
        prepend-icon="mdi-pencil-outline"
        v-bind="activatorProps"
      >
        编辑
      </v-btn>
    </template>

    <v-card title="编辑">
      <v-card-text>
        <v-form v-model="valid">
          <v-text-field
            v-model="subName"
            :rules="inputRules"
            label="订阅名称"
          ></v-text-field>

          <v-text-field
            v-model="subLink"
            :rules="inputRules"
            label="订阅链接"
          ></v-text-field>
        </v-form>
        <v-switch
          v-model="autoUpdate"
          label="自动更新"
          color="primary"
          hide-details
          inset
        ></v-switch>
      </v-card-text>

      <v-card-actions>
        <v-spacer></v-spacer>

        <v-btn
          color="primary"
          text="保存"
          variant="tonal"
          :disabled="!valid"
          @click="rwOne(item.id)"
        ></v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
  <v-btn
    variant="outlined"
    rounded
    color="primary"
    prepend-icon="mdi-update"
    :loading="loading"
    @click.stop="upOne(item.id)"
  >
    更新
  </v-btn>
  <v-btn
    variant="outlined"
    rounded
    color="error"
    prepend-icon="mdi-trash-can-outline"
    @click.stop="rmOne(item.id)"
  >
    删除
  </v-btn>
</template>

<script setup>
import { ref, reactive, computed } from "vue";
import { rm_sub, rw_sub, up_sub } from "@/api/home";
import emitter from "@/utils/emitter";

const { item } = defineProps(["item"]);

let dialog = ref(false);
let valid = ref(true);
let inputRules = reactive([(value) => !!value]);
let subName = ref(item.name);
let subLink = ref(item.link);
let loading = ref(false);
let autoUpdate = ref(item.auto_update);

let autoUpdateToNum = computed(() => {
  if (autoUpdate.value) {
    return 1;
  } else {
    return 0;
  }
});

async function rmOne(id) {
  await rm_sub(id);
  emitter.emit("reloadData");
}

async function rwOne(id) {
  await rw_sub(id, {
    name: subName.value,
    link: subLink.value,
    auto_update: autoUpdateToNum.value,
  });
  emitter.emit("reloadData");
  dialog.value = false;
}

async function upOne(id) {
  loading.value = true;
  await up_sub(id);
  loading.value = false;
}
</script>

<style lang="css" scoped></style>
