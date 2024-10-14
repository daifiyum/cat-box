import { defineStore } from "pinia";
import { ref } from "vue";

export const useSubStore = defineStore("sub", () => {
  const subNum = ref(0);
  return { subNum };
});
