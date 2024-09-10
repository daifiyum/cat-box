<template>
  <v-navigation-drawer v-model="drawer" location="right" temporary>
    <v-toolbar>
      <v-toolbar-title>设置</v-toolbar-title>

      <v-spacer></v-spacer>

      <v-btn icon="mdi-close" @click="drawer = !drawer"></v-btn>
    </v-toolbar>
    <v-container>
      <v-text-field
        v-model="updateDelay"
        label="更新延时"
        clearable
      ></v-text-field>
    </v-container>
    <template v-slot:append>
      <v-divider></v-divider>
      <div class="pa-4 flex-column">
        <v-btn
          color="primary"
          block
          prepend-icon="mdi-content-save-check-outline"
          variant="flat"
          @click="setUpdateDelay"
        >
          保存设置
        </v-btn>
      </div>
    </template>
  </v-navigation-drawer>
  <SnackBar ref="snackbarRef" />
</template>

<script setup>
import { ref, onMounted } from "vue";
import SnackBar from "../../../components/Snackbar.vue";
import { getDelay, setDelay } from "@/api/config";
import { storeToRefs } from "pinia";
import { useDrawerStore } from "@/stores/drawer";
const drawerStore = useDrawerStore();
const { drawer } = storeToRefs(drawerStore);
const snackbarRef = ref(null);

// 设置抽屉
const updateDelay = ref("");
onMounted(async () => {
  const { data } = await getDelay();
  updateDelay.value = data.update_delay;
});

async function setUpdateDelay() {
  await setDelay({ update_delay: updateDelay.value });
  showSnackbar("新设置已保存");
}

function showSnackbar(msg) {
  snackbarRef.value.openSnackbar(msg);
}
</script>

<style lang="scss" scoped></style>
